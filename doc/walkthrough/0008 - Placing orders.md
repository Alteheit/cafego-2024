# Placing orders

2024-07-28

We have two large use cases remaining: we still need to allow users to add items to their cart, and we need to allow users to check out their cart. These are large when taken as a whole, but as usual, we will do our best to break down the requirements.

In this section, we'll focus on implementing the first feature for letting users add products to their cart. We have two objectives:

- Add a form to each product detail page that will allow the user to add products to their cart.
- Have the server listen to these form submissions and update the database appropriately.

## Start with the database

Whenever we have a database, we should usually start adding a new feature by thinking first whether we have the tables to support it. Let's do a quick inventory of what we have:

- User
- Session
- Product

It looks like we're short a few tables. I will assert that we need at least one new table here: "Cart Item." Each row in Cart Item will represent a line item stored in the user's cart, so for example: 2 Americanos in user Melinoe's cart.

Let's add the table definition to our `initDB` function in `database.go`. We only have to create a table here; we don't need to populate it with seed data. So there will be only one new line:

```go
queries := []string{
    "CREATE TABLE IF NOT EXISTS cgo_user (username TEXT, password TEXT)",
    "CREATE TABLE IF NOT EXISTS cgo_product (name TEXT, price INTEGER, description TEXT)",
    "CREATE TABLE IF NOT EXISTS cgo_session (token TEXT, user_id INTEGER)",
    // This one
    "CREATE TABLE IF NOT EXISTS cgo_cart_item (product_id INTEGER, quantity INTEGER, user_id INTEGER)",
}
```

If you reload your app, you should see when you enter your database that there is a new table `cgo_cart_item`.

> A slight tangent: if you, for any reason, need to delete your database to start over, you can simply delete the `db` file beside your `main.go` file. That is the entire SQLite database.

Now, it is time to add a form to the product detail page. Remember that a "form" is just an HTML element that allows users to send input back to the server. We've already made a form on the login page, so this should feel somewhat familiar.

```html
<h2>Add to cart</h2>
<form action="" method="post">
    <input type="hidden" name="product_id" value="{{ .Id }}">
    <label for="quantity">Quantity</label>
    <input type="number" name="quantity">
    <input type="submit" value="Add to cart">
</form>
```

As a refresher: the "action" of a form is the URL that the form will send the HTTP request to. If it is empty, the form will send the request to the same URL that the page is already on. The "method" of a form is the HTTP verb that it will use to send the HTTP request. Most of the time, that verb is POST.

Remember that we are trying to collect three pieces of data from our user:

- Their user ID.
- The ID of the product they want to add to cart.
- The quantity of the product they want to add to cart.

We know we can collect the user's ID from the session cookie that they will pass to the browser when they send the form. We have no way of collecting the quantity other than by having them enter it into the form.

What might look a bit strange is this "hidden" input we have at the top of the form. This is the standard way to include data in a form without having to ask the user to input something. We know that this form will only render on pages that are about a specific product, and we know that the product is accessible in our Go templates through the `product` object. We can simply attach the product's ID, which we know, to the form as a hidden field so that our server will have access to it later.

## The route handler

We can now write our route handler. We can start simple by collecting the three pieces of data we want and then echoing it back to the user.

I moved the entire existing body of `productHandler` into a GET branch and put this new code in a POST branch, but I will paste the entire thing here for your convenience.

```go
func productHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Get the product ID
		reqPath := r.URL.Path
		splitPath := strings.Split(reqPath, "/")
		elemCount := len(splitPath)
		// Do note that this will be a string.
		productId := splitPath[elemCount-1]
		// Need to convert from string to int
		intId, err := strconv.Atoi(productId)
		if err != nil {
			log.Fatal(err)
		}
		// Predeclare a product
		var product Product
		// Check each product for whether it matches the given ID
		for _, p := range getProducts() {
			if p.Id == intId {
				product = p
				break
			}
		}
		// If the for loop failed, then product will be the "zero-value" of the Product struct
		if product == (Product{}) {
			log.Fatal("Can't find product with that ID")
		}
		// Template rendering
		tmpl, err := template.ParseFiles("./templates/product.html")
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, product)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == "POST" {
		// Get user
		// This is copy pasted from indexHandler, so you might want to consider extracting this into its own function. I will keep this as is.
		cookies := r.Cookies()
		var sessionToken string
		for _, cookie := range cookies {
			if cookie.Name == "cafego_session" {
				sessionToken = cookie.Value
				break
			}
		}
		user := getUserFromSessionToken(sessionToken)
		userId := user.Id
		// Get product ID
		sProductId := r.FormValue("product_id")
		productId, err := strconv.Atoi(sProductId)
		if err != nil {
			log.Fatal(err)
		}
		// Get quantity
		sQuantity := r.FormValue("quantity")
		quantity, err := strconv.Atoi(sQuantity)
		if err != nil {
			log.Fatal(err)
		}
		// Echo values
		fmt.Fprintf(w, "User ID: %v; Product ID: %v; Quantity: %v", userId, productId, quantity)
	}
}
```

If you run this code, you should correctly see that the expected values are echoed back to you, in which case you may proceed.

## Back to the database

We will need a new database function `createCartItem` to put this data into our database.

```go
func createCartItem(userId int, productId int, quantity int) {
	q := "INSERT INTO cgo_cart_item (user_id, product_id, quantity) VALUES (?, ?, ?)"
	_, err := database.Exec(q, userId, productId, quantity)
	if err != nil {
		log.Fatal(err)
	}
}
```

Now, you can replace the last `fmt.Fprintf` line in our POST branch with:

```go
// Create a cart item
createCartItem(userId, productId, quantity)
http.Redirect(w, r, "/", http.StatusFound)
```

That's actually it. Add a few items to cart and check your database to see if it works.

## Checkpoint

Add a few items to cart. Submit a screenshot of what appears when you run `SELECT * FROM cgo_cart_item` in your SQLite shell.

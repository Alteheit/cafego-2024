# The cart page

2024-07-28

We are actually almost done. The next logical step here is to give the users a page on which they can view their cart. You will see that with all the setup we've done, this will actually be a rather straightforward addition to our app.

## Start with the database

Once again, we need to check if we have the tables to support our new feature. If we go into SQLite, we can see what tables we have:

```sqlite
sqlite> .tables
cgo_cart_item  cgo_product    cgo_session    cgo_user
```

We can use cgo_user and cgo_cart_item for this. We don't need to add a new table (or modify an existing table) for this feature.

## Add the template

Let's write a quick template that will display a user's cart items.

Remember that a Go template can be fed data from the route handler. In this case, let's say that we expect two pieces of data: first, the user, and second, an slice of Cart Items that belong to the user.

Note that the fields on the Go template's Cart Items do not necessarily have to correspond one-to-one with the fields in `cgo_cart_item`. It will become clear soon why this is useful.

```html
<h1>CafeGo</h1>

<p>Welcome, {{ .User.Username }}!</p>

<a href="/">Back to home page</a>

<h2>My Cart</h2>

<ul>
    {{ range $ci := .CartItems }}
        <li>{{ $ci.Quantity }} - {{ $ci.ProductName }}</li>
    {{ end }}
</ul>
```

We should also include a link on the index page that will take us to the cart page.

```html
<!-- In index.html. Place this inside the user block so it won't render if there's no user. -->
<a href="/cart">View my cart</a>
```

Great. There's still some more setup to do. Now that we will be handling Cart Items directly in our application code, we need to write a struct type for it in `database.go`.

```go
type CartItem struct {
	Id          int
	UserId      int
	ProductId   int
	Quantity    int
	ProductName string
}
```

You'll see that this doesn't exactly match the table structure. This is fine. The application type should serve the needs of the application. It will be the database layer's job to translate between the physical data, as represented in the database, and what the application needs.

While we're defining types, we know that we expect to pass a User and a slice of Cart Items to our cart page, so let's define that struct in `main.go` as well.

```go
type CartPageData struct {
	CartItems []CartItem
	User      User
}
```

If you're coming from the JavaScript version of this tutorial, you might notice that there's a much greater emphasis on setting up our data types in the Go version. This is the difference in development approach you must take in a statically typed language versus a dynamically typed language.

We can write our route handler now. I'll preemptively wrap the code in a GET branch, because it is evident that we will have to submit a form to this route, too.

```go
func cartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/cart.html")
		if err != nil {
			log.Fatal(err)
		}
		// Set to nil for now
		pageData := CartPageData{}
		tmpl.Execute(w, pageData)
	}
}
```

Remember to register your route handler:

```go
http.HandleFunc("/cart/", cartHandler)
```

Run your code. Your cart page _should_ render, but there won't be anything on it. We have yet to pass our template an actual User and slice of Cart Items.

Getting the user should be simple enough. We'll use the cookie method from the other handlers.

```go
cookies := r.Cookies()
var sessionToken string
for _, cookie := range cookies {
    if cookie.Name == "cafego_session" {
        sessionToken = cookie.Value
        break
    }
}
user := getUserFromSessionToken(sessionToken)
```

Getting the cart items is a little trickier. We can suggest the existence of a function to help us like this:

```go
cartItems := getCartItemsByUser(user)
```

Now, we need to implement that function. Note that now we run into a similar problem from before. The cart item table doesn't actually have the product name, so we need to look it up for each cart item. Do we do this in the query, or in the application logic? To remain consistent, I will do it in the query. You can do what you want, though I recommend doing it in the query.

```go
func getCartItemsByUser(user User) []CartItem {
	userId := user.Id
	q := `
	SELECT
		cgo_cart_item.rowid,
		cgo_cart_item.user_id,
		cgo_cart_item.product_id,
		cgo_cart_item.quantity,
		cgo_product.name
	FROM cgo_cart_item
	LEFT JOIN cgo_product ON cgo_cart_item.product_id = cgo_product.rowid
	WHERE cgo_cart_item.user_id = ?
	`
	rows, err := database.Query(q, userId)
	if err == sql.ErrNoRows {
		return []CartItem{}
	} else if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var result []CartItem
	for rows.Next() {
		var cartItem CartItem
		rows.Scan(&cartItem.Id, &cartItem.UserId, &cartItem.ProductId, &cartItem.Quantity, &cartItem.ProductName)
		result = append(result, cartItem)
	}
	return result
}
```

That's actually it. We can proceed to the checkpoint.

## Checkpoint

Add some Cart Items to your cart if you have not done so already. Submit a screenshot of your cart page with a few Cart Items in it.

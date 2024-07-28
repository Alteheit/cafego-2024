# Product details

2024-07-27

CafeGo is progressing nicely. However, we have a long way to go. Before we implement the exciting features like users/login and adding products to a cart, we need to at least let CafeGo function as a brochure website. As it is, CafeGo doesn't even display the details of products.

This section will focus on adding one more route to your app. If a user visits the route `/product/{productId}`, they should be shown a page specifically for that product.

## Adding details to your products

My personal approach to building web apps, no matter what framework I use, is to always start with the data. Let's pretend that we have one more field for each entry in our product catalog: a description. Let us also add a numeric ID to each product.

```go
// In database.go

type Product struct {
	Id          int
	Name        string
	Price       int
	Description string
}

func getProducts() []Product {
	return []Product{
		{Id: 1, Name: "Americano", Price: 100, Description: "Espresso, diluted for a lighter experience"},
		{Id: 2, Name: "Cappuccino", Price: 110, Description: "Espresso with steamed milk"},
		{Id: 3, Name: "Espresso", Price: 90, Description: "A strong shot of coffee"},
	}
}
```

Run your server and visit the home page. You should see that nothing has changed. Even though our template has access to all the product data, it only uses two fields: the name and the price.

We will keep our index page as is with respect to how much data it shows. However, each product should have its own page where we can see all its data.

## Adding a route

Return to `main.go`. We will need to add a new route handler for the path `/product/{productId}`.

Here's a slight problem. Other frameworks have a way to "capture" the variable part of a path, which would be `{productId}` in our case. But we're using Go, so we have to capture this variable ourselves. Oh well.

Let's write another route handler. We will expect this handler to run when users visit a URL similar to `/product/{productId}`. For now, all we will make the handler do is return the product ID that the user gives us.

```go
func productHandler(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	splitPath := strings.Split(reqPath, "/")
	elemCount := len(splitPath)
	productId := splitPath[elemCount-1]
	fmt.Fprint(w, productId)
}
```

Register the route handler. This one's a bit tricky, because Go has its own rules for how the HTTP server will _match_ requests to patterns. After tweaking the code a little, I found that it works if you register the handler to `/product/`, but not `/product`, so keep this in mind.

```go
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/product/", productHandler)
	http.ListenAndServe(":5000", nil)
}
```

Visit `http://localhost:5000/product/2` now. It should return `2` to you. If it does, you may proceed.

## Making the template

Think back to our intent for this new route. When a user visits this route, we want to show them everything about the product whose ID they passed to the server.

This means two things:

- We will need to render a new web page. This page will expect to receive the details of only _one_ product.
- We will need to find a way to retrieve the product details from only the product ID.

Let's tackle the first one in this subsection. We'll need to create a new template, `product.html`.

```html
<h1>CafeGo</h1>

<a href="/">All products</a>

<h2>{{ .Name }}</h2>
<p>Price: PHP {{ .Price }}</p>
<p>{{ .Description }}</p>
```

These fields should look familiar. We are going to pass an instance of a Product struct to our template.

We can easily extend the template rendering code from `indexHandler` to fit the needs of `productHandler`. The new part is getting the product we need from the Product slice. Well, since this is Go, we can just run a for loop over it and stop when you've found the Product you need. Once you accept that there really isn't a better way to do it than the dumb, obvious way, the better you'll be at Go.

```go
func productHandler(w http.ResponseWriter, r *http.Request) {
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
}
```

Yup, welcome to Go. I still think this is better than NodeJS. Do not mistake the quality of the experience of writing the code for the quality of the codebase that you produced.

## Add links to the index page

We still have to do one thing: add links to the product details from our index page. To do this, we'll need to know the ID of each product, since the ID goes at the end of each product link. Thankfully, we have easy access to it with `$p.Id` in our loop.

Replace each list item with this:

```html
<li><a href="/product/{{ $p.Id }}">{{ $p.Name }}</a> - PHP {{ $p.Price }}</li>
```

An ugly solution, and I'm not sure there's a better approach. But that's the Go way.

Reload your app. It should now at least function as a brochure website.

## Checkpoint

Submit two screenshots. The first screenshot should be of your home page. The second screenshot should be of your Macchiato page.

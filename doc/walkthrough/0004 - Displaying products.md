# Displaying products

2024-07-27

Our online shopping website wouldn't be very good without any products. In this section, we will add some. Our only objective for this section is to display a list of products on the index page, but there's a lot going on here that needs some discussion.

## Hardcoding your products

Given what we have so far, the simplest way to display a list of products is to simply hardcode the product details in the app's source code and then pass the product data to a template file.

Let's start by writing out our product list. We'll keep it simple. Add this code just above your index handler function:

```go
type Product struct {
	Name        string
	Price       int
}

func getProducts() []Product {
	return []Product{
		{Name: "Americano", Price: 100},
		{Name: "Cappuccino", Price: 110},
		{Name: "Espresso", Price: 90},
	}
}
```

To review what we just did:

- We defined the structure of what a "Product" record should look like by using a Go struct. We say here that a Product should have two properties: a name and a price.
- We wrote a function whose only job it is to give us back a _slice_, which is something like a list or an array, of Product records. The syntax to represent the type of a "slice of Something" is `[]Something`.
- Inside the function, we instantiate a Product slice, and inside, we store three instances of Product.

Having to define the shape of data upfront is another of Go's quirks. Well, if we're being honest, it isn't that Go is quirky, it's that we may be unfamiliar with _static typing_. Many other respectable languages like Java, Rust, and Haskell pay very close attention to the type of their data. Languages that are popular nowadays like Python and JavaScript pay much less attention to the types.

We will expect to pass this slice of Products into our index page somehow, so let's add a `Products` field to our `IndexPage` struct.

```go
type IndexPageData struct {
	Username string
	Products []Product
}
```

It should be reasonably clear what should come next. We need to retrieve our products from within our index handler and add it to our page data.

```go
sampleProducts := getProducts()
// Where ... represents the other fields that were added
samplePageData := IndexPageData{..., Products: sampleProducts}
```

Great. We are now passing the product data into our template. If we load our page, though, there will be nothing displayed, because our template doesn't actually make use of the products yet.

We can add an unordered list to our template, and then we can loop within the unordered list to add list items for each product.

```html
<ul>
    {{ range $p := .Products }}
        <li>{{ $p.Name }} - PHP {{ $p.Price }}</li>
    {{ end }}
</ul>
```

You will see shortly that this non-HTML syntax is the Go-specific template syntax for a loop. It's slightly different from a regular Go loop (for example, the $ is required here). The semantics of the construct should be very clear, though.

That's actually it. If you reload your app, you should be able to see the products rendered on the page.

## Cleanup

We probably shouldn't be storing the data in this file. It's best to keep concerns separate in the code.

Create a new file, `database.go`. We will move our data here.

```go
// database.go
package main

type Product struct {
	Name        string
	Price       int
	Description string
}

func getProducts() []Product {
	return []Product{
		{Name: "Americano", Price: 100},
		{Name: "Cappuccino", Price: 110},
		{Name: "Espresso", Price: 90},
	}
}
```

Delete `Product` and `getProducts` from `main.go`. You should see that the code will actually still work as intended. This is because we have not gone so far as to make a new _package._ Everything is still in the main package, so the symbols in `database.go` are still available to code in `main.go`.

## Checkpoint

Add a new product to the products slice, "Macchiato", with price 120. Submit a screenshot of your index page with the new product added.

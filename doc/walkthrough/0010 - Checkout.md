# Checkout

2024-07-28

Welcome to the home stretch of CafeGo. Every online shopper knows that a cart must be checked out for you to receive your goods.

## Start with the database

We have settled into a good rhythm. In the previous sections, we started with the database. We then go to our template. Once the structures at the two ends have been established, we ping pong between them to build up the middle.

Nothing about that will change here. Do we have the tables we need to implement transactions?

```sql
sqlite> .tables
cgo_cart_item  cgo_product    cgo_session    cgo_user
```

Not quite. `cgo_cart_item` will obviously play a critical role here, but there are two missing tables. The first is a table for a Transaction, and the second is a table for Line Items. A Transaction is the "header item" for one or more Line Items.

We will need to add some new tables to our database initialization function.

```go
queries := []string{
    "CREATE TABLE IF NOT EXISTS cgo_user (username TEXT, password TEXT)",
    "CREATE TABLE IF NOT EXISTS cgo_product (name TEXT, price INTEGER, description TEXT)",
    "CREATE TABLE IF NOT EXISTS cgo_session (token TEXT, user_id INTEGER)",
    "CREATE TABLE IF NOT EXISTS cgo_cart_item (product_id INTEGER, quantity INTEGER, user_id INTEGER)",
    // These two
    "CREATE TABLE IF NOT EXISTS cgo_transaction (user_id INTEGER, created_at TEXT)",
    "CREATE TABLE IF NOT EXISTS cgo_line_item (transaction_id INTEGER, product_id INTEGER, quantity INTEGER)",
}
```

On our `cgo_transaction` table, we have an interesting field `created_at`. This is a timestamp meant to hold the date and time of a transaction's creation. Other databases have dedicated timestamp modules to handle this, but SQLite does not, so we store it as text. [It's in their docs.](https://www.sqlite.org/datatype3.html)

## Add the template

Well, we already have the template. We will initiate checkout via our cart page. The interesting part is that we now need to add a form for the user to send the instruction to check out.

What sort of information would we need to send to the server other than the user's intent? Honestly, for this one, we only need to tell the server who's doing the checkout. We can either pass the user's ID in as a hidden field on the form or retrieve it from the session cookie. Since we already have a good "framework" for getting our user from the cookie, let's get them from the cookie. We thus don't need to include the user ID in any way here.

```html
<form action="" method="post">
    <input type="submit" value="Checkout">
</form>
```

## The route handler

This will issue a post request to `/cart/`, so let's add a branch for that in the cart handler function.

```go
// In a POST block

cookies := r.Cookies()
var sessionToken string
for _, cookie := range cookies {
    if cookie.Name == "cafego_session" {
        sessionToken = cookie.Value
        break
    }
}
user := getUserFromSessionToken(sessionToken)
// The rest of the function...
```

You can add a sanity check here, but we've done this enough to not really need one.

## Back to the database

Our entire intent with this section is to transform a set of Cart Items into a Transaction and its associated set of Line Items. If we were to do this from our route handler, we might expect the call to the database to look like this:

```go
checkoutItemsForUser(user)
```

We'll have to go back to our database file to implement the function. This one will take a little bit of brain juice, but it's entirely doable.

```go
func checkoutItemsForUser(user User) {
	// Fetch cart items first
	// We want to transform each of these into a line item
	cartItems := getCartItemsByUser(user)
	// Create a new transaction
	now := time.Now().UTC()
	q := "INSERT INTO cgo_transaction (user_id, created_at) VALUES (?, ?)"
	// We need the first return value for once for the ID of the new transaction
	res, err := database.Exec(q, user.Id, now)
	if err != nil {
		log.Fatal(err)
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	// Transform each cart item into a line item
	for _, ci := range cartItems {
		var q string
		q = "INSERT INTO cgo_line_item (transaction_id, product_id, quantity) VALUES (?, ?, ?)"
		_, err = database.Exec(q, lastInsertId, ci.ProductId, ci.Quantity)
		if err != nil {
			log.Fatal(err)
		}
		q = "DELETE FROM cgo_cart_item WHERE rowid = ?"
		_, err = database.Exec(q, ci.Id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
```

To round things off, redirect the user in the route handler once this is all done.

```go
http.Redirect(w, r, "/", http.StatusFound)
```

Congratulations. Run your app and check out your cart. You won't get a visual signal, but your database should have changed.

## Checkpoint

Take a very short video of you going through the whole checkout process. Show me especially how the cart page changes before and after you checkout. Then show me the database, specifically the output of `SELECT * FROM cgo_line_item;`.

## A challenge

Your clients want you to implement a "transaction history" feature. Go do it.

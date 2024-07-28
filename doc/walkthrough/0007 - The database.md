# The database

2024-07-28

We have survived this far with our "fake" database. However, it should have become clear in the previous section that our fake database will not work for much longer. It is time to explore our options for a database.

This is likely to be one of the longest and most difficult subsections. You will see why I think that Django is a much better framework to teach beginners. However, we will not back down from this. Even if it takes a while, mastering the concepts in this subsection is the key to mastering Go and web development as a whole.

## Our options

A "database" is a program whose purpose it is to manage access to data that is stored persistently on the computer's disk. The _persistence of data_ is the key feature of a database. This is what we currently do not have with our `database.go` setup.

It is entirely possible to persist your data by writing to simple files on disk. This is actually a valid strategy for some small-scale applications. However, most software professionals (correctly) consider it best practice to use a real database management system, especially if the data under management is sufficiently voluminous or complex.

### Relational databases

A "relational" database is one that stores data in terms of tables. This is, by far, the most common form of database.

Each table in a relational database is about one _entity_. Each row in a table is about an _instance_ of the entity. In our app so far, `products` (or `product` according to the naming conventions of most relational database administrators) would be a table. The `"Americano"` record would be a row in the `product` table.

Tables in relational databases also tend to have references to rows in other tables. These references are called "joining keys." In our app so far, the `session` table would have a `user_id` column as a joining key to the `user` table.

The most popular relational database system nowadays is PostgreSQL. There are a number of other relational databases on the market like MySQL and MariaDB, or if you're in the enterprise market, Oracle and Microsoft SQL Server. These are all rather heavyweight systems, though. For smaller use cases, a lightweight relational database called SQLite is sometimes all you need.

The language used for interacting with relational databases is called "Structured Query Language," or SQL. Compared to most modern languages, SQL is frankly unpleasant to write and read, so a lot of new developers shy away from relational databases because of SQL. If you had a `line_item` table and wanted to sum the sales per product, your SQL would look like this:

```sql
SELECT
    product,
    SUM(price * quantity) AS sales
FROM line_item
GROUP BY product
ORDER BY sales DESC;
```

It isn't great. SQL was designed many decades ago to work with the relational data model. This data model unfortunately also conflicts with the way most modern languages store data, which is usually in the form of "objects." There have been attempts to bridge the two data models with libraries known as "object-relational mappers," or ORMs, but this space is notoriously difficult to get right. It is still best to understand how SQL works.

If you are going to have a career anywhere in the vicinity of data, there is no escaping SQL. I won't excuse its flaws as a language, but take it from me that you'll eventually get used to it and come to appreciate its power regardless of its warts.

### "NoSQL" databases

Any database that does not follow the relational data model (i.e., the table model) is perhaps unfairly grouped into what is known as "NoSQL." The two most common non-relational models are the key-value model, which is very similar to using a Python dictionary or a JavaScript object, and the document model, which is very similar to using _nested_ dictionaries/objects.

Most data models have use cases for which they are appropriate. However, I caution against engaging deeply with the NoSQL world as a beginner if only because there was previously, and perhaps still is, a lot of marketing hype around using non-relational models for _everything_. Most use cases are handled perfectly well by the relational model. That's why the relational model has survived so long. I wrote an entire article about it [here](https://joeilagan.com/article/relational-data-ite).

### Our database

For the purposes of this walkthrough, we will choose SQLite. It is a _relational_ database, which by now should be obvious is my preference. It is also much easier to install than a heavier relational database like Postgres, so we can get started almost right away.

This is the first time we'll need to install something other than the basic Go toolchain. Please run this command to download a database driver for Go:

```zsh
go get "github.com/mattn/go-sqlite3"
```

This is a "database driver." Go, as expected, has a built-in interface for working with SQL, but Go doesn't come equipped with specific code for every flavor of SQL database, so we need to download this separately.

We'll need to set up our database. To keep our design decision, we should only have to change things in `database.go`. Ideally, we shouldn't have to change the way our route handlers interact with our database at all.

First, we need to initialize and connect to our database every time we run our app. We'll store our initialization instructions in a function called `initDB`.

```go
// NOTE: You may need to import go-sqlite3 explicitly.
// Go won't know how to import this on its own
import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func initDB() {
	db, err := sql.Open("sqlite3", "./db")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
    database = db
}
```

Let's review what we just wrote.

- Up top, we say that there should be a pointer to e `sql.DB` handle available throughout the package. Of course, at the start, this will be uninitialized.
- We first open a "handle" on our database with `sql.Open`. The first argument to `sql.Open` is our database type, which is SQLite in this case. The second argument to `sql.Open` is the set of connection details to our database server. Since we're using SQLite, this is just a file.
- We then "ping" the database to make sure that it's up and running. If something goes wrong trying to connect to it, it will error, which we can then handle appropriately. (In this app, we will just crash if we find an error.)
- We then assign the top-level `database` variable to be the `db` variable we created in this function.

Then, in your `main` function in `main.go`, call `initDB` before anything else.

```go
func main() {
	initDB()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/product/", productHandler)
	http.HandleFunc("/login/", loginHandler)
	http.ListenAndServe(":5000", nil)
}
```

The next time you run your app, it should generate a `db` file in `cafegoroot`.

## Seed data

So far so good. The next thing on our database agenda is to not only connect to the database but also to set up the required tables. Let's extend `initDB` to do just that.

```go
queries := []string{
    "CREATE TABLE IF NOT EXISTS cgo_user (username TEXT, password TEXT)",
    "CREATE TABLE IF NOT EXISTS cgo_product (name TEXT, price INTEGER, description TEXT)",
    "CREATE TABLE IF NOT EXISTS cgo_session (token TEXT, user_id INTEGER)",
}
for _, q := range queries {
    _, err := db.Exec(q)
    if err != nil {
        log.Fatal(err)
    }
}
```

The strings in the `queries` slice are SQL. We simply iterate over a set of commands, and for each, we execute it. We don't expect to receive any data back, so we use `db.Exec` and not `db.Query`, which we will use later. We also prefix each of our tables with `cgo_` because if we were to use the base names of our entities, like `user` and eventually `transaction`, we would run into name conflicts with reserved SQL keywords.

This is pretty straightforward so far. The next part of our initialization code should insert seed data into our tables, but only if the tables don't already have any data in them.

Let's clean up a little before we insert seed data. Extract the hardcoded users and products from their functions and put them in top-level variables.

```go
var seedUsers = []User{
	{
		Id:       1,
		Username: "zagreus",
		Password: "cerberus",
	},
	{
		Id:       2,
		Username: "melinoe",
		Password: "b4d3ec1",
	},
}
var seedProducts = []Product{
	{Id: 1, Name: "Americano", Price: 100, Description: "Espresso, diluted for a lighter experience"},
	{Id: 2, Name: "Cappuccino", Price: 110, Description: "Espresso with steamed milk"},
	{Id: 3, Name: "Espresso", Price: 90, Description: "A strong shot of coffee"},
    {Id: 4, Name: "Macchiato", Price: 120, Description: "Espresso with a small amount of milk"},
}
```

Once we've done that, we can proceed to using the seed data in our database initialization.

```go
// Seed data
var q string
var count int
// cgo_user
q = "SELECT COUNT(*) FROM cgo_user"
err = db.QueryRow(q).Scan(&count)
if err != nil {
    log.Fatal(err)
}
if count == 0 {
    q = "INSERT INTO cgo_user (username, password) VALUES (?, ?)"
    for _, u := range seedUsers {
        _, err = db.Exec(q, u.Username, u.Password)
        if err != nil {
            log.Fatal(err)
        }
    }
}
// cgo_product
q = "SELECT COUNT(*) FROM cgo_product"
err = db.QueryRow(q).Scan(&count)
if err != nil {
    log.Fatal(err)
}
if count == 0 {
    q = "INSERT INTO cgo_product (name, price, description) VALUES (?, ?, ?)"
    for _, p := range seedProducts {
        _, err = db.Exec(q, p.Name, p.Price, p.Description)
        if err != nil {
            log.Fatal(err)
        }
    }
}
```

Run your code a few times, and you should see that it works just fine. If you go into your SQLite shell with `sqlite3 db`, and if you `SELECT * FROM cgo_user` and `SELECT * FROM cgo_product`, you'll see that they are populated the the data, as expected.

Let's review some of the potentially confusing items here.

- `db.QueryRow(q).Scan(&count)` might look weird, but if you recall what we've seen so far, this is really just putting the result of a query into the address of an existing variable. Good thing we learned about pointers, right?
- The `(?, ?, ..., ?)` syntax here is how SQLite accepts variable input. Notice that we have to provide a list of arguments to the `db.Exec` call to fill these placeholders. We should not simply `fmt.Sprint` a string here, because that's how you make yourself vulnerable to an attack called "SQL injection," which is where attackers manipulate the SQL itself to do _anything_ with your database.

The code we have so far is getting quite long, but believe me: I just finished writing the JavaScript version of this tutorial a few days ago, and the code is nightmarish. Go handles this a lot better.

## Refactoring users and products

We can now refactor our `getUsers` and `getProducts` functions to use the database instead of some hardcoded data. Let's start with getUsers.

```go
func getUsers() []User {
	var result []User
	q := "SELECT rowid, username, password FROM cgo_user"
	rows, err := database.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, user)
	}
	return result
}
```

To break it down:

- We basically want to transform a database query result into a slice of Users. We start by writing out the query. In SQL, the order that the columns are written in defines the order in which they will be displayed.
- We run the query. After the error check, we have to _defer_ closing the result. The `defer` keyword makes the statement run after the function ends. Closing the result is necessary because connecting to the database isn't free.
- We go through each of the rows, and for each row, we _scan_ the results (which we expect to be ID, username, and password in that order) into the fields of an empty User struct.

We can use a very similar pattern for `getProducts`.

```go
func getProducts() []Product {
	var result []Product
	q := "SELECT rowid, name, price, description FROM cgo_product"
	rows, err := database.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Description)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, product)
	}
	return result
}
```

Run your app. You should see that it actually behaves exactly the same.

## Refactoring sessions

We have two more functions to refactor: `setSession` and `getUserFromSessionToken`. (We will get to `getSessions` eventually.)

`setSession` is relatively straightforward, since we already know how to insert data.

```go
func setSession(token string, user User) {
	q := "INSERT INTO cgo_session (token, user_id) VALUES (?, ?)"
	_, err := database.Exec(q, token, user.Id)
	if err != nil {
		log.Fatal(err)
	}
}
```

It's `getUserFromSessionToken` that might be more difficult. There are two general approaches that make sense: keep the database query simple and do the logic in the application, or do the logic in the database query and keep the application code simple. My personal opinion is that we should aim to do complicated queries as close to the database as possible, since databases tend to be optimized for such things.

```go
func getUserFromSessionToken(token string) User {
	q := `
	SELECT
		cgo_session.user_id,
		cgo_user.username,
		cgo_user.password
	FROM cgo_session
	INNER JOIN cgo_user
	ON cgo_session.user_id = cgo_user.rowid
	WHERE cgo_session.token = ?
	LIMIT 1;
	`
	var user User
	err := database.QueryRow(q, token).Scan(&user.Id, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return User{}
	} else if err != nil {
		log.Fatal(err)
	}
	return user
}
```

You _can_ use the first approach, but I have it on good authority queries like this are really the database's job.

## Checkpoint

Submit one screenshot of your SQLite shell after running `SELECT * FROM cgo_product`.

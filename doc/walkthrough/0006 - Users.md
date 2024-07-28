# Users

2024-07-28

This will be a long section. Authentication -- knowing who someone is -- is a critical part of most non-trivial web apps. Most web frameworks like Rails and Django have a way to handle this. Since we're not using a framework at all in Go, we don't have a built-in way to handle users.

Thus it is part of our task in this walkthrough to explore the fundamentals of how a browser and a server can work together to remember who a user is. This will go deeper in the weeds than I originally planned, but you will emerge from these weeds as a stronger programmer.

## Cookies

A bit of background first.

Browsers and web servers talk to each other using the HyperText Transfer Protocol (HTTP). HTTP is designed to be "stateless," which just means that the server does not need to remember anything about what past HTTP requests looked like or contained.

The browser, however, _can_ remember what it was sent via a web server. There are a number of ways a browser can remember data, but for our purposes, we will use the "cookie."

A [cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies) is a small piece of data a server sends to a browser. In an HTTP response, a server may instruct the browser to store a cookie with this header:

```http
Set-Cookie: <cookie-name>=<cookie-value>
```

When a browser sends another HTTP request to the same website, it usually includes all its cookies for that website by default in the `Cookie` header as such:

```http
Cookie: yummy_cookie=choco; tasty_cookie=strawberry
```

You will typically not deal with the headers directly like this. Most web frameworks will expose cookies to you in interfaces that are idiomatic to the programming language they are written in. Underneath those abstractions, though, this is what is happening.

One of the most common use cases for a cookie, then, is to have the user's browser store a tiny token that identifies the user. In most web apps, this token will not represent something static and well-known. It will not be a username or an email. It will instead be a meaningless series of random bytes that is associated with the user only in the server's database and is useless in any other context. This is called a _session_.

```http
Set-Cookie: my_web_app_session=98a51b998f5ec044cfd5f6a2bf5fd2bb
```

The rest of this section will focus on getting our Go app to the point where we can set cookies to remember users via their sessions.

## Creating some users

To set up our first users, we will follow the same hacky approach we used for our products. Go to `database.go` and create the following struct type:

```go
type User struct {
	Id       int
	Username string
	Password string
}
```

Let's also create a function whose sole purpose is to return some users.

```go
func getUsers() []User {
	return []User{
		{
			Id:       1,
			Username: "zagreus",
			Password: "cerberus",
		},
		{
			Id: 2,
			Username: "melinoe",
			Password: "b4d3ec1",
		}
	}
}
```

I should note that it is a terrible idea to store passwords like this. We will tolerate it for the purposes of the walkthrough, but passwords should _always_ be hashed (not encrypted, hashed! they are not the same) and stored securely in a database.

Anyway, let's proceed. For our users to be able to log in, they will need a login page. We will make that next.

## A login page

This part should be easy enough to understand. To create a login page, we will just create a new Go template file that has a login form.

```html
<h1>CafeGo</h1>

<h2>Login</h2>

<form action="" method="post">
    <label for="username">Username</label>
    <input type="text" name="username">
    <label for="password">Password</label>
    <input type="password" name="password">
    <input type="submit" value="Login">
</form>
```

To render this page, we will write another route handler in `main.go`.

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/login.html")
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Register your login handler in main()
http.HandleFunc("/login/", loginHandler)
```

Let's do a brief review of what we've just written.

A "form" is an HTML element that collects data from a user and sends it to the server. In this case, we collect two pieces of data from the user: their username and their password.

Forms typically use a different HTTP "verb" called POST. Most HTTP requests that browsers send are GET requests, which are meant to fetch data. POST requests are different in that they are also meant to _send_ data to a browser.

This particular form thus collects a user's username and password then sends a POST request to the same route (i.e., `/login`). We can write another route handler for `/login` that handles POST requests specifically.

Let's do a sanity check. Our route handler for POST to `/login` will simply echo the user's username and password.

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/login.html")
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == "POST" {
		rUsername := r.FormValue("username")
		rPassword := r.FormValue("password")
		fmt.Fprintf(w, "Username: %v; Password: %v", rUsername, rPassword)
	}
}
```

Now, reload the app and try to log in. The app should simply echo what you entered back to your browser.

## Setting a cookie

We now have a working login form. Our ultimate intent is to use the details from the login form to create a session, tie it to a user, and store the session in a cookie. However, before we do that, we should check whether we can set a cookie at all.

We can test this by simply setting a cookie `cafego_username` to the username that the user submitted via the login form. Do not do this in a serious app, of course.

Change the POST branch of your login handler to this:

```go
rUsername := r.FormValue("username")
cookie := http.Cookie{Name: "cafego_username", Value: rUsername, Path: "/"}
http.SetCookie(w, &cookie)
// Then, redirect user to the home page
http.Redirect(w, r, "/", http.StatusFound)
```

Clear enough. There are two potential footguns here:

- You explicit need to set the `Path` property of an `http.Cookie` to "/" if you want the cookie to be retrievable on all paths. When writing this for the first time, I was confused why the cookie was not showing up on the index page even though I had set it correctly on the login page. This is why.
- `http.SetCookie` needs the _address_ of an `http.Cookie`, not the cookie itself. We will go over what this means in the next subsection here.

### A tangent on pointers

We should go over what this `&` symbol means.

Many languages, including Go, differentiate between _values_ and _references_ to those values. By default, most things in Go are values.

```go
myInt := 1
```

Sometimes, though, you don't need the value itself, but the _memory address_ of the value.

```go
myIntAddress := &myInt
```

The `&` symbol in Go asks for the address of the underlying value in memory. It very much does not store the value of the data itself.

Since you have a reference to the underlying data, you can _de-reference_ it with the `*` operator to use the underlying data.

```go
*myIntAddress := 2
```

Now, `myInt` will contain the value 2. You changed its value through its address.

I saw this explanation on Reddit some time ago. I hope they don't take this image down, because it was very helpful.

![An image of a meme explaning the difference between the ampersand and the asterisk when working with pointers.](https://external-preview.redd.it/FRCv4nqtau6Hpk-GRiA1UOjvn9JGn-ueImzf3O1oUbo.jpg?auto=webp&s=93893cd51c79e629cbb627ec25c476032bd52bce)

#### When does this matter?

In a lot of Go programs, pointers pop up sometimes, but not all the time. You will occasionally see things like how `http.SetCookie` takes a `*http.Cookie`, not an `http.Cookie`. This means that the function is asking for a _pointer_ to the data, not the data itself. That is why we need to use the `&` operator to fetch the address of our cookie.

This also explains what `r` is in our route handlers. `r` is actually a _pointer_ to an `http.Request`, not the `http.Request` itself. We have not had to de-reference the pointer to `r` so far because Go usually automatically de-references things for us.

Pointers are a huge part of lower-level (i.e., closer to the machine) languages like C, but higher-level languages usually try to limit their use. Overusing references tends to tie or complect things together. It is best to stick to bare values where possible.

## Reading the cookie

Remember our index page? We have a small snippet that displays a username. Let's swap it out from the dummy username to use the username stored in the cookie.

In `templates/index.html`, add the following code to only render the username snippet if it isn't a zero-value (i.e., "" or nil). While we're at it, we may as well add a link to login.

```html
{{ if .Username }}
    <p>Welcome, {{ .Username }}!</p>
{{ end }}
<a href="/login">Login</a>
```

Change the index handler to read the cookie. This one will be a bit verbose, because we'll need to actually go through the cookies slice to find what we need. _Go figure..._

```go
// Remember that if sampleUsername is not shadowed, it will be the zero-value of a string, "".
var sampleUsername string
for _, cookie := range cookies {
    if cookie.Name == "cafego_username" {
        sampleUsername = cookie.Value
        break
    }
}
```

If you log in now, you should be able to see the username you input.

That's not a very secure login. We can somewhat trivially add a guard to our login handler to check if the user's credentials at least match any that are in our database.

```go
// In the POST arm of `loginHandler`
rUsername := r.FormValue("username")
rPassword := r.FormValue("password")
var user User
for _, u := range getUsers() {
    if (rUsername == u.Username) && (rPassword == u.Password) {
        user = u
    }
}
if user == (User{}) {
    fmt.Fprint(w, "Invalid login. Please go back and try again.")
    return
}
cookie := http.Cookie{Name: "cafego_username", Value: rUsername, Path: "/"}
http.SetCookie(w, &cookie)
http.Redirect(w, r, "/", http.StatusFound)
```

It's not great, but it's something.

## Sessions

We've established that we can set and get cookies, but our cookie is currently the username of a user. This is horribly insecure. If we base our security decisions on whether we can find a username cookie, an attacker can simply set the cookie to someone's username to act as their account.

In this subsection, we'll try to refactor our web app to use sessions instead of usernames as cookie data. Sessions are supposed to be meaningless, easy to control, easy to revoke, and difficult to guess, so they are a much more secure way to identify users across requests.

In our database file, we can create a new type to represent sessions.

```go
type Session struct {
	Token  string
	UserId int
}
```

Now comes something different. We can define a package-level slice of Session objects that our other functions will be able to manipulate when needed.

```go
// Starts empty
var sessions = []Session{}
```

I should note that this is a very poor practice in a real app, especially if many threads are expected to read from and write to the data at the same time, but it will serve our needs for now. We will transfer everything to a real database in the next section, anyway.

To keep with the design principles so far, let's write a function that just returns the sessions to us. We will have to be disciplined and refrain from interacting with the underlying session data except through these functions.

```go
func getSessions() []Session {
	return sessions
}
```

We'll have to be careful with this one, because we're returning a slice, which is not "copied" in the way that you'd expect. If you get a slice from this function, if you change the slice from the calling function, the underlying slice will change.

This isn't a Go thing, it's just a programming thing in general. You have to be conscious of value types and reference types.

We know our intent is to associate sessions with user IDs, so let's write two more functions. The first function is to set a new session, and the second function is to get a user from a session token.

```go
// In database.go

func setSession(token string, user User) {
	sessions = append(sessions, Session{Token: token, UserId: user.Id})
}

func getUserFromSessionToken(token string) User {
	var userId int
	for _, session := range sessions {
		if session.Token == token {
			userId = session.UserId
		}
	}
	for _, u := range getUsers() {
		if u.Id == userId {
			return u
		}
	}
	return User{}
}
```

We can now rewrite sections of our login handler and our index handler to use sessions.

```go
// In main.go

// Here's how we'll generate a random session token.
func generateSessionToken() string {
	rawBytes := make([]byte, 16)
	_, err := rand.Read(rawBytes)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(rawBytes)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	cookies := r.Cookies()
	var sessionToken string
	for _, cookie := range cookies {
		if cookie.Name == "cafego_session" {
			sessionToken = cookie.Value
			break
		}
	}
	user := getUserFromSessionToken(sessionToken)
	sampleProducts := getProducts()
	samplePageData := IndexPageData{Username: user.Username, Products: sampleProducts}
	err = tmpl.Execute(w, samplePageData)
	if err != nil {
		log.Fatal(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("./templates/login.html")
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == "POST" {
		rUsername := r.FormValue("username")
		rPassword := r.FormValue("password")
		var user User
		for _, u := range getUsers() {
			if (rUsername == u.Username) && (rPassword == u.Password) {
				user = u
			}
		}
		if user == (User{}) {
			fmt.Fprint(w, "Invalid login. Please go back and try again.")
			return
		}
		// Set a session instead of a username
		token := generateSessionToken()
		setSession(token, user)
		cookie := http.Cookie{Name: "cafego_session", Value: token, Path: "/"}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
```

Now, if you log in with one of the users in your database, it should redirect you to the index page and greet you appropriately.

Please also note that since we do not have a _persistent_ database, every time you reload your app, the sessions object will be reset to nothing. This won't affect us that badly right now, but it is something that production web apps must remedy.

## Checkpoint

Take a screenshot of your home page while logged in as `melinoe`.

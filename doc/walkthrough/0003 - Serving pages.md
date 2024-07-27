# Serving pages

2024-07-27

We saw in the previous section that we can just send HTML-formatted strings to our response writers. This will work for a while. My [personal site](https://joeilagan.com) actually renders its HTML this way. If we're making a real web app, though, we'll need to move on.

Most web frameworks store their HTML in separate _template files_ that the web server reads, injects data into, and finally sends to the client. We will create a similar setup for ourselves now using only Go's standard library.

## Template files

A template file in Go is just an `.html` file that optionally contains some special syntax for making use of data that Go can inject into it.

Let's start by making a very simple template. Its only job will be to display the "CafeGo" header.

Create a new directory `templates/` in `cafegoroot/`. Inside `templates/`, create a new file `index.html`, and write the following.

```html
<h1>CafeGo</h1>
```

We may now read it from our `main.go` file. Before we do, let's clean up a little. Rename `handler` to `indexHandler` so that the entire file looks like this:

```go
package main

import (
    "fmt"
	"net/http"
)

type IndexPageData struct{}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "<h1>CafeGo</h1>")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":5000", nil)
}
```

So far so good. Now, let's learn how to read a template.

Go provides a package `html/template` that we can use to read template files and inject data into them. We can read a template file like so:

```go
tmpl, err := template.ParseFiles("./templates/index.html")
```

(Remember that VS Code's Go extension can automatically import packages for you if you just use them in a source code file and save it.)

The `.ParseFiles` function from `html/template` produces two values: first, a reference to a `template.Template` struct, and second, an error. This is how Go handles operations that can potentially fail. Instead of crashing the program, risky operations produce a second return result to hold an error.

This does mean that we need to handle errors after every risky operation. This is typically done like so:

```go
if err != nil {
    log.Fatal(err)
}
```

All this does is crash the program and print the error if a risky operation produced an error. Of course, there's no _requirement_ to handle the error like this, but if you want the default `error => crash` experience, this is how to get it in Go.

Once you have your template struct, you can make the `template` package render it to a file-like interface like so:

```go
// Another Go quirk. := means "new assignment" and = means "re-assignment."
// This is an = and not a := because we already have a variable named err, we are just shadowing it.
err = tmpl.Execute(w, nil)
// And, again, we need to handle the error
if err != nil {
    log.Fatal(err)
}
```

This function `.Execute` takes two arguments: a file-like interface and some optional data. Since we are not trying to inject data into our template yet, we can submit `nil` to the second argument.

That's actually it. Your entire program should now look like this:

```go
package main

import (
	"html/template"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":5000", nil)
}
```

There are a lot of lines involved here, but you may be beginning to get a feel for Go's philosophy. You may need to write a lot of code, but each individual line you write is dumb simple. This will be a huge benefit to us later.

## Injecting data

Of course, one of the main points of templates is to be able to inject data into them. Let's prototype this now by pretending that we have some user data already.

Revisit `templates/index.html` and add this line:

```html
<p>Welcome, {{ .Username }}!</p>
```

Where will our template get this field `Username`? `.Execute` requires us to pass in data as an instance of a _struct_, which you can think of as just a strict definition for a type of hashmap/dictionary/object if you come from a dynamic language like Python or JavaScript. We can define a type of struct like this:

```go
type IndexPageData struct {
	Username string
}
```

We can now create _instances_ of this struct elsewhere in our program. If we pretend that our user is named "Matthew", we can create an instance of `IndexPageData` like this:

```go
sampleUsername := "Matthew"
samplePageData := IndexPageData{Username: sampleUsername}
```

And finally, we can pass the page data to our template like this:

```go
err = tmpl.Execute(w, samplePageData)
```

So now you can see the flow of data. We create a _struct_ with one field, Username. We _instantiate_ the struct with the specific username "Matthew". We pass the instance of that struct to our template executor. The template now has access to the struct data, which it can access with the syntax {{ .Username }}, or more generally, {{ .FieldName }}.

## Checkpoint

Submit a screenshot of your index page, but the with the name Vic instead of Matthew.

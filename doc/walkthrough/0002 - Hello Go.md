# Hello Go

2024-07-27

In this section of the CafeGo walkthrough, we will set up our Go project.

The rest of this walkthrough assumes that you are familiar with your terminal or that you are at least capable of learning it. We will also be using UNIX file path conventions. Directories will be separated by forward slashes (e.g., `cafegoroot/main.go`). Mac OS follows UNIX file path conventions. If you are on Windows, forward slashes will be replaced by backward slashes (e.g., `cafegoroot\main.go`).

## Hello world

In `cafegoroot`, initialize a Go project with `go mod init example.com/cafego`.

Almost right away, we can see one of Go's quirks. `go mod init` makes enough sense as a project init command, but what is `example.com/cafego`? Go packages are identified not only with their name, which is `cafego` in this case, but also a domain. This domain doesn't actually need to be _real_, which is why we can use `example.com`, but programmers who intend to host Go packages tend to make this either `github.com` or a domain that their company controls. This will help Go packages be uniquely identifiable on the internet, but it doesn't actually matter in this case.

The "hello world" of Go is straightforward. Create a new file `main.go` and write the following code.

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello world")
}
```

Run `go run .` in your terminal. This should print out `Hello world`. If this works, then you may proceed.

This might look a little strange if you come from Python or JavaScript, but it should look perfectly normal if you come from C or Java. Each Go project has a "main" function that acts as the entry point for the program. Code will start running here.

The "hello world" HTTP server of Go is also straightforward to seasoned programmers. Change the code in `main.go` to the following.

```go
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":5000", nil)
}

```

Let's pause for a moment here to dissect what we see.

- The top line, `package main`, declares what package the current file belongs to. The "main" package is a special package that more or less indicates that the current package is meant to be run as its own program, not a library for another program to import.
- The `import` block brings other packages like `fmt` and `net/http` into scope. These two specifically are part of Go's standard library.
    - `fmt` stands for "format", and it is used for creating "formatted input and output." We typically use this to print output, save strings to variables, and write to files and other file-like interfaces.
    - `net/http` stands for "networking" and "hypertext transfer protocol" respectively. `http` is used for both HTTP clients and servers.
- Our `handler` function should look very familiar to anyone who's ever done HTTP before. This is just a function that is meant to run when an HTTP server receives an HTTP request.
    - The `w` argument here is of type `http.ResponseWriter`. It is a file-like interface, which means that to write to it, you need to treat it like a file.
    - The `r` argument here is a _reference_ to a value of type `http.Request`. This is where you can read the properties of an HTTP request like the query parameters, POST bodies, and form data.
- Within our `handler` function, we simply write `"Hello world"` to the Response Writer. Remember that a Response Writer is a file-like interface, so we have to use `fmt.Fprint` (i.e., file-print, or "write this text to a file") to get data to it.
- In our `main` function, we register our handler function to handle the index route `/`. We then tell Go to spin up an HTTP server on port 5000.

(As an aside: If you have VS Code and the Go extension, you can simply use the name of standard libraries like `fmt` and `http` without importing them. When you save your code, VS Code will automatically import them for you. When you remove all references to them, VS Code will automatically remove them from the import block.)

You may have noticed that uncanny feeling I mentioned in the README, especially if you come from another popular contemporary language. These explanations of what Go does make enough sense on their own, but there are things that just seem strange. Let's go over two.

- Why do we import packages with strings? In Python, you just use the symbol itself, like `import requests`. In JavaScript, you also import with strings, but (at least using the require syntax) it feels like a normal function call. This method of importing just feels off.
- Why do these package methods have an uppercase first letter? Well, if you want to export a function from a package, you have to name it with an uppercase first letter. Most other languages have some sort of `export` keyword that does this instead of making capitalization a meaningful syntax difference.

I do think that these strange things will start to make sense eventually, but I cannot deny that they really add up when you're learning. I am normally loath to recommend using AI to learn programming, but in Go's case, combining AI (hopefully one that gives you sources like Perplexity at least) and your personal library of snippets might be a very effective way to get a handle on this language.

Here's something interesting, though. We're actually done with the hello world. Never mind how strange Go can feel to write -- it is extremely straightforward as a language, and it tends to get out of your way when you're driving towards business value.

Let's make one small change before we move on to the next section. Go works best if you know the fundamentals. One such fundamental is that _HTTP messages are just text_. We can actually send a string, formatted as HTML, directly to the browser, and the browser will try to render the HTML.

Change the code in your `handler` function to this:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	markup := "<h1>CafeGo</h1>"
	fmt.Fprint(w, markup)
}
```

There's absolutely nothing special about that string other than that it represents valid HTML. If you `fmt.Fprint` it to the response writer, it will send the string as-is, and the browser will try to interpret it as HTML.

## Checkpoint

Run the server. Submit a screenshot of your browser, pointed at `http://localhost:3000`.

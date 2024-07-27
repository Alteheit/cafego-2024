# CafeGo

2024-07-26

Welcome to CafeGo. This is an introduction to web application programming. Our goal in this walkthrough is to teach you the very basics of how a "web application," as opposed to a mere "web page," works.

You can think of a web application (henceforth "web app") as a web page that stores data. This sounds like a simple addition, but it isn't. Having to manage data is where almost all the complexity of modern web development comes from. Thankfully, special programming libraries called frameworks are available to help us design and run such applications easily.

In this walkthrough, we will be building a simple e-commerce application for a company that sells coffee. Users should be able to browse products, add products to their cart, and check out their orders. There are a lot of requirements in that one sentence.

We will build CafeGo with very little other than the standard library of Go. This is what makes Go a very interesting language choice. With other languages like Python and JavaScript, much of what you need must come from third-party libraries. When people try to sell their language of choice, they usually focus on the strength of their ecosystems. Go takes a different approach. Philosophically, the maintainers of Go believe in having a strong standard library that can handle almost anything, so you very rarely need to bring in external libraries. You still can, of course, but it is liberating to not be dependent on the good work of a web of strangers.

The repository you are in right now is a snapshot of the code we had after writing the documents in the walkthrough folder. Note that you might not be able to download and run this repository directly. It shouldn't matter. You are meant to build CafeGo from scratch in your own repository. Use this repo as a guide, not as something to copy in its entirety.

## Objectives

After completing this walkthrough, you should be able to:

1. Build a simple web app from scratch with Go,
2. Reason about how to add a new feature to a Go web app.

I should note that I still think that Python is the best language for this exercise. While I respect Go as a software professional, I have a feeling that it abstracts too little for beginners to make good use of it. Regardless, we will push on. Go's relative lack of abstraction will make things more difficult but by no means _impossible_, and I have faith in your desire to learn.

We will also have to assume that you have some familiarity with programming in general for the rest of this walkthrough. I will explain things when relevant, especially around Go's quirks, but I cannot teach you programming here.

## Prework

### Installations

Bear with us, there is some setup that you need to do before running a Go project.

We are fortunate that the Go team places great emphasis on making Go easy to install. Whether you are on Windows or on Mac, go to the [Go website](https://go.dev) and navigate to the Downloads page. Download the appropriate package for your platform and follow the wizard.

Again, whether you are on Windows or Mac, open a terminal/PowerShell/etc and check if Go was installed.

```zsh
go version
```

You should see something similar to `go version go1.22.1 {platform}`. If you do, then you may proceed.

Make an empty directory and set it as your new working directory. This new directory will serve as the "root" of our entire Go project. We will call this directory `cafegoroot` for this tutorial.

### VS Code

I also highly recommend the use of VS Code and the corresponding Go extensions for this project. Go is rather annoying to type and format. VS Code and the Go extension handle this rather well, though they can be a bit opinionated.

### Theory

This part is optional, but highly helpful. A web app needs four parts to function:

- Routing HTTP requests based on their content and path ("routing")
- Storing data in a database and making it accessible to the rest of the app ("models")
- Rendering HTML based on data that may vary ("views")
- The glue logic between models and views ("controllers")

Web frameworks like Django, Ruby on Rails, and more all use these four basic concepts. Microframeworks like Express still use these four basic concepts, but they tend to leave the minutiae to you. I'm not sure it's appropriate to call Go a framework or a microframework _per se_, but it's definitely the sort of environment that leaves the decisions to you, so it will greatly help to understand the conceptual framework of what we're doing here.

I wrote about these in this series of articles:

- https://joeilagan.com/article/2024-itm-web-apps-1
- https://joeilagan.com/article/2024-itm-web-apps-2
- https://joeilagan.com/article/2024-itm-web-apps-3

## Checkpoint

Please take a screenshot of your Terminal/PowerShell window after running `go version`.

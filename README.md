# CafeGo

2024-07-26

This is the repository for the CafeGo walkthrough. This is the Go-based version of Digital Cafe, which is the scenario we use to teach students of ITMGT how to build a web application.

If you are a student aiming to complete the CafeGo assignment, please head to `doc/walkthrough` and go through each of the documents starting from 0001.

## Why Go

Go is not a common language in classrooms. As far as languages are, it is relatively new. It first appeared in 2009. It hasn't had the time to sneak into the mindshare of academe yet, unlike Python (1991) and JavaScript (1995).

Maybe I should say that it isn't common _yet_. I think it will eventually become an attractive beginner's language, probably because it's an extremely _dumb_ language, and I say that with the utmost respect. Go is honestly not fun to write at all. Its syntax feels clunky. It has very few features compared to its contemporaries. You won't ever feel smart while writing it, except maybe when you do concurrent programming with goroutines and channels. But that's what I think makes it a _wise_ choice for a programming language.

I have written a few small services and side projects with Go, and they have been by far my most successful non-Python projects. I just dislike the experience of writing Go so much that I focus almost entirely on finishing the project, which I end up doing in record time because there's often only one sensible way to do what a project needs. I have occasionally had to return to my Go projects to add a new feature, but even after being away from the codebase for a while, I find that the syntax is so straightforward that I can grok it with very little effort. If you've programmed at all before, you will know how desirable of a property this is. Imagine doing that with almost any other language.

Once Go code is written, it's incredibly stable, too. Part of its philosophy is that you should not have to depend on external packages, so the core functionality of my Go programs depends only on the language itself, which is incredible in today's environment. It helps that the maintainers of the language have committed to not breaking the interface of existing Go code even as they add new features to the language.

That might sound like a dubious sales pitch, but I do encourage you to give it a try. I still think that Python is the king, but I've come to believe that Go is its hand.

## Fundamentals

It should be clear by now that I like Go. I do understand, though, that Go might not be the _easiest_ choice for a beginner, mostly because you actually need to know things about how computers work to get the most out of it.

Python, being as flexible as it is, has spawned a host of libraries that competently hide the details of how things work. Go's ecosystem does not have quite the same level of polish. This is going to get particularly interesting, I think, when we head into CafeGo's sections on HTTP cookies and database calls. A Python framework like Django hides these details really well, but we Go programmers will need to study the underlying structure of HTTP before we can remember a user. Perhaps worse, we will need to study SQL* before we use a database.

There's another thing about Go that's a little off-putting. Even though its syntax is easy, I can't help but feel that it's somewhat uncanny in a way that makes it difficult to feel like you're really _learning_ the language. I can't really explain why. I do recommend keeping a library of "how-to" snippets as we go along. Index your snippets by use case, like "how to open a file" and "how to write an HTTP handler."

*I'm kidding about SQL being horrifying, to be clear. It's not that bad. It's unwieldy, but it's unavoidable, so you'll have to get over it at some point.

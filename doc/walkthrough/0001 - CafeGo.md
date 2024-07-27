# CafeGo

2024-07-26

Welcome to CafeGo. This is an introduction to web application programming. Our goal in this walkthrough is to teach you the very basics of how a "web application," as opposed to a mere "web page," works.

You can think of a web application (henceforth "web app") as a web page that stores data. This sounds like a simple addition, but it isn't. Having to manage data is where almost all the complexity of modern web development comes from. Thankfully, special programming libraries called frameworks are available to help us design and run such applications easily.

In this walkthrough, we will be building a simple e-commerce application for a company that sells coffee. Users should be able to browse products, add products to their cart, and check out their orders. There are a lot of requirements in that one sentence.

We will build CafeGo with very little other than the standard library of Go. This is what makes Go a very interesting language choice. With other languages like Python and JavaScript, much of what you need must come from third-party libraries. When people try to sell their language of choice, they usually focus on the strength of their ecosystems. Go takes a different approach. Philosophically, the maintainers of Go believe in having a strong standard library that can handle almost anything, so you very rarely need to bring in external libraries. You still can, of course, but it is liberating to not be dependent on the good work of a web of strangers.

The repository you are in right now is a snapshot of the code we had after writing the documents in the walkthrough folder. Note that you might not be able to download and run this repository directly. It shouldn't matter. You are meant to build CafeGo from scratch in your own repository. Use this repo as a guide, not as something to copy in its entirety.

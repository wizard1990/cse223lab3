## Machines

We have set up a cluster of 10 machines. You should use those for
all of the lab assignments.

- c08-11.sysnet.ucsd.edu
- c08-12.sysnet.ucsd.edu
- c08-13.sysnet.ucsd.edu
- c08-14.sysnet.ucsd.edu
- c08-15.sysnet.ucsd.edu
- c08-16.sysnet.ucsd.edu
- c08-17.sysnet.ucsd.edu
- c08-18.sysnet.ucsd.edu
- c08-19.sysnet.ucsd.edu
- c08-20.sysnet.ucsd.edu

## Programming Language

You will write the labs in Google's [golang](http://golang.org).  It
is a young language with a language syntax at somewhere between C/C++
and Python. It comes with a very rich standard library, and also
language-level support for light-weight but powerful concurrency
semantics like *go routines* and *channels*.

Here is some key documentation on the language:

- [Go Language Documentation Page](http://golang.org/doc/)
- [Effective Go](http://golang.org/doc/effective_go.html)
- [Go Language Spec](http://golang.org/ref/spec)

While you should be able to find a lot of documents about Go language
on the web, especially from the official site. If you know C++
already, here are some hints that might help you bootstrap.

- Go language code is organized in many separate *packages*.
- Different from C/C++, when defining a *variable* or *constant*, the
  *type* of it is written after the variable name.
- Go language has pointers, but has no pointer
  arithmetics. For example, you cannot increase a pointer by 1, to
  point the next element in memory.
- Go language has fixed length *arrays*.
  However, arrays are not very commonly used.  For most of the time,
  people use *slices*, which is a sliced view of an underlying array
  that is often implicitly declared.
- *maps* are built-in hash-based dictionaries.
- A function can have multiple return values.
- Exceptions are called `panic` and `recover`. However it is not
  encouraged to use that for error handling.
- `for` is the only loop keyword.
- *Foreach* is implemented with `range` keyword.
- Semicolons at the end of statements are optional.
- On the other hand though, trailing comma in a list is a must.
- Variables are garbage collected. The language is hence
  type safe and pointer safe. When you have a pointer,
  the content it points to is always
  valid.
- Identifier that starts with
  a capital letter is *public* and visible to other packages; others
  are *private* and only visible inside its own package.
- *Inheritance* is done by compositions of anonymous members.
- Virtual functions are binded via *interfaces*. Unlike Java,
  *interface* does not require explicit binding (via the *implements*
  keyword). As long as the type has the set of methods implemented, it
  can be automatically assigned to an inteface. As a result, it is
  okay to write the implementation first and declare the interface
  afterwards.
- Circular package dependency is not allowed.

## The Tribbler Story

Here is the story: some cowboy programmer wrote a
simple online microblogging service called Tribbler, and leveraging
the power of the Web, it becomes quite popular. However,
the program runs in one single process; it does not scale,
cannot support many concurrent connections,
and is vulnerable to machine crashes. Knowing that you
are taking the distributed computing system course at UCSD, he asks
you for help. You answered his call and started this project.

Your goal is to refactor Tribbler into a distributed system,
make it robust and scalable.

## Getting Started

The Tribbler project is written in golang and stored in a git
repository now. To get started, run these commands in command line:

```
$ cd                       # go to your home directory
$ mkdir -p gopath/src      # the path you use for storing golang src
$ cd gopath/src
$ git clone /classes/cse223b/sp14/labs/trib -b lab1
$ git clone /classes/cse223b/sp14/labs/triblab -b lab1
$ export GOPATH=~/gopath
$ go install ./...
```

Do some basic testing see if the framework is in good shape:

```
$ go test ./trib/...
```

Now The basic Tribbler service should be installed on
the system in your home directory. Let's give it a try:

```
$ ~/gopath/bin/trib-front -init -addr=:rand
```

The program should show that it serves on a port (which is randomly
generated).

Now open your browser and type in the address. For example, if the
machine you logged in was `c08-11.sysnet.ucsd.edu`, and Tribbler is
running on port 27944, then open `http://c08-11.sysnet.ucsd.edu:27944`.  You should see a list of Tribbler users, where you can view their tribs and login as them (with no authentication).

This is how the Tribbler service looks like to the user clients.
It is a single Web page the performs AJAX calls (a type of RPC
that is widely used in Web 2.0) to the web server behind. The
webserver then in turn calls the Tribbler logic functions
and returns the results back to the Web page in the
browser.

If you find it difficult to access the lab machines outside UCSD
campus, you need to setup a UCSD VPN or ssh tunnel.

## Source Code Organization

The source code in the `trib` package repository is organized as follow:

- `trib` defines the common Tribbler interfaces and data structures.
- `trib/tribtest` provides several basic test cases for the
  interfaces.
- `trib/cmd/trib-front` is the web-server launcher that you just run.
- `trib/cmd/kv-client` is a command line key-value RPC client
  for quick testing.
- `trib/cmd/kv-server` runs a key-value service as an RPC server.
- `trib/cmd/bins-client` is a bin storage service client.
- `trib/cmd/bins-back` is a bin storage service back-end launcher.
- `trib/cmd/bins-keeper` is a bin stroage service keeper launcher.
- `trib/cmd/bins-mkrc` generates a bin storage configuration file.
- `trib/entries` defines helper several functions on
  constructing a Tribbler front-end or a back-end.
- `trib/ref` is a reference monolithic implementation of the
  `trib.Server` interface. All the server logic runs in one single process.
  It is not scalable and vulnerable to machine crashes.
- `trib/store` contains an in-memory thread-safe implementation of the
  `trib.Store` interface. We will use this as
  the basic building block for our back-end storage system.
- `trib/randaddr` provides helper functions that generate a network
  address with a random port number.
- `trib/local` provides helper functions that check if an address
  belongs to the machine that the program is running.
- `trib/colon` provides helper functions that escape and unescape
  colons in a string.
- `trib/www` contains the static files (html, css, js, etc.) for the
  web front-end.

Don't be scared by the number of packages. Most of the packages are
very small. In fact, all Go language files under `trib` directory is
less than 2500 lines in total (the beauty of Go!).

Through the entire lab, you do not need to (and should not) modify anything in
this `trib` repository. If you feel that you have to change some code to
complete your lab, please discuss with the TA. You are always welcome to read
the code in `trib` repository. If you find any bug and reported it, you might
get some bonus credit.

## Your Job

Your job is to complete the implementation of the `triblab` package.
It is in the second repo that we checked out.

It would be a good practice for you to periodically commit your code
into your own `triblab` git repo. Only files commited in that repo
will be submitted for grading.  

## Lab Roadmap

- **Lab 1**. Wrap the key-value pair service interface with RPC, so
  that a remote client can call the service via network connections.
- **Lab 2**. Reimplement the Tribbler service, split the current
  Tribbler logic into stateless scalable front-ends and key-value
  pair scalable back-ends. The front-ends will call the back-ends via
  RPCs that is implemented in Lab 1. When this lab is done, we should
  have both the front-end and the back-end scalable.
- **Lab 3**. We make the back-ends fault-tolerent, by applying
  techniques like distributed hash table and replications. As a result, at the
  end of this lab, back-end servers can join, leave, or be killed, without
  breaking down the entire service.

By the end of the labs, you will have a new Tribbler service
implementation that is scalable and fault-tolerant.

## Misc

For convenience, you might set environment variables in your `.bashrc`
and/or `.bash_profile`:

```
export GOPATH=$HOME/gopath
export PATH=$PATH:$GOPATH/bin
```

We should have Vim and Emaces installed on the machines. If you need
to install other utility packages, ask the TA. Note that you do not
have `sudo` permissions on any of the machines; any `sudo` attempt
will be automatically reported, so please don't even try it.

You could also write your code on your own machine if you want to.
See Go language's [install](http://golang.org/doc/install) page for
more information. However, you should test your code on the lab
machines.

## Ready?

If you feel comfortable with the lab setup now,
go forward and read [Lab1](./lab1.html).

## Lab1

Welcome to Lab1. The goal of this lab is to implement a
key-value pair service that can be called via RPC. In particular
you need to:

1. Implement a key-value storage server type that wraps a `trib.Store`
interface object and takes http RPC requests from the network.
2. Implement a key-value storage client type that fits `trib.Store`
interface and RPCs a remote key-value pair server.

More specifically, you need to implement two entry functions that are
defined in `triblab/lab1.go` file: `ServeBack()` and `NewClient()`.
Now, they are both implemented with `panic("todo")`.

## Get Your Repo Up-to-date

If you cloned your source folder before Tuesday April, 1st.
You might need to first get your repo up-to-date:

```
$ cd ~/gopath/src/trib
$ git pull /classes/cse223b/sp14/labs/trib lab1
$ cd ~/gopath/src/triblab
$ git pull /classes/cse223b/sp14/labs/triblab lab1
```

The instructions here assume you used the the default directory
setup.

## The Key-value Pair Service Interface

The goal of Lab1 is to wrap a key-value pair interface
with RPC. Although you don't need to implement the key-value
pair storage by
yourself, you need to use it extensively in later labs, so it will be
good for you to understand the service semantics here.

The data structure and interfaces for the key-value pair service is
defined in `trib/kv.go` file (in the `trib` repository). The main
interface is `trib.Storage` interface, which consists of three parts.

First is the key-string pair part.

```
// Key-value pair interfaces
// Default value for all keys is empty string
type KeyString interface {
	// Gets a value. Empty string by default.
	Get(key string, value *string) error

	// Set kv.Key to kv.Value. Set succ to true when no error.
	Set(kv *KeyValue, succ *bool) error

	// List all the keys of non-empty pairs where the key matches
	// the given pattern.
	Keys(p *Pattern, list *List) error
}
```

`Pattern` is a prefix-suffix tuple. It has a `Match(string)` function
that returns true when the string matches the pattern.

Second is the key-list pair part.

```
// Key-list interfaces.
// Default value for all lists is an empty list.
// After the call, list.L should never by nil.
type KeyList interface {
	// Get the list.
	ListGet(key string, list *List) error

	// Append a string to the list. Set succ to true when no error.
	ListAppend(kv *KeyValue, succ *bool) error

	// Removes all elements that equal to kv.Value in list kv.Key
	// n is set to the number of elements removed.
	ListRemove(kv *KeyValue, n *int) error

	// List all the keys of non-empty lists, where the key matches
	// the given pattern.
	ListKeys(p *Pattern, list *List) error
}
```

And finally we put it together, and add an auto-incrementing clock service:

```
type Storage interface {
	// Returns an auto-incrementing clock, the returned value
	// will be no smaller than atLeast, and it will be
	// strictly larger than the value returned last time,
	// unless it was math.MaxUint64.
	Clock(atLeast uint64, ret *uint64) error

	KeyString
	KeyList
}
```

Note that the function signature of these methods are all RPC
friendly. You should directly implement the RPC interface with Go
language's [`rpc`](http://golang.org/pkg/net) package.  By doing this,
another person's client that speaks the same protocol will be able to
talk to your server as well.

Under the definition of the execution logic, all the methods will
always return `nil` error. Hence all errors you see from this
interface will be communication errors. You can assume that each call
(on the same key) is an atomic transaction; two concurrent writes
won't give the key a weird value that comes from nowhere.  However,
when an error occurs, the caller won't know if the transaction is
committed or not, because the error might occur before or after the
transaction.

## Entry Functions

These are the two entry functions you need to implement for
this Lab.
This is how other people's code (and your own code in later
labs) will use your code.

```
func ServeBack(b *trib.Back) error
```

This function creates an instance of a back-end server based on
configuration `b *trib.Back`. Structure `trib.Back` is defined in
`trib/config.go` file.  In the struct type, it has several fields:

- `Addr` is the address the server should listen on, in the form of
  `<host>:<port>`. Go language uses this address in its [standard
  `net` package] (http://golang.org/pkg/net), so you should be able to
  use it directly on opening connections.  
- `Store` is the storage device you will use for storing stuff. In
  fact, You should not store persistent data anywhere else.
  `Store` will never be nil.
- `Ready` is a channel for notifying the other parts in the program
  that the server is ready to accept RPC calls from the network
  (by sending value `true`) or the server failed to setup the
  connection (by sending value `false`). `Ready` might be nil (means
  the caller does not care about when it is ready).

This function should be a blocking call. It does not return until an
error (like the network is shutdown) occurred.

Note that you don't need to (and should not) implement the key-value
pair storage by yourself.  You only need to wrap the given `Store`
with RPC, so that a remote client can access it via the network.

***

```
func NewClient(addr string) trib.Stroage
```

This function takes `addr` as a TCP address in the form of
`<host>:<port>`, and connects to this address for an http
RPC server. It returns an implementation of `trib.Storage`, which
will provide the interface, but all calls will be actually RPCs
to the server. You can assume `addr` will always be a valid TCP
address.

Note that when `NewClient()` is called, the server might not start
up yet. While it is okay to make a try to connect the server at this
time, you should not report any error if your attempt failed.  It
might be better to establish the connection when you need to perform
your first RPC function call.

## The RPC Package

Go language comes with its own
[`net/rpc`](http://golang.org/pkg/net/rpc) package in its standard
library, and we will just use that.  Note that the `trib.Store`
interface is already in its "RPC friendly" form.

Your RPC needs to use the default encoding `encoding/gob`, listen on
the given address, and serve as an http RPC server. The server
needs to register the back-end key-value pair object under the
name `Storage`.

## Testing

Both the `trib` and `triblab` repository comes with a makefile with
some handy command line shorthands, and also some basic testing code.

Under the `trib` directory, if you type `make test`, you should see
that the test runs and all passed.

Under the `triblab` directory, if you type `make test` however, you
would see the test fails with a todo panic if you have not implement.

You should implement the logic and try to pass those test cases. If
you pass those, you should be fairly confident that you can get at
least 30% of the credits for Lab1 (unless you are cheating in some
way).

However, the test that comes with the repository is very basic and
simple.  Though you don't have to, you should really write more test
cases to make sure your implementation matches the specification.

For more information on writing test cases in Go language, please read
the [testing](http://golang.org/pkg/testing/) package document page.

## Starting Hints

While you are free to do the project in your own way as long as
it fits the specification, matches the interfaces and passes the
tests, here are some suggested steps for you to start.

First, create a `client.go` file under `triblab` repo, and declare a
new struct type called `client`:

```
package triblab

type client struct {
    // your private fields will go here
}
```

Then add method functions to this new `client` type so that
it matches `trib.Storage` interface. For example, for the `Get()`
function:

```
func (self *client) Get(key string, value *string) error {
    panic("todo")
}
```

After you added all the functions, you can write a line for compile
time checking if all the functions are implemented:

```
var _ trib.Storage = new(client)
```

This creates a zero-filled `client` and assigns it to an anonymous
object of `trig.Storage` interface. Your code hence only compiles when
`client` satisfies the interface. (Since this zero-filled variable is
anonymous and nobody can access it, it will be removed as dead code by
the compiler optimizer and hence has no negative effect to the
run-time execution.)

Now add a field into `client` called `addr`, which will save the
server address. Now `client` looks like this:

```
type client struct {
    addr string
}
```

Now that we have a client type that satisfies `trib.Stroage`, we
can return this type in our entry function `NewClient()`. So remove
the `panic("todo")` line in `NewClient()`, and replace it by
returning a new `client` object. Now the `NewClient()` function
should somehow look like this:

```
func NewClient(addr string) trib.Storage {
    return &client{addr: addr}
}
```

Now we have the code skeleton for the RPC client, and we will fill
in the actual logic that performs the RPC calls.

To do an RPC call, we need to import the `rpc` package, so at the
start of `client.go` file, lets import that after the package name
statement.

```
import (
    "net/rpc"
)
```

And following the examples given in the `rpc` package, we can
write the RPC client logic. For example, the `Get()` method
should somehow look like this:

```
func (self *client) Get(key string, value *string) error {
    // connect to the server
    conn, e := rpc.DialHTTP("tcp", self.addr)
    if e != nil {
        return e
    }

    // perform the call
    e = conn.Call("Storage.Get", key, value)
    if e != nil {
        conn.Close()
        return e
    }

    // close the connection
    return conn.Close()
}
```

However, note that if you do it this way, you will open a new HTTP
connection for every RPC call. It is okay but obviously not the most
efficient way to do so.  I will leave it for yourself to figure out
how to maintain a persistent RPC connection.

That was the client side. You also need to wrap the server side in the
`ServeBack()` function using the `rpc` library. This should be pretty
straight-forward by creating an RPC server, registering the `Store`
member field in `b *trib.Config` parameter under the name of
`Storage`, and serving it as an HTTP server. The code should be
similar to one of the examples given in the
[`rpc`](http://golang.org/pkg/net/rpc) package documentation. Just
remember that you need to register as `Storage` and also need to send
a `true` over the `Ready` channel when the service is ready (when
`Ready` is not `nil`), but send a `false` when you encounter any error
on starting your service.

When all of those are done, you should pass the test cases written in
`back_test.go` file. It calls the `CheckStorage()` function defined
in `trib/tribtest` package, and performs some basic checks on if an
RPC client and a server (that runs on the same host) will satisfy the
specification of a key-value pair service (as a local
`trib/store.Storage` will do without RPC).

## Playing with It

To do some simple test with your own implementation, you can use the
`kv-client` and `kv-server` command line launcher.

First make sure your code compiles.

Then run the server.

```
$ kv-server
```

*(You might need to add `$GOPATH/bin` into your `$PATH` to run this.)*

And you should see an address printing out, say it is
`localhost:12086`. (Note that you can also specify your own address
via command line. The default address is `localhost:rand`.)

Now you can play with your server via the `kv-client` program.
For example:

```
$ kv-client localhost:12086 get hello

$ kv-client localhost:12086 set h8liu run
true
$ kv-client localhost:12086 get h8liu
run
$ kv-client localhost:12086 keys h8
h8liu
$ kv-client localhost:12086 list-get hello
$ kv-client localhost:12086 list-get h8liu
$ kv-client localhost:12086 list-append h8liu something
true
$ kv-client localhost:12086 list-get h8liu
something
$ kv-client localhost:12086 clock
0
$ kv-client localhost:12086 clock
1
$ kv-client localhost:12086 clock
2
$ kv-client localhost:12086 clock 200
200
```

## Requirements

- When the network and the storage is errorless, RPC to your server
  should not return any error.
- When the network has error (like the back-end server crashed, and
  the client hence cannot connect), your RPC client should return
  error. However when the server is back up running, your RPC client
  should act as normal again (without the need of creating a new
  client). 
- When the server and the clients are running on the lab machines, for
  each function call, the latency introduced by your RPC (comparing
  with direct local function calls) should be less than 0.1 second.

## Turning In

First, make sure that you have committed every piece of your code into
the repository `triblab`. Then just type `make turnin` under the root
of the repository.  It will generate a `turnin.zip` that contains
everything in your git repository, and will then copy the zip file to
a place where only the lab instructors can read.

## Happy Lab1!

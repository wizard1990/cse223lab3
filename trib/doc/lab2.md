## Lab2

Welcome to Lab2. The goal of this lab is to use the RPC key-value pair
we built in Lab 1 to implement scalable Tribbler front-ends and
back-ends. In particular, you need to first wrap the key-value
pair server into a distributed key-value pair service called
*bin storage*, where your implementation should fit `trib.BinStorage`
interface, and then on top of this storage layer,
you need to implement a stateless Tribbler front-end type
that fits `trib.Server` interface.

## Get Your Repo Up-to-date

First we need to update the `trib` repo:

```
$ cd ~/gopath/src/trib
$ git pull /classes/cse223b/sp14/labs/trib lab2
```

If you have not changed anything in `trib` repo, this should be
painless. However, if you changed stuff, you need to merge the
changes.

Also update the `triblab` repo to merge branch `lab2`. There will be
several changes:

- Some line changes in `makefile`.
- Some added lines in `lab2.go`.
- A new file called `server_test.go`.
- A new file called `bins_test.go`.

If you have not touched those files and have not created file called
`server_test.go` or `bins_test.go` by yourself, the merge should be
painless:

```
$ cd ~/gopath/src/triblab
$ git pull /classes/cse223b/sp14/labs/triblab lab2
```

If you have made changes to those files. Then you need to merge the
changes properly.

If you have trouble on merging the changes (on any of the repos),
and don't know how to solve the conflicts, please ask the TA for help.

## System Architecture

The overall system architecture looks like this:

![System Arch](./arch.png)

As you can see from the figure, our system is split into two parts,
where the upper part is the Tribbler service, and the lower part is a
general purpose distributed key-value store called "Bin Storage".  The
Tribbler front-ends take AJAX calls from the users and translate the
calls into bin storage calls, where the bin storage client will issue
the calls into the "dumb" key-value back-ends via RPC. Each back-end
here is a simple independent key-value store (which we just
implemented in the previous lab).  There is also a keeper instance
running in the background that maintains the bin storage in a coherent
state. Both the bin storage clients and the keeper use the same RPC
interface to communicate with the back-ends. The interface semantics
defined in `trib.Storage` should be sufficient to perform all the bin
storage tasks.

## Bin Storage

**Bin Storage** is a distributed key-value pair storage that combines
a bunch of `trib.Storage` back-ends and presents a single and
scalable storage service over the network. Bin storage service is
provided via the `trib.BinStorage` interface.

On one hand, conceptually, bin storage combines an infinite set of
separated `trib.Storage` instances called *bins*. Each bin has a
different name, and a caller can fetch a bin via the
`BinStorage.Bin()` function -- the only function in the
`trib.BinStorage` interface.  Initially, all bins are empty.

On the other hand, physically, a bin storage is built on top of a
finite set of key-value storage back-ends (which we built in the last
lab), where a bin (a virtual KV store) will be mapped to one of the
key-value back-ends (physical KV stores). As a result, multiple bins
might share a same back-end by splitting its key space with prefixes
(or suffixes) that encode the bin names.  As shown in the figure
below, an example bin storage hashes 6 bins into 3 back-ends, storing
each bin with different prefixes inside.  When a caller calls
`Get("some-key")` on the bin `"alice"` for example, the bin storage
client will translate that call into a `Get("alice::some-key")` RPC
call on back-end 0.

![Bin Storage](./bins.png)

In addition to key-value stores in bins, a bin storage also coarsely
synchronizes the logical clock across *all* the bins, so that the
clocks won't be offset too much over time. Since the back-ends are
"dumb" and they do not talk to each other, this clock synchronization
is done by a background keeper process. Note that the clocks are NOT
real-time synchronized (which would be too expensive for a distributed
system), but only coarsely synchronized. In particular it means that
if a `Clock()` is issued at least 3 seconds after another `Clock()`
call, no matter which bins they are issued in, bin storage always
guarantees that the later one always returns a clock value no smaller
than the earlier one's.  Note that by definition, a bin storage with
one single back-end will automatically satisfy this requirement, but
when there are multiple back-ends, the keeper must do some maintenance
from time to time.

As said, in Lab1, we already implemented the back-ends, the APIs for
both the server and the client, so we will just reuse them without
changes (assuming that you implemented them correctly). Both the bin
storage client and the keeper will communicate with the "dumb"
back-ends via the RPC calls we implemented, and cooperatively they
present a coherent bin storage interface to upper layer applications.

## Tribble

Before we look into how the Tribbler service works, let's first define
what a tribble is.

A Tribble is a structure type that has 4 fields:

```
type Trib struct {
    User    string    // who posted this trib
    Message string    // the content of the trib
    Time    time.Time // the physical timestamp
    Clock   uint64    // the logical clock
}
```

`Time` is what the front-end claims when this tribble is created, by
reading the front-end's own physical time on the machine when `Post()`
in a `trib.Server` is called.  However, to sort tribbles in a globally
consistent and *reasonable* order, we cannot sort the tribbles only by
this timestamp, because different front-ends always have different
physical time readings. For sorting, Tribbler service needs to
maintain a distributed logical `Clock` in `uint64`. 

When sorting many tribbles into a single timeline, you should sort by
the fields following this priority:

1. `Clock` The logical timestamp.
2. `Time` The physical timestamp.
3. `User` The user id.
4. `Message` The message content.

We call this the *Tribble Order*.

## Tribbler Service Interface

The Tribbler service logic is defined by the `trib.Server` interface
(in `trib/trib.go` file). This is how the Web user interface interacts
with a Tribbler server.

***

```
SignUp(user string) error
```

Creates a new user. After a user is created, it will never disappear
in the system.  

A valid user name must be no longer than `trib.MaxUsernameLen=15`
characters but not empty, must start with a lower-case letter, and can
only contain lower-case letters or numbers.  There is a helper
function called `trib.IsValidUsername(string)` which you can use to
check if a username is valid.

Returns error when the username is invalid or the user already exists.
Concurrent sign-ups might both succeed.

***

```
ListUsers() ([]string, error)
```

Lists at least `trib.MinListUser=20` different registered users. When
there are less than 20 users that have ever signed up, list all of
them. The returned usernames should be sorted in alphabetical order.

This is just for showing some users on the front page; this is not for
listing all the users that have ever signed up, because that would be
too expensive in a scalable system.

***

```
Post(who, post string, clock uint64) error
```

Posts a tribble. `clock` is the maximum clock value this user client
has ever seen so far by reading tribbles (via `Home()` and `Tribs()`).
It returns error when the user does not exist or the post is too long
(longer than `trib.MaxTribLen=140`).

***

```
Tribs(user string) ([]*Trib, error)
```

Lists the recent `trib.MaxTribFetch=100` tribbles that a user posted.
Tribbles needs to be sorted in the Tribble Order. Also, it should make
sure that the order is the same order that the user posted the
tribbles.

***

```
Follow(who, whom string) error
Unfollow(who, whom string) error
IsFollowing(who, whom string) (bool, error)
Following(who string) ([]string, error)
```

These are functions to follow/unfollow, check following and listing
all following users of a user. A user can never follow or unfollow
himself. When calling with `who` equals to `whom`, the functions
return error. When the user does not exist, the functions return
error.

A user can follow at most `trib.MaxFollowing=2000` users. Returns
error when trying to follow more than that.

***

```
Home(user string) ([]*Trib, error)
```

List the recent `trib.MaxTribFetch=100` tribbles that are posted on
the user's following timeline in Tribble Order.  In addition, the
ordering should always satisfy that:

1. If a tribble A is posted after a tribble B is posted, and they are
posted by the same user, A always shows after B.
2. If a tribble A is posted 10 seconds after a tribble B is posted,
even if they are posted by different users, A always shows after B.
3. If a tribble A is posted after a user client sees tribble B, A
always shows after B.

A is *posted after* B means B calls `Post()` after A's `Post()`
returned.

It returns error when the user does not exist.

***

In addition to normal errors, these calls might also return IO errors
if the implementation needs to communicate to a remote part.
Returning a nil error means that the call is successfully executed;
returning a non-nil error means that the call might be successfully
executed or not.

Note that `trib/ref` package contains a reference implementation for
`trib.Server` (which use tried in the lab setup).

## Entry Functions

Your job for Lab2 is two-folded.

1. Implement the bin storage service by implementing the keeper (the
`ServeKeeper()` function) and the bin storage client logic (the
`NewBinClient()` function).
2. Reimplement the Tribbler service by translating all the Tribbler
service calls into bin storage service calls (the `NewFront()`
function).

You can find these entry functions in `lab2.go` file under `triblab`
repo:

```
func NewBinClient(backs []string) trib.BinStorage
```

This function is similar to `NewClient()` in `lab1.go` but instead
returns a `trib.BinStorage` interface.  `trib.BinStorage` has only one
function called `Bin()`, which takes a bin name and returns a
`trib.Storage`. A bin storage provides another layer of mapping, where
the caller will first get a key-value storage (the bin) with a bin
name, and then perform key-value function calls on the returned
storage. Different bin names should logically return completely
separated key-value storage spaces, but note that multiple bins might
share a single physical storage underlying by appending the bin name
as a prefix (or a suffix) in the keys.

For your convenience, we have a new small package called `trib/colon`
which you can use for escaping all the colons in a string. The
escaping rule is simple, all `'|'` runes are escaped into `"||"`, and
all `':'` runes are escaped into `"|:"`.  Please use this package
wisely.

***

```
func ServeKeeper(b *trib.KeeperConfig) error
```

This function is a blocking function, similar to `ServeBack()`.  It
will spawn a keeper instance that maintains the distributed back-ends
in consistent states. For Lab2, there will be only one keeper and the
keeper do not need to do much, but in Lab3, they will be responsible
of handling all the back-end joining, leaving, crashing and related
key migrations. In Lab2, there will be only one keeper, but in Lab3,
there will be multiple keepers for fault-tolerent.

(Note: A blocking function here means that it never returns in normal
states, and whenever it does return, it always returns an error. Some
students implemented `ServeBack()` with `go http.Serve(...)` or
`return nil` in Lab1; that was not correct. The right thing should be
doing something like `return http.Serve(...)`.) 

The `trib.KeeperConfig` structure contains all the back-end serving
addresses and also some peering information for the keepers which you
would need in Lab3:

- `Backs []string` These are the addresses of the back-ends.  These
  are the back-ends that the keeper needs to maintain.
- `Keepers []string` These are the addresses that the keepers might
  listen on. Lab2 will have only one keeper; Lab3 will have least
  three keepers.
- `This int` The index of this keeper (in the `Keepers` list).  For
  Lab2, it will always be zero.
- `Id int64` A non-zero incarnation identifier for this keeper,
  usually derived from system clock. For Lab2, since the keeper will
  never crash, this fields does not matter.
- `Ready` A ready signal channel. It works in a way similar to how
  `Ready` works in `trib.BackConfig`. When a `true` is received on
  this channel, the distributed bin storage should be ready to serve.
  Therefore, if you need to initialize the physical back-ends in some
  way, make sure you do it before you send a signal over `Ready`, and
  don't forget to send a `false` to `Ready` when the initialization
  fails.

The keeper can do whatever it wants to do, but the keeper should do no
more than maintaining the bin storage in a consistent state. It is 
okay for a keeper to understand how a bin storage client translates the keys,
but should not need to parse anything further in the keys or values.
This means that with `NewBinClient()`, `ServeBack()` (implemented in
Lab1) and `ServeKeeper()` calls, they should together provide a
general distributed key-value pair bin storage layer, where it could
work for any kinds of service (including but not only Tribbler).

The keeper does not need to listen on the address given. The keeper
address here serves more like a keeper configuration name; it tells
the keeper launcher if a keeper should be running on the current
machine.

***

```
func NewFront(s trib.BinStorage) trib.Server
```

This function takes a bin storage, and returns an implementation of
`trib.Server`. The returned instance then will serve as a service
front-end that takes Tribbler service function calls, and translates
them into key-value pair bin storage calls. This front-end should be
stateless, concurrency safe, and ready to be killed at any time.  This
means that at any time during its execution on any call, the back-end
key-value pair storage always needs to stay in a consistent state.
Also, note that one front-end might be taking multiple concurrent
requests from the Web, and there might be multiple front-ends talking
to the same back-end, so make sure it handles all the concurrency
issues correctly.

Also, be aware that the `trib.BinStorage` parameter might be *not* one
that you just implemented in the previous entry functions, but as long
as it satisfies the bin storage interface specification, your Tribbler
server should work just fine.  This means that you cannot rely on the
bin storage keeper to perform Tribbler related garbage cleaning. The
front-ends might spawn back-ground routines that do the garbage
cleaning by themselves.

## Playing with It

First is to setup the bin storage. Since we might have multiple parts
running at the same time probably on different machines, we will have
a configuration file that specifies the serving addresses of the
back-ends and the keeper for the distributed bin storage.  The config
file by default will use `bins.rc` as its file name, where `bins`
stands for bin storage.

`bins.rc` is saved in json format, marshalling a `RC` structure type
(defined in `trib/rc.go` file).  We have a utility program called
`bins-mkrc` that can generate a `bins.rc` file automatically.

Find a directory as your working directory (like `triblab`).

```
$ bins-mkrc -local -nback=3
```

This will generate a file called `bins.rc` under the current
directory, and also print the file content to stdout.  `-local` means
that all addresses will be on `localhost`.  `-nback=3` means there
will be in total 3 back-ends.  If you remove `-local`, then it will
generate back-ends starting from `172.22.14.211` to `172.22.14.220`,
which are the IP address of our lab machines. For `bins-mkrc`, there
can be at most 10 backends and 10 keepers (since we only have 10 lab
machines).  However, you are free to create your own `bins.rc` file
that has more back-ends and keepers. For this case, we use 3 back-ends
and 1 keeper (the default value for `-nkeep`).

With this configuration file generated, we can now launch the
back-ends:

```
$ bins-back
```

This will read and parse the `bins.rc` file, and spawn all the
back-ends which serving address is on this host. Since all the
back-ends we generate here are on `localhost`, so all the 3 back-ends
are created for this case (in different go routines). You should see
three log lines showing that three back-ends started, but listening on
different ports. Besides that `bins-back` reads from the configuration
file, it is not much different from the `kv-serve` program.  You can
also manually specify the back-ends you would like to start from the
command line. For example, you can run the following to start only the
first two back-ends:

```
$ bins-back 0 1
```

By the way, see how this program starts several independent servers at
the same time in different go routines but in the same process, and
now you should better understand why using the `rpc.DefaultServer` is
not the right thing to do in Lab1.

After the back-ends are ready, we can now start the keeper.

```
$ bins-keeper
```

If should print a message log that shows that the bin storage is
ready to serve.

To play with this distributed bin storage, we have another toy client
program called `bins-client`:

```
$ bins-client
(working on bin "")
> bin a
(working on bin "a")
> get a

> set a b
true
> get a
b
> bin t
(working on bin "t")
> get a b

> bin a
(working on bin "a")
> get a b
b
...
```

This program reads the back-end addresses from `bins.rc` and can
switch between different bins with `bin` command. The default bin
is the bin with a name of an empty string.

Now with the bin storage working, we can finally launch our
Tribbler front-end that uses this distributed storage:

```
$ trib-front -init -addr=:rand -lab
```

You have used this utility before. The only new thing here is the
`-lab` flag, which tells it to read the `bins.rc` file and use our lab
implementation. This will start a stateless front-end (which you
implemented in this lab) that will connect to the back-ends service
addresses specified in `bins.rc`.

Again `-init` will populate the service with some sample data.

Now you can open your browser, connect to the front-end machine and
play with your own implementation.

If you want to use some other configuration file, use the `-rc` flag.
It is supported in all the utilities above.

Again, when you completes this lab, it should be perfectly fine to
have multiple front-ends that connects to the set of back-ends.
Also both the front-ends and the back-ends should be scalable.

## Assumptions

These are some unreal assumptions you can have for Lab2.

- No network communication error will happen.
- Once a back-end or a keeper starts, it will remain online forever.
- The system will always start in the following order:
  all the back-ends, the keeper, then all the front-ends.
- The `trib.Storage` used in the backend will return every `Clock()`
  call in less than 1 second.
- In the `trib.Storage` used in the backend, all storage IOs
  on a single backend are serialized (and hence provides sequential
  consistency). Each key visiting (checking if the key exist, locating
  its corresponding value, or as a process of iterating all keys) will
  take less than 1 millisecond. Read and write 1MB of data on the
  value part (in a list or a string) will take less than 1
  millisecond.  Note that `Keys()` and `ListKeys()` might take longer
  time to complete because it needs to scan over all the keys.
- All the front-ends, the back-ends and the keeper will be running
  on the lab machines.
- Although the Tribbler front-ends can be killed at any time, the
  killing won't happen very often (less than once per second).

Note that some of them won't stay true in Lab3, so try not to rely on
the assumptions too much if possible.

## Requirements

In addition to the requirements specified by the interfaces, your
implementation should also satisfy the following requirements:

- When the Tribbler service function call has valid arguments, the
  function call should not return any error.
- The front-end part should be stateless and hence ready to be killed
  at anytime.
- The back-ends should be scalable, and the front-end should use the
  back-ends in a scalable way. This means that when the back-end is
  the system throughput bottleneck, adding more back-ends should (with
  high probability) mitigate the bottleneck and lead to better overall
  system performance.
- When running on the lab machines, with more than 5 back-ends
  supporting (and assuming all the back-ends satisfies the performance
  assumptions), every Tribbler service call should return within 3
  seconds.
- Each back-end should maintain the same general key-value pair
  semantics as they were in Lab1. As a result, all test cases that
  pass for Lab1 should also pass for Lab2. This means that, the
  back-ends do not need to understand anything about the bin storage
  or the front-ends (like how the keys will be structured and
  organized, or how to parse the values).

## Building Hints

While you are free to build the system in your own way, here are some
suggested hints:

- For each service call in the front-end, if it updates anything in
  the back-end storage, use only one write-RPC call for the
  commitment. This will make sure it the call either succeed or fail.
  You might issue more write calls afterwards, but those should be
  only soft hints, where if they did not succeed, it does not leave
  the storage in an inconsistent state.
- Hash the tribbles and other information into all the back-ends based
  on username. You may find the package `hash/fnv` helpful for
  hashing.
- Synchronize the logical clocks among all the back-ends every second.
  (This will also serve as a heart-beat signal, which will be useful
  for implementing Lab3.) However, you should not try to synchronize
  the clocks for every post, because that will be not scalable.
- Do some garbage collection when one user have too many tribbles
  saved in the storage.
- Keep multiple caches for the ListUsers() call when the users are
  many. Note that when the user count is more than 20, you don't
  need to track new registered users anymore.

## Possible Mistakes

Here are some possible mistakes that a lazy and quick but incorrect
implementation might do:

- **Read-modify-write**. For example, a tribbler might read a counter
  from the key-value store, increase it by one and then write it back
  (at the same key).  This will introduce racing condition among the
  front-ends.
- **Not handling errors**. A tribbler service call might require
  several RPC calls to the backend. It is important to properly handle
  *any* error returned by these calls. It is okay to tell the webpage
  that an error occurred. However, it is often not a good idea to leave
  the back-end in inconsistent state.
- **Sorting by the timestamps first**. Again, the Tribble Order means
  that the logic clock is the first field to consider on sorting.
- **Misuse the clock argument in Post()**. For example, you
  might directly use that argument as the new post's clock field.
  Technically, you can do that in your code internally as long as you
  can satisfy the ordering requirements specified for `Home()` and
  `Tribs()` (you might find it very hard).  Nonetheless, intuitively,
  the clock argument tells the *oldest* tribble a user have seen
  (which might be 0 if the user has not seen any tribble yet), hence
  the new posted tribble seems to better have a clock value that is
  larger than the argument.
- **Generate the clock from the timestamp**. While 64-bit can cover a
  very wide time range even in the unit of nanoseconds, you should
  keep in mind that the front-ends are running on different servers
  with arbitrary physical time differences, so it is not wise to
  generate the logical *clock* from the physical *time*.
- **Not handling old tribbles**. Note that only the most recent 100
  tribbles of a user matter. Not handling old tribbles might lead to
  worse and worse performance over time and eventually break the
  performance promise.

## Turning In

First, make sure that you have committed every piece of your code into
the repository `triblab`. Then just type `make turnin` under the root
of the repository.  It will generate a `turnin.zip` that contains
everything in your git repository, and will then copy the zip file to
a place where only the lab instructors can read.

## Happy Lab2

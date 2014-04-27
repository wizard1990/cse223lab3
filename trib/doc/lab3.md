## Lab3

Welcome to Lab3. The goal of this lab is to take the bin storage that
we implemented in Lab2 and make it fault-tolerant.

Lab3 can be submitted in teams of up to 3 people.

## Get Your Repo Up-to-date

```
$ cd ~/gopath/src/trib
$ git branch lab3
$ git checkout lab3
$ git pull /classes/cse223b/sp14/labs/trib lab3
$
$ cd ~/gopath/src/triblab
$ git branch lab3
$ git checkout lab3
$ git pull /classes/cse223b/sp14/labs/triblab lab3
```

Not many changes, only some small things, should be painless.

It does not come with more unit tests (because it is not very easy to
cleanly spawn and kill processes in unit tests). You need to test by
yourself with `bins-*` tools.

## System Scale and Failure Model

There could be up to 300 back-ends. Back-ends may join and leave at
will, but at any time there will be at least 1 back-end online. Also,
you can assume that each back-end join/leave event will have a time
interval of 30 seconds in between, and this time duration will be
enough for you to migrate storage.

There will be at least 3 and up to 10 keepers. Keepers may join and
leave at will, but at any time there will be at least 1 keeper online.
Also, you can assume that each keeper join/leave event will have a
time interval of 1 minute in between. When it says "leave" here, it
assumes that the process of the back-end or the keeper is killed;
everything in that process will be lost.  Each time the keeper comes
back at the same `Index`, although all states are lost, it will get a
new `Id` field in the `KeeperConfig` structure.

For starting, we will start at least one back-end, and then at least one
keeper. After the keeper sends `true` to the `Ready` channel, a
front-end may now start and issue `BinStorage` calls.

## Consistency Model

To tolerate failures, you have to save the data of each key on
multiple places, and we will have a slightly relaxed consistency
model. 

`Clock()` and the key-value interface calles (`Set()`, `Get()` and
`Keys()`) will remain the same semantics.

When concurrent `ListAppend()` happens, when calling `ListGet()`, the
caller might see the values that are currently being added appear in
arbitrary order. However, after all the concurrent `ListAppend()`'s
successfully returned, `ListGet()` should always return the list with
a consistent order.

Here is an example of an valid call and return sequence:

- Initially, the list `"k"` is empty.
- A invokes `ListAppend("k", "a")`
- B invokes `ListAppend("k", "b")`
- C calls `ListGet("k")` and gets `["b"]`, note that how `"b"` appears
  first in the list here.
- D calls `ListGet("k")` and gets `["a", "b"]`, note that although
  `"b"` appears first in time, it appears at the second position in
  the list.
- A's `ListAppend()` call returns
- B's `ListAppend()` call returns
- C calls `ListGet("k")` again and gets `["a", "b"]`
- D calls `ListGet("k")` again and gets `["a", "b"]`

`ListRemove()` removes all matched values that are appended into
the list in the past, and sets the `n` field propoerly.
When (and only when) concurrent `ListRemove()` on the same key and 
value is called, it is okay to double count on `n`.

`ListKeys()` remains the same semantics.

## Entry Functions

The entry functions will remain exactly the same as they are in Lab2,
but only that the `KeeperConfig` might now have multiple keepers.

## Additional Assumptions

- No network error; when a TCP connection is lost (RPC client
  returning `ErrShutdown`), you can assume that the RPC server
  crashed.
- When a bin-client, back-end, or keeper is killed, all data in that
  process will be lost; nothing will be carried over a respawn.
- It will take less than 20 seconds to read all data stored on a
  back-end and write it to another back-end.

## Requirement

- If at all times, there will always be at least 3 back-ends online
  (might be different three ones at any moment in time), there should
  be no data loss.
- Key-value storage call always returns without an error, even when a
  node and/or a keeper just joined or left.

## Building Hints

- You can use the logging technique to store everything (in lists on
  the back-ends, even for values).
- You need to replicate each piece of data.
- Let the keeper(s) keep track on the status of all the nodes, and do
  the data migration when a back-end joins or leaves.
- Keepers should also keep track on the status of each other.

For the ease of debugging, you can maintain some log messages (by
using `log` package, or by writing to a TCP socket or a log file).
However, for the convenience of grading, please turn them off by
default when you turn in your code.

Also, try use a machine different than c08-11 for testing and debugging,
this will lower your probability of running into a port collision.

## Turning In

If you are submitting as a team, please create a file called
`teammates` under the root of `triblab` repo that lists the login ids
of the members of your team in each line.

Make sure that you have committed every piece of your code (and the
`teammates` file) into the repository `triblab`. Then just type 
`make turnin-lab3` under the root of your repository. It will generate a
`turnin.zip` that contains everything in your gitt repo, and will then
copy the zip file to a place where only the lab instructors can read.

## Happy Lab3. :)

# line2go

Package `line2go` provides low level interfaces to LINE Thrift protocol and servers

much wow

such line protocol


### Installation

because this is private repository so it can't be installed directly using `go get`
workaround is:

    $ cd $GOPATH/src
    $ mkdir github.com/Choestelus/ && cd github.com/Choestelus/line2go
    $ git clone github.com/Choestelus/line2go
    $ go install


update using:

    git pull

### Documentation
run `godoc` as background process is recommended

example:

    $ godoc -http=:6060 &

and connect to `http://localhost:6060/pkg/line2go`


### Wishlists

 - Ensure data race free
 - Synchronize via communicating instead of mutexes
 - More Compact thrift `TCompactProtocol` & `THttpPostClient` (need to be rewritten)
 - Proper tests

Feedbacks, Issues and pull request are welcome!

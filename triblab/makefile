.PHONY: all rall fmt tags test testv lc doc \
	turnin-lab1 turnin-lab2 turnin-lab3 turnin-zip

all:
	go install ./... trib/...

rall:
	go build -a ./... trib/...

fmt:
	gofmt -s -w -l .

tags:
	gotags `find . -name "*.go"` > tags

test:
	go test ./...

testv:
	go test -v ./...

lc:
	wc -l `find . -name "*.go"`

doc:
	godoc -http=:8000

turnin-zip:
	git archive -o turnin.zip HEAD
	chmod 600 turnin.zip

turnin-lab1: turnin-zip
	@ echo "Turning in lab1 for `whoami`"
	cp turnin.zip /classes/cse223b/sp14/labs/turnin/lab1/`whoami`.zip

turnin-lab2: turnin-zip
	@ echo "Turning in lab2 for `whoami`"
	cp turnin.zip /classes/cse223b/sp14/labs/turnin/lab2/`whoami`.zip

turnin-lab3: turnin-zip
	@ echo "Turning in lab3 for `whoami`"
	cp turnin.zip /classes/cse223b/sp14/labs/turnin/lab3/`whoami`.zip


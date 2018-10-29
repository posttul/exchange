test:
	go test -v

dep:
	@which dep 2>/dev/null || $(GO) get github.com/golang/dep/cmd/dep
	dep ensure

build: dep
	docker build -t exchange .

run: 
	docker run --rm -p 4000:4000  -v ${PWD}/data/:/go/src/app exchange
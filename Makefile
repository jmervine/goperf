VERSION=$(shell cat VERSION)
BUILDBOX=$(shell uname -a)

# tests without -tabs for go tip
travis: get .PHONY
	# Run Test Suite
	go test -test.v=true

test: format get .PHONY
	# Run Test Suite
	-go test -test.v=true
	-cd results; go test -test.v=true
	-cd connector; go test -test.v=true

build: test .PHONY
	cd bin; go build -o '../pkg/goperf-$(VERSION)' -v -a -race
	@echo "goperf $(VERSION)" > pkg/build-$(VERSION)
	@echo "" >> pkg/build-$(VERSION)
	@echo "Built on:" >> pkg/build-$(VERSION)
	@echo "---------" >> pkg/build-$(VERSION)
	@echo "$(BUILDBOX)" >> pkg/build-$(VERSION)


get:
	# Go Get Deps
	go get github.com/jmervine/GoT

docs: format .PHONY
	@godoc -ex=true | sed -e 's/func /\nfunc /g' | less
	@#                                         ^ add a little spacing for readability

readme: test
	# generating readme
	godoc -ex -v -templates "$(PWD)/docs" . > README.md

format: .PHONY
	# Gofmt Source
	gofmt -tabs=false -tabwidth=4 -w=true -l=true *.go

.PHONY:

VERSION=$(shell cat VERSION)
BUILDBOX=$(shell uname -a)

# tests without -tabs for go tip
travis: get .PHONY
	# Run Test Suite
	go test -test.v=true . ./results ./connector

test: format lint .PHONY
	go test . ./results ./connector

build: test .PHONY
	cd bin; go build -o '../_pkg/goperf-$(VERSION)' -v -a -race
	@echo "goperf $(VERSION)" > _pkg/build-$(VERSION)
	@echo "" >> _pkg/build-$(VERSION)
	@echo "Built on:" >> _pkg/build-$(VERSION)
	@echo "---------" >> _pkg/build-$(VERSION)
	@echo "$(BUILDBOX)" >> _pkg/build-$(VERSION)


get:
	# Go Get Deps
	go get github.com/jmervine/GoT

docs: format .PHONY
	@godoc -ex=true | sed -e 's/func /\nfunc /g' | less
	@#                                         ^ add a little spacing for readability

readme: test
	# generating readme
	godoc -ex -v -templates "$(PWD)/_support" . > README.md

lint: .PHONY
	# Run Linter
	-golint . || echo "WARNING: golint isn't installed, skipping linting"

format: .PHONY
	# Gofmt Source
	gofmt -tabs=false -tabwidth=4 -w=true -l=true *.go

.PHONY:

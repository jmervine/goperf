CHECK_TARGET=$(GOPATH)/src/github.com/jmervine/check
VERSION=$(shell cat VERSION)
BUILDBOX=$(shell uname -a)

# tests without -tabs for go tip
travis: get .PHONY
	# Run Test Suite
	go test -test.v=true -check.v

test: format get .PHONY
	# Run Test Suite
	-go test -test.v=true -check.v

build: test .PHONY
	cd bin; go build -o '../pkg/goperf-$(VERSION)' -v -a -race
	@echo "goperf $(VERSION)" > pkg/build-$(VERSION)
	@echo "" >> pkg/build-$(VERSION)
	@echo "Built on:" >> pkg/build-$(VERSION)
	@echo "---------" >> pkg/build-$(VERSION)
	@echo "$(BUILDBOX)" >> pkg/build-$(VERSION)


get:
	# Go Get Deps
	@test -d $(CHECK_TARGET) || \
		git clone --branch v1 https://github.com/jmervine/check.git $(CHECK_TARGET)

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

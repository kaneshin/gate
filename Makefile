GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
LINTIGNOREDEPS='vendor/.+\.go'
VETIGNOREDEPS='vendor/.+\.go'
CYCLOIGNOREDEPS='vendor/.+\.go'
TARGET_ONLY_PKGS=$(shell go list ./... 2> /dev/null | grep -v "/misc/" | grep -v "/vendor/")
INTERNAL_BIN=.bin
HAVE_GLIDE:=$(shell which glide)
HAVE_GOLINT:=$(shell which golint)
HAVE_GOCYCLO:=$(shell which gocyclo)
HAVE_GOCOV:=$(shell which gocov)
VERSION=$(patsubst "%",%,$(lastword $(shell grep 'const Version' gate.go)))
COMMITISH=$(shell git rev-parse HEAD)

.PHONY:

init: install-deps

build: install-deps

install: install-deps

unit: lint vet cyclo build test
unit-report: lint vet cyclo build test-report

lint: golint
	@echo "go lint"
	@lint=`golint ./...`; \
		lint=`echo "$$lint" | grep -E -v -e ${LINTIGNOREDEPS}`; \
		echo "$$lint"; if [ "$$lint" != "" ]; then exit 1; fi

vet:
	@echo "go vet"
	@vet=`go tool vet -all -structtags -shadow $(shell ls -d */ | grep -v "vendor") 2>&1`; \
		vet=`echo "$$vet" | grep -E -v -e ${VETIGNOREDEPS}`; \
		echo "$$vet"; if [ "$$vet" != "" ]; then exit 1; fi

cyclo: gocyclo
	@echo "gocyclo -over 20"
	@cyclo=`gocyclo -over 20 .`; \
		cyclo=`echo "$$cyclo" | grep -E -v -e ${CYCLOIGNOREDEPS}/`; \
		echo "$$cyclo"; if [ "$$cyclo" != "" ]; then exit 1; fi

test:
	@go test $(TARGET_ONLY_PKGS)

coverage: gocov
	@gocov test $(TARGET_ONLY_PKGS) | gocov report

test-report:
	@echo "Invoking test and coverage"
	@echo "" > coverage.txt
	@for d in $(TARGET_ONLY_PKGS); do \
		go test -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out || true; done

install-deps: glide
	@echo "Installing all dependencies"
	@PATH=$(INTERNAL_BIN):$(PATH) glide i

golint:
ifndef HAVE_GOLINT
	@echo "Installing linter"
	@go get -u github.com/golang/lint/golint
endif

gocyclo:
ifndef HAVE_GOCYCLO
	@echo "Installing gocyclo"
	@go get -u github.com/fzipp/gocyclo
endif

gocov:
ifndef HAVE_GOCOV
	@echo "Installing gocov"
	@go get -u github.com/axw/gocov/gocov
endif

glide:
ifndef HAVE_GLIDE
	@echo "Installing glide"
	@mkdir -p $(INTERNAL_BIN)
	@GLIDE_VERSION='v0.12.2'; \
		wget -q -O - https://github.com/Masterminds/glide/releases/download/$$GLIDE_VERSION/glide-$$GLIDE_VERSION-$(GOOS)-$(GOARCH).tar.gz | tar xvz
	@mv $(GOOS)-$(GOARCH)/glide $(INTERNAL_BIN)/glide
	@rm -rf $(GOOS)-$(GOARCH)
endif

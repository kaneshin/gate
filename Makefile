GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
TARGET_ONLY_PKGS=$(shell go list ./... 2> /dev/null | grep -v "/vendor/")
IGNORE_DEPS_GOLINT='vendor/.+\.go'
IGNORE_DEPS_GOVET='vendor/.+\.go'
IGNORE_DEPS_GOCYCLO='vendor/.+\.go'
HAVE_GOLINT:=$(shell which golint)
HAVE_GOCYCLO:=$(shell which gocyclo)
HAVE_GOCOV:=$(shell which gocov)
HAVE_GHR:=$(shell which ghr)
HAVE_GOX:=$(shell which gox)
PROJECT_REPONAME=$(notdir $(abspath ./))
PROJECT_USERNAME=$(notdir $(abspath ../))
OBJS=$(notdir $(TARGETS))
LDFLAGS=-ldflags="-s -w"
COMMITISH=$(shell git rev-parse HEAD)
ARTIFACTS_DIR=artifacts
TARGETS=$(addprefix github.com/$(PROJECT_USERNAME)/$(PROJECT_REPONAME)/cmd/,gate gatecli)
VERSION=$(patsubst "%",%,$(lastword $(shell grep 'const Version' gate.go)))

all: $(TARGETS)

$(TARGETS):
	@go install $(LDFLAGS) -v $@

.PHONY: build release clean
build: gox
	@mkdir -p $(ARTIFACTS_DIR)/$(VERSION) && cd $(ARTIFACTS_DIR)/$(VERSION); \
		gox $(LDFLAGS) $(TARGETS)

release: ghr verify-github-token build
	@ghr -c $(COMMITISH) -u $(PROJECT_USERNAME) -r $(PROJECT_REPONAME) -t $$GITHUB_TOKEN \
		--replace $(VERSION) $(ARTIFACTS_DIR)/$(VERSION)

clean:
	$(RM) -r $(ARTIFACTS_DIR)

.PHONY: unit unit-report
unit: lint vet cyclo test
unit-report: lint vet cyclo test-report

.PHONY: lint vet cyclo test coverage test-report
lint: golint
	@echo "go lint"
	@lint=`golint ./...`; \
		lint=`echo "$$lint" | grep -E -v -e ${IGNORE_DEPS_GOLINT}`; \
		echo "$$lint"; if [ "$$lint" != "" ]; then exit 1; fi

vet:
	@echo "go vet"
	@vet=`go tool vet -all -structtags -shadow $(shell ls -d */ | grep -v "vendor") 2>&1`; \
		vet=`echo "$$vet" | grep -E -v -e ${IGNORE_DEPS_GOVET}`; \
		echo "$$vet"; if [ "$$vet" != "" ]; then exit 1; fi

cyclo: gocyclo
	@echo "gocyclo -over 30"
	@cyclo=`gocyclo -over 30 . 2>&1`; \
		cyclo=`echo "$$cyclo" | grep -E -v -e ${IGNORE_DEPS_GOCYCLO}`; \
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

.PHONY: verify-github-token
verify-github-token:
	@if [ -z "$$GITHUB_TOKEN" ]; then echo '$$GITHUB_TOKEN is required'; exit 1; fi

.PHONY: golint gocyclo gocov ghr gox
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

ghr:
ifndef HAVE_GHR
	@echo "Installing ghr to upload binaries for release page"
	@go get -u github.com/tcnksm/ghr
endif

gox:
ifndef HAVE_GOX
	@echo "Installing gox to build binaries for Go cross compilation"
	@go get -u github.com/mitchellh/gox
endif


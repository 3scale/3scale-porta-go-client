PACKAGE_CLIENT = github.com/3scale/3scale-porta-go-client/client

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))

.DEFAULT_GOAL := help
.PHONY : help
help: Makefile
	@sed -n 's/^##//p' $<

## test: Run unit tests
.PHONY: test
ifdef TEST_NAME
test: TEST_PATTERN := --run $(TEST_NAME)
endif
test:
	go test -v $(PACKAGE_CLIENT) -test.coverprofile="coverage.txt" $(TEST_PATTERN)

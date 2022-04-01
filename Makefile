PACKAGE_CLIENT = github.com/3scale/3scale-porta-go-client/client

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))

# find or download gotest
# download gotest if necessary
GOTEST=$(PROJECT_PATH)/bin/gotest
$(GOTEST):
	$(call go-get-tool,$(GOTEST),github.com/rakyll/gotest)

.DEFAULT_GOAL := help
.PHONY : help
help: Makefile
	@sed -n 's/^##//p' $<

## test: Run unit tests
.PHONY: test
test: $(GOTEST)
	$(GOTEST) -v $(PACKAGE_CLIENT) -test.coverprofile="coverage.txt"

# go-get-tool will 'go get' any package $2 and install it to $1.
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_PATH)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

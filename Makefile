PACKAGE_CLIENT = github.com/3scale/3scale-porta-go-client/client

.DEFAULT_GOAL := help
.PHONY : help
help: Makefile
	@sed -n 's/^##//p' $<

## test: Run unit tests
test:
	go test -v $(PACKAGE_CLIENT) -test.coverprofile="coverage.txt"

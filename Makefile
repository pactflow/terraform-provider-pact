TEST?=./...

.DEFAULT_GOAL := ci

ci:: deps bin test integration-test

local: build clean
	terraform init && \
	TF_LOG=DEBUG TF_LOG_PATH=log/tf.log terraform apply -auto-approve

build:
	go build -o bin/terraform-provider-pact
	mkdir -p ~/.terraform.d/plugins/
	cp bin/terraform-provider-pact ~/.terraform.d/plugins/

clean:
	mkdir -p ./log && \
	touch terraform.tfstate terraform.tfstate.backup log/tf.log && \
	rm terraform.tf* log/tf.log

docker:
	docker-compose up

bin:
	gox -os="darwin" -arch="amd64" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	gox -os="windows" -arch="386" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	gox -os="linux" -arch="386" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	gox -os="linux" -arch="amd64" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	@echo "==> Results:"
	ls -hl bin/

deps:
	@echo "--- 🐿  Fetching build dependencies "
	go get github.com/axw/gocov/gocov
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/cover
	go get github.com/modocache/gover
	go get github.com/mitchellh/gox

goveralls:
	goveralls -service="travis-ci" -coverprofile=coverage.txt -repotoken $(COVERALLS_TOKEN)

test:
	@echo "--- ✅ Running tests"
	@if [ -f coverage.txt ]; then rm coverage.txt; fi;
	@echo "mode: count" > coverage.txt
	@for d in $$(go list ./... | grep -v vendor | grep -v examples); \
		do \
			go test -race -coverprofile=profile.out -covermode=atomic $$d; \
			if [ -f profile.out ]; then \
					cat profile.out | tail -n +2 >> coverage.txt; \
					rm profile.out; \
			fi; \
	done; \

	go tool cover -func coverage.txt

integration-test:
	mkdir -p ~/.terraform.d/plugins
	cp bin/terraform-provider-pact_linux_amd64 ~/.terraform.d/plugins/terraform-provider-pact
	terraform init

release:
	echo "--- 🚀 Releasing it"
	"$(CURDIR)/scripts/release.sh"

.PHONY: build clean local bin deps goveralls release integration-test docker
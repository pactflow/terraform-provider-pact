TEST?=./...

.DEFAULT_GOAL := ci
TRAVIS_COMMIT?=1
export TF_VAR_build_number=$(TRAVIS_COMMIT)

ci:: docker deps vet bin test acceptance-test

local-no-clean: build
	terraform init && \
	TF_LOG=DEBUG TF_LOG_PATH=log/tf.log terraform apply -auto-approve

local: build clean
	terraform init && \
	TF_LOG=DEBUG TF_LOG_PATH=log/tf.log terraform apply -auto-approve

local-destroy:
	terraform destroy -auto-approve

build:
	go build -o bin/terraform-provider-pact
	mkdir -p ~/.terraform.d/plugins/github.com/pactflow/pact/0.0.1/darwin_amd64
	cp bin/terraform-provider-pact ~/.terraform.d/plugins/github.com/pactflow/pact/0.0.1/darwin_amd64/

clean:
	mkdir -p ./log && \
	touch terraform.tfstate terraform.tfstate.backup log/tf.log && \
	rm terraform.tf* log/tf.log

clean-acceptance:
	mkdir -p ./acceptance/pactflow/.terraform && \
	cd ./acceptance/pactflow/ && \
	touch terraform.tfstate terraform.tfstate.backup .terraform.lock.hcl log/tf.log && \
	rm -rf terraform.tf* log/tf.log .terraform*

docker:
	docker-compose up -d

bin:
	$$(go env GOPATH)/bin/gox -os="darwin" -arch="amd64" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	$$(go env GOPATH)/bin/gox -os="windows" -arch="386" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	$$(go env GOPATH)/bin/gox -os="linux" -arch="386" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	$$(go env GOPATH)/bin/gox -os="linux" -arch="amd64" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
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

oss-acceptance-test:
	@echo "--- Running OSS acceptance tests"
	cd acceptance/oss && \
		terraform init && \
		terraform apply -auto-approve && \
		terraform destroy -auto-approve

pactflow-acceptance-test:
	@echo "--- Running Pactflow acceptance tests"
	cd acceptance/pactflow && \
		terraform init && \
		TF_LOG=DEBUG TF_LOG_PATH=log/tf.log terraform apply -auto-approve && \
		TF_LOG=DEBUG TF_LOG_PATH=log/tf.log terraform destroy -auto-approve

binary-acceptance-test:
	@echo "--- Checking binary acceptance test"
	mkdir -p ~/.terraform.d/plugins/github.com/pactflow/pact/0.0.1/linux_amd64
	cp bin/terraform-provider-pact_linux_amd64 ~/.terraform.d/plugins/github.com/pactflow/pact/0.0.1/linux_amd64/terraform-provider-pact
	terraform init

acceptance-test: binary-acceptance-test oss-acceptance-test pactflow-acceptance-test
	@echo "--- ✅ Acceptance tests complete"

release:
	@echo "--- 🚀 Releasing it"
	"$(CURDIR)/scripts/release.sh"

vet:
	@echo "--- ✅ Running go vet"
	go vet -all ./...

.PHONY: build clean local bin deps goveralls release acceptance-test docker oss-acceptance-test
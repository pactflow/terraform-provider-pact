TEST?=./...

.DEFAULT_GOAL := ci
GITHUB_RUN_ID?=1
PACT_CLI="docker run --rm -v ${PWD}:${PWD} -e PACT_BROKER_BASE_URL -e PACT_BROKER_TOKEN pactfoundation/pact-cli:latest"

export TF_VAR_build_number=$(GITHUB_RUN_ID)
export TF_VAR_api_token=$(ACCEPTANCE_PACT_BROKER_TOKEN)
export TF_VAR_broker_base_url=$(ACCEPTANCE_PACT_BROKER_BASE_URL)

ci:: clean docker deps pact-go vet bin test pact publish acceptance-test

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
	mkdir -p ./log acceptance/oss/log acceptance/pactflow/log client/pacts && \
	touch terraform.tfstate terraform.tfstate.backup log/tf.log client/pacts/pact.json && \
	rm terraform.tf* log/tf.log client/pacts/*.json

clean-acceptance:
	mkdir -p ./acceptance/pactflow/.terraform && \
	cd ./acceptance/pactflow/ && \
	touch terraform.tfstate terraform.tfstate.backup .terraform.lock.hcl log/tf.log && \
	rm -rf terraform.tf* log/tf.log .terraform*

docker:
	docker compose up -d

bin:
	$$(go env GOPATH)/bin/gox -os="darwin" -arch="arm64" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	$$(go env GOPATH)/bin/gox -os="darwin" -arch="amd64" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	$$(go env GOPATH)/bin/gox -os="windows" -arch="386" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	$$(go env GOPATH)/bin/gox -os="linux" -arch="386" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	$$(go env GOPATH)/bin/gox -os="linux" -arch="amd64" -output="bin/terraform-provider-pact_{{.OS}}_{{.Arch}}"
	@echo "==> Results:"
	ls -hl bin/

deps:
	@echo "--- 🐿  Fetching build dependencies "
	cd /tmp; \
	go install github.com/axw/gocov/gocov@latest; \
	go install github.com/mattn/goveralls@latest; \
	go install golang.org/x/tools/cmd/cover@latest; \
	go install github.com/modocache/gover@latest; \
	go install github.com/mitchellh/gox@latest; \
	cd -
	go install github.com/pact-foundation/pact-go/v2@v2.0.5;

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

pact-go:
	echo "--- 🐿 Installing Pact FFI dependencies"
	pact-go -l DEBUG install --libDir /tmp

pact: clean pact-go
	@echo "--- 🤝 Running Pact tests"
	go test -tags=consumer -count=1 -v github.com/pactflow/terraform/client/...

publish:
	@echo "--- 🤝 Transforming Broken Keys in Pact File"
	"$(PWD)/scripts/transform-broken-keys.sh"
	@echo "--- 🤝 Publishing Pact"
	"${PACT_CLI}" publish ${PWD}/client/pacts --consumer-app-version ${GITHUB_SHA} --tag ${GITHUB_BRANCH}

can-i-deploy:
	@echo "--- 🤝 Can I Deploy?"
	# @"${PACT_CLI}" broker can-i-deploy \
	#   --pacticipant "terraform-client" \
	#   --version ${GITHUB_SHA} \
	#   --to prod

oss-acceptance-test:
	@echo "--- Running OSS acceptance tests"
	cd acceptance/oss && \
		terraform init && \
		terraform apply -auto-approve && \
		terraform destroy -auto-approve

pactflow-acceptance-test:
	@echo "--- Running Pactflow acceptance tests"
	cd acceptance/pactflow && \
		mkdir -p ./log && \
		terraform init && \
		TF_LOG=DEBUG TF_LOG_PATH=log/tf.log terraform apply -auto-approve && \
		mv pactflow.tf pactflow.tf.old && mv pactflow-update.tf.next pactflow-update.tf && \
		TF_LOG=DEBUG TF_LOG_PATH=log/tf.log terraform apply  -auto-approve && \
		TF_LOG=DEBUG TF_LOG_PATH=log/tf.log terraform destroy -auto-approve && \
		mv pactflow.tf.old pactflow.tf && mv pactflow-update.tf pactflow-update.tf.next

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
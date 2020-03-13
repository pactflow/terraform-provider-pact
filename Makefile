TEST?=./...

.DEFAULT_GOAL := ci

ci:: deps bin test acceptance-test

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
	docker-compose up -d

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

oss-acceptance-test: docker
	terraform apply -auto-approve -state=acceptance/oss/terraform.tfstate acceptance/oss/
	terraform destroy -auto-approve -state=acceptance/oss/terraform.tfstate acceptance/oss/

pactflow-acceptance-test: docker
	terraform apply -auto-approve -state=acceptance/pactflow/terraform.tfstate acceptance/pactflow/
	terraform destroy -auto-approve -state=acceptance/pactflow/terraform.tfstate acceptance/pactflow/

binary-acceptance-test:
	mkdir -p ~/.terraform.d/plugins
	cp bin/terraform-provider-pact_linux_amd64 ~/.terraform.d/plugins/terraform-provider-pact
	terraform init

acceptance-test: binary-acceptance-test oss-acceptance-test pactflow-acceptance-test
	echo "--- ✅ Acceptance tests complete"

release:
	echo "--- 🚀 Releasing it"
	"$(CURDIR)/scripts/release.sh"

.PHONY: build clean local bin deps goveralls release acceptance-test docker oss-acceptance-test
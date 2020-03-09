# Test local provider
local: build clean
	terraform init && \
	TF_LOG=DEBUG TF_LOG_PATH=/tmp/tf.log terraform apply -auto-approve

build:
	go build -o terraform-provider-pact

clean:
	touch terraform.tfstate terraform.tfstate.backup /tmp/tf.log && \
	rm terraform.tf* /tmp/tf.log

docker:
	docker-compose up

.PHONY: build clean local
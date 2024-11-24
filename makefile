# setup local cluster 
setup:
	docker compose up -d
	make setup-localstack

# setup localstack
setup-localstack:
	export AWS_SHARED_CREDENTIALS_FILE=$(pwd)/.aws/credentials
	export AWS_CONFIG_FILE=$(pwd)/.aws/config

	aws --endpoint-url=http://0.0.0.0:4566 s3 mb s3://image-poster
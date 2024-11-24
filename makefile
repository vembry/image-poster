# setup local cluster 
up:
	docker compose up -d --build --remove-orphans
	make setup-localstack

# setup localstack
setup-localstack:
	export AWS_SHARED_CREDENTIALS_FILE=$(pwd)/.aws/credentials
	export AWS_CONFIG_FILE=$(pwd)/.aws/config

	aws --endpoint-url=http://0.0.0.0:4566 s3 mb s3://public-images

# tear down local cluster
down:
	docker compose -f "compose.yml" down
	(echo "y" | docker volume prune)
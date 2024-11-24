# setup localstack
setup-localstack:
	# initialize s3
	. ./.docker/localstack/init-s3.sh

# setup local cluster 
up:
	docker compose up -d --build --remove-orphans
	make setup-localstack

# tear down local cluster
down:
	docker compose -f "compose.yml" down
	(echo "y" | docker volume prune)

# run down and up in that
start:
	make down
	make up
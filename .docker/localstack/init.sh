#!/bin/bash

# force point AWS CLI to local credentials and config
export AWS_SHARED_CREDENTIALS_FILE=$(pwd)/.aws/credentials
export AWS_CONFIG_FILE=$(pwd)/.aws/config

# create the bucket
aws --endpoint-url=http://0.0.0.0:4566 s3 mb s3://public-images

# NOTE: add more buckets definition here...


# create sqs queue url 
aws --endpoint-url=http://0.0.0.0:4566 sqs create-queue --queue-name image_transformer

# NOTE: add more sqs queue-url definition here... 
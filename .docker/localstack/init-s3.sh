#!/bin/bash

# Point AWS CLI to local credentials and config
export AWS_SHARED_CREDENTIALS_FILE=$(pwd)/.aws/credentials
export AWS_CONFIG_FILE=$(pwd)/.aws/config

# Create the bucket
aws --endpoint-url=http://0.0.0.0:4566 s3 mb s3://image-poster

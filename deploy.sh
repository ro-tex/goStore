#!/bin/bash

echo "> Compiling..."
GOOS=linux go build -ldflags="-s -w" -o main

echo "> Compressing..."
zip package.zip main

echo "> Deploying..."
aws lambda update-function-code \
  --function-name goStore \
  --zip-file fileb://$PWD/package.zip \
  --publish \
  --no-dry-run \
  --profile ro-tex \
  --region eu-west-1

rm main package.zip

echo "> Done."

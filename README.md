# AWS Lambda Extension

## Cache secret server extension

The cache server is written in Golang and it uses the same code to run inside a lambda extension and to run the sidecar container. To be able to run the same extension in different lambda’s runtime, we needed a language that could generate a binary in the end, so we are able to execute the binary directly without caring about the lambda runtime language/version.

The server runs on port 8015 and receives the secret name and the refresh parameters. The refresh parameter is used to invalidate the cache.

This helps us to reduce the cost of calling Secrets Manage for every request.

The extension was based on this code: https://github.com/aws-samples/aws-lambda-extensions
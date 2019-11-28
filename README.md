# goStore

This is a simple lambda function written in Go.
It saves payloads to DynamoDb and reads them back. With time it might start doing more.
It's a learning project to let me explore writing lambda functions in Go.

This is not the best way to approach the issue, as I'm implementing almost every single thing myself instead of using
existing libraries and frameworks. The idea is to explore and to gain deeper understanding of what the challenges are
and what needs to be there in terms of code/functionality.

**At this point I am more or less retiring this project, as I've reached a point where I would really benefit from built-in
middlewares, routing, etc. In other words it's time to start exploring the libraries and frameworks. I'll do that [here](https://gitlab.com/ro-tex/golambda).**

## Notes

* Do not forget to set "Use Lambda Proxy integration" to `true` for the root method in API Gateway and then inherit from that.
Otherwise your request object will not get populated.


## Compiling for AWS Lambda:

With debug info:

```
$ GOOS=linux go build -o main
$ zip deployment.zip main
```

Without debug info:

```
$ GOOS=linux go build -ldflags="-s -w" -o main
$ zip deployment.zip main
```

Build flags:
`-a` Force rebuild
`-o` Custom output name

Linking flags:
`-w` Omit the DWARF symbol table.
`-s` Omit the symbol table and debug information.

## Deploying to AWS Lambda:

```
aws lambda update-function-code \
  --function-name YourLambdaName \
  --zip-file fileb://$PWD/deployment.zip \
  --publish \
  --no-dry-run
```

## Further compress the binary:

```
upx --brute main
```

Be advised that compressing a binary in this way:

1. Takes some time.
2. Incurs a small delay when starting the binary (~150ms, ballpark value).

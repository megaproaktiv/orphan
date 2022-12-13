# Delete orphaned AWS Lambda Log Groups

Checks whether all Log Groups "/aws/lambda/$name"  have a Lambda Function $name.

## Call

List all orphans

```bash
orphan
```

Delete all orphans

```bash
orphan --no-dry-run
```

## Install

- Download release binary
- Copy binary in path

## Run with go

```bash
run main/main.go
```

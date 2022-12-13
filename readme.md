# Delete orphaned AWS Lambda Log Groups

Checks whether all Log Groups "/aws/lambda/name"  have a Lambda Function $name.

**Use at your own risk**

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

## Example

╰─ aws logs describe-log-groups --query "logGroups[].logGroupName"
[
    "/aws/lambda/hellodockerarm",
    "/aws/lambda/xraystarter-BucketNotificationsHandler050a0587b754-z8E1DEE5a2Zy",
    "/aws/lambda/xraystarter-LogRetentionaae0aa3c5b4d4f87b02d85b201-2QZ9sjniwX2i",
    "/aws/lambda/xraystarter-go",
    "/aws/lambda/xraystarter-py",
    "/aws/lambda/xraystarter-ts"
]


╰─ orphan
/aws/lambda/hellodockerarm
/aws/lambda/xraystarter-BucketNotificationsHandler050a0587b754-z8E1DEE5a2Zy
/aws/lambda/xraystarter-LogRetentionaae0aa3c5b4d4f87b02d85b201-2QZ9sjniwX2i
/aws/lambda/xraystarter-go
/aws/lambda/xraystarter-py
/aws/lambda/xraystarter-ts


╰─ orphan --no-dry-run
Deleted: /aws/lambda/hellodockerarm
Deleted: /aws/lambda/xraystarter-BucketNotificationsHandler050a0587b754-z8E1DEE5a2Zy
Deleted: /aws/lambda/xraystarter-LogRetentionaae0aa3c5b4d4f87b02d85b201-2QZ9sjniwX2i
Deleted: /aws/lambda/xraystarter-go
Deleted: /aws/lambda/xraystarter-py
Deleted: /aws/lambda/xraystarter-ts

╰─ aws logs describe-log-groups --query "logGroups[].logGroupName"
[]

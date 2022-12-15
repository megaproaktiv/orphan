#!/bin/bash
# Create many loggroups
for i in 0 1 2 3 4 5 6 7 8 9
do
    for k in  0 1 2 3 4 5 6 7 8 9
    do
        aws logs create-log-group --log-group-name "/aws/lambda/testgroup-${i}-${k}"
    done
done

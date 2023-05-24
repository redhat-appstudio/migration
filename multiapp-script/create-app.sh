#!/bin/bash

set -e

file=$1

# Each line should have a repo URL like this: https://github.com/ernesgonzalez33/hello-secret
while read line; do
    appName=${line##*/}
    sed -e "s/appName/$appName/g" -e "s@appUrl@$line@" values.yaml > values-$appName.yaml
    helm template ../appstudio-chart -f values-$appName.yaml | oc apply -f -
done < $file
#!/usr/bin/env bash

# (C) 2020 Hagen Huebel, ITinance GmbH https://itinance.com, dedified GmbH https://dedified.io
# All rights reserved.

# fix-swagger-models
# swaggo/swag has a bug that will prevent renaming of Models from "model.Account" ino "Account"
# we are going to fix this generation with this command
# the "-i ''" is a fix for sed required on Mac to avoid the auto-creation of backup-files

if [ "$(uname)" == "Darwin" ]; then
    sed -i '' -e 's/model\.//g' docs/*
else
    sed -i -e 's/model\.//g' docs/*
fi
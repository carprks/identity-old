#!/usr/bin/env bash
bash .travis-deploy/ecs.sh -c $CLUSTER -n $APP -i "$AWS_ECR/$APP:latest" -r $AWS_DB_REGION -t 240
#!/bin/sh
 
set -e
 
ENV=dev
echo $PROFILE
mv ${ENV}-app.env app.env
./server

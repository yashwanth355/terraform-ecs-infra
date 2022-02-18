#!/bin/sh
 
set -e
 
echo $APP_ENV=dev
cp ${APP_ENV}-app.env app.env
./server

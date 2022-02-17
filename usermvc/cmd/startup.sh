#!/bin/sh
 
set -e
 
echo $APP_ENV
mv "${APP_ENV}"-app.env app.env
./server

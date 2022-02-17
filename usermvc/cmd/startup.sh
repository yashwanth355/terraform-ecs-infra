set -e
 
echo $APP_ENV  dev
mv ${APP_ENV}-app.env app.env
./server

# GeoCalc
## Run container

```shell
docker run -d -t -i -e APP_NAME='Geo calculator' \ 
-e APP_PORT='8080' \
-e APP_VERSION='v1.0' \
-e APP_DEPLOY_FLAG='green' \
-p 8080:8080
```
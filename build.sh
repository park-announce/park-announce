
docker build -t park-announce/pa-api pa-api
docker build -t park-announce/pa-service pa-service
docker build -t park-announce/pa-web pa-web
docker build -t park-announce/pa-db pa-db

docker-compose up -d

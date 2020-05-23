//Vessel-service/Makefile

echo "Build proto file"
protoc -I. --go_out=plugins=micro:. proto/consignment/consignment.proto

echo "Build locally our app"
env CGO_ENABLED=0  GOOS=linux go build -a -installsuffix cgo -o ./builds/shippy-service-consignment

echo "Building container"
docker build -t shippy-service-consignment .


echo "Running docker container"
docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 shippy-service-consignment
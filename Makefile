PROJECTNAME=$(shell basename "$(PWD)")


proto-file:
	#Build proto file
	protoc -I. --go_out=plugins=micro:. proto/consignment/consignment.proto

build:
	#Build locally our app
	env CGO_ENABLED=0  GOOS=linux go build -a -installsuffix cgo -o ./builds/$(PROJECTNAME)

	#Building container
	docker build -t $(PROJECTNAME) .
run:
	#Running docker container
	docker run -p  50051:50051 \
				--network=host \
	 		    -e DB_HOST=mongodb://127.0.0.1:27017\
			    -e MICRO_ADDRESS=:50051\
			    -e MICRO_REGISTRY=mdns\
			    $(PROJECTNAME)

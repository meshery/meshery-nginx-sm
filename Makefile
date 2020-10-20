protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

docker:
	docker build -t layer5/meshery-nginx .

docker-run:
	(docker rm -f meshery-nginx) || true
	docker run --name meshery-nginx -d \
	-p 10007:10007 \
	-e DEBUG=true \
	layer5/meshery-nginx

run:
	DEBUG=true go run main.go
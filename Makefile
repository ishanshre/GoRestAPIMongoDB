createMongoDBcontainer:
	docker run -d -p 27017:27017 --name=GoRestApiMongoDB -v mongo_data:/data/db mongo


run:
	go run cmd/api/main.go
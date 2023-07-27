createMongoDBcontainer:
	docker run -d -p 27017:27017 --name=GoRestApiMongoDB -v mongo_data:/data/db mongo

createRedisContainer:
	docker run -d --name GoRestApiMongoRedis -p 6379:6379 redis:latest

createRedisInsightContainer:
	docker run -d --name redis-insight -p 8001:8001 redislabs/redisinsight:latest

startContainer:
	docker start GoRestApiMongoDB GoRestApiMongoRedis



stopContainer:
	docker stop GoRestApiMongoDB GoRestApiMongoRedis

run:
	go run cmd/api/main.go
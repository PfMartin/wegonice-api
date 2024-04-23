#!make
include .env

db-start:
	wegonice-db start

db-connect-admin:
	docker exec -it wegonice-db mongosh admin -u ${MONGO_INITDB_ROOT_USERNAME} -p ${MONGO_INITDB_ROOT_PASSWORD}

db-create-user:
	docker exec -it wegonice-db \
	sh -c 'mongosh admin -u ${MONGO_INITDB_ROOT_USERNAME} -p ${MONGO_INITDB_ROOT_PASSWORD} \
	--eval "use $(WEGONICE_DB)" \
	--eval "db.createUser({user: \"$(WEGONICE_USER)\", pwd: \"$(WEGONICE_PWD)\", roles: [{role: \"readWrite\", db: \"$(WEGONICE_DB)\"}]})"'

db-connect:
	docker exec -it wegonice-db mongosh "mongodb://${WEGONICE_USER}:${WEGONICE_PWD}@localhost:27017/wegonice?authSource=${WEGONICE_DB}"

get-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

docs:
	swag init -o ./api/v1/docs

mock-db:
	mockgen -destination db/mock/store.go github.com/PfMartin/wegonice-api/db DBStore
#!make
include .env

project_root := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

db-start:
	wegonice-db start

db-connect-admin:
	docker exec -it wegonice-db mongosh admin -u ${MONGO_INITDB_ROOT_USERNAME} -p ${MONGO_INITDB_ROOT_PASSWORD}
	
	// use ${WEGONICE_DB}
  // db.createUser({user: '${WEGONICE_USER}', pwd: '${WEGONICE_PWD}', roles: [{role: 'readWrite', db:'wegonice'}]})

db-connect:
	docker exec -it wegonice-db mongosh ${WEGONICE_DB} -u ${WEGONICE_USER} -p ${WEGONICE_PWD}

db-create-user:
	docker exec -it wegonice-db mongosh 

dirname:
	echo ${project_root}
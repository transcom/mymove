NAME = ppp
DB_DOCKER_CONTAINER = db-dev
export GOPATH = $(CURDIR)/server
export GOBIN = $(GOPATH)/bin
export PATH := $(PATH):$(GOBIN)
export PGPASSWORD=mysecretpassword

# This target ensures that the pre-commit hook is installed and kept up to date
# if pre-commit updates.
pre-commit: .git/hooks/pre-commit
.git/hooks/pre-commit: /usr/local/bin/pre-commit
	pre-commit install

server/bin/golint:
	go get -u github.com/golang/lint/golint
golint: server/bin/golint

deps: golint pre-commit client_deps server_deps

client_deps:
	cd client && \
	yarn install
client_build: client_deps
	cd client && \
	yarn build
client_run_dev:
	cd client && \
	yarn start
client_run: client_run_dev

glide_update:
	cd server/src/dp3 && glide update
server_deps:
	go get github.com/Masterminds/glide
	cd server/src/dp3 && glide install
	go get github.com/markbates/pop/soda
	go install github.com/markbates/pop/soda
server_build_only:
	cd server/src/dp3/cmd/webserver && \
	go install
server_build: server_deps server_build_only
server_run_only: db_dev_run
	./server/bin/webserver \
		-entry client/build/index.html \
		-build client/build \
		-port :8080 \
		-debug_logging
server_run: server_build client_build server_run_only
server_run_dev: server_build_only server_run_only
server_test: db_dev_run
	# Initialize a test database if we're not in a CircleCI environment.
	[ -z "$(CIRCLECI)" ] && \
		dropdb -p 5432 -h localhost -U postgres --if-exists test_db && \
		createdb -p 5432 -h localhost -U postgres test_db || \
		echo "Relying on CircleCI's test database setup."
	DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
		bin/wait-for-db
	DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
		go test -v dp3/pkg/api

db_dev_init:
	docker run --name $(DB_DOCKER_CONTAINER) \
		-e \
		POSTGRES_PASSWORD=$(PGPASSWORD) \
		-d \
		-p 5432:5432 \
		postgres:latest
	bin/wait-for-db
	createdb -p 5432 -h localhost -U postgres dev_db
db_dev_run:
	# We don't want to utilize Docker to start the database if we're
	# in the CircleCI environment. It has its own configuration to launch
	# a DB.
	[ -z "$(CIRCLECI)" ] && \
		docker start $(DB_DOCKER_CONTAINER) || \
		echo "Relying on CircleCI's database container."
db_dev_reset:
	echo "Attempting to reset local dev database..."
	docker kill $(DB_DOCKER_CONTAINER) &&	\
		docker rm $(DB_DOCKER_CONTAINER) || \
		echo "No dev database"
db_dev_migrate: db_dev_run
	echo "TODO: make some database migrations"

.PHONY: pre-commit deps db_dev_migrate

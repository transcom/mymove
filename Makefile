NAME = ppp
DB_DOCKER_CONTAINER = db-dev
export PGPASSWORD=mysecretpassword

# This target ensures that the pre-commit hook is installed and kept up to date
# if pre-commit updates.
pre-commit: .git/hooks/pre-commit
.git/hooks/pre-commit: /usr/local/bin/pre-commit
	pre-commit install

deps: pre-commit client_deps server_deps
test: client_test server_test

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
client_test:
	cd client && \
	yarn test

server_deps_update:
	dep ensure -update
server_deps: .server_deps.stamp
.server_deps.stamp: Gopkg.lock
	dep ensure
	go install ./vendor/github.com/markbates/pop/soda
	go install ./vendor/github.com/golang/lint/golint
	touch .server_deps.stamp
server_build: server_deps
	go build -o build/mymove ./server/cmd/mymove
server_run_only: db_dev_run
	./build/mymove \
		-entry client/build/index.html \
		-build client/build \
		-port :8080 \
		-debug_logging
server_run: client_build server_build server_run_only
server_run_dev: server_build server_run_only

server_build_docker:
	docker build . -t ppp:dev
server_run_only_docker: db_dev_run
	docker run --name ppp -p 8080:8080 ppp:dev

server_test: db_dev_run db_test_reset
	DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
		cd server/ && go test ./...

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
	cd server/ && \
	soda migrate up
db_dev_migrate_down: db_dev_run
	cd server/ && \
	soda migrate down

db_test_reset:
	# Initialize a test database if we're not in a CircleCI environment.
	[ -z "$(CIRCLECI)" ] && \
		dropdb -p 5432 -h localhost -U postgres --if-exists test_db && \
		createdb -p 5432 -h localhost -U postgres test_db || \
		echo "Relying on CircleCI's test database setup."
	DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
		bin/wait-for-db
	cd server/ && \
	soda -e test migrate up

.PHONY: pre-commit deps db_dev_migrate

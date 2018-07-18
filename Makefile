NAME = ppp
DB_DOCKER_CONTAINER = db-dev
export PGPASSWORD=mysecretpassword

# if S3 access is enabled, wrap webserver in aws-vault command
# to pass temporary AWS credentials to the binary.
ifeq ($(STORAGE_BACKEND),s3)
	USE_AWS:=true
endif
ifeq ($(EMAIL_BACKEND),ses)
	USE_AWS:=true
endif

ifeq ($(USE_AWS),true)
  AWS_VAULT:=aws-vault exec $(AWS_PROFILE) --
endif


# This target ensures that the pre-commit hook is installed and kept up to date
# if pre-commit updates.
ensure_pre_commit: .git/hooks/pre-commit
.git/hooks/pre-commit: /usr/local/bin/pre-commit
	pre-commit install

prereqs: .prereqs.stamp
.prereqs.stamp: bin/prereqs
	bin/prereqs
	touch .prereqs.stamp

go_version: .go_version.stamp
.go_version.stamp: bin/check_go_version
	bin/check_go_version
	touch .go_version.stamp

deps: prereqs ensure_pre_commit client_deps server_deps
test: client_test server_test e2e_test

spellcheck:
	node_modules/.bin/mdspell --ignore-numbers --ignore-acronyms --en-us \
		`find . -type f -name "*.md" \
			-not -path "./vendor/*" \
			-not -path "./node_modules/*" \
			-not -path "./docs/adr/index.md"`

client_deps_update:
	yarn upgrade
client_deps: .client_deps.stamp
.client_deps.stamp: yarn.lock
	yarn install
	bin/copy_swagger_ui.sh
	touch .client_deps.stamp
.client_build.stamp: $(shell find src -type f)
	yarn build
	touch .client_build.stamp
client_build: client_deps .client_build.stamp
client_run: client_deps
	yarn start
client_test: client_deps
	yarn test
client_test_coverage : client_deps
	yarn test:coverage

office_client_run: client_deps
	HOST=officelocal yarn start

tsp_client_run: client_deps
	HOST=tsplocal yarn start

server_deps_update: server_generate
	dep ensure -v -update
server_deps: go_version .server_deps.stamp
.server_deps.stamp: Gopkg.lock
	bin/check_gopath.sh
	dep ensure -vendor-only
	# Unfortunately, dep ensure blows away ./vendor every time so these builds always take a while
	go install ./vendor/github.com/golang/lint/golint # golint needs to be accessible for the pre-commit task to run, so `install` it
	go build -i -o bin/gas ./vendor/github.com/GoASTScanner/gas/cmd/gas
	go build -i -o bin/gin ./vendor/github.com/codegangsta/gin
	go build -i -o bin/soda ./vendor/github.com/gobuffalo/pop/soda
	go build -i -o bin/swagger ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
	touch .server_deps.stamp
server_generate: server_deps .server_generate.stamp
.server_generate.stamp: $(shell find swagger -type f -name *.yaml)
	bin/gen_server.sh
	touch .server_generate.stamp
server_build: server_deps server_generate
	go build -i -o bin/webserver ./cmd/webserver
# This command is for running the server by itself, it will serve the compiled frontend on its own
server_run_standalone: client_build server_build db_dev_run
	DEBUG_LOGGING=true $(AWS_VAULT) ./bin/webserver
# This command runs the server behind gin, a hot-reload server
server_run: server_deps server_generate db_dev_run
	INTERFACE=localhost DEBUG_LOGGING=true \
	$(AWS_VAULT) ./bin/gin --build ./cmd/webserver \
		--bin /bin/webserver \
		--port 8080 --appPort 8081 \
		--excludeDir vendor --excludeDir node_modules \
		-i --buildArgs "-i"
# This is just an alais for backwards compatibility
server_run_dev: server_run

server_build_docker:
	docker build . -t ppp:web-dev
server_run_only_docker: db_dev_run
	docker stop web || true
	docker rm web || true
	docker run --name web -p 8080:8080 ppp:web-dev

tools_build: server_deps
	go build -i -o bin/tsp-award-queue ./cmd/tsp_award_queue
	go build -i -o bin/generate-test-data ./cmd/generate_test_data
	go build -i -o bin/rateengine ./cmd/demo/rateengine.go
	go build -i -o bin/make-office-user ./cmd/make_office_user
	go build -i -o bin/load-office-data ./cmd/load_office_data
	go build -i -o bin/make-tsp-user ./cmd/make_tsp_user
	go build -i -o bin/load-user-gen ./cmd/load_user_gen
	go build -i -o bin/paperwork ./cmd/paperwork

tsp_run: tools_build db_dev_run
	./bin/tsp-award-queue

tsp_build_docker:
	docker build . -f Dockerfile.tsp -t ppp:tsp-dev
tsp_run_only_docker: db_dev_run
	docker stop tsp || true
	docker rm tsp || true
	docker run --name tsp ppp:tsp-dev

build: server_build tools_build client_build

server_test: server_deps server_generate db_dev_run db_test_reset
	# Don't run tests in /cmd or /pkg/gen & pass `-short` to exclude long running tests
	# Use -test.parallel 1 to test packages serially and avoid database collisions
	# Disable test caching with `-count 1` - caching was masking local test failures
	go test -p 1 -count 1 -short $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

server_test_all: server_deps server_generate db_dev_run db_test_reset
	# Like server_test but runs extended tests that may hit external services.
	go test -p 1 -count 1 $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

server_test_coverage: server_deps server_generate db_dev_run db_test_reset
	# Don't run tests in /cmd or /pkg/gen
	# Use -test.parallel 1 to test packages serially and avoid database collisions
	# Disable test caching with `-count 1` - caching was masking local test failures
	# Add coverage tracker via go cover
	# Then open coverage tracker in HTML
	go test -coverprofile=coverage.out -p 1 -count 1 $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)
	go tool cover -html=coverage.out

e2e_test: server_deps server_generate client_build db_e2e_init
	$(AWS_VAULT) ./bin/run-e2e-test

db_dev_run:
	# The version of the postgres container should match production as closely
	# as possible.
	# https://github.com/transcom/ppp-infra/blob/1578df6e6bc6bb45d43fdc7762228afdd17a4144/modules/aws-app-environment/database/main.tf#L87
	docker start $(DB_DOCKER_CONTAINER) || \
		(docker run --name $(DB_DOCKER_CONTAINER) \
			-e \
			POSTGRES_PASSWORD=$(PGPASSWORD) \
			-d \
			-p 5432:5432 \
			postgres:10.1 && \
		DB_NAME=postgres bin/wait-for-db && \
		createdb -p 5432 -h localhost -U postgres dev_db)
# This is just an alias for backwards compatibility
db_dev_init: db_dev_run
db_dev_reset:
	echo "Attempting to reset local dev database..."
	docker kill $(DB_DOCKER_CONTAINER) &&	\
		docker rm $(DB_DOCKER_CONTAINER) || \
		echo "No dev database"
db_dev_migrate: server_deps db_dev_run
	# We need to move to the bin/ directory so that the cwd contains `apply-secure-migration.sh`
	cd bin && \
		./soda -c ../config/database.yml -p ../migrations migrate up
db_dev_migrate_down: server_deps db_dev_run
	# We need to move to the bin/ directory so that the cwd contains `apply-secure-migration.sh`
	cd bin && \
		./soda -c ../config/database.yml -p ../migrations migrate down
db_build_docker:
	docker build -f Dockerfile.migrations -t ppp-migrations:dev .

db_e2e_init: tools_build db_dev_run db_test_reset
	DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
		./bin/soda -e test migrate -c config/database.yml -p cypress/migrations up

db_e2e_reset: db_dev_run
	DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
		./bin/soda -e test migrate -c config/database.yml -p cypress/migrations reset

db_test_reset:
	# Initialize a test database if we're not in a CircleCI environment.
ifndef CIRCLECI
	dropdb -p 5432 -h localhost -U postgres --if-exists test_db
	createdb -p 5432 -h localhost -U postgres test_db
else
	echo "Relying on CircleCI's test database setup."
endif
	DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
		bin/wait-for-db
	# We need to move to the bin/ directory so that the cwd contains `apply-secure-migration.sh`
	cd bin && \
		DB_HOST=localhost DB_PORT=5432 DB_NAME=test_db \
			./soda -e test migrate -c ../config/database.yml -p ../migrations up

adr_update:
	yarn run adr-log

pre_commit_tests:
	pre-commit run --all-files

clean:
	rm .*.stamp || true
	rm -rf ./node_modules
	rm -rf ./vendor
	rm -rf ./pkg/gen
	rm -rf $$GOPATH/pkg/dep/sources

.PHONY: pre-commit deps test client_deps client_build client_run client_test prereqs
.PHONY: server_deps_update server_generate server_deps server_build server_run_standalone server_run server_run_dev server_build_docker server_run_only_docker server_test
.PHONY: db_dev_init db_dev_run db_dev_reset db_dev_migrate db_dev_migrate_down db_test_reset
.PHONY: clean

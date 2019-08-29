NAME = ppp
DB_NAME_DEV = dev_db
DB_NAME_DEPLOYED_MIGRATIONS = deployed_migrations
DB_NAME_TEST = test_db
DB_DOCKER_CONTAINER_DEV = milmove-db-dev
DB_DOCKER_CONTAINER_DEPLOYED_MIGRATIONS = milmove-db-deployed-migrations
DB_DOCKER_CONTAINER_TEST = milmove-db-test
# The version of the postgres container should match production as closely
# as possible.
# https://github.com/transcom/ppp-infra/blob/7ba2e1086ab1b2a0d4f917b407890817327ffb3d/modules/aws-app-environment/database/variables.tf#L48
DB_DOCKER_CONTAINER_IMAGE = postgres:10.9
TASKS_DOCKER_CONTAINER = tasks
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

# Convenience for LDFLAGS
WEBSERVER_LDFLAGS=-X main.gitBranch=$(shell git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(shell git rev-list -1 HEAD)
GC_FLAGS=-trimpath=$(GOPATH)
DB_PORT_DEV=5432
DB_PORT_DEPLOYED_MIGRATIONS=5434
DB_PORT_DOCKER=5432
ifdef CIRCLECI
	DB_PORT_TEST=5432
	LDFLAGS=-linkmode external -extldflags -static
endif

ifdef GOLAND
	GOLAND_GC_FLAGS=all=-N -l
endif

.PHONY: help
help:  ## Print the help documentation
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#
# ----- END PREAMBLE -----
#

#
# ----- START CHECK TARGETS -----
#

# This target ensures that the pre-commit hook is installed and kept up to date
# if pre-commit updates.
.PHONY: ensure_pre_commit
ensure_pre_commit: .git/hooks/pre-commit ## Ensure pre-commit is installed
.git/hooks/pre-commit: /usr/local/bin/pre-commit
	pre-commit install
	pre-commit install-hooks

.PHONY: prereqs
prereqs: .prereqs.stamp ## Check that pre-requirements are installed
.prereqs.stamp: scripts/prereqs
	scripts/prereqs
	touch .prereqs.stamp

.PHONY: check_hosts
check_hosts: .check_hosts.stamp ## Check that hosts are in the /etc/hosts file
.check_hosts.stamp: scripts/check-hosts-file
ifndef CIRCLECI
	scripts/check-hosts-file
else
	@echo "Not checking hosts on CircleCI."
endif
	touch .check_hosts.stamp

.PHONY: check_go_version
check_go_version: .check_go_version.stamp ## Check that the correct Golang version is installed
.check_go_version.stamp: scripts/check-go-version
	scripts/check-go-version
	touch .check_go_version.stamp

.PHONY: check_gopath
check_gopath: .check_gopath.stamp ## Check that $GOPATH exists in $PATH
.check_gopath.stamp:
	scripts/check-gopath
	touch .check_gopath.stamp

.PHONY: check_bash_version
check_bash_version: .check_bash_version.stamp ## Check that the correct Bash version is installed
.check_bash_version.stamp: scripts/check-bash-version
ifndef CIRCLECI
	scripts/check-bash-version
else
	@echo "No need to check bash version on CircleCI"
endif
	touch .check_bash_version.stamp

.PHONY: deps
deps: prereqs check_hosts check_go_version check_gopath check_bash_version ensure_pre_commit client_deps server_deps ## Run all checks and install all depdendencies

.PHONY: test
test: client_test server_test e2e_test ## Run all tests

#
# ----- END CHECK TARGETS -----
#

#
# ----- START CLIENT TARGETS -----
#

.PHONY: client_deps_update
client_deps_update: ## Update client dependencies
	yarn upgrade

.PHONY: client_deps
client_deps: .check_hosts.stamp .client_deps.stamp ## Install client dependencies
.client_deps.stamp: yarn.lock
	yarn install
	scripts/copy-swagger-ui
	touch .client_deps.stamp

.client_build.stamp: $(shell find src -type f)
	yarn build
	touch .client_build.stamp

.PHONY: client_build
client_build: .client_deps.stamp .client_build.stamp ## Build the client

build/index.html: ## milmove serve requires this file to boot, but it isn't used during local development
	mkdir -p build
	touch build/index.html

.PHONY: client_run
client_run: .client_deps.stamp ## Run MilMove Service Member client
	HOST=milmovelocal yarn start

.PHONY: client_test
client_test: .client_deps.stamp ## Run client unit tests
	yarn test

.PHONY: client_test_coverage
client_test_coverage : .client_deps.stamp ## Run client unit test coverage
	yarn test:coverage

.PHONY: office_client_run
office_client_run: .client_deps.stamp ## Run MilMove Office client
	HOST=officelocal yarn start

.PHONY: tsp_client_run
tsp_client_run: .client_deps.stamp ## Run MilMove TSP client
	HOST=tsplocal yarn start

.PHONY: admin_client_run
admin_client_run: .client_deps.stamp ## Run MilMove Admin client
	HOST=adminlocal yarn start

#
# ----- END CLIENT TARGETS -----
#

#
# ----- START BIN TARGETS -----
#

### Go Tool Targets

bin/callgraph: .check_go_version.stamp .check_gopath.stamp
	go build -o bin/callgraph golang.org/x/tools/cmd/callgraph

bin/chamber: .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/chamber github.com/segmentio/chamber

bin/gin: .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/gin github.com/codegangsta/gin

bin/soda: .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/soda github.com/gobuffalo/pop/soda

bin/swagger: .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/swagger github.com/go-swagger/go-swagger/cmd/swagger

# No static linking / $(LDFLAGS) because mockery is only used for testing
bin/mockery: .check_go_version.stamp .check_gopath.stamp
	go build -o bin/mockery github.com/vektra/mockery/cmd/mockery

### Cert Targets

bin/rds-combined-ca-bundle.pem:
	mkdir -p bin/
	curl -sSo bin/rds-combined-ca-bundle.pem https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem

### MilMove Targets

# server_deps and server_generate required for this binary, because go build expects
# github.com/transcom/mymove/pkg/gen/internalmessages, even though it is not used for this program.
bin/compare-secure-migrations: server_deps server_generate
	go build -ldflags "$(LDFLAGS)" -o bin/compare-secure-migrations ./cmd/compare-secure-migrations

bin/ecs-deploy-task-container: server_deps server_generate
	go build -ldflags "$(LDFLAGS)" -o bin/ecs-deploy-task-container ./cmd/ecs-deploy-task-container

bin/ecs-service-logs:
	go build -ldflags "$(LDFLAGS)" -o bin/ecs-service-logs ./cmd/ecs-service-logs

bin/generate-1203-form: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/generate-1203-form ./cmd/generate_1203_form

bin/generate-access-codes: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/generate-access-codes ./cmd/generate_access_codes

bin/generate-shipment-edi: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/generate-shipment-edi ./cmd/generate_shipment_edi

bin/generate-test-data: pkg/assets/assets.go .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/generate-test-data ./cmd/generate-test-data

bin_linux/generate-test-data: pkg/assets/assets.go .server_generate_linux.stamp
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin_linux/generate-test-data ./cmd/generate-test-data

bin/health-checker:
	go build -ldflags "$(LDFLAGS)" -o bin/health-checker ./cmd/health-checker

bin/iws:
	go build -ldflags "$(LDFLAGS)" -o bin/iws ./cmd/iws/iws.go

bin/load-office-data: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/load-office-data ./cmd/load_office_data

bin/load-user-gen: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/load-user-gen ./cmd/load_user_gen

bin/make-dps-user: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/make-dps-user ./cmd/make_dps_user

bin/make-office-user: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/make-office-user ./cmd/make_office_user

bin/milmove: .server_generate.stamp
	go build -gcflags="$(GOLAND_GC_FLAGS) $(GC_FLAGS)" -asmflags=-trimpath=$(GOPATH) -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/milmove ./cmd/milmove

bin_linux/milmove: .server_generate_linux.stamp
	GOOS=linux GOARCH=amd64 go build -gcflags="$(GOLAND_GC_FLAGS) $(GC_FLAGS)" -asmflags=-trimpath=$(GOPATH) -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin_linux/milmove ./cmd/milmove

bin/renderer:
	# do not build with LDFLAGS since errors on alpine and dynamic linking is fine
	# throws errors loadinternal: cannot find runtime/cgo
	go build -o bin/renderer ./cmd/renderer

bin/save-fuel-price-data: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/save-fuel-price-data ./cmd/save_fuel_price_data

bin_linux/save-fuel-price-data: .server_generate_linux.stamp
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin_linux/save-fuel-price-data ./cmd/save_fuel_price_data

bin/send-to-gex: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/send-to-gex ./cmd/send_to_gex

bin/tsp-award-queue: .server_generate.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/tsp-award-queue ./cmd/tsp_award_queue

pkg/assets/assets.go: .check_go_version.stamp .check_gopath.stamp
	go-bindata -o pkg/assets/assets.go -pkg assets pkg/paperwork/formtemplates/

#
# ----- END BIN TARGETS -----
#

#
# ----- START SERVER TARGETS -----
#

.PHONY: go_deps_update
go_deps_update: server_deps server_generate mocks_generate ## Update golang dependencies
	go run cmd/update_deps/main.go

.PHONY: server_deps
server_deps: .check_gopath.stamp bin/callgraph bin/chamber bin/gin bin/swagger bin/mockery bin/rds-combined-ca-bundle.pem ## Install or Build server dependencies

.PHONY: server_generate
server_generate: .check_go_version.stamp .check_gopath.stamp .server_generate.stamp ## Generate golang server code from Swagger files
.server_generate.stamp: pkg/assets/assets.go bin/swagger $(shell find swagger -type f -name *.yaml)
	scripts/gen-server
	touch .server_generate.stamp

.PHONY: server_generate_linux
server_generate_linux: .check_go_version.stamp .check_gopath.stamp pkg/assets/assets.go bin/swagger .server_generate_linux.stamp ## Generate golang server code from Swagger files (linux)
.server_generate_linux.stamp: pkg/assets/assets.go bin/swagger $(shell find swagger -type f -name *.yaml)
	scripts/gen-server
	touch .server_generate_linux.stamp

.PHONY: server_build
server_build: server_deps server_generate bin/milmove ## Build the server

.PHONY: server_build_linux
server_build_linux: server_generate_linux ## Build the server (linux)
	# These don't need to go in bin_linux/ because local devs don't use them
	# Additionally it would not work with the default Dockerfile
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin_linux/chamber github.com/segmentio/chamber
	GOOS=linux GOARCH=amd64 go build -gcflags="$(GOLAND_GC_FLAGS) $(GC_FLAGS)" -asmflags=-trimpath=$(GOPATH) -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/milmove ./cmd/milmove

# This command is for running the server by itself, it will serve the compiled frontend on its own
# Note: Don't double wrap with aws-vault because the pkg/cli/vault.go will handle it
server_run_standalone: server_build client_build db_dev_run
	DEBUG_LOGGING=true ./bin/milmove serve

# This command will rebuild the swagger go code and rerun server on any changes
server_run:
	find ./swagger -type f -name "*.yaml" | entr -c -r make server_run_default
# This command runs the server behind gin, a hot-reload server
# Note: Gin is not being used as a proxy so assigning odd port and laddr to keep in IPv4 space.
# Note: The INTERFACE envar is set to configure the gin build, milmove_gin, local IP4 space with default port 8080.
server_run_default: .check_hosts.stamp .check_go_version.stamp .check_gopath.stamp bin/gin build/index.html server_generate db_dev_run
	INTERFACE=localhost DEBUG_LOGGING=true \
	$(AWS_VAULT) ./bin/gin \
		--build ./cmd/milmove \
		--bin /bin/milmove_gin \
		--laddr 127.0.0.1 --port 9001 \
		--excludeDir node_modules \
		--immediate \
		--buildArgs "-i -ldflags=\"$(WEBSERVER_LDFLAGS)\"" \
		serve

.PHONY: server_run_debug
server_run_debug: ## Debug the server
	$(AWS_VAULT) dlv debug cmd/milmove/main.go cmd/milmove/logger.go -- serve

.PHONY: build_tools
build_tools: server_deps \
	bin/compare-secure-migrations \
	bin/ecs-deploy-task-container \
	bin/ecs-service-logs \
	bin/generate-1203-form \
	bin/generate-access-codes \
	bin/generate-shipment-edi \
	bin/generate-test-data \
	bin/health-checker \
	bin/iws \
	bin/load-office-data \
	bin/load-user-gen \
	bin/make-dps-user \
	bin/make-office-user \
	bin/renderer \
	bin/save-fuel-price-data \
	bin/send-to-gex \
	bin/tsp-award-queue ## Build all tools

.PHONY: build
build: server_build build_tools client_build ## Build the server, tools, and client

# webserver_test runs a few acceptance tests against a local or remote environment.
# This can help identify potential errors before deploying a container.
.PHONY: webserver_test
webserver_test: bin/rds-combined-ca-bundle.pem server_generate mocks_generate bin/chamber  ## Run acceptance tests
ifndef TEST_ACC_ENV
	@echo "Running acceptance tests for webserver using local environment."
	@echo "* Use environment XYZ by setting environment variable to TEST_ACC_ENV=XYZ."
	TEST_ACC_HONEYCOMB=0 \
	TEST_ACC_CWD=$(PWD) \
	go test -v -p 1 -count 1 -short $$(go list ./... | grep \\/cmd\\/milmove)
else
ifndef CIRCLECI
	@echo "Running acceptance tests for webserver with environment $$TEST_ACC_ENV."
	TEST_ACC_HONEYCOMB=0 \
	TEST_ACC_CWD=$(PWD) \
	DISABLE_AWS_VAULT_WRAPPER=1 \
	aws-vault exec $(AWS_PROFILE) -- \
	bin/chamber -r $(CHAMBER_RETRIES) exec app-$(TEST_ACC_ENV) -- \
	go test -v -p 1 -count 1 -short $$(go list ./... | grep \\/cmd\\/milmove)
else
	@echo "Running acceptance tests for webserver with environment $$TEST_ACC_ENV."
	TEST_ACC_HONEYCOMB=0 \
	TEST_ACC_CWD=$(PWD) \
	bin/chamber -r $(CHAMBER_RETRIES) exec app-$(TEST_ACC_ENV) -- \
	go test -v -p 1 -count 1 -short $$(go list ./... | grep \\/cmd\\/milmove)
endif
endif

.PHONY: mocks_generate
mocks_generate: .mocks_generate.stamp ## Generate mockery mocks for tests
.mocks_generate.stamp : bin/mockery
	go generate $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)
	touch .mocks_generate.stamp

.PHONY: server_test
server_test: server_deps server_generate mocks_generate db_test_reset db_test_migrate ## Run server unit tests
	# Don't run tests in /cmd or /pkg/gen & pass `-short` to exclude long running tests
	# Disable test caching with `-count 1` - caching was masking local test failures
	DB_PORT=$(DB_PORT_TEST) go test -count 1 -short $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

server_test_build:
	# Try to compile tests, but don't run them.
	go test -run=nope -count 1 $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

.PHONY: server_test_all
server_test_all: server_deps server_generate mocks_generate db_dev_reset db_dev_migrate ## Run all server unit tests
	# Like server_test but runs extended tests that may hit external services.
	DB_PORT=$(DB_PORT_TEST) go test -p 1 -count 1 $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

.PHONY: server_test_coverage_generate
server_test_coverage_generate: server_deps server_generate mocks_generate db_test_reset db_test_migrate ## Run server unit test coverage
	# Don't run tests in /cmd or /pkg/gen
	# Use -test.parallel 1 to test packages serially and avoid database collisions
	# Disable test caching with `-count 1` - caching was masking local test failures
	# Add coverage tracker via go cover
	DB_PORT=$(DB_PORT_TEST) go test -coverprofile=coverage.out -covermode=count -p 1 -count 1 -short $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

.PHONY: server_test_coverage
server_test_coverage: server_deps server_generate mocks_generate db_test_reset db_test_migrate server_test_coverage_generate ## Run server unit test coverage with html output
	DB_PORT=$(DB_PORT_TEST) go tool cover -html=coverage.out

#
# ----- END SERVER TARGETS -----
#

#
# ----- START DB_DEV TARGETS -----
#

.PHONY: db_dev_destroy
db_dev_destroy: ## Destroy Dev DB
ifndef CIRCLECI
	@echo "Destroying the ${DB_DOCKER_CONTAINER_DEV} docker database container..."
	docker rm -f $(DB_DOCKER_CONTAINER_DEV) || echo "No database container"
	rm -fr mnt/db_dev # delete mount directory if exists
else
	@echo "Relying on CircleCI's database setup to destroy the DB."
endif

.PHONY: db_dev_start
db_dev_start: ## Start Dev DB
ifndef CIRCLECI
	brew services stop postgresql 2> /dev/null || true
endif
	@echo "Starting the ${DB_DOCKER_CONTAINER_DEV} docker database container..."
	# If running do nothing, if not running try to start, if can't start then run
	docker start $(DB_DOCKER_CONTAINER_DEV) || \
		docker run -d --name $(DB_DOCKER_CONTAINER_DEV) \
			-e POSTGRES_PASSWORD=$(PGPASSWORD) \
			-p $(DB_PORT_DEV):$(DB_PORT_DOCKER)\
			$(DB_DOCKER_CONTAINER_IMAGE)

.PHONY: db_dev_create
db_dev_create: ## Create Dev DB
	@echo "Create the ${DB_NAME_DEV} database..."
	DB_NAME=postgres scripts/wait-for-db && DB_NAME=postgres psql-wrapper "CREATE DATABASE $(DB_NAME_DEV);" || true

.PHONY: db_dev_run
db_dev_run: db_dev_start db_dev_create ## Run Dev DB (start and create)

.PHONY: db_dev_reset
db_dev_reset: db_dev_destroy db_dev_run ## Reset Dev DB (destroy and run)

.PHONY: db_dev_migrate_standalone ## Migrate Dev DB directly
db_dev_migrate_standalone: bin/milmove
	@echo "Migrating the ${DB_NAME_DEV} database..."
	bin/milmove migrate -p "file://migrations;file://local_migrations" -m migrations_manifest.txt

.PHONY: db_dev_migrate
db_dev_migrate: server_deps db_dev_migrate_standalone ## Migrate Dev DB

.PHONY: db_dev_psql
db_dev_psql: ## Open PostgreSQL shell for Dev DB
	scripts/psql-dev

#
# ----- END DB_DEV TARGETS -----
#

#
# ----- START DB_DEPLOYED_MIGRATIONS TARGETS -----
#

.PHONY: db_deployed_migrations_destroy
db_deployed_migrations_destroy: ## Destroy Deployed Migrations DB
ifndef CIRCLECI
	@echo "Destroying the ${DB_DOCKER_CONTAINER_DEPLOYED_MIGRATIONS} docker database container..."
	docker rm -f $(DB_DOCKER_CONTAINER_DEPLOYED_MIGRATIONS) || echo "No database container"
	rm -fr mnt/db_deployed_migrations # delete mount directory if exists
else
	@echo "Relying on CircleCI's database setup to destroy the DB."
endif

.PHONY: db_deployed_migrations_start
db_deployed_migrations_start: ## Start Deployed Migrations DB
ifndef CIRCLECI
	brew services stop postgresql 2> /dev/null || true
endif
	@echo "Starting the ${DB_DOCKER_CONTAINER_DEPLOYED_MIGRATIONS} docker database container..."
	# If running do nothing, if not running try to start, if can't start then run
	docker start $(DB_DOCKER_CONTAINER_DEPLOYED_MIGRATIONS) || \
		docker run -d --name $(DB_DOCKER_CONTAINER_DEPLOYED_MIGRATIONS) \
			-e POSTGRES_PASSWORD=$(PGPASSWORD) \
			-p $(DB_PORT_DEPLOYED_MIGRATIONS):$(DB_PORT_DOCKER)\
			$(DB_DOCKER_CONTAINER_IMAGE)

.PHONY: db_deployed_migrations_create
db_deployed_migrations_create: ## Create Deployed Migrations DB
	@echo "Create the ${DB_NAME_DEPLOYED_MIGRATIONS} database..."
	DB_NAME=postgres DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) scripts/wait-for-db && \
		createdb -p $(DB_PORT_DEPLOYED_MIGRATIONS) -h localhost -U postgres $(DB_NAME_DEPLOYED_MIGRATIONS) || true

.PHONY: db_deployed_migrations_run
db_deployed_migrations_run: db_deployed_migrations_start db_deployed_migrations_create ## Run Deployed Migrations DB (start and create)

.PHONY: db_deployed_migrations_reset
db_deployed_migrations_reset: db_deployed_migrations_destroy db_deployed_migrations_run ## Reset Deployed Migrations DB (destroy and run)

.PHONY: db_deployed_migrations_migrate_standalone
db_deployed_migrations_migrate_standalone: bin/milmove ## Migrate Deployed Migrations DB with local migrations
	@echo "Migrating the ${DB_NAME_DEPLOYED_MIGRATIONS} database..."
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) bin/milmove migrate -p "file://migrations;file://local_migrations" -m migrations_manifest.txt

.PHONY: db_deployed_migrations_migrate
db_deployed_migrations_migrate: server_deps db_deployed_migrations_migrate_standalone ## Migrate Deployed Migrations DB

.PHONY: db_deployed_psql
db_deployed_psql: ## Open PostgreSQL shell for Deployed Migrations DB
	scripts/psql-deployed-migrations

#
# ----- END DB_DEPLOYED_MIGRATIONS TARGETS -----
#

#
# ----- START DB_TEST TARGETS -----
#

.PHONY: db_test_destroy
db_test_destroy: ## Destroy Test DB
ifndef CIRCLECI
	@echo "Destroying the ${DB_DOCKER_CONTAINER_TEST} docker database container..."
	docker rm -f $(DB_DOCKER_CONTAINER_TEST) || \
		echo "No database container"
else
	@echo "Relying on CircleCI's database setup to destroy the DB."
endif

.PHONY: db_test_start
db_test_start: ## Start Test DB
ifndef CIRCLECI
	brew services stop postgresql 2> /dev/null || true
endif
	@echo "Starting the ${DB_DOCKER_CONTAINER_TEST} docker database container..."
	docker start $(DB_DOCKER_CONTAINER_TEST) || \
		docker run --name $(DB_DOCKER_CONTAINER_TEST) \
			-e \
			POSTGRES_PASSWORD=$(PGPASSWORD) \
			-d \
			-p $(DB_PORT_TEST):$(DB_PORT_DOCKER)\
			$(DB_DOCKER_CONTAINER_IMAGE)\
			-c fsync=off\
			-c full_page_writes=off

.PHONY: db_test_create
db_test_create: ## Create Test DB
ifndef CIRCLECI
	@echo "Create the ${DB_NAME_TEST} database..."
	DB_NAME=postgres DB_PORT=$(DB_PORT_TEST) scripts/wait-for-db && \
		createdb -p $(DB_PORT_TEST) -h localhost -U postgres $(DB_NAME_TEST) || true
else
	@echo "Relying on CircleCI's database setup to create the DB."
endif

.PHONY: db_test_run
db_test_run: db_test_start db_test_create ## Run Test DB

.PHONY: db_test_reset
db_test_reset: db_test_destroy db_test_run ## Reset Test DB (destroy and run)

.PHONY: db_test_migrate_standalone
db_test_migrate_standalone: bin/milmove ## Migrate Test DB directly
ifndef CIRCLECI
	@echo "Migrating the ${DB_NAME_TEST} database..."
	DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_TEST) bin/milmove migrate -p "file://migrations;file://local_migrations" -m migrations_manifest.txt
else
	@echo "Migrating the ${DB_NAME_TEST} database..."
	DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_DEV) bin/milmove migrate -p "file://migrations;file://local_migrations" -m migrations_manifest.txt
endif

.PHONY: db_test_migrate
db_test_migrate: server_deps db_test_migrate_standalone ## Migrate Test DB

.PHONY: db_test_migrations_build
db_test_migrations_build: .db_test_migrations_build.stamp ## Build Test DB Migrations Docker Image
.db_test_migrations_build.stamp: server_generate_linux bin_linux/milmove bin_linux/generate-test-data
	@echo "Build the docker migration container..."
	docker build -f Dockerfile.migrations_local --tag e2e_migrations:latest .

.PHONY: db_test_psql
db_test_psql: ## Open PostgreSQL shell for Test DB
	scripts/psql-test

#
# ----- END DB_TEST TARGETS -----
#

#
# ----- START E2E TARGETS -----
#

.PHONY: e2e_test
e2e_test: bin/gin server_generate server_build client_build db_e2e_init ## Run e2e (end-to-end) integration tests
	$(AWS_VAULT) ./scripts/run-e2e-test

.PHONY: e2e_test_docker
e2e_test_docker: ## Run e2e (end-to-end) integration tests with docker
	$(AWS_VAULT) ./scripts/run-e2e-test-docker

.PHONY: e2e_test_docker_mymove
e2e_test_docker_mymove: ## Run e2e (end-to-end) Service Member integration tests with docker
	$(AWS_VAULT) SPEC=cypress/integration/mymove/**/* ./scripts/run-e2e-test-docker

.PHONY: e2e_test_docker_office
e2e_test_docker_office: ## Run e2e (end-to-end) Office integration tests with docker
	$(AWS_VAULT) SPEC=cypress/integration/office/**/* ./scripts/run-e2e-test-docker

.PHONY: e2e_test_docker_tsp
e2e_test_docker_tsp: ## Run e2e (end-to-end) TSP integration tests with docker
	$(AWS_VAULT) SPEC=cypress/integration/tsp/**/* ./scripts/run-e2e-test-docker

.PHONY: e2e_test_docker_api
e2e_test_docker_api: ## Run e2e (end-to-end) API integration tests with docker
	$(AWS_VAULT) SPEC=cypress/integration/api/**/* ./scripts/run-e2e-test-docker

.PHONY: e2e_clean
e2e_clean: ## Clean e2e (end-to-end) files and docker images
	rm -f .*_linux.stamp
	rm -rf cypress/results
	rm -rf cypress/screenshots
	rm -rf cypress/videos
	rm -rf bin_linux/
	docker rm -f cypress || true

.PHONY: db_e2e_up
db_e2e_up: bin/generate-test-data ## Truncate Test DB and Generate e2e (end-to-end) data
	@echo "Truncate the ${DB_NAME_TEST} database..."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_TEST)/$(DB_NAME_TEST)?sslmode=disable -c 'TRUNCATE users CASCADE;'
	@echo "Populate the ${DB_NAME_TEST} database..."
	DB_PORT=$(DB_PORT_TEST) bin/generate-test-data --named-scenario="e2e_basic" --db-env="test"

.PHONY: db_e2e_init
db_e2e_init: db_test_reset db_test_migrate db_e2e_up ## Initialize e2e (end-to-end) DB (reset, migrate, up)

.PHONY: db_dev_e2e_populate
db_dev_e2e_populate: db_dev_reset db_dev_migrate build_tools ## Populate Dev DB with generated e2e (end-to-end) data
	@echo "Populate the ${DB_NAME_DEV} database with docker command..."
	bin/generate-test-data --named-scenario="e2e_basic" --db-env="development"

.PHONY: db_test_e2e_populate
db_test_e2e_populate: db_test_reset db_test_migrate build_tools db_e2e_up ## Populate Test DB with generated e2e (end-to-end) data

.PHONY: db_test_e2e_backup
db_test_e2e_backup: ## Backup Test DB as 'e2e_test'
	DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_TEST) ./scripts/db-backup e2e_test

.PHONY: db_test_e2e_restore
db_test_e2e_restore: ## Restore Test DB from 'e2e_test'
	DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_TEST) ./scripts/db-restore e2e_test

.PHONY: db_test_e2e_cleanup
db_test_e2e_cleanup: ## Clean up Test DB backup `e2e_test`
	./scripts/db-cleanup e2e_test


#
# ----- END E2E TARGETS -----
#

#
# ----- START SCHEDULED TASK TARGETS -----
#

.PHONY: tasks_clean
tasks_clean: ## Clean Scheduled Task files and docker images
	rm -f .db_test_migrations_build.stamp
	rm -rf bin_linux/
	docker rm -f tasks || true

.PHONY: tasks_build
tasks_build: server_generate bin/save-fuel-price-data ## Build Scheduled Task dependencies

.PHONY: tasks_build_docker
tasks_build_docker: bin/chamber server_generate bin/save-fuel-price-data ## Build Scheduled Task dependencies and Docker image
	@echo "Build the docker scheduled tasks container..."
	docker build -f Dockerfile.tasks --tag $(TASKS_DOCKER_CONTAINER):latest .

.PHONY: tasks_build_linux_docker
tasks_build_linux_docker: bin_linux/save-fuel-price-data ## Build Scheduled Task binaries (linux) and Docker image (local)
	@echo "Build the docker scheduled tasks container..."
	docker build -f Dockerfile.tasks_local --tag $(TASKS_DOCKER_CONTAINER):latest .

.PHONY: tasks_save_fuel_price_data
tasks_save_fuel_price_data: tasks_build_linux_docker ## Run save-fuel-price-data from inside docker container
	@echo "Saving the fuel price data to the ${DB_NAME_DEV} database with docker command..."
	DB_NAME=$(DB_NAME_DEV) DB_DOCKER_CONTAINER=$(DB_DOCKER_CONTAINER_DEV) scripts/wait-for-db-docker
	docker run \
		-t \
		-e DB_HOST="database" \
		-e DB_NAME \
		-e DB_PORT \
		-e DB_USER \
		-e DB_PASSWORD \
		-e EIA_KEY \
		-e EIA_URL \
		--link="$(DB_DOCKER_CONTAINER_DEV):database" \
		--rm \
		$(TASKS_DOCKER_CONTAINER):latest \
		save-fuel-price-data

#
# ----- END SCHEDULED TASK TARGETS -----
#

#
# ----- START Deployed MIGRATION TARGETS -----
#

.PHONY: run_prod_migrations
run_prod_migrations: server_deps bin/milmove db_deployed_migrations_reset ## Run Prod migrations against Deployed Migrations DB
	@echo "Migrating the prod-migrations database with prod migrations..."
	MIGRATION_PATH="s3://transcom-ppp-app-prod-us-west-2/secure-migrations;file://migrations" \
	DB_HOST=localhost \
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) \
	DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) \
	bin/milmove migrate

.PHONY: run_staging_migrations
run_staging_migrations: server_deps bin/milmove db_deployed_migrations_reset ## Run Staging migrations against Deployed Migrations DB
	@echo "Migrating the prod-migrations database with staging migrations..."
	MIGRATION_PATH="s3://transcom-ppp-app-staging-us-west-2/secure-migrations;file://migrations" \
	DB_HOST=localhost \
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) \
	DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) \
	bin/milmove migrate

.PHONY: run_experimental_migrations
run_experimental_migrations: server_deps bin/milmove db_deployed_migrations_reset ## Run Experimental migrations against Deployed Migrations DB
	@echo "Migrating the prod-migrations database with experimental migrations..."
	MIGRATION_PATH="s3://transcom-ppp-app-experimental-us-west-2/secure-migrations;file://migrations" \
	DB_HOST=localhost \
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) \
	DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) \
	bin/milmove migrate

#
# ----- END PROD_MIGRATION TARGETS -----
#

#
# ----- START DEPENDENCY UPDATE TARGETS -----
#

.PHONY: dependency_update
dependency_update: go_deps_update client_deps_update ## Update golang and client dependencies
	git --no-pager status
	git --no-pager diff --ignore-all-space --color

.PHONY: dependency_update_test
dependency_update_test: ## Test dependency updater
	docker build . -f Dockerfile.dep_updater

#
# ----- END DEPENDENCY UPDATE TARGETS -----
#

#
# ----- START MAKE TEST TARGETS -----
#

.PHONY: make_test
make_test: ## Test make targets not checked by CircleCI
	scripts/make-test

#
# ----- END MAKE TEST TARGETS -----
#

#
# ----- START RANDOM TARGETS -----
#

.PHONY: adr_update
adr_update: .client_deps.stamp ## Update ADR Log
	yarn run adr-log

.PHONY: gofmt
gofmt:  ## Run go fmt over all Go files
	go fmt $$(go list ./...) >> /dev/null

.PHONY: pre_commit_tests
pre_commit_tests: .server_generate.stamp .mocks_generate.stamp .client_deps.stamp ## Run pre-commit tests
	pre-commit run --all-files

.PHONY: pretty
pretty: gofmt ## Run code through JS and Golang formatters
	npx prettier --write --loglevel warn "src/**/*.{js,jsx}"

.PHONY: prune_images
prune_images:  ## Prune docker images
	@echo '****************'
	docker image prune -a

.PHONY: prune_containers
prune_containers:  ## Prune docker containers
	@echo '****************'
	docker container prune

.PHONY: prune_volumes
prune_volumes:  ## Prune docker volumes
	@echo '****************'
	docker volume prune

.PHONY: prune
prune: prune_images prune_containers prune_volumes ## Prune docker containers, images, and volumes

.PHONY: clean
clean: ## Clean all generated files
	rm -f .*.stamp
	rm -f coverage.out
	rm -rf ./bin
	rm -rf ./bin_linux
	rm -rf ./build
	rm -rf ./node_modules
	rm -rf ./pkg/gen
	rm -f ./pkg/assets/assets.go
	rm -rf ./public/swagger-ui/*.{css,js,png}
	rm -rf ./tmp/secure_migrations
	rm -rf ./tmp/storage
	find ./pkg -type d -name "mocks" -exec rm -rf {} +

.PHONY: spellcheck
spellcheck: .client_deps.stamp ## Run interactive spellchecker
	node_modules/.bin/mdspell --ignore-numbers --ignore-acronyms --en-us \
		`find . -type f -name "*.md" \
			-not -path "./node_modules/*" \
			-not -path "./vendor/*" \
			-not -path "./docs/adr/index.md"`

.PHONY: storybook
storybook: ## Start the storybook server
	yarn run storybook

.PHONY: build_storybook
build_storybook: ## Build static storybook site
	yarn run build-storybook

#
# ----- END RANDOM TARGETS -----
#

#
# ----- START DOCKER COMPOSE TARGETS -----
#

.PHONY: docker_compose_setup
docker_compose_setup: .check_hosts.stamp ## Install requirements to use docker-compose
	brew install -f bash git docker docker-compose direnv || true
	brew cask install -f aws-vault || true

.PHONY: docker_compose_up
docker_compose_up: ## Bring up docker-compose containers
	aws ecr get-login --no-include-email --region us-west-2 --no-include-email | sh
	scripts/update-docker-compose
	docker-compose up

.PHONY: docker_compose_down
docker_compose_down: ## Destroy docker-compose containers
	docker-compose down
	# Instead of using `--rmi all` which might destroy postgres we just remove the AWS containers
	docker rmi $(shell docker images --filter=reference='*amazonaws*/*:*' --format "{{.ID}}")
	git checkout docker-compose.yml

#
# ----- END DOCKER COMPOSE TARGETS -----
#

#
# ----- START ANTI VIRUS TARGETS -----
#

.PHONY: anti_virus
anti_virus: ## Scan repo with anti-virus service
	scripts/anti-virus

#
# ----- END ANTI VIRUS TARGETS -----
#

default: help

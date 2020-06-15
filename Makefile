DB_NAME_DEV = dev_db
DB_NAME_DEPLOYED_MIGRATIONS = deployed_migrations
DB_NAME_TEST = test_db
DB_DOCKER_CONTAINER_DEV = milmove-db-dev
DB_DOCKER_CONTAINER_DEPLOYED_MIGRATIONS = milmove-db-deployed-migrations
DB_DOCKER_CONTAINER_TEST = milmove-db-test
# The version of the postgres container should match production as closely
# as possible.
# https://github.com/transcom/transcom-infrasec-com/blob/c32c45078f29ea6fd58b0c246f994dbea91be372/transcom-com-legacy/app-prod/main.tf#L62
DB_DOCKER_CONTAINER_IMAGE = postgres:12.4
REDIS_DOCKER_CONTAINER_IMAGE = redis:5.0.6
REDIS_DOCKER_CONTAINER = milmove-redis
TASKS_DOCKER_CONTAINER = tasks
WEBHOOK_CLIENT_DOCKER_CONTAINER = webhook-client
export PGPASSWORD=mysecretpassword

# if S3 access is enabled, wrap webserver in aws-vault command
# to pass temporary AWS credentials to the binary.
ifeq ($(STORAGE_BACKEND),s3)
	USE_AWS:=true
endif
ifeq ($(USE_AWS),true)
  AWS_VAULT:=aws-vault exec $(AWS_PROFILE) --
endif

# Convenience for LDFLAGS
GIT_BRANCH ?= $(shell git branch | grep \* | cut -d ' ' -f2)
GIT_COMMIT ?= $(shell git rev-list -1 HEAD)
export GIT_BRANCH GIT_COMMIT
WEBSERVER_LDFLAGS=-X main.gitBranch=$(GIT_BRANCH) -X main.gitCommit=$(GIT_COMMIT)
GC_FLAGS=-trimpath=$(GOPATH)
DB_PORT_DEV=5432
DB_PORT_TEST=5433
DB_PORT_DEPLOYED_MIGRATIONS=5434
DB_PORT_DOCKER=5432
REDIS_PORT=6379
REDIS_PORT_DOCKER=6379
ifdef CIRCLECI
	DB_PORT_DEV=5432
	DB_PORT_TEST=5432
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		LDFLAGS=-linkmode external -extldflags -static
	endif
endif

ifdef GOLAND
	GOLAND_GC_FLAGS=all=-N -l
endif

SCHEMASPY_OUTPUT=./tmp/schemaspy

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

.PHONY: check_node_version
check_node_version: .check_node_version.stamp ## Check that the correct Node version is installed
.check_node_version.stamp: scripts/check-node-version
	scripts/check-node-version
	touch .check_node_version.stamp

.PHONY: check_docker_size
check_docker_size: ## Check the amount of disk space used by docker
	scripts/check-docker-size

.PHONY: deps
deps: prereqs ensure_pre_commit client_deps redis_pull bin/rds-ca-2019-root.pem bin/rds-ca-us-gov-west-1-2017-root.pem ## Run all checks and install all depdendencies

.PHONY: test
test: client_test server_test e2e_test ## Run all tests

.PHONY: diagnostic
diagnostic: .prereqs.stamp check_docker_size ## Run diagnostic scripts on environment

.PHONY: check_log_dir
check_log_dir: ## Make sure we have a log directory
	mkdir -p log

#
# ----- END CHECK TARGETS -----
#

#
# ----- START CLIENT TARGETS -----
#

.PHONY: client_deps_update
client_deps_update: .check_node_version.stamp ## Update client dependencies
	yarn upgrade

.PHONY: client_deps
client_deps: .check_hosts.stamp .check_node_version.stamp .client_deps.stamp ## Install client dependencies
.client_deps.stamp: yarn.lock
	yarn install
	scripts/copy-swagger-ui
	touch .client_deps.stamp

.client_build.stamp: .check_node_version.stamp $(shell find src -type f)
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

bin/gin: .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/gin github.com/codegangsta/gin

bin/soda: .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/soda github.com/gobuffalo/pop/v5/soda

# No static linking / $(LDFLAGS) because go-junit-report is only used for building the CirlceCi test report
bin/go-junit-report: .check_go_version.stamp .check_gopath.stamp
	go build -o bin/go-junit-report github.com/jstemmer/go-junit-report

# No static linking / $(LDFLAGS) because mockery is only used for testing
bin/mockery: .check_go_version.stamp .check_gopath.stamp
	go build -o bin/mockery github.com/vektra/mockery/cmd/mockery

### Cert Targets

bin/rds-ca-2019-root.pem:
	mkdir -p bin/
	curl -sSo bin/rds-ca-2019-root.pem https://s3.amazonaws.com/rds-downloads/rds-ca-2019-root.pem

bin/rds-ca-us-gov-west-1-2017-root.pem:
	mkdir -p bin/
	curl -sSo bin/rds-ca-us-gov-west-1-2017-root.pem https://s3.us-gov-west-1.amazonaws.com/rds-downloads/rds-ca-us-gov-west-1-2017-root.pem

### MilMove Targets

bin/big-cat: cmd/big-cat
	go build -ldflags "$(LDFLAGS)" -o bin/big-cat ./cmd/big-cat

bin/model-vet: cmd/model-vet
	go build -ldflags "$(LDFLAGS)" -o bin/model-vet ./cmd/model-vet

bin/generate-deploy-notes: cmd/generate-deploy-notes
	go build -ldflags "$(LDFLAGS)" -o bin/generate-deploy-notes ./cmd/generate-deploy-notes

bin/ecs-deploy: cmd/ecs-deploy
	go build -ldflags "$(LDFLAGS)" -o bin/ecs-deploy ./cmd/ecs-deploy

bin/generate-access-codes: cmd/generate_access_codes
	go build -ldflags "$(LDFLAGS)" -o bin/generate-access-codes ./cmd/generate_access_codes

bin/generate-shipment-summary: cmd/generate_shipment_summary
	go build -ldflags "$(LDFLAGS)" -o bin/generate-shipment-summary ./cmd/generate_shipment_summary

bin/generate-test-data: cmd/generate-test-data
	go build -ldflags "$(LDFLAGS)" -o bin/generate-test-data ./cmd/generate-test-data

bin/ghc-pricing-parser: cmd/ghc-pricing-parser
	go build -ldflags "$(LDFLAGS)" -o bin/ghc-pricing-parser ./cmd/ghc-pricing-parser

bin/ghc-transit-time-parser: cmd/ghc-transit-time-parser
	go build -ldflags "$(LDFLAGS)" -o bin/ghc-transit-time-parser ./cmd/ghc-transit-time-parser

bin/health-checker: cmd/health-checker
	go build -ldflags "$(LDFLAGS)" -o bin/health-checker ./cmd/health-checker

bin/iws: cmd/iws
	go build -ldflags "$(LDFLAGS)" -o bin/iws ./cmd/iws/iws.go

bin/milmove: cmd/milmove
	go build -gcflags="$(GOLAND_GC_FLAGS) $(GC_FLAGS)" -asmflags=-trimpath=$(GOPATH) -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/milmove ./cmd/milmove

bin/milmove-tasks: cmd/milmove-tasks
	go build -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/milmove-tasks ./cmd/milmove-tasks

bin/prime-api-client: cmd/prime-api-client
	go build -ldflags "$(LDFLAGS)" -o bin/prime-api-client ./cmd/prime-api-client

bin/webhook-client: cmd/webhook-client
	go build -ldflags "$(LDFLAGS)" -o bin/webhook-client ./cmd/webhook-client

bin/read-alb-logs: cmd/read-alb-logs
	go build -ldflags "$(LDFLAGS)" -o bin/read-alb-logs ./cmd/read-alb-logs

bin/report-ecs: cmd/report-ecs
	go build -ldflags "$(LDFLAGS)" -o bin/report-ecs ./cmd/report-ecs

bin/send-to-gex: pkg/gen/ cmd/send-to-gex
	go build -ldflags "$(LDFLAGS)" -o bin/send-to-gex ./cmd/send-to-gex

bin/send-to-syncada-via-sftp: pkg/gen/ cmd/send-to-syncada-via-sftp
	go build -ldflags "$(LDFLAGS)" -o bin/send-to-syncada-via-sftp ./cmd/send-to-syncada-via-sftp

bin/fetch-from-syncada-via-sftp: pkg/gen/ cmd/fetch-from-syncada-via-sftp
	go build -ldflags "$(LDFLAGS)" -o bin/fetch-from-syncada-via-sftp ./cmd/fetch-from-syncada-via-sftp

bin/tls-checker: cmd/tls-checker
	go build -ldflags "$(LDFLAGS)" -o bin/tls-checker ./cmd/tls-checker

bin/generate-payment-request-edi: cmd/generate-payment-request-edi
	go build -ldflags "$(LDFLAGS)" -o bin/generate-payment-request-edi ./cmd/generate-payment-request-edi

pkg/assets/assets.go:
	scripts/gen-assets

#
# ----- END BIN TARGETS -----
#

#
# ----- START SERVER TARGETS -----
#

.PHONY: server_generate
server_generate: .check_go_version.stamp .check_gopath.stamp pkg/gen/ ## Generate golang server code from Swagger files
pkg/gen/: pkg/assets/assets.go $(shell find swagger -type f -name *.yaml)
	scripts/gen-server

.PHONY: server_build
server_build: bin/milmove ## Build the server

# This command is for running the server by itself, it will serve the compiled frontend on its own
# Note: Don't double wrap with aws-vault because the pkg/cli/vault.go will handle it
server_run_standalone: check_log_dir server_build client_build db_dev_run redis_run
	./bin/milmove serve 2>&1 | tee -a log/dev.log

# This command will rebuild the swagger go code and rerun server on any changes
server_run:
	find ./swagger -type f -name "*.yaml" | entr -c -r make server_run_default
# This command runs the server behind gin, a hot-reload server
# Note: Gin is not being used as a proxy so assigning odd port and laddr to keep in IPv4 space.
# Note: The INTERFACE envar is set to configure the gin build, milmove_gin, local IP4 space with default port GIN_PORT.
server_run_default: .check_hosts.stamp .check_go_version.stamp .check_gopath.stamp .check_node_version.stamp check_log_dir bin/gin build/index.html server_generate db_dev_run redis_run
	INTERFACE=localhost \
		./bin/gin \
		--build ./cmd/milmove \
		--bin /bin/milmove_gin \
		--laddr 127.0.0.1 --port "$(GIN_PORT)" \
		--excludeDir node_modules \
		--immediate \
		--buildArgs "-i -ldflags=\"$(WEBSERVER_LDFLAGS)\"" \
		serve \
		2>&1 | tee -a log/dev.log

.PHONY: server_run_debug
server_run_debug: .check_hosts.stamp .check_go_version.stamp .check_gopath.stamp .check_node_version.stamp check_log_dir build/index.html server_generate db_dev_run redis_run ## Debug the server
	scripts/kill-process-on-port 8080
	scripts/kill-process-on-port 9443
	DISABLE_AWS_VAULT_WRAPPER=1 \
	AWS_REGION=us-gov-west-1 \
	aws-vault exec transcom-gov-dev -- \
	dlv debug cmd/milmove/*.go -- serve 2>&1 | tee -a log/dev.log

.PHONY: build_tools
build_tools: bin/gin \
	bin/mockery \
	bin/rds-ca-2019-root.pem \
	bin/rds-ca-us-gov-west-1-2017-root.pem \
	bin/big-cat \
	bin/generate-deploy-notes \
	bin/ecs-deploy \
	bin/generate-access-codes \
	bin/generate-test-data \
	bin/ghc-pricing-parser \
	bin/ghc-transit-time-parser \
	bin/health-checker \
	bin/iws \
	bin/milmove-tasks \
	bin/model-vet \
	bin/prime-api-client \
	bin/read-alb-logs \
	bin/report-ecs \
	bin/send-to-gex \
	bin/tls-checker ## Build all tools

.PHONY: build
build: server_build build_tools client_build ## Build the server, tools, and client

# acceptance_test runs a few acceptance tests against a local or remote environment.
# This can help identify potential errors before deploying a container.
.PHONY: acceptance_test
acceptance_test: bin/rds-ca-2019-root.pem bin/rds-ca-us-gov-west-1-2017-root.pem ## Run acceptance tests
ifndef TEST_ACC_ENV
	@echo "Running acceptance tests for webserver using local environment."
	@echo "* Use environment XYZ by setting environment variable to TEST_ACC_ENV=XYZ."
	TEST_ACC_CWD=$(PWD) \
	SERVE_ADMIN=true \
	SERVE_SDDC=true \
	SERVE_ORDERS=true \
	SERVE_DPS=true \
	SERVE_API_INTERNAL=true \
	SERVE_API_GHC=true \
	MUTUAL_TLS_ENABLED=true \
	go test -v -count 1 -short $$(go list ./... | grep \\/cmd\\/milmove)
else
ifndef CIRCLECI
	@echo "Running acceptance tests for webserver with environment $$TEST_ACC_ENV."
	TEST_ACC_CWD=$(PWD) \
	DISABLE_AWS_VAULT_WRAPPER=1 \
	aws-vault exec $(AWS_PROFILE) -- \
	chamber -r $(CHAMBER_RETRIES) exec app-$(TEST_ACC_ENV) -- \
	go test -v -count 1 -short $$(go list ./... | grep \\/cmd\\/milmove)
else
	go build -ldflags "$(LDFLAGS)" -o bin/chamber github.com/segmentio/chamber/v2
	@echo "Running acceptance tests for webserver with environment $$TEST_ACC_ENV."
	TEST_ACC_CWD=$(PWD) \
	bin/chamber -r $(CHAMBER_RETRIES) exec app-$(TEST_ACC_ENV) -- \
	go test -v -count 1 -short $$(go list ./... | grep \\/cmd\\/milmove)
endif
endif

.PHONY: mocks_generate
mocks_generate: bin/mockery ## Generate mockery mocks for tests
	go generate $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

.PHONY: server_test
server_test: db_test_reset db_test_migrate redis_reset server_test_standalone ## Run server unit tests

.PHONY: server_test_standalone
server_test_standalone: ## Run server unit tests with no deps
	NO_DB=1 scripts/run-server-test

.PHONY: server_test_build
server_test_build:
	NO_DB=1 DRY_RUN=1 scripts/run-server-test

.PHONY: server_test_all
server_test_all: db_dev_reset db_dev_migrate redis_reset ## Run all server unit tests
	# Like server_test but runs extended tests that may hit external services.
	LONG_TEST=1 scripts/run-server-test

.PHONY: server_test_coverage_generate
server_test_coverage_generate: db_test_reset db_test_migrate redis_reset server_test_coverage_generate_standalone ## Run server unit test coverage

.PHONY: server_test_coverage_generate_standalone
server_test_coverage_generate_standalone: ## Run server unit tests with coverage and no deps
	# Add coverage tracker via go cover
	NO_DB=1 SERVER_REPORT=1 COVERAGE=1 scripts/run-server-test

.PHONY: server_test_coverage
server_test_coverage: db_test_reset db_test_migrate redis_reset server_test_coverage_generate ## Run server unit test coverage with html output
	DB_PORT=$(DB_PORT_TEST) go tool cover -html=coverage.out

#
# ----- END SERVER TARGETS -----
#

#
# ----- START REDIS TARGETS -----
#

.PHONY: redis_pull
redis_pull: ## Pull redis image
	docker pull $(REDIS_DOCKER_CONTAINER_IMAGE)

.PHONY: redis_destroy
redis_destroy: ## Destroy Redis
	@echo "Destroying the ${REDIS_DOCKER_CONTAINER} docker redis container..."
	docker rm -f $(REDIS_DOCKER_CONTAINER) || echo "No Redis container"

.PHONY: redis_run
redis_run: ## Run Redis
ifndef CIRCLECI
		@echo "Stopping the Redis brew service in case it's running..."
		brew services stop redis 2> /dev/null || true
endif
	@echo "Starting the ${REDIS_DOCKER_CONTAINER} docker redis container..."
	docker start $(REDIS_DOCKER_CONTAINER) || \
		docker run -d --name $(REDIS_DOCKER_CONTAINER) \
			-p $(REDIS_PORT):$(REDIS_PORT_DOCKER) \
			$(REDIS_DOCKER_CONTAINER_IMAGE)

.PHONY: redis_reset
redis_reset: redis_destroy redis_run ## Reset Redis

#
# ----- END REDIS TARGETS -----
#

#
# ----- START DB_DEV TARGETS -----
#

.PHONY: db_pull
db_pull: ## Pull db image
	docker pull $(DB_DOCKER_CONTAINER_IMAGE)

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
	DB_DEBUG=0 bin/milmove migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"

.PHONY: db_dev_migrate
db_dev_migrate: db_dev_migrate_standalone ## Migrate Dev DB

.PHONY: db_dev_psql
db_dev_psql: ## Open PostgreSQL shell for Dev DB
	scripts/psql-dev

.PHONY: db_dev_fresh
db_dev_fresh: db_dev_reset db_dev_migrate
	@echo "Ensure that you're running the correct APPLICATION..."
	./scripts/ensure-application app
	@echo "Populate the ${DB_NAME_DEV} database..."
	go run github.com/transcom/mymove/cmd/generate-test-data --named-scenario="dev_seed" --db-env="development"

.PHONY: db_dev_e2e_populate
db_dev_e2e_populate: db_dev_migrate ## Populate Dev DB with generated e2e (end-to-end) data
	@echo "Ensure that you're running the correct APPLICATION..."
	./scripts/ensure-application app
	@echo "Truncate the ${DB_NAME_DEV} database..."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_DEV)/$(DB_NAME_DEV)?sslmode=disable -c 'TRUNCATE users CASCADE; TRUNCATE uploads CASCADE; TRUNCATE webhook_subscriptions;'
	@echo "Populate the ${DB_NAME_DEV} database..."
	go run github.com/transcom/mymove/cmd/generate-test-data --named-scenario="dev_seed" --db-env="development"

## Alias for db_dev_bandwidth_up
## We started with `db_bandwidth_up`, which some folks are already using, and
## then renamed it to `db_dev_bandwidth_up`. To allow folks to keep using the
## name they're familiar with, we've added this alias to the renamed command.
.PHONY: db_bandwidth_up
db_bandwidth_up: db_dev_bandwidth_up

.PHONY: db_dev_bandwidth_up
db_dev_bandwidth_up: bin/generate-test-data	 ## Truncate Dev DB and Generate data for bandwidth tests
	@echo "Ensure that you're running the correct APPLICATION..."
	./scripts/ensure-application app
	@echo "Truncate the ${DB_NAME_DEV} database..."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_DEV)/$(DB_NAME_DEV)?sslmode=disable -c 'TRUNCATE users CASCADE; TRUNCATE uploads CASCADE; TRUNCATE webhook_subscriptions;'
	@echo "Populate the ${DB_NAME_DEV} database..."
	DB_PORT=$(DB_PORT_DEV) go run github.com/transcom/mymove/cmd/generate-test-data --named-scenario="bandwidth" --db-env="development"
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
db_deployed_migrations_migrate_standalone: bin/milmove ## Migrate Deployed Migrations DB with local secure migrations
	@echo "Migrating the ${DB_NAME_DEPLOYED_MIGRATIONS} database..."
	DB_DEBUG=0 DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) bin/milmove migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"

.PHONY: db_deployed_migrations_migrate
db_deployed_migrations_migrate: db_deployed_migrations_migrate_standalone ## Migrate Deployed Migrations DB

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
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_TEST)?sslmode=disable -c 'DROP DATABASE IF EXISTS $(DB_NAME_TEST);'
endif

.PHONY: db_test_start
db_test_start: ## Start Test DB
ifndef CIRCLECI
	brew services stop postgresql 2> /dev/null || true
	@echo "Starting the ${DB_DOCKER_CONTAINER_TEST} docker database container..."
	docker start $(DB_DOCKER_CONTAINER_TEST) || \
		docker run --name $(DB_DOCKER_CONTAINER_TEST) \
			-e \
			POSTGRES_PASSWORD=$(PGPASSWORD) \
			-d \
			-p $(DB_PORT_TEST):$(DB_PORT_DOCKER)\
			--mount type=tmpfs,destination=/var/lib/postgresql/data \
			$(DB_DOCKER_CONTAINER_IMAGE)
else
	@echo "Relying on CircleCI's database setup to start the DB."
endif

.PHONY: db_test_create
db_test_create: ## Create Test DB
ifndef CIRCLECI
	@echo "Create the ${DB_NAME_TEST} database..."
	DB_NAME=postgres DB_PORT=$(DB_PORT_TEST) scripts/wait-for-db && \
		createdb -p $(DB_PORT_TEST) -h localhost -U postgres $(DB_NAME_TEST) || true
else
	@echo "Relying on CircleCI's database setup to create the DB."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_TEST)?sslmode=disable -c 'CREATE DATABASE $(DB_NAME_TEST);'
endif

.PHONY: db_test_run
db_test_run: db_test_start db_test_create ## Run Test DB

.PHONY: db_test_reset
db_test_reset: db_test_destroy db_test_run ## Reset Test DB (destroy and run)

.PHONY: db_test_migrate_standalone
db_test_migrate_standalone: bin/milmove ## Migrate Test DB directly
ifndef CIRCLECI
	@echo "Migrating the ${DB_NAME_TEST} database..."
	DB_DEBUG=0 DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_TEST) bin/milmove migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"
else
	@echo "Migrating the ${DB_NAME_TEST} database..."
	DB_DEBUG=0 DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_DEV) bin/milmove migrate -p "file://migrations/${APPLICATION}/secure;file://migrations/${APPLICATION}/schema" -m "migrations/${APPLICATION}/migrations_manifest.txt"
endif

.PHONY: db_test_migrate
db_test_migrate: db_test_migrate_standalone ## Migrate Test DB

.PHONY: db_test_migrations_build
db_test_migrations_build: .db_test_migrations_build.stamp ## Build Test DB Migrations Docker Image
.db_test_migrations_build.stamp:
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
e2e_test: db_test_migrate db_e2e_up ## Run e2e (end-to-end) integration tests
	$(AWS_VAULT) ./scripts/run-e2e-test

.PHONY: e2e_test_fresh ## Build everything from scratch before running tests
e2e_test_fresh: bin/gin server_generate server_build client_build db_e2e_init
	$(AWS_VAULT) ./scripts/run-e2e-test

.PHONY: e2e_mtls_test_docker
e2e_mtls_test_docker: ## Run e2e (end-to-end) integration tests with docker
	$(AWS_VAULT) ./scripts/run-e2e-mtls-test-docker

.PHONY: e2e_test_docker
e2e_test_docker: ## Run e2e (end-to-end) integration tests with docker
	$(AWS_VAULT) ./scripts/run-e2e-test-docker

.PHONY: e2e_test_docker_mymove
e2e_test_docker_mymove: ## Run e2e (end-to-end) Service Member integration tests with docker
	$(AWS_VAULT) SPEC=cypress/integration/mymove/**/* ./scripts/run-e2e-test-docker

.PHONY: e2e_test_docker_office
e2e_test_docker_office: ## Run e2e (end-to-end) Office integration tests with docker
	$(AWS_VAULT) SPEC=cypress/integration/office/**/* ./scripts/run-e2e-test-docker

.PHONY: e2e_test_docker_api
e2e_test_docker_api: ## Run e2e (end-to-end) API integration tests with docker
	$(AWS_VAULT) SPEC=cypress/integration/api/**/* ./scripts/run-e2e-test-docker

.PHONY: e2e_clean
e2e_clean: ## Clean e2e (end-to-end) files and docker images
	rm -f .*_linux.stamp
	rm -rf cypress/results
	rm -rf cypress/screenshots
	rm -rf cypress/videos
	rm -rf cypress/reports
	docker rm -f cypress || true

.PHONY: db_e2e_up
db_e2e_up: bin/generate-test-data ## Truncate Test DB and Generate e2e (end-to-end) data
	@echo "Ensure that you're running the correct APPLICATION..."
	./scripts/ensure-application app
	@echo "Truncate the ${DB_NAME_TEST} database..."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_TEST)/$(DB_NAME_TEST)?sslmode=disable -c 'TRUNCATE users CASCADE; TRUNCATE uploads CASCADE;'
	@echo "Populate the ${DB_NAME_TEST} database..."
	DB_PORT=$(DB_PORT_TEST) go run github.com/transcom/mymove/cmd/generate-test-data --named-scenario="e2e_basic" --db-env="test"

.PHONY: rerun_e2e_tests_with_new_data
rerun_e2e_tests_with_new_data: db_e2e_up
	$(AWS_VAULT) ./scripts/run-e2e-test

.PHONY: db_e2e_init
db_e2e_init: db_test_reset db_test_migrate redis_reset db_e2e_up ## Initialize e2e (end-to-end) DB (reset, migrate, up)

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
	docker rm -f tasks || true

.PHONY: tasks_build
tasks_build: server_generate bin/milmove-tasks ## Build Scheduled Task dependencies

.PHONY: tasks_build_docker
tasks_build_docker: server_generate bin/milmove-tasks ## Build Scheduled Task dependencies and Docker image
	@echo "Build the docker scheduled tasks container..."
	docker build -f Dockerfile.tasks --tag $(TASKS_DOCKER_CONTAINER):latest .

.PHONY: tasks_build_linux_docker
tasks_build_linux_docker:  ## Build Scheduled Task binaries (linux) and Docker image (local)
	@echo "Build the docker scheduled tasks container..."
	docker build -f Dockerfile.tasks_local --tag $(TASKS_DOCKER_CONTAINER):latest .

.PHONY: tasks_connect_to_gex_via_sftp
tasks_connect_to_gex_via_sftp: tasks_build_linux_docker ## Run connect-to-gex-via-sftp from inside docker container
	@echo "Connecting to GEX via SFTP with docker command..."
	DB_NAME=$(DB_NAME_DEV) DB_DOCKER_CONTAINER=$(DB_DOCKER_CONTAINER_DEV) scripts/wait-for-db-docker
	docker run \
		-t \
		-e DB_HOST="database" \
		-e DB_NAME \
		-e DB_PORT \
		-e DB_USER \
		-e DB_PASSWORD \
		-e GEX_SFTP_HOST \
		-e GEX_SFTP_HOST_KEY \
		-e GEX_SFTP_IP_ADDRESS \
		-e GEX_SFTP_PASSWORD \
		-e GEX_SFTP_PORT \
		-e GEX_SFTP_USER_ID \
		--link="$(DB_DOCKER_CONTAINER_DEV):database" \
		--rm \
		$(TASKS_DOCKER_CONTAINER):latest \
		milmove-tasks connect-to-gex-via-sftp

.PHONY: tasks_process_edis
tasks_process_edis: tasks_build_linux_docker ## Run process-edis from inside docker container
	@echo "Processing EDIs with docker command..."
	DB_NAME=$(DB_NAME_DEV) DB_DOCKER_CONTAINER=$(DB_DOCKER_CONTAINER_DEV) scripts/wait-for-db-docker
	docker run \
		-t \
		-e DB_HOST="database" \
		-e DB_NAME \
		-e DB_PORT \
		-e DB_USER \
		-e DB_PASSWORD \
		--link="$(DB_DOCKER_CONTAINER_DEV):database" \
		--rm \
		$(TASKS_DOCKER_CONTAINER):latest \
		milmove-tasks process-edis

.PHONY: tasks_save_ghc_fuel_price_data
tasks_save_ghc_fuel_price_data: tasks_build_linux_docker ## Run save-ghc-fuel-price-data from inside docker container
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
		milmove-tasks save-ghc-fuel-price-data

tasks_send_post_move_survey: tasks_build_linux_docker ## Run send-post-move-survey from inside docker container
	@echo "sending post move survey with docker command..."
	DB_NAME=$(DB_NAME_DEV) DB_DOCKER_CONTAINER=$(DB_DOCKER_CONTAINER_DEV) scripts/wait-for-db-docker
	docker run \
		-t \
		-e DB_HOST="database" \
		-e DB_NAME \
		-e DB_PORT \
		-e DB_USER \
		-e DB_PASSWORD \
		--link="$(DB_DOCKER_CONTAINER_DEV):database" \
		--rm \
		$(TASKS_DOCKER_CONTAINER):latest \
		milmove-tasks send-post-move-survey


tasks_send_payment_reminder: tasks_build_linux_docker ## Run send-payment-reminder from inside docker container
	@echo "sending payment reminder with docker command..."
	DB_NAME=$(DB_NAME_DEV) DB_DOCKER_CONTAINER=$(DB_DOCKER_CONTAINER_DEV) scripts/wait-for-db-docker
	docker run \
		-t \
		-e DB_HOST="database" \
		-e DB_NAME \
		-e DB_PORT \
		-e DB_USER \
		-e DB_PASSWORD \
		--link="$(DB_DOCKER_CONTAINER_DEV):database" \
		--rm \
		$(TASKS_DOCKER_CONTAINER):latest \
		milmove-tasks send-payment-reminder

tasks_post_file_to_gex: tasks_build_linux_docker ## Run post-file-to-gex from inside docker container
	@echo "sending payment reminder with docker command..."
	DB_NAME=$(DB_NAME_DEV) DB_DOCKER_CONTAINER=$(DB_DOCKER_CONTAINER_DEV) scripts/wait-for-db-docker
	docker run \
		-t \
		-e DB_HOST="database" \
		-e DB_NAME \
		-e DB_PORT \
		-e DB_USER \
		-e DB_PASSWORD \
		--link="$(DB_DOCKER_CONTAINER_DEV):database" \
		--rm \
		$(TASKS_DOCKER_CONTAINER):latest \
		milmove-tasks post-file-to-gex
#
# ----- END SCHEDULED TASK TARGETS -----
#

#
# ----- START Deployed MIGRATION TARGETS -----
#

.PHONY: run_prd_migrations
run_prd_migrations: bin/milmove db_deployed_migrations_reset ## Run GovCloud prd migrations against Deployed Migrations DB
	@echo "Migrating the prd-migrations database with prd migrations..."
	MIGRATION_PATH="s3://transcom-gov-milmove-prd-app-us-gov-west-1/secure-migrations;file://migrations/$(APPLICATION)/schema" \
	DB_HOST=localhost \
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) \
	DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) \
	DB_DEBUG=0 \
	DISABLE_AWS_VAULT_WRAPPER=1 \
	AWS_REGION=us-gov-west-1 \
	aws-vault exec transcom-gov-milmove-prd \
	bin/milmove migrate

.PHONY: run_stg_migrations
run_stg_migrations: bin/milmove db_deployed_migrations_reset ## Run GovCloud stg migrations against Deployed Migrations DB
	@echo "Migrating the stg-migrations database with stg migrations..."
	MIGRATION_PATH="s3://transcom-gov-milmove-stg-app-us-gov-west-1/secure-migrations;file://migrations/$(APPLICATION)/schema" \
	DB_HOST=localhost \
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) \
	DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) \
	DB_DEBUG=0 \
	DISABLE_AWS_VAULT_WRAPPER=1 \
	AWS_REGION=us-gov-west-1 \
	aws-vault exec transcom-gov-milmove-stg \
	bin/milmove migrate

.PHONY: run_exp_migrations
run_exp_migrations: bin/milmove db_deployed_migrations_reset ## Run GovCloud exp migrations against Deployed Migrations DB
	@echo "Migrating the exp-migrations database with exp migrations..."
	MIGRATION_PATH="s3://transcom-gov-milmove-exp-app-us-gov-west-1/secure-migrations;file://migrations/$(APPLICATION)/schema" \
	DB_HOST=localhost \
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) \
	DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) \
	DB_DEBUG=0 \
	DISABLE_AWS_VAULT_WRAPPER=1 \
	AWS_REGION=us-gov-west-1 \
	aws-vault exec transcom-gov-milmove-exp \
	bin/milmove migrate

#
# ----- END PROD_MIGRATION TARGETS -----
#

#
# ----- START WEBHOOK CLIENT TARGETS -----
#

.PHONY: webhook_client_docker
webhook_client_docker:
	docker build -f Dockerfile.webhook_client_local -t $(WEBHOOK_CLIENT_DOCKER_CONTAINER):latest .

.PHONY: webhook_client_start_reset_db
webhook_client_start_reset_db: db_dev_e2e_populate webhook_client_start

.PHONY: webhook_client_start
webhook_client_start:
	@echo "Starting the webhook client..."
	# Note regarding the use of 172.17.0.1: the default network bridge in Docker is 172.17.0.1
	# according to https://docs.docker.com/network/network-tutorial-standalone/. Based on Internet
	# searches, this IP address seems to be fairly static. Therefore, this address can be used to
	# serve as the IP address for various local hostnames used during testing. If this stops
	# working some day, dynamic resolution may be required to look up the host.docker.internal IP.
	# However, by using a static IP, it allows developers to restart Docker after building an image
	# and still have the container work if the /etc/hosts file had already been updated with this
	# script. While brainstorming a usable approach, we tried to introduce the use of
	# host.docker.internal by using the HOSTALIASES environment variable
	# (https://man7.org/linux/man-pages/man7/hostname.7.html), but this resulted in a segfaults
	# occurring during DNS lookup from within the container about 70% of the time. So we opted for
	# the more stable hardcoded Docker gateway IP approch.
	#
	# If this fails, make sure you are running both the mTLS server and the database in containers.
	# Due to the fact that on MacOS, Docker containers are run within a virtual machine rather than
	# directly from the host system, if the webhook client is running inside a container, it won't
	# be able to see anything outside of that Docker-managed virtual machine. Therefore, you need
	# to run the mTLS server and database from inside a container. This is possible with a command
	# like the following:
	#
	#   docker-compose -f docker-compose.mtls_local.yml --compatibility up --remove-orphans
	#
	# For more information about this, please see the following page:
	# https://docs.docker.com/docker-for-mac/networking/#known-limitations-use-cases-and-workarounds
	docker run \
		--add-host "adminlocal:172.17.0.1" \
		--add-host "milmovelocal:172.17.0.1" \
		--add-host "officelocal:172.17.0.1" \
		--add-host "orderslocal:172.17.0.1" \
		--add-host "primelocal:172.17.0.1" \
		-e DB_HOST=172.17.0.1 \
		-e DB_NAME \
		-e DB_PORT \
		-e DB_USER \
		-e DB_PASSWORD \
		-e GEX_MTLS_CLIENT_CERT \
		-e GEX_MTLS_CLIENT_KEY \
		-e LOGGING_LEVEL=debug \
		-e PERIOD \
		$(WEBHOOK_CLIENT_DOCKER_CONTAINER):latest

.PHONY: webhook_client_test
webhook_client_test: db_e2e_init webhook_client_test_standalone

.PHONY: webhook_client_test
webhook_client_test_standalone:
	go test -v -count 1 -short ./cmd/webhook-client/webhook

#
# ----- END WEBHOOK CLIENT TARGETS -----
#

#
# ----- START PRIME TARGETS -----
#

.PHONY: run_prime_docker
run_prime_docker: ## Runs the docker that spins up the Prime API and data to test with
	scripts/run-prime-docker

#
# ----- END PRIME TARGETS -----
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
adr_update: ## Update ADR Log
	pre-commit run -a gen-docs

.PHONY: gofmt
gofmt:  ## Run go fmt over all Go files
	go fmt $$(go list ./...) >> /dev/null

.PHONY: pre_commit_tests
pre_commit_tests: .client_deps.stamp ## Run pre-commit tests
	pre-commit run --all-files

.PHONY: pretty
pretty: gofmt ## Run code through JS and Golang formatters
	npx prettier --write --loglevel warn "src/**/*.{js,jsx}"

.PHONY: docker_circleci
docker_circleci: ## Run CircleCI container locally with project mounted
	docker pull milmove/circleci-docker:milmove-app-990c528cc6bfd9e9693fa28aae500d0f577075f6
	docker run -it --rm=true -v $(PWD):$(PWD) -w $(PWD) -e CIRCLECI=1 milmove/circleci-docker:milmove-app bash

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
	rm -rf ./build
	rm -rf ./node_modules
	rm -rf ./public/swagger-ui/*.{css,js,png}
	rm -rf ./tmp/secure_migrations
	rm -rf ./tmp/storage
	rm -rf $(SCHEMASPY_OUTPUT)
	rm -rf ./storybook-static
	rm -rf ./coverage
	rm -rf ./log

.PHONY: spellcheck
spellcheck: ## Run interactive spellchecker
	@which mdspell -s || (echo "Install mdspell with yarn global add markdown-spellcheck" && exit 1)
	/usr/local/bin/mdspell --ignore-numbers --ignore-acronyms --en-us --no-suggestions \
		`find . -type f -name "*.md" \
			-not -path "./node_modules/*" \
			-not -path "./vendor/*" \
			-not -path "./docs/adr/index.md" | sort`

.PHONY: storybook
storybook: ## Start the storybook server
	yarn run storybook

.PHONY: storybook_build
storybook_build: ## Build static storybook site
	yarn run build-storybook

.PHONY: schemaspy
schemaspy: db_test_reset db_test_migrate ## Generates database documentation using schemaspy
	rm -rf $(SCHEMASPY_OUTPUT)
	docker run -v $(PWD)/$(SCHEMASPY_OUTPUT):/output schemaspy/schemaspy:latest \
		-t pgsql11 -host host.docker.internal -port $(DB_PORT_TEST) -db $(DB_NAME_TEST) -u postgres -p $(PGPASSWORD) \
		-norows -nopages
	@echo "Schemaspy output can be found in $(SCHEMASPY_OUTPUT)"

.PHONY: reviewapp_docker
reviewapp_docker:
	docker-compose -f docker-compose.reviewapp.yml up

.PHONY: reviewapp_docker_build
reviewapp_docker_build:
# remove bin to maybe speed up docker builds by removing it from
# docker context
	rm -rf ./bin
	docker-compose -f docker-compose.reviewapp.yml build

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

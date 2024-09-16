DB_NAME_DEV = dev_db
DB_NAME_DEPLOYED_MIGRATIONS = deployed_migrations
DB_NAME_TEST = test_db
DB_DOCKER_CONTAINER_DEV = milmove-db-dev
DB_DOCKER_CONTAINER_DEPLOYED_MIGRATIONS = milmove-db-deployed-migrations
DB_DOCKER_CONTAINER_TEST = milmove-db-test
DB_DOCKER_CONTAINER_IMAGE = postgres:16.4
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

export DEVSEED_SUBSCENARIO

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
ensure_pre_commit: .git/hooks/pre-commit install_pre_commit ## Ensure pre-commit hooks are installed
.git/hooks/pre-commit: .check_pre-commit_installed.stamp
.check_pre-commit_installed.stamp: ## Ensure pre-commit is installed
ifeq (, $(shell which pre-commit))
	$(error pre-commit is not installed. Install with `brew install pre-commit`.)
else
	@echo "pre-commit is installed"
	touch .check_pre-commit_installed.stamp
endif

.PHONY: install_pre_commit
install_pre_commit:  ## Installs pre-commit hooks
	pre-commit install
	pre-commit install-hooks

.PHONY: prereqs
prereqs: ## Check that pre-requirements are installed, includes dependency scripts
	scripts/prereqs

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
.check_go_version.stamp: scripts/check-go-version .tool-versions
	scripts/check-go-version
	touch .check_go_version.stamp

.PHONY: check_gopath
check_gopath: .check_gopath.stamp ## Check that $GOPATH exists in $PATH
.check_gopath.stamp: scripts/check-gopath go.sum # Make sure any go binaries rebuild if version possibly changes
ifndef CIRCLECI
	scripts/check-gopath
else
	@echo "No need to check go path on CircleCI."
endif
	touch .check_gopath.stamp

.PHONY: check_node_version
check_node_version: .check_node_version.stamp ## Check that the correct Node version is installed
.check_node_version.stamp: scripts/check-node-version .tool-versions
	scripts/check-node-version
	touch .check_node_version.stamp

.PHONY: check_docker_size
check_docker_size: ## Check the amount of disk space used by docker
	scripts/check-docker-size

.PHONY: deps
deps: prereqs ensure_pre_commit deps_shared ## Run all checks and install all dependencies

.PHONY: setup
setup:
	scripts/setup

.PHONY: deps_nix
deps_nix: install_pre_commit deps_shared ## Nix equivalent (kind of) of `deps` target.

.PHONY: deps_shared
deps_shared: client_deps bin/rds-ca-2019-root.pem bin/rds-ca-rsa4096-g1.pem ## install dependencies

.PHONY: test
test: client_test server_test e2e_test ## Run all tests

.PHONY: diagnostic
diagnostic: check_docker_size ## Run diagnostic scripts on environment

.PHONY: check_log_dir
check_log_dir: ## Make sure we have a log directory
	mkdir -p log

.PHONY: check_app
check_app: ## Make sure you're running the correct APP
	@echo "Ensure that you're running the correct APPLICATION..."
	./scripts/ensure-application app
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
client_deps: .check_hosts.stamp .client_deps.stamp ## Install client dependencies
.client_deps.stamp: yarn.lock .check_node_version.stamp
	yarn install
	scripts/copy-swagger-ui
	touch .client_deps.stamp

.client_build.stamp: .client_deps.stamp $(shell find src -type f)
	REACT_APP_GIT_COMMIT=$(GIT_COMMIT) \
	REACT_APP_GIT_BRANCH=$(GIT_BRANCH) \
	yarn build
	touch .client_build.stamp

.PHONY: client_build
client_build: .client_build.stamp ## Build the client

build/index.html: build/downloads ## milmove serve requires this file to boot, but it isn't used during local development
	mkdir -p build
	touch build/index.html

build/downloads: public/downloads
	mkdir -p build
	rm -rf build/downloads
	cp -r public/downloads build/downloads

.PHONY: client_run
client_run: .client_deps.stamp ## Run MilMove Service Member client
	REACT_APP_GIT_COMMIT=$(GIT_COMMIT) \
	REACT_APP_GIT_BRANCH=$(GIT_BRANCH) \
	HOST=milmovelocal \
	yarn start

.PHONY: client_test
client_test: .client_deps.stamp ## Run client unit tests
	yarn test

.PHONY: client_test_coverage
client_test_coverage : .client_deps.stamp ## Run client unit test coverage
	yarn test:coverage

.PHONY: office_client_run
office_client_run: .client_deps.stamp ## Run MilMove Office client
	REACT_APP_GIT_COMMIT=$(GIT_COMMIT) \
	REACT_APP_GIT_BRANCH=$(GIT_BRANCH) \
	HOST=officelocal \
	yarn start

.PHONY: admin_client_run
admin_client_run: .client_deps.stamp ## Run MilMove Admin client
	REACT_APP_GIT_COMMIT=$(GIT_COMMIT) \
	REACT_APP_GIT_BRANCH=$(GIT_BRANCH) \
	HOST=adminlocal \
	yarn start



#
# ----- END CLIENT TARGETS -----
#

#
# ----- START BIN TARGETS -----
#

### Go Tool Targets

bin/gin: .check_go_version.stamp .check_gopath.stamp pkg/tools/tools.go
	go build -ldflags "$(LDFLAGS)" -o bin/gin github.com/codegangsta/gin

bin/soda: .check_go_version.stamp .check_gopath.stamp pkg/tools/tools.go
	go build -ldflags "$(LDFLAGS)" -o bin/soda github.com/gobuffalo/pop/v6/soda

# No static linking / $(LDFLAGS) because gotestsum is only used for building the CirlceCi test report
bin/gotestsum: .check_go_version.stamp .check_gopath.stamp pkg/tools/tools.go
	go build -o bin/gotestsum gotest.tools/gotestsum

# No static linking / $(LDFLAGS) because mockery is only used for testing
bin/mockery: .check_go_version.stamp .check_gopath.stamp pkg/tools/tools.go
	go build -o bin/mockery github.com/vektra/mockery/v2

# No static linking / $(LDFLAGS) because swagger is only used for code generation
bin/swagger: .check_go_version.stamp .check_gopath.stamp pkg/tools/tools.go
	go build -o bin/swagger github.com/go-swagger/go-swagger/cmd/swagger

### Cert Targets
# AWS is only providing a bundle for the 2022 cert, which includes 2017? and rds-ca-rsa4096-g1
bin/rds-ca-rsa4096-g1.pem:
	mkdir -p bin/
	curl -sSo bin/rds-ca-rsa4096-g1.pem https://truststore.pki.us-gov-west-1.rds.amazonaws.com/us-gov-west-1/us-gov-west-1-bundle.pem

bin/rds-ca-2019-root.pem:
	mkdir -p bin/
	curl -sSo bin/rds-ca-2019-root.pem https://s3.amazonaws.com/rds-downloads/rds-ca-2019-root.pem

### MilMove Targets

bin/big-cat: cmd/big-cat
	go build -ldflags "$(LDFLAGS)" -o bin/big-cat ./cmd/big-cat

bin/model-vet: cmd/model-vet
	go build -ldflags "$(LDFLAGS)" -o bin/model-vet ./cmd/model-vet

bin/generate-deploy-notes: cmd/generate-deploy-notes
	go build -ldflags "$(LDFLAGS)" -o bin/generate-deploy-notes ./cmd/generate-deploy-notes

bin/ecs-deploy: cmd/ecs-deploy
	go build -ldflags "$(LDFLAGS)" -o bin/ecs-deploy ./cmd/ecs-deploy

bin/generate-shipment-summary: cmd/generate-shipment-summary
	go build -ldflags "$(LDFLAGS)" -o bin/generate-shipment-summary ./cmd/generate-shipment-summary

bin/generate-test-data: cmd/generate-test-data
	@echo "WARNING: devseed data is being deprecated on 11/08/2023. This function will be deleted after this date."
	go build -ldflags "$(LDFLAGS)" -o bin/generate-test-data ./cmd/generate-test-data

bin/ghc-pricing-parser: cmd/ghc-pricing-parser
	go build -ldflags "$(LDFLAGS)" -o bin/ghc-pricing-parser ./cmd/ghc-pricing-parser

bin/ghc-transit-time-parser: cmd/ghc-transit-time-parser
	go build -ldflags "$(LDFLAGS)" -o bin/ghc-transit-time-parser ./cmd/ghc-transit-time-parser

bin/health-checker: cmd/health-checker
	go build -ldflags "$(LDFLAGS)" -o bin/health-checker ./cmd/health-checker

bin/iws: cmd/iws
	go build -ldflags "$(LDFLAGS)" -o bin/iws ./cmd/iws/iws.go

PKG_GOSRC := $(shell find pkg -name '*.go')

bin/milmove: $(shell find cmd/milmove -name '*.go') $(PKG_GOSRC) .check_go_version.stamp .check_gopath.stamp
	go build -gcflags="$(GOLAND_GC_FLAGS) $(GC_FLAGS)" -asmflags=-trimpath=$(GOPATH) -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/milmove ./cmd/milmove

bin/milmove-tasks: $(shell find cmd/milmove-tasks -name '*.go') $(PKG_GOSRC) .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/milmove-tasks ./cmd/milmove-tasks

bin/prime-api-client: $(shell find cmd/prime-api-client -name '*.go') $(PKG_GOSRC) .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/prime-api-client ./cmd/prime-api-client

bin/webhook-client: $(shell find cmd/webhook-client -name '*.go') $(PKG_GOSRC) .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/webhook-client ./cmd/webhook-client

bin/read-alb-logs: $(shell find cmd/read-alb-logs -name '*.go') $(PKG_GOSRC) .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/read-alb-logs ./cmd/read-alb-logs

bin/send-to-gex: $(shell find cmd/send-to-gex -name '*.go') $(PKG_GOSRC) .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/send-to-gex ./cmd/send-to-gex

bin/tls-checker: $(shell find cmd/tls-checker -name '*.go') $(PKG_GOSRC) .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/tls-checker ./cmd/tls-checker

bin/generate-payment-request-edi: $(shell find cmd/generate-payment-request-edi -name '*.go') $(PKG_GOSRC) .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/generate-payment-request-edi ./cmd/generate-payment-request-edi

bin/simulate-process-tpps: $(shell find cmd/simulate-process-tpps -name '*.go') $(PKG_GOSRC) .check_go_version.stamp .check_gopath.stamp
	go build -ldflags "$(LDFLAGS)" -o bin/simulate-process-tpps ./cmd/simulate-process-tpps

#
# ----- END BIN TARGETS -----
#

#
# ----- START SERVER TARGETS -----
#

swagger_generate: .swagger_build.stamp ## Check that the build files haven't been manually edited to prevent overwrites

# If any swagger files (source or generated) have changed, re-run so
# we can warn on improperly modified files. Look for any files so that
# if API docs have changed, swagger regeneration will capture those
# changes. On Circle CI, or if the user has set
# SWAGGER_AUTOREBUILD, rebuild automatically without asking
ifdef CIRCLECI
SWAGGER_AUTOREBUILD=1
endif
SWAGGER_FILES = $(shell find swagger swagger-def -type f)
.swagger_build.stamp: $(SWAGGER_FILES)
ifeq ($(SWAGGER_AUTOREBUILD),0)
ifneq ("$(shell find swagger -type f -name '*.yaml' -newer .swagger_build.stamp)","")
	@echo "Unexpected changes found in swagger build files. Code may be overwritten."
	@read -p "Continue with rebuild? [y/N] : " ANS && test "$${ANS}" == "y" || (echo "Exiting rebuild."; false)
endif
endif
	./scripts/openapi bundle -o swagger/ ## Bundles the API definition files into a complete specification
	touch .swagger_build.stamp

server_generate: .server_generate.stamp

.server_generate.stamp: .check_go_version.stamp .check_gopath.stamp .swagger_build.stamp bin/swagger $(wildcard swagger/*.yaml) ## Generate golang server code from Swagger files
	scripts/gen-server
	touch .server_generate.stamp

.PHONY: server_build
server_build: bin/milmove ## Build the server

# This command is for running the server by itself, it will serve the compiled frontend on its own
# Note: Don't double wrap with aws-vault because the pkg/cli/vault.go will handle it
server_run_standalone: check_log_dir server_build client_build db_dev_run redis_run
	./bin/milmove serve 2>&1 | tee -a log/dev.log

# This command will rebuild the swagger go code and rerun server on any changes
server_run:
	find ./swagger-def -type f | entr -c -r make server_run_default
# This command runs the server behind gin, a hot-reload server
# Note: Gin is not being used as a proxy so assigning odd port and laddr to keep in IPv4 space.
# Note: The INTERFACE envar is set to configure the gin build, milmove_gin, local IP4 space with default port GIN_PORT.
server_run_default: .check_hosts.stamp .check_go_version.stamp .check_gopath.stamp .check_node_version.stamp check_log_dir bin/gin build/index.html server_generate db_dev_run db_dev_migrate redis_run
	INTERFACE=localhost \
		./bin/gin \
		--build ./cmd/milmove \
		--bin /bin/milmove_gin \
		--laddr 127.0.0.1 --port "$(GIN_PORT)" \
		--excludeDir node_modules \
		--immediate \
		--buildArgs "-ldflags=\"$(WEBSERVER_LDFLAGS)\"" \
		serve \
		2>&1 | tee -a log/dev.log

.PHONY: server_run_debug
server_run_debug: .check_hosts.stamp .check_go_version.stamp .check_gopath.stamp .check_node_version.stamp check_log_dir build/index.html server_generate db_dev_run redis_run ## Debug the server
	scripts/kill-process-on-port 8080
	scripts/kill-process-on-port 9443
	DISABLE_AWS_VAULT_WRAPPER=1 \
	AWS_REGION=us-gov-west-1 \
	aws-vault exec transcom-gov-dev -- \
	dlv debug -l 127.0.0.1:38697 --headless cmd/milmove/*.go -- serve 2>&1 | tee -a log/dev.log

.PHONY: build_tools
build_tools: bin/gin \
	bin/mockery \
	bin/swagger \
	bin/rds-ca-rsa4096-g1.pem \
	bin/rds-ca-2019-root.pem \
	bin/big-cat \
	bin/generate-deploy-notes \
	bin/ecs-deploy \
	bin/generate-payment-request-edi \
	bin/generate-shipment-summary \
	bin/generate-test-data \
	bin/ghc-pricing-parser \
	bin/ghc-transit-time-parser \
	bin/health-checker \
	bin/iws \
	bin/milmove-tasks \
	bin/model-vet \
	bin/prime-api-client \
	bin/webhook-client \
	bin/read-alb-logs \
	bin/send-to-gex \
	bin/simulate-process-tpps \
	bin/tls-checker ## Build all tools

.PHONY: build
build: server_build build_tools client_build ## Build the server, tools, and client

.PHONY: mocks_generate
mocks_generate: bin/mockery ## Generate mockery mocks for tests
	go generate $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

.PHONY: server_test_setup
server_test_setup: db_test_reset db_test_migrate redis_reset db_test_truncate

.PHONY: server_test
server_test: server_test_setup server_test_standalone ## Run server unit tests

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
ifdef CIRCLECI
	@echo "Relying on CircleCI to setup redis."
else
	docker pull $(REDIS_DOCKER_CONTAINER_IMAGE)
endif

.PHONY: redis_destroy
redis_destroy: ## Destroy Redis
ifdef CIRCLECI
	@echo "Relying on CircleCI to setup redis."
else
	@echo "Destroying the ${REDIS_DOCKER_CONTAINER} docker redis container..."
	docker rm -f $(REDIS_DOCKER_CONTAINER) || echo "No Redis container"
endif

.PHONY: redis_run
redis_run: redis_pull ## Run Redis
ifdef CIRCLECI
	@echo "Relying on CircleCI to setup redis."
else
		@echo "Stopping the Redis brew service in case it's running..."
		brew services stop redis 2> /dev/null || true
	@echo "Starting the ${REDIS_DOCKER_CONTAINER} docker redis container..."
	docker start $(REDIS_DOCKER_CONTAINER) || \
		docker run -d --name $(REDIS_DOCKER_CONTAINER) \
			-p $(REDIS_PORT):$(REDIS_PORT_DOCKER) \
			$(REDIS_DOCKER_CONTAINER_IMAGE)
endif

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
	@echo "Starting the ${DB_DOCKER_CONTAINER_DEV} docker database container..."
	# If running do nothing, if not running try to start, if can't start then run
	docker start $(DB_DOCKER_CONTAINER_DEV) || \
		docker run -d --name $(DB_DOCKER_CONTAINER_DEV) \
			-e POSTGRES_PASSWORD=$(PGPASSWORD) \
			-p $(DB_PORT_DEV):$(DB_PORT_DOCKER)\
			$(DB_DOCKER_CONTAINER_IMAGE)
else
	@echo "Relying on CircleCI's database setup to start the DB."
endif

.PHONY: db_dev_create
db_dev_create: ## Create Dev DB
ifndef CIRCLECI
	@echo "Create the ${DB_NAME_DEV} database..."
	DB_NAME=postgres scripts/wait-for-db && DB_NAME=postgres psql-wrapper "CREATE DATABASE $(DB_NAME_DEV);" || true
else
	@echo "Relying on CircleCI's database setup to create the DB."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT)?sslmode=disable -c 'CREATE DATABASE $(DB_NAME_DEV);'
endif


.PHONY: db_dev_run
db_dev_run: db_dev_start db_dev_create ## Run Dev DB (start and create)

.PHONY: db_dev_reset
db_dev_reset: db_dev_destroy db_dev_run ## Reset Dev DB (destroy and run)

.PHONY: db_dev_init
db_dev_init: db_dev_reset db_dev_migrate ## Init Dev DB (destroy, run, migrate)

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
db_dev_fresh: check_app db_dev_reset db_dev_migrate ## Recreate dev db from scratch and populate with devseed data
	@echo "WARNING: Devseed data is being deprecated on 11/08/2023. This function will be deleted after that date."
	@echo "Populate the ${DB_NAME_DEV} database..."
	go run github.com/transcom/mymove/cmd/generate-test-data --named-scenario="dev_seed" --db-env="development" --named-sub-scenario="${DEVSEED_SUBSCENARIO}"

.PHONY: db_dev_truncate
db_dev_truncate: ## Truncate dev db
	@echo "Truncate the ${DB_NAME_DEV} database..."
	psql postgres://postgres:$(PGPASSWORD)@${DB_HOST}:$(DB_PORT_DEV)/$(DB_NAME_DEV)?sslmode=disable -c 'TRUNCATE users, uploads, webhook_subscriptions, storage_facilities CASCADE'

.PHONY: db_dev_e2e_populate
db_dev_e2e_populate: check_app db_dev_migrate db_dev_truncate ## Migrate dev db and populate with devseed data
	@echo "WARNING: Devseed data is being deprecated on 11/08/2023. This function will be deleted after that date."
	@echo "Populate the ${DB_NAME_DEV} database..."
	go run github.com/transcom/mymove/cmd/generate-test-data --named-scenario="dev_seed" --db-env="development" --named-sub-scenario="${DEVSEED_SUBSCENARIO}"

## Alias for db_dev_bandwidth_up
## We started with `db_bandwidth_up`, which some folks are already using, and
## then renamed it to `db_dev_bandwidth_up`. To allow folks to keep using the
## name they're familiar with, we've added this alias to the renamed command.
.PHONY: db_bandwidth_up
db_bandwidth_up: db_dev_bandwidth_up

.PHONY: db_dev_bandwidth_up
db_dev_bandwidth_up: check_app bin/generate-test-data db_dev_truncate ## Truncate Dev DB and Generate data for bandwidth tests
	@echo "WARNING: devseed data is being deprecated on 11/08/2023. This function will be deleted after this date."
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
		createdb -p $(DB_PORT_DEPLOYED_MIGRATIONS) -h $(DB_HOST) -U postgres $(DB_NAME_DEPLOYED_MIGRATIONS) || true

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
		createdb -p $(DB_PORT_TEST) -h $(DB_HOST) -U postgres $(DB_NAME_TEST) || true
else
	@echo "Relying on CircleCI's database setup to create the DB."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_TEST)?sslmode=disable -c 'CREATE DATABASE $(DB_NAME_TEST);'
endif

.PHONY: db_test_run
db_test_run: db_test_start db_test_create ## Run Test DB

.PHONY: db_test_reset
db_test_reset: db_test_destroy db_test_run ## Reset Test DB (destroy and run)

.PHONY: db_test_truncate
db_test_truncate:
	@echo "Truncating ${DB_NAME_TEST} database..."
	DB_PORT=$(DB_PORT_TEST) DB_NAME=$(DB_NAME_TEST) ./scripts/db-truncate

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
e2e_test: db_dev_run db_dev_truncate ## Run e2e (end-to-end) integration tests
	$(AWS_VAULT) ./scripts/run-e2e-test

.PHONY: e2e_test_fresh ## Build everything from scratch before running tests
e2e_test_fresh: db_dev_init
	$(AWS_VAULT) ./scripts/run-e2e-test

.PHONY: e2e_clean
e2e_clean: ## Clean e2e (end-to-end) files
	rm -f .*_linux.stamp
	rm -rf playwright-report

.PHONY: db_dev_e2e_backup
db_dev_e2e_backup: ## Backup Dev DB as 'e2e_dev'
	DB_NAME=$(DB_NAME_DEV) DB_PORT=$(DB_PORT_DEV) ./scripts/db-backup e2e_dev

.PHONY: db_dev_e2e_restore
db_dev_e2e_restore: ## Restore Dev DB from 'e2e_dev'
	DB_NAME=$(DB_NAME_DEV) DB_PORT=$(DB_PORT_DEV) ./scripts/db-restore e2e_dev

.PHONY: db_dev_e2e_cleanup
db_dev_e2e_cleanup: ## Clean up Dev DB backup `e2e_dev`
	./scripts/db-cleanup e2e_dev

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
		-e GEX_PRIVATE_KEY \
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

.PHONY: run_demo_migrations
run_demo_migrations: bin/milmove db_deployed_migrations_reset ## Run GovCloud demo migrations against Deployed Migrations DB
	@echo "Migrating the demo-migrations database with demo migrations..."
	MIGRATION_PATH="s3://transcom-gov-milmove-demo-app-us-gov-west-1/secure-migrations;file://migrations/$(APPLICATION)/schema" \
	DB_HOST=localhost \
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) \
	DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) \
	DB_DEBUG=0 \
	DISABLE_AWS_VAULT_WRAPPER=1 \
	AWS_REGION=us-gov-west-1 \
	aws-vault exec transcom-gov-milmove-demo \
	bin/milmove migrate

.PHONY: run_loadtest_migrations
run_loadtest_migrations: bin/milmove db_deployed_migrations_reset ## Run GovCloud loadtest migrations against Deployed Migrations DB
	@echo "Migrating the loadtest-migrations database with loadtest migrations..."
	MIGRATION_PATH="s3://transcom-gov-milmove-loadtest-app-us-gov-west-1/secure-migrations;file://migrations/$(APPLICATION)/schema" \
	DB_HOST=localhost \
	DB_PORT=$(DB_PORT_DEPLOYED_MIGRATIONS) \
	DB_NAME=$(DB_NAME_DEPLOYED_MIGRATIONS) \
	DB_DEBUG=0 \
	DISABLE_AWS_VAULT_WRAPPER=1 \
	AWS_REGION=us-gov-west-1 \
	aws-vault exec transcom-gov-milmove-loadtest \
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
		-e MOVE_MIL_INTEGRATIONS_DOD_TLS_CERT \
		-e MOVE_MIL_INTEGRATIONS_DOD_TLS_KEY \
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
	docker run -it --pull=always --rm=true -v $(PWD):$(PWD) -w $(PWD) -e CIRCLECI=1 milmove/circleci-docker:milmove-app-ab729849a08a773ea2557b19b67f378551d1ad3d bash

.PHONY: docker_local_ssh_server_with_password
docker_local_ssh_server_with_password:
	docker run --rm \
  --name sshd \
  -e USER_NAME=testu \
  -e USER_PASSWORD=testp \
  -e PASSWORD_ACCESS=true \
  -p 2222:2222 \
  -v some_local_upload_dir:/config/uploads \
 linuxserver/openssh-server


.PHONY: docker_local_ssh_server_with_key
docker_local_ssh_server_with_key:
	docker run --rm \
  --name sshd \
	-e PUBLIC_KEY="${TEST_GEX_PUBLIC_KEY}" \
	-e USER_NAME=testu \
  -p 2222:2222 \
  -v some_local_upload_dir:/config/uploads \
 linuxserver/openssh-server

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

.PHONY: clean_server
clean_server:
	rm -f bin/milmove

.PHONY: clean
# yarn clean removes node_modules
clean: ## Clean all generated files
	rm -f .*.stamp
	rm -f coverage.out
	rm -rf ./bin
	rm -rf ./build
	yarn clean
	rm -rf ./public/swagger-ui/*.{css,js,png}
	rm -rf ./tmp/secure_migrations
	rm -rf ./tmp/storage
	rm -rf $(SCHEMASPY_OUTPUT)
	rm -rf ./storybook-static
	rm -rf ./coverage
	rm -rf ./log

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

reviewapp_docker_destroy:
	docker-compose -f docker-compose.reviewapp.yml down

.PHONY: telemetry_docker
telemetry_docker:
	docker-compose -f docker-compose.telemetry.yml up

.PHONY: feature_flag_docker
feature_flag_docker:
	docker-compose -f docker-compose.feature_flag.yml up
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

#
# ----- START NON-ATO DEPLOYMENT TARGETS -----
#

.PHONY: nonato_deploy_prepare
nonato_deploy_prepare:  ## Replace placeholders in config to deploy to a non-ATO env. Requires DEPLOY_ENV to be set to exp, loadtest, or demo.
ifeq ($(DEPLOY_ENV), exp)
	@echo "Preparing for deploy to experimental"
else ifeq ($(DEPLOY_ENV), loadtest)
	@echo "Preparing for deploy to loadtest"
else ifeq ($(DEPLOY_ENV), demo)
	@echo "Preparing for deploy to demo"
else
	$(error DEPLOY_ENV must be exp, loadtest, or demo)
endif
	sed -E -i '' "s#(&dp3-branch) placeholder_branch_name#\1 $(GIT_BRANCH)#" .circleci/config.yml
	sed -E -i '' "s#(&integration-ignore-branch) placeholder_branch_name#\1 $(GIT_BRANCH)#" .circleci/config.yml
	sed -E -i '' "s#(&integration-mtls-ignore-branch) placeholder_branch_name#\1 $(GIT_BRANCH)#" .circleci/config.yml
	sed -E -i '' "s#(&client-ignore-branch) placeholder_branch_name#\1 $(GIT_BRANCH)#" .circleci/config.yml
	sed -E -i '' "s#(&server-ignore-branch) placeholder_branch_name#\1 $(GIT_BRANCH)#" .circleci/config.yml
	sed -E -i '' "s#(&dp3-env) placeholder_env#\1 $(DEPLOY_ENV)#" .circleci/config.yml
	@git --no-pager diff .circleci/config.yml
	@echo "Please make sure to commit the changes in .circleci/config.yml in order to have CircleCI deploy $(GIT_BRANCH) to the Non-ATO $(DEPLOY_ENV) environment."

.PHONY: nonato_deploy_restore
nonato_deploy_restore:  ## Restore placeholders in config after deploy to a non-ATO env
	sed -E -i '' "s#(&dp3-branch) $(GIT_BRANCH)#\1 placeholder_branch_name#" .circleci/config.yml
	sed -E -i '' "s#(&integration-ignore-branch) $(GIT_BRANCH)#\1 placeholder_branch_name#" .circleci/config.yml
	sed -E -i '' "s#(&integration-mtls-ignore-branch) $(GIT_BRANCH)#\1 placeholder_branch_name#" .circleci/config.yml
	sed -E -i '' "s#(&client-ignore-branch) $(GIT_BRANCH)#\1 placeholder_branch_name#" .circleci/config.yml
	sed -E -i '' "s#(&server-ignore-branch) $(GIT_BRANCH)#\1 placeholder_branch_name#" .circleci/config.yml
	sed -E -i '' "s#(&dp3-env) (exp|loadtest|demo)#\1 placeholder_env#" .circleci/config.yml

#
# ----- END NON-ATO DEPLOYMENT TARGETS -----
#

#
# ----- START SETUP MULTI BRANCH -----
#

HAS_ENVRC_LOCAL := $(shell [ -f .envrc.local ] && echo 1 || echo 0)
HAS_ENVRC_CLONED := $(shell [ -f $(CURDIR)2/.envrc.local ] && echo 1 || echo 0)

check_local_env:
ifeq (HAS_ENVRC_LOCAL,1)
	@echo "Local .envrc.local found."; \
	if [ -z "$$(grep -E '^export GIN_PORT=.*' $(CURDIR)/.envrc.local)" ]; then \
		echo "export GIN_PORT=s" >> "$(CURDIR)/.envrc.local"; \
	fi; \
	sed -i '' -e 's/^export GIN_PORT=.*/export GIN_PORT=9001/' "$(CURDIR)/.envrc.local"
else
	@echo "Local .envrc.local NOT found. Creating file."; \
	echo "export GIN_PORT=9001" > "$(CURDIR)/.envrc.local"
endif

check_cloned_env:
ifeq (HAS_ENVRC_CLONED,1)
	@echo "Cloned .envrc.local found."; \
	if [ -z "$$(grep -E '^export GIN_PORT=.*' $(CURDIR)2/.envrc.local)" ]; then \
		echo "export GIN_PORT=s" >> "$(CURDIR)2/.envrc.local"; \
	fi; \
	sed -i '' -e 's/^export GIN_PORT=.*/export GIN_PORT=9002/' "$(CURDIR)2/.envrc.local"
else
	@echo "Cloned .envrc.local NOT found. Creating file."; \
	echo "export GIN_PORT=9002" > "$(CURDIR)2/.envrc.local"
endif

clone_repo:
	@if [ -d "$(CURDIR)2" ]; then \
		echo "Error: Folder $(CURDIR)2 already exists."; \
		exit 1; \
	fi; \
	git clone https://github.com/transcom/mymove.git "$(CURDIR)2";

success_message:
	@echo "2 independent project folders created successfully."

.PHONY: multi_branch
multi_branch: check_local_env clone_repo check_cloned_env success_message ## Sets up 2 folders which can each target a different branch on the repo

#
# ----- END SETUP MULTI BRANCH -----
#

default: help

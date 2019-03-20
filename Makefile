NAME = ppp
DB_NAME_DEV = dev_db
DB_NAME_TEST = test_db
DB_DOCKER_CONTAINER_DEV = milmove-db-dev
DB_DOCKER_CONTAINER_TEST = milmove-db-test
# The version of the postgres container should match production as closely
# as possible.
# https://github.com/transcom/ppp-infra/blob/7ba2e1086ab1b2a0d4f917b407890817327ffb3d/modules/aws-app-environment/database/variables.tf#L48
DB_DOCKER_CONTAINER_IMAGE = postgres:10.5
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
DB_PORT_DEV=5432
DB_PORT_DOCKER=5432
ifdef CIRCLECI
	DB_PORT_TEST=5432
	LDFLAGS=-linkmode external -extldflags -static
endif

#
# ----- END PREAMBLE -----
#

#
# ----- START CHECK TARGETS -----
#

# This target ensures that the pre-commit hook is installed and kept up to date
# if pre-commit updates.
.PHONY: ensure_pre_commit
ensure_pre_commit: .git/hooks/pre-commit
.git/hooks/pre-commit: /usr/local/bin/pre-commit
	pre-commit install

.PHONY: prereqs
prereqs: .prereqs.stamp
.prereqs.stamp: bin/prereqs
	bin/prereqs
	touch .prereqs.stamp

.PHONY: check_hosts
check_hosts: .check_hosts.stamp
.check_hosts.stamp: bin/check-hosts-file
ifndef CIRCLECI
	bin/check-hosts-file
else
	@echo "Not checking hosts on CircleCI."
endif
	touch .check_hosts.stamp

.PHONY: go_version
go_version: .go_version.stamp
.go_version.stamp: bin/check_go_version
	bin/check_go_version
	touch .go_version.stamp

.PHONY: bash_version
bash_version: .bash_version.stamp
.bash_version.stamp: bin/check_bash_version
ifndef CIRCLECI
	bin/check_bash_version
else
	@echo "No need to check bash version on CircleCI"
endif
	touch .bash_version.stamp

.PHONY: deps
deps: prereqs check_hosts ensure_pre_commit client_deps server_deps

.PHONY: test
test: client_test server_test e2e_test

#
# ----- END CHECK TARGETS -----
#

#
# ----- START CLIENT TARGETS -----
#

.PHONY: client_deps_update
client_deps_update:
	yarn upgrade

.PHONY: client_deps
client_deps: check_hosts .client_deps.stamp
.client_deps.stamp: yarn.lock
	yarn install
	bin/copy_swagger_ui.sh
	touch .client_deps.stamp
.client_build.stamp: $(shell find src -type f)
	yarn build
	touch .client_build.stamp

.PHONY: client_build
client_build: client_deps .client_build.stamp

.PHONY: client_run
client_run: client_deps
	HOST=milmovelocal yarn start

.PHONY: client_test
client_test: client_deps
	yarn test

.PHONY: client_test_coverage
client_test_coverage : client_deps
	yarn test:coverage

.PHONY: office_client_run
office_client_run: client_deps
	HOST=officelocal yarn start

.PHONY: tsp_client_run
tsp_client_run: client_deps
	HOST=tsplocal yarn start

#
# ----- END CLIENT TARGETS -----
#

#
# ----- START SERVER TARGETS -----
#

.PHONY: go_deps_update
go_deps_update:
	dep ensure -v -update

.PHONY: go_deps
go_deps: go_version .go_deps.stamp
.go_deps.stamp: Gopkg.lock
	bin/check_gopath.sh
	dep ensure -vendor-only

.PHONY: build_chamber
build_chamber: go_deps .build_chamber.stamp
.build_chamber.stamp:
	go build -i -ldflags "$(LDFLAGS)" -o bin/chamber ./vendor/github.com/segmentio/chamber
	touch .build_chamber.stamp

.PHONY: build_soda
build_soda: go_deps .build_soda.stamp
.build_soda.stamp:
	go build -i -ldflags "$(LDFLAGS)" -o bin/soda ./vendor/github.com/gobuffalo/pop/soda
	touch .build_soda.stamp

.PHONY: build_generate_test_data
build_generate_test_data: go_deps
	go build -i -ldflags "$(LDFLAGS)" -o bin/generate-test-data ./cmd/generate_test_data

.PHONY: build_callgraph
build_callgraph: go_deps .build_callgraph.stamp
.build_callgraph.stamp:
	go build -i -o bin/callgraph ./vendor/golang.org/x/tools/cmd/callgraph
	touch .build_callgraph.stamp

.PHONY: get_goimports
get_goimports: go_deps .get_goimports.stamp
.get_goimports.stamp:
	go get -u golang.org/x/tools/cmd/goimports
	touch .get_goimports.stamp

.PHONY: download_rds_certs
download_rds_certs: .download_rds_certs.stamp
.download_rds_certs.stamp:
	curl -o bin/rds-combined-ca-bundle.pem https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem
	touch .download_rds_certs.stamp

.PHONY: server_deps
server_deps: check_hosts go_deps build_chamber build_soda build_callgraph get_goimports download_rds_certs .server_deps.stamp
.server_deps.stamp:
	# Unfortunately, dep ensure blows away ./vendor every time so these builds always take a while
	go install ./vendor/golang.org/x/lint/golint # golint needs to be accessible for the pre-commit task to run, so `install` it
	go build -i -ldflags "$(LDFLAGS)" -o bin/gosec ./vendor/github.com/securego/gosec/cmd/gosec
	go build -i -ldflags "$(LDFLAGS)" -o bin/gin ./vendor/github.com/codegangsta/gin
	go build -i -ldflags "$(LDFLAGS)" -o bin/swagger ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
	touch .server_deps.stamp

.PHONY: server_deps_linux
server_deps_linux: go_deps .server_deps_linux.stamp
.server_deps_linux.stamp:
	go build -i -ldflags "$(LDFLAGS)" -o bin/swagger ./vendor/github.com/go-swagger/go-swagger/cmd/swagger

.PHONY: server_generate
server_generate: server_deps server_go_bindata .server_generate.stamp
.server_generate.stamp: $(shell find swagger -type f -name *.yaml)
	bin/gen_server.sh
	touch .server_generate.stamp

.PHONY: server_generate_linux
server_generate_linux: server_deps_linux server_go_bindata .server_generate_linux.stamp
.server_generate_linux.stamp: $(shell find swagger -type f -name *.yaml)
	bin/gen_server.sh
	touch .server_generate_linux.stamp

.PHONY: server_go_bindata
server_go_bindata: pkg/assets/assets.go
pkg/assets/assets.go: pkg/paperwork/formtemplates/*
	go-bindata -o pkg/assets/assets.go -pkg assets pkg/paperwork/formtemplates/

.PHONY: server_build
server_build: server_deps server_generate
	go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) -i -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/webserver ./cmd/webserver

.PHONY: server_build_linux
server_build_linux: server_deps_linux server_generate_linux
	# These don't need to go in bin_linux/ because local devs don't use them
	# Additionally it would not work with the default Dockerfile
	GOOS=linux GOARCH=amd64 go build -i -ldflags "$(LDFLAGS)" -o bin/chamber ./vendor/github.com/segmentio/chamber
	GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) -i -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/webserver ./cmd/webserver

# This command is for running the server by itself, it will serve the compiled frontend on its own
server_run_standalone: client_build server_build db_dev_run
	DEBUG_LOGGING=true $(AWS_VAULT) ./bin/webserver
# This command will rebuild the swagger go code and rerun server on any changes
server_run:
	find ./swagger -type f -name "*.yaml" | entr -c -r make server_run_default
# This command runs the server behind gin, a hot-reload server
server_run_default: server_deps server_generate db_dev_run
	INTERFACE=localhost DEBUG_LOGGING=true \
	$(AWS_VAULT) ./bin/gin --build ./cmd/webserver \
		--bin /bin/webserver \
		--port 8080 --appPort 8081 \
		--excludeDir vendor --excludeDir node_modules \
		-i --buildArgs "-i -ldflags=\"$(WEBSERVER_LDFLAGS)\""

.PHONY: server_run_debug
server_run_debug:
	INTERFACE=localhost DEBUG_LOGGING=true \
	$(AWS_VAULT) dlv debug cmd/webserver/main.go

.PHONY: build_tools
build_tools: bash_version server_deps server_generate build_generate_test_data
	go build -i -ldflags "$(LDFLAGS)" -o bin/compare-secure-migrations ./cmd/compare_secure_migrations
	go build -i -ldflags "$(LDFLAGS)" -o bin/ecs-service-logs ./cmd/ecs-service-logs
	go build -i -ldflags "$(LDFLAGS)" -o bin/generate-1203-form ./cmd/generate_1203_form
	go build -i -ldflags "$(LDFLAGS)" -o bin/generate-shipment-edi ./cmd/generate_shipment_edi
	go build -i -ldflags "$(LDFLAGS)" -o bin/generate-shipment-summary ./cmd/generate_shipment_summary
	go build -i -ldflags "$(LDFLAGS)" -o bin/health_checker ./cmd/health_checker
	go build -i -ldflags "$(LDFLAGS)" -o bin/iws ./cmd/demo/iws.go
	go build -i -ldflags "$(LDFLAGS)" -o bin/load-office-data ./cmd/load_office_data
	go build -i -ldflags "$(LDFLAGS)" -o bin/load-user-gen ./cmd/load_user_gen
	go build -i -ldflags "$(LDFLAGS)" -o bin/make-dps-user ./cmd/make_dps_user
	go build -i -ldflags "$(LDFLAGS)" -o bin/make-office-user ./cmd/make_office_user
	go build -i -ldflags "$(LDFLAGS)" -o bin/make-tsp-user ./cmd/make_tsp_user
	go build -i -ldflags "$(LDFLAGS)" -o bin/paperwork ./cmd/paperwork
	go build -i -ldflags "$(LDFLAGS)" -o bin/save-fuel-price-data ./cmd/save_fuel_price_data
	go build -i -ldflags "$(LDFLAGS)" -o bin/send-to-gex ./cmd/send_to_gex
	go build -i -ldflags "$(LDFLAGS)" -o bin/tsp-award-queue ./cmd/tsp_award_queue

.PHONY: build
build: server_build build_tools client_build

# webserver_test runs a few acceptance tests against a local or remote environment.
# This can help identify potential errors before deploying a container.
.PHONY: webserver_test
webserver_test: server_generate build_chamber
ifndef TEST_ACC_ENV
	@echo "Running acceptance tests for webserver using local environment."
	@echo "* Use environment XYZ by setting environment variable to TEST_ACC_ENV=XYZ."
	TEST_ACC_HONEYCOMB=0 \
	TEST_ACC_CWD=$(PWD) \
	go test -v -p 1 -count 1 -short $$(go list ./... | grep \\/cmd\\/webserver)
else
ifndef CIRCLECI
	@echo "Running acceptance tests for webserver with environment $$TEST_ACC_ENV."
	TEST_ACC_HONEYCOMB=0 \
	TEST_ACC_CWD=$(PWD) \
	DISABLE_AWS_VAULT_WRAPPER=1 \
	aws-vault exec $(AWS_PROFILE) -- \
	bin/chamber -r $(CHAMBER_RETRIES) exec app-$(TEST_ACC_ENV) -- \
	go test -v -p 1 -count 1 -short $$(go list ./... | grep \\/cmd\\/webserver)
else
	@echo "Running acceptance tests for webserver with environment $$TEST_ACC_ENV."
	TEST_ACC_HONEYCOMB=0 \
	TEST_ACC_CWD=$(PWD) \
	bin/chamber -r $(CHAMBER_RETRIES) exec app-$(TEST_ACC_ENV) -- \
	go test -v -p 1 -count 1 -short $$(go list ./... | grep \\/cmd\\/webserver)
endif
endif

.PHONY: server_test
server_test: server_deps server_generate db_test_reset db_test_migrate
	# Don't run tests in /cmd or /pkg/gen & pass `-short` to exclude long running tests
	# Use -test.parallel 1 to test packages serially and avoid database collisions
	# Disable test caching with `-count 1` - caching was masking local test failures
	DB_PORT=$(DB_PORT_TEST) go test -p 1 -count 1 -short $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

server_test_all: server_deps server_generate db_dev_reset db_dev_migrate
	# Like server_test but runs extended tests that may hit external services.
	DB_PORT=$(DB_PORT_TEST) go test -p 1 -count 1 $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

server_test_coverage_generate: server_deps server_generate db_test_reset db_test_migrate
	# Don't run tests in /cmd or /pkg/gen
	# Use -test.parallel 1 to test packages serially and avoid database collisions
	# Disable test caching with `-count 1` - caching was masking local test failures
	# Add coverage tracker via go cover
	DB_PORT=$(DB_PORT_TEST) go test -coverprofile=coverage.out -covermode=count -p 1 -count 1 -short $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

server_test_coverage: server_deps server_generate db_test_reset db_test_migrate server_test_coverage_generate
	DB_PORT=$(DB_PORT_TEST) go tool cover -html=coverage.out

#
# ----- END SERVER TARGETS -----
#

#
# ----- START DB_DEV TARGETS -----
#

.PHONY: db_dev_destroy
db_dev_destroy:
ifndef CIRCLECI
	@echo "Destroying the ${DB_DOCKER_CONTAINER_DEV} docker database container..."
	docker rm -f $(DB_DOCKER_CONTAINER_DEV) || \
		echo "No database container"
else
	@echo "Relying on CircleCI's database setup to destroy the DB."
endif

.PHONY: db_dev_start
db_dev_start:
ifndef CIRCLECI
	brew services stop postgresql 2> /dev/null || true
endif
	@echo "Starting the ${DB_DOCKER_CONTAINER_DEV} docker database container..."
	# If running do nothing, if not running try to start, if can't start then run
	docker start $(DB_DOCKER_CONTAINER_DEV) || \
		docker run --name $(DB_DOCKER_CONTAINER_DEV) \
			-e \
			POSTGRES_PASSWORD=$(PGPASSWORD) \
			-d \
			-p $(DB_PORT_DEV):$(DB_PORT_DOCKER)\
			$(DB_DOCKER_CONTAINER_IMAGE)

.PHONY: db_dev_create
db_dev_create:
	@echo "Create the ${DB_NAME_DEV} database..."
	DB_NAME=postgres bin/wait-for-db && \
		createdb -p $(DB_PORT_DEV) -h localhost -U postgres $(DB_NAME_DEV) || true

.PHONY: db_dev_run
db_dev_run: db_dev_start db_dev_create

.PHONY: db_dev_reset
db_dev_reset: db_dev_destroy db_dev_run

.PHONY: db_dev_migrate_standalone
db_dev_migrate_standalone:
	@echo "Migrating the ${DB_NAME_DEV} database..."
	# We need to move to the bin/ directory so that the cwd contains `apply-secure-migration.sh`
	cd bin && \
		./soda -c ../config/database.yml -p ../migrations migrate up

.PHONY: db_dev_migrate
db_dev_migrate: server_deps db_dev_migrate_standalone

#
# ----- END DB_DEV TARGETS -----
#

#
# ----- START DB_TEST TARGETS -----
#

.PHONY: db_test_destroy
db_test_destroy:
ifndef CIRCLECI
	@echo "Destroying the ${DB_DOCKER_CONTAINER_TEST} docker database container..."
	docker rm -f $(DB_DOCKER_CONTAINER_TEST) || \
		echo "No database container"
else
	@echo "Relying on CircleCI's database setup to destroy the DB."
endif

.PHONY: db_test_start
db_test_start:
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
db_test_create:
ifndef CIRCLECI
	@echo "Create the ${DB_NAME_TEST} database..."
	DB_NAME=postgres DB_PORT=$(DB_PORT_TEST) bin/wait-for-db && \
		createdb -p $(DB_PORT_TEST) -h localhost -U postgres $(DB_NAME_TEST) || true
else
	@echo "Relying on CircleCI's database setup to create the DB."
endif

.PHONY: db_test_create_docker
db_test_create_docker:
	@echo "Create the ${DB_NAME_TEST} database with docker command..."
	DB_NAME=postgres DB_DOCKER_CONTAINER=$(DB_DOCKER_CONTAINER_TEST) bin/wait-for-db-docker && \
		docker exec $(DB_DOCKER_CONTAINER_TEST) createdb -p $(DB_PORT_DOCKER) -h localhost -U postgres $(DB_NAME_TEST) || true

.PHONY: db_test_run
db_test_run: db_test_start db_test_create

.PHONY: db_test_run_docker
db_test_run_docker: db_test_start db_test_create_docker

.PHONY: db_test_reset
db_test_reset: db_test_destroy db_test_run

.PHONY: db_test_reset_docker
db_test_reset_docker: db_test_destroy db_test_run_docker

.PHONY: db_test_migrate_standalone
db_test_migrate_standalone:
ifndef CIRCLECI
	@echo "Migrating the ${DB_NAME_TEST} database..."
	# We need to move to the bin/ directory so that the cwd contains `apply-secure-migration.sh`
	cd bin && \
		DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_TEST)\
			./soda -c ../config/database.yml -p ../migrations migrate up
else
	@echo "Migrating the ${DB_NAME_TEST} database..."
	# We need to move to the bin/ directory so that the cwd contains `apply-secure-migration.sh`
	cd bin && \
		DB_NAME=$(DB_NAME_TEST) DB_PORT=$(DB_PORT_DEV) \
			./soda -c ../config/database.yml -p ../migrations migrate up
endif

.PHONY: db_test_migrate
db_test_migrate: server_deps db_test_migrate_standalone

.PHONY: db_test_migrations_build
db_test_migrations_build: .db_test_migrations_build.stamp
.db_test_migrations_build.stamp: server_deps_linux server_generate_linux
	@echo "Build required binaries for the docker migration container..."
	mkdir -p bin_linux/
	GOOS=linux GOARCH=amd64 go build -i -ldflags "$(LDFLAGS)" -o bin_linux/soda ./vendor/github.com/gobuffalo/pop/soda
	GOOS=linux GOARCH=amd64 go build -i -ldflags "$(LDFLAGS)" -o bin_linux/generate-test-data ./cmd/generate_test_data
	@echo "Build the docker migration container..."
	docker build -f Dockerfile.migrations_local --tag e2e_migrations:latest .

.PHONY: db_test_migrate_docker
db_test_migrate_docker: db_test_migrations_build
	@echo "Migrating the ${DB_NAME_TEST} database with docker command..."
	DB_NAME=$(DB_NAME_TEST) DB_DOCKER_CONTAINER=$(DB_DOCKER_CONTAINER_TEST) bin/wait-for-db-docker
	docker run \
		-t \
		-e DB_NAME=$(DB_NAME_TEST) \
		-e DB_HOST=database \
		-e DB_PORT=$(DB_PORT_DOCKER) \
		-e DB_USER=postgres \
		-e DB_PASSWORD=$(PGPASSWORD) \
		--link="$(DB_DOCKER_CONTAINER_TEST):database" \
		--rm \
		--entrypoint soda \
		e2e_migrations:latest \
		migrate -c /migrate/database.yml -p /migrate/migrations up

#
# ----- END DB_TEST TARGETS -----
#

#
# ----- START E2E TARGETS -----
#

.PHONY: e2e_test
e2e_test: server_deps server_generate server_build client_build db_e2e_init
	$(AWS_VAULT) ./bin/run-e2e-test

.PHONY: e2e_test_docker
e2e_test_docker:
	$(AWS_VAULT) ./bin/run-e2e-test-docker

.PHONY: e2e_clean
e2e_clean:
	rm -f .*_linux.stamp
	rm -f .db_test_migrations_build.stamp
	rm -rf cypress/results
	rm -rf cypress/screenshots
	rm -rf cypress/videos
	rm -rf bin_linux/
	docker rm -f cypress || true
	docker rm -f e2e || true
	docker rm -f e2e_migrations || true

.PHONY: db_e2e_up
db_e2e_up: build_generate_test_data
	@echo "Truncate the ${DB_NAME_TEST} database..."
	psql postgres://postgres:$(PGPASSWORD)@localhost:$(DB_PORT_TEST)/$(DB_NAME_TEST)?sslmode=disable -c 'TRUNCATE users CASCADE;'
	@echo "Populate the ${DB_NAME_TEST} database..."
	DB_PORT=$(DB_PORT_TEST) bin/generate-test-data -named-scenario="e2e_basic" -env="test"

.PHONY: db_e2e_up_docker
db_e2e_up_docker:
	@echo "Truncate the ${DB_NAME_TEST} database with docker command..."
	docker run \
		--link="$(DB_DOCKER_CONTAINER_TEST):database" \
		--rm \
		--entrypoint psql \
		e2e_migrations:latest \
		postgres://postgres:$(PGPASSWORD)@database:$(DB_PORT_DOCKER)/$(DB_NAME_TEST)?sslmode=disable -c 'TRUNCATE users CASCADE;'
	@echo "Populate the ${DB_NAME_TEST} database with docker command..."
	docker run \
		-t \
		-e DB_NAME=$(DB_NAME_TEST) \
		-e DB_HOST=database \
		-e DB_PORT=$(DB_PORT_DOCKER) \
		-e DB_USER=postgres \
		-e DB_PASSWORD=$(PGPASSWORD) \
		--link="$(DB_DOCKER_CONTAINER_TEST):database" \
		--rm \
		--entrypoint generate-test-data \
		e2e_migrations:latest \
		-config-dir /migrate -named-scenario e2e_basic

.PHONY: db_e2e_init
db_e2e_init: db_test_reset db_test_migrate db_e2e_up

.PHONY: db_e2e_init_docker
db_e2e_init_docker: db_test_reset_docker db_test_migrate_docker db_e2e_up_docker

.PHONY: db_dev_e2e_populate
db_dev_e2e_populate: db_dev_reset db_dev_migrate build_tools
	@echo "Populate the ${DB_NAME_DEV} database with docker command..."
	bin/generate-test-data -named-scenario="e2e_basic" -env="development"

.PHONY: db_test_e2e_populate
db_test_e2e_populate: db_test_reset_docker db_test_migrate_docker build_tools db_e2e_up_docker

#
# ----- END E2E TARGETS -----
#

#
# ----- START RANDOM TARGETS -----
#

.PHONY: 1203_form
1203_form:
	find ./cmd/generate_1203_form -type f -name "main.go" | entr -c -r go run ./cmd/generate_1203_form/main.go

.PHONY: adr_update
adr_update:
	yarn run adr-log

.PHONY: tsp_run
tsp_run: build_tools db_dev_run
	./bin/tsp-award-queue

.PHONY: pre_commit_tests
pre_commit_tests:
	pre-commit run --all-files

.PHONY: pretty
pretty:
	npx prettier --write --loglevel warn "src/**/*.{js,jsx}"
	gofmt pkg/ >> /dev/null

.PHONY: clean
clean:
	rm -f .*.stamp
	rm -rf ./node_modules
	rm -rf ./vendor
	rm -rf ./pkg/gen
	rm -rf ./public/swagger-ui/*.{css,js,png}
	rm -rf $$GOPATH/pkg/dep/sources

.PHONY: spellcheck
spellcheck:
	node_modules/.bin/mdspell --ignore-numbers --ignore-acronyms --en-us \
		`find . -type f -name "*.md" \
			-not -path "./vendor/*" \
			-not -path "./node_modules/*" \
			-not -path "./docs/adr/index.md"`

#
# ----- END RANDOM TARGETS -----
#

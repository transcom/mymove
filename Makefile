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

# Convenience for LDFLAGS
WEBSERVER_LDFLAGS=-X main.gitBranch=$(shell git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(shell git rev-list -1 HEAD)
ifndef CIRCLECI
	LDFLAGS=
else
	LDFLAGS=-linkmode external -extldflags -static
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

check_hosts: .check_hosts.stamp
.check_hosts.stamp: bin/check-hosts-file
ifndef CIRCLECI
	bin/check-hosts-file
else
	@echo "Not checking hosts on CircleCI."
endif
	touch .check_hosts.stamp

go_version: .go_version.stamp
.go_version.stamp: bin/check_go_version
	bin/check_go_version
	touch .go_version.stamp

deps: prereqs check_hosts ensure_pre_commit client_deps server_deps
test: client_test server_test e2e_test

spellcheck:
	node_modules/.bin/mdspell --ignore-numbers --ignore-acronyms --en-us \
		`find . -type f -name "*.md" \
			-not -path "./vendor/*" \
			-not -path "./node_modules/*" \
			-not -path "./docs/adr/index.md"`

client_deps_update:
	yarn upgrade
client_deps: check_hosts .client_deps.stamp
.client_deps.stamp: yarn.lock
	yarn install
	bin/copy_swagger_ui.sh
	touch .client_deps.stamp
.client_build.stamp: $(shell find src -type f)
	yarn build
	touch .client_build.stamp
client_build: client_deps .client_build.stamp
client_run: client_deps
	HOST=milmovelocal yarn start
client_test: client_deps
	yarn test
client_test_coverage : client_deps
	yarn test:coverage

office_client_run: client_deps
	HOST=officelocal yarn start

tsp_client_run: client_deps
	HOST=tsplocal yarn start

go_deps_update:
	dep ensure -v -update

go_deps: go_version .go_deps.stamp
.go_deps.stamp: Gopkg.lock
	bin/check_gopath.sh
	dep ensure -vendor-only

build_chamber: go_deps .build_chamber.stamp
.build_chamber.stamp:
	go build -i -ldflags "$(LDFLAGS)" -o bin/chamber ./vendor/github.com/segmentio/chamber
	touch .build_chamber.stamp

build_soda: go_deps .build_soda.stamp
.build_soda.stamp:
	go build -i -ldflags "$(LDFLAGS)" -o bin/soda ./vendor/github.com/gobuffalo/pop/soda
	touch .build_soda.stamp

server_deps: check_hosts go_deps build_chamber build_soda .server_deps.stamp
.server_deps.stamp:
	# Unfortunately, dep ensure blows away ./vendor every time so these builds always take a while
	go install ./vendor/github.com/golang/lint/golint # golint needs to be accessible for the pre-commit task to run, so `install` it
	go build -i -ldflags "$(LDFLAGS)" -o bin/gosec ./vendor/github.com/securego/gosec/cmd/gosec
	go build -i -ldflags "$(LDFLAGS)" -o bin/gin ./vendor/github.com/codegangsta/gin
	go build -i -ldflags "$(LDFLAGS)" -o bin/swagger ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
	touch .server_deps.stamp
server_deps_linux: go_deps .server_deps_linux.stamp
.server_deps_linux.stamp:
	go build -i -ldflags "$(LDFLAGS)" -o bin/swagger ./vendor/github.com/go-swagger/go-swagger/cmd/swagger

server_generate: server_deps server_go_bindata .server_generate.stamp
.server_generate.stamp: $(shell find swagger -type f -name *.yaml)
	bin/gen_server.sh
	touch .server_generate.stamp
server_generate_linux: server_deps_linux server_go_bindata .server_generate_linux.stamp
.server_generate_linux.stamp: $(shell find swagger -type f -name *.yaml)
	bin/gen_server.sh
	touch .server_generate_linux.stamp

server_go_bindata: pkg/assets/assets.go
pkg/assets/assets.go: pkg/paperwork/formtemplates/*
	go-bindata -o pkg/assets/assets.go -pkg assets pkg/paperwork/formtemplates/

server_build: server_deps server_generate
	go build -gcflags=-trimpath=$(GOPATH) -asmflags=-trimpath=$(GOPATH) -i -ldflags "$(LDFLAGS) $(WEBSERVER_LDFLAGS)" -o bin/webserver ./cmd/webserver
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

server_run_debug:
	INTERFACE=localhost DEBUG_LOGGING=true \
	$(AWS_VAULT) dlv debug cmd/webserver/main.go

build_tools: server_deps server_generate
	go build -i -ldflags "$(LDFLAGS)" -o bin/generate-1203-form ./cmd/generate_1203_form
	go build -i -ldflags "$(LDFLAGS)" -o bin/generate-shipment-summary ./cmd/generate_shipment_summary
	go build -i -ldflags "$(LDFLAGS)" -o bin/generate-test-data ./cmd/generate_test_data
	go build -i -ldflags "$(LDFLAGS)" -o bin/health_checker ./cmd/health_checker
	go build -i -ldflags "$(LDFLAGS)" -o bin/iws ./cmd/demo/iws.go
	go build -i -ldflags "$(LDFLAGS)" -o bin/load-office-data ./cmd/load_office_data
	go build -i -ldflags "$(LDFLAGS)" -o bin/load-user-gen ./cmd/load_user_gen
	go build -i -ldflags "$(LDFLAGS)" -o bin/make-dps-user ./cmd/make_dps_user
	go build -i -ldflags "$(LDFLAGS)" -o bin/make-office-user ./cmd/make_office_user
	go build -i -ldflags "$(LDFLAGS)" -o bin/make-tsp-user ./cmd/make_tsp_user
	go build -i -ldflags "$(LDFLAGS)" -o bin/paperwork ./cmd/paperwork
	go build -i -ldflags "$(LDFLAGS)" -o bin/tsp-award-queue ./cmd/tsp_award_queue

tsp_run: build_tools db_dev_run
	./bin/tsp-award-queue

build: server_build build_tools client_build

# webserver_test runs a few acceptance tests against a local or remote environment.
# This can help identify potential errors before deploying a container.
webserver_test: server_generate build_chamber
ifndef TEST_ACC_ENV
	@echo "Running acceptance tests for webserver using local environment."
	@echo "* Use environment XYZ by setting environment variable to TEST_ACC_ENV=XYZ."
	TEST_ACC_DATABASE=0 TEST_ACC_HONEYCOMB=0 \
	go test -v -p 1 -count 1 -short $$(go list ./... | grep \\/cmd\\/webserver)
else
ifndef CIRCLECI
	@echo "Running acceptance tests for webserver with environment $$TEST_ACC_ENV."
	TEST_ACC_DATABASE=0 TEST_ACC_HONEYCOMB=0 \
	TEST_ACC_CWD=$(PWD) \
	DISABLE_AWS_VAULT_WRAPPER=1 \
	aws-vault exec $(AWS_PROFILE) -- \
	bin/chamber exec app-$(TEST_ACC_ENV) -- \
	go test -v -p 1 -count 1 -short $$(go list ./... | grep \\/cmd\\/webserver)
else
	@echo "Running acceptance tests for webserver with environment $$TEST_ACC_ENV."
	TEST_ACC_DATABASE=0 TEST_ACC_HONEYCOMB=0 \
	TEST_ACC_CWD=$(PWD) \
	bin/chamber exec app-$(TEST_ACC_ENV) -- \
	go test -v -p 1 -count 1 -short $$(go list ./... | grep \\/cmd\\/webserver)
endif
endif

server_test: server_deps server_generate db_test_reset db_test_migrate
	# Don't run tests in /cmd or /pkg/gen & pass `-short` to exclude long running tests
	# Use -test.parallel 1 to test packages serially and avoid database collisions
	# Disable test caching with `-count 1` - caching was masking local test failures
	go test -p 1 -count 1 -short $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

server_test_all: server_deps server_generate db_dev_reset db_dev_migrate
	# Like server_test but runs extended tests that may hit external services.
	go test -p 1 -count 1 $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)

server_test_coverage: server_deps server_generate db_dev_reset db_dev_migrate
	# Don't run tests in /cmd or /pkg/gen
	# Use -test.parallel 1 to test packages serially and avoid database collisions
	# Disable test caching with `-count 1` - caching was masking local test failures
	# Add coverage tracker via go cover
	# Then open coverage tracker in HTML
	go test -coverprofile=coverage.out -p 1 -count 1 $$(go list ./... | grep -v \\/pkg\\/gen\\/ | grep -v \\/cmd\\/)
	go tool cover -html=coverage.out

e2e_test: server_deps server_generate server_build client_build db_e2e_init
	$(AWS_VAULT) ./bin/run-e2e-test

e2e_test_docker:
	$(AWS_VAULT) ./bin/run-e2e-test-docker

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

db_run:
ifndef CIRCLECI
	brew services stop postgresql 2> /dev/null || true
endif
	@echo "Starting the local docker database container..."
	# The version of the postgres container should match production as closely
	# as possible.
	# https://github.com/transcom/ppp-infra/blob/7ba2e1086ab1b2a0d4f917b407890817327ffb3d/modules/aws-app-environment/database/variables.tf#L48
	docker start $(DB_DOCKER_CONTAINER) || \
		docker run --name $(DB_DOCKER_CONTAINER) \
			-e \
			POSTGRES_PASSWORD=$(PGPASSWORD) \
			-d \
			-p 5432:5432 \
			postgres:10.5

db_destroy:
ifndef CIRCLECI
	@echo "Destroying the local docker database container..."
	docker rm -f $(DB_DOCKER_CONTAINER) || \
		echo "No database container"
else
	@echo "Relying on CircleCI's database setup to destroy the DB."
endif

db_dev_create:
	@echo "Create the dev database..."
	DB_NAME=postgres bin/wait-for-db && \
		createdb -p 5432 -h localhost -U postgres dev_db || true

db_dev_run: db_run db_dev_create

db_dev_reset: db_destroy db_dev_run

db_dev_migrate: server_deps
	@echo "Migrating the ${DB_NAME} database..."
	# We need to move to the bin/ directory so that the cwd contains `apply-secure-migration.sh`
	cd bin && \
		./soda -c ../config/database.yml -p ../migrations migrate up

db_test_create:
	@echo "Create the test database..."
	DB_NAME=postgres bin/wait-for-db && \
		createdb -p 5432 -h localhost -U postgres test_db || true

db_test_run: db_run db_test_create

db_test_create_docker:
	@echo "Create the test database with docker command..."
	DB_NAME=postgres bin/wait-for-db-docker && \
		docker exec $(DB_DOCKER_CONTAINER) createdb -p 5432 -h localhost -U postgres test_db || true

db_test_run_docker: db_run db_test_create_docker

db_test_migrations_build: .db_test_migrations_build.stamp
.db_test_migrations_build.stamp: server_deps_linux server_generate_linux
	@echo "Build required binaries for the docker migration container..."
	mkdir -p bin_linux/
	GOOS=linux GOARCH=amd64 go build -i -ldflags "$(LDFLAGS)" -o bin_linux/soda ./vendor/github.com/gobuffalo/pop/soda
	GOOS=linux GOARCH=amd64 go build -i -ldflags "$(LDFLAGS)" -o bin_linux/generate-test-data ./cmd/generate_test_data
	@echo "Build the docker migration container..."
	docker build -f Dockerfile.migrations_local --tag e2e_migrations:latest .

db_test_migrate: server_deps
	@echo "Migrating the test database..."
	# We need to move to the bin/ directory so that the cwd contains `apply-secure-migration.sh`
	cd bin && \
		DB_NAME=test_db \
		./soda -c ../config/database.yml -p ../migrations migrate up

db_test_migrate_docker: db_test_migrations_build
	@echo "Migrating the test database with docker command..."
	DB_NAME=test_db bin/wait-for-db-docker
	docker run \
		-t \
		-e DB_NAME=test_db \
		-e DB_HOST=database \
		-e DB_PORT=5432 \
		-e DB_USER=postgres \
		-e DB_PASSWORD=$(PGPASSWORD) \
		--link="$(DB_DOCKER_CONTAINER):database" \
		--rm \
		--entrypoint soda \
		e2e_migrations:latest \
		migrate -c /migrate/database.yml -p /migrate/migrations up

db_test_reset: db_destroy db_test_run

db_test_reset_docker: db_destroy db_test_run_docker

db_e2e_up:
	@echo "Truncate the test database..."
	psql postgres://postgres:$(PGPASSWORD)@localhost:5432/test_db?sslmode=disable -c 'TRUNCATE users CASCADE;'
	@echo "Populate the test database..."
	bin/generate-test-data -named-scenario="e2e_basic" -env="test"

db_e2e_up_docker:
	@echo "Truncate the test database with docker command..."
	docker run \
		--link="$(DB_DOCKER_CONTAINER):database" \
		--rm \
		--entrypoint psql \
		e2e_migrations:latest \
		postgres://postgres:$(PGPASSWORD)@database:5432/test_db?sslmode=disable -c 'TRUNCATE users CASCADE;'
	@echo "Populate the test database with docker command..."
	docker run \
		-t \
		-e DB_NAME=test_db \
		-e DB_HOST=database \
		-e DB_PORT=5432 \
		-e DB_USER=postgres \
		-e DB_PASSWORD=$(PGPASSWORD) \
		--link="$(DB_DOCKER_CONTAINER):database" \
		--rm \
		--entrypoint generate-test-data \
		e2e_migrations:latest \
		-config-dir /migrate -named-scenario e2e_basic

db_e2e_init: db_test_reset db_test_migrate db_e2e_up

db_e2e_init_docker: db_test_reset_docker db_test_migrate_docker db_e2e_up_docker

db_e2e_reset: db_e2e_init
	@echo "\033[0;31mUse 'make db_e2e_init' instead please\033[0m"

db_dev_e2e_populate: db_dev_reset db_dev_migrate build_tools
	@echo "Populate the dev database with docker command..."
	bin/generate-test-data -named-scenario="e2e_basic" -env="development"

db_test_e2e_populate: db_test_reset_docker db_test_migrate_docker build_tools
	@echo "Populate the test database with docker command..."
	bin/generate-test-data -named-scenario="e2e_basic" -env="test"

# Backwards compatibility
db_populate_e2e: db_dev_e2e_populate
	@echo "\033[0;31mUse 'make db_dev_e2e_populate' instead please\033[0m"

1203_form:
	find ./cmd/generate_1203_form -type f -name "main.go" | entr -c -r go run ./cmd/generate_1203_form/main.go

adr_update:
	yarn run adr-log

pre_commit_tests:
	pre-commit run --all-files

pretty:
	npx prettier --write --loglevel warn "src/**/*.{js,jsx}"
	gofmt pkg/ >> /dev/null

clean:
	rm -f .*.stamp
	rm -rf ./node_modules
	rm -rf ./vendor
	rm -rf ./pkg/gen
	rm -rf ./public/swagger-ui/*.{css,js,png}
	rm -rf $$GOPATH/pkg/dep/sources

.PHONY: pre-commit deps test client_deps client_build client_run client_test prereqs check_hosts
.PHONY: server_run_standalone server_run server_run_default server_test webserver_test
.PHONY: go_deps_update server_go_bindata
.PHONY: build_chamber build_soda
.PHONY: server_generate server_deps server_build
.PHONY: server_generate_linux server_deps_linux server_build_linux
.PHONY: db_run db_destroy
.PHONY: db_dev_run db_dev_create db_dev_reset db_dev_migrate db_dev_e2e_populate
.PHONY: db_test_run db_test_create db_test_reset db_test_migrate db_test_e2e_populate
.PHONY: db_test_run_docker db_test_create_docker db_test_reset_docker db_test_migrations_build db_test_migrate_docker
.PHONY: db_populate_e2e db_e2e_up db_e2e_up_docker db_e2e_init db_e2e_init_docker db_e2e_reset
.PHONY: e2e_test e2e_test_docker e2e_clean
.PHONY: clean pretty

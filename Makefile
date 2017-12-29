NAME = ppp
export GOPATH = $(CURDIR)/server
export GOBIN = $(GOPATH)/bin

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
	cd server/src/dp3 && glide install
server_build_only:
	cd server/src/dp3/cmd/webserver && \
	go install
server_build: server_deps server_build_only
server_run_only:
	./server/bin/webserver \
		-entry client/build/index.html \
		-static client/build/static \
		-port :8080 \
		-debug_logging
server_run: server_build client_build server_run_only
server_run_dev: server_build_only server_run_only
server_test:
	go test -v dp3/pkg/api

.PHONY: pre-commit deps

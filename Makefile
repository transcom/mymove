NAME = ppp

# This target ensures that the pre-commit hook is installed and kept up to date
# if pre-commit updates.
pre-commit: .git/hooks/pre-commit
.git/hooks/pre-commit: /usr/local/bin/pre-commit
	pre-commit install

client_deps:
	cd client && \
	yarn --frozen-lockfile install
client_build: client_deps
	cd client && \
	yarn build
client_run_dev:
	cd client && \
	yarn start
client_run: client_run_dev


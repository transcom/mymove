#! /usr/bin/env bash
# Executes an SQL file from S3 against the environment's database.
#
# If `SECURE_MIGRATION_SOURCE=local` then we look for a similarly named file in the
# local repository, instead of pulling from S3.

if [ -z "${SECURE_MIGRATION_SOURCE:-}" ]; then
  echo "error: \$SECURE_MIGRATION_SOURCE needs to be set"
  exit 1
fi

if [ -z "${DB_USER:-}" ]; then
  echo "error: \$DB_USER needs to be set"
  exit 1
fi

if [ -z "${DB_PASSWORD:-}" ]; then
  echo "error: \$DB_PASSWORD needs to be set"
  exit 1
fi

# -e : immediately exit script on any error (exit code > 0)
# -u : treat unset variables as errors
# -x : print each command to stderr before running
# -o pipefail : have pipe propogate exit code to end
set -euxo pipefail

readonly migration_file="${1:-}"

case $SECURE_MIGRATION_SOURCE in
  local)

    sslmode=${PSQL_SSL_MODE:-prefer}

    if [ -z "${SECURE_MIGRATION_DIR:-}" ]; then
      echo "error: \$SECURE_MIGRATION_DIR needs to be set"
      exit 1
    fi

    echo "Running secure migrations from local source..."

    readonly migration="${SECURE_MIGRATION_DIR}/${migration_file}"

    if [ ! -f "$migration" ]; then
      echo "Migration file not found: $migration"
      exit 1
    fi

    echo "Applying secure migrations from local filesystem with file ${migration_file}"

    # +x : don't print commands to stderr anymore
    # Don't share the database password
    set +x

    # Use pipe for local source like we do for S3 source
    # Run the migrations file with the following options:
    # - The migration is wrapped in a single transaction
    # - Any errors in the migration file will cause a failure
    # shellcheck disable=SC2002
    cat "$migration" | psql \
      --single-transaction \
      --variable "ON_ERROR_STOP=1" \
      "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${sslmode}" > /dev/null

    ;;
  s3)

    sslmode=${PSQL_SSL_MODE:-require}

    if [ -z "${AWS_S3_BUCKET_NAME:-}" ]; then
      echo "error: \$AWS_S3_BUCKET_NAME needs to be set"
      exit 1
    fi

    readonly url="s3://${AWS_S3_BUCKET_NAME}/secure-migrations/$1"

    echo "Applying secure migrations from S3 using url $url"

    # +x : don't print commands to stderr anymore
    # Don't share the database password
    set +x

    # Use pipe for local source like we do for S3 source
    # Run the migrations file with the following options:
    # - The migration is wrapped in a single transaction
    # - Any errors in the migration file will cause a failure
    aws s3 cp --quiet "$url" - | psql \
    --single-transaction \
    --variable "ON_ERROR_STOP=1" \
    "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${sslmode}" > /dev/null

    ;;
  *)

    echo "Unknown migration source (${SECURE_MIGRATION_SOURCE}). Exiting."
    exit 1

    ;;
esac

#!/bin/sh
# Executes an SQL file from S3 against the environment's database.
#
# If `ENVIRONMENT=devlocal` then we look for a similarly named file in the
# local repository, instead of pulling from S3.


# sh doesn't have `-o pipefail`
set -eux

readonly migration_file="${1:-}"
psql_ssl_mode=""

if [ -z "${SECURE_MIGRATION_DIR:-}" ]; then
  echo "error: \$SECURE_MIGRATION_DIR needs to be set"
  exit 1
fi

if [ -z "${SECURE_MIGRATION_SOURCE:-}" ]; then
  echo "error: \$SECURE_MIGRATION_SOURCE needs to be set"
  exit 1
fi

download_migration_from_s3() {
  readonly url="s3://${AWS_S3_BUCKET_NAME}/secure-migrations/$1"
  echo "Downloading from S3: $url"
  aws s3 cp "${url}" "${SECURE_MIGRATION_DIR}/${migration_file}"
}

case $SECURE_MIGRATION_SOURCE in
  local)
    echo "Running secure migrations from local source..."
    ;;
  s3)
    echo "Running secure migrations from S3..."
    download_migration_from_s3 "$migration_file"
    psql_ssl_mode="?sslmode=require"
    ;;
  *)
    echo "Unknown migration source (${SECURE_MIGRATION_SOURCE}). Exiting."
    exit 1
    ;;
esac

readonly migration="${SECURE_MIGRATION_DIR}/${migration_file}"

if [ ! -f "$migration" ]; then
  echo "Migration file not found: $migration"
  exit 1
fi

echo "Applying secure migrations: ${migration_file}"

# Don't share the database password
set +x

# Run the migrations file with the following options:
# - The migration is wrapped in a single transaction
# - Any errors in the migration file will cause a failure
psql \
  --single-transaction \
  --variable "ON_ERROR_STOP=1" \
  --file="$migration" \
  postgres://"${DB_USER}":"$DB_PASSWORD"@"$DB_HOST":"$DB_PORT"/"$DB_NAME""$psql_ssl_mode"

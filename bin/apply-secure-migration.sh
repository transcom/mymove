#!/bin/sh
# Executes an SQL file from S3 against the environment's database.
#
# If `ENVIRONMENT=devlocal` then we look for a similarly named file in the
# local repository, instead of pulling from S3.


# sh doesn't have `-o pipefail`
set -eux

readonly migration_file="${1:-}"
ssl_mode=""

if [ -z "${SECURE_MIGRATION_DIR:-}" ]; then
  echo "error: \$SECURE_MIGRATION_DIR needs to be set"
  exit 1
fi

download_migration_from_s3() {
  readonly url="s3://${AWS_S3_BUCKET_NAME}/secure-migrations/$1"
  echo "Downloading from S3: $url"
  aws s3 cp "${url}" "${SECURE_MIGRATION_DIR}/${migration_file}"
}

case $ENVIRONMENT in
  devlocal)
    echo "Running in local development environment..."
    ;;
  *)
    echo "Running in a CI environment..."
    download_migration_from_s3 "$migration_file"
    ssl_mode="?sslmode=require"
    ;;
esac

readonly migration="${SECURE_MIGRATION_DIR}/${migration_file}"

if [ ! -f "$migration" ]; then
  echo "Migration file not found: $migration"
  exit 1
fi

echo "Applying secure migrations: ${migration_file}"

psql postgres://"${DB_USER}":"$DB_PASSWORD"@"$DB_HOST":"$DB_PORT"/"$DB_NAME""$ssl_mode" < "$migration"

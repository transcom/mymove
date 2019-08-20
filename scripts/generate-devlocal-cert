#! /usr/bin/env bash
# Convenience script for creating a new certificate signed by the devlocal CA.

set -eo pipefail

function usage() {
  echo "Usage: $0 -f <filename prefix> [-d <output dir>] [-c <country>] [-s <state>] [-l <city>] -o <organization> -u <organizational unit> [-n <common name>]"
  echo "If omitted, the optional parameters are set as follows:"
  echo "  -d (the current working directory)"
  echo "  -c US"
  echo "  -s DC"
  echo "  -l Washington"
  echo "  -n localhost"
  exit 1
}

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
TLSDIR="$( cd "$SCRIPTDIR/../config/tls" >/dev/null 2>&1 && pwd )"
OUTDIR="$( pwd )"

# CA files
DEVLOCAL_CA_PEM="$TLSDIR/devlocal-ca.pem"
DEVLOCAL_CA_KEY="$TLSDIR/devlocal-ca.key"

# Default subject fields
C="US"
ST="DC"
L="Washington"
CN="localhost"

while getopts "?c:d:f:l:n:o:s:u:" opt; do
  case ${opt} in
    c ) # C (country) field in certificate subject
      C=$OPTARG
      ;;
    d ) # Output directory
      OUTDIR=$OPTARG
      ;;
    f ) # Filename prefix
      PREFIX=$OPTARG
      ;;
    l ) # L (locality) field in certificate subject
      L=$OPTARG
      ;;
    n ) # CN (common name) field in certificate subject
      CN=$OPTARG
      ;;
    o ) # O (organization) field in certificate subject
      O=$OPTARG
      ;;
    s ) # ST (state) field in certificate subject
      ST=$OPTARG
      ;;
    u ) # OU (organizational unit) field in certificate subject
      OU=$OPTARG
      ;;
    \? )
      usage
      ;;
  esac
done

# Before using set -u to detect unset variables, check that the command line
# options are set explicitly, allowing us to show a user-friendly usage message
# if one or more are missing
if [[ -z $PREFIX || -z $OUTDIR || -z $C || -z $ST || -z $L || -z $O || -z $OU || -z $CN ]]; then
  usage
fi

set -u

openssl req -nodes -new -keyout "$OUTDIR/$PREFIX.key" -out "$OUTDIR/$PREFIX.csr" -subj "/C=$C/ST=$ST/L=$L/O=$O/OU=$OU/CN=$CN"
openssl x509 -req -in "$OUTDIR/$PREFIX.csr" -CA "$DEVLOCAL_CA_PEM" -CAkey "$DEVLOCAL_CA_KEY" -CAcreateserial -out "$OUTDIR/$PREFIX.cer" -days 3652 -sha256
rm "$OUTDIR/$PREFIX.csr"
echo -n "SHA256 digest: "
openssl x509 -outform der -in "$OUTDIR/$PREFIX.cer" | openssl dgst -sha256

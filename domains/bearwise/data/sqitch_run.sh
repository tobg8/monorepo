#!/bin/bash

# Check if the environment argument is provided
if [ -z "$2" ]; then
  echo "Error: No environment specified. Usage: $0 {deploy|revert|verify} {local|qa3|preprod|prod}"
  exit 1
fi

# Set the appropriate environment variable based on the argument
case "$2" in
  local)
    PG_URI=$LOCAL_PG_URI
    ;;
  qa3)
    PG_URI=$QA3_PG_URI
    ;;
  preprod)
    PG_URI=$PREPROD_PG_URI
    ;;
  prod)
    PG_URI=$PROD_PG_URI
    ;;
  *)
    echo "Error: Unknown environment '$2'. Usage: $0 {deploy|revert|verify} {local|qa3|preprod|prod}"
    exit 1
    ;;
esac

# Ensure the PG_URI is set
if [ -z "$PG_URI" ]; then
  echo "Error: PG_URI for environment '$2' is not set."
  exit 1
fi

# Check the first argument (command to run)
COMMAND=$1

# Shift the arguments to pass the remaining to Sqitch
shift 2

case "$COMMAND" in
  deploy)
    sqitch deploy $PG_URI "$@"
    ;;
  revert)
    sqitch revert $PG_URI "$@"
    ;;
  verify)
    sqitch verify $PG_URI "$@"
    ;;
  *)
    echo "Usage: $0 {deploy|revert|verify} {local|qa3|preprod|prod} [additional sqitch options]"
    exit 1
    ;;
esac
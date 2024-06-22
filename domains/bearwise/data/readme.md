# Sqitch

    apt-get update
    apt-get install sqitch libdbd-pg-perl

## Permissions
Ensure scripts are executable:

    chmod +x ./setup_local_pg_docker.sh
    chmod +x ./sqitch_run.sh

## Setting up Docker
Make sure Docker is installed locally. Confirm by running docker ps -a.
Then execute ./setup_local_pg_docker.sh to set up the local PostgreSQL database.

## Usage Instructions
Before using Sqitch, configure database URIs in .zshrc or .bashrc; contact an administrator for these details.

Once URIs are set, execute Sqitch commands via the script:

    Commands: {deploy|revert|verify}
    Environments: {local|qa3|preprod|prod}

## Examples
    Deploy locally: ./sqitch_run.sh deploy local
    Revert locally: ./sqitch_run.sh revert local
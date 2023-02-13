#!/bin/bash
echo '=================================='
echo 'START create-tablespace.sh'
echo '=================================='

mkdir -p /var/lib/postgresql/tablespaces/decision_helper_i
mkdir -p /var/lib/postgresql/tablespaces/decision_helper_d
set -e
# todo there is one password try change to environment
psql -v ON_ERROR_STOP=1 --username "decision_helper" --dbname "decision_helper" <<-EOSQL
  CREATE TABLESPACE decision_helper_d OWNER decision_helper LOCATION '/var/lib/postgresql/tablespaces/decision_helper_d';
  CREATE TABLESPACE decision_helper_i OWNER decision_helper LOCATION '/var/lib/postgresql/tablespaces/decision_helper_i';
  CREATE SCHEMA decision_helper;
  ALTER DATABASE decision_helper SET search_path to decision_helper;
EOSQL

echo '=================================='
echo 'FINISH create-tablespace.sh'
echo '=================================='
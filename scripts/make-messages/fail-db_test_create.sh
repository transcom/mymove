#! /usr/bin/env bash

echo -e "\e[31mSomething went wrong!\e[0m \e[33mRead more below:\e[0m"

cat <<DBTESTFAIL_MSG
  Oh no! There was something wrong with creating the test database.
  Try running the following Make commands to fix the issue.
  *********************************************
  *                                           *
  * >_ make db_test_reset db_test_migrate     *
  *                                           *
  *********************************************
DBTESTFAIL_MSG

exit 0

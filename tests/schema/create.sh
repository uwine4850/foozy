#!/bin/bash

sudo docker exec -i mysql bash -c "mysql --defaults-extra-file=/schema/mysql.cnf < ./schema/foozy_test.sql"
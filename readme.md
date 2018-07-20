# Go SQL Executer - simple tool for bulk execution of sql DDL and DML scripts

this has so far only been used with mysql but support for others should be as simple as adding the driver.

Example config

database mytest {
  driver        = "mysql"
  connectionstr = "batchscriptuser:password@tcp(localhost:3306)/testdb"
}

a test file is required with a list of scripts to execute

## Running it

gosqlexecuter -dbname=mytest -ScriptsListFile=myscriptslist.txt

## Building it

go build 
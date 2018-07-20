package main

/**
	This file is part of Go SQL Executer.

	Go SQL Executer - simple tool for bulk execution of sql DDL and DML scripts

    Go SQL Executer is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Go SQL Executer is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with Go SQL Executer.  If not, see <http://www.gnu.org/licenses/>.

**/

import (
	"bufio"
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/rhysbryant/gosqlexecuter/scriptexecutor"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	cfg, err := loadConfig("database.hcl")
	if err != nil {
		log.Fatalf("error [%s] when loading config file", err)
	}

	var dbname, scriptsListFile string

	flag.StringVar(&dbname, "dbname", "default", "The database profile to use")
	flag.StringVar(&scriptsListFile, "scriptslistfile", "", "the path to the scripts list")

	flag.Parse()

	if scriptsListFile == "" {
		flag.Usage()
	}

	var dbcfg database
	var ok bool

	if dbcfg, ok = cfg.Database[dbname]; !ok {
		log.Fatalln("no " + dbname + " database def found in config")
	}

	db, err := sql.Open(dbcfg.Driver, dbcfg.ConnectionStr)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(scriptsListFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scriptExecutor := script.NewScriptExecutor(db)
	C := consoleProgressOutput{}
	scriptExecutor.SetProgressExecutionHandler(C)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileToExecute := scanner.Text()
		if err := scriptExecutor.ExecuteScript(fileToExecute); err != nil {
			log.Fatalf("error [%s] when executing file [%s]\n", err, fileToExecute)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

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
	"database/sql"
	"fmt"
	"os"
	"time"
)

type consoleProgressOutput struct {
}

func (c consoleProgressOutput) BeginStatementExecution(id, statement string) {
	fmt.Printf("%s executing [%s][%s]", time.Now().String(), id, statement)
}

func (c consoleProgressOutput) StatementExecutionSucceeded(duration time.Duration, statement string, result sql.Result) {

	rowsCount, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nunable to get RowsAffected error %s\n", err)
		return
	}

	fmt.Printf(": %d rows affected, %4.3fs\n", rowsCount, duration.Seconds())
}

func (c consoleProgressOutput) StatementExecutionFailed(statement string, err error) {
	fmt.Fprintf(os.Stderr, "error %s \n", err)
}

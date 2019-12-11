package script

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
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rhysbryant/gosqlexecuter/parser"
)

type StatementExecution interface {
	StatementExecutionFailed(statement string, err error)
	StatementExecutionSucceeded(totalTime time.Duration, statement string, result sql.Result)
	BeginStatementExecution(id, statement string)
}

type DatabaseQueryExecuter interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Executor struct {
	progressInfo StatementExecution
	dbCon        DatabaseQueryExecuter
}

func NewScriptExecutor(dbCon DatabaseQueryExecuter) *Executor {
	s := Executor{}
	s.dbCon = dbCon

	return &s
}

func (se *Executor) SetProgressExecutionHandler(pei StatementExecution) {
	se.progressInfo = pei
}

func (se *Executor) ExecuteScript(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error [%s] opening script [%s]", err, filename)
	}

	defer f.Close()

	s := parser.NewSQLStatementParser(f)

	for {
		line, err := s.NextStatement()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if line == "" {
			continue
		}
		if se.progressInfo != nil {
			id := filepath.FromSlash(filename)
			se.progressInfo.BeginStatementExecution(id, line)
		}
		start := time.Now()
		result, err := se.dbCon.Exec(line)
		if err != nil {
			if se.progressInfo != nil {
				se.progressInfo.StatementExecutionFailed(line, fmt.Errorf("[%s]", err))
			}

			return fmt.Errorf("script execution failed")
		}

		if se.progressInfo != nil {
			duration := time.Now().Sub(start)
			se.progressInfo.StatementExecutionSucceeded(duration, line, result)
		}

	}

	return nil
}

package parser

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
	"bytes"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserBasicStatement(t *testing.T) {
	stmt := "select * from;"
	b := bytes.NewBufferString(stmt)
	s := NewSQLStatementParser(b)
	r, err := s.NextStatement()
	assert.NoError(t, err, "NextStatement returned an error")
	assert.Equal(t, "select * from", r)

}

func TestParserBasicStatementLeadingLines(t *testing.T) {
	stmt := "\r\n \r\n select * from;"
	b := bytes.NewBufferString(stmt)
	s := NewSQLStatementParser(b)
	r, err := s.NextStatement()
	assert.NoError(t, err, "NextStatement returned an error")
	assert.Equal(t, "select * from", r)

}

func TestParserBasicStatementLeadingLinesAndComment(t *testing.T) {
	stmt := "\r\n \r\n /* test */ select * from;"
	b := bytes.NewBufferString(stmt)
	s := NewSQLStatementParser(b)
	r, err := s.NextStatement()
	assert.NoError(t, err, "NextStatement returned an error")
	assert.Equal(t, " select * from", r)

}

func TestParserBasicStatementCommentedout(t *testing.T) {
	stmt := "-- CREATE INDEX P_01104_ORG_ID_INDEX ON P_01104(ORG_ID);\r\n"
	b := bytes.NewBufferString(stmt)
	s := NewSQLStatementParser(b)
	r, err := s.NextStatement()
	assert.Error(t, err, "EOF not returned")
	assert.Equal(t, "", r)

}

func TestParserBasicStatementCommentedoutNoLeadingNewLine(t *testing.T) {
	stmt := "\r\n-- CREATE INDEX P_01104_ORG_ID_INDEX ON P_01104(ORG_ID);\r\n;\r\n"
	b := bytes.NewBufferString(stmt)
	s := NewSQLStatementParser(b)
	r, err := s.NextStatement()
	assert.Error(t, err, "EOF not returned")
	assert.Equal(t, "", r)

}

func TestParserBasicStatementEmptyLine(t *testing.T) {
	stmt := "\r\n /* test */\r\n "
	b := bytes.NewBufferString(stmt)
	s := NewSQLStatementParser(b)
	r, err := s.NextStatement()
	assert.Error(t, err, "EOF not returned")
	assert.Equal(t, "", r)

}

func TestParsermultiLineStatement(t *testing.T) {
	stmt := `
	update sometable
	set n=9
	where m=7;
	-- test
	`
	b := bytes.NewBufferString(stmt)
	s := NewSQLStatementParser(b)
	r, err := s.NextStatement()
	assert.NoError(t, err, "NextStatement returned an error")
	assert.Equal(t, "update sometable\n\tset n=9\n\twhere m=7", r)

}

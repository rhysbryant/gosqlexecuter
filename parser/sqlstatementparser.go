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
	"errors"
	"io"
)

var (
	//ErrSyntex returned when there is tailing text
	ErrSyntex = errors.New("Syntax Error: unterminated statement, missing ;?")
)

//SQLStatementParser struct
type SQLStatementParser struct {
	src    io.Reader
	buffer []byte
	offset int
	size   int
}

const scriptBufferSize = 1024 * 1024

//NewSQLStatementParser creates a new statement parser (block of text ending ; skipping single line (--) and multi line (/* */) comments )
func NewSQLStatementParser(src io.Reader) *SQLStatementParser {
	s := SQLStatementParser{}
	s.buffer = make([]byte, scriptBufferSize)
	s.src = src
	s.offset = 0
	return &s
}

func (s *SQLStatementParser) consumeCommentType1(data []byte) int {
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			return i + 2
		} else if data[i] == '\r' || data[i] == '\n' {
			return i + 1
		}
	}
	return 1
}

func (s *SQLStatementParser) consumeCommentType2(data []byte) int {
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '*' && data[i+1] == '/' {
			return i + 3
		}
	}
	return 0
}

func (s *SQLStatementParser) consumeLeadingWhitespace(data []byte) int {
	for i := 0; i < len(data)-1; i++ {
		if data[i] > ' ' {
			return i + 1
		}
	}
	return 1
}

func (s *SQLStatementParser) loadBuffer() error {
	if s.offset >= s.size {
		var err error
		s.size, err = s.src.Read(s.buffer)
		return err
	}
	return nil
}

//NextStatement returns the next block ending in ;
func (s *SQLStatementParser) NextStatement() (string, error) {
	var lastChar byte
	var i = s.offset
	var err error

	s.loadBuffer()

	var buf bytes.Buffer
	for {

		for i < s.size {
			switch s.buffer[i] {
			case '\r':
				fallthrough
			case '\n':
				if buf.Len() == 0 { //strip leeding whitespace
					i += s.consumeLeadingWhitespace(s.buffer[i:s.size])

					lastChar = 0
				} else {
					i++
				}
				break
			case '-':
				if lastChar == '-' { //strip -- some text\r\n
					i += s.consumeCommentType1(s.buffer[i-1 : s.size])
					//i++
					lastChar = 0
					break
				} else {
					i++
				}
			case '*':
				if lastChar == '/' { //strip /* some text */
					i += 2
					i += s.consumeCommentType2(s.buffer[i:s.size])
					lastChar = 0
					break
				}
				i++
				s.offset = i
			case ';':
				s.offset = i + 1
				if lastChar != 0 {
					buf.WriteByte(lastChar)
				}

				return buf.String(), nil
			default:
				i++

			}
			if lastChar != 0 {
				buf.WriteByte(lastChar)
			}
			if i > 0 {
				lastChar = s.buffer[i-1]
			}

		}
		s.offset = i
		if err = s.loadBuffer(); err != nil {
			if err == io.EOF && buf.Len() > 0 {
				return "", ErrSyntex
			}
			return "", err
		}
		s.offset = 0
		i = 0
	}

}

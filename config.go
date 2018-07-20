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
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

type database struct {
	Driver        string `hcl:"driver"`
	ConnectionStr string `hcl:"connectionstr"`
}

type config struct {
	Database        map[string]database `hcl:"database"`
	ScriptsListFile string              `hcl:"scriptslistfile"`
}

func loadConfig(fileName string) (*config, error) {
	config := config{}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	e := hcl.Unmarshal(file, &config)
	if e != nil {
		return nil, e
	}
	return &config, nil
}

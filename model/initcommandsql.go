package model

import (
	"io/ioutil"
	"os"
)

//InitSQL creates tables in sqlite database
func InitSQL() error {
	file, err := os.Open("./sqlite/command.sql")
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if _, err = Db.Exec(string(data)); err != nil {
		return err
	}

	return nil
}

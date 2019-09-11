package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

// CreateDB - sql to create db
const CreateDB = "CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4;\n"

// Source -
type Source struct {
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	NoData       bool     `toml:"no-data"`
	IgnoreTables []string `toml:"ignore-table"`
}

// Destination -
type Destination struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// Config -
type Config struct {
	Source      Source
	Destination Destination
}

func main() {
	doc, err := ioutil.ReadFile("/app/config.toml")
	if err != nil {
		log.Fatal("Fail to load config")
	}

	config := Config{}
	toml.Unmarshal(doc, &config)
	fmt.Println(config)

	if config.Source.Database == "" {
		log.Fatal("database should not be empty")
	}

	if config.Destination.Database == "" {
		config.Destination.Database = config.Source.Database
	}

	file := path.Join("/data", config.Source.Database+".sql")
	if _, err := os.Create(file); err != nil {
		log.Fatal(err)
	}

	var args []string
	args = append(args, "-h"+config.Source.Host)
	args = append(args, "-P"+config.Source.Port)
	args = append(args, "-u"+config.Source.User)
	args = append(args, "-p"+config.Source.Password)
	args = append(args, "-r"+file)

	if config.Source.NoData {
		args = append(args, "-d")
	}

	if config.Source.IgnoreTables != nil {
		for _, t := range config.Source.IgnoreTables {
			args = append(args, "--ignore-table="+config.Source.Database+"."+t)
		}
	}

	args = append(args, config.Source.Database)

	runCmd(exec.Command("mysqldump", args...))

	// need to import the sql into target database
	if config.Destination.Host != "" &&
		config.Destination.User != "" {

		sql := fmt.Sprintf(CreateDB, config.Destination.Database)

		var args []string
		args = append(args, "-h"+config.Destination.Host)
		args = append(args, "-P"+config.Destination.Port)
		args = append(args, "-u"+config.Destination.User)
		args = append(args, "-p"+config.Destination.Password)

		var args2 = make([]string, 4)
		var args3 = make([]string, 4)
		copy(args2, args)
		copy(args3, args)

		args2 = append(args2, "-e "+sql+"")
		runCmd(exec.Command("mysql", args2...))

		args3 = append(args3, config.Destination.Database)
		args3 = append(args3, "<")
		args3 = append(args3, file)

		cmd := "mysql " + strings.Join(args3, " ")
		runCmd(exec.Command("/bin/sh", "-c", cmd))
	}
}

func runCmd(cmd *exec.Cmd) {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	fmt.Println(cmd.String())
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		log.Fatal(err)
	}
}

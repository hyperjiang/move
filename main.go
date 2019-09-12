package main

import (
	"bytes"
	"flag"
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
const CreateDB = "CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4;"

// Option represents options
type Option struct {
	NoData       bool     `toml:"no-data"`
	IgnoreTables []string `toml:"ignore-table"`
}

// DSN represents data source name
type DSN struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// Config represents the config in toml
type Config struct {
	Source      DSN
	Destination DSN
	Option      Option
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "/app/config.toml", "Config file path")
	flag.Parse()

	doc, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Fail to load config")
	}

	config := Config{}
	toml.Unmarshal(doc, &config)
	fmt.Println(config)

	if config.Source.Database == "" {
		log.Fatal("database should not be empty")
	}

	file := path.Join("/data", config.Source.Database+".sql")
	if _, err := os.Create(file); err != nil {
		log.Fatal(err)
	}

	var args = config.Source.buildArgs()
	args = append(args, "-r"+file)

	if config.Option.NoData {
		args = append(args, "-d")
	}

	if config.Option.IgnoreTables != nil {
		for _, t := range config.Option.IgnoreTables {
			args = append(args, "--ignore-table="+config.Source.Database+"."+t)
		}
	}

	args = append(args, config.Source.Database)

	runCmd(buildCmd("mysqldump", args))

	// need to import the sql into target database
	if config.Destination.Host != "" &&
		config.Destination.User != "" {

		if config.Destination.Database == "" {
			config.Destination.Database = config.Source.Database
		}

		var args = config.Destination.buildArgs()
		var args2 = make([]string, 4)
		var args3 = make([]string, 4)
		copy(args2, args)
		copy(args3, args)

		// create database if not exists
		sql := fmt.Sprintf(CreateDB, config.Destination.Database)
		args2 = append(args2, "-e '"+sql+"'")
		runCmd(buildCmd("mysql", args2))

		// import sql into database
		args3 = append(args3, config.Destination.Database)
		args3 = append(args3, "<")
		args3 = append(args3, file)
		runCmd(buildCmd("mysql", args3))
	}
}

func (dsn DSN) buildArgs() []string {
	var args []string
	args = append(args, "-h"+dsn.Host)
	args = append(args, "-P"+dsn.Port)
	args = append(args, "-u"+dsn.User)
	args = append(args, "-p"+dsn.Password)
	return args
}

func buildCmd(name string, args []string) *exec.Cmd {
	cmd := name + " " + strings.Join(args, " ")
	return exec.Command("/bin/sh", "-c", cmd)
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

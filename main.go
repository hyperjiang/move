package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

const sqlCreateDB = "CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4;"

// Option represents options
type Option struct {
	NoData         bool     `toml:"no-data"`
	NoCreateInfo   bool     `toml:"no-create-info"`
	SkipLockTables bool     `toml:"skip-lock-tables"`
	IgnoreTables   []string `toml:"ignore-tables"`
	Tables         []string `toml:"tables"`
}

// DSN represents data source name
type DSN struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// Rule represents a rule
type Rule struct {
	Name        string
	Source      DSN
	Destination DSN
	Option      Option
}

// Config represents the config in toml
type Config struct {
	Rules []Rule `toml:"rule"`
}

func main() {
	var configFile string
	var ruleName string
	flag.StringVar(&configFile, "c", "/app/config.toml", "Config file path")
	flag.StringVar(&ruleName, "r", "", "Specify which rule to run, if not specify then all the rules in config.toml will be run")
	flag.Parse()

	config := Config{}
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	// check if all the rules are valid before real run
	for _, rule := range config.Rules {
		rule.check()
	}

	// run rules
	for _, rule := range config.Rules {
		if ruleName != "" && ruleName != rule.Name {
			continue
		}
		rule.handle()
	}
}

func (rule Rule) check() {
	if rule.Name == "" {
		log.Fatal("rule name should not be empty")
	}

	if rule.Source.Host == "" {
		log.Fatal("source host should not be empty in rule " + rule.Name)
	}

	if rule.Source.User == "" {
		log.Fatal("source user should not be empty in rule " + rule.Name)
	}

	if rule.Source.Database == "" {
		log.Fatal("source database should not be empty in rule " + rule.Name)
	}
}

func (rule Rule) handle() {
	file := path.Join("/data", rule.Name+".sql")
	if _, err := os.Create(file); err != nil {
		log.Fatal(err)
	}

	var args = rule.Source.buildArgs()
	args = append(args, "-r"+file)

	// do not dump table contents, equals --no-data
	if rule.Option.NoData {
		args = append(args, "-d")
	}

	// do not write CREATE TABLE statements which re-create each dumped table, equals --no-create-info
	if rule.Option.NoCreateInfo {
		args = append(args, "-t")
	}

	// do not lock tables
	if rule.Option.SkipLockTables {
		args = append(args, "--skip-lock-tables")
	}

	if rule.Option.IgnoreTables != nil {
		for _, t := range rule.Option.IgnoreTables {
			args = append(args, "--ignore-table="+rule.Source.Database+"."+t)
		}
	}

	args = append(args, rule.Source.Database)

	if rule.Option.Tables != nil {
		for _, t := range rule.Option.Tables {
			args = append(args, t)
		}
	}

	runCmd(buildCmd("mysqldump", args))

	// need to import the sql into target database
	if rule.Destination.Host != "" &&
		rule.Destination.User != "" {

		if rule.Destination.Database == "" {
			rule.Destination.Database = rule.Source.Database
		}

		var args = rule.Destination.buildArgs()
		var args2 = make([]string, 4)
		var args3 = make([]string, 4)
		copy(args2, args)
		copy(args3, args)

		// create database if not exists
		sql := fmt.Sprintf(sqlCreateDB, rule.Destination.Database)
		args2 = append(args2, "-e '"+sql+"'")
		runCmd(buildCmd("mysql", args2))

		// import sql into database
		args3 = append(args3, rule.Destination.Database)
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

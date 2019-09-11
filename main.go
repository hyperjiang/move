package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/BurntSushi/toml"
)

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

	file := path.Join("/data", config.Source.Database+".sql")
	_, err = os.Create(file)
	if err != nil {
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

	cmd := exec.Command("mysqldump", args...)

	fmt.Println(cmd.String())

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		log.Fatal(err)
	}
}

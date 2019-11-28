package main

import (
	"github.com/mholt/archiver/v3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	config := new(config)
	zip := archiver.Zip{
		OverwriteExisting: true,
	}

	// open file
	file, err := os.Open(".serve.yml")

	if err != nil {
		log.Fatal(err)
	}

	// close file
	defer func() {
		err = file.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	// read file
	b, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	// to yaml
	err = yaml.Unmarshal(b, &config)

	if err != nil {
		log.Fatal(err)
	}

	err = zip.Archive(config.Files, "deploy.zip")

	if err != nil {
		log.Fatal(err)
	}
}

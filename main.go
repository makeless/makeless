package main

import (
	"bytes"
	"github.com/mholt/archiver/v3"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
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

	// to docker-compose.yml
	y, err := yaml.Marshal(map[string]interface{}{
		"version": "3",
		"services": map[string]interface{}{
			config.Service: config.Compose,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	// write docker-compose.yml file
	err = ioutil.WriteFile("docker-compose.yml", y, 0644)

	if err != nil {
		log.Fatal(err)
	}

	// append docker-compose.yml file
	config.Files = append(config.Files, "docker-compose.yml")

	// zip
	err = zip.Archive(config.Files, "build.zip")

	if err != nil {
		log.Fatal(err)
	}

	err = postFile("build.zip", "http://localhost:8080/deploy")
}

func postFile(filename string, targetUrl string) error {
	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, err := bodyWriter.CreateFormFile("file", filename)

	if err != nil {
		return err
	}

	// open file handle
	file, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = io.Copy(fileWriter, file)

	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()

	// close bodyWriter
	err = bodyWriter.Close()

	if err != nil {
		return err
	}

	// post file
	resp, err := http.Post(targetUrl, contentType, bodyBuffer)

	if err != nil {
		return err
	}

	// close response body
	defer func() {
		err := resp.Body.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	// print response body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	log.Println(string(body))
	return nil
}
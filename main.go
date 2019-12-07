package main

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

	// build compose
	compose := map[string]interface{}{
		"version": "3",
		"services": map[string]interface{}{
			config.Name: config.Service,
		},
	}

	// assign shared values
	for key, value := range config.Shared {
		compose[key] = value
	}

	// to docker-compose.yml
	y, err := yaml.Marshal(compose)

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
	err = zip.Archive(config.Files, "deploy.zip")

	if err != nil {
		log.Fatal(err)
	}

	err = post(config, "deploy.zip", fmt.Sprintf("http://%s/deploy", config.Host))

	if err != nil {
		log.Fatal(err)
	}
}

func getSignedToken(name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
	})

	return token.SignedString([]byte(os.Getenv("TOKEN")))
}

func post(config *config, filename string, targetUrl string) error {
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

	// name field
	fieldWriter, err := bodyWriter.CreateFormField("name")

	if err != nil {
		log.Fatal(err)
	}

	_, err = fieldWriter.Write([]byte(config.Name))

	if err != nil {
		log.Fatal(err)
	}

	contentType := bodyWriter.FormDataContentType()

	// close bodyWriter
	err = bodyWriter.Close()

	if err != nil {
		return err
	}

	// get signed token
	signedToken, err := getSignedToken(config.Name)

	if err != nil {
		return err
	}

	// request
	req, err := http.NewRequest("POST", targetUrl, bodyBuffer)

	if err != nil {
		return err
	}

	// add headers
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", signedToken))

	// client
	client := new(http.Client)

	// post
	resp, err := client.Do(req)

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

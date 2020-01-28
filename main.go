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
	"strings"
)

func main() {
	config := new(config)
	zip := archiver.Zip{
		OverwriteExisting: true,
	}

	// open file
	file, err := os.Open(".makeless.yml")

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
	configBytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	// to yaml
	err = yaml.Unmarshal(configBytes, &config)

	if err != nil {
		log.Fatal(err)
	}

	// build compose
	var compose = make(map[string]interface{})

	// compose -> version
	compose["version"] = "3"

	// compose -> services
	if config.Service != nil {
		compose["services"] = map[string]interface{}{
			config.Name: config.Service,
		}
	}

	// compose -> shared
	for key, value := range config.Shared {
		if key == "services" {
			continue
		}

		compose[key] = value
	}

	// to docker-compose.yml
	y, err := yaml.Marshal(compose)

	if err != nil {
		log.Fatal(err)
	}

	// replace placeholders
	yStr := string(y)

	// --> %build_dir%
	yStr = strings.ReplaceAll(
		yStr,
		"%build_dir%",
		fmt.Sprintf("/home/makeless/containers/%s/latest", config.Name),
	)

	// write docker-compose.yml file
	err = ioutil.WriteFile("docker-compose.yml", []byte(yStr), 0644)

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

	err = post(config, configBytes, "deploy.zip", fmt.Sprintf("http://%s/deploy", config.Host))

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

func post(config *config, configBytes []byte, filename string, targetUrl string) error {
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

	// config field
	configFieldWriter, err := bodyWriter.CreateFormField("config")

	if err != nil {
		return err
	}

	_, err = configFieldWriter.Write(configBytes)

	if err != nil {
		return err
	}

	// content type
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

	log.Printf("%s", string(body))
	return nil
}

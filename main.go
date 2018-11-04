package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

type Gojanus struct {
	AdminURL    string
	AdminSecret string
}

func (g *Gojanus) GenerateToken() (string, error) {
	log.Info("GenerateToken")
	newUuid := uuid.NewV4()
	transactionId := fmt.Sprintf("%s", newUuid)
	newUuid = uuid.NewV4()
	hash, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s", newUuid)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(hash)

	var jsonStr = []byte(`{
        "janus" : "add_token",
        "token": "` + token + `",
        "transaction": "` + string(transactionId) + `",
        "admin_secret": "` + g.AdminSecret + `"
	  }`)
	req, err := http.NewRequest("POST", g.AdminURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(body))
	}

	return token, nil
}

func (g *Gojanus) RemoveToken(token string) error {
	log.Info("RemoveToken ", token)
	newUuid := uuid.NewV4()
	transactionId := fmt.Sprintf("%s", newUuid)

	var jsonStr = []byte(`{
        "janus" : "remove_token",
        "token": "` + token + `",
        "transaction": "` + string(transactionId) + `",
        "admin_secret": "` + g.AdminSecret + `"
	  }`)
	req, err := http.NewRequest("POST", g.AdminURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}

func (g *Gojanus) ListToken() ([]string, error) {
	log.Info("ListToken ")
	tokens := []string{}
	newUuid := uuid.NewV4()
	transactionId := fmt.Sprintf("%s", newUuid)

	var jsonStr = []byte(`{
        "janus" : "list_tokens",
        "transaction": "` + string(transactionId) + `",
        "admin_secret": "` + g.AdminSecret + `"
	  }`)
	req, err := http.NewRequest("POST", g.AdminURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return tokens, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return tokens, errors.New(string(body))
	}
	var objmap map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &objmap); err != nil {
		return tokens, err
	}
	length := len(objmap["data"].(map[string]interface{})["tokens"].([]interface{}))
	index := 1
	for index <= length {
		tokens = append(tokens, objmap["data"].(map[string]interface{})["tokens"].([]interface{})[index-1].(map[string]interface{})["token"].(string))
		index++
	}
	return tokens, nil
}

func main() {
	gojanus := &Gojanus{
		AdminURL:    "http://herpiko-devbox.tarsius.id:7088/admin",
		AdminSecret: "janusoverlord",
	}
	token, err := gojanus.GenerateToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	tokens, err := gojanus.ListToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(tokens))
	err = gojanus.RemoveToken(tokens[0])
	if err != nil {
		panic(err)
	}
	tokens, err = gojanus.ListToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(tokens))
}

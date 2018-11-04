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

func (g *Gojanus) ListTokens() ([]string, error) {
	log.Info("ListTokens")
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
    if objmap["janus"].(string) != "success" {
		return tokens, errors.New(objmap["error"].(map[string]interface{})["reason"].(string))
    }
	length := len(objmap["data"].(map[string]interface{})["tokens"].([]interface{}))
	index := 1
	for index <= length {
		tokens = append(tokens, objmap["data"].(map[string]interface{})["tokens"].([]interface{})[index-1].(map[string]interface{})["token"].(string))
		index++
	}
	return tokens, nil
}

func (g *Gojanus) ListSessions() ([]string, error) {
	log.Info("ListSessions")
	sessions := []string{}
	newUuid := uuid.NewV4()
	transactionId := fmt.Sprintf("%s", newUuid)

	var jsonStr = []byte(`{
        "janus" : "list_sessions",
        "transaction": "` + string(transactionId) + `",
        "admin_secret": "` + g.AdminSecret + `"
	  }`)
	req, err := http.NewRequest("POST", g.AdminURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return sessions, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return sessions, errors.New(string(body))
	}
	var objmap map[string]interface{}
	if err := json.Unmarshal([]byte(string(body)), &objmap); err != nil {
		return sessions, err
	}
    if objmap["janus"].(string) != "success" {
		return sessions, errors.New(objmap["error"].(map[string]interface{})["reason"].(string))
    }
	length := len(objmap["sessions"].([]interface{}))
	index := 1
	for index <= length {
		sessions = append(sessions, objmap["sessions"].([]interface{})[index-1].(string))
		index++
	}
	return sessions, nil
}

func main() {
	gojanus := &Gojanus{
		AdminURL:    "http://localhost:7088/admin",
		AdminSecret: "janusoverlord",
	}
	token, err := gojanus.GenerateToken()
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	tokens, err := gojanus.ListTokens()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(tokens))
	err = gojanus.RemoveToken(tokens[1])
	if err != nil {
		panic(err)
	}
    sessions, err := gojanus.ListSessions()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(sessions))
}

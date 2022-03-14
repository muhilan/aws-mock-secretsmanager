package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var wg sync.WaitGroup

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type GetSecretsResponse struct {
	Name          string
	VersionId     string
	SecretString  string
	VersionStages []string
	CreatedDate   int64
	ARN           string
}

type GetSecretsErrorResponse struct {
	Type    string `json:"__type"`
	Message string
}

type GetSecretsRequest struct {
	SecretId string
}

var dataMap MyMap

type MyMap struct {
	m    sync.Mutex
	data map[string]string
}

func main() {
	dataMap = MyMap{data: make(map[string]string)}

	baseDir, err := baseDir()

	if err != nil {
		return
	}
	createMap(baseDir)

	http.HandleFunc("/", postHandler)
	http.ListenAndServe(":8080", nil)
}
func postHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var r GetSecretsRequest
		err := decoder.Decode(&r)
		if err != nil {
			panic(err)
		}

		response := GetSecretsResponse{}
		value := dataMap.data[r.SecretId]

		if value == "" {
			json.NewEncoder(w).Encode(GetSecretsErrorResponse{Type: "ResourceNotFoundException", Message: "Secrets Manager can’t find the specified secret."})
			return
		}
		response.SecretString = value
		response.ARN = fmt.Sprintf("arn:aws:secretsmanager:us-west-2:1234567789:secret:%s-%s", r.SecretId, RandStringRunes(6))
		response.VersionId = RandStringRunes(6)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)

	} else {
		io.WriteString(w, "Method not supported")
	}
}

func load(file string, name string) {
	defer wg.Done()
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
	}
	var key = name[0 : len(name)-len(filepath.Ext(name))]
	var finalKey = strings.Replace(key, ".", "/", -1)
	dataMap.m.Lock()
	dataMap.data[finalKey] = string(content)
	dataMap.m.Unlock()

}

func createMap(baseDir string) {
	start := time.Now()
	// Traverse filepath and update data map
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error %v\n", err)
			return err
		}
		// Max AWS secret size is 10 KB
		if !info.IsDir() && info.Size() <= 10000 {
			go load(path, info.Name())
			wg.Add(1)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Added %d secrets in %s", len(dataMap.data), elapsed)
}

func baseDir() (string, error) {
	baseDir := os.Getenv("DATA_DIR")

	if baseDir == "" {
		baseDir = "/data"
	}

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		log.Printf("Path \"%s\" does not exist", baseDir)
		return "", err
	}
	return baseDir, nil
}

package main

import (
	"io"
	"net/http"
	"encoding/json"
	"log"
	"os"
	// "bufio"
	// "strings"
	"path/filepath"
	// "net/http/httputil"
	"io/ioutil"
	"sync"
	"math/rand"
	"fmt"

)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

type GetSecretsResponse struct {
	 Name string
	 VersionId string
	 SecretString string
	 VersionStages []string
	 CreatedDate int64
	 ARN string
}

type GetSecretsRequest struct {
	 SecretId string
}

var dataMap MyMap

type MyMap struct {
	lock sync.Mutex
	data map[string]string
}


func main() {
	dataMap = MyMap{data: make(map[string]string)}

	baseDir := os.Getenv("DATA_DIR")+ "/test/"
	log.Println("baseDir is "+ baseDir)


	// Traverse filepath and update data map
    err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
    	if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		log.Println(path)
	    log.Println(info.Name())
    	if info.IsDir() {
			log.Printf("skipping a dir without errors: %+v \n", info.Name())
			// return filepath.SkipDir
		} else {
			go load(path, info.Name())
		}
	    
	    
	    return nil
    })


    if err != nil {
        panic(err)
    }
	
	log.Println(dataMap)

	http.HandleFunc("/", helloWorldHandler)
	http.ListenAndServe(":8080", nil)
}
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	// dataMap := map[string]string{
 //        "Name": "public-cert",
 //        "VersionId" : "yu87678",
 //        "SecretString" : "879879",
 //        "VersionStages" : [
 //        	"AWSCURRENT"
 //        ],
 //        "CreatedDate" : 158787987.78,
 //        "ARN" : "arn:aws:secretsmanager:us-east-1:1234567789:secret:public-cert-89cXyz",
 //    }
	if r.Method == http.MethodPost {
      


	// requestDump, err := httputil.DumpRequest(r, true)
	// if err != nil {
	//   log.Println(err)
	// }
	// log.Println(string(requestDump))

    decoder := json.NewDecoder(r.Body)
    var r GetSecretsRequest
    err := decoder.Decode(&r)
    if err != nil {
        panic(err)
    }
    log.Println(r)
	response := GetSecretsResponse{}
	value := dataMap.data[r.SecretId]
	if value == "" {
		log.Println("key not found " + r.SecretId)
		io.WriteString(w, "Not found")
		return
	}
	response.SecretString = value;
	response.ARN = fmt.Sprintf("arn:aws:secretsmanager:us-east-1:1234567789:secret:%s-%s", r.SecretId, RandStringRunes(6))
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(response)

    } else {
		io.WriteString(w, "Method not supported")
	}
}

func load(file string, name string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
	}
	dataMap.lock.Lock()
	dataMap.data[name[0:len(name)-len(filepath.Ext(name))]] = string(content)
	dataMap.lock.Unlock()
	log.Println(dataMap)
	
}
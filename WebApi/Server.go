package WebApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"unicorn.dev.web-scrap/Regestery/GovRu"
)

const INVALID_REQUEST_METHOD_POST string = "Only POST request allowed"
const INVALID_REQUEST_METHOD_GET string = "Only GET request allowed"
const UNABLE_TO_GET_REQUEST_BODY string = "Unable to get request body"
const INVALID_REQUEST_DATA string = "Invalid request data"

type damiaConf struct {
	Key    string `json:"key"`
	Active bool
}

var damiaConfig damiaConf

type apiKeysFile struct {
	Damia damiaConf `json:"damia"`
}

type ServerConfiguration struct {
	UseTls      bool
	TlsCrtFile  string
	TlsKeyFile  string
	ApiKeysFile string
}

func ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, INVALID_REQUEST_METHOD_GET, http.StatusMethodNotAllowed)
		return
	}

	_, err := fmt.Fprint(w, "pong")
	if err != nil {
		fmt.Errorf("Unable to send response to ", r.RemoteAddr)
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, INVALID_REQUEST_METHOD_POST, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UNABLE_TO_GET_REQUEST_BODY, http.StatusInternalServerError)
		return
	}

	var searchRequest SearchRequest
	err = json.Unmarshal(body, &searchRequest)
	if err != nil {
		http.Error(w, INVALID_REQUEST_DATA, http.StatusBadRequest)
		return
	}

	govRuTask := createTaskContext()
	govRuSearchQueru := GovRu.NewSearchQuery()
	govRuSearchQueru.Keywords = searchRequest.Keywords
	if damiaConfig.Active {
		go GovRu.Search(damiaConfig.Key, govRuSearchQueru, govRuTask)
	}

	response := ResponseStatus{
		Status: "OK",
	}
	response.IDs = append(response.IDs, govRuTask.Id)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Errorf("Unable to marshal response JSON")
	}

	http.Error(w, string(jsonResponse), http.StatusOK)
}

func handleTaskStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, INVALID_REQUEST_METHOD_GET, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UNABLE_TO_GET_REQUEST_BODY, http.StatusInternalServerError)
		return
	}

	var taskStatusRequest TaskStatusRequest
	err = json.Unmarshal(body, &taskStatusRequest)
	if err != nil {
		http.Error(w, INVALID_REQUEST_DATA, http.StatusBadRequest)
		return
	}

	task := getTaskContext(taskStatusRequest.Id)
	if task == nil {
		response := ResponseStatus{
			Status: "Not exists",
			IDs:    make([]uint64, 0),
		}

		jsonResponse, _ := json.Marshal(response)
		http.Error(w, string(jsonResponse), http.StatusNotFound)
		return
	}

	jsonResponse, _ := json.Marshal(*task)
	http.Error(w, string(jsonResponse), http.StatusOK)
}

func StartServer(configuration *ServerConfiguration) {

	var err error

	if configuration == nil {
		log.Fatal("Unable to run server")
	}

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/status/task", handleTaskStatus)

	if len(configuration.ApiKeysFile) > 0 {
		value, err := ioutil.ReadFile(configuration.ApiKeysFile)
		if err != nil {
			log.Fatal("Unable to read " + configuration.ApiKeysFile)
			return
		}

		f := apiKeysFile{}
		err = json.Unmarshal(value, &f)
		if err != nil {
			log.Fatal("Unable to parse " + configuration.ApiKeysFile)
			return
		}

		if len(f.Damia.Key) > 0 {
			damiaConfig = f.Damia
			damiaConfig.Active = true
		}
	}

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	initTaskContext()

	if configuration.UseTls {
		err = s.ListenAndServeTLS(configuration.TlsCrtFile, configuration.TlsKeyFile)
	} else {
		err = s.ListenAndServe()
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

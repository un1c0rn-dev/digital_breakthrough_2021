package WebApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"unicorn.dev.web-scrap/Regestery/GovRu"
	"unicorn.dev.web-scrap/Tasks"
)

const (
	InvalidRequestMethodPost string = "Only POST request allowed"
	InvalidRequestMethodGet  string = "Only GET request allowed"
	UnableToGetRequestBody   string = "Unable to get request body"
	InvalidRequestData       string = "Invalid request data"
)

type apiKeysFile struct {
	Damia GovRu.DamiaConf `json:"damia"`
}

type ServerConfiguration struct {
	UseTls      bool
	TlsCrtFile  string
	TlsKeyFile  string
	ApiKeysFile string
}

func ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, InvalidRequestMethodGet, http.StatusMethodNotAllowed)
		return
	}

	_, err := fmt.Fprint(w, "pong")
	if err != nil {
		fmt.Errorf("Unable to send response to ", r.RemoteAddr)
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, InvalidRequestMethodPost, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UnableToGetRequestBody, http.StatusInternalServerError)
		return
	}

	var searchRequest SearchRequest
	err = json.Unmarshal(body, &searchRequest)
	if err != nil {
		http.Error(w, InvalidRequestData, http.StatusBadRequest)
		return
	}

	govRuTask := createTaskContext()
	govRuSearchQueru := GovRu.NewSearchQuery()
	govRuSearchQueru.Keywords = searchRequest.Keywords
	go GovRu.Search(govRuSearchQueru, govRuTask)

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

	if r.Method != http.MethodPost {
		http.Error(w, InvalidRequestMethodPost, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UnableToGetRequestBody, http.StatusInternalServerError)
		return
	}

	var taskStatusRequest TaskStatusRequest
	err = json.Unmarshal(body, &taskStatusRequest)
	if err != nil {
		http.Error(w, InvalidRequestData, http.StatusBadRequest)
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

func handleDataCollect(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, InvalidRequestMethodGet, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UnableToGetRequestBody, http.StatusInternalServerError)
		return
	}

	var collectDataRequest CollectDataRequest
	err = json.Unmarshal(body, &collectDataRequest)
	if err != nil {
		http.Error(w, InvalidRequestData, http.StatusBadRequest)
		return
	}

	responseCollectData := ResponseCollectData{}
	for _, id := range collectDataRequest.Ids {
		task := getTaskContext(id)
		if task == nil {
			response := ResponseStatus{
				Status: "Not exists",
				IDs:    make([]uint64, 0),
			}

			jsonResponse, _ := json.Marshal(response)
			http.Error(w, string(jsonResponse), http.StatusNotFound)
			return
		}

		if task.Status == Tasks.TaskStatusError {
			response := ResponseStatus{
				Status: "Task failed",
				IDs:    make([]uint64, 0),
			}

			jsonResponse, _ := json.Marshal(response)
			http.Error(w, string(jsonResponse), http.StatusInternalServerError)
			return
		}

		if task.Status != Tasks.TaskStatusDone {
			response := ResponseStatus{
				Status: "Not ready",
				IDs:    make([]uint64, 0),
			}

			jsonResponse, _ := json.Marshal(response)
			http.Error(w, string(jsonResponse), http.StatusTooEarly)
			return
		}

		if task.Result == nil {
			response := ResponseStatus{
				Status: "Empty result",
				IDs:    make([]uint64, 0),
			}

			jsonResponse, _ := json.Marshal(response)
			http.Error(w, string(jsonResponse), http.StatusInternalServerError)
			return
		}

		responseCollectData.Data[strconv.FormatUint(id, 10)] = task.Result
	}

	jsonResponse, _ := json.Marshal(&responseCollectData)
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
	http.HandleFunc("/data/collect", handleDataCollect)

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

		GovRu.Configure(f.Damia)
	}

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	InitContext()

	if configuration.UseTls {
		err = s.ListenAndServeTLS(configuration.TlsCrtFile, configuration.TlsKeyFile)
	} else {
		err = s.ListenAndServe()
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

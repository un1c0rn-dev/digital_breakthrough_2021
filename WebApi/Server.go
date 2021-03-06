package WebApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	UseTls       bool
	TlsCrtFile   string
	TlsKeyFile   string
	ApiKeysFile  string
	Port         string
	TemplatesDir string
}

func setupCORS(w *http.ResponseWriter, r *http.Request) bool {
	fmt.Println("Setting up CORS for " + r.RemoteAddr)
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		http.Error(*w, "POST,GET", http.StatusOK)
		return false
	}

	return true
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

	fmt.Println(r.Method + " /search " + r.RemoteAddr)

	if !setupCORS(&w, r) {
		return
	}

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
	if searchRequest.MaxRequests > 0 {
		govRuSearchQueru.MaxRequests = searchRequest.MaxRequests
	}
	if searchRequest.MinPrice > 0 {
		govRuSearchQueru.MinPrice = searchRequest.MinPrice
	}
	if searchRequest.MaxPrice > 0 {
		govRuSearchQueru.MaxPrice = searchRequest.MaxPrice
	}
	if len(searchRequest.Etp) > 0 {
		govRuSearchQueru.Etp = searchRequest.Etp
	}
	if len(searchRequest.Placing) > 0 {
		govRuSearchQueru.Placing = searchRequest.Placing
	}
	if len(searchRequest.Region) > 0 {
		govRuSearchQueru.Region = searchRequest.Region
	}
	if len(searchRequest.Okpd) > 0 {
		govRuSearchQueru.Okpd = searchRequest.Okpd
	}
	if searchRequest.Status > 0 {
		govRuSearchQueru.Status = GovRu.SearchStatusCode(searchRequest.Status)
	}
	if searchRequest.Fz > 0 {
		govRuSearchQueru.Fz = GovRu.SearchFZ(searchRequest.Fz)
	}
	if searchRequest.ToDateYMD[0] != 0 && searchRequest.ToDateYMD[1] != 0 && searchRequest.ToDateYMD[2] != 0 {
		govRuSearchQueru.ToDateYMD = searchRequest.ToDateYMD
	}
	if searchRequest.FromDateYMD[0] != 0 && searchRequest.FromDateYMD[1] != 0 && searchRequest.FromDateYMD[2] != 0 {
		govRuSearchQueru.FromDateYMD = searchRequest.FromDateYMD
	}
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

func handleTasksStatus(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Method + " /status/tasks " + r.RemoteAddr)

	if !setupCORS(&w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, InvalidRequestMethodPost, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UnableToGetRequestBody, http.StatusInternalServerError)
		return
	}

	var tasksStatusRequest TaskStatusRequest
	err = json.Unmarshal(body, &tasksStatusRequest)
	if err != nil {
		http.Error(w, InvalidRequestData, http.StatusBadRequest)
		return
	}

	data := make([]Tasks.Task, 0)
	for _, id := range tasksStatusRequest.Id {
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

		data = append(data, *task)
	}

	jsonResponse, _ := json.Marshal(data)
	http.Error(w, string(jsonResponse), http.StatusOK)
}

func handleDataCollect(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Method + " /data/collect " + r.RemoteAddr)

	if !setupCORS(&w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, InvalidRequestMethodPost, http.StatusMethodNotAllowed)
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
	responseCollectData.Data = make(map[string][]Tasks.TaskResult)

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

		if task.Status != Tasks.TaskStatusDone && task.Status != Tasks.TaskStatusError {
			response := ResponseStatus{
				Status: "Not ready",
				IDs:    make([]uint64, 0),
			}

			jsonResponse, _ := json.Marshal(response)
			http.Error(w, string(jsonResponse), http.StatusTooEarly)
			return
		}

		if task.Result == nil {
			removeTaskContext(task)
			continue
		}

		responseCollectData.Data[strconv.FormatUint(id, 10)] = task.Result
		removeTaskContext(task)
	}

	jsonResponse, _ := json.Marshal(&responseCollectData)
	http.Error(w, string(jsonResponse), http.StatusOK)
}

func handleMailTemplate(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Method + " /mail/from_template " + r.RemoteAddr)

	if !setupCORS(&w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, InvalidRequestMethodPost, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UnableToGetRequestBody, http.StatusInternalServerError)
		return
	}

	var templateMailRequest MailTemplate
	err = json.Unmarshal(body, &templateMailRequest)
	if err != nil {
		http.Error(w, InvalidRequestData, http.StatusBadRequest)
		return
	}

	filePath := serverConfig.TemplatesDir + "/" + templateMailRequest.Template + ".template"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Template doesn't exist", http.StatusBadRequest)
		return
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Can't read template", http.StatusInternalServerError)
		return
	}

	response := make(map[string]string)
	for _, fillResult := range templateMailRequest.FillResults {
		template := string(file)
		if len(fillResult.ContactPersons) == 0 {
			fmt.Errorf("Empty persons")
			continue
		}

		if len(fillResult.Emails) == 0 {
			fmt.Errorf("Empty mails")
			continue
		}

		person := fillResult.ContactPersons[0]
		if len(person) == 0 {
			fmt.Errorf("Empty persons")
			template = strings.Replace(template, "{{??????}}", fillResult.CompanyName, -1)
			template = strings.Replace(template, "{{??????????????}}", "", -1)
		} else {
			personNameParts := strings.Split(person, " ")
			name := personNameParts[0]
			surname := personNameParts[1]

			template = strings.Replace(template, "{{??????}}", name, -1)
			template = strings.Replace(template, "{{??????????????}}", surname, -1)
		}

		template = strings.Replace(template, "{{??????}}", "2021", -1)
		template = strings.Replace(template, "{{??????????}}", templateMailRequest.Product, -1)
		template = strings.Replace(template, "{{????????}}", fillResult.AverageCapitalization, -1)

		log.Println("Template render will be sent to " + fillResult.Emails[0])
		log.Println("Text: " + template)

		response[fillResult.Emails[0]] = template
	}

	jsonResponse, _ := json.Marshal(response)
	http.Error(w, string(jsonResponse), http.StatusOK)
}

var serverConfig ServerConfiguration

func StartServer(configuration *ServerConfiguration) {

	var err error

	if configuration == nil {
		log.Fatal("Unable to run server")
	}

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/status/tasks", handleTasksStatus)
	http.HandleFunc("/data/collect", handleDataCollect)
	http.HandleFunc("/mail/from_template", handleMailTemplate)

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
		Addr:           ":" + configuration.Port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	InitContext()

	serverConfig = *configuration

	if configuration.UseTls {
		err = s.ListenAndServeTLS(configuration.TlsCrtFile, configuration.TlsKeyFile)
	} else {
		err = s.ListenAndServe()
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

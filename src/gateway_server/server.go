package main

import (
  "log"
  "fmt"
  "os"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "github.com/cloudfoundry-incubator/bbs/models"
  "gateway_server/bbs_client"
  "github.com/gorilla/mux"
)

type Configuration struct {
  BBSAddress  string
  CertFile    string
  KeyFile     string
}

var conf Configuration

func Log(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	      handler.ServeHTTP(w, r)
    })
}

func main() {

  // file, _ := os.Open("/Users/danhigham/workspace/go/src/bbs-client/conf.json")
  file, _ := os.Open("/var/vcap/jobs/gateway_server/config/conf.json")
  decoder := json.NewDecoder(file)
  conf = Configuration{}
  err := decoder.Decode(&conf)

  if err != nil {
    log.Fatal("Error loading configuration: ", err)
  }

  rtr := mux.NewRouter()
  rtr.PathPrefix("/firehose/").Handler(http.StripPrefix("/firehose/", http.FileServer(http.Dir("/var/vcap/packages/gateway_server/firehose"))))
  // rtr.PathPrefix("/firehose/").Handler(http.StripPrefix("/firehose/", http.FileServer(http.Dir("./firehose/"))))
  rtr.HandleFunc("/tasks/{id:[a-z|0-9|-]+}", getTask).Methods("GET")
  rtr.HandleFunc("/tasks", postTask).Methods("POST")
  http.Handle("/", rtr)

  if err := http.ListenAndServe(":8080", Log(http.DefaultServeMux)); err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}

func postTask(w http.ResponseWriter, r *http.Request) {
  body, _ := ioutil.ReadAll(r.Body)

  var task models.TaskDefinition
  err := json.Unmarshal(body, &task)

  if err != nil {
    http.Error(w, fmt.Sprintf("500 : %+v", err), 500)
    return
  }

  client, clientErr := bbs_client.NewBBSClient(conf.BBSAddress, conf.CertFile, conf.KeyFile)
  if clientErr != nil {
    http.Error(w, fmt.Sprintf("500 : %+v", clientErr), 500)
    return
  }

  guid := client.DesireTask(&task)

  w.Write([]byte(guid))
}

func getTask(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  id := params["id"]

  log.Printf("Getting task %s", id)

  client, clientErr := bbs_client.NewBBSClient(conf.BBSAddress, conf.CertFile, conf.KeyFile)

  if clientErr != nil {
    http.Error(w, fmt.Sprintf("500 : %+v", clientErr), 404)
    return
  }

  task, err := client.GetTask(id)

  if err != nil {
    http.Error(w, "404 Not Found : Error'd while getting task.", 404)
    return
  }

  json, _ := json.Marshal(task)
  w.Write(json)
}

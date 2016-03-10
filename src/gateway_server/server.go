package main

import (
  "log"
  "os"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "github.com/cloudfoundry-incubator/bbs/models"
  "bbs-client/bbs_client"
)

type Configuration struct {
  BBSAddress  string
  CertFile    string
  KeyFile     string
}

var conf Configuration

func main() {

  file, _ := os.Open("conf.json")
  decoder := json.NewDecoder(file)
  conf = Configuration{}
  err := decoder.Decode(&conf)

  http.Handle("/firehose/", http.StripPrefix("/firehose/", http.FileServer(http.Dir("./firehose"))))
  http.HandleFunc("/tasks/", handleTask)

  if err := http.ListenAndServe(":8080", nil); err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}

func handleTask(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      // Serve the resource.
    case "POST":
      body, _ := ioutil.ReadAll(r.Body)

      var task models.TaskDefinition
      err := json.Unmarshal(body, &task)

      if err != nil {
        log.Println("error:", err)
      }

      client := bbs_client.NewBBSClient(conf.BBSAddress, conf.CertFile, conf.KeyFile)
      client.DesireTask(&task)

    case "PUT":
       // Update an existing record.
    case "DELETE":
       // Remove the record.
    default:
       // Give an error message.
  }
}

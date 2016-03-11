package main

import (
  "net/http"
  "crypto/tls"
  "io/ioutil"
  "os"
  "bytes"
  "log"
)

func main() {
  url := os.Args[1]
  json_file := os.Args[2]

  json, _ := ioutil.ReadFile(json_file)

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
  req.Header.Set("Content-Type", "application/json")

  tr := &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }
  client := &http.Client{Transport: tr}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  log.Println("response Status:", resp.Status)
  log.Println("response Headers:", resp.Header)
  body, _ := ioutil.ReadAll(resp.Body)
  log.Println("response Body:", string(body))
}

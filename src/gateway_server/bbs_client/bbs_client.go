package bbs_client

import (
	"github.com/cloudfoundry-incubator/bbs"
	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/nu7hatch/gouuid"
)

type BBSClient struct {
  clientCertFile  string
  clientKeyFile   string
  bbsUrl          string
  client          bbs.Client
}

func (b *BBSClient) Connect() {
  client, _ := bbs.NewSecureSkipVerifyClient(b.bbsUrl, b.clientCertFile, b.clientKeyFile, 8192, 10)
	b.client = client
}


func (b *BBSClient) DesireTask(task *models.TaskDefinition) string {
  u, _ := uuid.NewV4()
  guid := u.String()

  err := b.client.DesireTask(guid, "some-domain", task)
	if err == nil {
		return guid
	}

  return ""
}

func (b *BBSClient) GetTask(id string) *models.Task {
  task, err := b.client.TaskByGuid(id)

  if err == nil  {
    return task
  }

	return &models.Task{}
}

func NewBBSClient(url, cert, key string) BBSClient {
  client := BBSClient{bbsUrl: url, clientCertFile: cert, clientKeyFile: key}
  client.Connect()

	return client
}

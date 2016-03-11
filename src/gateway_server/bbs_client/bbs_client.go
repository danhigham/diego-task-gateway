package bbs_client

import (
	"github.com/cloudfoundry-incubator/bbs"
	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/nu7hatch/gouuid"
	"fmt"
)

type BBSClient struct {
  clientCertFile  string
  clientKeyFile   string
  bbsUrl          string
  client          bbs.Client
}

func (b *BBSClient) Connect() error {
  client, err := bbs.NewSecureSkipVerifyClient(b.bbsUrl, b.clientCertFile, b.clientKeyFile, 8192, 10)

	if err != nil {
		return err
	}

	b.client = client
	return nil
}

func (b *BBSClient) DesireTask(task *models.TaskDefinition) string {
  u, _ := uuid.NewV4()
  guid := u.String()

  err := b.client.DesireTask(guid, "some-domain", task)
	if err == nil {
		return guid
	}

  return fmt.Sprintf("%+v", err)
}

func (b *BBSClient) GetTask(id string) (*models.Task, error) {
  task, err := b.client.TaskByGuid(id)

  if err == nil  {
    return task, nil
  }

	return nil, err
}

func NewBBSClient(url, cert, key string) (*BBSClient, error) {
  client := &BBSClient{bbsUrl: url, clientCertFile: cert, clientKeyFile: key}
  err := client.Connect()

	if err != nil {
		return nil, err
	}

	return client, nil
}

package main

import (
	"context"

	"http_service/pkg/api"
	"http_service/pkg/repository"
	"http_service/pkg/service"

	"github.com/sirupsen/logrus"
)

func main() {
	client, err := repository.NewDateBase(context.Background(), repository.Config{Host: "localhost", Port: "27017"})
	if err != nil {
		logrus.Error(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
	repos := repository.NewRepository(client)
	services := service.NewService(repos)
	api := api.NewApiUser(context.Background(), services)

	if err := api.Run("8080"); err != nil {
		logrus.Error(err)
	}

}

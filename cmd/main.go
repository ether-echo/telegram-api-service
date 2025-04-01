package main

import (
	"context"
	"github.com/ether-echo/telegram-api-service/internal/handler"
	"github.com/ether-echo/telegram-api-service/pkg/config"
	"github.com/ether-echo/telegram-api-service/pkg/logger"
)

var (
	log = logger.Logger().Named("main").Sugar()
)

func main() {

	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hand := handler.NewHandler(ctx, conf)

	hand.RegisterHandler(ctx)

}

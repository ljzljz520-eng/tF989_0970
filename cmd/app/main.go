package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"followuplabel/internal/controller"
	"followuplabel/internal/repository"
	"followuplabel/internal/service"
)

func main() {
	store := repository.NewMemory(repository.FixtureTags())
	plugin := service.NewPlugin(store, store, nil)
	if err := plugin.Install(context.Background()); err != nil {
		log.Fatal(err)
	}

	tagService := service.NewTags(store, store)
	address := os.Getenv("LABEL_PLUGIN_ADDR")
	if address == "" {
		address = ":8080"
	}

	displayAddress := address
	if strings.HasPrefix(address, ":") {
		displayAddress = "localhost" + address
	}
	log.Printf("客户回访标签插件已启动: http://%s", displayAddress)
	if err := http.ListenAndServe(address, controller.NewHTTP(tagService)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

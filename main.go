package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goapi/config"
	"goapi/resources"
	"goapi/services"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func configureRouter() *gin.Engine {
	router := gin.Default()
	resources.RegisterHandlers(router, services.NewDocumentServiceImpl())
	return router
}

func retrieveConfig() (*config.Config, error) {
	//open the configuration file
	f, err := os.Open("config.yml")
	if err != nil {
		log.Print("Cannot find config.yml")
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err == nil {
			log.Print("Cannot close config.yml")
		}
	}(f)

	//read the configuration file
	var cfg config.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Print(fmt.Sprintf("Error parsing config.yml %v", err))
		return nil, err
	}
	return &cfg, nil
}

func main() {
	//retrieve the configuration
	configuration, err := retrieveConfig()
	if err != nil {
		return
	}

	//configure the router
	router := configureRouter()

	//run the server
	err = router.Run(":" + configuration.Server.Port)
	if err != nil {
		return
	}
}

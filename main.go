package main

import (
	"bytes"
	"context"
	"goapi/config"
	"goapi/database"
	_ "goapi/docs/apis"
	"goapi/kafka"
	"goapi/repositories/repodocuments"
	"goapi/resources/documents"
	"goapi/resources/emails"
	"goapi/services/servicedocuments"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leekchan/gtf"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gopkg.in/yaml.v2"
)

func configureRouter(configuration *config.Config) *gin.Engine {
	router := gin.Default()
	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{"*"}
	configCors.AllowHeaders = []string{"*"}
	configCors.AllowCredentials = true
	router.Use(cors.New(configCors))

	//register document resource endpoints
	documents.RegisterHandlers(router, servicedocuments.NewDocumentServiceImpl(repodocuments.CreateDocumentRepository(configuration)))
	//register Email resource
	emails.RegisterHandlers(router, configuration)

	// @title Swagger REST API Documentation
	// @version 1.0
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

func retrieveConfig() (*config.Config, error) {
	//that's a pain in the ass, but the yaml module spec do not include the handling of env vars
	//so use go template module to parse the yaml configuration file, and inject the variable values
	data := struct {
		MONGO_SERVER_HOST     string
		MONGO_SERVER_PORT     string
		STORAGE_MEMORY        string
		EMAIL_CONSUMERS       string
		KAFKA_SERVER_HOST     string
		KAFKA_SERVER_PORT     string
		EMAIL_SERVER_HOST     string
		EMAIL_SERVER_PORT     string
		EMAIL_SERVER_USERNAME string
		EMAIL_SERVER_PASSWORD string
		EMAIL_SERVER_STARTTLS string
	}{
		MONGO_SERVER_HOST:     os.Getenv("MONGO_SERVER_HOST"),
		MONGO_SERVER_PORT:     os.Getenv("MONGO_SERVER_PORT"),
		STORAGE_MEMORY:        os.Getenv("STORAGE_MEMORY"),
		EMAIL_CONSUMERS:       os.Getenv("EMAIL_CONSUMERS"),
		KAFKA_SERVER_HOST:     os.Getenv("KAFKA_SERVER_HOST"),
		KAFKA_SERVER_PORT:     os.Getenv("KAFKA_SERVER_PORT"),
		EMAIL_SERVER_HOST:     os.Getenv("EMAIL_SERVER_HOST"),
		EMAIL_SERVER_PORT:     os.Getenv("EMAIL_SERVER_PORT"),
		EMAIL_SERVER_USERNAME: os.Getenv("EMAIL_SERVER_USERNAME"),
		EMAIL_SERVER_PASSWORD: os.Getenv("EMAIL_SERVER_PASSWORD"),
		EMAIL_SERVER_STARTTLS: os.Getenv("EMAIL_SERVER_STARTTLS"),
	}

	fileData, _ := ioutil.ReadFile("config.yml")
	var finalData bytes.Buffer
	t := template.New("config")
	t, err := t.Funcs(gtf.GtfTextFuncMap).Parse(string(fileData))
	if err != nil {
		panic(err)
	}
	err = t.Execute(&finalData, data)
	if err != nil {
		panic(err)
	}
	var cfg config.Config
	err = yaml.Unmarshal(finalData.Bytes(), &cfg)
	if err != nil {
		panic(err)
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
	router := configureRouter(configuration)

	//start everything
	startEveryThing(configuration)

	srv := &http.Server{
		Addr:    ":" + configuration.ServerConfig.Port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutting down server...")

	stopEveryThing(configuration)

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Info("Server exiting")
}

func startEveryThing(configuration *config.Config) {
	//Create a connection to DB if needed
	if !configuration.StorageInMemory {
		database.GetMongoDatabaseHandler().TryOrRetryCreateConnection(&configuration.DbConfig)
	}
	//start consumers
	kafka.GetInstanceKafkaConsumers(configuration).StartConsumers(configuration.EmailConsumers, kafka.EmailConsumer)
}

func stopEveryThing(configuration *config.Config) {
	//stop consumers
	kafka.GetInstanceKafkaConsumers(configuration).StopConsumers(configuration.EmailConsumers, kafka.EmailConsumer)
	//stop DB client
	if !configuration.StorageInMemory {
		database.GetMongoDatabaseHandler().Close()
	}
}

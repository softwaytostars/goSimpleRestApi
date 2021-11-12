package main

import (
	"bytes"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leekchan/gtf"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"goapi/config"
	_ "goapi/docs/resourcedocument"
	"goapi/repositories"
	"goapi/resources"
	"goapi/services"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
)

func configureRouter(configuration *config.Config) *gin.Engine {
	router := gin.Default()
	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{"*"}
	configCors.AllowHeaders = []string{"*"}
	configCors.AllowCredentials = true
	router.Use(cors.New(configCors))

	//register document resource endpoints
	resources.RegisterHandlers(router, services.NewDocumentServiceImpl(repositories.CreateDocumentRepository(configuration)))

	// @title Swagger REST API Documentation
	// @version 1.0
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

func retrieveConfig() (*config.Config, error) {
	//that's a pain in the ass, but the yaml module spec do not include the handling of env vars
	//so use go template to parse the yaml configuration file, and inject the variable values
	data := struct {
		MONGO_SERVER_HOST string
		MONGO_SERVER_PORT string
		STORAGE_MEMORY    string
	}{
		MONGO_SERVER_HOST: os.Getenv("MONGO_SERVER_HOST"),
		MONGO_SERVER_PORT: os.Getenv("MONGO_SERVER_PORT"),
		STORAGE_MEMORY:    os.Getenv("STORAGE_MEMORY"),
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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
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

func stopEveryThing(configuration *config.Config) {
	repositories.CloseRepositories(configuration)
}

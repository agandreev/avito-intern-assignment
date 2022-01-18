package main

import (
	"context"
	"fmt"
	"github.com/agandreev/avito-intern-assignment/internal/controller"
	"github.com/agandreev/avito-intern-assignment/internal/handlers"
	"github.com/agandreev/avito-intern-assignment/internal/repository"
	"github.com/agandreev/avito-intern-assignment/internal/service"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	logFile    = "logs.txt"
	configPath = "config.env"
	apiKeyTag  = "API_KEY"
	dbUser     = "DB_USER"
	dbPSWD     = "DB_PSWD"
	dbName     = "DB_NAME"
	dbPort     = "DB_PORT"
	srvPort    = "SRV_PORT"
)

// @title Balance control API
// @version 1.0
// @description This is a multi-user balance control system.

// @host localhost:8000
// @BasePath /
func main() {
	logger := logrus.New()
	// add file logger
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	defer file.Close()
	mw := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(mw)

	// load config
	apiKey, port, dbConfig, err := loadConfig()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	// create service and run server
	exchange := service.NewExchangeAPI(apiKey)
	gbStorage := repository.NewGrossBookStorage(*dbConfig)
	if err = gbStorage.Connect(); err != nil {
		logger.Fatal(err)
	}

	gb := service.NewGrossBook(gbStorage, exchange, logger)
	handler := handlers.NewHandler(gb, logger)
	srv := controller.NewServer(*handler)
	go func() {
		if err = srv.Run(port); err != nil {
			logger.Fatalf("ERROR: running server is failed <%s>", err)
		}
	}()
	logger.Print("Server is running")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Print("Server is shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		logger.Errorf("ERROR: graceful shutdown is broken <%s>", err.Error())
	}
}

// loadString loads a string value from config
func loadString(name string) (string, error) {
	value, ok := viper.Get(name).(string)
	if !ok {
		return "", fmt.Errorf("invalid %s type assertion", name)
	}
	return value, nil
}

// loadConfig loads all values from config
func loadConfig() (string, string, *repository.ConnectionConfig, error) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return "", "", nil, fmt.Errorf("can't load config: %w", err)
	}
	key, err := loadString(apiKeyTag)
	if err != nil {
		return "", "", nil, fmt.Errorf("can't load api key: %w", err)
	}
	port, err := loadString(srvPort)
	if err != nil {
		return "", "", nil, fmt.Errorf("can't load server port: %w", err)
	}
	dbConfig, err := loadDBVars()
	if err != nil {
		return "", "", nil, fmt.Errorf("can't load db vars: %w", err)
	}
	return key, port, dbConfig, nil
}

// loadDBVars loads all db values from config as repository.ConnectionConfig
func loadDBVars() (*repository.ConnectionConfig, error) {
	user, err := loadString(dbUser)
	if err != nil {
		return nil, err
	}
	password, err := loadString(dbPSWD)
	if err != nil {
		return nil, err
	}
	name, err := loadString(dbName)
	if err != nil {
		return nil, err
	}
	port, err := loadString(dbPort)
	if err != nil {
		return nil, err
	}
	orderConfig := &repository.ConnectionConfig{
		Username: user,
		Password: password,
		NameDB:   name,
		Port:     port,
	}
	return orderConfig, nil
}

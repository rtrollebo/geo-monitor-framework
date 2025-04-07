package system

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func ReadFile[T any](name string) ([]T, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jsonData []T
	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func WriteFile[T any](data []T, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func InitApp(appname string) context.Context {
	logFileName := appname + ".log"
	fileLogging, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer fileLogging.Close()

	var LogInfo *log.Logger
	var LogWarning *log.Logger
	var LogError *log.Logger

	LogInfo = log.New(fileLogging, fmt.Sprintf("%s - INFO:", appname), log.Ldate|log.Ltime)
	LogWarning = log.New(fileLogging, fmt.Sprintf("%s - WARNING:", appname), log.Ldate|log.Ltime)
	LogError = log.New(fileLogging, fmt.Sprintf("%s - ERROR:", appname), log.Ldate|log.Ltime)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "loginfo", LogInfo)
	ctx = context.WithValue(ctx, "logwarning", LogWarning)
	ctx = context.WithValue(ctx, "logerror", LogError)
	return ctx
}

func GetLog(ctx context.Context, logName string) *log.Logger {
	logInfo := ctx.Value(logName)
	if logInfo == nil {
		return nil
	}
	return logInfo.(*log.Logger)
}

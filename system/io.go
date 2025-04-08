package system

import (
	"context"
	"encoding/json"
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

func InitLogging(logFilePath string) (*log.Logger, *log.Logger, *log.Logger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, nil, err
	}

	infoLogger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	warnLogger := log.New(logFile, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger := log.New(logFile, "ERROR: ", log.Ldate|log.Ltime)

	return infoLogger, warnLogger, errorLogger, nil
}

func GetLog(ctx context.Context, logName string) *log.Logger {
	logInfo := ctx.Value(logName)
	if logInfo == nil {
		return nil
	}
	return logInfo.(*log.Logger)
}

func SetLog(ctx context.Context, logName string, logger *log.Logger) {
	ctx = context.WithValue(ctx, logName, logger)
}

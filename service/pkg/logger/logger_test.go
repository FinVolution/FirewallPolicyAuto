package logger

import (
	"fmt"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logFileName := "./log.log"
	logLevel := "DEBUG"
	fileSize := 100
	backups := 2
	age := 1

	InitLogger(logFileName, logLevel, fileSize, backups, age)

	startTime := time.Now()
	for i := 0; i < 10; i++ {
		Infof("testing infof log")
	}
	for i := 0; i < 10; i++ {
		Printf("testing printf log")
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("time used:", duration)
}

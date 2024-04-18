package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"syscall"
)

var logFile *lumberjack.Logger

func init_log() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logFile = &lumberjack.Logger{
		Filename:   "task.log",
		MaxSize:    1, // MB
		MaxBackups: 5,
		Compress:   false,
	}

	logrus.SetOutput(logFile)
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.Level(Conf.Log.Level))

	logrus.SetReportCaller(true)
}

func main() {

	err := LoadConf()
	if err != nil {
		return
	}

	init_log()
	defer func() {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	DispatchTask()

	//time.Sleep(time.Second * 3600)

	waitForQuit()
	logrus.Info("exit.......")
}

func waitForQuit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	<-done

	fmt.Println("quit ....")
}

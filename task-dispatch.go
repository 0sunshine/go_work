package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

func DispatchTask() {
	for _, task := range Conf.Tasks {
		taskSeq := 1
		for taskSeq <= task.Replicas {
			go doTask(taskSeq, task)
			taskSeq++
		}
	}
}

func doTask(taskSeq int, task TasksYamlConf) {
	repeat := task.Repeat

	logHeader := fmt.Sprintf("[task name: %s, task seq: %d]--", task.Name, taskSeq)
	for repeat != 0 {
		if repeat > 0 {
			repeat--
		}

		for _, step := range task.Steps {
			err := doStep(logHeader, taskSeq, step)
			if err != nil {
				logrus.Errorf("step error,quit ............ %v", step)
				return
			}
		}
	}
}

func doStep(logHeader string, taskSeq int, step StepsTasksYamlConf) error {

	//命令填空，用于暂停
	if len(step.Cmd) == 0 {
		logrus.Infof("%s sleep: %dms", logHeader, time.Duration(step.Duration))
		time.Sleep(time.Millisecond * time.Duration(step.Duration))
		return nil
	}

	parameters := make(map[string]int)
	parameters["taskSeq"] = taskSeq

	cmd, err := FormatCmd(step.Cmd, parameters)
	if err != nil {
		logrus.Errorf("%s FormatCmd err, step.Cmd: %s, err: %s", logHeader, step.Cmd, err)
		return err
	}
	logrus.Infof("%s use cmd：%s", logHeader, cmd)

	//cmdObj := exec.Command(cmd)
	cmd_slice := safeSplit(cmd)

	cmdObj := exec.Command(cmd_slice[0], cmd_slice[1:]...)

	if step.Tostdio == 1 {
		cmdObj.Stderr = os.Stderr
		cmdObj.Stdout = os.Stdout
	}

	err = cmdObj.Start()

	if err != nil {
		logrus.Errorf("%s err:%s", logHeader, err)
		return err
	}

	if step.Duration > 0 {
		time.Sleep(time.Millisecond * time.Duration(step.Duration))

		if cmdObj.Process != nil {
			cmdObj.Process.Kill()
		}
	}

	err = cmdObj.Wait()
	logrus.Debugf("%sCommand finished with error: %v", logHeader, err)

	//var mutex sync.Mutex
	//isQuit := false
	//
	//go func() {
	//	time.Sleep(time.Millisecond * time.Duration(step.Duration))
	//	mutex.Lock()
	//	isQuit = true
	//	if
	//	logrus.Debugf("%s send ", logHeader, err)
	//	mutex.Unlock()
	//}()
	//
	//time.Sleep(time.Millisecond * time.Duration(step.Duration))
	//
	//err = cmdObj.Wait()
	//mutex.Lock()
	//isQuit = true
	//mutex.Unlock()
	//logrus.Debugf("%sCommand finished with error: %v", logHeader, err)
	//
	//endTime := time.Now().UnixMilli()
	//needSleepTime := endTime - startTime
	//
	//if needSleepTime > 0 { //过早结束，继续等待
	//	time.Sleep(time.Millisecond * time.Duration(needSleepTime))
	//}
	//
	//time.Sleep(time.Millisecond * time.Duration(step.Duration))

	return nil
}

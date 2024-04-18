package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type StepsTasksYamlConf struct {
	Cmd      string `yaml:"cmd"`
	Duration int    `yaml:"duration"`
	Tostdio  int    `yaml:"tostdio"`
}

type TasksYamlConf struct {
	Name     string               `yaml:"name"`
	Replicas int                  `yaml:"replicas"`
	Repeat   int                  `yaml:"repeat"`
	Delay    int                  `yaml:"delay"`
	Steps    []StepsTasksYamlConf `yaml:"steps"`
}

type LogYamlConf struct {
	Level int `yaml:"level"`
}

type YamlConf struct {
	Log   LogYamlConf     `yaml:"log"`
	Tasks []TasksYamlConf `yaml:"tasks"`
}

var Conf YamlConf

func LoadConf() error {
	f, err := os.Open("./conf.yaml")
	if err != nil {
		fmt.Println(err)
		return err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = yaml.Unmarshal(b, &Conf)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("%#v\n", Conf)

	//Conf = YamlConf{
	//	Play: []PlayYamlConf{
	//		PlayYamlConf{},
	//	},
	//	Log: LogYamlConf{Level: 5},
	//}
	//bytes, err := yaml.Marshal(Conf)
	//str := string(bytes)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//fmt.Println(str)

	return nil
}

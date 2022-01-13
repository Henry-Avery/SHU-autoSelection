package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

var Conf Configuration

type Configuration struct {
	TermId   string   `json:"termId" yaml:"termId"`
	Username string   `json:"username" yaml:"username"`
	Password string   `json:"password" yaml:"password"`
	Courses  []Course `json:"courses" yaml:"courses"`
	Multithreading bool `json:"multithreading" yaml:"multithreading"`
	LoopTrail      bool `json:"loopTrail" yaml:"loopTrail"`
}

type Course struct {
	CourseId  string `json:"courseId" yaml:"courseId"`
	TeacherNo string `json:"teacherNo" yaml:"teacherNo"`
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() >> 1) // Prevent explosion

	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		panic(err)
	}

	if Conf.TermId == "3" {
		Conf.TermId = strconv.Itoa(time.Now().Year()-1) + Conf.TermId
	} else {
		Conf.TermId = strconv.Itoa(time.Now().Year()) + Conf.TermId
	}
}

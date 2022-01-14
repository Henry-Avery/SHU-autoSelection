package main

import (
	"log"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

const LoginUrl = "https://oauth.shu.edu.cn/login/eyJ0aW1lc3RhbXAiOjE2MjE2MDUzMTQ1ODcxMTM5MDMsInJlc3BvbnNlVHlwZSI6ImNvZGUiLCJjbGllbnRJZCI6InlSUUxKZlVzeDMyNmZTZUtOVUN0b29LdyIsInNjb3BlIjoiIiwicmVkaXJlY3RVcmkiOiJodHRwOi8veGsuYXV0b2lzcC5zaHUuZWR1LmNuL3Bhc3Nwb3J0L3JldHVybiIsInN0YXRlIjoiIn0="
const TermSelectUrl = "http://xk.autoisp.shu.edu.cn/Home/TermSelect"
const CourseSelectionSaveUrl = "http://xk.autoisp.shu.edu.cn/CourseSelectionStudent/CourseSelectionSave"
const QueryCourseCheckUrl = "http://xk.autoisp.shu.edu.cn/CourseSelectionStudent/QueryCourseCheck"

var times int64
var selected = make(map[string]bool)

func main() {
	c := colly.NewCollector(
	// colly.Debugger(&debug.LogDebugger{}),
	)
	extensions.RandomUserAgent(c)
	Login(c)

	OnQueryCallbacks(c) // Register hook functions

	c.Async = Conf.Multithreading // Use multithreading: true or false

	for Conf.EndlessLoop {
		QueryCourse(c) // start loop trailing
	}
	QueryCourse(c) // just once
}

// Login try to log in to the xk.autoisp.shu.edu.cn.
func Login(c *colly.Collector) {
	// Check login status
	//c.OnResponse(func(r *colly.Response) {
	//	if strings.Contains(string(r.Body), "上海大学本硕博一体化选课系统") {
	//		log.Println("Login successful")
	//	}
	//})

	err := c.Post(LoginUrl, map[string]string{
		"username": Conf.Username,
		"password": EncryptPassword(Conf.Password),
	})
	if err != nil {
		panic(err)
	}

	err = c.Post(TermSelectUrl, map[string]string{"termId": Conf.TermId})
	if err != nil {
		panic(err)
	}
}

// OnQueryCallbacks registers a function.
// It will save the course on every query if the course is not full.
func OnQueryCallbacks(c *colly.Collector) {
	c.OnHTML("#tblcoursecheck > tbody > tr:nth-child(2) > td:nth-child(2)",
		func(e *colly.HTMLElement) {
			defer func() {
				if info := recover(); info != nil {
					log.Println(info)
				}
				c.Async = Conf.Multithreading
			}()
			for _, course := range Conf.Courses {
				if !strings.Contains(e.DOM.Text(), course.CourseId) || selected[course.CourseId] {
					continue
				}
				c.Async = false
				err := c.Post(CourseSelectionSaveUrl, map[string]string{
					"cids": course.CourseId,
					"tnos": course.TeacherNo,
				})
				if err != nil {
					panic(err)
				}
				c.Async = Conf.Multithreading
				atomic.AddInt64(&times, 1)
				log.Println("=== The", times, "times trying,", course.CourseId, "selection successful! ===")
				selected[course.CourseId] = true
			}
		})
}

// QueryCourse will try to query every course status.
// If any course is able to save, it will be hooked by QueryCallbacks.
func QueryCourse(c *colly.Collector) {
	for _, course := range Conf.Courses {
		err := c.Post(QueryCourseCheckUrl, map[string]string{
			"CID":            course.CourseId,
			"TeachNo":        course.TeacherNo,
			"FunctionString": "LoadData",
			"IsNotFull":      strconv.FormatBool(Conf.EndlessLoop),
			"PageIndex":      "1",
			"PageSize":       "10",
		})
		if err != nil {
			panic(err)
		}
		atomic.AddInt64(&times, 1)
		log.Println("The", times, "times trying, the course", course.CourseId, "is full")
	}
}

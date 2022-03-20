package main

import (
	"strings"
	"sync"
	"testing"
)

func TestProcess(t *testing.T) {
	cases := []string{
		"adjust.com",
		"google.com",
		"facebook.com",
		"yahoo.com",
		"yandex.com",
		"reddit.com/r/funny",
		"reddit.com/r/notfunny",
		"baroquemusiclibrary.com",
	}

	limitChan := make(chan struct{}, 10)
	resultChan := make(chan string, len(cases))
	var waitGroup sync.WaitGroup

	for _, testCase := range cases {
		limitChan <- struct{}{}

		waitGroup.Add(1)
		processHost(testCase, limitChan, resultChan, &waitGroup)

		select {
		case res, ok := <-resultChan:
			if ok {
				if !strings.HasPrefix(testCase, "http://") {
					testCase = "http://" + testCase
				}
				if !strings.HasPrefix(res, testCase) {
					t.Errorf("incorrect result for host %s result: %s", testCase, res)
				}
			} else {
				t.Error("result channel closed")
			}
		default:
			t.Errorf("response wasnt't getted for host %s", testCase)
		}
	}
}

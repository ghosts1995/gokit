/*
 * A toolkit for Golang development
 * https://www.likexian.com/
 *
 * Copyright 2019, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */

package workqueue

import (
	"fmt"
	"github.com/likexian/gokit/assert"
	"github.com/likexian/gokit/xfile"
	"github.com/likexian/gokit/xhttp"
	"os"
	"testing"
)

func TestVersion(t *testing.T) {
	assert.NotEqual(t, Version(), "")
	assert.NotEqual(t, Author(), "")
	assert.NotEqual(t, License(), "")
}

func TestMathSum(t *testing.T) {
	mathPlus := func(t Task) Task {
		return t.(int) + 1
	}

	mathSum := func(r Task, t Task) Task {
		return r.(int) + t.(int)
	}

	wq := New(100)
	wq.SetWorker(mathPlus, 10)
	wq.SetMerger(mathSum, 0)

	for i := 0; i < 1000; i++ {
		wq.Add(i)
	}

	result := wq.Wait()
	assert.Equal(t, result, 500500)
}

func TestFileLine(t *testing.T) {
	defer os.RemoveAll("tmp")

	lineCount := func(t Task) Task {
		ls, _ := xfile.ReadLines(t.(string), 0)
		return len(ls)
	}

	lineSum := func(r Task, t Task) Task {
		return r.(int) + t.(int)
	}

	wq := New(0)
	wq.SetWorker(lineCount, 0)
	wq.SetMerger(lineSum, 0)

	for i := 0; i < 100; i++ {
		xfile.WriteText(fmt.Sprintf("tmp/%d.txt", i), "0\n1\n2\n3\n4\n5\n6\n7\n8\n9")
	}

	files, err := xfile.ListDir("tmp", "file", -1)
	assert.Nil(t, err)
	for _, v := range files {
		wq.Add(v[1])
	}

	result := wq.Wait()
	assert.Equal(t, result, 1000)
}

func TestHttpStatus(t *testing.T) {
	getStatus := func(t Task) Task {
		rsp, err := xhttp.New().Do("GET", fmt.Sprintf("https://httpbin.org/status/%d", t.(int)))
		if err != nil {
			return 0
		}

		defer rsp.Close()
		return rsp.Response.StatusCode
	}

	sumStatus := func(r Task, t Task) Task {
		tt := t.(int)
		rr := r.(map[int]int)

		if _, ok := rr[tt]; !ok {
			rr[tt] = 0
		}
		rr[tt] += 1

		return r
	}

	wq := New(0)
	wq.SetWorker(getStatus, 100)
	wq.SetMerger(sumStatus, map[int]int{})

	tasks := map[int]int{200: 5, 206: 4, 401: 3, 403: 2, 405: 1}
	for k, v := range tasks {
		for i := 0; i < v; i++ {
			wq.Add(k)
		}
	}

	result := wq.Wait()
	assert.Equal(t, result, tasks)
}
package main

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIndexPage(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World", string(data))
}

func TestDecoHandler(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	buf := &bytes.Buffer{}
	log.SetOutput(buf)//표준 Log를 위한 destination set, 원래는 화면인데 버퍼로..

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)


	r := bufio.NewReader(buf)
	line, _, err := r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER2] Started")

	line, _, err = r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER1] Started")

	line, _, err = r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER1] Completed")

	line, _, err = r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER2] Completed")
}


//로그투가 로그원을 가지고 있고 로그원이 핸들러를 가지고 있어서 로그투가 먼저 불리고 그다음로그원스타티드 그다음 로그원컴플리트
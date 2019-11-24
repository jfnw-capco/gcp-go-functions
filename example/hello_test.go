package main

import (
	"net/http"
	"testing"

	common "github.com/capcodigital/gcpfunctions"
	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {

	_, recorder := common.ExecuteTestRequest(t, http.MethodGet, "/ping", "", common.Ping)
	assert.Equal(t, 200, recorder.Code)
}

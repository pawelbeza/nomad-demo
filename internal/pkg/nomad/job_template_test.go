package nomad

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFetcherCmd_NoScript(t *testing.T) {
	url := "www.test.com"
	script := false

	expectedCmd := "wget -T 5 -O - www.test.com > ${NOMAD_ALLOC_DIR}/index.html"
	actualCmd := CreateFetcherCmd(url, script)
	assert.Equal(t, expectedCmd, actualCmd)
}

func TestCreateFetcherCmd_Script(t *testing.T) {
	url := "www.test.com"
	script := true

	expectedCmd := "sh <(wget -T 5 -O - www.test.com) 2>&1 > ${NOMAD_ALLOC_DIR}/index.html"
	actualCmd := CreateFetcherCmd(url, script)
	assert.Equal(t, expectedCmd, actualCmd)
}

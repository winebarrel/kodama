package kodama_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/kodama"
)

func TestOptions_AfterApply_OK(t *testing.T) {
	assert := assert.New(t)

	options := &kodama.Options{
		NS: map[string]string{"ns.example.com": "203.0.113.0"},
	}

	err := options.AfterApply()
	assert.NoError(err)
}

func TestOptions_AfterApply_Err(t *testing.T) {
	assert := assert.New(t)

	options := &kodama.Options{
		NS: map[string]string{"ns.example.com": "999.999.999.999"},
	}

	err := options.AfterApply()
	assert.ErrorContains(err, "invalid IP address: 999.999.999.999")
}

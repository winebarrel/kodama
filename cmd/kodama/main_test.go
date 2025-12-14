package main_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	main "github.com/winebarrel/kodama/cmd/kodama"
)

func TestMain(t *testing.T) {
	assert := assert.New(t)

	orgArgs := os.Args
	t.Cleanup(func() { os.Args = orgArgs })
	os.Args = []string{"kodama", "--version"}

	defer func() {
		v := recover()
		assert.Contains(v, "os.Exit(0)")
	}()

	main.Main()
}

func TestParseArgs(t *testing.T) {
	assert := assert.New(t)

	orgArgs := os.Args
	t.Cleanup(func() { os.Args = orgArgs })
	os.Args = []string{"kodama", "--version"}

	defer func() {
		v := recover()
		assert.Contains(v, "os.Exit(0)")
	}()

	main.ParseArgs()
}

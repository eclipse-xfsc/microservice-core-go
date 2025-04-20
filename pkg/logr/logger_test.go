package logr

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
)

var testErrorEntry = "test error entry"
var testDebugEntry = "test debug entry"
var testInfoEntry = "test info entry"

func testDebug(logger *Logger) {
	logger.Debug(testDebugEntry)
	logger.Info(testInfoEntry)
	logger.Error(fmt.Errorf(testErrorEntry), "")
}

func setEnv() func() {
	initProd := viper.GetBool("RELEASE")
	initLev := viper.GetString("LOG_LEVEL")

	viper.Set("RELEASE", false)
	viper.Set("LOG_LEVEL", "DEBUG")

	return func() {
		viper.Set("RELEASE", initProd)
		viper.Set("LOG_LEVEL", initLev)
	}
}
func TestGetLogger(t *testing.T) {
	resetFn := setEnv()
	defer resetFn()

	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	logger, err := New(viper.GetString("LOG_LEVEL"), true, writer)
	if err != nil {
		t.Error(err)
	}

	testDebug(logger)
	err = writer.Flush()
	if err != nil {
		t.Error(err)
	}

	logs := strings.Split(buffer.String(), "\n")

	var res = make([]bool, 3)
	for i, l := range logs {
		if l == "" {
			logs = append(logs[:i], logs[i+1:]...)
		}
		if strings.HasPrefix(l, `{"L":"DEBUG"`) {
			res[i] = strings.Contains(l, testDebugEntry)
		}
		if strings.HasPrefix(l, `{"L":"INFO"`) {
			res[i] = strings.Contains(l, testInfoEntry)
		}
		if strings.HasPrefix(l, `{"L":"ERROR"`) {
			res[i] = strings.Contains(l, testErrorEntry)
		}
	}

	assert.Equal(t, len(logs), 3)
	assert.Equal(t, res, []bool{true, true, true})
}

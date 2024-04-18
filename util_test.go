package main

import (
	"testing"
)

func TestFormatCmd(test *testing.T) {
	parameters := make(map[string]int)
	parameters["taskSeq"] = 1

	cmd := "./ffmpeg -i ./pcr_error.ts -avcodec copy -f mpegts udp://226.0.0.{taskSeq}:{taskSeq+10010}"

	newCmd, err := FormatCmd(cmd, parameters)

	if err != nil {
		test.Logf("got err '%v'\n", err)
		test.Fail()
	}

	if newCmd != "./ffmpeg -i ./pcr_error.ts -avcodec copy -f mpegts udp://226.0.0.1:10011" {
		test.Logf("not expect")
		test.Fail()
	}
}

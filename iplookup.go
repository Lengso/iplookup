package iplookup

import (
	"bytes"
	"context"
	"github.com/Lengso/iplookup/pkg/runner"
	"io"
	"io/ioutil"
	"log"
)

func GetDomain(ip string) []string {
	options := runner.ParsePkgOptions()
	newRunner, err := runner.NewRunner(options)
	buf := bytes.Buffer{}
	err = newRunner.EnumerateSingleDomain(context.Background(), ip, []io.Writer{&buf})
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(&buf)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) != 0 {
		return []string{string(data)}
	}
	return nil
}

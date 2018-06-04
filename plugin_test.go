package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
)

func Test_CanGenerateCamelCaseFieldNames(t *testing.T) {
	fixture, _ := ioutil.ReadFile("./fixtures/money.pb")
	fds := new(descriptor.FileDescriptorSet)
	proto.Unmarshal(fixture, fds)

	req := new(plugin_go.CodeGeneratorRequest)
	req.ProtoFile = fds.GetFile()
	req.FileToGenerate = append(req.FileToGenerate, fds.GetFile()[0].GetName())

	plugin := &plugin{out: &bytes.Buffer{}}
	res, _ := plugin.Generate(req)
	content := res.GetFile()[0].GetContent()

	if !strings.Contains(content, "amountCents") {
		t.Errorf("Expected generated schema to contain camelCase field names, but got this instead: %v", content)
	}

	if strings.Contains(content, "amount_cents") {
		t.Errorf("Expected generated schema not to contain snake_case field names, but got this instead: \n%v", content)
	}
}
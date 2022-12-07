package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
)

func Test_SupportNewTimestampFormat(t *testing.T) {
	fixture, _ := ioutil.ReadFile("./fixtures/money.pb")
	fds := new(descriptor.FileDescriptorSet)
	proto.Unmarshal(fixture, fds)

	req := new(plugin_go.CodeGeneratorRequest)
	req.ProtoFile = fds.GetFile()
	req.FileToGenerate = append(req.FileToGenerate, fds.GetFile()[0].GetName())

	plugin := &plugin{out: &bytes.Buffer{}}
	res, _ := plugin.Generate(req)
	content := res.GetFile()[0].GetContent()

	if !strings.Contains(content, "scalar Timestamp") {
		t.Errorf("Expected the schema to include a scalar definition for Time, %s", content)
	}

	if !strings.Contains(content, "createdAt: Timestamp") {
		t.Errorf("Expected the schema to reflect the ISO8601 json serialization format, %s", content)
	}
}

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

func Test_CanGenerate(t *testing.T) {
	fixture, _ := ioutil.ReadFile("./fixtures/money.pb")
	fds := new(descriptor.FileDescriptorSet)
	proto.Unmarshal(fixture, fds)

	req := new(plugin_go.CodeGeneratorRequest)
	req.ProtoFile = fds.GetFile()
	req.FileToGenerate = append(req.FileToGenerate, fds.GetFile()[0].GetName())

	plugin := &plugin{out: &bytes.Buffer{}}
	res, _ := plugin.Generate(req)
	content := res.GetFile()[0].GetContent()

	t.Log(content)
}

func Test_CanAddComments(t *testing.T) {
	fixture, _ := ioutil.ReadFile("./fixtures/money.pb")
	fds := new(descriptor.FileDescriptorSet)
	proto.Unmarshal(fixture, fds)

	req := new(plugin_go.CodeGeneratorRequest)
	req.ProtoFile = fds.GetFile()
	req.FileToGenerate = append(req.FileToGenerate, fds.GetFile()[0].GetName())

	plugin := &plugin{out: &bytes.Buffer{}}
	res, _ := plugin.Generate(req)
	content := res.GetFile()[0].GetContent()

	if !strings.Contains(content, "foozles are the best") {
		t.Errorf("Expected generated schema to include comments, but got %s", content)
	}
}

func Test_FieldDeprecated(t *testing.T) {
	fixture, _ := ioutil.ReadFile("./fixtures/money.pb")
	fds := new(descriptor.FileDescriptorSet)
	proto.Unmarshal(fixture, fds)

	req := new(plugin_go.CodeGeneratorRequest)
	req.ProtoFile = fds.GetFile()
	req.FileToGenerate = append(req.FileToGenerate, fds.GetFile()[0].GetName())

	plugin := &plugin{out: &bytes.Buffer{}}
	res, _ := plugin.Generate(req)
	content := res.GetFile()[0].GetContent()

	if !strings.Contains(content, "foobar: String @deprecated") {
		t.Errorf("Expected generated schema to include deprecation descriptors, but got %s", content)
	}
}

func Test_InputDeprecated(t *testing.T) {
	fixture, _ := ioutil.ReadFile("./fixtures/money.pb")
	fds := new(descriptor.FileDescriptorSet)
	proto.Unmarshal(fixture, fds)

	req := new(plugin_go.CodeGeneratorRequest)
	req.ProtoFile = fds.GetFile()
	req.FileToGenerate = append(req.FileToGenerate, fds.GetFile()[0].GetName())

	plugin := &plugin{out: &bytes.Buffer{}}
	res, _ := plugin.Generate(req)
	content := res.GetFile()[0].GetContent()

	if strings.Contains(content, "yeOldMoneyBox: Input_core_apimessages_OldMoneyBox @deprecated") {
		t.Errorf("Expected generated schema to exclude input deprecation, but got %s", content)
	}
}

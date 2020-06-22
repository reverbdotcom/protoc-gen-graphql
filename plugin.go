package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"
)

var scalars = map[string]string{
	".google.protobuf.Timestamp": "Timestamp",
}

var primitives = map[descriptor.FieldDescriptorProto_Type]string{
	descriptor.FieldDescriptorProto_TYPE_BOOL:    "Boolean",
	descriptor.FieldDescriptorProto_TYPE_INT32:   "Int",
	descriptor.FieldDescriptorProto_TYPE_INT64:   "Int",
	descriptor.FieldDescriptorProto_TYPE_BYTES:   "String",
	descriptor.FieldDescriptorProto_TYPE_FLOAT:   "Float",
	descriptor.FieldDescriptorProto_TYPE_STRING:  "String",
	descriptor.FieldDescriptorProto_TYPE_DOUBLE:  "Float",
	descriptor.FieldDescriptorProto_TYPE_FIXED32: "Float",
	descriptor.FieldDescriptorProto_TYPE_FIXED64: "Float",
	descriptor.FieldDescriptorProto_TYPE_SINT32:  "Int",
	descriptor.FieldDescriptorProto_TYPE_SINT64:  "Int",
	descriptor.FieldDescriptorProto_TYPE_UINT32:  "Int",
	descriptor.FieldDescriptorProto_TYPE_UINT64:  "Int",
}

type plugin struct {
	out *bytes.Buffer
}

func (p *plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	params := parseParams(req.GetParameter())
	descriptors := protokit.ParseCodeGenRequest(req)
	resp := &plugin_go.CodeGeneratorResponse{}

	for _, d := range descriptors {
		p.printFile(d)
	}

	resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
		Name:    proto.String(params["file_out"]),
		Content: proto.String(p.out.String()),
	})

	return resp, nil
}

func (p *plugin) printFile(file *protokit.FileDescriptor) {
	for _, msg := range file.GetMessages() {
		p.printDescriptor(msg)
	}

	for _, s := range file.GetServices() {
		p.printService(s)
	}

	for _, e := range file.GetEnums() {
		p.printEnum(e)
	}
}

func (p *plugin) printEnum(enum *protokit.EnumDescriptor) {
	fmt.Fprintf(p.out, "enum %s {\n", underscore(enum.GetFullName()))

	for _, val := range enum.GetValues() {
		fmt.Fprintf(p.out, "  %s\n", val.GetName())
	}

	fmt.Fprintf(p.out, "}\n\n")
}

func (p *plugin) printService(srv *protokit.ServiceDescriptor) {
	fmt.Fprintf(p.out, "type %s {\n", underscore(srv.GetFullName()))

	for _, meth := range srv.GetMethod() {
		in := underscore(meth.GetInputType()[1:])
		out := underscore(meth.GetOutputType()[1:])
		fmt.Fprintf(p.out, "  %s(input: Input_%s): %s\n", meth.GetName(), in, out)
	}

	fmt.Fprintf(p.out, "}\n\n")
}

func (p *plugin) printDescriptor(desc *protokit.Descriptor) {
	for _, nested := range desc.GetMessages() {
		p.printDescriptor(nested)
	}

	if desc.GetFullName() == "google.protobuf.Timestamp" {
		fmt.Fprintf(p.out, "scalar Timestamp\n\n")
		return
	}

	for _, e := range desc.GetEnums() {
		p.printEnum(e)
	}

	for t, prefix := range map[string]string{
		"type":  "",
		"input": "Input_",
	} {
		name := fmt.Sprintf("%s%s", prefix, underscore(desc.GetFullName()))

		if len(desc.GetField()) == 0 {
			fmt.Fprintf(p.out, "%s %s { \n  _: Boolean\n}\n\n", t, name)
		} else {
			fmt.Fprintf(p.out, "%s %s {\n", t, name)

			for _, field := range desc.GetField() {
				fmt.Fprintf(p.out, "  %s: %s\n", field.GetJsonName(), typeName(field, prefix))
			}

			fmt.Fprintf(p.out, "}\n\n")
		}
	}
}

func parseParams(p string) map[string]string {
	params := make(map[string]string)

	if len(p) == 0 {
		return params
	}

	for _, part := range strings.Split(p, ",") {
		kv := strings.Split(part, "=")
		params[kv[0]] = kv[1]
	}

	return params
}

func resolveType(field *descriptor.FieldDescriptorProto, prefix string) string {
	if t, isPrimitive := primitives[field.GetType()]; isPrimitive {
		return t
	}

	if t, isScalar := scalars[field.GetTypeName()]; isScalar {
		return t
	}

	return fmt.Sprintf("%s%s", prefix, underscore(field.GetTypeName()[1:]))
}

func typeName(field *descriptor.FieldDescriptorProto, prefix string) string {
	if field.GetType() == descriptor.FieldDescriptorProto_TYPE_ENUM {
		prefix = ""
	}

	name := resolveType(field, prefix)

	if field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		return fmt.Sprintf("[%s]", name)
	}

	return name
}

func underscore(name string) string {
	return strings.Replace(name, ".", "_", -1)
}

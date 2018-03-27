package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"
)

func main() {
	// all the heavy lifting done for you!
	if err := protokit.RunPlugin(&plugin{out: &bytes.Buffer{}}); err != nil {
		log.Fatal(err)
	}
}

type plugin struct {
	out *bytes.Buffer
}

func (p *plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	descriptors := protokit.ParseCodeGenRequest(req)

	resp := &plugin_go.CodeGeneratorResponse{}

	// TODO: configurable
	fileName := "my.graphql"

	for _, d := range descriptors {
		p.printFile(d)
	}

	resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
		Name:    proto.String(fileName),
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

	fmt.Fprintln(p.out, "}\n")
}

func (p *plugin) printService(srv *protokit.ServiceDescriptor) {
	fmt.Fprintf(p.out, "type %s {\n", underscore(srv.GetFullName()))

	for _, meth := range srv.GetMethod() {
		in := underscore(meth.GetInputType()[1:])
		out := underscore(meth.GetOutputType()[1:])
		fmt.Fprintf(p.out, "  %s(input: Input_%s): %s\n", meth.GetName(), in, out)
	}

	fmt.Fprintln(p.out, "}\n")
}

func (p *plugin) printDescriptor(desc *protokit.Descriptor) {
	for _, t := range []string{"input", "type"} {
		name := underscore(desc.GetFullName())

		if len(desc.GetField()) == 0 {
			fmt.Fprintf(p.out, "%s %s\n", t, name)
		} else {
			fmt.Fprintf(p.out, "%s %s {\n", t, name)

			for _, field := range desc.GetField() {
				fmt.Fprintf(p.out, "  %s: %s\n", field.GetName(), typeName(field))
			}

			fmt.Fprintln(p.out, "}\n")
		}
	}
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

func typeName(field *descriptor.FieldDescriptorProto) string {
	var name string

	t, isPrimitive := primitives[field.GetType()]
	if isPrimitive {
		name = t
	} else {
		name = underscore(field.GetTypeName()[1:])
	}

	if field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		return fmt.Sprintf("[%s]", name)
	}

	return name
}

func underscore(name string) string {
	return strings.Replace(name, ".", "_", -1)
}

package svcdef

import (
	"reflect"
	"testing"

	"github.com/TuneLab/go-truss/gengokit/gentesthelper"
	"github.com/davecgh/go-spew/spew"
)

func basicFromString(t *testing.T) *Svcdef {
	defStr := `
		syntax = "proto3";

		// General package
		package general;

		import "google/api/annotations.proto";

		message SumRequest {
			int64 a = 1;
			int64 b = 2;
		}

		message SumReply {
			int64 v = 1;
			string err = 2;
		}

		service SumSvc {
			rpc Sum(SumRequest) returns (SumReply) {
				option (google.api.http) = {
					get: "/sum/{a}"
				};
			}
		}
	`
	sd, err := NewFromString(defStr)

	if err != nil {
		t.Fatal("Failed to create a svcdef from the definition string:", err)
	}
	return sd
}

func TestMessages(t *testing.T) {
	sd := basicFromString(t)
	expected := []*Message{
		&Message{
			Name: "SumRequest",
			Fields: []*Field{
				&Field{
					Name: "A",
					Type: &FieldType{
						Name:      "int64",
						Enum:      nil,
						Message:   nil,
						Map:       nil,
						StarExpr:  false,
						ArrayType: false,
					},
				},
				&Field{
					Name: "B",
					Type: &FieldType{
						Name:      "int64",
						Enum:      nil,
						Message:   nil,
						Map:       nil,
						StarExpr:  false,
						ArrayType: false,
					},
				},
			},
		},
		&Message{
			Name: "SumReply",
			Fields: []*Field{
				&Field{
					Name: "V",
					Type: &FieldType{
						Name:      "int64",
						Enum:      nil,
						Message:   nil,
						Map:       nil,
						StarExpr:  false,
						ArrayType: false,
					},
				},
				&Field{
					Name: "Err",
					Type: &FieldType{
						Name:      "string",
						Enum:      nil,
						Message:   nil,
						Map:       nil,
						StarExpr:  false,
						ArrayType: false,
					},
				},
			},
		},
	}

	if got, want := sd.Messages, expected; !reflect.DeepEqual(got, want) {
		diff := gentesthelper.DiffStrings(spew.Sdump(got), spew.Sdump(want))
		t.Errorf("got != want; methods differ: %v\n", diff)
	}
}

func TestHTTPBinding(t *testing.T) {
	sd := basicFromString(t)
	expected := []*HTTPBinding{
		&HTTPBinding{
			Verb: "get",
			Path: "/sum/{a}",
			Params: []*HTTPParameter{
				&HTTPParameter{
					Location: "path",
					Field: &Field{
						Name: "A",
						Type: &FieldType{
							Name: "int64",
						},
					},
				},
				&HTTPParameter{
					Location: "query",
					Field: &Field{
						Name: "B",
						Type: &FieldType{
							Name: "int64",
						},
					},
				},
			},
		},
	}
	output := sd.Service.Methods[0].Bindings
	if got, want := output, expected; !reflect.DeepEqual(got, want) {
		diff := gentesthelper.DiffStrings(spew.Sdump(got), spew.Sdump(want))
		t.Errorf("got != want; methods differ: %v\n", diff)
	}
}

func TestNoHTTPBinding(t *testing.T) {
	defstr := `
		syntax = "proto3";

		// General package
		package general;

		import "google/api/annotations.proto";

		message SumRequest {
			int64 a = 1;
			int64 b = 2;
		}

		message SumReply {
			int64 v = 1;
			string err = 2;
		}

		service SumSvc {
			rpc Sum(SumRequest) returns (SumReply) {}
		}
	`
	_, err := NewFromString(defstr)
	if err != nil {
		t.Fatal("Failed to create svcdef from string:", err)
	}
}

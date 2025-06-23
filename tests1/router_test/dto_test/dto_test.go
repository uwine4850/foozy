package dto_test

import (
	"os"
	"testing"

	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/router/rest"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

type TestMessage struct {
	rest.ImplementDTOMessage
	TypTestMessage rest.TypeId `dto:"-typeid"`
	Id             int         `dto:"Id"`
	Name           string      `dto:"Name"`
	Ok             bool        `dto:"Ok"`
}

type Test1Message struct {
	rest.ImplementDTOMessage
	TypTest1Message rest.TypeId `dto:"-typeid"`
	Id1             int         `dto:"Id1"`
	Name1           string      `dto:"Name1"`
	Ok1             bool        `dto:"Ok1"`
}

type TestNotSafeMessage struct {
	rest.ImplementDTOMessage
	TypTestNotSafeMessage rest.TypeId `dto:"-typeid"`
	Id                    int         `dto:"Id"`
}

type TestDeepNotSafeMessage struct {
	rest.ImplementDTOMessage
	TypTestDeepNotSafeMessage rest.TypeId             `dto:"-typeid"`
	Message                   TestInnerNotSafeMessage `dto:"Message"`
}

type TestInnerNotSafeMessage struct {
	rest.ImplementDTOMessage
	TypTestInnerNotSafeMessage rest.TypeId `dto:"-typeid"`
	Id                         int         `dto:"Id"`
}

var newDTO = rest.NewDTO()
var messages = map[string][]irest.IMessage{
	"test.ts": {
		TestMessage{},
		Test1Message{},
	},
}
var allowMessages = []rest.AllowMessage{
	{
		Package: "dto_test",
		Name:    "TestMessage",
	},
	{
		Package: "dto_test",
		Name:    "Test1Message",
	},
}

func TestMain(m *testing.M) {
	newDTO.AllowedMessages([]rest.AllowMessage{})
	newDTO.Messages(messages)
	newDTO.AllowedMessages(allowMessages)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGenerate(t *testing.T) {
	if err := newDTO.Generate(); err != nil {
		t.Error(err)
	}
	ok, err := filesAreEqual("test.ts", "expected_file.ts")
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("the generated file and the expected file do not match")
	}
}

func TestIsSafeMessage(t *testing.T) {
	var msg TestMessage
	if err := rest.IsSafeMessage(typeopr.Ptr{}.New(&msg), newDTO.GetAllowedMessages()); err != nil {
		t.Error(err)
	}
}

func TestIsSafeMessageError(t *testing.T) {
	var msg TestNotSafeMessage
	if err := rest.IsSafeMessage(typeopr.Ptr{}.New(&msg), newDTO.GetAllowedMessages()); err != nil {
		if err.Error() != "dto_test.TestNotSafeMessage message is unsafe" {
			t.Errorf("the error doesn't match the expected error. Unexpected error: %s", err.Error())
		}
	}
}

func TestDeepCheckSafeMessage(t *testing.T) {
	newDTO.Messages(map[string][]irest.IMessage{
		"test.ts": {
			TestDeepNotSafeMessage{},
		},
	})
	newDTO.AllowedMessages([]rest.AllowMessage{
		{
			Package: "dto_test",
			Name:    "TestDeepNotSafeMessage",
		},
	})
	var msg TestDeepNotSafeMessage
	if err := mapper.DeepCheckDTOSafeMessage(newDTO, typeopr.Ptr{}.New(&msg)); err != nil {
		if err.Error() != "dto_test.TestInnerNotSafeMessage message is unsafe" {
			t.Errorf("the error doesn't match the expected error. Unexpected error: %s", err.Error())
		}
	}
	newDTO.Messages(map[string][]irest.IMessage{
		"test.ts": {
			TestDeepNotSafeMessage{},
			TestInnerNotSafeMessage{},
		},
	})
	defer newDTO.Messages(messages)
	newDTO.AllowedMessages([]rest.AllowMessage{
		{
			Package: "dto_test",
			Name:    "TestDeepNotSafeMessage",
		},
		{
			Package: "dto_test",
			Name:    "TestInnerNotSafeMessage",
		},
	})
	defer newDTO.AllowedMessages(allowMessages)
	if err := mapper.DeepCheckDTOSafeMessage(newDTO, typeopr.Ptr{}.New(&msg)); err != nil {
		t.Error(err)
	}
}

func TestJsonToDTOMessage(t *testing.T) {
	json := map[string]any{"Id": 1, "Name": "name", "Ok": true}
	var out TestMessage
	if err := mapper.JsonToDTOMessage(json, newDTO, &out); err != nil {
		t.Error(err)
	}
}

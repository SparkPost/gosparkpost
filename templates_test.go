package gosparkpost_test

import (
	"fmt"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
)

func TestTemplateValidation(t *testing.T) {
	fromStruct := sp.From{"a@b.com", "A B"}
	f, err := sp.ParseFrom(fromStruct)
	if err != nil {
		t.Error(err)
		return
	}
	if fromStruct.Email != f.Email {
		t.Error(fmt.Errorf("expected email [%s] didn't match actual [%s]",
			fromStruct.Email, f.Email))
		return
	}
	if fromStruct.Name != f.Name {
		t.Error(fmt.Errorf("expected name [%s] didn't match actual [%s]",
			fromStruct.Name, f.Name))
		return
	}

	addrStruct := sp.Address{"a@b.com", "A B", "c@d.com"}
	f, err = sp.ParseFrom(addrStruct)
	if err != nil {
		t.Error(err)
		return
	}
	if addrStruct.Email != f.Email {
		t.Error(fmt.Errorf("expected email [%s] didn't match actual [%s]",
			addrStruct.Email, f.Email))
		return
	}
	if addrStruct.Name != f.Name {
		t.Error(fmt.Errorf("expected name [%s] didn't match actual [%s]",
			addrStruct.Name, f.Name))
		return
	}

	fromString := "a@b.com"
	f, err = sp.ParseFrom(fromString)
	if err != nil {
		t.Error(err)
		return
	}
	if fromString != f.Email {
		t.Error(fmt.Errorf("expected email [%s] didn't match actual [%s]",
			fromString, f.Email))
		return
	}
	if "" != f.Name {
		t.Error(fmt.Errorf("expected name to be blank"))
		return
	}
	fromString = ""
	_, err = sp.ParseFrom(fromString)
	if err == nil {
		t.Error(fmt.Errorf("Content.From should not be allowed!"))
	}

	fromMap1 := map[string]interface{}{
		"name":  "A B",
		"email": "a@b.com",
	}
	f, err = sp.ParseFrom(fromMap1)
	if err != nil {
		t.Error(err)
		return
	}
	// ParseFrom will bail if these aren't strings
	fromString, _ = fromMap1["email"].(string)
	if fromString != f.Email {
		t.Error(fmt.Errorf("expected email [%s] didn't match actual [%s]",
			fromString, f.Email))
		return
	}
	nameString, _ := fromMap1["name"].(string)
	if nameString != f.Name {
		t.Error(fmt.Errorf("expected name [%s] didn't match actual [%s]",
			nameString, f.Name))
		return
	}

	fromMap1["name"] = 1
	f, err = sp.ParseFrom(fromMap1)
	if err == nil {
		t.Error(fmt.Errorf("failed to detect non-string name"))
		return
	}

	fromMap2 := map[string]string{
		"name":  "A B",
		"email": "a@b.com",
	}
	f, err = sp.ParseFrom(fromMap2)
	if err != nil {
		t.Error(err)
		return
	}
	if fromMap2["email"] != f.Email {
		t.Error(fmt.Errorf("expected email [%s] didn't match actual [%s]",
			fromMap2["email"], f.Email))
		return
	}
	if fromMap2["name"] != f.Name {
		t.Error(fmt.Errorf("expected name [%s] didn't match actual [%s]",
			fromMap2["name"], f.Name))
		return
	}

	fromBytes := []byte("a@b.com")
	f, err = sp.ParseFrom(fromBytes)
	if err == nil {
		t.Error(fmt.Errorf("failed to detect unsupported type"))
		return
	}

}

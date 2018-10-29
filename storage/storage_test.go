package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func Test_ReadFileStorage(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		log.Printf("Up's tempfile testing -> %s", err.Error())
		t.Fail()
	}
	fstore, err := NewFileStorage(f.Name())
	if err != nil {
		log.Printf("Error on fileStorage -> %s", err.Error())
		t.Fail()
	}
	expected := Response{
		Rates: map[string]Rate{"operation_1": {
			Value:      5.00,
			LastUpdate: time.Now().UTC(),
		}},
	}
	bts, _ := json.Marshal(expected)
	if err := fstore.Write(bts); err != nil {
		log.Printf("Error wiritng to fileStorage -> %s", err.Error())
		t.Fail()
	}

	gotInterface, _ := fstore.Read(&Response{})
	got, ok := gotInterface.(*Response)
	if !ok {
		log.Printf("Wrong type ")
		t.Fail()
	}
	if got.Rates["operation_1"] != expected.Rates["operation_1"] {
		log.Printf("Error wiritng to fileStorage expcted -> \n %+v \n got -> \n %+v", expected, *got)
		t.Fail()
	}
}

package structhash

import (
	// "fmt"
	"testing"
)

type First struct {
	Bool   bool    `version:"1"`
	String string  `version:"2"`
	Int    int     `version:"1" lastversion:"1"`
	Struct *Second `version:"1"`
}

type Second struct {
	Map   map[string]string `version:"1"`
	Slice []int             `version:"1"`
}

func dataSetup() *First {
	tmpmap := make(map[string]string)
	tmpmap["foo"] = "bar"
	tmpmap["baz"] = "go"
	tmpslice := make([]int, 3)
	tmpslice[0] = 0
	tmpslice[1] = 1
	tmpslice[2] = 2
	return &First{
		Bool:   true,
		String: "test",
		Int:    123456789,
		Struct: &Second{
			Map:   tmpmap,
			Slice: tmpslice,
		},
	}
}

func TestHash(t *testing.T) {
	v1Hash := "v1_19a51ad2e72bc84e8c35dcdfe55e4d79"
	v2Hash := "v2_00a9b29c4d344aec8114c76b324757c0"

	data := dataSetup()
	v1, err := Hash(data, 1)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(v1)
	if v1 != v1Hash {
		t.Errorf("%s is not %s", v1, v1Hash)
	}
	v2, err := Hash(data, 2)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(v2)
	if v2 != v2Hash {
		t.Errorf("%s is not %s", v2, v2Hash)
	}
}

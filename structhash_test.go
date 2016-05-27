package structhash

import (
	"fmt"
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

type Tags1 struct {
	Int int    `hash:"-"`
	Str string `hash:"Foo, version(1) lastversion(2)"`
	Bar string `hash:",version(1)"`
}

type Tags2 struct {
	Foo string
	Bar string
}

type Tags3 struct {
	Bar string
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

	v1md5 := fmt.Sprintf("v1_%x", Md5(data, 1))
	if v1md5 != v1Hash {
		t.Errorf("%s is not %s", v1md5, v1Hash[3:])
	}
	v2md5 := fmt.Sprintf("v2_%x", Md5(data, 2))
	if v2md5 != v2Hash {
		t.Errorf("%s is not %s", v2md5, v2Hash[3:])
	}
}

func TestTags(t *testing.T) {
	t1 := Tags1{11, "foo", "bar"}
	t1x := Tags1{22, "foo", "bar"}
	t2 := Tags2{"foo", "bar"}
	t3 := Tags3{"bar"}

	t1_dump := string(Dump(t1, 1))
	t1x_dump := string(Dump(t1x, 1))
	if t1_dump != t1x_dump {
		t.Errorf("%s is not %s", t1_dump, t1x_dump)
	}

	t2_dump := string(Dump(t2, 1))
	if t1_dump != t2_dump {
		t.Errorf("%s is not %s", t1_dump, t2_dump)
	}

	t1v3_dump := string(Dump(t1, 3))
	t3v3_dump := string(Dump(t3, 3))
	if t1v3_dump != t3v3_dump {
		t.Errorf("%s is not %s", t1v3_dump, t3v3_dump)
	}
}

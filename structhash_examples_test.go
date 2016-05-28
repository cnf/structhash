package structhash

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

func ExampleHash() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
		Spouse *Person
	}
	bill := &Person{
		Name:   "Bill",
		Age:    24,
		Emails: []string{"bob@foo.org", "bob@bar.org"},
		Extra: map[string]string{
			"facebook": "Bob42",
		},
	}
	bob := &Person{
		Name:   "Bob",
		Age:    42,
		Emails: []string{"bob@foo.org", "bob@bar.org"},
		Extra: map[string]string{
			"facebook": "Bob42",
		},
		Spouse: bill,
	}

	hash, err := Hash(bob, 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", hash)
	// Output:
	// v1_d00068b9441e09d87689c7cb06a646a1
}

func ExampleHash_tags() {
	type Person struct {
		Ignored string            `hash:"-"`
		NewName string            `hash:"name:OldName version:1"`
		Age     int               `hash:"version:1"`
		Emails  []string          `hash:"version:1"`
		Extra   map[string]string `hash:"version:1 lastversion:2"`
		Spouse  *Person           `hash:"version:2"`
	}
	bill := &Person{
		NewName: "Bill",
		Age:     24,
		Emails:  []string{"bob@foo.org", "bob@bar.org"},
		Extra: map[string]string{
			"facebook": "Bob42",
		},
	}
	bob := &Person{
		NewName: "Bob",
		Age:     42,
		Emails:  []string{"bob@foo.org", "bob@bar.org"},
		Extra: map[string]string{
			"facebook": "Bob42",
		},
		Spouse: bill,
	}
	hashV1, err := Hash(bob, 1)
	if err != nil {
		panic(err)
	}
	hashV2, err := Hash(bob, 2)
	if err != nil {
		panic(err)
	}
	hashV3, err := Hash(bob, 3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", hashV1)
	fmt.Printf("%s\n", hashV2)
	fmt.Printf("%s\n", hashV3)
	// Output:
	// v1_a4500d206f830e75bb4b362705ee6240
	// v2_67caf9d9f9d2922ecc6f997bace6f06c
	// v3_a10f69ec95d652fc16f5f744a554e624
}

func ExampleDump() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
		Spouse *Person
	}
	bill := &Person{
		Name:   "Bill",
		Age:    24,
		Emails: []string{"bob@foo.org", "bob@bar.org"},
		Extra: map[string]string{
			"facebook": "Bob42",
		},
	}
	bob := &Person{
		Name:   "Bob",
		Age:    42,
		Emails: []string{"bob@foo.org", "bob@bar.org"},
		Extra: map[string]string{
			"facebook": "Bob42",
		},
		Spouse: bill,
	}

	fmt.Printf("md5:  %x\n", md5.Sum(Dump(bob, 1)))
	fmt.Printf("sha1: %x\n", sha1.Sum(Dump(bob, 1)))
	// Output:
	// md5:  d00068b9441e09d87689c7cb06a646a1
	// sha1: 24c19cd7a9fcfd4d4394e1f6a1874bd8751645e3
}

func ExampleVersion() {
	// A hash string gotten from Hash(). Returns the version as an int.
	i := Version("v1_55743877f3ffd5fc834e97bc43a6e7bd")
	fmt.Printf("%d", i)
	// Output:
	// 1
}

func ExampleVersion_errors() {
	// A hash string gotten from Hash(). Returns -1 on error.
	i := Version("va_55743877f3ffd5fc834e97bc43a6e7bd")
	fmt.Printf("%d", i)
	// Output:
	// -1
}

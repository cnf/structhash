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
		Ignored  string            `hash:"-"`
		FullName string            `hash:"Name, version(1)"`
		Age      int               `hash:",version(1)"`
		Emails   []string          `hash:",version(1)"`
		Extra    map[string]string `hash:",version(1) lastversion(2)"`
		Spouse   *Person           `hash:",version(2)"`
	}
	bill := &Person{
		FullName: "Bill",
		Age:      24,
		Emails:   []string{"bob@foo.org", "bob@bar.org"},
		Extra: map[string]string{
			"facebook": "Bob42",
		},
	}
	bob := &Person{
		FullName: "Bob",
		Age:      42,
		Emails:   []string{"bob@foo.org", "bob@bar.org"},
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
	// v1_461558d2570e10f79693e34ea309d1ad
	// v2_d00068b9441e09d87689c7cb06a646a1
	// v3_b5b651c6650939ef4d063d05caa5c778
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

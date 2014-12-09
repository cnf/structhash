package structhash

import (
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
	// v1_55743877f3ffd5fc834e97bc43a6e7bd
}

func ExampleHash_version() {
	type Person struct {
		Name   string            `version:"1"`
		Age    int               `version:"1"`
		Emails []string          `version:"1"`
		Extra  map[string]string `version:"1" lastversion:"2"`
		Spouse *Person           `version:"2"`
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

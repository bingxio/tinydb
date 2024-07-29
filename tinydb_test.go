package tinydb

import (
	"fmt"
	"os"
	"testing"
)

type user struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

func TestTinyDB(t *testing.T) {
	file, err := os.Create("test.json")
	if err != nil {
		t.Fatal(err)
	}

	var db Database

	db.Insert("user", user{Name: "kate", Age: 34})
	db.Insert("user", user{Name: "Gee", Age: 13})
	db.Insert("user", user{Name: "Oifd", Age: 54})
	db.Insert("user", user{Name: "Josin", Age: 22})

	result := db.Query("user", func(u any) bool {
		var user = u.(user)
		return user.Age >= 18
	})
	fmt.Printf("result: %v\n", result)

	modified := db.Delete("user", func(u any) bool {
		var user = u.(user)
		return user.Name == "Gee"
	})
	fmt.Printf("modified: %v\n", modified)

	modified = db.Update(
		"user",
		func(u any) bool {
			var user = u.(user)
			return user.Name == "Josin"
		},
		user{Name: "Updated"},
	)
	fmt.Printf("modified: %v\n", modified)

	result = db.Query("user", nil)
	fmt.Printf("result: %v\n", result)

	if err = db.WriteStorage(file); err != nil {
		t.Fatal(err)
	}
	file.Close()
}

func TestOpenDB(t *testing.T) {
	db, err := OpenDB("test.json")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("len(db.Table): %v\n", len(db.Table))

	result := db.Query("user", nil)
	for _, v := range result {
		fmt.Println(v)
	}
}

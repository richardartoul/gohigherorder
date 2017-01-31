package main

import "testing"

type User struct {
	Name string
}

//gohigherorder
type UserSlice []User // (["[]"], "User")

/*
UserPtrSlice is a slice of User pointers.
gohigherorder
*/
type UserPtrSlice []*User // (["[]", "*"], "User")

// Type without gohigherorder on newline is ignored
type UserSliceIgnore1 []User

func TestUserSliceFilter(t *testing.T) {
	userSlice := UserSlice{User{"Richie"}, User{"Bradfield"}}
	filteredSlice := userSlice.Filter(func(u User) bool {
		return u.Name == "Richie"
	})

	if len(filteredSlice) != 1 {
		t.Error("Filter receiver did not filter properly.")
	}
	if !(filteredSlice[0].Name == "Richie") {
		t.Error("Filter receiver did not filter properly.")
	}
}

func TestUserPtrSliceFilter(t *testing.T) {
	userSlice := UserPtrSlice{&User{"Richie"}, &User{"Bradfield"}}
	filteredSlice := userSlice.Filter(func(u *User) bool {
		return u.Name == "Richie"
	})

	if len(filteredSlice) != 1 {
		t.Error("Filter receiver did not filter properly.")
	}
	if !(filteredSlice[0].Name == "Richie") {
		t.Error("Filter receiver did not filter properly.")
	}
}

func TestUserSliceMap(t *testing.T) {
	userSlice := UserSlice{User{"Richie"}, User{"Bradfield"}}
	mappedSlice := userSlice.Map(func(u User) User {
		u.Name = u.Name + "1"
		return u
	})

	if len(mappedSlice) != 2 {
		t.Error("Filter receiver did not filter properly.")
	}
	if !(mappedSlice[0].Name == "Richie1") {
		t.Error("Filter receiver did not filter properly.")
	}
	if !(mappedSlice[1].Name == "Bradfield1") {
		t.Error("Filter receiver did not filter properly.")
	}
}

func TestUserPtrSliceMap(t *testing.T) {
	userSlice := UserPtrSlice{&User{"Richie"}, &User{"Bradfield"}}
	mappedSlice := userSlice.Map(func(u *User) *User {
		u.Name = u.Name + "1"
		return u
	})

	if len(mappedSlice) != 2 {
		t.Error("Filter receiver did not filter properly.")
	}
	if !(mappedSlice[0].Name == "Richie1") {
		t.Error("Filter receiver did not filter properly.")
	}
	if !(mappedSlice[1].Name == "Bradfield1") {
		t.Error("Filter receiver did not filter properly.")
	}
}

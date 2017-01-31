package main

//salice++
type User struct {
  name string
}

//slice++
type Users []User     // (["[]"], "User")
//slice++
type Users []*User    // (["[]", "*"], "User")
//slice++
type Users []**User   // (["[]", "*", "*"], "User")
//slice++
type Users [][]User   // (["[]", "[]"], "User")
//slice++
type Users [][]*User  // (["[]", "[]", "*"], "User")

var x int = 10

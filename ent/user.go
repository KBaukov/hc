package ent

type User struct {
	Login string
	Name  string
	Type  string
}

//func NewUser(l string, n string, t string) *User {

//	//	if l == nil {
//	//		panic("Login cannot be nil")
//	//	}

//	//	if n == nil {
//	//		panic("Name cannot be nil")
//	//	}

//	//	if t == nil {
//	//		panic("User type cannot be nil")
//	//	}

//	return &User{l, n, t}
//}

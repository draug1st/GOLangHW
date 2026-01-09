package main

import "fmt"

type User struct {
	Name string
}

func (u User) Greet() string {
	return fmt.Sprintf("Hello, %s!", u.Name)
}

type Admin struct {
	User
	banned     bool
	banMessage string
}

func (a *Admin) GetBanMessage() string {
	return a.banMessage
}

func (a *Admin) Ban(message string) {
	a.banned = true
	a.banMessage = fmt.Sprintf("user %s has been banned for %s", a.Name, message)
}

func NewAdmin(user User) *Admin {
	return &Admin{user, false, ""}
}

func (a *Admin) Greet() string {
	if a.banned {
		return a.banMessage
	}
	return a.User.Greet()
}

type Greeter interface { // 2 implementations
	Greet() string // 2 implementations
}

func main() {
	user1 := User{Name: "Alice"}
	admin := NewAdmin(user1)
	admin.Ban("spamming")
	fmt.Println(admin.Greet())

	user2 := User{Name: "Bob"}
	admin2 := NewAdmin(user2)
	fmt.Println(admin2.Greet())
}

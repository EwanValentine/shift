package main

import (
	"fmt"
	"log"

	shift "github.com/EwanValentine/shift/pkg/service"
)

// User -
type User struct {
	Name string
	Age  uint
}

// Test -
type Test struct {
	Date uint64
}

// ExampleService -
type ExampleService struct{}

// Greet - Example Greet:(*User, *Test):(string)
func (svc *ExampleService) Greet(user *User, test *Test) (string, error) {
	log.Println("Hello, this is test", user)
	log.Println("Received at:", test.Date)
	return fmt.Sprintf("Hello %s I am %d", user.Name, user.Age), nil
}

func main() {
	service := shift.NewService()
	service.Register(&ExampleService{}, "http://localhost:5002")
	service.RegisterType("*User", User{})
	service.RegisterType("*Test", Test{})
	service.Emit(shift.Event{
		Signature: "Greet:(*User, *Test):(string, error)",
		Body:      []byte(`{ "name": "Ewan", "age": 29 }`),
	})
	service.Run(":5002")
}

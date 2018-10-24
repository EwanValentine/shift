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

// Greet - Example Greet:(*User, *Test)
func (svc *ExampleService) Greet(user *User, test *Test) {
	log.Println("Hello, this is test", user)
	log.Println("Received at:", test.Date)
	fmt.Printf("Hello %s I am %d", user.Name, user.Age)
}

func main() {

	// Create a new service
	service := shift.NewService()

	// Register a struct against an address
	service.Register(&ExampleService{}, "http://localhost:5002")

	// Register types (will ideally phase out the need for this)
	service.RegisterType("*User", User{})
	service.RegisterType("*Test", Test{})

	// Trigger an event to anything listening with the given
	// method signature.
	service.Emit(shift.Event{
		Signature: "Greet:(*User, *Test)",
		Body:      []byte(`{ "name": "Ewan", "age": 29 }`),
	})

	// Run the service, as a http server
	service.Run(":5002")
}

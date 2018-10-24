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

// ExampleService -
type ExampleService struct{}

// Greet - Example Greet:(string):(string)
func (svc *ExampleService) Greet(user *User) (string, error) {
	log.Println("Hello, this is test", user)
	return fmt.Sprintf("Hello %s I am %d", user.Name, user.Age), nil
}

func main() {
	service := shift.NewService()
	service.Register(&ExampleService{}, "http://localhost:5002")
	service.Emit(shift.Event{
		Signature: "Greet:(*User):(string, error)",
		Body:      []byte(`{ "name": "Ewan", "age": 29 }`),
	})
	service.Run(":5002")
}

![Shift Logo](./logo.png)

# Shift
De-coupled, event-driven microservice architecture. Based around Golangs interface design.

## Note
This is a purely theoretical currently, although the examples function. It currently doesn't have service discovery, you're currently hard-coding network locations, and those locations aren't currently shared across multiple instances. 

If this design is validated, then I'll start to create integrations with services such as Consul, ZooKeeper etc. Oh and different transport methods, such as gRPC, NATS, etc.

## Design
This project is inspired by Golangs interface system. I was thinking of ways in which you can trigger events without having to worry about registering handlers, thinking up event names etc. I enjoy the experience of using gRPC and similar tools, where talking to another service, just feels like executing code.

I also enjoy tools such as GraphQL, whereby you use a rich domain specific language to perform tasks and operations. However, I felt that gRPC tended to be good for coupled services. In other words, you make a request to another service, and you get a response back. 

## How it works
With Shift, you emit events, as you would with a typical event-driven system. However, instead of thinking in terms of event names, you think in terms of code. Your event name becomes:

```
Greet:(*User, *Test)
```

When you register a service, for example: 

```
shift.Register(&MyService{}, "http://localhost:5000")
```

`&MyService{}` is scanned, its methods are analysed using reflection. Each method on that struct is then registered against its signature, under that host address.

This means that if you emitted an event with this signature: `Greet:(*User)`. Shift will use service discovery to find which host address/network location `Greet:(*User)` exists in. There maybe be one or more of these functions present. Thus Shift iterates through each. Which is similar to how interfaces work in Go. A single struct can satisfy different interfaces, by the signature of that method.

This allows us to emit de-coupled events across services, using the behaviour of your code and the signature of your methods and types to match what should be listening to that event. 

Possible variations on this, currently we're looking at the method, the argument signature, and the return types. We could also take into account the struct name if we wanted to be more focussed, or even talk directly to another service. Or we could be more loose and purely base our events on argument signature and return types. For example, emit this event to anything which takes a `*User` argument and returns a `*CreatedUser` object.

## Code Example
```go
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
```

## Run Example
```
$ cd example && vgo run test.go
```

## Todo
- Integrate service discovery.
- Integrate different transports.
- Tidy up code.
- Add tests, I know, I've already written the code.
- Benchmark and optimise.

## Feedback welcome
[Email Me](ewan.valentine89@gmail.com)

[Tweet Me](https://twitter.com/Ewan_Valentine)

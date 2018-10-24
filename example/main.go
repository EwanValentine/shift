package main

import "log"

type MyService struct{}

func (s *MyService) PrintUser(user *User) {
	log.Println(user)
}

func main() {
	srv := shift.NewService()
	srv.Expose(MyService{})
	srv.Emit(shift.Event{
		Signature: "PrintOrder:(order *Order)",
		Body:      []byte(`{ "order": "abc123" }`),
	})
	srv.Run(":8080")
}

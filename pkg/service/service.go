package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"

	"github.com/EwanValentine/shift/pkg/parser"
	"github.com/EwanValentine/shift/pkg/registry"
)

// Svc is a reference to another service
// @todo - this is poorly named
type Svc struct {
	Signature string
	Addr      string
}

// Parser to parse the DSL
type Parser interface {
	Parse(method, signature string) string
	Unmarshal(signature string) parser.Signature
}

// Registry for types/strings
type Registry interface {
	Register(name string, t interface{})
	MakeInstance(name string) interface{}
}

// Service is the base struct
type Service struct {
	mu       sync.Mutex
	Services []Svc
	Service  interface{}
	parser   Parser
	registry Registry
}

// NewService returns a Service instance
func NewService() *Service {
	p := parser.NewParser()
	r := registry.NewRegistry()
	return &Service{
		parser:   p,
		registry: r,
	}
}

// Register a service
func (s *Service) Register(svc interface{}, addr string) error {
	t := reflect.TypeOf(svc)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		svcSig := s.parser.Parse(method.Name, method.Type.String())
		s.mu.Lock()
		updated := append(s.Services, Svc{
			Signature: svcSig,
			Addr:      addr,
		})
		s.Services = updated
		s.Service = svc
		s.mu.Unlock()
	}
	return nil
}

func (s *Service) callService(svc Svc, event Event) (string, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	body := bytes.NewReader([]byte(data))
	res, err := http.Post(svc.Addr, "application/json", body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

// Call a service by its signature
func (s *Service) Call(event Event) []string {
	var wg sync.WaitGroup
	var results []string
	resChan := make(chan string)
	for _, svc := range s.Services {
		if svc.Signature == event.Signature {
			go func() {
				wg.Add(1)
				defer wg.Done()
				res, err := s.callService(svc, event)
				// @todo - handle error properly here...
				log.Println(err)
				resChan <- res
			}()
		}
	}
	wg.Wait()
	for res := range resChan {
		results = append(results, res)
	}
	return results
}

// Emit an event
func (s *Service) Emit(event Event) error {
	var wg sync.WaitGroup
	for _, svc := range s.Services {
		if svc.Signature == event.Signature {
			go func() {
				wg.Add(1)
				defer wg.Done()
				s.callService(svc, event)
			}()
		}
	}
	wg.Wait()
	return nil
}

// Invoke a function by name
func (s *Service) invoke(any interface{}, name string, body []byte, signature parser.Signature) {
	args := signature.Args
	log.Println(args[0])
	data := s.registry.MakeInstance(signature.Args[0])
	json.Unmarshal(body, data)
	inputs := make([]reflect.Value, len(args))
	for i := range signature.Args {
		log.Println(reflect.TypeOf(data))
		log.Println(reflect.ValueOf(data))
		inputs[i] = reflect.ValueOf(data)
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}

func (s *Service) handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event Event
	if err := decoder.Decode(&event); err != nil {
		log.Println(err)
	}
	signature := s.parser.Unmarshal(event.Signature)
	s.invoke(s.Service, signature.Method, event.Body, signature)
	w.Write([]byte("Ok"))
}

// Run the webserver
func (s *Service) Run(port string) error {
	http.HandleFunc("/", s.handler)
	return http.ListenAndServe(port, nil)
}

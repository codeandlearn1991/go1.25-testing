package main

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"encoding/json/jsontext"
	"encoding/json/v2"
)

func main() {
	emojiMarshaler := json.MarshalToFunc(func(enc *jsontext.Encoder, val string) error {
		if val == "yes" || val == "true" {
			return enc.WriteToken(jsontext.String("🙂"))
		}
		if val == "no" || val == "false" {
			return enc.WriteToken(jsontext.String("🤨"))
		}
		// SkipFunc is a special type of error that tells Go to skip
		// the current marshaler and move on to the next one. In our case,
		// the next one will be the default marshaler for strings.
		return json.SkipFunc
	})

	boolMarshaler := json.MarshalToFunc(
		func(enc *jsontext.Encoder, val bool) error {
			if val {
				return enc.WriteToken(jsontext.String("Oh yeah!"))
			}
			return enc.WriteToken(jsontext.String("Oh no!"))
		},
	)

	marshalers := json.JoinMarshalers(boolMarshaler, emojiMarshaler)

	vals := []any{true, "off", "no", "hello"}

	data, _ := json.Marshal(vals, json.WithMarshalers(marshalers))
	fmt.Println(string(data))

	type Data struct {
		Items []string       `json:"items"`
		Meta  map[string]any `json:"meta"`
	}

	d := Data{}

	b1, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Default v2 (nil -> {} or []):", string(b1)) // Output: {"items":[],"meta":{}}

	b2, err := json.Marshal(d, json.FormatNilSliceAsNull(true), json.FormatNilMapAsNull(true))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("With options (nil -> null):", string(b2)) // Output: {"items":null,"meta":null}

	type Person struct {
		FirstName string `json:"firstName"`
		Age       int    `json:"age,case:ignore"`
	}

	// Case Sensitivity in Unmarshaling
	jsonInput := `{"firstname": "Bob", "AGE": 40}`
	var p1 Person

	// Default v2 (case-sensitive) - will not unmarshal "firstname" or "AGE"
	err = json.Unmarshal([]byte(jsonInput), &p1)
	if err != nil {
		// v2 will produce an error about unknown fields or just not populate
		fmt.Println("Default v2 Unmarshal Error:", err)
	}
	fmt.Printf("Default v2 Unmarshaled: %+v\n", p1) // Output: {FirstName: Age:0} (or error)

	var p2 Person
	// With option to match case-insensitively (like v1)
	err = json.Unmarshal([]byte(jsonInput), &p2, json.MatchCaseInsensitiveNames(true))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("With options Unmarshaled: %+v\n", p2) // Output: {FirstName:Bob Age:40}

	// Binary data representation
	type BinaryData struct {
		Data []byte `json:"data"`
		Raw  []byte `json:"raw,format:array"` // Force v1 behavior
	}

	bd := BinaryData{Data: []byte("Hello"), Raw: []byte("World")}

	b, err := json.Marshal(bd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Binary Data:", string(b))
	// Expected Output (v2): {"data":"SGVsbG8=","raw":[87,111,114,108,100]}

	// Time formatting

	type Event struct {
		Name      string    `json:"name"`
		Timestamp time.Time `json:"timestamp,format:RFC3339"`
		UnixEpoch time.Time `json:"unix_ts,format:unix"`
	}

	event := Event{
		Name:      "Meeting",
		Timestamp: time.Now(),
		UnixEpoch: time.Now(),
	}

	b, err = json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Event JSON:", string(b))
	// Example output: {"name":"Meeting","timestamp":"2025-07-16T12:45:27Z","unix_ts":1721157927}

	var wg sync.WaitGroup

	wg.Go(func() {
		// wg.Done()
		fmt.Println("go is awesome")
	})
	wg.Go(func() {
		fmt.Println("go is awesome")
	})
	wg.Go(func() {
		fmt.Println("go is awesome")
	})

	wg.Wait()

	// mux := http.NewServeMux()
	// mux.HandleFunc("GET /get", func(w http.ResponseWriter, req *http.Request) {
	// 	io.WriteString(w, "ok\n")
	// })
	// mux.HandleFunc("POST /post", func(w http.ResponseWriter, req *http.Request) {
	// 	io.WriteString(w, "ok\n")
	// })

	// // Configure protection against CSRF attacks.
	// antiCSRF := http.NewCrossOriginProtection()
	// antiCSRF.AddTrustedOrigin("https://example.com")
	// antiCSRF.AddTrustedOrigin("https://*.example.com")

	// // Add CSRF protection to all handlers.
	// srv := http.Server{
	// 	Addr:    ":8080",
	// 	Handler: antiCSRF.Handler(mux),
	// }
	// log.Fatal(srv.ListenAndServe())

	var myInterface interface{} = "Hello, Go 1.25!"

	// Old way (allocates a new string for the assertion)
	val := reflect.ValueOf(myInterface)
	strVal, ok := val.Interface().(string)
	fmt.Println("Old way:", strVal)

	// New way with reflect.TypeAssert (avoids intermediate allocation)
	newStrVal, ok := reflect.TypeAssert[string](val)
	if !ok {
		fmt.Println("Type assertion failed!")
	} else {
		fmt.Println("New way (TypeAssert):", newStrVal)
	}

	var myIntf interface{} = 123
	newIntVal, ok := reflect.TypeAssert[int](reflect.ValueOf(myIntf))
	if !ok {
		fmt.Println("Integer type assertion failed!")
	} else {
		fmt.Println("New way (TypeAssert) for int:", newIntVal)
	}
}

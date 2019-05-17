package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// Import book_store generated code
	// Note that Go requires fully qualified paths from the module root. In this case, the module name is defined as
	// 'web' in go.mod. Convention is that this is the same as the directory name, but this is not required, the source
	// of truth is go.mod.
	"web/book_store"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	defaultName = "world"
)

func getBookFromBackend(backendEndAddress string) *book_store.Book {
	// Set up a connection to the server.
	conn, err := grpc.Dial(backendEndAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := book_store.NewBookStoreClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	book, err := c.GetBook(ctx, &book_store.BookReference{Title: "foo"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Got Book from Backend: %s", book.Title)
	return book
}

// https://stackoverflow.com/questions/40326540/how-to-assign-default-value-if-env-var-is-empty
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func main() {
	version := getEnv("SERVICE_VERSION", "1.0.0")
	port := getEnv("SERVICE_PORT", "1234")
	backendAddress := getEnv("SERVICE_BACKEND", "localhost:50051")
	hostname := getEnv("HOSTNAME", "Hostname undefined")

	log.Println("Hello!")
	log.Println("Version:", version)
	log.Println("Listening port:", port)
	log.Println("Backend:", backendAddress)

	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		log.Println("[GET] /hello")
		book := getBookFromBackend(backendAddress)
		fmt.Fprintf(w, "Book: %s\n", book.Title)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("[GET] /")
		fmt.Fprintf(w, "This is the index page! (Version: %s, Host: %s)\n", version, hostname)
	})
	http.ListenAndServe(":" + port, nil)

}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// Import book_store generated code
	// Note that Go requires fully qualified paths from the module root. In this case, the module name is defined as
	// 'web' in go.mod. Convention is that this is the same as the directory name, but this is not required, the source
	// of truth is go.mod.
	pb "web/book_store"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	defaultName = "world"
)

func getBookFromBackend() *pb.Book {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBookStoreClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	book, err := c.GetBook(ctx, &pb.BookReference{Title: "foo"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Got Book from Backend: %s", book.Title)
	return book
}

func main() {
	log.Println("Hello!")
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		book := getBookFromBackend()
		fmt.Fprintf(w, "Book: %s\n", book.Title)
	})
	http.ListenAndServe(":1234", nil)

}

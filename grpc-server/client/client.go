package main

import (
	context "context"
	"log"
	"time"

	pb "github.com/chancegraff/grpc2/gravatar"
	grpc "google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGravatarServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Generate(ctx, &pb.GravatarRequest{Email: "name", Size: 10})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.Url)
}

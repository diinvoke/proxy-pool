package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	rpc "github.com/mingcheng/proxypool/protobuf"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := rpc.NewProxyPoolClient(conn)

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	if proxy, err := client.Random(ctx, &empty.Empty{}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(proxy)
	}

	if proxies, err := client.All(context.TODO(), &empty.Empty{}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(proxies)
	}
}

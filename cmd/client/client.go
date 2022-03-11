package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/yur1ss/fc-grpc.git/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Yuri",
		Email: "y@y.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Yuri",
		Email: "y@y.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "Y1",
			Name:  "Yuri 1",
			Email: "y1@email.com",
		},
		&pb.User{
			Id:    "Y2",
			Name:  "Yuri 2",
			Email: "y2@email.com",
		},
		&pb.User{
			Id:    "Y3",
			Name:  "Yuri 3",
			Email: "y3@email.com",
		},
		&pb.User{
			Id:    "Y4",
			Name:  "Yuri 4",
			Email: "y4@email.com",
		},
		&pb.User{
			Id:    "Y5",
			Name:  "Yuri 5",
			Email: "y5@email.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "Y1",
			Name:  "Yuri 1",
			Email: "y1@email.com",
		},
		&pb.User{
			Id:    "Y2",
			Name:  "Yuri 2",
			Email: "y2@email.com",
		},
		&pb.User{
			Id:    "Y3",
			Name:  "Yuri 3",
			Email: "y3@email.com",
		},
		&pb.User{
			Id:    "Y4",
			Name:  "Yuri 4",
			Email: "y4@email.com",
		},
		&pb.User{
			Id:    "Y5",
			Name:  "Yuri 5",
			Email: "y5@email.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.GetName())
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}
			fmt.Printf("Recebendo user %v com status: %v", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}

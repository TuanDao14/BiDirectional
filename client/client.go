package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/tuanda/BiDirectional/BiDirectionalpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalf(" err while dial %v", err)
	}
	defer cc.Close()

	client := BiDirectionalpb.NewCalculatorServiceClient(cc)

	log.Printf("service client %f", client)
	callSum(client)
	callPND(client)
	callAverage(client)
	callFindMax(client)
}

func callSum(c BiDirectionalpb.CalculatorServiceClient) {
	log.Println("calling sum api")
	resp, err := c.Sum(context.Background(), &BiDirectionalpb.SumRequest{
		Num1: 7,
		Num2: 6,
	})

	if err != nil {
		log.Fatalf("call sum api err %v", err)
	}

	log.Printf("sum api response %v\n", resp.GetResult())
}

func callPND(c BiDirectionalpb.CalculatorServiceClient) {
	log.Println("calling pnd api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &BiDirectionalpb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("callPND err %v", err)
	}

	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Println("server finish streaming")
			return
		}

		if recvErr != nil {
			log.Fatalf("callPND recvErr %v", recvErr)
		}

		log.Printf("prime number %v", resp.GetResult())
	}
}

func callAverage(c BiDirectionalpb.CalculatorServiceClient) {
	log.Println("calling average api")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("call average err %v", err)
	}

	listReq := []BiDirectionalpb.AverageRequest{
		BiDirectionalpb.AverageRequest{
			Num: 5,
		},
		BiDirectionalpb.AverageRequest{
			Num: 10,
		},
		BiDirectionalpb.AverageRequest{
			Num: 12,
		},
		BiDirectionalpb.AverageRequest{
			Num: 3,
		},
		BiDirectionalpb.AverageRequest{
			Num: 4.2,
		},
	}

	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("send average request err %v", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("receive average response err %v", err)
	}

	log.Printf("average response %+v", resp)
}

func callFindMax(c BiDirectionalpb.CalculatorServiceClient) {
	log.Println("calling find max ...")
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("call find max err %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		//gui nhieu request
		listReq := []BiDirectionalpb.FindMaxRequest{
			BiDirectionalpb.FindMaxRequest{
				Num: 5,
			},
			BiDirectionalpb.FindMaxRequest{
				Num: 10,
			},
			BiDirectionalpb.FindMaxRequest{
				Num: 12,
			},
			BiDirectionalpb.FindMaxRequest{
				Num: 3,
			},
			BiDirectionalpb.FindMaxRequest{
				Num: 4,
			},
		}
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("send find max request err %v", err)
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("ending find max api ...")
				break
			}
			if err != nil {
				log.Fatalf("recv find max err %v", err)
				break
			}

			log.Printf("max: %v\n", resp.GetMax())
		}
		close(waitc)
	}()

	<-waitc
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/tuanda/BiDirectional/BiDirectionalpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *BiDirectionalpb.SumRequest) (*BiDirectionalpb.SumResponse, error) {
	log.Println("sum called...")
	resp := &BiDirectionalpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func (*server) PrimeNumberDecomposition(req *BiDirectionalpb.PNDRequest,
	stream BiDirectionalpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Println("PrimeNumberDecomposition called...")
	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			//send to client
			stream.Send(&BiDirectionalpb.PNDResponse{
				Result: k,
			})
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
			log.Printf("k increase to %v", k)
		}
	}
	return nil
}

func (*server) Average(stream BiDirectionalpb.CalculatorService_AverageServer) error {
	log.Println("Average called..")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//tinh trung binh va return cho client
			resp := &BiDirectionalpb.AverageResponse{
				Result: total / float32(count),
			}

			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatalf("err while Recv Average %v", err)
			return err
		}
		log.Printf("receive req %v", req)
		total += req.GetNum()
		count++
	}
}

func (*server) FindMax(stream BiDirectionalpb.CalculatorService_FindMaxServer) error {
	log.Println("Find max called ....")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF...")
			return nil
		}
		if err != nil {
			log.Fatalf("err while Recv FindMax %v", err)
			return err
		}

		num := req.GetNum()
		log.Printf("recv num %v\n", num)
		if num > max {
			max = num
		}
		err = stream.Send(&BiDirectionalpb.FindMaxResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("send max err %v", err)
			return err
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer()

	BiDirectionalpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("calculator is running...")
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("err while serve %v", err)
	}
}

package service

import (
	"context"
	"math/rand"
	"time"

	v1 "streaming-repro/api/helloworld/v1"
	"streaming-repro/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}

	flip := rand.Intn(2)
	if flip == 0 {
		panic("Oh no!")
	}

	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

func (s *GreeterService) KeepSayingHello(in *v1.HelloRequest, conn v1.Greeter_KeepSayingHelloServer) error {
	g, err := s.uc.CreateGreeter(conn.Context(), &biz.Greeter{Hello: in.Name})
	if err != nil {
		return err
	}

	for err == nil {
		err = conn.SendMsg(&v1.HelloReply{
			Message: "Hello " + g.Hello,
		})
		time.Sleep(1 * time.Second)
	}

	panic(err)
}

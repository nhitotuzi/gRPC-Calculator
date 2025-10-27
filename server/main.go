package main

import (
	"context"
	"log"
	"net"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection" // grpcurl
	"google.golang.org/grpc/status"
	pb "calculator/protobuf/operations" // Protobuf compilado
)

// server implementa a interface CalculatorServer gerada.
type server struct {
	pb.UnimplementedCalculatorServer // Compatibilidade com mudanças
}

func (*server) Add(ctx context.Context, req *pb.OperationRequest) (*pb.OperationResponse, error) {
	result := req.GetA() + req.GetB()
	return &pb.OperationResponse{Result: result}, nil
}

func (*server) Sub(ctx context.Context, req *pb.OperationRequest) (*pb.OperationResponse, error) {
	result := req.GetA() - req.GetB()
	return &pb.OperationResponse{Result: result}, nil
}

func (*server) Mul(ctx context.Context, req *pb.OperationRequest) (*pb.OperationResponse, error) {
	result := req.GetA() * req.GetB()
	return &pb.OperationResponse{Result: result}, nil
}

func (*server) Div(ctx context.Context, req *pb.OperationRequest) (*pb.OperationResponse, error) {
	b := req.GetB()
	if b == 0 {
		// Erro gRPC padrão para argumentos inválidos.
		return nil, status.Errorf(codes.InvalidArgument, "Divisor não pode ser zero")
	}
	result := req.GetA() / b
	return &pb.OperationResponse{Result: result}, nil
}

// Registra informações sobre cada chamada RPC.
func loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	log.Printf("--> Chamada recebida: %s | %v", info.FullMethod, req)

	// Repassa a chamada para o handler real (a implementação do método)
	resp, err := handler(ctx, req)

	duration := time.Since(start)
	statusCode := codes.OK
	if err != nil {
		if s, ok := status.FromError(err); ok {
			statusCode = s.Code()
		}
	}

	log.Printf("<-- Chamada finalizada: %s | %v | Resultado: %v | Duração: %s | Status: %s",
		info.FullMethod,
		req,
		resp,
		duration.Round(time.Microsecond),
		statusCode,
	)

	return resp, err
}

func main() {
	// Listener na porta 8081
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Falha ao escutar na porta: %v", err)
	}

	// Nova instância do servidor com interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	// Registro do serviço no servidor gRPC
	pb.RegisterCalculatorServer(s, &server{})

	// Ativando server reflection
	reflection.Register(s)

	log.Println("Servidor gRPC iniciado na porta :8081")
	if err := s.Serve(lis);
	err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
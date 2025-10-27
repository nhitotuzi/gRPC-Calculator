package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "calculator/protobuf/operations"
)

var (
	addr = flag.String("addr", "localhost:8081", "endereço do servidor")
)

func main() {
	flag.Parse()
	// Configurar a conexão com o servidor.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalculatorClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60)*time.Minute)
	defer cancel()

	// Interface para realizar as operações.
	for {
		var op string
		fmt.Print("Digite a operação (add, sub, mul, div) ou 'sair' para encerrar: ")
		_, err := fmt.Scanln(&op)
		if err != nil {
			log.Printf("Erro ao ler a operação: %v", err)
			continue
		}
		if op == "sair" {
			break
		}

		switch op {
		case "add", "sub", "mul", "div":
			// Válido
		default:
			log.Printf("Operação inválida: %s", op)
			continue
		}

		// Resposta da operação.
		var r *pb.OperationResponse

		var a, b float64

		fmt.Print("Digite o primeiro operando: ")
		fmt.Scanln(&a)
		fmt.Print("Digite o segundo operando: ")
		fmt.Scanln(&b)

		// Criar a requisição.
		or := &pb.OperationRequest{
			A: a,
			B: b,
		}

		switch op {
		case "add":
			r, err = c.Add(ctx, or)
		case "sub":
			r, err = c.Sub(ctx, or)
		case "mul":
			r, err = c.Mul(ctx, or)
		case "div":
			r, err = c.Div(ctx, or)
		}
		if err != nil {
			log.Printf("Erro ao realizar a operação: %v", err)
			continue
		}
		fmt.Printf("Resultado da operação: %.2f\n", r.GetResult())
	}

}

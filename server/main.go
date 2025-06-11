package main

import (
	"os"
	"time"
)

func main() {
	ReadBlocks()
	check := CheckBlocks()
	if !check {
		// Finaliza o sistema imediatamente
		println("Existe inconsistência de dados na blockchain lida.")
		os.Exit(1)
	}

	porta := os.Getenv("PORT") //Pega a porta do docker-compose
	startingREST(porta)

	//cria-se uma goroutine para sicronizar sempre as blockchains a cada 15 segundos
	go func() {
		for {
			time.Sleep(15 * time.Second)
			syncBlockchain(porta)
		}
	}()

	select {} // manter o programa rodando
}

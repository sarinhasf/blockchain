package main

import (
	"os"
	"time"
)

func main() {
	ReadBlocks()
	porta := os.Getenv("PORT") //Pega a porta do docker-compose
	startingREST(porta)

	//cria-se uma goroutine para sicronizar sempre as blockchains a cada 5 segundos
	go func() {
		for {
			time.Sleep(10 * time.Second)
			println("\nSincronizando blockchain com outros clientes...")
			syncBlockchain(porta)
		}
	}()

	BuildingMenu()
	select {} // manter o programa rodando
}

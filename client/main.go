package main

import (
	"os"
)

func main() {
	StartNode()
}

func StartNode() {
	ReadBlocks()
	porta := os.Getenv("PORT") //Pega a porta do docker-compose
	startingREST(porta)

	//cria-se uma goroutine para sicronizar sempre as blockchains a cada 5 segundos
	//go func() {
	syncBlockchain(porta)
	//time.Sleep(5 * time.Second)
	//}()

	//time.Sleep(3 * time.Second)
	BuildingMenu()
	select {} // manter o programa rodando
}

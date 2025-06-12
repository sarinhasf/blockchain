package main

import (
	"os"
	"time"
)

func main() {
	porta := os.Getenv("PORT") //Pega a porta do docker-compose

	ReadBlocks() //Lê os blocos do arquivo JSON

	if !isChainValid(blockchain.Blocos) {
		println("-----------------------------------------------------")
		println("ALERTA: Blockchain local detectada como inválida.")
		println("Tentando correção automática a partir da rede...")
		println("-----------------------------------------------------")

		syncBlockchain(porta) // Tenta sincronizar a blockchain com os outros nós

		// Verifica novamente após a tentativa de correção
		if !isChainValid(blockchain.Blocos) {
			println("*****************************************************")
			println("FALHA CRÍTICA: Não foi possível corrigir a blockchain.")
			println("*****************************************************")
			os.Exit(1) // Finaliza se a correção falhou
		}

		println("-----------------------------------------------------")
		println("SUCESSO: Blockchain corrigida a partir da rede!")
		println("-----------------------------------------------------")
	} else {
		println("Blockchain local verificada com sucesso.")
	}

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

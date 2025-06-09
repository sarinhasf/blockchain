package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func BuildingTitle() {
	fmt.Println("      ____             ____________________      ")
	fmt.Println("   __/  |_\\_          |                    |     ")
	fmt.Println("  |  _     _``-.      | S I S T E M A  D E |     ")
	fmt.Println("  '-(_)---(_)--'      |    R E C A R G A   |     ")
	fmt.Println("                      |____________________|     ")
	fmt.Println("")
	fmt.Println("")
}

func ConvertToNum(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Erro ao converter:", err)
	}
	return num
}

func GeneratingPlate() string {
	const letras = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	placa := make([]byte, 7)
	for i := 0; i < 3; i++ {
		placa[i] = letras[rand.Intn(len(letras))]
	}
	for i := 3; i < 7; i++ {
		placa[i] = byte(rand.Intn(10) + '0')
	}

	return string(placa)
}

/*
BuildingMenu: construindo Menu para exibir para os clientes.
*/
func BuildingMenu() {
	placa := GeneratingPlate() //gerando placa automaticamente
	fmt.Printf("\nIniciando sistema e gerando placa do veículo...\n")
	fmt.Printf("Cliente com a placa [%s] entrou no sistema!\n\n", placa)

	BuildingTitle()
	ReadPoints()

	pontos := dataPoints.PontosDeRecarga

	fmt.Println("<<SELECIONE O PONTO DE RECARGA QUE DESEJA RESERVAR>>")
	for n, ponto := range pontos {
		fmt.Printf("[%d]: %s\n", n+1, ponto.Nome)
	}

	leitor := bufio.NewReader(os.Stdin)
	var opcao string
	for {
		fmt.Print("Escolha uma opção (1 a 9): ")
		opcao, _ = leitor.ReadString('\n')
		opcao = strings.TrimSpace(opcao)

		if opcao >= "1" && opcao <= "9" && len(opcao) == 1 {
			break
		}

		fmt.Println("Opção inválida! Por favor, digite um número de 1 a 9.")
	}
	numChoose := ConvertToNum(opcao)
	posicao := numChoose - 1
	fmt.Println("\nVocê escolheu o ponto:", pontos[posicao].Nome)
	conteudo := "[placa: " + placa + ", ponto: " + pontos[posicao].Nome + "]"
	addNewBlock("reserva", conteudo)
	CarChargingSimulator(placa, pontos[posicao])
}

/*
ValueRandom: gera valor de pagamento aleatório.
*/
func ValueRandom() string {
	rand.Seed(time.Now().UnixNano())
	valor := 10 + rand.Float64()*(500-10)
	return fmt.Sprintf("%.2f", valor)
}

/*
CarChargingSimulator: função para simular carro indo até o ponto e carregando, e após isso pagando.
*/
func CarChargingSimulator(placa string, ponto Point) {
	//simulando carregamento
	fmt.Printf("\nCarro [%s] se deslocando até o %s...", placa, ponto.Nome)
	time.Sleep(5 * time.Second)
	fmt.Printf("\nCarro [%s] chegando ao %s...", placa, ponto.Nome)
	time.Sleep(5 * time.Second)
	fmt.Printf("\nCarro [%s] carregando no %s...\n", placa, ponto.Nome)
	conteudoRecarga := "[placa: " + placa + ", ponto: " + ponto.Nome + "]"
	addNewBlock("recarga", conteudoRecarga) //criando novo bloco

	//simulando pagamento
	valorRecarga := ValueRandom() //gerando valor aleatório de pagamento
	fmt.Printf("\nRecarga do carro [%s] no %s finalizada! O valor total foi: R$ %s.", placa, ponto.Nome, valorRecarga)
	fmt.Printf("\nCarro [%s] efetuando pagamento...\n", placa)
	conteudoPagamento := "[placa: " + placa + ", ponto: " + ponto.Nome + ", valor: " + valorRecarga + "]"
	addNewBlock("pagamento", conteudoPagamento) //criando novo bloco

	//serverLocal := fmt.Sprintf("http://%s:%s", ipLocal, portaLocal)
	//getBlockchainFrom(serverLocal)
}

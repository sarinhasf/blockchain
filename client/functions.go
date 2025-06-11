package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// AQUI SETAMOS OS SERVIDORES DE TODAS MAQUINAS
var servidores []string = []string{
	"http://server1:8091",
	"http://server2:8092",
	"http://server3:8093",
}

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

func VerifyCompany(choose int) int {
	empresas := dataCompanies.Empresas
	//verificando qual empresa mandar de acordo com o ponto escolhido
	//fmt.Println("Entrando no for de verifyCompany...")
	//fmt.Printf("Escolha: %d\n", choose)
	for _, empresa := range empresas {
		for _, p := range empresa.Pontos {
			//fmt.Printf("Empresa %s, Ponto %d\n", empresa.Nome, p)
			if p == choose {
				//fmt.Printf("O id da empresa é %d\n", empresa.ID)
				return empresa.ID
			}
		}
	}
	return 0
}

func getUrlById(id int) string {
	var url string
	switch id {
	case 1:
		url = servidores[0]
	case 2:
		url = servidores[1]
	case 3:
		url = servidores[2]

	}
	return url
}

func sendMessageForServer(mensagem string, url string) {
	fmt.Printf("Mensagem enviada ao servidor %s: %s\n", url, mensagem)
	msg := Message{Content: mensagem}
	urlSend := url + "/mensagem"

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Erro ao converter mensagem:", err)
		return
	}

	resp, err := http.Post(urlSend, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Erro ao enviar requisição:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Resposta do servidor: ", resp.Status)
}

/*
BuildingMenu: construindo Menu para exibir para os clientes.
*/
func BuildingMenu() {
	placa := GeneratingPlate() //gerando placa automaticamente
	fmt.Printf("\nIniciando sistema e gerando placa do veículo...\n")
	fmt.Printf("Cliente com a placa [%s] entrou no sistema!\n\n", placa)

	BuildingTitle()

	//lendo dados dos arquivos
	ReadPoints()
	ReadCompanies()

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

	//Estrutura da mensagem: tipo de processo (reserva, pagamento...), placa, nome do ponto
	mensagem := "reserva," + placa + "," + pontos[posicao].Nome

	idEmpresa := VerifyCompany(numChoose)
	var url string
	if idEmpresa == 1 || idEmpresa == 2 || idEmpresa == 3 { //verificando se o id retornado é 1, 2 ou 3
		url = getUrlById(idEmpresa)
		sendMessageForServer(mensagem, url)
	} else {
		fmt.Println("Id da empresa não encontrado!")
		return
	}

	//simula carro carregando e pagando
	CarChargingSimulator(placa, pontos[posicao], url)
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
func CarChargingSimulator(placa string, ponto Point, url string) {

	//simulando carregamento
	fmt.Printf("\nCarro [%s] se deslocando até o %s...", placa, ponto.Nome)
	time.Sleep(5 * time.Second)
	fmt.Printf("\nCarro [%s] chegando ao %s...", placa, ponto.Nome)
	time.Sleep(5 * time.Second)
	fmt.Printf("\nCarro [%s] carregando no %s...\n", placa, ponto.Nome)
	mensagemRecarga := "recarga," + placa + "," + ponto.Nome
	sendMessageForServer(mensagemRecarga, url)

	//simulando pagamento
	valorRecarga := ValueRandom() //gerando valor aleatório de pagamento
	fmt.Printf("\nRecarga do carro [%s] no %s finalizada! O valor total foi: R$ %s.", placa, ponto.Nome, valorRecarga)
	fmt.Printf("\nCarro [%s] efetuando pagamento...\n", placa)
	mensagemPagamento := "pagamento," + placa + "," + ponto.Nome
	sendMessageForServer(mensagemPagamento, url)
}

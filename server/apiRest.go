package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// AQUI SETAMOS OS SERVIDORES DE TODAS MAQUINAS
// AQUI SETAMOS OS SERVIDORES DE TODAS MAQUINAS
var servidores []string = []string{
	"http://server1:8091",
	"http://server2:8092",
	"http://server3:8093",
}

// O endereço setado como 0.0.0.0 serve para escutar conexões de qualquer
// endereço IP disponível na máquina
var enderecoUniversal = "0.0.0.0:"

// IP da maquina local (PRECISA SER MODIFICADO!!)
var ipLocal = os.Getenv("HOSTNAME")

// Porta da maquina
var portaLocal string

// Função responsável por iniciar a API REST em cada server
func startingREST(porta string) {
	portaLocal = porta

	createEndpoints()
	fmt.Printf("[API REST] Endpoints criados e iniciado API na porta %s\n", porta)
	endereco := enderecoUniversal + porta

	//ListenAndServe: escuta conexões na porta especificada e entrega as requisições para os handlers registrados
	//Cria-se uma goroutine para ficar sempre rodando a ApiRest
	go func() {
		if err := http.ListenAndServe(endereco, nil); err != nil {
			//fmt.Printf("Erro ao tentar iniciar API REST: %v\n", err)
		}
	}()
}

func createEndpoints() {
	http.HandleFunc("/blockchain", getBlockchain) //get
	http.HandleFunc("/add-block", postBlock)      //post
	http.HandleFunc("/status", getStatus)         //get
	http.HandleFunc("/mensagem", postMessage)     //post
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Erro ao decodificar a mensagem", http.StatusBadRequest)
		return
	}

	content := msg.Content
	fmt.Printf("Mensagem recebida: %s\n", content)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mensagem recebida com sucesso!")) //enviando retorno

	// Separar a string em uma lista
	listString := strings.Split(content, ",")
	if len(listString) < 3 {
		println("Mensagem incorreta!")
		return
	}

	tipoTransacao := listString[0]
	placa := listString[1]
	ponto := listString[2]

	conteudoParaSalvar := "[Placa do Veiculo: " + placa + ", Ponto de Recarga: " + ponto + "]"
	addNewBlock(tipoTransacao, conteudoParaSalvar)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func isPeerAlive(peer string) bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(peer + "/status")
	if err != nil {
		return false // Não respondeu = provavelmente offline
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	ReadBlocks()
	json.NewEncoder(w).Encode(blockchain.Blocos)
}

func postBlock(w http.ResponseWriter, r *http.Request) {
	var newBlock Block
	if err := json.NewDecoder(r.Body).Decode(&newBlock); err != nil {
		http.Error(w, "Bloco inválido", http.StatusBadRequest)
		return
	}

	lastBlock := blockchain.Blocos[len(blockchain.Blocos)-1]
	if isBlockValid(newBlock, lastBlock) {
		blockchain.Blocos = append(blockchain.Blocos, newBlock)
		SaveJSONBlockchain()
		fmt.Fprintln(w, "Bloco adicionado com sucesso.")
	} else {
		http.Error(w, "Bloco inválido ou fora de ordem", http.StatusBadRequest)
	}
}

// main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AQUI SETAMOS OS SERVIDORES DE TODAS MAQUINAS
var servidores []string = []string{
	"http://172.22.208.1:8091",
	"http://172.22.208.1:8093",
	"http://172.22.208.1:8092",
}

// O endereço setado como 0.0.0.0 serve para escutar conexões de qualquer
// endereço IP disponível na máquina
var enderecoUniversal = "0.0.0.0:"

// IP da maquina local (PRECISA SER MODIFICADO!!)
var ipLocal = "172.22.208.1"

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
	http.HandleFunc("/add-data", postData)        // <- novo endpoint
	http.HandleFunc("/status", getStatus)         //get
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func isPeerAlive(peer string) bool {
	if isTestMode() {
		return false // evita qualquer tentativa de rede em teste
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(peer)
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

func postData(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	addNewBlock(payload.Type, payload.Data)
	fmt.Fprintln(w, "Bloco adicionado com sucesso.")
}

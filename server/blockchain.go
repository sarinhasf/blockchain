package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Calcula o hash de um bloco
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Data + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Cria o bloco gênesis (o primeiro da cadeia)
func createGenesisBlock(tipo string, conteudo string) Block {
	genesis := Block{
		Index:     0,
		Timestamp: time.Now().String(),
		Type:      tipo,
		Data:      conteudo,
		PrevHash:  "",
	}
	genesis.Hash = calculateHash(genesis)
	return genesis
}

// Gera um novo bloco baseado no anterior
func generateBlock(oldBlock Block, data string, tipo string) Block {
	newBlock := Block{
		Index:     oldBlock.Index + 1,
		Timestamp: time.Now().String(),
		Type:      tipo,
		Data:      data,
		PrevHash:  oldBlock.Hash,
	}
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

// Verifica se o bloco é válido
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func addNewBlock(tipo string, conteudo string) {
	println("\nSincronizando com a rede antes de criar um novo bloco...")
	syncBlockchain(portaLocal) // Usa a variável global portaLocal
	println("Sincronização concluída. Prosseguindo com a criação do bloco.")

	vazio := false
	if len(blockchain.Blocos) == 0 {
		vazio = true
	}

	if vazio {
		// Cria o bloco gênesis (o primeiro block)
		genesisBlock := createGenesisBlock(tipo, conteudo)
		blockchain.Blocos = append(blockchain.Blocos, genesisBlock)
		SaveJSONBlockchain()            // Salva localmente primeiro//////////
		broadcastNewBlock(genesisBlock) // Depois, transmite para os outros nós

	} else {
		// Já existem blocos, vamos adicionar mais um na ponta da cadeia correta
		lastBlock := blockchain.Blocos[len(blockchain.Blocos)-1]
		newBlock := generateBlock(lastBlock, conteudo, tipo)

		if isBlockValid(newBlock, lastBlock) {
			blockchain.Blocos = append(blockchain.Blocos, newBlock)
			SaveJSONBlockchain()        // Salva localmente primeiro
			broadcastNewBlock(newBlock) // Depois, transmite para os outros nós
		} else {
			// Isso não deveria acontecer com a sincronização prévia, mas é uma boa proteção
			fmt.Println("[Erro] Bloco gerado localmente é inválido após a sincronização. Abortando.")
		}
	}
}

// Função para sincronizar blockchain com os nós existentes
func syncBlockchain(porta string) {
	println("Sicronizando blockchain entre os nós...")
	var longestChain []Block
	serverLocal := fmt.Sprintf("http://%s:%s", ipLocal, porta)
	//println(serverLocal)
	for _, server := range servidores {
		//println(server)
		if server == serverLocal {
			continue
		}

		//verificar se o nó está on ou of
		if !isPeerAlive(server) {
			fmt.Println(server, " está offline.")
			continue
		}
		fmt.Println(server, " está online. Buscando blockchain...")

		//pega todo arquivo blockchain do nó
		resp, err := http.Get(server + "/blockchain")
		if err != nil {
			fmt.Printf("[Sync] Erro ao acessar %s: %v\n", server, err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("[Sync] Erro ao ler resposta de %s: %v\n", server, err)
			continue
		}

		var remoteChain []Block
		if err := json.Unmarshal(body, &remoteChain); err != nil {
			fmt.Printf("[Sync] Erro ao decodificar blockchain de %s\n", server)
			continue
		}

		// Se a cadeia recebida for maior e válida, armazenar
		if len(remoteChain) > len(longestChain) && isChainValid(remoteChain) {
			longestChain = remoteChain
		}
	}

	localChainIsValid := isChainValid(blockchain.Blocos)
	networkChainIsValid := len(longestChain) > 0 // Verifica se a cadeia recebida é válida

	// Condição para substituir a blockchain local:
	// 1. A cadeia da rede é válida E é maior que a nossa. (Consenso normal)
	// OU
	// 2. A nossa cadeia local é INVÁLIDA e encontramos uma cadeia válida na rede. (Autocorreção)
	if (networkChainIsValid && len(longestChain) > len(blockchain.Blocos)) || (!localChainIsValid && networkChainIsValid) {
		fmt.Println("[Sync] Blockchain local desatualizada ou inválida. Substituindo pela versão da rede.")
		blockchain.Blocos = longestChain
		SaveJSONBlockchain() // Salva a versão correta
	} else if !localChainIsValid && !networkChainIsValid {
		fmt.Println("[Sync] ALERTA: A blockchain local está inválida e não foi possível encontrar uma cópia válida na rede.")
	} else {
		fmt.Println("[Sync] Blockchain local já está sincronizada e válida.")
	}
}

// Verifica se toda a cadeia é válida de acordo com hash anterior
func isChainValid(chain []Block) bool {
	for i := 1; i < len(chain); i++ { // Começa do segundo bloco
		if !isBlockValid(chain[i], chain[i-1]) { // Verifica se o bloco atual é válido em relação ao anterior
			return false
		}
	}
	return true
}

// adicionar novo bloco enviando para todos os nós
func broadcastNewBlock(block Block) {
	serverLocal := fmt.Sprintf("http://%s:%s", ipLocal, portaLocal)

	for _, server := range servidores {
		if server == serverLocal { // Evita enviar para si mesmo
			continue
		}

		//verificar se o nó está on ou of
		if isPeerAlive(server) {
			fmt.Println(server, " está online.")
		} else {
			fmt.Println(server, " está offline.")
			continue
		}

		jsonBlock, _ := json.Marshal(block)
		resp, err := http.Post(server+"/add-block", "application/json", bytes.NewBuffer(jsonBlock))
		if err != nil {
			fmt.Printf("[Broadcast] Falha ao enviar bloco para %s: %v\n", server, err)
			continue
		}
		resp.Body.Close()
	}
}

func getBlockchainFrom(peer string) {
	resp, err := http.Get(peer + "/blockchain")
	if err != nil {
		fmt.Printf("Erro ao obter blockchain de %s: %v\n", peer, err)
		return
	}
	defer resp.Body.Close()

	var blocos []Block
	if err := json.NewDecoder(resp.Body).Decode(&blocos); err != nil {
		fmt.Printf("Erro ao decodificar resposta de %s: %v\n", peer, err)
		return
	}

	fmt.Printf("\nBlockchain de %s:\n", peer)
	for i, bloco := range blocos {
		fmt.Printf("Bloco #%d: %+v\n", i, bloco)
	}
}

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func isTestMode() bool {
	return os.Getenv("IS_TEST") == "1"
}

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
	if isTestMode() {
		// 1) pega índice e prevHash sob lock
		blockchainMutex.Lock()
		idx := len(blockchain.Blocos)
		var prevHash string
		if idx > 0 {
			prevHash = blockchain.Blocos[idx-1].Hash
		}
		blockchainMutex.Unlock()

		// 2) monta o bloco (fora de qualquer lock)
		newBlock := Block{
			Index:     idx,
			Timestamp: "TESTE",
			Type:      tipo,
			Data:      conteudo,
			PrevHash:  prevHash,
		}
		newBlock.Hash = calculateHash(newBlock)

		// 3) só a região de append fica sob lock
		blockchainMutex.Lock()
		blockchain.Blocos = append(blockchain.Blocos, newBlock)
		blockchainMutex.Unlock()

		// 4) salva sem deadlock (SaveJSON faz seu próprio lock)
		if err := SaveJSONBlockchain(); err != nil {
			fmt.Printf("Erro ao salvar blockchain de teste: %v\n", err)
		}
		return
	}
	blockchainMutex.Lock()
	defer blockchainMutex.Unlock()

	println("\nVerificando os outros clientes conectados para sicronizar...")
	vazio := false
	if len(blockchain.Blocos) == 0 {
		vazio = true
	}

	if vazio {
		//fmt.Println("O arquivo blockchain.json está vazio.")
		// Cria o bloco gênesis (o primeiro block)
		genesisBlock := createGenesisBlock(tipo, conteudo)
		blockchain.Blocos = append(blockchain.Blocos, genesisBlock)
		broadcastNewBlock(genesisBlock) //salva o bloco nos outros nós

	} else {
		//fmt.Println("O arquivo blockchain.json contém dados.")
		// Já existe blocos, vamos add mais
		newBlock := generateBlock(blockchain.Blocos[len(blockchain.Blocos)-1], conteudo, tipo)
		if isBlockValid(newBlock, blockchain.Blocos[len(blockchain.Blocos)-1]) {
			blockchain.Blocos = append(blockchain.Blocos, newBlock) //verifica se esse bloco é válido comparado com o anterior
		}
		broadcastNewBlock(newBlock) //salva o bloco nos outros nós
	}

	//fmt.Println("Printando arquivo BLOCKCHAIN:")
	//for _, block := range blockchain.Blocos {
	//fmt.Printf("Index: %d\nTimestamp: %s\nData: %s\nHash: %s\nPrevHash: %s\n\n",
	//block.Index, block.Timestamp, block.Data, block.Hash, block.PrevHash)
	//}
	SaveJSONBlockchain() //salvando dados no json local
}

// Função para sincronizar blockchain com os nós existentes
func syncBlockchain(porta string) {
	if isTestMode() {
		fmt.Println("[TESTE] Ignorando syncBlockchain em ambiente de teste")
		return
	}

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
		if isPeerAlive(server) {
			fmt.Println(server, " está online.")
		} else {
			fmt.Println(server, " está offline.")
			continue
		}

		//pega todo arquivo blockchain do nó
		resp, err := http.Get(server + "/blockchain")
		if err != nil {
			fmt.Printf("[Sync] Erro ao acessar %s: %v\n", server, err)
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
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

	// Se encontrou uma cadeia mais longa válida, substituir a atual
	if len(longestChain) > len(blockchain.Blocos) {
		fmt.Println("[Sync] Substituindo blockchain local pela maior válida recebida")
		blockchain.Blocos = longestChain
		SaveJSONBlockchain()
	}
}

// Verifica se toda a cadeia é válida de acordo com hash anterior
func isChainValid(chain []Block) bool {
	for i := 1; i < len(chain); i++ {
		if !isBlockValid(chain[i], chain[i-1]) {
			return false
		}
	}
	return true
}

// adicionar novo bloco enviando para todos os nós
func broadcastNewBlock(block Block) {
	if isTestMode() {
		fmt.Println("[TESTE] Ignorando broadcastNewBlock em ambiente de teste")
		return
	}

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

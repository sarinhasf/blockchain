package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

const localBlockchainPath = "../../dados/blockchain.json"

/*
teste para a função addNewBlock
que adiciona blocos à blockchain de forma concorrente
*/
func TestConcurrentAddNewBlock(t *testing.T) {
	// 1. Ambiente de teste
	os.Setenv("IS_TEST", "1")
	os.Setenv("BLOCKCHAIN_PATH", localBlockchainPath)

	// 2. Preparar diretório e arquivo
	if err := os.MkdirAll("../../dados", os.ModePerm); err != nil {
		t.Fatalf("Erro ao criar pasta dados: %v", err)
	}
	if err := ioutil.WriteFile(localBlockchainPath, []byte(`{"blocos":[]}`), 0644); err != nil {
		t.Fatalf("Erro ao limpar blockchain.json: %v", err)
	}

	// 3. Inicializar a blockchain e criar o gênesis
	ReadBlocks()
	for i := 1; i < len(blockchain.Blocos); i++ {
		if blockchain.Blocos[i].Index != i {
			t.Errorf("Bloco fora de ordem: esperado índice %d, encontrado %d", i, blockchain.Blocos[i].Index)
		}
		if blockchain.Blocos[i].PrevHash != blockchain.Blocos[i-1].Hash {
			t.Errorf("PrevHash inválido no bloco %d", i)
		}
		expectedHash := calculateHash(blockchain.Blocos[i])
		if blockchain.Blocos[i].Hash != expectedHash {
			t.Errorf("Hash inválido no bloco %d", i)
		}
	}

	genesis := createGenesisBlock("init", "bloco inicial")
	blockchain.Blocos = []Block{genesis}
	if err := SaveJSONBlockchain(); err != nil {
		t.Fatalf("Erro ao salvar bloco gênesis: %v", err)
	}

	// 4. Disparar N goroutines adicionando blocos
	const numVeiculos = 10
	var wg sync.WaitGroup
	wg.Add(numVeiculos)
	for i := 0; i < numVeiculos; i++ {
		go func(i int) {
			defer wg.Done()
			addNewBlock("veiculo", "veiculo-"+fmt.Sprint(i))
		}(i)
	}
	wg.Wait()

	// 5) Verificar resultado
	ReadBlocks()
	expected := numVeiculos + 1 // bloco gênesis + veículos
	if got := len(blockchain.Blocos); got != expected {
		t.Errorf("Esperado %d blocos, mas encontrou %d", expected, got)
	}
}

// Testa concorrência de leitura e escrita na blockchain
// onde múltiplas goroutines tentam ler e escrever blocos simultaneamente.
func TestConcurrentReadAndWrite(t *testing.T) {
	/*
		testa se addNewBlock (escrita) e len(blockchain.Blocos) (leitura) podem rodar em paralelo sem corromper memória
		testa se todos os blocos esperados são adicionados corretamente
		testa se nenhum data race é reportado */

	os.Setenv("IS_TEST", "1")
	os.Setenv("BLOCKCHAIN_PATH", localBlockchainPath)

	// Limpa e reinicia a blockchain
	if err := os.MkdirAll("../../dados", os.ModePerm); err != nil {
		t.Fatalf("Erro ao criar pasta dados: %v", err)
	}
	if err := ioutil.WriteFile(localBlockchainPath, []byte(`{"blocos":[]}`), 0644); err != nil {
		t.Fatalf("Erro ao limpar blockchain.json: %v", err)
	}

	ReadBlocks()
	blockchain.Blocos = []Block{createGenesisBlock("init", "bloco inicial")}
	if err := SaveJSONBlockchain(); err != nil {
		t.Fatalf("Erro ao salvar bloco gênesis: %v", err)
	}

	const numWriters = 10
	const numReaders = 10
	var wg sync.WaitGroup

	// Lançar writers
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(i int) {
			defer wg.Done()
			addNewBlock("veiculo", fmt.Sprintf("veiculo-%d", i))
		}(i)
	}

	// Lançar readers
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func(i int) {
			defer wg.Done()
			blockchainMutex.Lock()
			_ = len(blockchain.Blocos) // leitura protegida
			blockchainMutex.Unlock()
		}(i)
	}

	wg.Wait()

	// Verificação final
	ReadBlocks()
	expected := numWriters + 1 // bloco gênesis + writers
	if got := len(blockchain.Blocos); got != expected {
		t.Errorf("Esperado %d blocos, mas encontrou %d", expected, got)
	}
}

// Testa a persistência da blockchain
// onde blocos são adicionados, salvos e recarregados do arquivo JSON.
func TestBlockchainPersistence(t *testing.T) {
	/*
		1. Limpar JSON manualmente - Garantia de início limpo
		2. Salvar blocos - Se o JSON é corretamente persistido
		3. Recarregar com `ReadBlocks` - Se o disco -> memória funciona
		4. Validar conteúdo            - Se os blocos lidos são fiéis       */

	os.Setenv("IS_TEST", "1")
	os.Setenv("BLOCKCHAIN_PATH", localBlockchainPath)

	// Limpa blockchain
	if err := os.MkdirAll("../../dados", os.ModePerm); err != nil {
		t.Fatalf("Erro ao criar pasta dados: %v", err)
	}
	if err := ioutil.WriteFile(localBlockchainPath, []byte(`{"blocos":[]}`), 0644); err != nil {
		t.Fatalf("Erro ao limpar blockchain.json: %v", err)
	}

	// Inicia com gênesis
	ReadBlocks()
	blockchain.Blocos = []Block{createGenesisBlock("init", "bloco inicial")}
	if err := SaveJSONBlockchain(); err != nil {
		t.Fatalf("Erro ao salvar bloco gênesis: %v", err)
	}

	// Adiciona dois blocos
	addNewBlock("veiculo", "veiculo-A")
	addNewBlock("veiculo", "veiculo-B")

	// Simula reinício: recarrega do JSON
	blockchain.Blocos = nil // limpa memória
	ReadBlocks()

	// Verifica se blocos estão todos presentes
	if len(blockchain.Blocos) != 3 {
		t.Fatalf("Esperado 3 blocos (gênesis + 2), encontrado %d", len(blockchain.Blocos))
	}

	// Verifica se dados estão corretos
	expected := []string{"bloco inicial", "veiculo-A", "veiculo-B"}
	for i, blk := range blockchain.Blocos {
		if blk.Data != expected[i] {
			t.Errorf("Bloco %d com conteúdo errado: esperado '%s', obtido '%s'", i, expected[i], blk.Data)
		}
	}
}

/*
Testa a consistência do bloco gênesis
onde o bloco gênesis é criado, salvo e recarregado do arquivo JSON.
*/
func TestGenesisBlockConsistency(t *testing.T) {
	/* Garantir que o bloco gênesis:
	Sempre tenha índice 0
	Tenha PrevHash vazio
	Tenha um Hash válido
	Tenha os dados esperados*/
	os.Setenv("IS_TEST", "1")
	os.Setenv("BLOCKCHAIN_PATH", localBlockchainPath)

	// Limpa blockchain
	if err := os.MkdirAll("../../dados", os.ModePerm); err != nil {
		t.Fatalf("Erro ao criar pasta dados: %v", err)
	}
	if err := ioutil.WriteFile(localBlockchainPath, []byte(`{"blocos":[]}`), 0644); err != nil {
		t.Fatalf("Erro ao limpar blockchain.json: %v", err)
	}

	// Criar e salvar apenas o bloco gênesis
	ReadBlocks()
	genesis := createGenesisBlock("init", "bloco inicial")
	blockchain.Blocos = []Block{genesis}
	if err := SaveJSONBlockchain(); err != nil {
		t.Fatalf("Erro ao salvar bloco gênesis: %v", err)
	}

	// Simular reinício do sistema
	blockchain.Blocos = nil
	ReadBlocks()

	// Verificações do bloco gênesis
	if len(blockchain.Blocos) != 1 {
		t.Fatalf("Blockchain deveria conter apenas 1 bloco, encontrou %d", len(blockchain.Blocos))
	}

	blk := blockchain.Blocos[0]

	if blk.Index != 0 {
		t.Errorf("Índice do bloco gênesis incorreto: esperado 0, obtido %d", blk.Index)
	}
	if blk.PrevHash != "" {
		t.Errorf("PrevHash do bloco gênesis deve ser vazio, obtido '%s'", blk.PrevHash)
	}
	expectedHash := calculateHash(blk)
	if blk.Hash != expectedHash {
		t.Errorf("Hash do bloco gênesis incorreto: esperado %s, obtido %s", expectedHash, blk.Hash)
	}
	if blk.Data != "bloco inicial" {
		t.Errorf("Conteúdo do bloco gênesis incorreto: esperado 'bloco inicial', obtido '%s'", blk.Data)
	}
}

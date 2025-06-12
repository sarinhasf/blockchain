package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
Structs usadas para construção do projeto.
*/

type Block struct {
	Index     int
	Timestamp string //marcacao de tempo
	Type      string //tipo do conteudo armazenado
	Data      string //conteudo
	PrevHash  string //hash do bloco anterior
	Hash      string //hash do bloco
}

type Blockchain struct {
	Blocos []Block `json:"blocos"`
}

type Message struct {
	Content string `json:"content"`
}

/*
Váriaveis Globais.
*/
var blockchain Blockchain

/*
ReadJSONFile: lê qualquer arquivo JSON da pasta dados e deserializa para a struct fornecida.
*/
func ReadJSONFile(fileName string, target interface{}) error {
	fullPath := fileName

	// Abrir o arquivo
	file, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo %s: %w", fullPath, err)
	}
	defer file.Close()

	// Ler o conteúdo
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("erro ao ler o arquivo %s: %w", fullPath, err)
	}

	// Deserializar o JSON no alvo fornecido
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("erro ao decodificar o JSON do arquivo %s: %w", fullPath, err)
	}

	return nil
}

/*
ReadBlocks: lê todos pontos do arquivo blockchain da pasta dados.
*/
func ReadBlocks() {
	err := ReadJSONFile("/app/blockchain.json", &blockchain)
	if err != nil {
		fmt.Println("Erro:", err)
	} else {
		fmt.Printf("Pontos lidos do arquivo blockchain com sucesso!\n")
	}
}

/*
CheckBlocks: verifica se a cadeia de blocos está ok antes do servidor rodar.
*/
func CheckBlocks() bool {
	blocos := blockchain.Blocos

	// A verificação só é necessária se houver 2 ou mais blocos.
	// Itera a partir do segundo bloco (índice 1).
	for i := 1; i < len(blocos); i++ {
		blocoAtual := blocos[i]
		blocoAnterior := blocos[i-1]

		// Se o hash do bloco anterior não corresponder ao PrevHash do bloco atual, a cadeia é inválida.
		if blocoAtual.PrevHash != blocoAnterior.Hash {
			fmt.Printf("Inconsistência encontrada entre o Bloco %d e o Bloco %d\n", i-1, i)
			fmt.Printf("Hash do Bloco %d: %s\n", i-1, blocoAnterior.Hash)
			fmt.Printf("PrevHash do Bloco %d: %s\n", i, blocoAtual.PrevHash)
			return false // Encontrou um erro, a cadeia é inválida.
		}
	}

	// Se o laço terminar sem encontrar erros, a cadeia é válida.
	return true
}

/*
SaveBlockchain: salva os dados do blockchain no json novamente.
*/
func SaveJSONBlockchain() error {
	fullPath := filepath.Join("/app/", "blockchain.json")

	// Garantir que a pasta "dados" exista
	err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar diretório dados: %w", err)
	}

	// Serializar a variável blockchain para JSON
	bytes, err := json.MarshalIndent(blockchain, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao codificar blockchain para JSON: %w", err)
	}

	// Salvar no arquivo
	err = ioutil.WriteFile(fullPath, bytes, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo blockchain.json: %w", err)
	}

	fmt.Println("\nDados da Blockchain salva com sucesso!")
	return nil
}

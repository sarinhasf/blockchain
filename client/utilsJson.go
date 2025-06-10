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
type Point struct {
	ID        int     `json:"id"`
	Nome      string  `json:"nome"`
	Cidade    string  `json:"cidade"`
	Estado    string  `json:"estado"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Reservado string  `json:"reservado"`
}

type PointsWrapper struct {
	PontosDeRecarga []Point `json:"pontos_de_recarga"`
}

// Estrutura de um bloco
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

/*
Váriaveis Globais.
*/
var dataPoints PointsWrapper
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
ReadPoints: lê todos pontos do arquivo points da pasta dados.
*/
func ReadPoints() {
	err := ReadJSONFile("points.json", &dataPoints)
	if err != nil {
		fmt.Println("Erro:", err)
	} else {

		//fmt.Printf("Pontos lidos do arquivo points com sucesso!\n")
	}
}

/*
ReadBlocks: lê todos pontos do arquivo blockchain da pasta dados.
*/
func ReadBlocks() {
	err := ReadJSONFile("/app/blockchain.json", &blockchain)
	if err != nil {
		fmt.Println("Erro:", err)
	} else {
		//fmt.Printf("Pontos lidos do arquivo blockchain com sucesso!\n")
	}
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

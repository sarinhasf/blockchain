package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

type Empresa struct {
	ID     int    `json:"id"`
	Nome   string `json:"nome"`
	Pontos []int  `json:"pontos"`
}

type Companies struct {
	Empresas []Empresa `json:"empresas"`
}

type Message struct {
	Content string `json:"content"`
}

/*
Váriaveis Globais.
*/
var dataPoints PointsWrapper
var dataCompanies Companies

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
ReadCompanies: lê todos as empresas do arquivo companies da pasta dados.
*/
func ReadCompanies() {
	err := ReadJSONFile("companies.json", &dataCompanies)
	if err != nil {
		fmt.Println("Erro:", err)
	} else {
		//fmt.Printf("Pontos lidos do arquivo points com sucesso!\n")
	}
}

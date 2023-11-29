package database

import (
	"database/sql"
)

type ClientDB struct {
	DB *sql.DB
}

type Client struct {
	db      *sql.DB
	Cpf     string
	Nome    string
	Celular string
}

type Dados struct {
	db        *sql.DB
	CELULAR   string
	NuCpf     string
	NomeSeg   string
	DtNasc    string
	Especie   string
	Salario   string
	Banco     string
	Agencia   string
	Conta     string
	BcoEmp    string
	Contrato  string
	ValorEmp  string
	IniEmp    string
	FimEmp    string
	ParcEmp   string
	ValorParc string
	Endereco  string
	Bairro    string
	Municipio string
	Uf        string
	Cep       string
	Celular   string
	Tel1      string
	Tel2      string
	Tel3      string
	Idade     string
}

func (c *ClientDB) GetClient(cpf string) (*Client, error) {
	row := c.DB.QueryRow("SELECT cpf, nome, celular FROM pbnew WHERE cpf = $1", cpf)
	var client Client
	err := row.Scan(&client.Cpf, &client.Nome, &client.Celular)

	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *ClientDB) GetDados(nucpf string) (*Dados, error) {
	row := c.DB.QueryRow("SELECT nucpf, nomesegurado, esp FROM clientsnew WHERE nucpf = $1", nucpf)
	var dados Dados
	err := row.Scan(&dados.NuCpf, &dados.NomeSeg, &dados.Especie)

	if err != nil {
		return nil, err
	}
	return &dados, nil
}

package mysql

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/brunocsmaciel/rinha-backend-2024/types"
)

type TransacaoStore interface {
	RealizaTransacao(types.TransacaoInput) (*types.TransacaoResponse, error)
}

type MySqlTransacaoStore struct {
	db *sql.DB
}

func NewMySqlTransacaoStore(db *sql.DB) *MySqlTransacaoStore {
	return &MySqlTransacaoStore{
		db: db,
	}
}

type AppError struct {
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

const transacaoQuery = `INSERT INTO transacoes (cliente_id, valor, tipo, descricao) VALUES (?, ?, ?, ?)`
const atualizacaoSaldoQuery = `UPDATE saldos SET valor = ? WHERE cliente_id = ?`

func (mySql *MySqlTransacaoStore) RealizaTransacao(transacaoInput types.TransacaoInput) (*types.TransacaoResponse, error) {

	extrato, err := BuscarSaldo(NewMySqlExtratoStore(mySql.db), transacaoInput.ClienteId)
	if err != nil {
		msg := fmt.Sprintf("Cliente com ID %s não encontrado", transacaoInput.ClienteId)
		return nil, &AppError{Message: msg, StatusCode: 404}
	}

	limite := extrato.Limite
	saldoAtualizado := extrato.Total

	if transacaoInput.Tipo == "c" {
		saldoAtualizado += transacaoInput.Valor
	} else if transacaoInput.Tipo == "d" {
		valorControle := saldoAtualizado - transacaoInput.Valor

		saldoAbsoluto := math.Abs(float64(valorControle))
		if saldoAbsoluto > float64(limite) {
			mensagemErro := fmt.Sprintf("valor do débito %d maior que o permitido. Limite da conta = %d, saldo atual = %d",
				transacaoInput.Valor, limite, saldoAtualizado)
			return nil, &AppError{Message: mensagemErro, StatusCode: 422}
		}
		saldoAtualizado = valorControle
	}
	_, err = mySql.db.Exec(transacaoQuery,
		transacaoInput.ClienteId, transacaoInput.Valor, transacaoInput.Tipo, transacaoInput.Descricao)
	if err != nil {
		return nil, err
	}

	_, err = mySql.db.Exec(atualizacaoSaldoQuery, saldoAtualizado, transacaoInput.ClienteId)
	if err != nil {
		return nil, err
	}
	conexoes := mySql.db.Stats().InUse
	fmt.Println("Conexoes em uso", conexoes)

	return &types.TransacaoResponse{
		Limite: limite,
		Saldo:  saldoAtualizado,
	}, nil
}

package mysql

import (
	"database/sql"
	"log"
	"time"

	"github.com/brunocsmaciel/rinha-backend-2024/types"
	_ "github.com/go-sql-driver/mysql"
)

type ExtratoStore interface {
	BuscarExtrato(idCliente string) (*Extrato, error)
}

type MySqlExtratoStore struct {
	db *sql.DB
}

func NewMySqlExtratoStore(db *sql.DB) *MySqlExtratoStore {
	return &MySqlExtratoStore{
		db: db,
	}
}

type Extrato struct {
	Saldo             types.Saldo       `json:"saldo"`
	UltimasTransacoes []types.Transacao `json:"ultimas_transacoes"`
}

const saldoQuery = `SELECT saldos.valor, clientes.limite 
					FROM clientes 
					JOIN saldos 
					ON clientes.id = saldos.id
					WHERE clientes.id = ?
					LIMIT 1`

const transacoesQuery = `SELECT t.valor, t.tipo, t.descricao, t.realizada_em 
						FROM transacoes t
						WHERE t.cliente_id = ?
						ORDER BY t.realizada_em DESC
						LIMIT 10`

func (mySql *MySqlExtratoStore) BuscarExtrato(idCliente string) (*Extrato, error) {

	saldo, err := BuscarSaldo(mySql, idCliente)
	if err != nil {
		return nil, err
	}

	transacoes, err := BuscarTransacoes(mySql, idCliente)
	if err != nil {
		return nil, err
	}
	var transacoesValores []types.Transacao
	for _, transacao := range transacoes {
		transacoesValores = append(transacoesValores, *transacao)
	}
	return &Extrato{
		Saldo:             *saldo,
		UltimasTransacoes: transacoesValores,
	}, nil

}

func BuscarSaldo(mySql *MySqlExtratoStore, idCliente string) (*types.Saldo, error) {
	var saldo types.Saldo
	err := mySql.db.QueryRow(saldoQuery, idCliente).
		Scan(&saldo.Total, &saldo.Limite)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Cliente com ID %s não encontrado", idCliente)
			return nil, err
		}
		log.Fatal(err)
	}
	saldo.DataExtrato = time.Now().UTC().Format("2006-01-02T15:04:05.9999999Z07:00")

	return &saldo, nil

}

func BuscarTransacoes(mySql *MySqlExtratoStore, idCliente string) ([]*types.Transacao, error) {
	var transacoes []*types.Transacao
	res, err := mySql.db.Query(transacoesQuery, idCliente)
	if err != nil {
		log.Printf("Erro durante query %v", err)
		return nil, err

	}
	for res.Next() {
		var transacaoAtual types.Transacao
		res.Scan(&transacaoAtual.Valor, &transacaoAtual.Tipo, &transacaoAtual.Descricao, &transacaoAtual.RealizadaEm)
		transacoes = append(transacoes, &transacaoAtual)
	}

	if transacoes == nil {
		transacoes = []*types.Transacao{}
	}

	return transacoes, nil
}

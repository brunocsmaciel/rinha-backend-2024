package main

import (
	"github.com/brunocsmaciel/rinha-backend-2024/db/mysql"
	"github.com/brunocsmaciel/rinha-backend-2024/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	sql := mysql.Init()
	var extratoStore = mysql.NewMySqlExtratoStore(sql)
	var transacaoStore = mysql.NewMySqlTransacaoStore(sql)

	transacaoHandler := handlers.NewTransacaoHandler(transacaoStore)
	extratoHandler := handlers.NewClienteHandler(extratoStore)
	app.Get("/clientes/:id/extrato", extratoHandler.BuscarExtrato)
	app.Post("/clientes/:id/transacoes", transacaoHandler.ProcessaTransacao)

	app.Listen(":8080")
	defer sql.Close()
}

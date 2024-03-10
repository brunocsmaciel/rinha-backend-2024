package handlers

import (
	"fmt"

	"github.com/brunocsmaciel/rinha-backend-2024/db/mysql"
	"github.com/gofiber/fiber/v2"
)

type ExtratoHandler struct {
	extratoStore mysql.ExtratoStore
}

func NewClienteHandler(extratoStore mysql.ExtratoStore) *ExtratoHandler {

	return &ExtratoHandler{
		extratoStore: extratoStore,
	}
}

func (clienteH *ExtratoHandler) BuscarExtrato(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	cliente, err := clienteH.extratoStore.BuscarExtrato(id)
	if cliente == nil {
		return ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Cliente com ID %s não encontrado", id))
	}
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).
			SendString("erro interno do servidor")
	}

	return ctx.JSON(cliente)
}

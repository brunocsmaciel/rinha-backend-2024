package handlers

import (
	"fmt"

	"github.com/brunocsmaciel/rinha-backend-2024/db/mysql"
	"github.com/brunocsmaciel/rinha-backend-2024/types"
	"github.com/gofiber/fiber/v2"
)

type TransacaoHandler struct {
	transacaoStore mysql.TransacaoStore
}

func NewTransacaoHandler(transacaoStore mysql.TransacaoStore) *TransacaoHandler {
	return &TransacaoHandler{
		transacaoStore: transacaoStore,
	}

}

func (transacaoH *TransacaoHandler) ProcessaTransacao(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var transacaoInput types.TransacaoInput
	err := ctx.BodyParser(&transacaoInput)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	errosMap := validate(&transacaoInput)
	if errosMap != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"erros": errosMap})
	}

	transacaoInput.ClienteId = id

	transacao, err := transacaoH.transacaoStore.RealizaTransacao(transacaoInput)
	if err != nil {
		if appError, ok := err.(*mysql.AppError); ok {
			switch appError.StatusCode {
			case 404:
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": appError.Message})
			case 422:
				return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": appError.Message})
			default:
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
			}
		}
	}
	return ctx.JSON(transacao)
}

func validate(transacaoInput *types.TransacaoInput) map[string]string {
	errorMap := make(map[string]string)
	descricao := transacaoInput.Descricao
	if len(descricao) > 10 || len(descricao) < 1 {
		msg := fmt.Sprintf("descrição deve ter tamanho entre 1 e 10. Descrição %s é inválida", descricao)
		errorMap["descricao"] = msg
	}

	tipo := transacaoInput.Tipo
	if len(tipo) != 1 || (tipo != "c" && tipo != "d") {
		msg := fmt.Sprintf("tipos válidos: 'c' ou 'd'. Tipo %s é inválido", tipo)
		errorMap["tipo"] = msg
	}

	valor := transacaoInput.Valor
	if valor < 0 {
		msg := fmt.Sprintf("valor deve ser um número inteiro positivo. Valor %d é inválido", valor)
		errorMap["valor"] = msg
	}

	if len(errorMap) > 0 {
		return errorMap
	}

	return nil
}

package bankdetails

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/services"
)

func apiGetBankDetails(c *fiber.Ctx) error {
	bullionId := c.Query("bullionId")
	if bullionId == "" {
		return &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "bullionId is required",
			Name:       "INVALID_INPUT",
		}
	}
	if err := interfaces.ValidateBullionIdMatchingInToken(c, bullionId); err != nil {
		return err
	}
	entity, err := services.BankDetailsService.GetBankDetailsByBullionId(bullionId)
	if err != nil {
		return err
	} else {
		return c.JSON(entity)
	}
}

package validator

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	validator *validator.Validate
}

func NewValidator() *Service {
	return &Service{
		validator: validator.New(),
	}
}

func (v *Service) ParseAndValidate(c *fiber.Ctx, data interface{}) error {
	if err := c.BodyParser(data); err != nil {
		return err
	}

	if err := v.validator.Struct(data); err != nil {
		return err
	}

	return nil
}

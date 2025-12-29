package utils

import (
	"bytes"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thedevsaddam/govalidator"
)

func ValidateJSON(c *fiber.Ctx, data interface{}, rules govalidator.MapData) map[string][]string {
	body := c.Body()
	reader := bytes.NewReader(body)
	req, _ := http.NewRequest("POST", "/", reader)
	req.Header.Set("Content-Type", "application/json")

	opts := govalidator.Options{
		Request: req,
		Data:    data,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	errs := v.ValidateJSON()

	if len(errs) > 0 {
		errors := make(map[string][]string)
		for field, msgs := range errs {
			errors[field] = msgs
		}
		return errors
	}

	return nil
}

func ValidateJSONWithMessages(c *fiber.Ctx, data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	body := c.Body()
	reader := bytes.NewReader(body)
	req, _ := http.NewRequest("POST", "/", reader)
	req.Header.Set("Content-Type", "application/json")

	opts := govalidator.Options{
		Request:  req,
		Data:     data,
		Rules:    rules,
		Messages: messages,
	}

	v := govalidator.New(opts)
	errs := v.ValidateJSON()

	if len(errs) > 0 {
		errors := make(map[string][]string)
		for field, msgs := range errs {
			errors[field] = msgs
		}
		return errors
	}

	return nil
}

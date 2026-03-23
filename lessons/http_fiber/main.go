package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

/*
curl -XPOST localhost:8282/someRoute -H "Content-Type: application/json" -d "{\"name\":\"John\",\"age\":30}"
curl -XGET "localhost:8282/someOtherRoute/123?query=test" -H "Content-Type: application/json"
*/

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		for _, tag := range []string{"json", "uri", "query"} {
			if raw := fld.Tag.Get(tag); raw != "" {
				name := strings.SplitN(raw, ",", 2)[0]
				if name == "-" {
					return ""
				}
				if name != "" {
					return name
				}
			}
		}
		return fld.Name
	})

	return v
}

type SomeRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age"`
}

type SomeOtherRequestParams struct {
	Param string `uri:"param" validate:"required,alpha"`
}

type SomeOtherRequestQuery struct {
	Foo string `query:"foo" validate:"required"`
}

var V = NewValidator()

const (
	REQUEST_KEY = "request_key"
	QUERY_KEY   = "query_key"
	PARAM_KEY   = "param_key"
)

func messageField(fe validator.FieldError) string {
	switch fe.Tag() {
	case "alpha":
		return "Поле должно состоять из букв"
	case "required":
		return "Поле обязательно для заполнения"
	default:
		return "Поле не соответствует правилу " + fe.Tag()
	}
}

func ValidationErrors(c fiber.Ctx, err error) error {
	if err, ok := err.(validator.ValidationErrors); ok && len(err) > 0 {
		out := make(map[string]string)
		for _, e := range err {
			out[e.Field()] = messageField(e)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_errors": out})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
}

func validateSomeRequestMiddleware(c fiber.Ctx) error {
	var req = new(SomeRequest)
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if verrors := V.Struct(req); verrors != nil {
		return ValidationErrors(c, verrors)
	}

	c.Locals(REQUEST_KEY, req)
	return c.Next()
}

func validateSomeOtherRequestMiddleware(c fiber.Ctx) error {
	var req = new(SomeOtherRequestParams)
	req.Param = c.Params("param")

	var query = new(SomeOtherRequestQuery)
	query.Foo = c.Query("foo")

	if verrors := V.Struct(req); verrors != nil {
		return ValidationErrors(c, verrors)
	}

	if verrors := V.Struct(query); verrors != nil {
		return ValidationErrors(c, verrors)
	}

	c.Locals(PARAM_KEY, req)
	c.Locals(QUERY_KEY, query)
	return c.Next()
}

func SomeRouteHandler(c fiber.Ctx) error {
	var req = c.Locals(REQUEST_KEY).(*SomeRequest)

	return c.Status(fiber.StatusOK).JSON(req)
}

func SomeOtherRouteHandler(c fiber.Ctx) error {
	var req = c.Locals(PARAM_KEY).(*SomeOtherRequestParams)
	var query = c.Locals(QUERY_KEY).(*SomeOtherRequestQuery)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"param": req.Param,
		"foo":   query.Foo,
	})
}

func main() {
	app := fiber.New()
	app.Post("/someRoute", validateSomeRequestMiddleware, SomeRouteHandler)
	app.Get("/someOtherRoute/:param", validateSomeOtherRequestMiddleware, SomeOtherRouteHandler)

	err := app.Listen(":8282")

	if err != nil {
		fmt.Print(err)
	}
}

package main

import (
	"github.com/feryadialoi/x-order-service/contract"
	"github.com/feryadialoi/x-order-service/lib/echo-middleware"
	"github.com/feryadialoi/x-order-service/lib/kafka"
	"github.com/feryadialoi/x-order-service/lib/must"
	"github.com/feryadialoi/x-order-service/order/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(echomiddleware.ErrorMiddleware())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	producer := must.Must(kafka.NewSyncProducer[contract.OrderSubmitted](
		[]string{"localhost:9092"},
		"order-saga",
	))

	order := handler.NewOrder(producer)

	e.POST("/api/v1/orders", func(c echo.Context) error {
		var req handler.SubmitOrderRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		res, err := order.SubmitOrder(c.Request().Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(200, res)
	})

	must.Run(e.Start(":8080"))
}

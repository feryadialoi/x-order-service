package echomiddleware

import "github.com/labstack/echo/v4"

func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				return c.JSON(400, err)
			}
			return nil
		}
	}
}

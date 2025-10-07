package currency

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CurrencyAlarmRouter struct {
	createCurrencyAlarmUC CreateCurrencyAlarmUC
}

func NewCurrencyAlarmRouter(createCurrencyAlarmUC CreateCurrencyAlarmUC) *CurrencyAlarmRouter {
	return &CurrencyAlarmRouter{
		createCurrencyAlarmUC: createCurrencyAlarmUC,
	}
}

func (car CurrencyAlarmRouter) Route(e *echo.Echo) {
	group := e.Group("/currency-alarms")

	group.POST("/", func(c echo.Context) error {
		var dto CreateCurrencyAlarmDto
		if err := c.Bind(&dto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err := car.createCurrencyAlarmUC.Execute(c.Request().Context(), dto)	
		
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusCreated)
	})
}
package currency

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CurrencyAlarmRouter struct {
	createCurrencyAlarmUC CreateCurrencyAlarmUC
	deleteCurrencyAlarmUC DeleteCurrencyAlarmUC
}

func NewCurrencyAlarmRouter(createCurrencyAlarmUC CreateCurrencyAlarmUC, deleteCurrencyAlarmUC DeleteCurrencyAlarmUC) *CurrencyAlarmRouter {
	return &CurrencyAlarmRouter{
		createCurrencyAlarmUC: createCurrencyAlarmUC,
		deleteCurrencyAlarmUC: deleteCurrencyAlarmUC,
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

	group.DELETE("/:id", func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = car.deleteCurrencyAlarmUC.Execute(c.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error()) 
		}

		return c.NoContent(http.StatusNoContent)
	})
}
package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/josestg/problemdetail"
	"github.com/labstack/echo/v4"
)

func ProblemDetailHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	// If it's a ProblemDetail
	var pd *problemdetail.ProblemDetail
	if errors.As(err, &pd) {
		// write the RFC7807 response
		// set content type
		c.Response().Header().Set(echo.HeaderContentType, "application/problem+json")
		c.Response().WriteHeader(pd.Status)
		if err2 := json.NewEncoder(c.Response()).Encode(pd); err2 != nil {
			c.Logger().Error("failed to write problem detail:", err2)
		}
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		pd := problemdetail.New(
			"about:blank",
			problemdetail.WithTitle(http.StatusText(he.Code)),
			problemdetail.WithDetail(fmt.Sprintf("%v", he.Message)),
			problemdetail.WithInstance(c.Request().URL.Path),
		)
		pd.Status = he.Code
		c.Response().Header().Set(echo.HeaderContentType, "application/problem+json")
		c.Response().WriteHeader(he.Code)
		if err2 := json.NewEncoder(c.Response()).Encode(pd); err2 != nil {
			c.Logger().Error("failed to write problem detail:", err2)
		}
		return
	}

	pdFallback := problemdetail.New(
		"about:blank",
		problemdetail.WithTitle("Internal Server Error"),
		problemdetail.WithDetail(err.Error()),
		problemdetail.WithInstance(c.Request().URL.Path),
	)
	pdFallback.Status = http.StatusInternalServerError
	c.Response().Header().Set(echo.HeaderContentType, "application/problem+json")
	c.Response().WriteHeader(http.StatusInternalServerError)
	if err2 := json.NewEncoder(c.Response()).Encode(pdFallback); err2 != nil {
		c.Logger().Error("failed to write problem detail:", err2)
	}
}

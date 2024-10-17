package v1

import (
	"errors"
	"net/http"
	"vio_coding_challenge/internal/domain"
	"vio_coding_challenge/internal/service"

	"github.com/gin-gonic/gin"
)

type GeoCtrl struct {
	geo service.GeoServiceI
}

func NewGeoCtrl(geo service.GeoServiceI) GeoCtrlI {
	return &GeoCtrl{geo: geo}
}

// Get
// @Summary Get geolocation by IP address
// @Tags geolocation
// @Accept json
// @Produce json
// @Param ip_address path string true "ip address"
// @Success 200 {object} domain.Geolocation
// @failure 404 {object} domain.Error
// @failure 500 {object} domain.Error
// @Router /api/v1/geolocations/{ip_address} [get]
func (gc GeoCtrl) Get(c *gin.Context) {
	geo, err := gc.geo.GetByIP(c, c.Param(ipParam))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, domain.Error{Msg: "Not found geolocation by IP address."})
		default:
			c.JSON(http.StatusInternalServerError, domain.Error{Msg: err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, geo)
}

package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"vio_coding_challenge/internal/domain"
	mocks "vio_coding_challenge/mocks/service"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GeoCtrl_Get_Ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	geoService := mocks.NewMockGeoServiceI(ctrl)

	ip := "1.1.1.1"

	geo := domain.Geolocation{
		IPAddress:    ip,
		CountryCode:  "EN",
		Country:      "England",
		City:         "Hogsmead",
		Latitude:     10.1,
		Longitude:    -10.1,
		MysteryValue: "Hogwarts",
	}

	geoService.EXPECT().GetByIP(gomock.Any(), gomock.Eq(ip)).Return(&geo, nil)

	r := gin.Default()
	r.GET("/api/v1/geolocations/:"+ipParam, NewGeoCtrl(geoService).Get)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/geolocations/"+ip, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t,
		"{\"ip_address\":\"1.1.1.1\",\"country_code\":\"EN\",\"country\":\"England\","+
			"\"city\":\"Hogsmead\",\"latitude\":10.1,\"longitude\":-10.1,\"mystery_value\":\"Hogwarts\"}",
		w.Body.String(),
	)
}

func Test_GeoCtrl_Get_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	geoService := mocks.NewMockGeoServiceI(ctrl)

	ip := "1.1.1.1"

	r := gin.Default()
	r.GET("/api/v1/geolocations/:"+ipParam, NewGeoCtrl(geoService).Get)

	cases := []struct {
		name string
		data string
		err  error
		code int
		body string
	}{
		{
			name: "code 404",
			err:  domain.ErrNotFound,
			code: http.StatusNotFound,
			body: "{\"message\":\"Not found geolocation by IP address.\"}",
		},
		{
			name: "code 500",
			err:  errors.New("undefined error"),
			code: http.StatusInternalServerError,
			body: "{\"message\":\"undefined error\"}",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			geoService.EXPECT().GetByIP(gomock.Any(), gomock.Eq(ip)).Return(nil, c.err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/geolocations/"+ip, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, c.code, w.Code)
			assert.Equal(t, c.body, w.Body.String())
		})
	}
}

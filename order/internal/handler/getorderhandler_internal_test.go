package handler

import (
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"test/order/internal/config"
	"test/order/internal/svc"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

var configFile = flag.String("f", "../../etc/order.yaml", "the config file")

type HandlerSuite struct {
	suite.Suite
	ctx *svc.ServiceContext
}

func (s *HandlerSuite) SetupSuite() {
	var c config.Config

	if err := conf.LoadConfig(*configFile, &c); err != nil {
		s.T().Error(err)
	}

	s.ctx = svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	RegisterHandlers(server, s.ctx)
}

func (s *HandlerSuite) TestGetOrder() {
	handler := getOrderHandler(s.ctx)
	cases := []struct {
		method, url string
		pathvar     map[string]string
		in          io.Reader
		want        int
	}{
		{method: "GET", url: "/api/order/get/", pathvar: nil, in: nil, want: 400},
		{method: "GET", url: "/api/order/get/", pathvar: map[string]string{"id": "123"}, in: nil, want: 200},
	}

	for _, c := range cases {
		r, err := http.NewRequest(c.method, c.url, c.in)
		assert.NoError(s.T(), err)

		r = pathvar.WithVars(r, c.pathvar)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		assert.Equal(s.T(), c.want, w.Code, "response code should be equal")
	}
}

func TestMain(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

package pos_indonesia_test

import (
	"testing"
	"time"

	"github.com/BTechnopark/ipostal/pkg/cache"
	"github.com/BTechnopark/ipostal/src/client"
	"github.com/BTechnopark/ipostal/src/pos_indonesia"
	"github.com/BTechnopark/ipostal/src/session"
	"github.com/stretchr/testify/assert"
)

type configMock struct{}

// GetBaseUrl implements api.PostalConfig.
func (c *configMock) GetBaseUrl() string {
	return "https://kodepos.posindonesia.co.id"
}

func TestFindPostalCode(t *testing.T) {
	sess := session.NewSession("test_session")
	c := cache.NewCache(time.Minute)
	a := client.NewApi(sess)
	api := pos_indonesia.NewPosIndonesiaApi(&configMock{}, a, c)

	resp, err := api.SearchPostalCode("0")
	assert.Nil(t, err)
	assert.NotEmpty(t, resp)

	// t.Log(data, "asd")
}

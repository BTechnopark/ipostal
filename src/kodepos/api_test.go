package kodepos_test

import (
	"testing"

	"github.com/BTechnopark/ipostal/src/client"
	"github.com/BTechnopark/ipostal/src/kodepos"
	"github.com/BTechnopark/ipostal/src/session"
	"github.com/stretchr/testify/assert"
)

type configMock struct{}

// GetBaseUrl implements api.PostalConfig.
func (c *configMock) GetBaseUrl() string {
	return "https://kodepos.co.id"
}

func TestKodePost(t *testing.T) {
	sess := session.NewSession("test_session")
	a := client.NewApi(sess)
	api := kodepos.NewKodePosApi(&configMock{}, a)

	t.Run("test get province", func(t *testing.T) {
		provinces, err := api.Province()
		assert.Nil(t, err)
		assert.NotEmpty(t, provinces)

		// raw, err := json.MarshalIndent(provinces, "", "  ")
		// assert.Nil(t, err)
		// t.Log(string(raw))

		t.Run("test get region", func(t *testing.T) {
			p := provinces[0]
			region, err := api.Region(p.Key)
			assert.Nil(t, err)
			assert.NotEmpty(t, region)

			// raw, err := json.MarshalIndent(region, "", "  ")
			// assert.Nil(t, err)
			// t.Log(string(raw))

			t.Run("test get postal code", func(t *testing.T) {
				r := region[0]
				postalCode, err := api.PostalCode(p.Key, r.Key)
				assert.Nil(t, err)
				assert.NotEmpty(t, postalCode)

				// raw, err := json.MarshalIndent(postalCode, "", "  ")
				// assert.Nil(t, err)
				// t.Log(string(raw))
			})
		})
	})

	t.Run("test search postal code", func(t *testing.T) {
		postalCode, err := api.SearchPostalCode("154")
		assert.Nil(t, err)
		assert.NotEmpty(t, postalCode)
	})
}

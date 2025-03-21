package dns

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func Test_loadConfig(t *testing.T) {
	testCases := []struct {
		json      string
		config    ProviderConfig
		shouldErr bool
	}{
		{
			json: `{
  "apiKeySecretRef": {
    "name": "bunny-secret",
    "key": "api-key"
  }
}`,
			config: ProviderConfig{
				ApiKey: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "bunny-secret",
					},
					Key: "api-key",
				},
			},
			shouldErr: false,
		},
		{
			json: `{
  "dummy": }
}`,
			shouldErr: true,
		},
		{
			shouldErr: false,
			config:    ProviderConfig{},
		},
	}

	for _, test := range testCases {
		json := &extapi.JSON{
			Raw: []byte(test.json),
		}
		if test.json == "" {
			json = nil
		}
		config, err := loadConfig(json)
		if err != nil {
			if !test.shouldErr {
				t.Errorf("got error %v where no error was expected", err)
			}
		} else if test.shouldErr {
			t.Errorf("didn't get an error where an error was expected")
		}
		if !reflect.DeepEqual(config, test.config) {
			t.Errorf("Wrong config value: wanted %v got %v", test.config, config)
		}
	}
}

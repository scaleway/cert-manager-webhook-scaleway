package dns

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func Test_loadConfig(t *testing.T) {
	testCases := []struct {
		json   string
		config ProviderConfig
		err    error
	}{
		{
			json: `{
  "accessKeySecretRef": {
    "name": "scaleway-secret",
    "key": "SCW_ACCESS_KEY"
  },
  "secretKeySecretRef": {
    "name": "scaleway-secret",
    "key": "SCW_SECRET_KEY"
  }
}`,
			config: ProviderConfig{
				AccessKey: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "scaleway-secret",
					},
					Key: "SCW_ACCESS_KEY",
				},
				SecretKey: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "scaleway-secret",
					},
					Key: "SCW_SECRET_KEY",
				},
			},
			err: nil,
		},
	}

	for _, test := range testCases {
		config, err := loadConfig(&extapi.JSON{
			Raw: []byte(test.json),
		})
		if err != test.err {
			t.Errorf("Wrong error value: wanted %v got %v", test.err, err)
		}
		if !reflect.DeepEqual(config, test.config) {
			t.Errorf("Wrong config value: wanted %v got %v", test.config, config)
		}
	}
}

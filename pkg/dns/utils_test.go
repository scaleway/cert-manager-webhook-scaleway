package dns

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_loadConfig(t *testing.T) {
	testCases := []struct {
		json      string
		config    ProviderConfig
		shouldErr bool
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

func Test_getDomainAPI(t *testing.T) {
	jsonBothKey := &extapi.JSON{
		Raw: []byte(`{
  "accessKeySecretRef": {
    "name": "scaleway-secret",
    "key": "SCW_ACCESS_KEY"
  },
  "secretKeySecretRef": {
    "name": "scaleway-secret",
    "key": "SCW_SECRET_KEY"
  }
}`),
	}

	testCases := []struct {
		ch         *v1alpha1.ChallengeRequest
		env        map[string]string
		secret     *v1.Secret
		shouldErr  bool
		errMessage string
	}{
		{
			ch:         &v1alpha1.ChallengeRequest{},
			shouldErr:  true,
			errMessage: "failed to initialize scaleway client: scaleway-sdk-go: access key cannot be empty",
		},
		{
			ch: &v1alpha1.ChallengeRequest{},
			env: map[string]string{
				scw.ScwAccessKeyEnv: "SCWXXXXXXXXXXXXXXXXX",
			},
			shouldErr:  true,
			errMessage: "failed to initialize scaleway client: scaleway-sdk-go: secret key cannot be empty",
		},
		{
			ch: &v1alpha1.ChallengeRequest{},
			env: map[string]string{
				scw.ScwAccessKeyEnv: "SCWXXXXXXXXXXXXXXXXX",
				scw.ScwSecretKeyEnv: "66666666-7777-8888-9999-000000000000",
			},
			shouldErr: false,
		},
		{
			ch: &v1alpha1.ChallengeRequest{
				Config:            jsonBothKey,
				ResourceNamespace: "test",
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "scaleway-secret",
					Namespace: "test",
				},
				Data: map[string][]byte{
					scw.ScwAccessKeyEnv: []byte("SCWXXXXXXXXXXXXXXXXX"),
				},
			},
			shouldErr:  true,
			errMessage: "could not get key SCW_SECRET_KEY in secret scaleway-secret",
		},
		{
			ch: &v1alpha1.ChallengeRequest{
				Config:            jsonBothKey,
				ResourceNamespace: "test",
			},
			secret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "scaleway-secret",
					Namespace: "test",
				},
				Data: map[string][]byte{
					scw.ScwSecretKeyEnv: []byte("66666666-7777-8888-9999-000000000000"),
					scw.ScwAccessKeyEnv: []byte("SCWXXXXXXXXXXXXXXXXX"),
				},
			},
			shouldErr: false,
		},
	}

	for _, test := range testCases {
		fakeKubernetesClient := fake.NewSimpleClientset()
		pSolver := &ProviderSolver{
			client: fakeKubernetesClient,
		}

		if test.secret != nil {
			_, err := pSolver.client.CoreV1().Secrets(test.ch.ResourceNamespace).Create(context.Background(), test.secret, metav1.CreateOptions{})
			if err != nil {
				t.Errorf("failed to create kubernetes secret")
			}
		}
		for k, v := range test.env {
			os.Setenv(k, v)
		}
		_, err := pSolver.getDomainAPI(test.ch)
		if err != nil {
			if !test.shouldErr {
				t.Errorf("got error %v where no error was expected", err)
			}
			if err.Error() != test.errMessage {
				t.Errorf("expected error %s, got %s", test.errMessage, err.Error())
			}
		} else if test.shouldErr {
			t.Errorf("didn't get an error where an error was expected with message %s", test.errMessage)
		}
		for k := range test.env {
			os.Unsetenv(k)
		}
	}
}

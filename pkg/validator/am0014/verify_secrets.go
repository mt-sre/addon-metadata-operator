package am0014

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"gopkg.in/yaml.v3"
)

func init() {
	validator.Register(NewVerifySecretParams)
}

const (
	code = 14
	name = "verify_secret_params"
	desc = "Ensure that `secretParameters` present in app-interface's SaaS file match secrets in the addon metadata"
)

type Response struct {
	SecretParameters []SecretParameter `json:"secretParameters,omitempty" yaml:"secretParameters,omitempty"`
}

type SecretParameter struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Secret Secret `json:"secret,omitempty" yaml:"secret,omitempty"`
}

type Secret struct {
	Path  string `json:"path,omitempty" yaml:"path,omitempty"`
	Field string `json:"field,omitempty" yaml:"field,omitempty"`
}

func NewVerifySecretParams(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &VerifySecretParams{
		Base: base,
	}, nil
}

type VerifySecretParams struct {
	*validator.Base
}

func (i *VerifySecretParams) Run(ctx context.Context, mb types.MetaBundle) validator.Result {

	if mb.AddonMeta.Secrets == nil {
		return i.Success()
	}

	response := Response{}

	if value := ctx.Value(testutils.TestEnvKey); value != nil && value == true {
		response.SecretParameters = []SecretParameter{
			{
				Name: "secret-one",
				Secret: Secret{
					Path:  "mtsre/quay/osd-addons/secrets/random-operator-1/secret-one",
					Field: "secret-one.url",
				},
			},
			{
				Name: "secret-two",
				Secret: Secret{
					Path:  "mtsre/quay/osd-addons/secrets/random-operator-1/secret-two",
					Field: "secret-two.metadata",
				},
			},
		}
	} else {
		responseInBytes, err := getDataFromAppInterface("https://gitlab.cee.redhat.com/service/app-interface/-/raw/master/data/services/addons/cicd/ci-int/saas-mt-SelectorSyncSet.yaml")

		if err != nil {
			return i.Fail(fmt.Sprintf("Failed to get data from app-interface: %v", err))
		}

		if err := yaml.Unmarshal(responseInBytes, &response); err != nil {
			return i.Fail(fmt.Sprintf("Failed to unmarshal response from app-interface: %s", err))
		}
	}

	secretMap := make(map[string]SecretParameter)

	for _, sp := range response.SecretParameters {
		secretMap[sp.Name] = sp
	}

	var secretsNotMatched []string

	for _, secret := range *(mb.AddonMeta.Secrets) {
		if val, isPresent := secretMap[secret.Name]; !isPresent || val.Secret.Path != secret.VaultPath {
			secretsNotMatched = append(secretsNotMatched, secret.Name)
		}
	}

	if len(secretsNotMatched) == 0 {
		return i.Success()
	}

	return i.Fail(fmt.Sprintf("following secrets in addon metadata do not match secretParameters in app-interface's SaaS file: %s", secretsNotMatched))

}

func getDataFromAppInterface(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

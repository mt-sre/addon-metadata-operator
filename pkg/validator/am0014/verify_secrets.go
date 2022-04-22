package am0014

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
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
	Secret struct {
		Path  string `json:"path,omitempty" yaml:"path,omitempty"`
		Field string `json:"field,omitempty" yaml:"field,omitempty"`
	}
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

	responseInBytes, err := getDataFromAppInterface("https://gitlab.cee.redhat.com/service/app-interface/-/raw/master/data/services/addons/cicd/ci-int/saas-mt-SelectorSyncSet.yaml")

	if err != nil {
		return i.Fail(fmt.Sprintf("Failed to get data from app-interface: %v", err))
	}

	response := Response{}

	if err := yaml.Unmarshal(responseInBytes, &response); err != nil {
		return i.Fail(fmt.Sprintf("Failed to unmarshal response from app-interface: %s", err))
	}

	secretMap := make(map[string]SecretParameter)

	for i := range response.SecretParameters {
		secretMap[response.SecretParameters[i].Name] = response.SecretParameters[i]
	}

	var secretsNotMatched []string

	for _, secret := range mb.AddonMeta.Secrets {
		if _, isPresent := secretMap[secret.Name]; !isPresent || secretMap[secret.Name].Secret.Path != secret.VaultPath {
			secretsNotMatched = append(secretsNotMatched, secret.Name)
		}
	}

	if len(secretsNotMatched) == 0 {
		return i.Success()
	}

	return i.Fail(fmt.Sprintf("following secrets in addon metadata do not match secretParameters in app-interface's SaaS file: " + " " + strings.Join(secretsNotMatched, ", ")))

}

func getDataFromAppInterface(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

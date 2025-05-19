package azure

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/danjelhysenaj-dev/azure-keyvault-sdk-go/errors"
)

const (
	keyVaultURLFmt          = "https://%s.vault.azure.net"
	secretNotFoundErrMsgFmt = "A secret with name (%s) was not found in the KeyVault (%s)"
)

type (

	// Secret is a struct that represent to be set in the KeyVault.
	Secret struct {
		Name       string    `json:"name"`
		Value      string    `json:"value,omitempty"`
		Expiration time.Time `json:"expiration,omitempty"`
	}

	// AzSecretsClientProvider is an interface that represents the operations available for the Azure KeyVault Client Provider.
	AzSecretsClientProvider interface {
		NewClient(vaultURL string, credential azcore.TokenCredential, options *azsecrets.ClientOptions) (*azsecrets.Client, error)
	}

	// defaultAzSecretsClientProvider is a struct that represents the default Azure KeyVault Client Provider.
	defaultAzSecretsClientProvider struct{}

	// KeyVaultClient is a struct that represents the KeyVault client.
	KeyVaultClient struct {
		ctx  context.Context
		name string
		url  string

		secretsClient         AzKeyVaultSecretsClientOperations
		secretsClientProvider AzSecretsClientProvider

		// support interfaces
		Secret IKeyVaultSecret
	}

	// AzKeyVaultSecretsClientOperations defines the methods available from azure KEyVault for interacting with the SecretClient.
	AzKeyVaultSecretsClientOperations interface {
		SetSecret(ctx context.Context, name string, parameters azsecrets.SetSecretParameters, options *azsecrets.SetSecretOptions) (azsecrets.SetSecretResponse, error)
		GetSecret(ctx context.Context, name string, version string, options *azsecrets.GetSecretOptions) (azsecrets.GetSecretResponse, error)
		DeleteSecret(ctx context.Context, name string, options *azsecrets.DeleteSecretOptions) (azsecrets.DeleteSecretResponse, error)
		NewListSecretPropertiesPager(options *azsecrets.ListSecretPropertiesOptions) *runtime.Pager[azsecrets.ListSecretPropertiesResponse]
	}

	// IKeyVaultSecret defines the methods available for interacting with KeyVault.
	IKeyVaultSecret interface {
		List() ([]Secret, *errors.Error)
		Get(name string) (*Secret, *errors.Error)
		Set(secret Secret) *errors.Error
		Delete(name string) *errors.Error
	}

	// ListSecretPropertiesPager is an interface that represents the operations available for the ListSecretPropertiesPager.
	ListSecretPropertiesPager interface {
		More() bool
		NextPage(ctx context.Context) (azsecrets.ListSecretPropertiesResponse, error)
	}

	// KeyVaultSecretsManager is a struct that implements the IKeyVaultSecret Interface.
	KeyVaultSecretsManager struct {
		kvClient *KeyVaultClient
	}
)

// List all the secrets from the KeyVault.
// This function returns a slice of a secret names and an error if any.
// returns a list of secrets
func (ksm *KeyVaultSecretsManager) List() ([]Secret, error) {
	// create a slice of secrets
	var secrets []Secret

	// list all the secrets from the KeyVault
	pager := ksm.kvClient.secretsClient.NewListSecretPropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ksm.kvClient.ctx)
		if err != nil {
			return nil, checkAzErrResp(err)
		}
		for _, secret := range page.Value {
			secrets = append(secrets, Secret{
				Name:       secret.ID.Name(),
				Expiration: *secret.Attributes.Expires,
			})
		}
	}

	return secrets, nil
}

// Get a secret from the KeyVault
// The secretName is required to be set.
// The function returns the response payload and an error if any.
func (ksm *KeyVaultSecretsManager) Get(name string) (*Secret, *errors.Error) {
	// get the secret from the KeyVault
	resp, getErr := ksm.kvClient.secretsClient.GetSecret(ksm.kvClient.ctx, name, "", nil)
	if getErr != nil {
		err := checkAzErrResp(getErr)

		if err.Status == http.StatusNotFound {
			err.Message = fmt.Sprintf(secretNotFoundErrMsgFmt, name, ksm.kvClient.name)
		}
		return nil, err
	}
	// create secret object from response
	retrievedSecret := &Secret{
		Name:       name,
		Value:      *resp.Value,
		Expiration: *resp.Attributes.Expires,
	}

	// return secret and value
	return retrievedSecret, nil
}

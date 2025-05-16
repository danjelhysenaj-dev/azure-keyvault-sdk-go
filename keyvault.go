package azure

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
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
		List() ([]Secret, error)
		Get(name string) (*Secret, error)
		Set(secret Secret) error
		Delete(name string) error
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

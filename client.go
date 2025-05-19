package azure

import (
	"context"
	"encoding/json"
	goErr "errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/danjelhysenaj-dev/azure-keyvault-sdk-go/errors"
)

type (
	// AzCredentialProvider is an Interface that represent the operations available for the Azure Credential Provider.
	AzCredentialProvider interface {
		NewDefaultAzureCredential(options *azidentity.DefaultAzureCredentialOptions) (*azidentity.DefaultAzureCredential, error)
	}

	// defaultCredentialProvider is a struct that represents the default Azure Credential Provider.
	defaultCredentialProvider struct{}

	// Client represents the azure connection configuration.
	Client struct {
		ctx          context.Context
		cred         *azidentity.DefaultAzureCredential
		credProvider AzCredentialProvider
	}

	// ClientOption to configure API Client
	ClientOption func(*Client)

	// KeyVaultClientOption to configure KeyVault Client
	KeyVaultClientOption func(*KeyVaultClient)

	// AzErrorResponse is a struct that represents the base response with the common fields from Azure KeyVault.
	AzErrorResponse struct {
		Error AzError `json:"error"`
	}

	// AzError is a struct that represents the error object in the Microsoft KeyVault SDK.
	AzError struct {
		Code       string       `json:"code"`
		Message    string       `json:"message"`
		InnerError AzInnerError `json:"Innererror"`
	}

	// AzInnerError is a struct that represents the inner error object in the Microsoft KeyVault SDK.
	AzInnerError struct {
		Code string `json:"code"`
	}
)

// checkAzErrResp checks for the http status codes returned by the Azure Services.
// The method returns an error if:
// - There is an error while creating and making the http request.
// - If the server returns a http error status code.
func checkAzErrResp(err error) *errors.Error {
	if err == nil {
		return nil
	}
	var (
		statusCode int
		errMsg     string
	)

	azRawErr := new(azcore.ResponseError)
	// assert error
	if goErr.As(err, &azRawErr) {
		azErr := new(AzErrorResponse)
		if err := json.NewDecoder(azRawErr.RawResponse.Body).Decode(&azErr); err != nil {
			return errors.InternalServerErrorf("Unable to decode the azure error response: %v", err)
		}
		statusCode = azRawErr.RawResponse.StatusCode
		errMsg = azErr.Error.Message
	} else {
		errMsg = err.Error()
	}
	switch statusCode {
	case http.StatusNotFound:
		return errors.NotFoundError("")
	case http.StatusUnauthorized:
		return errors.UnAuthorizedError(errMsg)
	case http.StatusForbidden:
		return errors.ForbiddenError(errMsg)
	default:
		return errors.InternalServerError(errMsg)
	}
}

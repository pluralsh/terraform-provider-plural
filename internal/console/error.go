package console

import (
	"encoding/json"

	rawclient "github.com/Yamashou/gqlgenc/client"
	"github.com/pkg/errors"
)

func GetErrorResponse(err error, methodName string) error {
	if err == nil {
		return nil
	}

	errResponse := &rawclient.ErrorResponse{}
	newErr := json.Unmarshal([]byte(err.Error()), errResponse)
	if newErr != nil {
		return err
	}

	errList := errors.New(methodName)
	if errResponse.GqlErrors != nil {
		for _, err := range *errResponse.GqlErrors {
			errList = errors.Wrap(errList, err.Message)
		}
		errList = errors.Wrap(errList, "GraphQL error")
	}
	if errResponse.NetworkError != nil {
		errList = errors.Wrap(errList, errResponse.NetworkError.Message)
		errList = errors.Wrap(errList, "Network error")
	}

	return errList
}

package model

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Cloud represents supported providers
type Cloud string

func (c Cloud) String() string {
	return string(c)
}

func (c Cloud) Equals(cloud Cloud) bool {
	return c == cloud
}

const (
	CloudGCP   Cloud = "gcp"
	CloudAWS   Cloud = "aws"
	CloudAzure Cloud = "azure"
	CloudBYOK  Cloud = "byok"
)

var (
	CloudValidator = stringvalidator.OneOfCaseInsensitive(
		CloudBYOK.String(), CloudAWS.String(), CloudAzure.String(), CloudGCP.String(),
	)
)

func IsCloud(c string, cloud Cloud) bool {
	return c == string(cloud)
}

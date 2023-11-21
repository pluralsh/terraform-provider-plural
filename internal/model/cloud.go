package model

// Cloud represents supported providers.
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

func IsCloud(c string, cloud Cloud) bool {
	return c == string(cloud)
}

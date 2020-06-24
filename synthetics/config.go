package synthetics

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/synthetics"
	awsbase "github.com/hashicorp/aws-sdk-go-base"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
)

type Config struct {
	AccessKey     string
	SecretKey     string
	CredsFilename string
	Profile       string
	Token         string
	Region        string
	MaxRetries    int

	AssumeRoleARN         string
	AssumeRoleExternalID  string
	AssumeRoleSessionName string
	AssumeRolePolicy      string

	AllowedAccountIds   []string
	ForbiddenAccountIds []string

	Endpoints         map[string]string
	IgnoreTagPrefixes []string
	IgnoreTags        []string
	Insecure          bool

	SkipCredsValidation     bool
	SkipGetEC2Platforms     bool
	SkipRegionValidation    bool
	SkipRequestingAccountId bool
	SkipMetadataApiCheck    bool
	S3ForcePathStyle        bool

	terraformVersion string
}

type AWSClient struct {
	accountid          string
	dnsSuffix          string
	partition          string
	region             string
	supportedplatforms []string
	terraformVersion   string
	syntheticsconn     *synthetics.Synthetics
}

// PartitionHostname returns a hostname with the provider domain suffix for the partition
// e.g. PREFIX.amazonaws.com
// The prefix should not contain a trailing period.
func (client *AWSClient) PartitionHostname(prefix string) string {
	return fmt.Sprintf("%s.%s", prefix, client.dnsSuffix)
}

// RegionalHostname returns a hostname with the provider domain suffix for the region and partition
// e.g. PREFIX.us-west-2.amazonaws.com
// The prefix should not contain a trailing period.
func (client *AWSClient) RegionalHostname(prefix string) string {
	return fmt.Sprintf("%s.%s.%s", prefix, client.region, client.dnsSuffix)
}

// Client configures and returns a fully initialized AWSClient
func (c *Config) Client() (interface{}, error) {
	// Get the auth and region. This can fail if keys/regions were not
	// specified and we're attempting to use the environment.
	if !c.SkipRegionValidation {
		if err := awsbase.ValidateRegion(c.Region); err != nil {
			return nil, err
		}
	}

	log.Println("[INFO] Building AWS auth structure")
	awsbaseConfig := &awsbase.Config{
		AccessKey:               c.AccessKey,
		AssumeRoleARN:           c.AssumeRoleARN,
		AssumeRoleExternalID:    c.AssumeRoleExternalID,
		AssumeRolePolicy:        c.AssumeRolePolicy,
		AssumeRoleSessionName:   c.AssumeRoleSessionName,
		CredsFilename:           c.CredsFilename,
		DebugLogging:            logging.IsDebugOrHigher(),
		IamEndpoint:             c.Endpoints["iam"],
		Insecure:                c.Insecure,
		MaxRetries:              c.MaxRetries,
		Profile:                 c.Profile,
		Region:                  c.Region,
		SecretKey:               c.SecretKey,
		SkipCredsValidation:     c.SkipCredsValidation,
		SkipMetadataApiCheck:    c.SkipMetadataApiCheck,
		SkipRequestingAccountId: c.SkipRequestingAccountId,
		StsEndpoint:             c.Endpoints["sts"],
		Token:                   c.Token,
		UserAgentProducts: []*awsbase.UserAgentProduct{
			{Name: "APN", Version: "1.0"},
			{Name: "HashiCorp", Version: "1.0"},
			{Name: "Terraform", Version: c.terraformVersion,
				Extra: []string{"+https://www.terraform.io"}},
		},
	}

	sess, accountID, partition, err := awsbase.GetSessionWithAccountIDAndPartition(awsbaseConfig)
	if err != nil {
		return nil, err
	}

	if accountID == "" {
		log.Printf("[WARN] AWS account ID not found for provider. See https://www.terraform.io/docs/providers/aws/index.html#skip_requesting_account_id for implications.")
	}

	dnsSuffix := "amazonaws.com"
	if p, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), c.Region); ok {
		dnsSuffix = p.DNSSuffix()
	}

	client := &AWSClient{
		accountid:        accountID,
		dnsSuffix:        dnsSuffix,
		partition:        partition,
		region:           c.Region,
		terraformVersion: c.terraformVersion,
		syntheticsconn:   synthetics.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["synthetics"])})),
	}
	return client, nil
}

package synthetics

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/arn"
)

const (
	awsAccountIDRegexpPattern = `^(aws|\d{12})$`
	awsPartitionRegexpPattern = `^aws(-[a-z]+)*$`
	awsRegionRegexpPattern    = `^[a-z]{2}(-[a-z]+)+-\d$`
)

var awsAccountIDRegexp = regexp.MustCompile(awsAccountIDRegexpPattern)
var awsPartitionRegexp = regexp.MustCompile(awsPartitionRegexpPattern)
var awsRegionRegexp = regexp.MustCompile(awsRegionRegexpPattern)

func validateArn(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return
	}

	parsedARN, err := arn.Parse(value)

	if err != nil {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: %s", k, value, err))
		return
	}

	if parsedARN.Partition == "" {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: missing partition value", k, value))
	} else if !awsPartitionRegexp.MatchString(parsedARN.Partition) {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: invalid partition value (expecting to match regular expression: %s)", k, value, awsPartitionRegexpPattern))
	}

	if parsedARN.Region != "" && !awsRegionRegexp.MatchString(parsedARN.Region) {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: invalid region value (expecting to match regular expression: %s)", k, value, awsRegionRegexpPattern))
	}

	if parsedARN.AccountID != "" && !awsAccountIDRegexp.MatchString(parsedARN.AccountID) {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: invalid account ID value (expecting to match regular expression: %s)", k, value, awsAccountIDRegexpPattern))
	}

	if parsedARN.Resource == "" {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: missing resource value", k, value))
	}

	return
}

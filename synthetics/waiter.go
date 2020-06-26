package synthetics

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/synthetics"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	// Maximum amount of time to wait for a Canary to return Ready
	CanaryCreatedTimeout = 5 * time.Minute
)

func CanaryStatus(conn *synthetics.Synthetics, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &synthetics.GetCanaryInput{
			Name: aws.String(name),
		}

		output, err := conn.GetCanary(input)

		if err != nil {
			return nil, synthetics.CanaryStateError, err
		}

		if aws.StringValue(output.Canary.Status.State) == synthetics.CanaryStateError {
			return output, synthetics.CanaryStateError, fmt.Errorf("%s: %s", aws.StringValue(output.Canary.Status.StateReasonCode), aws.StringValue(output.Canary.Status.StateReason))
		}

		return output, aws.StringValue(output.Canary.Status.State), nil
	}
}

func CanaryReady(conn *synthetics.Synthetics, name string) (*synthetics.GetCanaryOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{synthetics.CanaryStateCreating, synthetics.CanaryStateUpdating, synthetics.CanaryStateStopping},
		Target:  []string{synthetics.CanaryStateReady, synthetics.CanaryStateRunning, synthetics.CanaryStateStopped},
		Refresh: CanaryStatus(conn, name),
		Timeout: CanaryCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*synthetics.GetCanaryOutput); ok {
		return v, err
	}

	return nil, err
}

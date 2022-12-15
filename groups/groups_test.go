package groups_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/megaproaktiv/awsmock"
	"github.com/megaproaktiv/orphan/groups"
	"gotest.tools/assert"
)


func TestListOrphans(t *testing.T){
	counter := 0
	DescribeLogGroupsFunc := func(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {		
		testfile := "../testdata/loggroups-1-50.json"
		data, err := os.ReadFile(testfile)
		if err != nil {
			fmt.Println("File reading error: ", err)
		}		
		out := &cloudwatchlogs.DescribeLogGroupsOutput{}
		err = json.Unmarshal(data, out); 
		if err != nil {
			t.Error(err)
		}
		if counter == 0 {
			out.NextToken = aws.String("/aws/lambda/testgroup-9-5")
		}
		counter ++
		return out,nil
	}

	ListFunctionsFunc := func(ctx context.Context, params *lambda.ListFunctionsInput) (*lambda.ListFunctionsOutput, error) {
		out := &lambda.ListFunctionsOutput{}
		
		return out,nil
	}

	// Create a Mock Handler
	mockCfg := awsmock.NewAwsMockHandler()
	// add a function to the handler
	// Routing per paramater types
	mockCfg.AddHandler(DescribeLogGroupsFunc)
	mockCfg.AddHandler(ListFunctionsFunc)

	// Create mocking client
	clientLogs := cloudwatchlogs.NewFromConfig(mockCfg.AwsConfig())
	clientLambda := lambda.NewFromConfig(mockCfg.AwsConfig())
	list, err := groups.ListOrphans(clientLambda, clientLogs)
	assert.NilError(t, err, "ListOrphans should return no error")
	assert.Equal(t, len(list), 50, "Length should be 50")
	assert.Equal(t, counter, 2, "Describefunction should be called twice")
}
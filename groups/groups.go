package groups

import (
	"context"
	"log"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

var ClientLambda *lambda.Client
var ClientLogs *cloudwatchlogs.Client

func init(){
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
			panic("configuration error, " + err.Error())
	}
	ClientLambda = lambda.NewFromConfig(cfg)
	ClientLogs = cloudwatchlogs.NewFromConfig(cfg)

}


func DeleteLogGroup( client cloudwatchlogs.Client, logGroupName *string){
	_, err := client.DeleteLogGroup(context.TODO(), 
				&cloudwatchlogs.DeleteLogGroupInput{
					LogGroupName: logGroupName,
				},
			)
			if err != nil{
				log.Println("Problem deleting log group")
				panic(err)
			}
}			

// todo paginate
// Get all log groups with corresponding Lambda Function
func ListOrphans(clientLambda *lambda.Client, clientLogs *cloudwatchlogs.Client) ([] *string, error){
	const maxCount = 1000

	orphans := make([]*string, 0,100)
	lambdaMap := map[string]string{}
	counter := 1
	paramsLambda := lambda.ListFunctionsInput{}
	for {
		resp, err := clientLambda.ListFunctions(context.TODO(), &paramsLambda)
		if err != nil{
			log.Println("Error calling lambda")
			panic(err)
		}
		for _, i := range resp.Functions {
			lambdaMap["/aws/lambda/"+*i.FunctionName] = *i.FunctionName
		}
		if counter > maxCount {
			break
		}
		if isNil(resp.NextMarker) {
			break
		}
		paramsLambda.Marker = resp.NextMarker
	}

	counter = 1
	params := cloudwatchlogs.DescribeLogGroupsInput{}
	for{
		counter++
		// Range logs
		// limit: Valid Range: Minimum value of 1. Maximum value of 50.
		respLogs, err := clientLogs.DescribeLogGroups(context.TODO(),&params)
		if err != nil{
			log.Println("Error calling cloudwatchlogs")
			return nil, err
		}
		// Token is just the name of the last group
		//fmt.Printf("Token: <%v>\n", *respLogs.NextToken)
		for _, i := range respLogs.LogGroups {
			logGroupName := *i.LogGroupName
			if !strings.HasPrefix(logGroupName, "/aws/lambda") {
				continue
			}
			_, isMapContainsKey := lambdaMap[logGroupName]
			if ! isMapContainsKey {
				// fmt.Printf("Found orphan: %v \n", logGroupName)
				// fmt.Printf("Delete log group: %v", logGroupName)
				orphans = append(orphans, &logGroupName)
			}
		}
		if counter > maxCount {
			break
		}
		if isNil(respLogs.NextToken) {
			break
		}
		params.NextToken = respLogs.NextToken
	}
	return orphans, nil
}	

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
 }
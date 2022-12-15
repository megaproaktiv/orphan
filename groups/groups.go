package groups

import (
	"context"
	"log"
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
	orphans := make([]*string, 0,100)
	lambdaMap := map[string]string{}
	resp, err := clientLambda.ListFunctions(context.TODO(), nil)
	if err != nil{
		log.Println("Error calling lambda")
		panic(err)
	}
	for _, i := range resp.Functions {
		lambdaMap["/aws/lambda/"+*i.FunctionName] = *i.FunctionName
	}

	// Range logs
	// limit: Valid Range: Minimum value of 1. Maximum value of 50.
	respLogs, err := clientLogs.DescribeLogGroups(context.TODO(),nil)
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
	return orphans, nil
}	
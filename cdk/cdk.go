package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewAppStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	announceTable := awsdynamodb.NewTableV2(
		stack, jsii.String("AnnounceTable"), &awsdynamodb.TablePropsV2{
			TableName: jsii.Sprintf("%s-AnnounceTable", *stack.StackName()),
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("announce_id"),
				Type: awsdynamodb.AttributeType_STRING,
			},
			SortKey: &awsdynamodb.Attribute{
				Name: jsii.String("published_at"),
				Type: awsdynamodb.AttributeType_STRING,
			},
			RemovalPolicy:      awscdk.RemovalPolicy_DESTROY,
			DeletionProtection: jsii.Bool(false),
		})
	awscdk.NewCfnOutput(
		stack, jsii.String("AnnounceTableName"), &awscdk.CfnOutputProps{
			Value: announceTable.TableName(),
		})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAppStack(app, "GolangDynamoDBSandbox", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}

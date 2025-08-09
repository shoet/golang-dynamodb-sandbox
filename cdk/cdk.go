package main

import (
	"cdk/resources"

	"github.com/aws/aws-cdk-go/awscdk/v2"
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

	announceTable := resources.NewAnnounceTable(stack, resources.AnnounceTableProps{})
	awscdk.NewCfnOutput(
		stack, jsii.String("AnnounceTableName"), &awscdk.CfnOutputProps{
			Value: announceTable.Table.TableName(),
		})
	return stack
}

func NewLocalStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	announceTable := resources.NewAnnounceTable(stack, resources.AnnounceTableProps{})
	awscdk.NewCfnOutput(
		stack, jsii.String("AnnounceTableName"), &awscdk.CfnOutputProps{
			Value: announceTable.Table.TableName(),
		})
	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewLocalStack(app, "GolangDynamoDBSandbox-local", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

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

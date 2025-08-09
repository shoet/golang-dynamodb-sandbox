package resources

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/jsii-runtime-go"
)

type AnnounceTableProps struct{}

type AnnounceTable struct {
	Table awsdynamodb.TableV2
}

func NewAnnounceTable(stack awscdk.Stack, props AnnounceTableProps) *AnnounceTable {
	announceTable := awsdynamodb.NewTableV2(
		stack, jsii.String("AnnounceTable"), &awsdynamodb.TablePropsV2{
			TableName: jsii.Sprintf("%s-AnnounceTable", *stack.StackName()),
			PartitionKey: &awsdynamodb.Attribute{ // uuid
				Name: aws.String("announce_id"),
				Type: awsdynamodb.AttributeType_STRING,
			},
			GlobalSecondaryIndexes: &[]*awsdynamodb.GlobalSecondaryIndexPropsV2{
				{
					IndexName: aws.String("StatusWithPublished"),
					PartitionKey: &awsdynamodb.Attribute{ // "published" | "draft"
						Name: jsii.String("status"),
						Type: awsdynamodb.AttributeType_STRING,
					},
					SortKey: &awsdynamodb.Attribute{ // RFC3999
						Name: jsii.String("published_at"),
						Type: awsdynamodb.AttributeType_STRING,
					},
				},
			},
			RemovalPolicy:      awscdk.RemovalPolicy_DESTROY,
			DeletionProtection: jsii.Bool(false),
		})
	return &AnnounceTable{
		Table: announceTable,
	}
}

func (t *AnnounceTable) GrantReadWrite(grantee ...awsiam.IGrantable) {
	for _, g := range grantee {
		t.GrantReadWrite(g)
	}
}

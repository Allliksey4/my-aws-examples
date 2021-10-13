package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
)

func main() {
	var dryRun = false

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// IAM adding user
	var outAIMUser *iam.CreateUserOutput
	userIAM := iam.New(sess)
	_, err := userIAM.GetUser(&iam.GetUserInput{
		UserName: aws.String("tester"),
	})
	if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == iam.ErrCodeNoSuchEntityException {
		outAIMUser, err := userIAM.CreateUser(&iam.CreateUserInput{
			UserName: &os.Args[1],
		})

		if err != nil {
			fmt.Println("CreateUser Error", err)
			return
		}

		fmt.Println("Success", outAIMUser)
	} else {
		fmt.Println("GetUser Error", err)
	}

	// VPC building
	svc := ec2.New(sess)
	outVpc, err := svc.CreateDefaultVpc(&ec2.CreateDefaultVpcInput{
		DryRun: &dryRun,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// ECS building
	cluster := ecs.New(sess)
	outCluster, err := cluster.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String("ECS Test cluster"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// EC2 building
	outEc2, err := svc.RunInstances(&ec2.RunInstancesInput{
		// Ubuntu 18.04
		ImageId: aws.String("ami-0b1deee75235aa4bb"),
		// Free instance type
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	fmt.Println(outAIMUser)
	fmt.Println(outCluster)
	fmt.Println(outVpc)
	fmt.Println(outEc2)
}

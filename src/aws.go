package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
)

func main() {

	os.Setenv("AWS_PROFILE", "default")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "credentials")
	file := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	os.Setenv("AWS_REGION", "us-east-2")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	fmt.Println("File:", file)
	fmt.Println("Cred:", sess)

	svc := ec2.New(sess)

	fmt.Println("Svc", svc)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}
	resp, _ := svc.DescribeInstances(params)

	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("instance", inst.InstanceId)
		}

	}

}

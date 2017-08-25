package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
)

func main() {

	// os.Setenv("AWS_PROFILE", "oam-devops")

	os.Setenv("AWS_PROFILE", "work")
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "credentials")
	os.Setenv("AWS_CONFIG_FILE", "config")

	// os.Setenv("AWS_REGION", "us-east-2")
	os.Setenv("AWS_REGION", "eu-central-1")

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	fmt.Println("Key:", os.Getenv("AWS_ACCESS_KEY_ID"))

	fmt.Println("File:", os.Getenv("AWS_SHARED_CREDENTIALS_FILE"))
	// fmt.Println("Config:", configfile)

	fmt.Println("CONFIG", os.Getenv("AWS_PROFILE"))

	fmt.Println("Cred:", sess)

	svc := ec2.New(sess)
	// fmt.Println("error", err)

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

	fmt.Println("Instance:", resp)

	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("instance", inst.InstanceId)
		}

	}

}

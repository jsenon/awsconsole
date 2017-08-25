package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/vaughan0/go-ini"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
printIds accepts an aws credentials file and a region, and prints out all
instances within the region in a format that's acceptable to us. Currently that
format is like this:

  instance_id name private_ip instance_type public_ip account

Any values that aren't available (such as public ip) will be printed out as
"None"

Because the "name" parameter is user-defined, we'll run QueryEscape on it so that
our output stays as a space-separated line.
*/
func printIds(creds aws.CredentialsProvider, account string, region string, wg *sync.WaitGroup) {
	defer wg.Done()

	svc := ec2.New(&aws.Config{
		Credentials: creds,
		Region:      region,
	})

	// Here we create an input that will filter any instances that aren't either
	// of these two states. This is generally what we want
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

	// TODO: Actually care if we can't connect to a host
	resp, _ := svc.DescribeInstances(params)
	// if err != nil {
	//      panic(err)
	// }

	// Loop through the instances. They don't always have a name-tag so set it
	// to None if we can't find anything.
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {

			// We need to see if the Name is one of the tags. It's not always
			// present and not required in Ec2.
			name := "None"
			for _, keys := range inst.Tags {
				if *keys.Key == "Name" {
					name = url.QueryEscape(*keys.Value)
				}
			}

			important_vals := []*string{
				inst.InstanceID,
				&name,
				inst.PrivateIPAddress,
				inst.InstanceType,
				inst.PublicIPAddress,
				&account,
			}

			// Convert any nil value to a printable string in case it doesn't
			// doesn't exist, which is the case with certain values
			output_vals := []string{}
			for _, val := range important_vals {
				if val != nil {
					output_vals = append(output_vals, *val)
				} else {
					output_vals = append(output_vals, "None")
				}
			}
			// The values that we care about, in the order we want to print them
			fmt.Println(strings.Join(output_vals, " "))
		}
	}
}

func main() {
	// Go for it!
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Make sure the config file exists
	config := "config"
	if _, err := os.Stat(config); os.IsNotExist(err) {
		fmt.Println("No config file found at: %s", config)
		os.Exit(1)
	}

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "credentials")

	var wg sync.WaitGroup

	file, err := ini.LoadFile(config)
	check(err)

	for key, values := range file {
		profile := strings.Fields(key)

		// Don't print the default or non-standard profiles
		if len(profile) != 2 {
			continue
		}

		// Where to find this host. The account isn't necessary for the creds
		// but it's something we expose to users when we print
		account := profile[1]
		key := values["aws_access_key_id"]
		fmt.Println("key", key)
		pass := values["aws_secret_access_key"]
		creds := aws.Credentials(key, pass, "")

		// Gather a list of all available AWS regions. Even though we're gathering
		// all regions, we still must use a region here for api calls.
		svc := ec2.New(&aws.Config{
			Credentials: creds,
			Region:      "eu-central-1",
		})

		// Iterate over every single stinking region to get a list of available
		// ec2 instances
		regions, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
		check(err)
		for _, region := range regions.Regions {
			wg.Add(1)
			go printIds(creds, account, *region.RegionName, &wg)
		}
	}

	// Allow the goroutines to finish printing
	wg.Wait()
}

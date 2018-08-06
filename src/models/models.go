package models

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Ec2Manager interface {
	StartInstance(ids ...string) error
	StopInstance(ids ...string) error
	DescribeInstances(ids ...string) (*ec2.DescribeInstancesOutput, error)
}

type GraphQLResult struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	LaunchTime string `json:"launchtime"`
	State      string `json:"state"`
}

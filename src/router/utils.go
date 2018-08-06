package router

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/ec2manager/src/models"
	"github.com/gin-gonic/gin"
)

func renderResponse(c *gin.Context, status int, title, message string) {
	c.JSON(status, gin.H{
		title: message,
	})
}

func prepareGraphQLoutput(output *ec2.DescribeInstancesOutput) models.GraphQLResult {
	result := models.GraphQLResult{}
	if output == nil {
		return result
	}

	if len(output.Reservations) <= 0 {
		return result
	}

	if len(output.Reservations[0].Instances) <= 0 {
		return result
	}

	if ID := output.Reservations[0].Instances[0].InstanceId; ID == nil {
		result.ID = "not available"
	} else {
		result.ID = *ID
	}

	result.Type = (string)(output.Reservations[0].Instances[0].InstanceType)

	if launchTime := output.Reservations[0].Instances[0].LaunchTime; launchTime == nil {
		result.LaunchTime = "not available"
	} else {
		result.LaunchTime = launchTime.String()
	}

	if state := output.Reservations[0].Instances[0].State; state == nil {
		result.State = "not available"
	} else {
		result.State = (string)(state.Name)
	}

	return result
}

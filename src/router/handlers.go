package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
)

const (
	ID    = "id"
	Query = "query"

	errorMsg   = "error"
	successMsg = "message"

	contextParam = "output"
)

var (
	instanceType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Instance",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"launchtime": &graphql.Field{
				Type: graphql.String,
			},
			"state": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"instance": &graphql.Field{
				Type: instanceType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if output, ok := p.Context.Value(contextParam).(*ec2.DescribeInstancesOutput); !ok {
						return nil, errors.New("can't cast to value")
					} else {
						return prepareGraphQLoutput(output), nil
					}
				},
			},
		},
	})

	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
)

func (s *server) startInstanceHandler(c *gin.Context) {
	id := c.Query(ID)
	if err := s.Ec2Manager.StartInstance(id); err != nil {
		renderResponse(c, http.StatusForbidden, errorMsg, err.Error())
	} else {
		renderResponse(c, http.StatusOK, successMsg, fmt.Sprintf("Successfully started instance with ID %s", id))
	}
}

func (s *server) stopInstanceHandler(c *gin.Context) {
	id := c.Query(ID)
	if err := s.Ec2Manager.StopInstance(id); err != nil {
		renderResponse(c, http.StatusForbidden, errorMsg, err.Error())
	} else {
		renderResponse(c, http.StatusOK, successMsg, fmt.Sprintf("Successfully stoped instance with ID %s", id))
	}
}

func (s *server) describeInstancesHandler(c *gin.Context) {
	id := c.Query(ID)
	if result, err := s.Ec2Manager.DescribeInstances(id); err != nil {
		renderResponse(c, http.StatusForbidden, errorMsg, err.Error())
	} else {
		c.JSON(http.StatusOK, &result)
	}
}

func (s *server) graphQLHandler(c *gin.Context) {
	id := c.Query(ID)
	output, err := s.Ec2Manager.DescribeInstances(id)
	if err != nil {
		renderResponse(c, http.StatusForbidden, errorMsg, err.Error())
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: c.Query(Query),
		Context:       context.WithValue(context.Background(), contextParam, output),
	})
	c.JSON(http.StatusOK, result)
}

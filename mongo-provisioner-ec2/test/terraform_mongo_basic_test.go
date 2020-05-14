package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestTerraformAwsExample(t *testing.T) {
	if t.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Parallel()
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/mongodb-in-public-subnet",
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"ebs_volume_id": "YOUR_VOLUME_ID",
		},
	}
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)
	mongodbPublicIp := terraform.Output(t, terraformOptions, "mongo_server_ip_address")
	fmt.Println("mongodb public ip: ", mongodbPublicIp)
	mongodbConnectUrl := "mongodb://" + mongodbPublicIp + ":27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbConnectUrl))
	assert.Nil(t, err)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	assert.Nil(t, err)

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	assert.Nil(t, err)
}

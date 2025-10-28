package testimpl

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testConfigsExamplesFolderDefault = "../../examples"
	infraTFVarFileNameDefault        = "test.tfvars"
	errDescribeClusterFmt            = "Failed to describe EKS cluster %s"
)

func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	eksClient := GetAWSEKSClient(t)

	t.Run("TestClusterExists", func(t *testing.T) {
		clusterName := terraform.Output(t, ctx.TerratestTerraformOptions(), "cluster_name")
		output, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
			Name: &clusterName,
		})
		require.NoErrorf(t, err, errDescribeClusterFmt, clusterName)
		require.NotNil(t, output.Cluster)
		assert.Equal(t, clusterName, *output.Cluster.Name)
	})

	t.Run("TestClusterProperties", func(t *testing.T) {
		clusterName := terraform.Output(t, ctx.TerratestTerraformOptions(), "cluster_name")
		output, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
			Name: &clusterName,
		})
		require.NoErrorf(t, err, errDescribeClusterFmt, clusterName)

		cluster := output.Cluster
		require.NotNil(t, cluster)

		assert.Equal(t, "ACTIVE", string(cluster.Status))
		assert.NotEmpty(t, *cluster.Endpoint)

		expectedVersion := terraform.Output(t, ctx.TerratestTerraformOptions(), "cluster_version")
		assert.Equal(t, expectedVersion, *cluster.Version)

		assert.NotNil(t, cluster.ResourcesVpcConfig)
		assert.NotEmpty(t, cluster.ResourcesVpcConfig.VpcId)
		assert.NotEmpty(t, cluster.ResourcesVpcConfig.SubnetIds)

		tags := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "cluster_tags")
		for key, value := range tags {
			assert.Equal(t, value, cluster.Tags[key])
		}
	})

	t.Run("TestClusterRoleArn", func(t *testing.T) {
		clusterName := terraform.Output(t, ctx.TerratestTerraformOptions(), "cluster_name")
		output, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
			Name: &clusterName,
		})
		require.NoErrorf(t, err, errDescribeClusterFmt, clusterName)

		expectedRoleArn := terraform.Output(t, ctx.TerratestTerraformOptions(), "cluster_role_arn")
		assert.Equal(t, expectedRoleArn, *output.Cluster.RoleArn)
	})
}

func TestModuleOutputs(t *testing.T, ctx types.TestContext) {
	terraformOptions := ctx.TerratestTerraformOptions()
	eksClient := GetAWSEKSClient(t)

	// Get cluster name and fetch details from AWS API
	clusterName := terraform.Output(t, terraformOptions, "cluster_name")
	cluster := GetClusterFromAWS(t, eksClient, clusterName)

	t.Run("TestResourceId", func(t *testing.T) {
		resourceId := terraform.Output(t, terraformOptions, "resource_id")
		assert.NotEmpty(t, resourceId)
		// Verify output matches actual cluster name in AWS
		assert.Equal(t, *cluster.Name, resourceId)
	})

	t.Run("TestResourceName", func(t *testing.T) {
		resourceName := terraform.Output(t, terraformOptions, "resource_name")
		assert.NotEmpty(t, resourceName)
		// Verify output matches actual cluster name in AWS
		assert.Equal(t, *cluster.Name, resourceName)
	})

	t.Run("TestClusterName", func(t *testing.T) {
		assert.NotEmpty(t, clusterName)
		// Verify output matches actual cluster name in AWS
		assert.Equal(t, *cluster.Name, clusterName)
	})

	t.Run("TestClusterArn", func(t *testing.T) {
		arn := terraform.Output(t, terraformOptions, "cluster_arn")
		assert.NotEmpty(t, arn)
		assert.Contains(t, arn, "arn:aws:eks:")
		// Verify output matches actual ARN from AWS API
		assert.Equal(t, *cluster.Arn, arn)
	})

	t.Run("TestClusterEndpoint", func(t *testing.T) {
		endpoint := terraform.Output(t, terraformOptions, "cluster_endpoint")
		assert.NotEmpty(t, endpoint)
		assert.Contains(t, endpoint, "https://")
		// Verify output matches actual endpoint from AWS API
		assert.Equal(t, *cluster.Endpoint, endpoint)
	})

	t.Run("TestClusterVersion", func(t *testing.T) {
		version := terraform.Output(t, terraformOptions, "cluster_version")
		assert.NotEmpty(t, version)
		assert.Equal(t, "1.34", version)
		// Verify output matches actual version from AWS API
		assert.Equal(t, *cluster.Version, version)
	})

	t.Run("TestClusterSecurityGroupId", func(t *testing.T) {
		sgId := terraform.Output(t, terraformOptions, "cluster_security_group_id")
		assert.NotEmpty(t, sgId)
		assert.Contains(t, sgId, "sg-")
		// Verify output matches actual security group from AWS API
		assert.Equal(t, *cluster.ResourcesVpcConfig.ClusterSecurityGroupId, sgId)
	})

	t.Run("TestClusterRoleArn", func(t *testing.T) {
		roleArn := terraform.Output(t, terraformOptions, "cluster_role_arn")
		assert.NotEmpty(t, roleArn)
		assert.Contains(t, roleArn, "arn:aws:iam:")
		assert.Contains(t, roleArn, ":role/")
		// Verify output matches actual role ARN from AWS API
		assert.Equal(t, *cluster.RoleArn, roleArn)
	})

	t.Run("TestResourceNamesGenerated", func(t *testing.T) {
		resourceNames := terraform.OutputMap(t, terraformOptions, "resource_names_generated")
		assert.NotEmpty(t, resourceNames)
		assert.Contains(t, resourceNames, "eks_cluster")
		assert.Contains(t, resourceNames, "iam_role")
		assert.Contains(t, resourceNames, "vpc")
		// Verify the generated EKS cluster name matches what's deployed
		assert.Equal(t, resourceNames["eks_cluster"], clusterName)
	})

	t.Run("TestClusterVpcConfiguration", func(t *testing.T) {
		ValidateClusterVpcConfiguration(t, cluster)
	})

	t.Run("TestClusterLogging", func(t *testing.T) {
		ValidateClusterLogging(t, cluster)
	})
}

func GetAWSEKSClient(t *testing.T) *eks.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)

	return eks.NewFromConfig(cfg)
}

// GetClusterFromAWS fetches cluster details from AWS EKS API
func GetClusterFromAWS(t *testing.T, eksClient *eks.Client, clusterName string) *ekstypes.Cluster {
	clusterOutput, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: &clusterName,
	})
	require.NoErrorf(t, err, errDescribeClusterFmt, clusterName)
	require.NotNil(t, clusterOutput.Cluster, "Cluster should not be nil")
	return clusterOutput.Cluster
}

// ValidateClusterVpcConfiguration validates VPC configuration from AWS API
func ValidateClusterVpcConfiguration(t *testing.T, cluster *ekstypes.Cluster) {
	vpcConfig := cluster.ResourcesVpcConfig
	require.NotNil(t, vpcConfig)

	// Verify endpoint access settings
	assert.True(t, vpcConfig.EndpointPrivateAccess, "Private endpoint access should be enabled")
	assert.True(t, vpcConfig.EndpointPublicAccess, "Public endpoint access should be enabled")

	// Verify subnets are configured
	assert.NotEmpty(t, vpcConfig.SubnetIds, "Cluster should have subnets configured")
	assert.GreaterOrEqual(t, len(vpcConfig.SubnetIds), 2, "Cluster should have at least 2 subnets")

	// Verify VPC ID is present
	assert.NotEmpty(t, vpcConfig.VpcId, "Cluster should be associated with a VPC")
}

// ValidateClusterLogging validates logging configuration from AWS API
func ValidateClusterLogging(t *testing.T, cluster *ekstypes.Cluster) {
	logging := cluster.Logging
	require.NotNil(t, logging)

	if len(logging.ClusterLogging) > 0 {
		// Verify that at least API and Audit logs are enabled (as per test.tfvars)
		enabledLogTypes := make(map[string]bool)
		for _, logSetup := range logging.ClusterLogging {
			if logSetup.Enabled != nil && *logSetup.Enabled {
				for _, logType := range logSetup.Types {
					enabledLogTypes[string(logType)] = true
				}
			}
		}

		assert.True(t, enabledLogTypes["api"], "API logging should be enabled")
		assert.True(t, enabledLogTypes["audit"], "Audit logging should be enabled")
	}
}

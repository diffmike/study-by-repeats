package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		account, err := aws.GetCallerIdentity(ctx)
		if err != nil {
			return err
		}

		region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
		if err != nil {
			return err
		}

		role, err := iam.NewRole(ctx, "task-exec-role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Sid": "",
					"Effect": "Allow",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}]
			}`),
		})
		if err != nil {
			return err
		}

		logPolicy, err := iam.NewRolePolicy(ctx, "lambda-log-policy", &iam.RolePolicyArgs{
			Role: role.Name,
			Policy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [{
                    "Effect": "Allow",
                    "Action": [
                        "logs:CreateLogGroup",
                        "logs:CreateLogStream",
                        "logs:PutLogEvents"
                    ],
                    "Resource": "arn:aws:logs:*:*:*"
                }]
            }`),
		})
		if err != nil {
			return err
		}

		networkPolicy, err := iam.NewRolePolicy(ctx, "network-policy", &iam.RolePolicyArgs{
			Role: role.Name,
			Policy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [{
					"Effect": "Allow",
					"Action": [
						"ec2:CreateNetworkInterface",
						"ec2:CreateNetworkInterfacePermission",
						"ec2:DescribeNetworkInterfaces",
						"ec2:DeleteNetworkInterface"
					],
					"Resource": "*"
                }]
            }`),
		})
		if err != nil {
			return err
		}

		c := config.New(ctx, "")

		rdsSg, err := ec2.NewSecurityGroup(ctx, "study-and-repeat-rds", &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(5432),
					ToPort:     pulumi.Int(5432),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		})
		if err != nil {
			return err
		}

		db, err := rds.NewInstance(ctx, "study-and-repeat", &rds.InstanceArgs{
			AllocatedStorage:    pulumi.Int(20),
			Engine:              pulumi.String("postgres"),
			EngineVersion:       pulumi.String("14.4"),
			InstanceClass:       pulumi.String("db.t3.micro"),
			DbName:              pulumi.String("studyAndRepeat"),
			Password:            c.RequireSecret("DB_PASSWORD"),
			SkipFinalSnapshot:   pulumi.Bool(true),
			Username:            pulumi.String("root"),
			PubliclyAccessible:  pulumi.Bool(true),
			VpcSecurityGroupIds: pulumi.StringArray{rdsSg.ID()},
		})
		if err != nil {
			return err
		}

		ctx.Export("DB Endpoint", pulumi.Sprintf("%s", db.Endpoint))

		function, err := lambda.NewFunction(
			ctx,
			"talkToMe",
			&lambda.FunctionArgs{
				Handler: pulumi.String("handler"),
				Role:    role.Arn,
				Runtime: pulumi.String("go1.x"),
				Code:    pulumi.NewFileArchive("./build/handler.zip"),
				Environment: &lambda.FunctionEnvironmentArgs{
					Variables: pulumi.StringMap{
						"TG_TOKEN":    c.RequireSecret("TG_TOKEN"),
						"DB_PASSWORD": c.RequireSecret("DB_PASSWORD"),
						"DB_USER":     db.Username,
						"DB_HOST":     db.Endpoint,
						"DB_NAME":     db.DbName,
					},
				},
			},
			pulumi.DependsOn([]pulumi.Resource{logPolicy, networkPolicy, db}),
		)
		if err != nil {
			return err
		}

		gateway, err := apigateway.NewRestApi(ctx, "BotGateway", &apigateway.RestApiArgs{
			Name:        pulumi.String("BotGateway"),
			Description: pulumi.String("An API Gateway for the Bot function"),
			Policy: pulumi.String(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    },
    {
      "Action": "execute-api:Invoke",
      "Resource": "*",
      "Principal": "*",
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}`)})
		if err != nil {
			return err
		}

		apiresource, err := apigateway.NewResource(ctx, "BotAPI", &apigateway.ResourceArgs{
			RestApi:  gateway.ID(),
			PathPart: pulumi.String("{proxy+}"),
			ParentId: gateway.RootResourceId,
		})
		if err != nil {
			return err
		}

		_, err = apigateway.NewMethod(ctx, "AnyMethod", &apigateway.MethodArgs{
			HttpMethod:    pulumi.String("ANY"),
			Authorization: pulumi.String("NONE"),
			RestApi:       gateway.ID(),
			ResourceId:    apiresource.ID(),
		})
		if err != nil {
			return err
		}

		_, err = apigateway.NewIntegration(ctx, "LambdaIntegration", &apigateway.IntegrationArgs{
			HttpMethod:            pulumi.String("ANY"),
			IntegrationHttpMethod: pulumi.String("POST"),
			ResourceId:            apiresource.ID(),
			RestApi:               gateway.ID(),
			Type:                  pulumi.String("AWS_PROXY"),
			Uri:                   function.InvokeArn,
		})
		if err != nil {
			return err
		}

		permission, err := lambda.NewPermission(ctx, "APIPermission", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  function.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, gateway.ID()),
		}, pulumi.DependsOn([]pulumi.Resource{apiresource}))
		if err != nil {
			return err
		}

		_, err = apigateway.NewDeployment(ctx, "APIDeployment", &apigateway.DeploymentArgs{
			Description:      pulumi.String("Lambda handler"),
			RestApi:          gateway.ID(),
			StageDescription: pulumi.String("Development"),
			StageName:        pulumi.String("dev"),
		}, pulumi.DependsOn([]pulumi.Resource{apiresource, function, permission}))
		if err != nil {
			return err
		}

		ctx.Export("Invocation URL", pulumi.Sprintf("https://%s.execute-api.%s.amazonaws.com/dev/{message}", gateway.ID(), region.Name))

		return nil
	})
}

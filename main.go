package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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

		args := &lambda.FunctionArgs{
			Handler: pulumi.String("handler"),
			Role:    role.Arn,
			Runtime: pulumi.String("go1.x"),
			Code:    pulumi.NewFileArchive("./handler/handler.zip"),
		}

		function, err := lambda.NewFunction(
			ctx,
			"talkToMe",
			args,
			pulumi.DependsOn([]pulumi.Resource{logPolicy}),
		)
		if err != nil {
			return err
		}

		// Create a new API Gateway.
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

		ctx.Export("invocation URL", pulumi.Sprintf("https://%s.execute-api.%s.amazonaws.com/dev/{message}", gateway.ID(), region.Name))

		//db, err := rds.NewInstance(ctx, "database", &rds.InstanceArgs{
		//	AllocatedStorage:   pulumi.Int(10),
		//	Engine:             pulumi.String("postgree"),
		//	EngineVersion:      pulumi.String("5.7"),
		//	InstanceClass:      pulumi.String("db.t3.micro"),
		//	DbName:             pulumi.String("mydb"),
		//	ParameterGroupName: pulumi.String("default.mysql5.7"),
		//	Password:           pulumi.String("foobarbaz"),
		//	SkipFinalSnapshot:  pulumi.Bool(true),
		//	Username:           pulumi.String("admin"),
		//	PubliclyAccessible: pulumi.Bool(true),
		//})
		//if err != nil {
		//	return err
		//}
		//
		//ctx.Export("DB Endpoint", pulumi.Sprintf("%s", db.Endpoint))

		return nil
	})
}

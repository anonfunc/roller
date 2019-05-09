// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// *** SAM Build ***

const ServiceName = "roller"

// Borrowing STAGE concept from serverless framework as a unique suffix of the
// cloudformation stack.  No correlation to API GW stage name, which appears
// to always be "Prod".
var Stage = os.Getenv("STAGE")
var Region = endpoints.UsWest2RegionID
var BucketName = os.Getenv("SAM_BUCKET") // Overrode later
var StackName = ServiceName + "-" + Stage

func DeploySAM() error {
	mg.Deps(BuildLinux)
	mg.SerialDeps(makeBucket, samPackage)

	cmd := exec.Command("sam", "deploy", "--template-file", "packaged.yml",
		"--stack-name", StackName, "--capabilities", "CAPABILITY_IAM")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("aws", "cloudformation", "describe-stacks", "--stack-name", StackName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RemoveSAM() error {
	cmd := exec.Command("aws", "cloudformation", "delete-stack", "--stack-name", StackName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunSAM() error {
	fmt.Println("http://localhost:3000/roll/2d6")
	cmd := exec.Command("sam", "local", "start-api")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func makeBucket() error {
	sess := awsSession()
	svcS3 := s3.New(sess)
	bucketName := bucketName()
	_, err := svcS3.HeadBucket(&s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		_, err := svcS3.CreateBucket(&s3.CreateBucketInput{Bucket: &bucketName})
		return err
	}
	return err
}

func samPackage() error {
	cmd := exec.Command("sam", "package",
		"--output-template-file", "packaged.yml",
		"--s3-bucket", bucketName())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func bucketName() string {
	if len(BucketName) != 0 {
		return BucketName
	}
	sess := awsSession()
	svcSts := sts.New(sess)
	output, err := svcSts.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		panic("could not get caller identity: " + err.Error())
	}
	BucketName = strings.ToLower(fmt.Sprintf("roller-%s-%s",
		Stage, *output.Account))
	fmt.Printf("Deployment bucket name is %s\n", BucketName)
	return BucketName
}

func awsSession() *session.Session {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: &Region,
		},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic("could not get AWS session: " + err.Error())
	}
	return sess
}

// *** GCP Cloud Function

func CloudFunction() error {
	return sh.Run("gcloud", "functions", "deploy","DiceRoller",
		"--trigger-http", "--runtime", "go111", "--source", "fn")
}


func RemoveCloudFunction() error {
	return sh.Run("gcloud", "functions", "delete","DiceRoller")
}

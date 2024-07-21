// Package main
package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"

	// _ "github.com/aws/aws-sdk-go-v2/aws"
	// _ "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	// _ "github.com/aws/aws-sdk-go/aws/awserr"
	// _ "github.com/aws/aws-sdk-go/aws/session"
	"github.com/ryotaro612/dpcli/internal"
	// _ "github.com/slack-go/slack"
	// _ "github.com/urfave/cli/v2" // imports as package "cli"
	// _ "google.golang.org/api/calendar/v3"

	_ "log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func main() {
	err := report(os.Args)
	if err != nil {
		internal.MakeLogger(nil).Error("Failure", "error", err)

	}
}

func report(args []string) error {
	// https://cli.urfave.org/
	app := &cli.App{
		Name:                 "dpreport",
		EnableBashCompletion: true,
		//ArgsUsage:            "doge",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "awsprofile",
				Aliases: []string{"p"},
				Usage:   "The AWS profile to get credentials.",
			},
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "A Go template file. Annotations in the template refer to activities like meeting and pull requests.",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Prints more information.",
			},
		},
		Action: func(ctx *cli.Context) error {
			aws := ctx.String("awsprofile")
			template := ctx.String("template")
			verbose := ctx.Bool("verbose")
			reporting := internal.MakeReporting(ctx.Context, aws, verbose, template)
			report, err := reporting.Report(ctx.Context)
			if err != nil {
				return err
			}
			fmt.Println(report)
			return nil
		},
	}
	if err := app.Run(args); err != nil {
		return err
	}
	return nil

}

func main2() {
	ctx := context.Background()
	svc, err := internal.Load(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	secretName := "secret"

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	//fmt.Printf("svc: %v, err: %v\n", svc, err)

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	//var secretString string = *result.SecretString
	// JSONになる
	fmt.Printf("%v\n", *result.SecretString)
	fmt.Printf("OK")
}

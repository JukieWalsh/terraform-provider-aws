// +build ignore

package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
)

const filename = `service_tags_gen.go`

// Representing types such as []*athena.Tag, []*ec2.Tag, ...
var sliceServiceNames = []string{
	"acm",
	"acmpca",
	"appmesh",
	"athena",
	/* "autoscaling", // includes extra PropagateAtLaunch, skip for now */
	"cloudformation",
	"cloudfront",
	"cloudhsmv2",
	"cloudtrail",
	"cloudwatch",
	"cloudwatchevents",
	"codebuild",
	"codedeploy",
	"codepipeline",
	"configservice",
	"databasemigrationservice",
	"datapipeline",
	"datasync",
	"dax",
	"devicefarm",
	"directconnect",
	"directoryservice",
	"dlm",
	"docdb",
	"dynamodb",
	"ec2",
	"ecr",
	"ecs",
	"efs",
	"elasticache",
	"elasticbeanstalk",
	"elasticsearchservice",
	"elb",
	"elbv2",
	"emr",
	"firehose",
	"fms",
	"fsx",
	"iam",
	"inspector",
	"iot",
	"iotanalytics",
	"iotevents",
	"kinesis",
	"kinesisanalytics",
	"kinesisanalyticsv2",
	"kms",
	"licensemanager",
	"lightsail",
	"mediastore",
	"neptune",
	"organizations",
	"ram",
	"rds",
	"redshift",
	"route53",
	"route53resolver",
	"s3",
	"sagemaker",
	"secretsmanager",
	"serverlessapplicationrepository",
	"servicecatalog",
	"sfn",
	"sns",
	"ssm",
	"storagegateway",
	"swf",
	"transfer",
	"waf",
	"workspaces",
}

var mapServiceNames = []string{
	"amplify",
	"apigateway",
	"apigatewayv2",
	"appstream",
	"appsync",
	"backup",
	"batch",
	"cloudwatchlogs",
	"codecommit",
	"cognitoidentity",
	"cognitoidentityprovider",
	"glacier",
	"glue",
	"guardduty",
	"kafka",
	"lambda",
	"mediaconnect",
	"mediaconvert",
	"medialive",
	"mediapackage",
	"mq",
	"opsworks",
	"pinpoint",
	"resourcegroups",
	"securityhub",
	"sqs",
}

type TemplateData struct {
	MapServiceNames   []string
	SliceServiceNames []string
}

func main() {
	// Always sort to reduce any potential generation churn
	sort.Strings(mapServiceNames)
	sort.Strings(sliceServiceNames)

	templateData := TemplateData{
		MapServiceNames:   mapServiceNames,
		SliceServiceNames: sliceServiceNames,
	}
	templateFuncMap := template.FuncMap{
		"TagType":           ServiceTagType,
		"TagTypeKeyField":   ServiceTagTypeKeyField,
		"TagTypeValueField": ServiceTagTypeValueField,
		"Title":             strings.Title,
	}

	tmpl, err := template.New("servicetags").Funcs(templateFuncMap).Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by generators/servicetags/main.go; DO NOT EDIT.

package keyvaluetags

import (
	"github.com/aws/aws-sdk-go/aws"
{{- range .SliceServiceNames }}
	"github.com/aws/aws-sdk-go/service/{{ . }}"
{{- end }}
)

// map[string]*string handling
{{- range .MapServiceNames }}

// {{ . | Title }}Tags returns {{ . }} service tags.
func (tags KeyValueTags) {{ . | Title }}Tags() map[string]*string {
	return aws.StringMap(tags.Map())
}

// {{ . | Title }}KeyValueTags creates KeyValueTags from {{ . }} service tags.
func {{ . | Title }}KeyValueTags(tags map[string]*string) KeyValueTags {
	return New(tags)
}
{{- end }}

// []*SERVICE.Tag handling
{{- range .SliceServiceNames }}

// {{ . | Title }}Tags returns {{ . }} service tags.
func (tags KeyValueTags) {{ . | Title }}Tags() []*{{ . }}.{{ . | TagType }} {
	result := make([]*{{ . }}.{{ . | TagType }}, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &{{ . }}.{{ . | TagType }}{
			{{ . | TagTypeKeyField }}:   aws.String(k),
			{{ . | TagTypeValueField }}: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// {{ . | Title }}KeyValueTags creates KeyValueTags from {{ . }} service tags.
func {{ . | Title }}KeyValueTags(tags []*{{ . }}.{{ . | TagType }}) KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.{{ . | TagTypeKeyField }})] = tag.{{ . | TagTypeValueField }}
	}

	return New(m)
}
{{- end }}
`

// ServiceTagType determines the service tagging tag type.
func ServiceTagType(serviceName string) string {
	switch serviceName {
	case "appmesh":
		return "TagRef"
	case "datasync":
		return "TagListEntry"
	case "fms":
		return "ResourceTag"
	case "swf":
		return "ResourceTag"
	default:
		return "Tag"
	}
}

// ServiceTagTypeKeyField determines the service tagging tag type key field.
func ServiceTagTypeKeyField(serviceName string) string {
	switch serviceName {
	case "kms":
		return "TagKey"
	default:
		return "Key"
	}
}

// ServiceTagTypeValueField determines the service tagging tag type value field.
func ServiceTagTypeValueField(serviceName string) string {
	switch serviceName {
	case "kms":
		return "TagValue"
	default:
		return "Value"
	}
}

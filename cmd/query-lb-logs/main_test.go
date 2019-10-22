package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestGenerateCustomCloudWatchLink(t *testing.T) {
	region := "us-west-2"
	start := "2019-06-19"
	end := "2019-06-21"
	clientIP := "127.0.0.1"
	urlString := fmt.Sprintf("https://%s.console.aws.amazon.com/cloudwatch/home?region=%s#logs-insights:queryDetail=~(end~'%s~start~'%s~timeType~'ABSOLUTE~tz~'Local~editorString~'fields*20*40timestamp*2c*20*40message*0a*7c*20filter*20*60x-forwarded-for*60*20*3d*20'%s'~isLiveTail~false~source~(~'ecs-tasks-app-prod))~",
		region, region, end, start, clientIP)

	expectedURL, _ := url.Parse(urlString)
	expected := expectedURL.String()

	result := generateCustomCloudWatchLink(region, end, start, clientIP, "prod")

	if result != expected {
		t.Errorf(`Result doesn't match expected.
			Expected:   %s
			Result: 	%s`,
			expected, result)
	}
}

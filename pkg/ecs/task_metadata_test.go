package ecs

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// body copied from https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v2.html
var body = `
{
  "Cluster": "default",
  "TaskARN": "arn:aws:ecs:us-east-2:012345678910:task/9781c248-0edd-4cdb-9a93-f63cb662a5d3",
  "Family": "nginx",
  "Revision": "5",
  "DesiredStatus": "RUNNING",
  "KnownStatus": "RUNNING",
  "Containers": [
    {
      "DockerId": "731a0d6a3b4210e2448339bc7015aaa79bfe4fa256384f4102db86ef94cbbc4c",
      "Name": "~internal~ecs~pause",
      "DockerName": "ecs-nginx-5-internalecspause-acc699c0cbf2d6d11700",
      "Image": "amazon/amazon-ecs-pause:0.1.0",
      "ImageID": "",
      "Labels": {
        "com.amazonaws.ecs.cluster": "default",
        "com.amazonaws.ecs.container-name": "~internal~ecs~pause",
        "com.amazonaws.ecs.task-arn": "arn:aws:ecs:us-east-2:012345678910:task/9781c248-0edd-4cdb-9a93-f63cb662a5d3",
        "com.amazonaws.ecs.task-definition-family": "nginx",
        "com.amazonaws.ecs.task-definition-version": "5"
      },
      "DesiredStatus": "RESOURCES_PROVISIONED",
      "KnownStatus": "RESOURCES_PROVISIONED",
      "Limits": {
        "CPU": 0,
        "Memory": 0
      },
      "CreatedAt": "2018-02-01T20:55:08.366329616Z",
      "StartedAt": "2018-02-01T20:55:09.058354915Z",
      "Type": "CNI_PAUSE",
      "Networks": [
        {
          "NetworkMode": "awsvpc",
          "IPv4Addresses": [
            "10.0.2.106"
          ]
        }
      ]
    },
    {
      "DockerId": "43481a6ce4842eec8fe72fc28500c6b52edcc0917f105b83379f88cac1ff3946",
      "Name": "nginx-curl",
      "DockerName": "ecs-nginx-5-nginx-curl-ccccb9f49db0dfe0d901",
      "Image": "nrdlngr/nginx-curl",
      "ImageID": "sha256:2e00ae64383cfc865ba0a2ba37f61b50a120d2d9378559dcd458dc0de47bc165",
      "Labels": {
        "com.amazonaws.ecs.cluster": "default",
        "com.amazonaws.ecs.container-name": "nginx-curl",
        "com.amazonaws.ecs.task-arn": "arn:aws:ecs:us-east-2:012345678910:task/9781c248-0edd-4cdb-9a93-f63cb662a5d3",
        "com.amazonaws.ecs.task-definition-family": "nginx",
        "com.amazonaws.ecs.task-definition-version": "5"
      },
      "DesiredStatus": "RUNNING",
      "KnownStatus": "RUNNING",
      "Limits": {
        "CPU": 512,
        "Memory": 512
      },
      "CreatedAt": "2018-02-01T20:55:10.554941919Z",
      "StartedAt": "2018-02-01T20:55:11.064236631Z",
      "Type": "NORMAL",
      "Networks": [
        {
          "NetworkMode": "awsvpc",
          "IPv4Addresses": [
            "10.0.2.106"
          ]
        }
      ]
    }
  ],
  "PullStartedAt": "2018-02-01T20:55:09.372495529Z",
  "PullStoppedAt": "2018-02-01T20:55:10.552018345Z",
  "AvailabilityZone": "us-east-2b"
}
`

func TestUnmarshalTaskMetadata(t *testing.T) {

	taskMetadata := &TaskMetadata{}
	err := json.Unmarshal([]byte(body), taskMetadata)
	assert.Nil(t, err)

	assert.Equal(t, taskMetadata.Cluster, "default")
	assert.Equal(t, taskMetadata.Family, "nginx")
	assert.Equal(t, taskMetadata.Revision, "5")

}

// copied from
// https://docs.aws.amazon.com/AmazonECS/latest/userguide/task-metadata-endpoint-v4-fargate.html
const bodyV4 = `
{
    "Cluster": "arn:aws:ecs:us-west-2:111122223333:cluster/default",
    "TaskARN": "arn:aws:ecs:us-west-2:111122223333:task/default/e9028f8d5d8e4f258373e7b93ce9a3c3",
    "Family": "curltest",
    "Revision": "3",
    "DesiredStatus": "RUNNING",
    "KnownStatus": "RUNNING",
    "Limits": {
        "CPU": 0.25,
        "Memory": 512
    },
    "PullStartedAt": "2020-10-08T20:47:16.053330955Z",
    "PullStoppedAt": "2020-10-08T20:47:19.592684631Z",
    "AvailabilityZone": "us-west-2a",
    "Containers": [
        {
            "DockerId": "e9028f8d5d8e4f258373e7b93ce9a3c3-2495160603",
            "Name": "curl",
            "DockerName": "curl",
            "Image": "111122223333.dkr.ecr.us-west-2.amazonaws.com/curltest:latest",
            "ImageID": "sha256:25f3695bedfb454a50f12d127839a68ad3caf91e451c1da073db34c542c4d2cb",
            "Labels": {
                "com.amazonaws.ecs.cluster": "arn:aws:ecs:us-west-2:111122223333:cluster/default",
                "com.amazonaws.ecs.container-name": "curl",
                "com.amazonaws.ecs.task-arn": "arn:aws:ecs:us-west-2:111122223333:task/default/e9028f8d5d8e4f258373e7b93ce9a3c3",
                "com.amazonaws.ecs.task-definition-family": "curltest",
                "com.amazonaws.ecs.task-definition-version": "3"
            },
            "DesiredStatus": "RUNNING",
            "KnownStatus": "RUNNING",
            "Limits": {
                "CPU": 10,
                "Memory": 128
            },
            "CreatedAt": "2020-10-08T20:47:20.567813946Z",
            "StartedAt": "2020-10-08T20:47:20.567813946Z",
            "Type": "NORMAL",
            "Networks": [
                {
                    "NetworkMode": "awsvpc",
                    "IPv4Addresses": [
                        "192.0.2.3"
                    ],
                    "IPv6Addresses": [
                        "2001:dB8:10b:1a00:32bf:a372:d80f:e958"
                    ],
                    "AttachmentIndex": 0,
                    "MACAddress": "02:b7:20:19:72:39",
                    "IPv4SubnetCIDRBlock": "192.0.2.0/24",
                    "IPv6SubnetCIDRBlock": "2600:1f13:10b:1a00::/64",
                    "DomainNameServers": [
                        "192.0.2.2"
                    ],
                    "DomainNameSearchList": [
                        "us-west-2.compute.internal"
                    ],
                    "PrivateDNSName": "ip-172-31-30-173.us-west-2.compute.internal",
                    "SubnetGatewayIpv4Address": "192.0.2.0/24"
                }
            ],
            "ClockDrift": {
                "ClockErrorBound": 0.5458234999999999,
                "ReferenceTimestamp": "2021-09-07T16:57:44Z",
                "ClockSynchronizationStatus": "SYNCHRONIZED"
             },
            "ContainerARN": "arn:aws:ecs:us-west-2:111122223333:container/1bdcca8b-f905-4ee6-885c-4064cb70f6e6",
            "LogOptions": {
                "awslogs-create-group": "true",
                "awslogs-group": "/ecs/containerlogs",
                "awslogs-region": "us-west-2",
                "awslogs-stream": "ecs/curl/e9028f8d5d8e4f258373e7b93ce9a3c3"
            },
            "LogDriver": "awslogs"
        }
    ],
    "LaunchType": "FARGATE"
}
`

func TestUnmarshalTaskMetadataV4(t *testing.T) {

	taskMetadata := &TaskMetadataV4{}
	err := json.Unmarshal([]byte(bodyV4), taskMetadata)
	assert.Nil(t, err)

	assert.Equal(t, "arn:aws:ecs:us-west-2:111122223333:cluster/default", taskMetadata.Cluster)
	assert.Equal(t, "curltest", taskMetadata.Family)
	assert.Equal(t, "3", taskMetadata.Revision)
	assert.Equal(t, "arn:aws:ecs:us-west-2:111122223333:task/default/e9028f8d5d8e4f258373e7b93ce9a3c3", taskMetadata.TaskARN)
	assert.Equal(t, "us-west-2a", taskMetadata.AvailabilityZone)
}

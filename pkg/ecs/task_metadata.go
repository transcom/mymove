package ecs

// TaskMetadata is used for parsing AWS ECS task metadata.
//
//   - https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v2.html
type TaskMetadata struct {
	Cluster  string `json:"cluster"`
	Family   string `json:"family"`
	Revision string `json:"revision"`
}

// TaskMetadataV4 is used for parsing AWS ECS task metadata.
//
//   - https://docs.aws.amazon.com/AmazonECS/latest/userguide/task-metadata-endpoint-v4-fargate.html
type TaskMetadataV4 struct {
	Cluster          string `json:"Cluster"`
	Family           string `json:"Family"`
	Revision         string `json:"Revision"`
	TaskARN          string `json:"TaskARN"`
	AvailabilityZone string `json:"AvailabilityZone"`
}

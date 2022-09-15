package ecs

// TaskMetadata is used for parsing AWS ECS task metadata.
//
//	- https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v2.html
type TaskMetadata struct {
	Cluster  string `json:"cluster"`
	Family   string `json:"family"`
	Revision string `json:"revision"`
}

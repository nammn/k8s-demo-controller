package common

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// RelayEvent indicate the informerEvent
type RelayEvent struct {
	Key            string
	Reason         string
	Message        string
	FirstTimestamp v1.Time
	LastTimestamp  v1.Time
	EventType      string
	Namespace      string
	ResourceType   string
}

type BackendTypes string

const (
	Local    BackendTypes = "local"
	Cloudant BackendTypes = "cloudant"
	Aurora   BackendTypes = "aurora"
)

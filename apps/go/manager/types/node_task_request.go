package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ------------------------------------------------------------------------------
// TaskRequestRecord
// ------------------------------------------------------------------------------

type RequesterArgs struct {
	Address string `json:"address"`
	Service string `json:"service"`
	Method  string `json:"method"`
	Path    string `json:"path"`
}

// A pending request already processed by the Sampler
type TaskRequestRecord struct {
	Id             primitive.ObjectID `bson:"_id"`
	RequesterArgs  RequesterArgs      `bson:"requester_args"`
	Framework      string             `bson:"framework"`
	Task           string             `bson:"task"`
	Blacklist      []int              `bson:"blacklist"`
	Qty            int                `bson:"qty"`
	TotalInstances int                `bson:"total_instances"`
	RequestType    string             `bson:"string"`
	Done           bool               `bson:"done"`
}

// ------------------------------------------------------------------------------
// InstanceRecord
// ------------------------------------------------------------------------------
type InstanceRecord struct {
	TaskID primitive.ObjectID `bson:"task_id"`
	DocIds []int              `bson:"doc_id"`
}
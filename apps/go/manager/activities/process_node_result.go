package activities

import (
	"context"
	"fmt"
	"manager/records"
	"manager/types"
	"packages/mongodb"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.temporal.io/sdk/temporal"
)

var AnalyzeResultName = "analyze_result"

func (aCtx *Ctx) AnalyzeResult(ctx context.Context, params types.AnalyzeResultParams) (*types.AnalyzeResultResults, error) {

	var result types.AnalyzeResultResults
	result.Success = false

	// Get logger
	l := aCtx.App.Logger
	l.Debug().
		Str("task_id", params.TaskID.String()).
		Msg("Analyzing task.")

	// Get results collection
	resultsCollection := aCtx.App.Mongodb.GetCollection(types.ResultsCollection)

	//------------------------------------------------------------------
	// Get Task data
	//------------------------------------------------------------------
	taskData, err := retrieveTaskData(params.TaskID, aCtx.App.Mongodb, l)
	if err != nil {
		return nil, err
	}
	// Extract data
	Node := types.NodeData{
		Address: taskData.RequesterArgs.Address,
		Service: taskData.RequesterArgs.Service,
	}

	l.Debug().
		Str("task_id", params.TaskID.String()).
		Str("address", Node.Address).
		Str("service", Node.Service).
		Str("framework", taskData.Framework).
		Str("task", taskData.Task).
		Msg("Analyzing result.")

	//------------------------------------------------------------------
	// Get stored data for this node
	//------------------------------------------------------------------
	var nodeData records.NodeRecord
	found, err := nodeData.FindAndLoadNode(Node, aCtx.App.Mongodb, l)
	if err != nil {
		return nil, err
	}

	if !found {
		err = temporal.NewApplicationErrorWithCause("unable to get node data", "FindAndLoadNode", fmt.Errorf("Node %s not found", Node.Address))
		l.Error().
			Str("address", Node.Address).
			Msg("Cannot retrieve node data")
		return nil, err
	}

	//------------------------------------------------------------------
	// Get stored data for this task
	//------------------------------------------------------------------
	taskType, err := records.GetTaskType(taskData.Framework, taskData.Task, aCtx.App.Config.Frameworks, l)
	if err != nil {
		return nil, err
	}
	thisTaskRecord, found := records.GetTaskData(nodeData.ID, taskType, taskData.Framework, taskData.Task, aCtx.App.Mongodb, l)

	if !found {
		err = temporal.NewApplicationErrorWithCause("unable to get task data", "GetTaskData", fmt.Errorf("Task %s not found", taskData.Task))
		l.Error().
			Str("address", nodeData.Address).
			Str("service", nodeData.Service).
			Str("framework", taskData.Framework).
			Str("task", taskData.Task).
			Msg("Requested task was not found.")
		return nil, err
	}

	thisTaskResults := thisTaskRecord.GetResultStruct()
	found = false
	found, err = thisTaskResults.FindAndLoadResults(params.TaskID,
		resultsCollection,
		l)
	if err != nil {
		return nil, err
	}
	if !found {
		l.Error().
			Str("address", nodeData.Address).
			Str("service", nodeData.Service).
			Str("framework", taskData.Framework).
			Str("task", taskData.Task).
			Msg("Requested result was not found.")
	}

	l.Debug().
		Str("address", nodeData.Address).
		Str("service", nodeData.Service).
		Str("framework", taskData.Framework).
		Str("task", taskData.Task).
		Str("task_id", params.TaskID.String()).
		Msg("Processing found results.")

	// If nothing is wrong with the result calculation
	if thisTaskResults.GetStatus() == 0 {
		l.Debug().
			Int("NumSamples", int(thisTaskResults.GetNumSamples())).
			Str("address", nodeData.Address).
			Str("service", nodeData.Service).
			Str("framework", taskData.Framework).
			Str("task", taskData.Task).
			Str("task_id", params.TaskID.String()).
			Msg("Inserting results into buffers.")
		// Add results to current task record
		for i := 0; i < int(thisTaskResults.GetNumSamples()); i++ {
			thisTaskRecord.InsertSample(time.Now(), thisTaskResults.GetSample(i), l)
		}
		// Update the last seen fields
		thisTaskRecord.UpdateLastHeight(thisTaskResults.GetResultHeight())
		thisTaskRecord.UpdateLastSeen(thisTaskResults.GetResultTime())
	} else {
		// TODO: handle status!=0
		l.Debug().
			Str("address", nodeData.Address).
			Str("service", nodeData.Service).
			Str("framework", taskData.Framework).
			Str("task", taskData.Task).
			Str("task_id", params.TaskID.String()).
			Msg("Status not zero.")
	}

	// Delete all MongoDB entries associated with this task ID
	errDel := RemoveTaskID(params.TaskID, aCtx.App.Mongodb, l)
	if errDel != nil {
		l.Debug().
			Str("delete_error", errDel.Error()).
			Str("task_id", params.TaskID.String()).
			Msg("Deletion error.")
	}

	//------------------------------------------------------------------
	// Calculate new metrics for this task
	//------------------------------------------------------------------
	thisTaskRecord.ProcessData(l)

	//------------------------------------------------------------------
	// Update task in DB
	//------------------------------------------------------------------

	_, err = thisTaskRecord.UpdateTask(nodeData.ID, taskData.Framework, taskData.Task, aCtx.App.Mongodb, l)
	if err != nil {
		return nil, err
	}

	result.Success = true

	return &result, nil
}

// Looks for an specific task in the TaskDB and retrieves all data
func retrieveTaskData(taskID primitive.ObjectID,
	mongoDB mongodb.MongoDb,
	l *zerolog.Logger) (tasksData types.TaskRequestRecord,
	err error) {

	// Get tasks collection
	tasksCollection := mongoDB.GetCollection(types.TaskCollection)

	// Set filtering for this task
	task_request_filter := bson.D{{Key: "_id", Value: taskID}}
	opts := options.FindOne()
	// Set mongo context
	ctxM, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Now retrieve all node task requests entries
	cursor := tasksCollection.FindOne(ctxM, task_request_filter, opts)
	var taskReq types.TaskRequestRecord
	if err := cursor.Decode(&taskReq); err != nil {
		l.Error().Msg("Could not decode task request data from MongoDB.")
		return taskReq, err
	}

	return taskReq, nil

}

// Given a TaskID from MongoDB, deletes all associated entries from the "tasks", "instances", "prompts", "responses" and "results" collections.
func RemoveTaskID(taskID primitive.ObjectID, mongoDB mongodb.MongoDb, l *zerolog.Logger) (err error) {

	//--------------------------------------------------------------------------
	//-------------------------- Instances -------------------------------------
	//--------------------------------------------------------------------------
	instancesCollection := mongoDB.GetCollection(types.InstanceCollection)
	// Set filtering for this node-service pair data
	task_request_filter := bson.D{{Key: "task_id", Value: taskID}}
	// Set mongo context
	ctxM, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Now retrieve all node task requests entries
	response, err := instancesCollection.DeleteMany(ctxM, task_request_filter)
	if err != nil {
		l.Warn().Msg("Could not delete instances data from MongoDB.")
		return err
	}

	l.Debug().Int("deleted_count", int(response.DeletedCount)).Str("TaskID", taskID.String()).Msg("deleted instances data from MongoDB")

	//--------------------------------------------------------------------------
	//-------------------------- Prompts ---------------------------------------
	//--------------------------------------------------------------------------
	promptsCollection := mongoDB.GetCollection(types.PromptsCollection)
	// Set filtering for this node-service pair data
	task_request_filter = bson.D{{Key: "task_id", Value: taskID}}
	// Set mongo context
	ctxM, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Now retrieve all node task requests entries
	response, err = promptsCollection.DeleteMany(ctxM, task_request_filter)
	if err != nil {
		l.Warn().Msg("Could not delete prompts data from MongoDB.")
		return err
	}

	l.Debug().Int("deleted", int(response.DeletedCount)).Str("TaskID", taskID.String()).Msg("deleted prompts data from MongoDB")

	//--------------------------------------------------------------------------
	//-------------------------- Responses -------------------------------------
	//--------------------------------------------------------------------------
	responsesCollection := mongoDB.GetCollection(types.ResponsesCollection)
	// Set filtering for this node-service pair data
	task_request_filter = bson.D{{Key: "task_id", Value: taskID}}
	// Set mongo context
	ctxM, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Now retrieve all node task requests entries
	response, err = responsesCollection.DeleteMany(ctxM, task_request_filter)
	if err != nil {
		l.Warn().Msg("Could not delete responses data from MongoDB.")
		return err
	}

	l.Debug().Int("deleted_count", int(response.DeletedCount)).Str("TaskID", taskID.String()).Msg("deleted responses data from MongoDB")

	//--------------------------------------------------------------------------
	//-------------------------- Results ---------------------------------------
	//--------------------------------------------------------------------------
	resultsCollection := mongoDB.GetCollection(types.ResultsCollection)
	// Set filtering for this node-service pair data
	task_request_filter = bson.D{{Key: "result_data.task_id", Value: taskID}}
	// Set mongo context
	ctxM, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Now retrieve all node task requests entries
	response, err = resultsCollection.DeleteMany(ctxM, task_request_filter)
	if err != nil {
		l.Warn().Msg("Could not delete results data from MongoDB.")
		return err
	}

	l.Debug().Int("deleted_count", int(response.DeletedCount)).Str("TaskID", taskID.String()).Msg("deleted results data from MongoDB")

	//--------------------------------------------------------------------------
	//-------------------------- Task ------------------------------------------
	//--------------------------------------------------------------------------
	tasksCollection := mongoDB.GetCollection(types.TaskCollection)
	// Set filtering for this node-service pair data
	task_request_filter = bson.D{{Key: "_id", Value: taskID}}
	// Set mongo context
	ctxM, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Now retrieve all node task requests entries
	response, err = tasksCollection.DeleteMany(ctxM, task_request_filter)
	if err != nil {
		l.Warn().Msg("Could not delete task data from MongoDB.")
		return err
	}

	l.Debug().Int("deleted_count", int(response.DeletedCount)).Str("TaskID", taskID.String()).Msg("deleted task data from MongoDB")

	return nil

}
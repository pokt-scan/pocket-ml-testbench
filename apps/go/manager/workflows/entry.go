package workflows

import (
	"manager/types"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Ctx represents a context struct that holds an instance of `app.App`
// This is created because sharing dependencies by context is not recommended
type Ctx struct {
	App *types.App
}

// Workflows represent the context for executing workflow logic.
var Workflows *Ctx

// SetAppConfig sets the provided app config to the Workflows global variable in the Ctx struct.
func SetAppConfig(ac *types.App) {
	Workflows = &Ctx{
		App: ac,
	}
}

// Register registers the workflows with the provided worker.
func (wCtx *Ctx) Register(w worker.Worker) {

	// Main workflow, containing the logic of:
	// - Staked nodes retrieval
	// - Analyze nodes data
	// - Triggering new evaluation tasks
	w.RegisterWorkflowWithOptions(wCtx.NodeManager, workflow.RegisterOptions{
		Name: NodeManagerName,
	})
}

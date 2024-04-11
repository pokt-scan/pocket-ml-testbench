package activities

import (
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"requester/types"
)

// Ctx represents a context struct that holds an instance of `app.App`
// This is created because sharing dependencies by context is not recommended
type Ctx struct {
	App *types.App
}

// Activities represent the context for executing activities logic.
var Activities *Ctx

// SetAppConfig sets the provided app configuration to the global Activities variable in the Ctx struct.
func SetAppConfig(ac *types.App) {
	Activities = &Ctx{
		App: ac,
	}
}

// Register registers a worker activity with the provided activity function in the Ctx struct.
func (aCtx *Ctx) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(aCtx.GetApp, activity.RegisterOptions{
		Name: GetAppName,
	})

	w.RegisterActivityWithOptions(aCtx.GetHeight, activity.RegisterOptions{
		Name: GetHeightName,
	})

	w.RegisterActivityWithOptions(aCtx.GetBlock, activity.RegisterOptions{
		Name: GetBlockName,
	})

	w.RegisterActivityWithOptions(aCtx.GetSession, activity.RegisterOptions{
		Name: GetSessionName,
	})

	w.RegisterActivityWithOptions(aCtx.LookupTaskRequest, activity.RegisterOptions{
		Name: LookupTaskRequestName,
	})

	w.RegisterActivityWithOptions(aCtx.Relayer, activity.RegisterOptions{
		Name: RelayerName,
	})
}

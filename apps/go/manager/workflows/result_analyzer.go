package workflows

import (
	"time"

	"manager/activities"
	"manager/types"

	"go.temporal.io/sdk/workflow"
)

var ResultAnalyzerName = "Manager-ResultAnalyzer"

// ResultAnalyzer - Is a method that processes the results generated by the evaluator.
func (wCtx *Ctx) ResultAnalyzer(ctx workflow.Context, params types.ResultAnalyzerParams) (*types.ResultAnalyzerResults, error) {

	l := wCtx.App.Logger
	l.Debug().Msg("Starting Result Analyzer Workflow.")

	// Create result
	result := types.ResultAnalyzerResults{Success: false}

	// -------------------------------------------------------------------------
	// -------------------- Analyze Result -------------------------------------
	// -------------------------------------------------------------------------
	// Set timeout to get staked nodes activity
	ctxTimeout := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 5,
		StartToCloseTimeout:    time.Minute * 5,
	})
	// Set activity input
	getAnalyzeResultInput := types.AnalyzeResultParams{
		TaskID: params.TaskID,
	}
	// Results will be kept logged by temporal
	var resultAnalysisData types.AnalyzeResultResults
	// Execute activity
	err := workflow.ExecuteActivity(ctxTimeout, activities.AnalyzeResultName, getAnalyzeResultInput).Get(ctx, &resultAnalysisData)
	if err != nil {
		return &result, err
	}

	return &result, nil
}
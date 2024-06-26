package tests

import (
	"errors"
	poktGoSdk "github.com/pokt-foundation/pocket-go/provider"
	"github.com/stretchr/testify/mock"
	"packages/pocket_rpc/samples"
	"reflect"
	"requester/activities"
)

// define a test suite struct
type GetAppUnitTestSuite struct {
	BaseSuite
}

func (s *GetAppUnitTestSuite) Test_GetApp_Activity() {
	app := samples.GetAppMock(s.app.Logger)
	getAppParams := activities.GetAppParams{
		Address: "f3abbe313689a603a1a6d6a43330d0440a552288",
		Service: "0001",
	}

	s.GetPocketRpcMock().
		On("GetApp", getAppParams.Address).
		Return(app, nil).
		Times(1)

	// Run the Activity in the test environment
	future, err := s.activityEnv.ExecuteActivity(activities.Activities.GetApp, getAppParams)
	// Check there was no error on the call to execute the Activity
	s.NoError(err)
	// rpc must be called once
	s.GetPocketRpcMock().AssertExpectations(s.T())
	// Check that there was no error returned from the Activity
	result := poktGoSdk.App{}
	s.NoError(future.Get(&result))
	// Check for the expected return value.
	s.True(reflect.DeepEqual(&result, app))
}

func (s *GetAppUnitTestSuite) Test_GetApp_Error_Activity() {
	getAppParams := activities.GetAppParams{
		Address: "f3abbe313689a603a1a6d6a43330d0440a552288",
		Service: "0001",
	}

	s.GetPocketRpcMock().
		On("GetApp", mock.Anything).
		Return(nil, errors.New("not found")).
		Times(1)

	// Run the Activity in the test environment
	_, err := s.activityEnv.ExecuteActivity(activities.Activities.GetApp, getAppParams)
	// Check there was no error on the call to execute the Activity
	s.Error(err)
	s.GetPocketRpcMock().AssertExpectations(s.T())
}

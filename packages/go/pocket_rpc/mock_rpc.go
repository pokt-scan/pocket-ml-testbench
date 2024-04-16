package pocket_rpc

import (
	"errors"
	poktGoSdk "github.com/pokt-foundation/pocket-go/provider"
	"github.com/stretchr/testify/mock"
)

const UnexpectedResponseType = "one of the response arguments does not match the expected type"

type MockRpc struct {
	mock.Mock
}

func NewMockRpc() *MockRpc {
	mockRpc := MockRpc{}
	return &mockRpc
}

func (rpc *MockRpc) GetClientPool() *ClientPool {
	return nil
}

func (rpc *MockRpc) SetClientPool(_ *ClientPool) {
}

func (rpc *MockRpc) GetApp(address string) (*poktGoSdk.App, error) {
	args := rpc.Called(address)
	firstResponseArg := args.Get(0)
	var response *poktGoSdk.App
	if firstResponseArg != nil {
		if v, ok := firstResponseArg.(*poktGoSdk.App); !ok {
			return nil, errors.New(UnexpectedResponseType)
		} else {
			response = v
		}
	}

	return response, args.Error(1)
}

func (rpc *MockRpc) GetNodes(service string) (nodes []*poktGoSdk.Node, e error) {
	args := rpc.Called(service)
	firstResponseArg := args.Get(0)
	var response []*poktGoSdk.Node
	if firstResponseArg != nil {
		if v, ok := firstResponseArg.([]*poktGoSdk.Node); !ok {
			return nil, errors.New(UnexpectedResponseType)
		} else {
			response = v
		}
	}

	return response, args.Error(1)
}

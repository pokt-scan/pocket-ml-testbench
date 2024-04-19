package activities

import (
	"context"
)

type GetStakedParams struct {
	Service string
}

type NodeData struct {
	Address string
	Service string
}

type GetStakedResults struct {
	Nodes []NodeData
}

var GetStakedName = "get_staked"

func (aCtx *Ctx) GetStaked(ctx context.Context, params GetStakedParams) (*GetStakedResults, error) {

	l := aCtx.App.Logger
	l.Debug().Msg("Collecting staked nodes from network.")

	result := GetStakedResults{}

	// Get all nodes in given chain
	l.Debug().Str("service", params.Service).Msg("Querying service...")
	nodes, err := aCtx.App.PocketRpc.GetNodes(params.Service)
	if err != nil {
		l.Error().Str("service", params.Service).Msg("Could not retrieve staked nodes.")
		return nil, err
	}
	if len(nodes) == 0 {
		l.Warn().Str("service", params.Service).Msg("No nodes found staked.")
	}
	for _, node := range nodes {
		if !node.Jailed {
			this_node := NodeData{
				Address: node.Address,
				Service: params.Service,
			}
			result.Nodes = append(result.Nodes, this_node)
		}
	}

	if len(result.Nodes) == 0 {
		l.Warn().Msg("No nodes were found on any of the given services")
	} else {
		l.Info().Int("nodes_staked", len(result.Nodes)).Msg("Successfully pulled staked node-services.")
	}

	// // cheap mock
	// for i := 0; i < 5; i++ {
	// 	thisNode := NodeData{Address: fmt.Sprint(i), Service: fmt.Sprint(i * 10)}
	// 	result.Nodes = append(result.Nodes, thisNode)
	// }

	return &result, nil
}
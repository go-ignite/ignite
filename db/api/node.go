package api

import "github.com/go-ignite/ignite/db"

func (api *API) GetAllNodes() ([]*db.Node, error) {
	var nodes []*db.Node
	return nodes, api.Find(&nodes)
}

func (api *API) UpsertNode(node *db.Node) (int64, error) {
	if node.Id == 0 {
		return api.Insert(node)
	}
	return api.ID(node.Id).Cols("name", "comment").Update(node)
}

func (api *API) GetNodeByID(id int64) (*db.Node, error) {
	node := &db.Node{}
	_, err := api.ID(id).Get(node)
	return node, err

}

func (api *API) DeleteNode(id int64) (int64, error) {
	return api.ID(id).Delete(new(db.Node))
}

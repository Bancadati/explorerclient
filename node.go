package explorerclient

import (
	"net/http"
	"net/url"
)

// Node represents a node
type Node struct {
	ID                int64          `json:"id"`
	Hostname          string         `json:"hostname"`
	NodeID            string         `json:"node_id"`
	NodeIDV1          string         `json:"node_id_v1"`
	FarmID            int64          `json:"farm_id"`
	OsVersion         string         `json:"os_version"`
	Created           Date           `json:"created"`
	Updated           Date           `json:"updated"`
	Uptime            int64          `json:"uptime"`
	Address           string         `json:"address"`
	Location          Location       `json:"location"`
	TotalResources    ResourceAmount `json:"total_resources"`
	UsedResources     ResourceAmount `json:"used_resources"`
	ReservedResources ResourceAmount `json:"reserved_resources"`
	Workloads         WorkloadAmount `json:"workloads"`
	Proofs            []Proof        `json:"proofs"`
	Ifaces            []Iface        `json:"ifaces"`
	PublicConfig      *PublicIface   `json:"public_config"`
	FreeToUse         bool           `json:"free_to_use"`
	Approved          bool           `json:"approved"`
	PublicKeyHex      string         `json:"public_key_hex"`
	WgPorts           []int64        `json:"wg_ports"`
}

// ListNodes lists all nodes from on the explorer
func (cl *Client) ListNodes(f *NodeFilter, page *Pager) ([]Node, error) {
	nodes := []Node{}
	query := url.Values{}
	f.Apply(query)
	page.apply(query)
	_, err := cl.get(cl.geturl("nodes"), query, &nodes, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetNode returns the details of a single node
func (cl *Client) GetNode(id string) (*Node, error) {
	node := &Node{}
	query := url.Values{}
	_, err := cl.get(cl.geturl("nodes", id), query, node, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return node, nil
}

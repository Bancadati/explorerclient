package explorerclient

import (
	"net/http"
	"net/url"
	"strconv"
)

// Farm represents a farm
type Farm struct {
	ID              int64               `json:"id"`
	ThreebotID      int64               `json:"threebot_id"`
	IyoOrganization string              `json:"iyo_organization"`
	Name            string              `json:"name"`
	WalletAddresses []WalletAddress     `json:"wallet_addresses"`
	Location        Location            `json:"location"`
	Email           string              `json:"email"`
	ResourcePrices  []NodeResourcePrice ` json:"resource_prices"`
	PrefixZero      IPRange             `json:"prefix_zero"`
}

// ListFarms lists all farms on the explorer
func (cl *Client) ListFarms(page *Pager) ([]Farm, error) {
	farms := []Farm{}
	query := url.Values{}
	page.apply(query)
	_, err := cl.get(cl.geturl("farms"), query, &farms, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return farms, nil
}

// GetFarm returns the details of a single farm
func (cl *Client) GetFarm(id int64) (*Farm, error) {
	farm := &Farm{}
	_, err := cl.get(cl.geturl("farms", strconv.FormatInt(id, 10)), nil, farm, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return farm, nil
}

// Package explorerclient provides a simple HTTP client for the Threefold explorer API
// This package simplifies and reduces the dependencies from the Threefold implementation (https://github.com/threefoldtech/tfexplorer/tree/master/client)
// This also means not all functionality is used and only the API endpoints needed for Bancadati are implemented.
package explorerclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	// DefaultBaseURL represents the default API URL for the Threefold Explorer
	DefaultBaseURL = "https://explorer.grid.tf/explorer/"
)

var (
	successCodes = []int{
		http.StatusOK,
		http.StatusCreated,
	}
)

// NewClient returns a new TF explorer client
func NewClient(baseURL string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("%s is an invalid base URL: %w", baseURL, err)
	}

	cl := &Client{
		baseURL: u,
		http:    &http.Client{},
	}

	return cl, nil
}

// Client represents a TF explorer client
type Client struct {
	http    *http.Client
	baseURL *url.URL
}

func (c *Client) geturl(p ...string) string {
	b := *c.baseURL
	b.Path = path.Join(b.Path, filepath.Join(p...))

	return b.String()
}

func (c *Client) get(u string, query url.Values, output interface{}, expect ...int) (*http.Response, error) {
	if len(query) > 0 {
		u = fmt.Sprintf("%s?%s", u, query.Encode())
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	response, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return response, c.process(response, output, expect...)
}

func (c *Client) process(response *http.Response, output interface{}, expect ...int) error {
	defer response.Body.Close()

	if len(expect) == 0 {
		expect = successCodes
	}

	in := func(i int, l []int) bool {
		for _, x := range l {
			if x == i {
				return true
			}
		}
		return false
	}

	dec := json.NewDecoder(response.Body)
	if !in(response.StatusCode, expect) {
		var output struct {
			E string `json:"error"`
		}

		if err := dec.Decode(&output); err != nil {
			return errors.Wrapf(HTTPError{
				err:  err,
				resp: response,
			}, "failed to load error while processing invalid return code of: %s", response.Status)
		}

		return HTTPError{
			err:  fmt.Errorf(output.E),
			resp: response,
		}
	}

	if output == nil {
		//discard output
		ioutil.ReadAll(response.Body)
		return nil
	}

	if err := dec.Decode(output); err != nil {
		return HTTPError{
			err:  errors.Wrap(err, "failed to load output"),
			resp: response,
		}
	}

	return nil
}

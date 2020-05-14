# Explorer client

This repo contains a simple and limited implementation of a ThreeFold explorer client.  
[https://explorer.grid.tf](https://explorer.grid.tf)


This client only implements endpoints needed for the purposes of Bancadati and is therefor very limited functionality.

## Usage

This is an example that lists all the nodes of a certain farm in pages of 10 nodes.

```go
package main

import (
	"fmt"
	"log"

	"github.com/bancadati/explorerclient"
)

func main() {
	cl, err := explorerclient.NewClient("https://explorer.grid.tf/explorer")
	if err != nil {
		log.Fatal(err)
	}

	f := &explorerclient.NodeFilter{}
	f.WithFarm(173636)

	pageSize := 10
	page := 1

	for {
		nodePager := explorerclient.NewPager(page, pageSize)
		nodes, err := cl.ListNodes(f, nodePager)
		if err != nil {
			log.Fatal(err)
		}
		for _, node := range nodes {
			fmt.Println(node.NodeID, node.Location)
		}

		if len(nodes) == pageSize {
			page++
			continue
		}

		break
	}
}

```
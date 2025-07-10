package snow

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

func GetID() (int64, error) {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// Generate a snowflake ID.
	id := node.Generate()
	return id.Int64(), nil

}

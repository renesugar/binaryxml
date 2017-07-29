package binaryxml

import (
	"encoding/xml"
)

type xmlTraversalNode struct {
	XMLName xml.Name
	Content []byte             `xml:",innerxml"`
	Nodes   []xmlTraversalNode `xml:",any"`
}

func walk(nodes []xmlTraversalNode, f func(xmlTraversalNode) bool) {
	for _, node := range nodes {
		if f(node) {
			walk(node.Nodes, f)
		}
	}
}

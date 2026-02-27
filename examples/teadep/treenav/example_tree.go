package main

import (
	"github.com/mikeschinkel/go-tealeaves/teadep"
)

// ExampleTree returns a sample dependency tree for demonstration.
// This represents a fictional Go project with several modules.
func ExampleTree() *teadep.Tree {
	return exampleTree
}

// Module kind constants
const (
	kindLib  = 1
	kindExe  = 2
	kindTest = 3
)

// exampleNode embeds BaseNode and adds a Kind field for the example
type exampleNode struct {
	*teadep.BaseNode
	Kind    int
	KindSet bool
}

func newNode(name string, kind int, deps ...teadep.Node) *exampleNode {
	return &exampleNode{
		BaseNode: teadep.NewBaseNode(name, &teadep.BaseNodeArgs{
			Dependencies: deps,
		}),
		Kind:    kind,
		KindSet: kind != 0,
	}
}

var exampleTree *teadep.Tree

func init() {
	// Leaf dependencies (no deps of their own)
	dtModule := newNode("github.com/mikeschinkel/go-dt", kindLib)
	logModule := newNode("github.com/mikeschinkel/go-logutil", kindLib)

	// Mid-level dependencies
	cliutilModule := newNode("github.com/mikeschinkel/go-cliutil", kindLib,
		dtModule,
		logModule,
	)
	cfgstoreModule := newNode("github.com/mikeschinkel/go-cfgstore", kindLib,
		cliutilModule,
		dtModule,
	)

	// Core library
	coreModule := newNode("github.com/example/myapp/core", kindLib,
		cfgstoreModule,
		cliutilModule,
		dtModule,
	)

	// CLI binary
	cliModule := newNode("github.com/example/myapp/cmd/myapp", kindExe,
		coreModule,
		cliutilModule,
		cfgstoreModule,
		dtModule,
	)

	// Test module
	testModule := newNode("github.com/example/myapp/test", kindTest,
		coreModule,
	)

	// Root: the repository containing all modules
	repoNode := teadep.NewBaseNode("git@github.com:example/myapp.git", &teadep.BaseNodeArgs{
		Dependencies: []teadep.Node{
			cliModule,
			coreModule,
			testModule,
		},
	})

	exampleTree = teadep.NewTree(repoNode)
}

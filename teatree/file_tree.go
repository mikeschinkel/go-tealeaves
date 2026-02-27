package teatree

import (
	"path/filepath"
	"sort"

	"github.com/mikeschinkel/go-dt"
)

type BuildFileTreeArgs struct {
	RootPath dt.PathSegment
}

// BuildTree creates a hierarchical tree from flat file list
// Returns top-level nodes (files and folders at root level)
func BuildFileTree(files []*File, args BuildFileTreeArgs) (nodes []*FileNode) {
	name := ","
	if args.RootPath != "" {
		name = string(args.RootPath)
	}
	root := NewNode(".", name, File{
		Path: dt.RelFilepath(name),
	})

	// Sort files by path first - this allows efficient tree building
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	// Build tree structure using path-based node map for O(1) lookups
	nodeMap := make(map[string]*FileNode, len(files)-1)
	nodeMap["."] = root

	for _, file := range files {
		realPath := file.Path
		if realPath == "" {
			continue
		}

		// Compute synthetic path for tree building (prefix + real path)
		var treePath dt.RelFilepath
		if args.RootPath != "" {
			treePath = dt.RelFilepathJoin(args.RootPath, realPath)
		} else {
			treePath = realPath
		}

		segments := treePath.Split("/")
		currentPath := ""

		// Create folder nodes for each segment (except last, which is the file)
		for i := 0; i < len(segments)-1; i++ {
			segment := segments[i]
			switch {
			case currentPath == "":
				currentPath = string(segment)
			default:
				currentPath = currentPath + "/" + string(segment)
			}

			// Check if this folder node already exists
			_, exists := nodeMap[currentPath]
			if !exists {
				// Create new folder node (id=fullPath, name=basename)
				folderNode := NewNode(
					currentPath,                // id
					filepath.Base(currentPath), // name (basename for display)
					File{
						Path: dt.RelFilepath(currentPath),
					},
				)

				// Set text to tree path for this folder (used for tree structure)
				folderNode.SetText(currentPath)

				// Find parent node
				parentPath := filepath.Dir(currentPath)
				if parentPath == "" {
					parentPath = "."
				}
				parentNode := nodeMap[parentPath]

				// Add to parent
				parentNode.AddChild(folderNode)

				// Add to map
				nodeMap[currentPath] = folderNode
			}
		}

		// Add file node (id=treePath, name=basename)
		fileNode := NewNode(
			string(treePath),        // id (synthetic path for tree structure)
			string(realPath.Base()), // name (basename for display)
			*file,                   // data (File.Path remains REAL)
		)

		// Set text to real path (used for file I/O)
		fileNode.SetText(string(realPath))

		// Collapse all folders (first level should be visible but collapsed)
		fileNode.expanded = false

		// Find parent node for this file (using tree path)
		parentPath := treePath.Dir()
		if parentPath == "" {
			parentPath = "."
		}
		parentNode := nodeMap[string(parentPath)]

		// Add file to parent folder
		parentNode.AddChild(fileNode)
	}

	// Always return children of temporary root
	// When RootPath is set, the prefixed paths naturally create a folder node (e.g., "go-dt")
	// which becomes the actual root in the returned tree
	return root.Children()
}

package teatree

import (
	"os"
	"time"

	"github.com/mikeschinkel/go-dt"
)

// FileMeta contains cached meta about a file for display in directory tables.
type FileMeta struct {
	Size        int64          // File size in bytes
	ModTime     time.Time      // Modification time
	Permissions os.FileMode    // Full permissions
	EntryStatus dt.EntryStatus // File, Dir, Symlink, etc.
	Data        any
}

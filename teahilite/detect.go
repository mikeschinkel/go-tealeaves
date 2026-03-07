package teahilite

import (
	"path/filepath"

	"github.com/alecthomas/chroma/v2/lexers"
)

// DetectLanguage returns the Chroma lexer alias for the given file path
// (e.g. "go", "python", "javascript") based on its extension. Returns "text"
// if no lexer matches. The generic constraint accepts dt.Filepath,
// dt.RelFilepath, dt.Filename, or plain string.
func DetectLanguage[S ~string](path S) (name string) {
	name = "text"

	ext := filepath.Ext(string(path))
	if ext == "" {
		goto end
	}

	{
		lexer := lexers.Match(string(path))
		if lexer == nil {
			goto end
		}

		config := lexer.Config()
		if config == nil {
			goto end
		}
		if len(config.Aliases) == 0 {
			goto end
		}

		name = config.Aliases[0]
	}

end:
	return name
}

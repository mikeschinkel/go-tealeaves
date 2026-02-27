package teatree

import (
	"os"

	"github.com/mikeschinkel/go-dt"
)

type File struct {
	Path    dt.RelFilepath
	content *string
	data    any
	YOffset int // Viewport scroll position

	// Cached file metadata (for directory table display)
	meta *FileMeta // nil if not yet loaded
}

func NewFile(path dt.RelFilepath, content *string) *File {
	return &File{
		Path:    path,
		content: content,
	}
}

func (f *File) IsEmpty() bool {
	return f.Path == ""
}

func (f *File) Data() any {
	if f.data == nil {
		panic("File.Data() called before ....???")
	}
	return f.data
}

func (f *File) HasData() bool {
	return f.data != nil
}

func (f *File) Content() string {
	if f.content == nil {
		panic("File.Content() called before ....???")
	}
	return *f.content
}

func (f *File) SetContent(content string) {
	f.content = &content
}

func (f *File) HasContent() bool {
	return f.content != nil
}

func (f *File) WithData(data any) *File {
	f.data = data
	return f
}

func (f *File) WithYOffset(yOfs int) *File {
	f.YOffset = yOfs
	return f
}

func (f *File) Meta() *FileMeta {
	if f.meta == nil {
		panic("File.Meta() called before File.LoadMeta()")
	}
	return f.meta
}

func (f *File) HasMeta() bool {
	return f.meta != nil
}

func (f *File) WithMeta(meta *FileMeta) *File {
	f.meta = meta
	return f
}

func (f *File) LoadMeta(root dt.DirPath) (err error) {
	f.meta, err = LoadFileMeta(root, f.Path)
	return err
}

func LoadFileMeta(root dt.DirPath, path dt.RelFilepath) (fm *FileMeta, err error) {
	var fp dt.Filepath
	var info os.FileInfo

	// Construct full path
	fp = dt.FilepathJoin(root, path)

	// Get f info using dt.Filepath.Stat() method
	info, err = fp.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			// File might be deleted after git status - return nil (not an error)
			fm = &FileMeta{}
			err = nil
			goto end
		}
		err = NewErr(dt.ErrFileSystem, dt.ErrFileStat, fp.ErrKV(), err)
		goto end
	}

	// Initialize metadata if not already present
	fm = &FileMeta{
		Size:        info.Size(),
		ModTime:     info.ModTime(),
		Permissions: info.Mode(),
		EntryStatus: dt.GetEntryStatus(info),
	}
end:
	return fm, err
}

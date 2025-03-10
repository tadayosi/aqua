package unarchive

import (
	"fmt"
	"io"

	"github.com/mholt/archiver/v3"
	"github.com/spf13/afero"
)

type unarchiverWithUnarchiver struct {
	unarchiver archiver.Unarchiver
	dest       string
}

func (unarchiver *unarchiverWithUnarchiver) Unarchive(fs afero.Fs, body io.Reader) error {
	dest := unarchiver.dest
	f, err := afero.TempFile(fs, "", "")
	if err != nil {
		return fmt.Errorf("create a temporal file: %w", err)
	}
	defer func() {
		f.Close()
		fs.Remove(f.Name()) //nolint:errcheck
	}()
	if _, err := io.Copy(f, body); err != nil {
		return fmt.Errorf("copy the file to the temporal file: %w", err)
	}
	return unarchiver.unarchiver.Unarchive(f.Name(), dest) //nolint:wrapcheck
}

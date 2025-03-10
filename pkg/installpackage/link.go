package installpackage

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func (inst *installer) createLink(linkPath, linkDest string, logE *logrus.Entry) error {
	if fileInfo, err := inst.linker.Lstat(linkPath); err == nil {
		switch mode := fileInfo.Mode(); {
		case mode.IsDir():
			// if file is a directory, raise error
			return fmt.Errorf("%s has already existed and is a directory", linkPath)
		case mode&os.ModeNamedPipe != 0:
			// if file is a pipe, raise error
			return fmt.Errorf("%s has already existed and is a named pipe", linkPath)
		case mode.IsRegular():
			// TODO if file is a regular file, remove it and create a symlink.
			return fmt.Errorf("%s has already existed and is a regular file", linkPath)
		case mode&os.ModeSymlink != 0:
			return inst.recreateLink(linkPath, linkDest, logE)
		default:
			return fmt.Errorf("unexpected file mode %s: %s", linkPath, mode.String())
		}
	}
	logE.WithFields(logrus.Fields{
		"link_file": linkPath,
		"new":       linkDest,
	}).Info("create a symbolic link")
	if err := inst.linker.Symlink(linkDest, linkPath); err != nil {
		return fmt.Errorf("create a symbolic link: %w", err)
	}
	return nil
}

func (inst *installer) recreateLink(linkPath, linkDest string, logE *logrus.Entry) error {
	lnDest, err := inst.linker.Readlink(linkPath)
	if err != nil {
		return fmt.Errorf("read a symbolic link (%s): %w", linkPath, err)
	}
	if linkDest == lnDest {
		return nil
	}
	// recreate link
	logE.WithFields(logrus.Fields{
		// TODO add version
		"link_file": linkPath,
		"old":       lnDest,
		"new":       linkDest,
	}).Debug("recreate a symbolic link")
	if err := inst.fs.Remove(linkPath); err != nil {
		return fmt.Errorf("remove a symbolic link (%s): %w", linkPath, err)
	}
	if err := inst.linker.Symlink(linkDest, linkPath); err != nil {
		return fmt.Errorf("create a symbolic link: %w", err)
	}
	return nil
}

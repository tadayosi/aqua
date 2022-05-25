package installpackage

import (
	"context"

	"github.com/aquaproj/aqua/pkg/checksum"
	"github.com/aquaproj/aqua/pkg/config"
	"github.com/aquaproj/aqua/pkg/download"
	"github.com/aquaproj/aqua/pkg/link"
	"github.com/aquaproj/aqua/pkg/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Installer interface {
	InstallPackage(ctx context.Context, pkgInfo *config.PackageInfo, pkg *config.Package, isTest, checksumEnabled bool, logE *logrus.Entry) error
	InstallPackages(ctx context.Context, cfg *config.Config, registries map[string]*config.RegistryContent, binDir string, onlyLink, isTest bool, logE *logrus.Entry) error
	InstallProxy(ctx context.Context, logE *logrus.Entry) error
	ReadChecksumFile(fs afero.Fs, p string) error
	UpdateChecksumFile(fs afero.Fs, p string) error
}

func New(param *config.Param, downloader download.PackageDownloader, rt *runtime.Runtime, fs afero.Fs, linker link.Linker) Installer {
	return &installer{
		rootDir:           param.RootDir,
		maxParallelism:    param.MaxParallelism,
		packageDownloader: downloader,
		runtime:           rt,
		fs:                fs,
		linker:            linker,
		checksums:         checksum.New(),
	}
}

// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package controller

import (
	"context"
	"github.com/aquaproj/aqua/pkg/config"
	"github.com/aquaproj/aqua/pkg/config-finder"
	"github.com/aquaproj/aqua/pkg/config-reader"
	"github.com/aquaproj/aqua/pkg/controller/initcmd"
	"github.com/aquaproj/aqua/pkg/controller/list"
	"github.com/aquaproj/aqua/pkg/download"
	"github.com/aquaproj/aqua/pkg/github"
	"github.com/aquaproj/aqua/pkg/install-registry"
	"github.com/aquaproj/aqua/pkg/installpackage"
	"github.com/aquaproj/aqua/pkg/log"
)

// Injectors from wire.go:

func NewController(ctx context.Context, aquaVersion string, param *config.Param) (*Controller, error) {
	rootDir := config.NewRootDir()
	configFinder := finder.NewConfigFinder()
	fileReader := reader.NewFileReader()
	configReader := reader.New(fileReader)
	logger := log.NewLogger(aquaVersion)
	repositoryService := github.NewGitHub(ctx)
	packageDownloader := download.NewPackageDownloader(repositoryService, logger)
	installer := installpackage.New(aquaVersion, logger, packageDownloader)
	registryDownloader := download.NewRegistryDownloader(repositoryService, logger)
	registryInstaller := registry.New(rootDir, logger, registryDownloader)
	controller, err := New(ctx, rootDir, configFinder, configReader, logger, installer, repositoryService, registryInstaller, param)
	if err != nil {
		return nil, err
	}
	return controller, nil
}

func InitializeListCommandController(ctx context.Context, aquaVersion string, param *config.Param) *list.Controller {
	configFinder := finder.NewConfigFinder()
	fileReader := reader.NewFileReader()
	configReader := reader.New(fileReader)
	rootDir := config.NewRootDir()
	logger := log.NewLogger(aquaVersion)
	repositoryService := github.NewGitHub(ctx)
	registryDownloader := download.NewRegistryDownloader(repositoryService, logger)
	installer := registry.New(rootDir, logger, registryDownloader)
	controller := list.NewController(configFinder, configReader, installer)
	return controller
}

func InitializeInitCommandController(ctx context.Context, aquaVersion string, param *config.Param) *initcmd.Controller {
	logger := log.NewLogger(aquaVersion)
	repositoryService := github.NewGitHub(ctx)
	controller := initcmd.New(logger, repositoryService)
	return controller
}

func InitializeGenerateCommandController(ctx context.Context, aquaVersion string, param *config.Param) *GenerateController {
	configFinder := finder.NewConfigFinder()
	fileReader := reader.NewFileReader()
	configReader := reader.New(fileReader)
	logger := log.NewLogger(aquaVersion)
	rootDir := config.NewRootDir()
	repositoryService := github.NewGitHub(ctx)
	registryDownloader := download.NewRegistryDownloader(repositoryService, logger)
	installer := registry.New(rootDir, logger, registryDownloader)
	generateController := NewGenerateController(configFinder, configReader, logger, installer, repositoryService)
	return generateController
}

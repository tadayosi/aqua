package generate

import (
	"fmt"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

type FuzzyFinder interface {
	Find(pkgs []*FindingPackage) ([]int, error)
}

type fuzzyFinder struct{}

func NewFuzzyFinder() FuzzyFinder {
	return &fuzzyFinder{}
}

func NewMockFuzzyFinder(idxs []int, err error) FuzzyFinder {
	return &mockFuzzyFinder{
		idxs: idxs,
		err:  err,
	}
}

type mockFuzzyFinder struct {
	idxs []int
	err  error
}

func (finder *mockFuzzyFinder) Find(pkgs []*FindingPackage) ([]int, error) {
	return finder.idxs, finder.err
}

func (finder *fuzzyFinder) Find(pkgs []*FindingPackage) ([]int, error) {
	return fuzzyfinder.FindMulti(pkgs, func(i int) string { //nolint:wrapcheck
		return find(pkgs[i])
	},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i < 0 {
				return "No package matches"
			}
			return getPreview(pkgs[i], i, w)
		}))
}

func find(pkg *FindingPackage) string {
	files := pkg.PackageInfo.GetFiles()
	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if file.Name == "" {
			continue
		}
		fileNames = append(fileNames, file.Name)
	}
	fileNamesStr := strings.Join(fileNames, ", ")
	aliases := make([]string, 0, len(pkg.PackageInfo.Aliases))
	for _, alias := range pkg.PackageInfo.Aliases {
		if alias.Name == "" {
			continue
		}
		aliases = append(aliases, alias.Name)
	}
	pkgName := pkg.PackageInfo.GetName()
	item := pkgName
	if len(aliases) != 0 {
		item += " (" + strings.Join(aliases, ", ") + ")"
	}
	item += " (" + pkg.RegistryName + ")"
	if strings.HasSuffix(pkgName, "/"+fileNamesStr) || pkgName == fileNamesStr {
		return item
	}
	return item + " (" + fileNamesStr + ")"
}

func getPreview(pkg *FindingPackage, i, w int) string {
	if i < 0 {
		return ""
	}
	return fmt.Sprintf("%s\n\n%s\n%s",
		pkg.PackageInfo.GetName(),
		pkg.PackageInfo.GetLink(),
		formatDescription(pkg.PackageInfo.GetDescription(), w/2-8)) //nolint:gomnd
}

func formatDescription(desc string, w int) string {
	descRune := []rune(desc)
	lenDescRune := len(descRune)
	lineWidth := w - len([]rune("\n"))
	numOfLines := (lenDescRune / lineWidth) + 1
	descArr := make([]string, numOfLines)
	for i := 0; i < numOfLines; i++ {
		start := i * lineWidth
		end := start + lineWidth
		if i == numOfLines-1 {
			end = lenDescRune
		}
		descArr[i] = string(descRune[start:end])
	}
	return strings.Join(descArr, "\n")
}

package exportproject

import (
	"io"
	"os"
	"path/filepath"

	"github.com/golangaccount/cmd.go.internal/load"
	gos "github.com/golangaccount/go-libs/os"
)

//ExportPackage 导出package code
//pkgs:需要进行导出的pkg
//dest:导出的文件夹地址
//filter:导出过滤
func ExportPackage(pkgs []string, dest string, filter func(pkg *load.Package) bool) error {
	disposepkg := map[string]interface{}{}
	pkgsinfo := load.Packages(pkgs)
	for _, item := range pkgsinfo {
		if err := exportPackage(item, disposepkg, dest, filter); err != nil {
			return err
		}
	}
	return nil
}

func exportPackage(pkg *load.Package, loadcache map[string]interface{}, dest string, filter func(pkg *load.Package) bool) error {
	if !filter(pkg) {
		return nil
	}
	if _, has := loadcache[pkg.ImportPath]; has {
		return nil
	}
	loadcache[pkg.ImportPath] = nil
	if err := copypackage(pkg, dest); err != nil {
		return err
	}
	for _, item := range pkg.Internal.Deps {
		if err := exportPackage(item, loadcache, dest, filter); err != nil {
			return err
		}
	}
	return nil
}

func copypackage(pkg *load.Package, dest string) error {
	files := []string{}
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.CFiles...)
	files = append(files, pkg.CXXFiles...)
	files = append(files, pkg.MFiles...)
	files = append(files, pkg.HFiles...)
	files = append(files, pkg.FFiles...)
	files = append(files, pkg.SFiles...)
	files = append(files, pkg.SwigCXXFiles...)
	files = append(files, pkg.SwigFiles...)
	files = append(files, pkg.SysoFiles...)
	for _, item := range files {
		if err := copyfile(filepath.Join(pkg.Dir, item), filepath.Join(dest, pkg.ImportPath, item)); err != nil {
			return err
		}
	}
	return nil
}

func copyfile(source, dest string) error {
	sourcef, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourcef.Close()
	destf, err := gos.Create(dest)
	if err != nil {
		return err
	}
	defer destf.Close()
	_, err = io.Copy(destf, sourcef)
	return err
}

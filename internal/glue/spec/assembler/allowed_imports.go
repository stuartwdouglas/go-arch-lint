package assembler

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/glue/yaml/spec"
	"github.com/fe3dback/go-arch-lint/internal/models"
)

type allowedImportsAssembler struct {
	rootDirectory string
	resolver      *resolver
}

func newAllowedImportsAssembler(
	rootDirectory string,
	resolver *resolver,
) *allowedImportsAssembler {
	return &allowedImportsAssembler{
		rootDirectory: rootDirectory,
		resolver:      resolver,
	}
}

func (aia *allowedImportsAssembler) assemble(
	spec *spec.ArchV1Document,
	componentNames []string,
	vendorNames []string,
) ([]models.ResolvedPath, error) {
	list := make([]models.ResolvedPath, 0)

	allowedComponents := make([]string, 0)
	allowedComponents = append(allowedComponents, componentNames...)
	allowedComponents = append(allowedComponents, spec.V1CommonComponents...)

	allowedVendors := make([]string, 0)
	allowedVendors = append(allowedVendors, vendorNames...)
	allowedVendors = append(allowedVendors, spec.V1CommonVendors...)

	for _, name := range allowedComponents {
		maskPath := spec.V1Components[name].V1LocalPath

		resolved, err := aia.resolver.resolveLocalPath(maskPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve mask '%s'", maskPath)
		}

		for _, resolvedPath := range resolved {
			list = append(list, resolvedPath)
		}
	}

	for _, name := range allowedVendors {
		importPath := spec.V1Vendors[name].V1ImportPath
		localPath := fmt.Sprintf("vendor/%s", importPath)
		absPath := fmt.Sprintf("%s/%s", aia.rootDirectory, localPath)

		list = append(list, models.ResolvedPath{
			ImportPath: importPath,
			LocalPath:  localPath,
			AbsPath:    absPath,
		})
	}

	return list, nil
}

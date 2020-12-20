package spec

import (
	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type (
	ArchV1 struct {
		V1Document  ArchV1Document
		V1Integrity []speca.Notice
	}

	ArchV1Document struct {
		reference                  models.Reference
		internalVendors            archV1InternalVendors
		internalComponents         archV1InternalComponents
		internalExclude            archV1InternalExclude
		internalExcludeFilesRegExp archV1InternalExcludeFilesRegExp
		internalCommonVendors      archV1InternalCommonVendors
		internalCommonComponents   archV1InternalCommonComponents
		internalDependencies       archV1InternalDependencies

		V1Version            speca.ReferableInt                     `yaml:"version" json:"version"`
		V1Allow              ArchV1Allow                            `yaml:"allow" json:"allow"`
		V1Vendors            map[arch.VendorName]ArchV1Vendor       `yaml:"vendors" json:"vendors"`
		V1Exclude            []speca.ReferableString                `yaml:"exclude" json:"exclude"`
		V1ExcludeFilesRegExp []speca.ReferableString                `yaml:"excludeFiles" json:"excludeFiles"`
		V1Components         map[arch.ComponentName]ArchV1Component `yaml:"components" json:"components"`
		V1Dependencies       map[arch.ComponentName]ArchV1Rules     `yaml:"deps" json:"deps"`
		V1CommonComponents   []speca.ReferableString                `yaml:"commonComponents" json:"commonComponents"`
		V1CommonVendors      []speca.ReferableString                `yaml:"commonVendors" json:"commonVendors"`
	}

	ArchV1Allow struct {
		reference models.Reference

		V1DepOnAnyVendor speca.ReferableBool `yaml:"depOnAnyVendor" json:"depOnAnyVendor"`
	}

	ArchV1Vendor struct {
		reference models.Reference

		V1ImportPath speca.ReferableString `yaml:"in" json:"in"`
	}

	ArchV1Component struct {
		reference models.Reference

		V1LocalPath speca.ReferableString `yaml:"in" json:"in"`
	}

	ArchV1Rules struct {
		reference models.Reference

		V1MayDependOn    []speca.ReferableString `yaml:"mayDependOn" json:"mayDependOn"`
		V1CanUse         []speca.ReferableString `yaml:"canUse" json:"canUse"`
		V1AnyProjectDeps speca.ReferableBool     `yaml:"anyProjectDeps" json:"anyProjectDeps"`
		V1AnyVendorDeps  speca.ReferableBool     `yaml:"anyVendorDeps" json:"anyVendorDeps"`
	}

	archV1InternalVendors struct {
		reference models.Reference
		data      map[arch.VendorName]ArchV1Vendor
	}

	archV1InternalComponents struct {
		reference models.Reference
		data      map[arch.ComponentName]ArchV1Component
	}

	archV1InternalExclude struct {
		reference models.Reference
		data      []speca.ReferableString
	}

	archV1InternalExcludeFilesRegExp struct {
		reference models.Reference
		data      []speca.ReferableString
	}

	archV1InternalCommonVendors struct {
		reference models.Reference
		data      []speca.ReferableString
	}

	archV1InternalCommonComponents struct {
		reference models.Reference
		data      []speca.ReferableString
	}

	archV1InternalDependencies struct {
		reference models.Reference
		data      map[arch.ComponentName]ArchV1Rules
	}
)

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (a ArchV1) Document() arch.Document {
	return a.V1Document
}

func (a ArchV1) Integrity() []speca.Notice {
	return a.V1Integrity
}

func (a ArchV1) applyReferences(resolver YamlSourceCodeReferenceResolver) ArchV1 {
	a.V1Document = a.V1Document.applyReferences(resolver)

	return a
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (doc ArchV1Document) Reference() models.Reference {
	return doc.reference
}

func (doc ArchV1Document) Version() speca.ReferableInt {
	return doc.V1Version
}

func (doc ArchV1Document) Options() arch.Options {
	return doc.V1Allow
}

func (doc ArchV1Document) ExcludedDirectories() arch.ExcludedDirectories {
	return doc.internalExclude
}

func (doc ArchV1Document) ExcludedFilesRegExp() arch.ExcludedFilesRegExp {
	return doc.internalExcludeFilesRegExp
}

func (doc ArchV1Document) Vendors() arch.Vendors {
	return doc.internalVendors
}

func (doc ArchV1Document) Components() arch.Components {
	return doc.internalComponents
}

func (doc ArchV1Document) CommonComponents() arch.CommonComponents {
	return doc.internalCommonComponents
}

func (doc ArchV1Document) CommonVendors() arch.CommonVendors {
	return doc.internalCommonVendors
}

func (doc ArchV1Document) Dependencies() arch.Dependencies {
	return doc.internalDependencies
}

func (doc ArchV1Document) applyReferences(resolver YamlSourceCodeReferenceResolver) ArchV1Document {
	// Version
	doc.V1Version = speca.NewReferableInt(
		doc.V1Version.Value(),
		resolver.Resolve("$.version"),
	)

	// Allow
	doc.V1Allow = doc.V1Allow.applyReferences(resolver)

	// Vendors
	doc.internalVendors = archV1InternalVendors{
		reference: resolver.Resolve("$.vendors"),
		data:      doc.V1Vendors,
	}

	// Exclude
	doc.internalExclude = archV1InternalExclude{
		reference: resolver.Resolve("$.exclude"),
		data:      doc.V1Exclude,
	}

	// ExcludeFilesRegExp
	doc.internalExcludeFilesRegExp = archV1InternalExcludeFilesRegExp{
		reference: resolver.Resolve("$.excludeFiles"),
		data:      doc.V1ExcludeFilesRegExp,
	}

	// Components
	doc.internalComponents = archV1InternalComponents{
		reference: resolver.Resolve("$.components"),
		data:      doc.V1Components,
	}

	// Dependencies
	doc.internalDependencies = archV1InternalDependencies{
		reference: resolver.Resolve("$.deps"),
		data:      doc.V1Dependencies,
	}

	// CommonComponents
	doc.internalCommonComponents = archV1InternalCommonComponents{
		reference: resolver.Resolve("$.commonComponents"),
		data:      doc.V1CommonComponents,
	}

	// CommonVendors
	doc.internalCommonVendors = archV1InternalCommonVendors{
		reference: resolver.Resolve("$.commonVendors"),
		data:      doc.V1CommonVendors,
	}

	return doc
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (opt ArchV1Allow) Reference() models.Reference {
	return opt.reference
}

func (opt ArchV1Allow) IsDependOnAnyVendor() speca.ReferableBool {
	return opt.V1DepOnAnyVendor
}

func (opt ArchV1Allow) applyReferences(resolver YamlSourceCodeReferenceResolver) ArchV1Allow {
	opt.reference = resolver.Resolve("$.allow")

	opt.V1DepOnAnyVendor = speca.NewReferableBool(
		opt.V1DepOnAnyVendor.Value(),
		resolver.Resolve("$.allow.depOnAnyVendor"),
	)

	return opt
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (v ArchV1Vendor) Reference() models.Reference {
	return v.reference
}

func (v ArchV1Vendor) ImportPath() speca.ReferableString {
	return v.V1ImportPath
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (c ArchV1Component) Reference() models.Reference {
	return c.reference
}

func (c ArchV1Component) LocalPath() speca.ReferableString {
	return c.V1LocalPath
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (rule ArchV1Rules) Reference() models.Reference {
	return rule.reference
}

func (rule ArchV1Rules) MayDependOn() []speca.ReferableString {
	return rule.V1MayDependOn
}

func (rule ArchV1Rules) CanUse() []speca.ReferableString {
	return rule.V1CanUse
}

func (rule ArchV1Rules) AnyProjectDeps() speca.ReferableBool {
	return rule.V1AnyProjectDeps
}

func (rule ArchV1Rules) AnyVendorDeps() speca.ReferableBool {
	return rule.V1AnyVendorDeps
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (a archV1InternalDependencies) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalDependencies) Map() map[arch.ComponentName]arch.DependencyRule {
	res := make(map[arch.ComponentName]arch.DependencyRule)
	for name, rules := range a.data {
		res[name] = rules
	}
	return res
}

func (a archV1InternalCommonComponents) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalCommonComponents) List() []speca.ReferableString {
	return a.data
}

func (a archV1InternalCommonVendors) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalCommonVendors) List() []speca.ReferableString {
	return a.data
}

func (a archV1InternalExcludeFilesRegExp) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalExcludeFilesRegExp) List() []speca.ReferableString {
	return a.data
}

func (a archV1InternalExclude) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalExclude) List() []speca.ReferableString {
	return a.data
}

func (a archV1InternalComponents) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalComponents) Map() map[arch.ComponentName]arch.Component {
	res := make(map[arch.ComponentName]arch.Component)
	for name, component := range a.data {
		res[name] = component
	}
	return res
}

func (a archV1InternalVendors) Reference() models.Reference {
	return a.reference
}

func (a archV1InternalVendors) Map() map[arch.VendorName]arch.Vendor {
	res := make(map[arch.VendorName]arch.Vendor)
	for name, vendor := range a.data {
		res[name] = vendor
	}
	return res
}

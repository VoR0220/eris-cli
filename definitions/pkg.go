package definitions

type PackageDefinition struct {
	Name    string   `json:"name" yaml:"name" toml:"name"`
	Package *Package `mapstructure:"eris" json:"eris" yaml:"eris" toml:"eris"`
}

type PackageManifest struct {
	// name of the package
	Name string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	// string ID of the package
	PackageID string `mapstructure:"package_id" json:"package_id" yaml:"package_id" toml:"package_id"`
	// environment variables required when running the package operations
	Environment map[string]string `mapstructure:"environment" json:"environment" yaml:"environment" toml:"environment"`
	// name of the chain to use (can utilize the $chain variable)
	ChainName string `mapstructure:"chain_name" json:"chain_name" yaml:"chain_name" toml:"chain_name"`
	// ID of the chain to use (currently this is not utilized)
	ChainID string `mapstructure:"chain_id" json:"chain_id" yaml:"chain_id" toml:"chain_id"`
	// ChainTypes the package is restricted to (currently this is not utilized)
	ChainTypes []string `mapstructure:"chain_types" json:"chain_types" yaml:"chain_types" toml:"chain_types"`
	// Dependencies to be booted before the package is ran
	Dependencies *Dependencies `mapstructure:"dependencies" json:"dependencies" yaml:"dependencies" toml:"dependencies"`

	Maintainer *Maintainer `json:"maintainer,omitempty" yaml:"maintainer,omitempty" toml:"maintainer,omitempty"`
	Location   *Location   `json:"location,omitempty" yaml:"location,omitempty" toml:"location,omitempty"`
	// AppType           *AppType    `json:"app_type,omitempty" yaml:"app_type,omitempty" toml:"app_type,omitempty"`
	Chain             *Chain
	Srvs              []*Service
	Operations        *Operation
	SkipContractsPath bool
	SkipABIPath       bool

	// from epm
	Account   string
	Jobs      []*Jobs
	Libraries map[string]string
}

type PackageLock struct {
	// The lock_file_version field defines the specification version that this document conforms to. 
	// All release lock files must include this field.
	LockFileVersion uint `mapstructure:"lock_file_version" json:"lock_file_version"`
	// The version field declares the version number of this release. 
	// This value must be included in all release lock files. 
	// This value should be conform to the semver version numbering specification.
	Version string `mapstructure:"version" json:"version"`
	// The license field declares the license under which this package is released. 
	// This value should be conform to the SPDX format. 
	// All release lock files should include this field.
	License string `mapstructure:"license" json:"license"`
	// Chain that the package is deployed on. Required for all addresses.
	// Need to figure out how to do this specifically for Eris. Should follow BIP122.
	Chain []string `mapstructure:"chain" json:"chain"`
	// The contracts field declares information about the deployed contracts included within this release.
	Contracts []*ContractInstance `mapstructure:"contracts" json:"contracts"`
}

type ContractInstance struct {
	// Name of the contract
	Name string `mapstructure:"contract_name" json:"contract_name"`
	// Address of the contract
	Address string `mapstructure:"address" json:"address"`
	// Bytecode
	Bytecode string `mapstructure:"bytecode" json:"bytecode"`
	// Abi
	Abi string `mapstructure:"abi" json:"abi"`
	// Compiler Information
	Compiler *CompilerSettings `mapstructure:"compiler" json:"compiler"`
	// Link dependencies (for libraries)
	LinkDependencies map[string]interface{} `mapstructure:"link_dependencies" json:"link_dependencies"`
}

type CompilerSettings struct {
	// Compiler version
	Version string `mapstructure:"version" json:"version"`
	// Optimize 
	Optimize bool `mapstructure:"optimize" json:"optimize"`
	// Optimization runs
	OptimizeRuns uint `mapstructure:"optimize_runs" json:"optimize_runs"`
}

func BlankPackageDefinition() *PackageDefinition {
	return &PackageDefinition{
		Package: BlankPackage(),
	}
}

func BlankPackage() *Package {
	return &Package{
		Dependencies: &Dependencies{},
		Location:     BlankLocation(),
		Operations:   BlankOperation(),
		// AppType:      BlankAppType(),
		Chain: BlankChain(),
	}
}

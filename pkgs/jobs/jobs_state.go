package jobs

// ------------------------------------------------------------------------
// State Jobs
// ------------------------------------------------------------------------

type DumpState struct {
	WithValidators bool   `mapstructure:"include-validators" json:"include-validators" yaml:"include-validators" toml:"include-validators"`
	ToIPFS         bool   `mapstructure:"to-ipfs" json:"to-ipfs" yaml:"to-ipfs" toml:"to-ipfs"`
	ToFile         bool   `mapstructure:"to-file" json:"to-file" yaml:"to-file" toml:"to-file"`
	IPFSHost       string `mapstructure:"ipfs-host" json:"ipfs-host" yaml:"ipfs-host" toml:"ipfs-host"`
	FilePath       string `mapstructure:"file" json:"file" yaml:"file" toml:"file"`
}

type RestoreState struct {
	FromIPFS bool   `mapstructure:"from-ipfs" json:"from-ipfs" yaml:"from-ipfs" toml:"from-ipfs"`
	FromFile bool   `mapstructure:"from-file" json:"from-file" yaml:"from-file" toml:"from-file"`
	IPFSHost string `mapstructure:"ipfs-host" json:"ipfs-host" yaml:"ipfs-host" toml:"ipfs-host"`
	FilePath string `mapstructure:"file" json:"file" yaml:"file" toml:"file"`
}
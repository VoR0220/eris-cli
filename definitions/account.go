package definitions

type ErisDBAccount struct {
	Name    string `mapstructure:"name" json:"name"`
	Address string `mapstructure:"address" json:"address"`
	PubKey  string `mapstructure:"" json:","`
	Tokens  int    `mapstructure:"," json:","`
	ToBond  int    `mapstructure:"," json:","`

	// [zr] from definitions/erisdb_chains.go
	//Address     string                    `json:"address"`
	Amount int `json:"amount"`
	//Name        string                    `json:"name"`
	Permissions *ErisDBAccountPermissions `json:"permissions"`

	Validator         bool
	PermissionsMap    map[string]int
	ErisDBPermissions *ErisDBAccountPermissions
	MintKey           *MintPrivValidator
}

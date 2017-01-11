package definitions

type ErisDBAccount struct {
	Name        string                    `mapstructure:"name" json:"name"`
	Address     string                    `mapstructure:"address" json:"address"`
	Amount      int                       `mapstructure:"amount" json:"amount"`
	Permissions *ErisDBAccountPermissions `mapstructure:"permissions" json:"permissions"`

	Validator         bool
	PermissionsMap    map[string]int
	ErisDBPermissions *ErisDBAccountPermissions
	MintKey           *MintPrivValidator
	PubKey            string
	ToBond            int
}

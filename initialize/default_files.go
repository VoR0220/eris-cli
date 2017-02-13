package initialize

import (
	"fmt"
	"os"
	"path"

	"github.com/eris-ltd/eris/definitions"
	"github.com/eris-ltd/eris/version"
)

const tomlHeader = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml
`

// ------------------ services ------------------

var ServiceDefinitions = []string{
	"compilers",
	"ipfs",
	"keys",
	// used by [eris chains start myChain --logrotate]
	// but its docker image is not pulled on [eris init]
	"logrotate",
}

const serviceHeader = `# For more information on configurations, see the services specification:
# https://monax.io/docs/documentation/cli/latest/services_specification/

# These fields marshal roughly into the [docker run] command, see:
# https://docs.docker.com/engine/reference/run/
`

const serviceDefinitionGeneral = `
name = "{{ .Name}}"

description = """
{{ .Description}}
"""

status = "{{ .Status}}"

[service]
image = "{{ .Service.Image}}"
data_container = {{ .Service.AutoData}}
ports = {{ .Service.Ports}}
exec_host = "{{ .Service.ExecHost}}"
volumes = {{ .Service.Volumes}}

[maintainer]
name = "{{ .Maintainer.Name}}"
email = "{{ .Maintainer.Email}}"
`

// ------------------ account types ------------------

var AccountTypeDefinitions = []string{
	"participant",
	"developer",
	"validator",
	"full",
	"root",
}

const accountTypeDefinitionGeneral = `
name = "{{ .Name}}"

description = """
{{ .Description}}
"""

typical_user = """
{{ .TypicalUser}}
"""

default_number = {{ .DefaultNumber}}

default_tokens = {{ .DefaultTokens}}

default_bond = {{ .DefaultBond}}

[perms]
root = {{index .Perms "root"}}
send = {{index .Perms "send"}}
call = {{index .Perms "call"}}
create_contract = {{index .Perms "create_contract"}}
create_account = {{index .Perms "create_account"}}
bond = {{index .Perms "bond"}}
name = {{index .Perms "name"}}
has_base = {{index .Perms "has_base"}}
set_base = {{index .Perms "set_base"}}
unset_base = {{index .Perms "unset_base"}}
set_global = {{index .Perms "set_global"}}
has_role = {{index .Perms "has_role"}}
add_role = {{index .Perms "add_role"}}
rm_role = {{index .Perms "rm_role"}}
`

// ------------------ chain types ------------------

var ChainTypeDefinitions = []string{
	"simplechain",
	"sprawlchain",
	"adminchain",
	"demochain",
	"gochain",
}

const chainTypeDefinitionGeneral = `
name = "{{ .Name}}"

description = """
{{ .Description}}
"""

[account_types]
Full = {{index .AccountTypes "Full"}}
Developer = {{index .AccountTypes "Developer"}}
Participant = {{index .AccountTypes "Participant"}}
Root = {{index .AccountTypes "Root"}}
Validator = {{index .AccountTypes "Validator"}}

[servers]

[erismint]

[tendermint]
`

func defaultServices(service string) *definitions.ServiceDefinition {

	serviceDefinition := definitions.BlankServiceDefinition()

	serviceDefinition.Maintainer.Name = "Monax Industries"
	serviceDefinition.Maintainer.Email = "support@monax.io"

	switch service {
	case "keys":

		serviceDefinition.Name = "keys"
		serviceDefinition.Description = `Eris keys is meant for quick prototyping. You must replace it with a hardened key signing daemon to use in production. Eris does not intend to harden this for production, but rather will keep it as a rapid prototyping server.

This service is usually linked to a chain and/or an application. Its functionality is wrapped by [eris keys].`
		serviceDefinition.Status = "unfit for production"
		serviceDefinition.Service.Image = path.Join(version.DefaultRegistry, version.ImageKeys)
		serviceDefinition.Service.AutoData = true
		serviceDefinition.Service.Ports = []string{`"4767:4767"`} // XXX these exposed ports are a gaping security flaw

	case "compilers":

		serviceDefinition.Name = "compilers"
		serviceDefinition.Description = `Monax's Solidity Compiler Server.

This eris service compiles smart contract languages.`
		serviceDefinition.Status = "beta"
		serviceDefinition.Service.Image = path.Join(version.DefaultRegistry, version.ImageCompilers)
		serviceDefinition.Service.AutoData = true
		//serviceDefinition.Service.Ports = []string{`"9090:9090"`}

	case "ipfs":

		port_to_use := os.Getenv("ERIS_CLI_TESTS_PORT")
		if port_to_use == "" {
			port_to_use = "8080"
		}
		serviceDefinition.Name = "ipfs"
		serviceDefinition.Description = `IPFS is The Permanent Web: A new peer-to-peer hypermedia protocol. IPFS uses content-based addressing versus http's location-based addressing.

This eris service is all but essential as part of the eris tool. The [eris files] relies upon this running service.`
		serviceDefinition.Status = "alpha"
		serviceDefinition.Service.Image = path.Join(version.DefaultRegistry, version.ImageIPFS)
		serviceDefinition.Service.AutoData = true
		serviceDefinition.Service.Ports = []string{`"4001:4001"`, `"5001:5001"`, `"` + port_to_use + `:` + port_to_use + `"`}

	case "logrotate":

		serviceDefinition.Name = "logrotate"
		serviceDefinition.Description = `Truncates docker container logs when the grow in size.

This eris service can also be run by adding the [--logrotate] flag on [eris chains start]

It is essential for long-running chain nodes.

Alternatively, use logspout to pipe logs to a service of you choosing"`
		serviceDefinition.Status = "ready"
		serviceDefinition.Service.Image = "tutum/logrotate"
		serviceDefinition.Service.AutoData = false
		serviceDefinition.Service.Volumes = []string{`"/var/lib/docker/containers:/var/lib/docker/containers:rw"`}

	default:
		panic(fmt.Errorf("not allowed"))
	}

	return serviceDefinition

}

func defaultAccountTypes(accountType string) *definitions.ErisDBAccountType {
	accountTypeDefinition := definitions.BlankAccountType()

	switch accountType {
	case "participant":

		accountTypeDefinition.Name = "Participant"
		accountTypeDefinition.Description = `Users who have a key which is registered with participant privileges can send
tokens; call contracts; and use the name registry`
		accountTypeDefinition.TypicalUser = `Generally the number of participants in your chain who do not need elevated
privileges should be given these keys.

Usually this group will have the most number of keys of all of the groups.`
		accountTypeDefinition.DefaultNumber = 25
		accountTypeDefinition.DefaultTokens = 9999999999
		accountTypeDefinition.DefaultBond = 0
		accountTypeDefinition.Perms = map[string]int{
			"root":            0,
			"send":            1,
			"call":            1,
			"create_contract": 0,
			"create_account":  0,
			"bond":            0,
			"name":            1,
			"has_base":        0,
			"set_base":        0,
			"unset_base":      0,
			"set_global":      0,
			"has_role":        1,
			"add_role":        0,
			"rm_role":         0,
		}

	case "developer":

		accountTypeDefinition.Name = "Developer"
		accountTypeDefinition.Description = `Users who have a key which is registered with developer privileges can send
tokens; call contracts; create contracts; create accounts; use the name registry;
and modify account's roles.`
		accountTypeDefinition.TypicalUser = `Generally the development team seeking to build the application on top of the
given chain would be within the group. If this is a multi organizational
chain then developers from each of the stakeholders should generally be registered
within this group, although that design is up to you.`
		accountTypeDefinition.DefaultNumber = 6
		accountTypeDefinition.DefaultTokens = 9999999999
		accountTypeDefinition.DefaultBond = 0
		accountTypeDefinition.Perms = map[string]int{
			"root":            0,
			"send":            1,
			"call":            1,
			"create_contract": 1,
			"create_account":  1,
			"bond":            0,
			"name":            1,
			"has_base":        0,
			"set_base":        0,
			"unset_base":      0,
			"set_global":      0,
			"has_role":        1,
			"add_role":        1,
			"rm_role":         1,
		}

	case "validator":

		accountTypeDefinition.Name = "Validator"
		accountTypeDefinition.Description = `Users who have a key which is registered with validator privileges can
only post a bond and begin validating the chain. This is the only privilege
this account group gets.`
		accountTypeDefinition.TypicalUser = `Generally the marmots recommend that you put your validator nodes onto a cloud
(IaaS) provider so that they will be always running.

We also recommend that if you are in a multi organizational chain then you would
have some separation of the validators to be ran by the different organizations
in the system.`
		accountTypeDefinition.DefaultNumber = 7
		accountTypeDefinition.DefaultTokens = 9999999999
		accountTypeDefinition.DefaultBond = 9999999998
		accountTypeDefinition.Perms = map[string]int{
			"root":            0,
			"send":            0,
			"call":            0,
			"create_contract": 0,
			"create_account":  0,
			"bond":            1,
			"name":            0,
			"has_base":        0,
			"set_base":        0,
			"unset_base":      0,
			"set_global":      0,
			"has_role":        0,
			"add_role":        0,
			"rm_role":         0,
		}

	case "full":

		accountTypeDefinition.Name = "Full"
		accountTypeDefinition.Description = `Users who have a key which is registered with root privileges can do everything
on the chain. They have all of the permissions possible. These users are also
bonded at the genesis time, so these should be used only for simple chains with
a few nodes who will be on during the prototyping session.`
		accountTypeDefinition.TypicalUser = `If you are making a small chain just to play around then usually you would
give all of the accounts needed for your experiment full accounts.

If you are making a more complex chain, don't use this account type.`
		accountTypeDefinition.DefaultNumber = 1
		accountTypeDefinition.DefaultTokens = 99999999999999
		accountTypeDefinition.DefaultBond = 9999999999
		accountTypeDefinition.Perms = map[string]int{
			"root":            1,
			"send":            1,
			"call":            1,
			"create_contract": 1,
			"create_account":  1,
			"bond":            1,
			"name":            1,
			"has_base":        1,
			"set_base":        1,
			"unset_base":      1,
			"set_global":      1,
			"has_role":        1,
			"add_role":        1,
			"rm_role":         1,
		}

	case "root":

		accountTypeDefinition.Name = "Root"
		accountTypeDefinition.Description = `Users who have a key which is registered with root privileges can do everything
on the chain. They have all of the permissions possible.`
		accountTypeDefinition.TypicalUser = `If you are making a small chain just to play around then usually you would
give all of the accounts needed for your experiment root privileges (unless you
were testing different) privilege types.

If you are making a more complex chain, then you would usually have a few
keys which have registered root permissions and as such will act in a capacity
similar to a network administrator in other data management situations.`
		accountTypeDefinition.DefaultNumber = 3
		accountTypeDefinition.DefaultTokens = 9999999999
		accountTypeDefinition.DefaultBond = 0
		accountTypeDefinition.Perms = map[string]int{
			"root":            1,
			"send":            1,
			"call":            1,
			"create_contract": 1,
			"create_account":  1,
			"bond":            1,
			"name":            1,
			"has_base":        1,
			"set_base":        1,
			"unset_base":      1,
			"set_global":      1,
			"has_role":        1,
			"add_role":        1,
			"rm_role":         1,
		}

	default:
		panic(fmt.Errorf("not allowed"))
	}

	return accountTypeDefinition
}

func defaultChainTypes(chainType string) *definitions.ChainType {
	chainTypeDefinition := definitions.BlankChainType()

	switch chainType {
	case "simplechain":

		chainTypeDefinition.Name = "simplechain"
		chainTypeDefinition.Description = `A simple chain type will build a single node chain. This chain type is usefulfor quick and easy prototyping. It should not be used for anything more than the most simple prototyping as it only has one node and the keys to that node could get lost or compromised and the chain would thereafter become useless.`
		chainTypeDefinition.AccountTypes = map[string]int{
			"Full":        1,
			"Developer":   0,
			"Participant": 0,
			"Root":        0,
			"Validator":   0,
		}

	case "sprawlchain":

		chainTypeDefinition.Name = "sprawlchain"
		chainTypeDefinition.Description = `A sprawlchain type has a little bit of everything. Modify as necessary for your ecosystem application. Will tolerate three nodes down. As with other chains, Validator accounts ought to go on the cloud. No Full accounts are provided since these are prefered for quick development only.`
		chainTypeDefinition.AccountTypes = map[string]int{
			"Full":        0,
			"Developer":   10,
			"Participant": 20,
			"Root":        3,
			"Validator":   10,
		}

	case "adminchain":

		chainTypeDefinition.Name = "adminchain"
		chainTypeDefinition.Description = `An adminchain type has settings for prototyping a larger chain from a sysadmin point of view. With four Validator and three Full account_types, at minimum of five nodes must be up for consensus to happen. This account combination is what we use to test long running chains on our CI system.`
		chainTypeDefinition.AccountTypes = map[string]int{
			"Full":        3,
			"Developer":   1,
			"Participant": 1,
			"Root":        1,
			"Validator":   4,
		}

	case "demochain":

		chainTypeDefinition.Name = "demochain"
		chainTypeDefinition.Description = `A demo chain is useful for setting up proof of concept demonstration chains. It is a quick and easy way to have multi-validator, multi-developer, multi-participant chains set up for your application. This chain will tolerate 2 validators becoming byzantine or going off-line while still moving forward. You should utilize 7 different cloud instances and deploy one of the validator chain directories to each.`
		chainTypeDefinition.AccountTypes = map[string]int{
			"Full":        0,
			"Developer":   5,
			"Participant": 20,
			"Root":        3,
			"Validator":   7,
		}

	case "gochain":

		chainTypeDefinition.Name = "gochain"
		chainTypeDefinition.Description = `A gochain type will build a three node chain. It is a quick and easy way to get started with a multi-validator chain. The Full account_type includes validator and deploy permissions, allowing for experimentation with setting up a chain and halting it by taking down a single node. This Full account should be deployed on your local machine and cloud nodes should have only Validator accounts. Use for prototyping only.`
		chainTypeDefinition.AccountTypes = map[string]int{
			"Full":        1,
			"Developer":   0,
			"Participant": 0,
			"Root":        0,
			"Validator":   2,
		}

	default:
		panic(fmt.Errorf("not allowed"))
	}
	return chainTypeDefinition
}

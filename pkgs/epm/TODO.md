TODO

Accept JSON --json jsonabi

abi storage

restructure core.go to be more useful.

hash flag

abi indexing
	- Resolve indexfile Key pair
	- ReadIndex (return index object)
	- WriteIndex (dowit)
	- EditIndex (set key value)

import abi

accept chainid, contract pair (chainid should be able to be an ENV)

=======================FINISHED LINE=====================

pack server

--- Push point

Tests

index functions added to server?

--- Push point

Cleanup + Documentation

--- Rework/Improvements ---

use Cobra?

utils--v
better file handling (if absolute path don't prepend cwd)

add completion functionality (chainid, hash ...)


________________________________________________________________
Organization:

/cmd/eris-abi/main.go - app creation
/cmd/eris-abi/cli.go - cli functions (calling on eris-abi)

/core.go - Eris-specific ABI management system stuff (package ebi)
/storage.go - Read and write functions for eris ABI storage system
/directories.go - Outlines Directory Tree structure (and creation)


_______________________________________________________________

Eris ABI Folder structure:

root: .eris/abi/

abi file stored in root named as sha256 hash of abi contents

indexes: .eris/abi/index

Index files are json formatted. the name of the index is the "outer key" in the case of chainid/contractaddr pairs. Index name should be chainid the internal mapping would then map contract addr to abi hash.  

//Directory Structure
//.eris
//   +-abi
//      +-index
//      |    + abi indexing jsons
//      |
//      +-raw
//           + abi files (hash named)

_______________________________________________________________
Core.go

Core.go contains the interface functions between abi package, the cli, and storage.go

It needs to be able to:

1) take abiData []byte output an ABI object
2) take ABI object convert to []byte? not needed?
3) Pack Transaction data into and ABI and return transaction ?put in abi.go?

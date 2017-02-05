package compilers

type Compiler interface {
	Compile(files []string, version string) (Return, error)
}

//Practicing inheritance, this struct gives us access to all types of returns
type Return struct {
	SolcReturn
	//Enter your return struct here...
}

type DefaultCompilers struct {
	defaultLanguage string `toml:"defaultLanguage"`
	details         map[string]CompilerDetails
}

type CompilerDetails struct {
	Name              string   `toml:"name"`
	DockerRepo        string   `toml:"docker_repo"`
	VersionsAvailable []string `toml:"versions_available"`
}

package brew

type DependencyType string

const (
	Brew     DependencyType = "BREW"
	BrewCask DependencyType = "BREW_CASK"
	Custom   DependencyType = "CUSTOM"
)

type DependencyRequirement struct {
	Name     string `json:"name"`
	Cask     string `json:"cask"`
	Download string `json:"download"`
}

type Dependency struct {
	Type                    DependencyType          `json:"type"`
	Name                    string                  `json:"name"`
	Aliases                 []string                `json:"aliases"`
	FullName                string                  `json:"full_name"`
	Description             string                  `json:"desc"`
	HomePage                string                  `json:"homepage"`
	Requirements            []DependencyRequirement `json:"requirements"`
	Caveats                 string                  `json:"caveats"`
	ConflictsWith           []string                `json:"conflicts_with"`
	Dependencies            []string                `json:"dependencies"`
	RecommendedDependencies []string                `json:"recommended_dependencies"`
	OptionalDependencies    []string                `json:"optional_dependencies"`
	BuildDependencies       []string                `json:"build_dependencies"`
}

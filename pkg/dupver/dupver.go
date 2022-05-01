package dupver

// MajorVersion is the major version of dupver that creates breaking changes
// MinorVersion is the major version of dupver that adds new functionality without breaking changes
// PatchVersion is the patch version of dupver that does not add new functionality
// RepoMajorVersion is the major version of the repository format that creates breaking changes
// RepoMinorVersion is the minor version of the repository format that does not create breaking changes
// PrefsMajorVersion is the major version of the global preferences format that creates breaking changes
// PrefsMinorVersion is the major version of the global preferences format that creates breaking changes
const (
	MajorVersion       = 4
	MinorVersion       = 0
	PatchVersion       = 0
	RepoMajorVersion   = 4
	RepoMinorVersion   = 0
	PrefsMajorVersion  = 2
	PrefsMinorVersion  = 0
)

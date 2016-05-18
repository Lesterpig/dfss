package entities

import (
	"dfss/dfssc/security"
)

// AuthContainer is global for performance reasons; singleton is not a problem.
// This variable should be loaded by dfsst/server package.
var AuthContainer *security.AuthContainer

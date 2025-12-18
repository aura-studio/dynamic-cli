package check

import (
	"github.com/aura-studio/dynamic-cli/config"
)

func CheckForProcedure(proc config.Procedure) {
	checker := NewChecker(proc)
	checker.Run()
}

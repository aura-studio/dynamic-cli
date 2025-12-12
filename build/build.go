package build

import "github.com/aura-studio/dynamic-cli/config"

// BuildForProcedure executes build operations for a given procedure.
// Stub implementation; integration will follow.
func BuildForProcedure(proc config.Procedure) {
	// Compose RenderData from Procedure
	// Name: namespace_package_version with underscores
	name := proc.Target.Namespace + "_" + proc.Target.Package + "_" + proc.Target.Version
	// Module: Source.Module
	pkg := proc.Source.Module
	// Version: Source.Version
	ver := proc.Source.Version
	// House: Warehouse.Local
	house := proc.Warehouse.Local
	// Environment: Toolchain.OS_Toolchain.Arch_Toolchain.Compiler_Toolchain.Variant
	env := proc.Toolchain.OS + "_" + proc.Toolchain.Arch + "_" + proc.Toolchain.Compiler + "_" + proc.Toolchain.Variant
	// Variant: Toolchain.Variant
	variant := proc.Toolchain.Variant
	// Dir: House/Environment/Name
	dir := house + "/" + env + "/" + name

	rd := &RenderData{
		Name:        name,
		Module:      pkg,
		Version:     ver,
		House:       house,
		Environment: env,
		Variant:     variant,
		Dir:         dir,
	}

	// Execute build using existing builder
	New(rd).Build()
}

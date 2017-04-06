package main

import (
	"golang/common"
	"golang/finalize"
	"os"

	"github.com/cloudfoundry/libbuildpack"
)

func main() {
	stager, err := libbuildpack.NewStager(os.Args[1:], libbuildpack.NewLogger())
	err = stager.CheckBuildpackValid()
	if err != nil {
		os.Exit(10)
	}

	err = libbuildpack.SetEnvironmentFromSupply(stager.DepsDir)
	if err != nil {
		stager.Log.Error("Unable to setup environment variables: %s", err.Error())
		os.Exit(11)
	}

	err = libbuildpack.WriteProfileDFromSupply(stager.DepsDir, stager.BuildDir)
	if err != nil {
		stager.Log.Error("Unable to write .profile.d supply script: %s", err.Error())
		os.Exit(12)
	}

	var godep common.Godep

	vendorTool, err := common.SelectVendorTool(stager, &godep)
	if err != nil {
		stager.Log.Error("Unable to select Go vendor tool: %s", err.Error())
		os.Exit(14)
	}

	goVersion, err := common.SelectGoVersion(stager, vendorTool, godep)
	if err != nil {
		stager.Log.Error("Unable to select Go version: %s", err.Error())
		os.Exit(15)
	}

	gf := finalize.Finalizer{
		Stager:     stager,
		Godep:      godep,
		GoVersion:  goVersion,
		VendorTool: vendorTool,
	}

	err = finalize.Run(&gf)
	if err != nil {
		os.Exit(16)
	}

	err = libbuildpack.RunAfterCompile(stager)
	if err != nil {
		stager.Log.Error("After Compile: %s", err.Error())
		os.Exit(17)
	}

	stager.StagingComplete()
}
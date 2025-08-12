package detectors

func RunAll(path string) {
	DetectUnusedParams(path)
	DetectUnusedConstants(path)
	DetectUnusedMessages(path)
	DetectUnDefinedMessageKeys(path)
	DetectCapitalVars(path)
	DetectDeprecatedPackages(path)
	DetectUnexportedFuncNames(path)
	RunDeadCode(path)
	// DetectFuncWithManyParams(path)
}

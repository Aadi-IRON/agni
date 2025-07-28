package detectors

func RunAll(path string) {
	DetectUnusedParams(path)
	DetectUnusedConstants(path)
	DetectUnusedMessages(path)
	DetectUnDefinedMessageKeys(path)
	DetectCapitalVars(path)
	RunDeadCode(path)
	// DetectFuncWithManyParams(path)
}

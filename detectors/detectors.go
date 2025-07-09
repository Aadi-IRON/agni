package detectors

func RunAll(path string) {
	DetectUnusedParams(path)
	DetectUnusedConstants(path)
	DetectUnusedMessages(path)
	RunDeadCode(path)
	// DetectFuncWithManyParams(path)
}

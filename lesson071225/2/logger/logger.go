package logger

type LoggerApi struct {
}

func (l LoggerApi) LogInfo(message IStringer) {
	// Implementation for logging info messages
}

func (l LoggerApi) LogWarning(message IStringer) {
	// Implementation for logging warning messages
}

func (l LoggerApi) LogError(message IStringer) {
	// Implementation for logging error messages
}

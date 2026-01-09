package logger

type IStringer interface {
	String() string
}

type ILogger interface {
	LogInfo(message IStringer)
	LogWarning(message IStringer)
	LogError(message IStringer)
}

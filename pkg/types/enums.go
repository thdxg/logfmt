package types

type LevelFormat string

const (
	LevelFormatFull  LevelFormat = "full"
	LevelFormatShort LevelFormat = "short"
	LevelFormatTiny  LevelFormat = "tiny"
)

func (f LevelFormat) IsValid() bool {
	switch f {
	case LevelFormatFull, LevelFormatShort, LevelFormatTiny:
		return true
	}
	return false
}

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
)

package logging

import "log/slog"


var _LEVEL_BY_VERBOSITY = map[bool]slog.Level{
	false: slog.LevelInfo,
	true: slog.LevelDebug,
}

func Configure(verbose bool) {
	var handler = _Handler{
		level: _LEVEL_BY_VERBOSITY[verbose],
	}
	var logger = slog.New(handler)
	slog.SetDefault(logger)
}

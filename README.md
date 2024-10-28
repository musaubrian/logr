# Logr

This is very much not to be used by anyone other than me.
**slog** is a much much better option for you

But if you want to try it out:

```go
// [LEVEL] <file:line_no> <time:optional> <message>

//...
log := logr.New(logr.LevelInfo)
// Defaults to time in the format: 15:04:06

log.Info("Hello world")
log.Debug("Hello world")
log.Warn("Hello world")
log.Error("Hello world")
log.Logf(logr.LevelInfo,"Some value %d", 45+90)

//...
```

## TODO

- [x] Align logs
- [ ] save to file?


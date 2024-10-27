# Logr

This is very much not to be used by anyone other than me.
**slog** is a much much better option for you

But if you want to try it out:

```go

//...
log := logr.New(logr.LevelInfo)
log.Info("Hello world")
log.Debug("Hello world")
log.Warn("Hello world")
log.Error("Hello world")
log.Logf(logr.LevelInfo,"Some value %d", 45+90)

//...
```

## TODO

- [ ] Align logs
- [ ] save to file?


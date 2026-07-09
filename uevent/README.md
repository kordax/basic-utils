# uevent

Channel watcher helpers.

Includes parallel and broadcast watchers for consuming values from channels and dispatching events to handlers.

## Example

```go
events := make(chan string)
watcher := uevent.NewParallelWatcher(events, func(ctx context.Context, value string) {
	handle(value)
})
watcher.Watch(ctx)
```

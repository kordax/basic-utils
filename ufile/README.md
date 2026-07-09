# ufile

File-system helpers for common read/write/list workflows.

Includes file existence checks, safe creation, JSON read/write helpers, directory creation, recursive walking, and panic-on-error read/write helpers.

## Example

```go
err := ufile.WriteJSON("config.json", config, 0644)
config, err := ufile.ReadJSON[Config]("config.json")
```

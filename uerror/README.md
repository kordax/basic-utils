# uerror

Error handling helpers.

Currently includes `Must`, a convenience helper for panic-on-error flows where failing fast is expected.

## Example

```go
data, err := os.ReadFile("config.json")
uerror.Must(data, err)
```

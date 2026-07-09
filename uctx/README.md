# uctx

Small context-related helpers and global context access.

This package is intentionally minimal.

## Example

```go
ctx := uctx.GetContext()
ctx.Set("request_id", requestID)
requestID := ctx.Get("request_id").(string)
```

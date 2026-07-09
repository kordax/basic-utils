# uonce

Helpers for one-time execution semantics.

Use this package when a value or function must be initialized or executed only once.

## Example

```go
getClient := uonce.Once(func() *Client {
	return NewClient()
})

client := getClient()
```

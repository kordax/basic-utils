# usrlz

Serialization helpers.

Currently focused on converting supported value types into byte representations.

## Example

```go
type Header struct {
	Version uint16
	Flags   uint16
}

bytes := usrlz.ToBytes(&Header{Version: 1, Flags: 2})
```

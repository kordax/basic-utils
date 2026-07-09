# usql

SQL helper types and functions.

Use this package for common SQL value handling that is shared across services.

## Example

```go
name := usql.NullString(input.Name)
createdAt := usql.NullTime(time.Now())
```

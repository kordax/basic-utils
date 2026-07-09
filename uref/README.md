# uref

Reference helpers for creating pointers to values and working with optional references.

## Example

```go
timeout := uref.Ref(30 * time.Second)
same := uref.Compare(timeout, uref.Ref(30*time.Second))
```

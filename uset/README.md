# uset

Set implementations and set algebra helpers.

Includes hash sets, synchronized sets, sharded concurrent sets, ordered sets, comparable-key sets, and helpers for union, intersection, difference, subset, and superset checks.

## Example

```go
set := uset.NewHashSet("read", "write")
set.Add("admin")
allowed := set.Contains("read")
```

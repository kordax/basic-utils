# umap

Map utilities and multi-map implementations.

Includes contains/equality helpers, copy/merge, keys/values, map filtering, key/value mapping, inversion, and several multi-map variants.

## Example

```go
usersByID := map[int]string{1: "Ann", 2: "Bob"}
names := umap.Values(usersByID)
active := umap.Filter(usersByID, func(id int, name string) bool {
	return id == 1
})
```

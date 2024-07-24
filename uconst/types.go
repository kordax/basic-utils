/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uconst

/*
Comparable entity
*/
type Comparable interface {
	Equals(other Comparable) bool
}

/*
UniqueKey specifies an abstract key with an ability to provide hash.
*/
type UniqueKey[T comparable] interface {
	Key() T // Key should return a unique item key. It can be a hash or just an index.
}

/*
Unique specifies an abstract key with an ability to provide hash.
Comparable interface can be issued to mitigate potential collisions, e.g. in hashmaps or other implementations.
*/
type Unique interface {
	Comparable
	UniqueKey[int64]
}

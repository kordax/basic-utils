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
Unique specifies an abstract key with an ability to provide hash.
*/
type Unique interface {
	Comparable
	Key() int64 // Key should return a unique item key. It can be a hash or just an index.
}

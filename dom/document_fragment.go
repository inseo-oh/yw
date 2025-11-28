package dom

type DocumentFragment interface {
	Node
	Host() Node
}

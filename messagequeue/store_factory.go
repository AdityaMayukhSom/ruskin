package messagequeue

type StoreFactory interface {
	Produce(string) Store
}

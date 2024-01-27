package messagequeue

type Store interface {
	Insert([]byte) (int, error)
	Extract(int) ([]byte, error)
	ExtractLatest() ([]byte, error)
}

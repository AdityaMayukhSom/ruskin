package store_relay_connector

/*
This interface acts as a bridge between some component and a store, where
data fetching is required, but the protocol might change or the underlying
implementation might differ. Works based on pull mechanism.
*/
type StoreConnector interface {
	Fetch(int) ([]byte, error)
	FetchAll() ([][]byte, error)
	FetchLatest() ([]byte, error)
}

type RelayConnector interface {
}

package load

type LoadDistributor struct {
	PartitionsCount int
	Partitions      []LoadPartition
}

package replica

type LoadDistributor struct {
	PartitionsCount int
	Partitions      []LoadPartition
}

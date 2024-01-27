package load

import "fmt"

type LoadDistributor struct {
	PartitionsCount int
	Partitions      []LoadPartition
}

func (ld *LoadDistributor) Listen() {
	agg := make(chan string)
	for _, ch := range chans {
		go func(c chan string) {
			for msg := range c {
				agg <- msg
			}
		}(ch)
	}

	select {
	case msg <- agg:
		fmt.Println("received ", msg)
	}
}

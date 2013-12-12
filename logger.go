package collector

import (
	"fmt"
)

func PrintWriterStats(elapsed int, writer Writer) {
	created, failed, transferred := writer.GetCounters()
	writer.ResetCounters()

	logFormat := "Created %d document(s), Failed %d times(s), %g"

	rate := float64(transferred) / 1000 / float64(elapsed)
	fmt.Println(fmt.Sprintf(logFormat, created, failed, rate))
}

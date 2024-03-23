package testing

import "testing"

/*
 * node api testing
 */

func BenchmarkAddNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mc.AddNode(nodeCfg)
	}
}
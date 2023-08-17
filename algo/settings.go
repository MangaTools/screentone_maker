package algo

type ClusterSettings struct {
	Size        int
	DotSettings DotSettings
}

type DotSettings struct {
	MinValue, MaxValue byte
	Size               int
}

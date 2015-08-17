package invoker

// encryption algorithms
type EncAlgorithm string
type Compress int

const (
	namespace = "blinker/invoker"

	A128GCM = EncAlgorithm("A128GCM")
	A192GCM = EncAlgorithm("A192GCM")
	A256GCM = EncAlgorithm("A256GCM")
)

const (
	COMPRESS_NONE = Compress(iota)
	COMPRESS_DEF
)

var CompressAlgorithm = COMPRESS_DEF

package pagetoken

type Checksumer interface {
	Checksum() (uint32, error)
}

type Encodable interface {
	Params() []string
}

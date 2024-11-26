package types

const (
	// ModuleName defines the module name
	ModuleName = "optio"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_optio"
)

var (
	ParamsKey = []byte("p_optio")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

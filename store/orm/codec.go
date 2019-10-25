package orm

import "github.com/cosmos/cosmos-sdk/codec"

// module codec
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for
// governance.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MultiRef{}, "cosmos-sdk/MultiRef", nil)
	cdc.RegisterConcrete(SimpleObj{}, "cosmos-sdk/SimpleObject", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}

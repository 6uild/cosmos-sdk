// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/cosmos/cosmos-sdk/x/nft/internal/keeper
// ALIASGEN: github.com/cosmos/cosmos-sdk/x/nft/internal/types
package nft

import (
	"github.com/cosmos/cosmos-sdk/x/nft/internal/keeper"
	"github.com/cosmos/cosmos-sdk/x/nft/internal/types"
)

const (
	QuerySupply           = keeper.QuerySupply
	QueryOwner            = keeper.QueryOwner
	QueryOwnerByDenom     = keeper.QueryOwnerByDenom
	QueryCollection       = keeper.QueryCollection
	QueryDenoms           = keeper.QueryDenoms
	QueryNFT              = keeper.QueryNFT
	DefaultCodespace      = types.DefaultCodespace
	CodeInvalidCollection = types.CodeInvalidCollection
	CodeUnknownCollection = types.CodeUnknownCollection
	CodeInvalidNFT        = types.CodeInvalidNFT
	CodeUnknownNFT        = types.CodeUnknownNFT
	CodeNFTAlreadyExists  = types.CodeNFTAlreadyExists
	CodeEmptyMetadata     = types.CodeEmptyMetadata
	ModuleName            = types.ModuleName
	StoreKey              = types.StoreKey
	QuerierRoute          = types.QuerierRoute
	RouterKey             = types.RouterKey
)

var (
	// functions aliases
	RegisterInvariants       = keeper.RegisterInvariants
	AllInvariants            = keeper.AllInvariants
	SupplyInvariant          = keeper.SupplyInvariant
	NewKeeper                = keeper.NewKeeper
	NewQuerier               = keeper.NewQuerier
	RegisterCodec            = types.RegisterCodec
	NewCollection            = types.NewCollection
	EmptyCollection          = types.EmptyCollection
	NewCollections           = types.NewCollections
	ErrInvalidCollection     = types.ErrInvalidCollection
	ErrUnknownCollection     = types.ErrUnknownCollection
	ErrInvalidNFT            = types.ErrInvalidNFT
	ErrNFTAlreadyExists      = types.ErrNFTAlreadyExists
	ErrUnknownNFT            = types.ErrUnknownNFT
	ErrEmptyMetadata         = types.ErrEmptyMetadata
	NewGenesisState          = types.NewGenesisState
	DefaultGenesisState      = types.DefaultGenesisState
	ValidateGenesis          = types.ValidateGenesis
	GetCollectionKey         = types.GetCollectionKey
	SplitOwnerKey            = types.SplitOwnerKey
	GetOwnersKey             = types.GetOwnersKey
	GetOwnerKey              = types.GetOwnerKey
	NewMsgTransferNFT        = types.NewMsgTransferNFT
	NewMsgEditNFTMetadata    = types.NewMsgEditNFTMetadata
	NewMsgMintNFT            = types.NewMsgMintNFT
	NewMsgBurnNFT            = types.NewMsgBurnNFT
	NewBaseNFT               = types.NewBaseNFT
	NewNFTs                  = types.NewNFTs
	NewIDCollection          = types.NewIDCollection
	NewOwner                 = types.NewOwner
	NewQueryCollectionParams = types.NewQueryCollectionParams
	NewQueryBalanceParams    = types.NewQueryBalanceParams
	NewQueryNFTParams        = types.NewQueryNFTParams

	// variable aliases
	ModuleCdc                = types.ModuleCdc
	EventTypeTransfer        = types.EventTypeTransfer
	EventTypeEditNFTMetadata = types.EventTypeEditNFTMetadata
	EventTypeMintNFT         = types.EventTypeMintNFT
	EventTypeBurnNFT         = types.EventTypeBurnNFT
	AttributeValueCategory   = types.AttributeValueCategory
	AttributeKeySender       = types.AttributeKeySender
	AttributeKeyRecipient    = types.AttributeKeyRecipient
	AttributeKeyOwner        = types.AttributeKeyOwner
	AttributeKeyNFTID        = types.AttributeKeyNFTID
	AttributeKeyDenom        = types.AttributeKeyDenom
	CollectionsKeyPrefix     = types.CollectionsKeyPrefix
	OwnersKeyPrefix          = types.OwnersKeyPrefix
)

type (
	Keeper                = keeper.Keeper
	Collection            = types.Collection
	Collections           = types.Collections
	CollectionJSON        = types.CollectionJSON
	CodeType              = types.CodeType
	GenesisState          = types.GenesisState
	MsgTransferNFT        = types.MsgTransferNFT
	MsgEditNFTMetadata    = types.MsgEditNFTMetadata
	MsgMintNFT            = types.MsgMintNFT
	MsgBurnNFT            = types.MsgBurnNFT
	NFT                   = types.NFT
	BaseNFT               = types.BaseNFT
	NFTs                  = types.NFTs
	NFTJSON               = types.NFTJSON
	IDCollection          = types.IDCollection
	IDCollections         = types.IDCollections
	Owner                 = types.Owner
	QueryCollectionParams = types.QueryCollectionParams
	QueryBalanceParams    = types.QueryBalanceParams
	QueryNFTParams        = types.QueryNFTParams
)

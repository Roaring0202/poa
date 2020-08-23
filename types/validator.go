package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
)

// Validator
type Validator struct {
	OperatorAddress sdk.ValAddress `json:"operator_address"`
	ConsensusPubkey string         `json:"consensus_pubkey"`
	Description     Description    `json:"description"`
}

func NewValidator(operator sdk.ValAddress, pubKey crypto.PubKey, description Description) Validator {
	var pkStr string
	if pubKey != nil {
		pkStr = sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubKey)
	}

	return Validator{
		OperatorAddress: operator,
		ConsensusPubkey: pkStr,
		Description:     description,
	}
}

func (v Validator) GetOperator() sdk.ValAddress {
	return v.OperatorAddress
}

func (v Validator) GetConsPubKey() crypto.PubKey {
	return sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, v.ConsensusPubkey)
}

func (v Validator) GetConsAddr() sdk.ConsAddress {
	return sdk.ConsAddress(v.GetConsPubKey().Address())
}

// Get a ABCI validator update object from the validator
func (v Validator) ABCIValidatorUpdateAppend() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: types.TM2PB.PubKey(v.GetConsPubKey()),
		Power:  1, // Same weight for all the validators
	}
}

// Get a ABCI validator update with no voting power from the validator
func (v Validator) ABCIValidatorUpdateRemove() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: types.TM2PB.PubKey(v.GetConsPubKey()),
		Power:  0,
	}
}

// Validator encoding functions
func MustMarshalValidator(cdc codec.Codec, validator Validator) []byte {
	return cdc.MustMarshalBinaryBare(&validator)
}

func MustUnmarshalValidator(cdc codec.Codec, value []byte) Validator {
	validator, err := UnmarshalValidator(cdc, value)
	if err != nil {
		panic(err)
	}

	return validator
}

func UnmarshalValidator(cdc codec.Codec, value []byte) (v Validator, err error) {
	err = cdc.UnmarshalBinaryBare(value, &v)
	return v, err
}

// Description defines a validator description
type Description struct {
	Moniker         string
	Identity        string
	Website         string
	SecurityContact string
	Details         string
}

// Create a new Description
func NewDescription(moniker, identity, website, securityContact, details string) Description {
	return Description{
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		SecurityContact: securityContact,
		Details:         details,
	}
}

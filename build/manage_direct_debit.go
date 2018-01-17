package build

import (
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
)

// ChangeTrust groups the creation of a new ChangeTrustBuilder with a call to Mutate.
func ManageDirectDebit(muts ...interface{}) (result ManageDirectDebitBuilder) {
	result.Mutate(muts...)
	return
}

// ChangeTrustMutator is a interface that wraps the
// MutateChangeTrust operation.  types may implement this interface to
// specify how they modify an xdr.ChangeTrustOp object
type ManageDirectDebitMutator interface {
	MutateDirectDebit(op *xdr.ManageDirectDebitOp) error
}

// ChangeTrustBuilder represents a transaction that is being built.
type ManageDirectDebitBuilder struct {
	O   xdr.Operation
	CT  xdr.ManageDirectDebitOp
	Err error
}

// Mutate applies the provided mutators to this builder's payment or operation.
func (b *ManageDirectDebitBuilder) Mutate(muts ...interface{}) {
	for _, m := range muts {
		var err error
		switch mut := m.(type) {
		case ManageDirectDebitMutator:
			err = mut.MutateDirectDebit(&b.CT)
		case OperationMutator:
			err = mut.MutateOperation(&b.O)
		default:
			err = errors.New("Mutator type not allowed")
		}

		if err != nil {
			b.Err = err
			return
		}
	}
}

// MutateChangeTrust for Asset sets the ChangeTrustOp's Line field

func (debitor Destination) MutateDirectDebit(o *xdr.ManageDirectDebitOp) error {

	return setAccountId(debitor.AddressOrSeed, &o.Debitor)

}

func (m CancelDebit) MutateDirectDebit(o *xdr.ManageDirectDebitOp) (err error) {
	o.CancelDebit = bool(m)
	return
}
func (asset Asset) MutateDirectDebit(o *xdr.ManageDirectDebitOp) (err error) {
	o.Asset, err = asset.ToXDR()
	return err
}

// Trust is a helper that creates ChangeTrustBuilder
func CreateDebit(code, issuer, debitor string, args ...interface{}) (result ManageDirectDebitBuilder) {
	var deb = Destination{debitor}

	var asset Asset
	if issuer == "" && code == "" {
		asset = NativeAsset()
	} else {
		asset = CreditAsset(code, issuer)
	}

	mutators := []interface{}{
		asset,
		deb,
		CancelDebit(false),
	}
	for _, mut := range args {
		mutators = append(mutators, mut)
	}

	return ManageDirectDebit(mutators...)
}

// RemoveTrust is a helper that creates ChangeTrustBuilder
func RemoveDebit(code, issuer, debitor string, args ...interface{}) (result ManageDirectDebitBuilder, err error) {
	var deb = Destination{debitor}

	var asset Asset
	if issuer == "" && code == "" {
		asset = NativeAsset()
	} else {
		asset = CreditAsset(code, issuer)
	}
	mutators := []interface{}{
		asset,
		deb,
		CancelDebit(true),
	}

	for _, mut := range args {
		mutators = append(mutators, mut)
	}

	return ManageDirectDebit(mutators...), err
}

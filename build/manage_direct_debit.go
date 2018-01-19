package build

import (
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/xdr"
)


func ManageDirectDebit(muts ...interface{}) (result ManageDirectDebitBuilder) {
	result.Mutate(muts...)
	return
}


type ManageDirectDebitMutator interface {
	MutateDirectDebit(op *xdr.ManageDirectDebitOp) error
}


type ManageDirectDebitBuilder struct {
	O   xdr.Operation
	MDD  xdr.ManageDirectDebitOp
	Err error
}


func (b *ManageDirectDebitBuilder) Mutate(muts ...interface{}) {
	for _, m := range muts {
		var err error
		switch mut := m.(type) {
		case ManageDirectDebitMutator:
			err = mut.MutateDirectDebit(&b.MDD)
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


func CreateDebit(asset Asset, debitor string, args ...interface{}) (result ManageDirectDebitBuilder) {
	var deb = Destination{debitor}

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


func DeleteDebit(asset Asset, debitor string, args ...interface{}) (result ManageDirectDebitBuilder) {
	var deb = Destination{debitor}

	mutators := []interface{}{
		asset,
		deb,
		CancelDebit(true),
	}

	for _, mut := range args {
		mutators = append(mutators, mut)
	}

	return ManageDirectDebit(mutators...)
}

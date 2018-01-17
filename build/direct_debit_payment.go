package build

import (
"github.com/stellar/go/amount"
"github.com/stellar/go/support/errors"
"github.com/stellar/go/xdr"

)

// Payment groups the creation of a new PaymentBuilder with a call to Mutate.
func DebitPayment(muts ...interface{}) (result DebitPaymentBuilder) {
	result.Mutate(muts...)
	return
}

// PaymentMutator is a interface that wraps the
// MutatePayment operation.  types may implement this interface to
// specify how they modify an xdr.PaymentOp object
type DebitPaymentMutator interface {
	MutateDebitPayment(interface{}) error
}

// PaymentBuilder represents a transaction that is being built.
type DebitPaymentBuilder struct {

	O           xdr.Operation
	DP          xdr.DirectDebitPaymentOp
	P           xdr.PaymentOp
	Err         error
}

// Mutate applies the provided mutators to this builder's payment or operation.
func (b *DebitPaymentBuilder) Mutate(muts ...interface{}) {
	for _, m := range muts {
		var err error
		switch mut := m.(type) {
		case DebitPaymentMutator:
			err = mut.MutateDebitPayment(&b.DP)
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

// MutatePayment for Asset sets the PaymentOp's Asset field
func (m CreditAmount) MutateDebitPayment(o interface{}) (err error) {
	switch o := o.(type) {
	default:
		err = errors.New("Unexpected operation type")
	case *xdr.PaymentOp:
		o.Amount, err = amount.Parse(m.Amount)
		if err != nil {
			return
		}

		o.Asset, err = createAlphaNumAsset(m.Code, m.Issuer)
	case *xdr.PathPaymentOp:
		o.DestAmount, err = amount.Parse(m.Amount)
		if err != nil {
			return
		}

		o.DestAsset, err = createAlphaNumAsset(m.Code, m.Issuer)
	}
	return
}

// MutatePayment for Destination sets the PaymentOp's Destination field
func (m Destination) MutateDebitPayment(o interface{}) error {
	switch o := o.(type) {
	default:
		return errors.New("Unexpected operation type")
	case *xdr.DirectDebitPaymentOp:
		return setAccountId(m.AddressOrSeed, &o.Creditor)

	}
}

// MutatePayment for NativeAmount sets the PaymentOp's currency field to
// native and sets its amount to the provided intege

// MutatePayment for PayWithPath sets the PathPaymentOp's SendAsset,
// SendMax and Path fields
func (m PaymentOp) MutateDebitPayment(o interface{}) (err error) {
	var Pay *xdr.DirectDebitPaymentOp
	var ok bool
	if Pay, ok = o.(*xdr.DirectDebitPaymentOp); !ok {
		return errors.New("Unexpected operation type")
	}

	// MaxAmount
	Pay.Payment.Amount, err = amount.Parse(m.amount)
	if err != nil {
		return
	}
	Pay.Payment.Asset,err = m.asset.ToXDR()
	if err != nil {
		return
	}
	err=setAccountId(m.dest,&Pay.Payment.Destination)
	if err != nil {
		return
	}

	return
}


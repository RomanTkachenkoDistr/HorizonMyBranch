package build

import (
"github.com/stellar/go/amount"
"github.com/stellar/go/support/errors"
"github.com/stellar/go/xdr"

)


func DebitPayment(muts ...interface{}) (result DebitPaymentBuilder) {
	result.Mutate(muts...)
	return
}


type DebitPaymentMutator interface {
	MutateDebitPayment(interface{}) error
}


type DebitPaymentBuilder struct {

	O           xdr.Operation
	DP          xdr.DirectDebitPaymentOp
	Err         error
}


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

func (m CreditAmount) MutateDebitPayment(o interface{}) (err error) {

		o.(*xdr.DirectDebitPaymentOp).Payment.Amount, err = amount.Parse(m.Amount)
		if err != nil {
			return
		}
		o.(*xdr.DirectDebitPaymentOp).Payment.Asset, err = createAlphaNumAsset(m.Code, m.Issuer)

		return

}

func (m Destination) MutateDebitPayment(o interface{}) error {
		return setAccountId(m.AddressOrSeed, &o.(*xdr.DirectDebitPaymentOp).Payment.Destination)
}


func (m PaymentOp) MutateDebitPayment(o interface{}) (err error) {
	var Pay *xdr.DirectDebitPaymentOp
	var ok bool
	if Pay, ok = o.(*xdr.DirectDebitPaymentOp); !ok {
		return errors.New("Unexpected operation type")
	}

	// MaxAmount
	Pay.Payment.Amount, err = amount.Parse(m.Amount)
	if err != nil {
		return
	}
	Pay.Payment.Asset,err = m.Asset.ToXDR()
	if err != nil {
		return
	}
	err=setAccountId(m.Dest,&Pay.Payment.Destination)
	if err != nil {
		return
	}

	return
}
func (m Creditor) MutateDebitPayment(o interface{})(err error){
	return setAccountId(m.AddressOrSeed, &o.(*xdr.DirectDebitPaymentOp).Creditor)
}
func DirectDebitPayment(pay PaymentOp, creditor Creditor,args ...interface{}) (result DebitPaymentBuilder){
	mutators := []interface{}{
		pay,
		creditor,
	}

	for _, mut := range args {
		mutators = append(mutators, mut)
	}
	return DebitPayment(mutators...)
}

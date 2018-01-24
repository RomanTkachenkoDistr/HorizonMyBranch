package build

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/xdr"

)

var _ = Describe("DebitPayment", func() {

	Describe("DebitPaymentBuilder", func() {
		var (
			subject DebitPaymentBuilder
			mut     interface{}
			debitor = "GBHBUTDLCF7OCNIWPH3XZVRD5HCBPNQMYV455W2EV3WMNWLAHUJ2A4QD"

			creditor = "GAXEMCEXBERNSRXOEKD4JAIKVECIXQCENHEBRVSPX2TTYZPMNEDSQCNQ"
			bad      = "foo"
			asset    = CreditAsset("EUR", "GAWSI2JO2CF36Z43UGMUJCDQ2IMR5B3P5TMS7XM7NUTU3JHG3YJUDQXA")
			dest     = "GB6Y4WL4KLFFV7XFLCMCG4WV2G2KZBSWLYKXWJHRRYKDM6TP5FB7V2FI"
			amount   = "77"
			pay      = PaymentOp{
				Asset:  asset,
				Amount: amount,
				Dest:   dest,
			}
		)

		JustBeforeEach(func() {
			subject = DebitPaymentBuilder{}
			subject.Mutate(mut)
		})

		Describe("CreateDebitPayment", func() {
			Context("creates debit payment properly", func() {
				It("sets values properly", func() {
					builder := DirectDebitPayment(pay, Creditor{creditor})

					Expect(builder.DP.Payment.Amount).To(Equal(xdr.Int64(770000000)))

					Expect(builder.DP.Payment.Asset.Type).To(Equal(xdr.AssetTypeAssetTypeCreditAlphanum4))
					Expect(builder.DP.Payment.Asset.AlphaNum4.AssetCode).To(Equal([4]byte{'E', 'U', 'R', 0}))
					var aid xdr.AccountId
					aid.SetAddress(asset.Issuer)
					Expect(builder.DP.Payment.Asset.AlphaNum4.Issuer.MustEd25519()).To(Equal(aid.MustEd25519()))
					Expect(builder.DP.Payment.Asset.AlphaNum12).To(BeNil())
					aid.SetAddress(creditor)
					Expect(builder.DP.Creditor.MustEd25519()).To(Equal(aid.MustEd25519()))
					aid.SetAddress(dest)
					Expect(builder.DP.Payment.Destination.MustEd25519()).To(Equal(aid.MustEd25519()))
				})
			})
		})

		Describe("SourceAccount", func() {
			Context("using a valid stellar address", func() {
				BeforeEach(func() { mut = SourceAccount{debitor} })

				It("succeeds", func() {
					Expect(subject.Err).NotTo(HaveOccurred())
				})

				It("sets the destination to the correct xdr.AccountId", func() {
					var aid xdr.AccountId
					aid.SetAddress(debitor)
					Expect(subject.O.SourceAccount.MustEd25519()).To(Equal(aid.MustEd25519()))
				})
			})

			Context("using an invalid value", func() {
				BeforeEach(func() { mut = SourceAccount{bad} })
				It("failed", func() { Expect(subject.Err).To(HaveOccurred()) })
			})
		})
		Describe("mutate invalid values",func(){
			Context("mutate invalid creditamount", func() {

				Context("issuer invalid", func() {
					BeforeEach(func() {
						mut = CreditAmount{"USD", bad, "50.0"}
					})

					It("failed", func() {
						Expect(subject.Err).To(HaveOccurred())
					})
				})

				Context("amount invalid", func() {
					BeforeEach(func() {
						mut = CreditAmount{"ABCDEF", creditor, "test"}
					})

					It("failed", func() {
						Expect(subject.Err).To(HaveOccurred())
					})
				})

				Context("asset code length invalid", func() {
					Context("empty", func() {
						BeforeEach(func() {
							mut = CreditAmount{"", creditor, "50.0"}
						})

						It("failed", func() {
							Expect(subject.Err).To(MatchError("Asset code length is invalid"))
						})
					})

					Context("too long", func() {
						BeforeEach(func() {
							mut = CreditAmount{"1234567890123", creditor, "50.0"}
						})

						It("failed", func() {
							Expect(subject.Err).To(MatchError("Asset code length is invalid"))
						})
					})
				})
			})
			Context("invalid destination",func() {
				BeforeEach(func() {
					mut = Destination{bad}
				})
				It("set invalid destination", func() {

					Expect(subject.Err).To(HaveOccurred())
				})
			})
			Context("invalid creditor",func() {
				BeforeEach(func() {
					mut = Creditor{bad}
				})
				It("set invalid destination", func() {

					Expect(subject.Err).To(HaveOccurred())
				})
			})
		})
	})

})
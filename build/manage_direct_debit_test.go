package build

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/xdr"
)

var _ = Describe("ManageDebit", func() {

	Describe("ManageDebitBuilder", func() {
		var (
			subject ManageDirectDebitBuilder
			mut     interface{}

			address= "GAXEMCEXBERNSRXOEKD4JAIKVECIXQCENHEBRVSPX2TTYZPMNEDSQCNQ"
			bad= "foo"
			asset= CreditAsset("EUR", "GAWSI2JO2CF36Z43UGMUJCDQ2IMR5B3P5TMS7XM7NUTU3JHG3YJUDQXA")
		)

		JustBeforeEach(func() {
			subject = ManageDirectDebitBuilder{}
			subject.Mutate(mut)
		})

		Describe("CreateDebit", func() {
			Context("creates debit properly", func() {
				It("sets values properly", func() {
					builder := CreateDebit(asset, address)
					Expect(builder.MDD.Asset.Type).To(Equal(xdr.AssetTypeAssetTypeCreditAlphanum4))
					Expect(builder.MDD.Asset.AlphaNum4.AssetCode).To(Equal([4]byte{'E', 'U', 'R', 0}))
					var aid xdr.AccountId
					aid.SetAddress(asset.Issuer)

					Expect(builder.MDD.Asset.AlphaNum4.Issuer.MustEd25519()).To(Equal(aid.MustEd25519()))
					Expect(builder.MDD.Asset.AlphaNum12).To(BeNil())
					aid.SetAddress(address)
					Expect(builder.MDD.Debitor).To(Equal(aid))
					Expect(builder.MDD.CancelDebit).To(Equal(false))
				})
			})

		})



		Describe("DeleteDebit", func() {
			Context("deletes the debit properly", func() {
				It("sets values properly", func() {
					builder := DeleteDebit(asset, address)

					Expect(builder.MDD.Asset.Type).To(Equal(xdr.AssetTypeAssetTypeCreditAlphanum4))
					Expect(builder.MDD.Asset.AlphaNum4.AssetCode).To(Equal([4]byte{'E', 'U', 'R', 0}))
					var aid xdr.AccountId
					aid.SetAddress(asset.Issuer)
					Expect(builder.MDD.Asset.AlphaNum4.Issuer.MustEd25519()).To(Equal(aid.MustEd25519()))
					Expect(builder.MDD.Asset.AlphaNum12).To(BeNil())
					aid.SetAddress(address)
					Expect(builder.MDD.Debitor.MustEd25519()).To(Equal(aid.MustEd25519()))
					Expect(builder.MDD.CancelDebit).To(Equal(true))

				})
			})

		})

		Describe("SourceAccount", func() {
			Context("using a valid stellar address", func() {
				BeforeEach(func() { mut = SourceAccount{address} })

				It("succeeds", func() {
					Expect(subject.Err).NotTo(HaveOccurred())
				})

				It("sets the destination to the correct xdr.AccountId", func() {
					var aid xdr.AccountId
					aid.SetAddress(address)
					Expect(subject.O.SourceAccount.MustEd25519()).To(Equal(aid.MustEd25519()))
				})
			})

			Context("using an invalid value", func() {
				BeforeEach(func() { mut = SourceAccount{bad} })
				It("failed", func() { Expect(subject.Err).To(HaveOccurred()) })
			})
		})
		Describe("mutate invalid values",func(){
			Context("mutate invalid asset", func() {
				Context("asset code length invalid", func() {
					BeforeEach(func() {
						mut = Asset{"", address, false}
					})
					Context("empty", func() {
						It("failed", func() {
							Expect(subject.Err).To(HaveOccurred())
						})
					})

					BeforeEach(func() {
						mut = Asset{"1234567890123", address, false}
					})
					Context("asset code too long", func() {
						It("failed", func() {
							Expect(subject.Err).To(HaveOccurred())
						})
					})
				})
				BeforeEach(func() {
					mut = Asset{"USD", bad, false}
				})
				Context("invalid asset issuer",func(){
					It("failed", func() {
						Expect(subject.Err).To(HaveOccurred())
					})
				})
			})

			Context("invalid destination",func() {
				BeforeEach(func() {
					mut = Destination{bad}
				})
				It("failed", func() {

					Expect(subject.Err).To(HaveOccurred())
				})
			})
		})

	})

})
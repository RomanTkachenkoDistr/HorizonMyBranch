package core
//
//import (
//	"errors"
//
//	sq "github.com/Masterminds/squirrel"
//	"github.com/stellar/go/xdr"
//)
//
//// AssetsForAddress loads `dest` as `[]xdr.Asset` with every asset the account
//// at `addy` can hold.
//func (q *Q) DebitForAddress(dest interface{}, addy string) error {
//	var deb []DirectDebit
//
//	err := q.DebitByAddress(&deb, addy)
//	if err != nil {
//		return err
//	}
//
//	dtl, ok := dest.(*[]xdr.AccountId)
//	if !ok {
//		return errors.New("Invalid destination")
//	}
//
//	result := make([]xdr.AccountId, len(deb)+1)
//	*dtl = result
//
//	for i, db := range deb {
//		result[i], err = DebitFromDB(db.Creditor)
//		if err != nil {
//			return err
//		}
//	}
//
//	result[len(result)-1], err = xdr.NewAsset(xdr.AssetTypeAssetTypeNative, nil)
//
//	return err
//}
//
//
//func (q *Q) DebitByDebitor(dest interface{}, addy string) error {
//	sql := selectDebit.Where("creditor = ?", addy)
//	return q.Select(dest, sql)
//}
//
//
//
//
//var selectDebit = sq.Select(
//	"deb.creditor",
//	"deb.debitor",
//	"deb.assetcode",
//	"deb.assetissuer",
//	"deb.assettype",
//
//).From("debits deb")
//

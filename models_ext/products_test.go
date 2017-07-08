package models_ext_test

// import (
// 	. "github.com/harrisbaird/dailyteedeals/models_ext"
// )

// func TestMarkProductsInactive(t *testing.T) {
// 	type args struct {
// 		siteID int
// 		deal   bool
// 	}

// 	testCases := []struct {
// 		name    string
// 		args    args
// 		wantIDS []int
// 	}{
// 		{"deal", args{1, true}, []int{2, 3, 4}},
// 		{"non-deal", args{1, false}, []int{1, 3, 4}},
// 	}

// 	for _, tt := range testCases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			RunInTestTransaction(func(db boil.Executor) {
// 				CreateProductFixtures(db)
// 				err := MarkProductsInactive(db, tt.args.siteID, tt.args.deal)
// 				boil.DebugMode = true
// 				products, err := models.Products(db,
// 					qm.Select("id"),
// 					qm.Where("active = ?", true),
// 					qm.OrderBy("id ASC")).All()

// 				spew.Dump(products)

// 				var ids []int
// 				for _, product := range products {
// 					ids = append(ids, product.ID)
// 				}

// 				st.Expect(t, err, nil)
// 				st.Expect(t, ids, tt.wantIDS)
// 				st.Expect(t, err != nil, false)
// 			})
// 		})
// 	}
// }

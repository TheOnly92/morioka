package interfaces

import ()

type MenuItem struct {
	Identifier string
	Href       string
	Text       string
	Icon       string
	Active     bool
	Collapsed  bool
	Nested     []*MenuItem
}

func GetUserLeftMenu() []*MenuItem {
	return []*MenuItem{
		&MenuItem{"DefaultIndexHandler", "/", "マイ家計簿", "home", false, false, []*MenuItem{}},
		&MenuItem{"InputTransactions", "#", "家計簿を書く", "pencil", false, true, []*MenuItem{
			&MenuItem{"TransactionCreateBatchHandler", "/transaction/create/batch", "複数登録", "", false, false, []*MenuItem{}},
			&MenuItem{"TransactionCreateHandler", "/transaction/create", "詳細入力", "", false, false, []*MenuItem{}},
			&MenuItem{"TransactionReceiptCreateHandler", "/transaction/receipt/create", "レシート入力", "", false, false, []*MenuItem{}},
			&MenuItem{"TransactionTransferCreateHandler", "/transaction/transfer/create", "お金移動・残高調整", "", false, false, []*MenuItem{}},
		}},
		&MenuItem{"ViewTransactions", "#", "家計簿を見る", "book", false, true, []*MenuItem{
			&MenuItem{"TransactionListHandler", "/transaction/list", "家計簿一覧", "", false, false, []*MenuItem{}},
			&MenuItem{"FinanceDetailHandler", "/finance/detail", "残高・入出金明細", "", false, false, []*MenuItem{}},
			&MenuItem{"FinanceOverviewHandler", "/finance/overview", "集計", "", false, false, []*MenuItem{}},
			&MenuItem{"FinanceGraphHandler", "/finance/graph", "グラフ", "", false, false, []*MenuItem{}},
		}},
		&MenuItem{"FinanceSettings", "#", "家計簿の設定", "list", false, true, []*MenuItem{
			&MenuItem{"AccountManageHandler", "/accounts/#", "口座一覧", "", false, false, []*MenuItem{}},
			&MenuItem{"CategoryManageHandler", "/category/manage", "項目一覧", "", false, false, []*MenuItem{}},
			&MenuItem{"ShopManageHandler", "/shop/manage", "お店一覧", "", false, false, []*MenuItem{}},
			&MenuItem{"ProductManageHandler", "/product/manage", "商品一覧", "", false, false, []*MenuItem{}},
		}},
	}
}

func ConstructMenu(identifier string) []*MenuItem {
	menus := GetUserLeftMenu()
	for _, item := range menus {
		if item.Identifier == identifier {
			item.Active = true
			break
		}
		if item.Collapsed {
			for _, nest := range item.Nested {
				if nest.Identifier == identifier {
					item.Active = true
					nest.Active = true
				}
			}
		}
	}
	return menus
}

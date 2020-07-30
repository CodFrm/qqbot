package taobaoopen

type MaterialSearchRespond struct {
	Respond struct {
		ResultList struct {
			MapData []*MaterialItem `json:"map_data"`
		} `json:"result_list"`
	} `json:"tbk_dg_material_optional_response"`
}

type MaterialItem struct {
	ShortTitle     string `json:"short_title"`
	Url            string `json:"url"`
	CouponShareUrl string `json:"coupon_share_url"`
	CouponStartFee string `json:"coupon_start_fee"`
	CouponAmount   string `json:"coupon_amount"`
	ZkFinalPrice   string `json:"zk_final_price"`
	ReservePrice   string `json:"reserve_price"`
	PictUrl        string `json:"pict_url"`
}

type GetSpreadRespond struct {
	Respond struct {
		Results struct {
			TbkSpread []*SpreadItem `json:"tbk_spread"`
		} `json:"results"`
	} `json:"tbk_spread_get_response"`
}

type SpreadItem struct {
	Content string `json:"content"`
	ErrMsg  string `json:"err_msg"`
}

type TpwdRespond struct {
	Respond struct {
		Data struct {
			Model string `json:"model"`
		} `json:"data"`
	} `json:"tbk_tpwd_create_response"`
}

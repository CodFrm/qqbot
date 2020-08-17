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

type ConverseTkl struct {
	TbkPrivilegeGetResponse struct {
		Result struct {
			Data struct {
				CategoryID        string `json:"category_id"`
				CouponClickURL    string `json:"coupon_click_url"`
				CouponEndTime     string `json:"coupon_end_time"`
				CouponInfo        string `json:"coupon_info"`
				CouponAmount      string `json:"coupon_amount"`
				CouponStartfee    string `json:"coupon_startfee"`
				CouponRemainCount string `json:"coupon_remain_count"`
				CouponStartTime   string `json:"coupon_start_time"`
				CouponTotalCount  string `json:"coupon_total_count"`
				CouponType        string `json:"coupon_type"`
				ItemID            string `json:"item_id"`
				ItemURL           string `json:"item_url"`
				MaxCommissionRate string `json:"max_commission_rate"`
				MinCommissionRate string `json:"min_commission_rate"`
				SCouponID         string `json:"s_coupon_id"`
				SCouponAmount     string `json:"s_coupon_amount"`
				SCouponStartfee   string `json:"s_coupon_startfee"`
				SCouponStartTime  string `json:"s_coupon_start_time"`
				SCouponEndTime    string `json:"s_coupon_end_time"`
				CatLeafName       string `json:"cat_leaf_name"`
				CatName           string `json:"cat_name"`
				Nick              string `json:"nick"`
				PictURL           string `json:"pict_url"`
				Provcity          string `json:"provcity"`
				SellerID          string `json:"seller_id"`
				SmallImages       string `json:"small_images"`
				Title             string `json:"title"`
				UserType          string `json:"user_type"`
				Volume            string `json:"volume"`
				ZkFinalPrice      string `json:"zk_final_price"`
				Shorturl          string `json:"shorturl"`
				Tkl               string `json:"tkl"`
			} `json:"data"`
		} `json:"result"`
		RequestID string `json:"request_id"`
	} `json:"tbk_privilege_get_response"`
}

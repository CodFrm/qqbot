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
	Status  int `json:"status"`
	Content []struct {
		Code                   string `json:"code"`
		TypeOneID              string `json:"type_one_id"`
		TaoID                  string `json:"tao_id"`
		Title                  string `json:"title"`
		Jianjie                string `json:"jianjie"`
		PictURL                string `json:"pict_url"`
		UserType               string `json:"user_type"`
		SellerID               string `json:"seller_id"`
		ShopDsr                string `json:"shop_dsr"`
		Volume                 string `json:"volume"`
		Size                   string `json:"size"`
		QuanhouJiage           string `json:"quanhou_jiage"`
		DateTimeYongjin        string `json:"date_time_yongjin"`
		Tkrate3                string `json:"tkrate3"`
		YongjinType            string `json:"yongjin_type"`
		CouponID               string `json:"coupon_id"`
		CouponStartTime        string `json:"coupon_start_time"`
		CouponEndTime          string `json:"coupon_end_time"`
		CouponInfoMoney        string `json:"coupon_info_money"`
		CouponTotalCount       string `json:"coupon_total_count"`
		CouponRemainCount      string `json:"coupon_remain_count"`
		CouponInfo             string `json:"coupon_info"`
		Juhuasuan              string `json:"juhuasuan"`
		Taoqianggou            string `json:"taoqianggou"`
		Haitao                 string `json:"haitao"`
		Jiyoujia               string `json:"jiyoujia"`
		Jinpaimaijia           string `json:"jinpaimaijia"`
		Pinpai                 string `json:"pinpai"`
		PinpaiName             string `json:"pinpai_name"`
		Yunfeixian             string `json:"yunfeixian"`
		Nick                   string `json:"nick"`
		SmallImages            string `json:"small_images"`
		WhiteImage             string `json:"white_image"`
		TaoTitle               string `json:"tao_title"`
		Provcity               string `json:"provcity"`
		ShopTitle              string `json:"shop_title"`
		ZhiboURL               string `json:"zhibo_url"`
		SellCount              string `json:"sellCount"`
		CommentCount           string `json:"commentCount"`
		Favcount               string `json:"favcount"`
		Score1                 string `json:"score1"`
		Score2                 string `json:"score2"`
		Score3                 string `json:"score3"`
		CreditLevel            string `json:"creditLevel"`
		ShopIcon               string `json:"shopIcon"`
		PcDescContent          string `json:"pcDescContent"`
		TaobaoURL              string `json:"taobao_url"`
		CategoryID             string `json:"category_id"`
		CategoryName           string `json:"category_name"`
		LevelOneCategoryID     string `json:"level_one_category_id"`
		LevelOneCategoryName   string `json:"level_one_category_name"`
		Tkfee3                 string `json:"tkfee3"`
		Biaoqian               string `json:"biaoqian"`
		Tag                    string `json:"tag"`
		DateTime               string `json:"date_time"`
		PresaleDiscountFeeText string `json:"presale_discount_fee_text"`
		PresaleTailEndTime     string `json:"presale_tail_end_time"`
		PresaleTailStartTime   string `json:"presale_tail_start_time"`
		PresaleEndTime         string `json:"presale_end_time"`
		PresaleStartTime       string `json:"presale_start_time"`
		PresaleDeposit         string `json:"presale_deposit"`
		MinCommissionRate      string `json:"min_commission_rate"`
		CouponClickURL         string `json:"coupon_click_url"`
		ItemURL                string `json:"item_url"`
		Shorturl               string `json:"shorturl"`
		Tkl                    string `json:"tkl"`
	} `json:"content"`
}

type ZtkError struct {
	Status  int    `json:"status"`
	Content string `json:"content"`
}

func (e *ZtkError) Error() string {
	return e.Content
}

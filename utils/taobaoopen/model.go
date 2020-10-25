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

type GetActiveInfo struct {
	Response struct {
		Data struct {
			WxQrcodeURL       string `json:"wx_qrcode_url"`
			ClickURL          string `json:"click_url"`
			ShortClickURL     string `json:"short_click_url"`
			TerminalType      string `json:"terminal_type"`
			MaterialOssURL    string `json:"material_oss_url"`
			PageName          string `json:"page_name"`
			PageStartTime     string `json:"page_start_time"`
			PageEndTime       string `json:"page_end_time"`
			WxMiniprogramPath string `json:"wx_miniprogram_path"`
		} `json:"data"`
	} `json:"tbk_activity_info_get_response"`
}

type SpreadItem struct {
	Content string `json:"content"`
	ErrMsg  string `json:"err_msg"`
}

type TpwdRespond struct {
	Respond struct {
		Data struct {
			Model          string `json:"model"`
			PasswordSimple string `json:"password_simple"`
		} `json:"data"`
	} `json:"tbk_tpwd_create_response"`
}

type ConverseTklContent struct {
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
}

type TklAddress struct {
	Status  int    `json:"status"`
	Content string `json:"content"`
	URLType string `json:"url_type"`
	URLID   string `json:"url_id"`
}

type ConverseTkl struct {
	Status  int                  `json:"status"`
	Content []ConverseTklContent `json:"content"`
}

type ZtkError struct {
	Status  int    `json:"status"`
	Content string `json:"content"`
}

func (e *ZtkError) Error() string {
	return e.Content
}

type OrderOption func(options *OrderOptions)

type OrderOptions struct {
	StartTime     string
	EndTime       string
	OrderScene    string
	PageNo        int
	PositionIndex int
}

type QueryOrderUrl struct {
	Url string `json:"url"`
}

type TaobaoError struct {
	ErrorResponse struct {
		Code      int    `json:"code"`
		Msg       string `json:"msg"`
		SubCode   string `json:"sub_code"`
		SubMsg    string `json:"sub_msg"`
		RequestID string `json:"request_id"`
	} `json:"error_response"`
}

func (e *TaobaoError) Error() string {
	return e.ErrorResponse.SubMsg
}

type OrderItem struct {
	AdzoneID                           int64  `json:"adzone_id"`
	AdzoneName                         string `json:"adzone_name"`
	AlimamaRate                        string `json:"alimama_rate"`
	AlimamaShareFee                    string `json:"alimama_share_fee"`
	AlipayTotalPrice                   string `json:"alipay_total_price"`
	ClickTime                          string `json:"click_time"`
	DepositPrice                       string `json:"deposit_price"`
	FlowSource                         string `json:"flow_source"`
	IncomeRate                         string `json:"income_rate"`
	IsLx                               string `json:"is_lx"`
	ItemCategoryName                   string `json:"item_category_name"`
	ItemID                             int64  `json:"item_id"`
	ItemImg                            string `json:"item_img"`
	ItemLink                           string `json:"item_link"`
	ItemNum                            int    `json:"item_num"`
	ItemPrice                          string `json:"item_price"`
	ItemTitle                          string `json:"item_title"`
	OrderType                          string `json:"order_type"`
	PubID                              int    `json:"pub_id"`
	PubShareFee                        string `json:"pub_share_fee"`
	PubSharePreFee                     string `json:"pub_share_pre_fee"`
	PubShareRate                       string `json:"pub_share_rate"`
	RefundTag                          int    `json:"refund_tag"`
	SellerNick                         string `json:"seller_nick"`
	SellerShopTitle                    string `json:"seller_shop_title"`
	SiteID                             int    `json:"site_id"`
	SiteName                           string `json:"site_name"`
	SubsidyFee                         string `json:"subsidy_fee"`
	SubsidyRate                        string `json:"subsidy_rate"`
	SubsidyType                        string `json:"subsidy_type"`
	TbDepositTime                      string `json:"tb_deposit_time"`
	TbPaidTime                         string `json:"tb_paid_time"`
	TerminalType                       string `json:"terminal_type"`
	TkCommissionFeeForMediaPlatform    string `json:"tk_commission_fee_for_media_platform"`
	TkCommissionPreFeeForMediaPlatform string `json:"tk_commission_pre_fee_for_media_platform"`
	TkCommissionRateForMediaPlatform   string `json:"tk_commission_rate_for_media_platform"`
	TkCreateTime                       string `json:"tk_create_time"`
	TkDepositTime                      string `json:"tk_deposit_time"`
	TkOrderRole                        int    `json:"tk_order_role"`
	TkPaidTime                         string `json:"tk_paid_time"`
	TkStatus                           int    `json:"tk_status"`
	TkTotalRate                        string `json:"tk_total_rate"`
	TotalCommissionFee                 string `json:"total_commission_fee"`
	TotalCommissionRate                string `json:"total_commission_rate"`
	TradeID                            string `json:"trade_id"`
	TradeParentID                      string `json:"trade_parent_id"`
	PayPrice                           string `json:"pay_price,omitempty"`
	TkEarningTime                      string `json:"tk_earning_time,omitempty"`
}

type OrderQueryRespond struct {
	TbkScOrderDetailsGetResponse struct {
		Data struct {
			HasNext       bool   `json:"has_next"`
			HasPre        bool   `json:"has_pre"`
			PageNo        int    `json:"page_no"`
			PageSize      int    `json:"page_size"`
			PositionIndex string `json:"position_index"`
			Results       struct {
				PublisherOrderDto []*OrderItem `json:"publisher_order_dto"`
			} `json:"results"`
		} `json:"data"`
		RequestID string `json:"request_id"`
	} `json:"tbk_sc_order_details_get_response"`
}

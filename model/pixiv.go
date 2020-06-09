package model

type PixivPicItem struct {
	Alt             string `json:"alt"`
	Id              string `json:"id"`
	ProfileImageUrl string `json:"profileImageUrl"`
	Url             string `json:"url"`
	UserId          string `json:"userId"`
	IllustTitle     string `json:"illustTitle"`
	UserName        string `json:"userName"`
	Title           string `json:"title"`
}

type IllustRespond struct {
	Body struct {
		Illust struct {
			Data  []*PixivPicItem `json:"data"`
			Total int             `json:"total"`
		} `json:"illust"`
		RelatedTags []string `json:"relatedTags"`
	} `json:"body"`
}

type PixivTags struct {
	Candidates []struct {
		TagName string `json:"tag_name"`
	} `json:"candidates"`
}

type PixivIllust struct {
	Body struct {
		Urls struct {
			Mini     string `json:"mini"`
			Regular  string `json:"regular"`
			Original string `json:"original"`
			Small    string `json:"small"`
		} `json:"urls"`
		UserIllusts map[string]interface{} `json:"userIllusts"`
		IllustTitle string                 `json:"illustTitle"`
		Title       string                 `json:"title"`
		UserName    string                 `json:"userName"`
		Tags        struct {
			Tags []struct {
				Tag string `json:"tag"`
			} `json:"tags"`
		} `json:"tags"`
	} `json:"body"`
}

type PixivRecommend struct {
	Body struct {
		Illusts []struct {
			Id string `json:"id"`
		} `json:"illusts"`
		NextIds []string `json:"nextIds"`
	} `json:"body"`
}

type PixivRankList struct {
	Contents []*struct {
		IllustId   int    `json:"illust_id"`
		ProfileImg string `json:"profile_img"`
		Url        string `json:"url"`
		UserId     int    `json:"user_id"`
		UserName   string `json:"user_name"`
		Title      string `json:"title"`
	} `json:"contents"`
}

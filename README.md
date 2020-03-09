## Table to Struct 工具 

将mysql表结构自动映射成标准xorm po层文件，节省时间  
支持自定义mysql连接串和直接从apollo读取

### Usage
```
 go install -a  github.com/echoingtech/t2s
```

```
Usage of t2s:
  -H string
    	mysql host (dsn模式必传)
  -P string
    	mysql port (dsn模式必传)
  -c string
    	配置类型 apollo (从apollo读取mysql配置)  dsn (自定义mysql配置) (default "apollo")
  -db string
    	数据库名
  -out string
    	生成文件地址,不指定则不生成
  -p string
    	mysql password (dsn必传模式)
  -package string
    	生成文件的packageName (default "po")
  -t string
    	表名, 多表按照 ',' 隔开
  -u string
    	mysql username (dsn必传模式)
```

### Example
```
自动生成payment库refund表的映射文件

apollo 模式需要内网环境

从apollo读取Mysql配置(单表生成)

t2s -db payment -t refund -out ./refund.go

从apollo读取Mysql配置(多表生成)
t2s -db payment -t refund,refund_detail,payment_config, -out ./refund.go

自定义Mysql配置

t2s -c dsn -H 127.0.0.1 -P 3306 -db payment -t refund -u username -p password -out ./refund.go

```
### Result
夸张的使用了以前开发遇到的一张奇葩表，这张表促使了这个工具的诞生
```
package mysql

type CbdGoods struct {
	GoodsId                int64   `xorm:"int(11) pk  autoincr  notnull 'goods_id'"` // 自增ID
	GoodsBzno              string  `xorm:"varchar(50) unique  null 'goods_bzno'"`    // 宝尊商品编码
	Cate                   int64   `xorm:"int(11) null 'cate'"`                      // 所属分类
	SecondCate             int64   `xorm:"int(11) null 'second_cate'"`               // 二级分类
	ThreeCate              int64   `xorm:"int(11) null 'three_cate'"`                // 三级分类
	BrandId                int64   `xorm:"int(11) notnull 'brand_id'"`               // 品牌ID，对应到cbd_brands表的brand_id字段
	Brand                  string  `xorm:"varchar(100) null 'brand'"`                // 所属品牌
	GoodsType              int64   `xorm:"int(11) null 'goods_type'"`                // 商品类型（1：福袋）
	GoodsBrand             string  `xorm:"varchar(100) null 'goods_brand'"`          // 卖客疯品牌
	GoodsNo                string  `xorm:"varchar(30) null 'goods_no'"`              // 商品编号
	GoodsName              string  `xorm:"varchar(200) null 'goods_name'"`           // 商品名称
	Content                string  `xorm:"text null 'content'"`                      // 文字介绍
	ContentHtml            string  `xorm:"text null 'content_html'"`                 // html文字介绍
	Videourl               string  `xorm:"varchar(200) null 'videourl'"`             // 视频url
	Coverurl               string  `xorm:"varchar(200) null 'coverurl'"`             // 视频封面图url
	Base                   int64   `xorm:"int(11) null 'base'"`                      // 初始销量
	Saletype               int     `xorm:"tinyint(1) null 'saletype'"`               // 销售类型 1：普通商品；2：限时抢购；3：限量抢购；4：内卖商品
	OpenTime               int64   `xorm:"int(11) null 'open_time'"`                 // 抢购开始时间
	CloseTime              int64   `xorm:"int(11) null 'close_time'"`                // 结束时间
	Activitycontent        string  `xorm:"text null 'activitycontent'"`              // 限时抢购首页文案：图片显示的五行字，数组序列化保存
	Realnum                int64   `xorm:"int(11) null 'realnum'"`                   // 限量抢购真实数量
	Viewnum                int64   `xorm:"int(11) null 'viewnum'"`                   // 限量抢购显示数量
	Surplus                int64   `xorm:"int(11) null 'surplus'"`                   // 限量抢购剩余件数
	Buyupnum               int64   `xorm:"int(2) null 'buyupnum'"`                   // 单笔限购数量
	Label                  int     `xorm:"tinyint(2) null 'label'"`                  // 商品标签：0:new;1~9折扣
	StartTime              int64   `xorm:"int(11) null 'start_time'"`                // 折扣或新品有效开始时间
	EndTime                int64   `xorm:"int(11) null 'end_time'"`
	IsTop                  int     `xorm:"tinyint(1) null 'is_top'"`                              // 是否置顶
	Topnum                 int64   `xorm:"int(11) null 'topnum'"`                                 // 置顶位置（数字）
	TopsTime               int64   `xorm:"int(11) null 'tops_time'"`                              // 置顶开始时间
	TopeTime               int64   `xorm:"int(11) null 'tope_time'"`                              // 置顶结束时间
	Remark                 string  `xorm:"text null 'remark'"`                                    // 操作记录
	Createtime             int64   `xorm:"int(10) null 'createtime'"`                             // 创建时间
	FreightMoney           float64 `xorm:"decimal(10,2) null 'freight_money'"`                    // 是否收取运费（0：不收取，大于0的数字标示每多少件收取一次运费
	Freight                int64   `xorm:"int(11) null 'freight'"`                                // 收取运费的标准，比如freight是4，表示每4件收取一次运费。
	Operator               int64   `xorm:"int(11) null 'operator'"`                               // 创建人（编辑此商品的人）
	Brankcode              string  `xorm:"varchar(50) null 'brankcode'"`                          // 品牌编码
	GoodsParam             string  `xorm:"text null 'goods_param'"`                               // 商品参数（json字符串）;
	Skuorder               string  `xorm:"text null 'skuorder'"`                                  // 所有sku项排序
	Praisenum              int64   `xorm:"int(11) null 'praisenum'"`                              // 被赞数
	BreakageText           string  `xorm:"varchar(255) null 'breakage_text'"`                     // 商品提示说明
	BreakageImg            string  `xorm:"varchar(255) null 'breakage_img'"`                      // 商品原图
	SpecialId              int64   `xorm:"int(11) null 'special_id'"`                             // 所属专题
	BreakageTitle          string  `xorm:"varchar(255) null 'breakage_title'"`                    // 商品实物图名称
	Buydownnum             int64   `xorm:"int(2) null 'buydownnum'"`                              // 限购下限
	SupplierId             int64   `xorm:"int(11) null 'supplier_id'"`                            // 所属供应商
	Bathch                 string  `xorm:"varchar(200) null 'bathch'"`                            // 批次
	Salenum                int64   `xorm:"int(11) null 'salenum'"`                                // 销售数量
	Buyer                  int64   `xorm:"int(11) null 'buyer'"`                                  // 买手
	SpecKey                string  `xorm:"varchar(60) null 'spec_key'"`                           // 对应获取图片的sku编号
	BatchKey               string  `xorm:"varchar(100) null 'batch_key'"`                         // 获取图片的批次号
	Gender                 int     `xorm:"tinyint(1) null 'gender'"`                              // 性别：0，不限；1，(成年)男；2，(成年)女；3，(成年)中性；4，童(已废弃)；5，男童；6，女童；7，(童)中性；
	BakBuydownnum          int64   `xorm:"int(11) null 'bak_buydownnum'"`                         // 限购下限备份
	ShopId                 int64   `xorm:"int(11) null 'shop_id'"`                                // oms店铺ID
	Consignment            int     `xorm:"tinyint(1) null 'consignment'"`                         // 代销模式：1，入库代销；2，FBS代销；3，付款经销
	MinExpDate             string  `xorm:"varchar(10) null 'min_exp_date'"`                       // 最小有效期
	MaxExpDate             string  `xorm:"varchar(10) null 'max_exp_date'"`                       // 最大有效期
	IsPack                 int     `xorm:"tinyint(1) null 'is_pack'"`                             // 是否为打包销售的商品[针对商品ID]（0：否，1：是）
	IsPackGoods            int     `xorm:"tinyint(1) null 'is_pack_goods'"`                       // 是否打包商品[针对打包ID]（0：否，1：是）
	BrandPic               string  `xorm:"varchar(100) null 'brand_pic'"`                         // 品牌图片(2,3,4)
	BrandContentPic        string  `xorm:"varchar(100) null 'brand_content_pic'"`                 // 品牌系列图片(2,3,4)
	BrandSpec              string  `xorm:"varchar(100) null 'brand_spec'"`                        // 品牌规格图片(2,3,4)
	ReturnRemark           string  `xorm:"text null 'return_remark'"`                             // 退换货说明
	PackGoods              string  `xorm:"text null 'pack_goods'"`                                // 礼包内商品集合
	CateTag                int64   `xorm:"int(11) null 'cate_tag'"`                               // 分类标签
	SelectVal              int64   `xorm:"int(1) null 'select_val'"`                              // 单笔限购下线：1--系统定义；2--人工定义
	GoodsGif               string  `xorm:"varchar(200) null 'goods_gif'"`                         // 商品GIF图片存放地址
	HDphoto                string  `xorm:"varchar(200) null 'hdphoto'"`                           // 商品高清图
	IsFbs                  int     `xorm:"tinyint(1) null 'is_fbs'"`                              // 是否FBS商品（0：否，1：是，按供应商拆单）
	FbsVersion             string  `xorm:"varchar(10) null 'fbs_version'"`                        // fbs版本号(目前只区分1.0和2.0)
	RandomPack             int     `xorm:"tinyint(1) null 'random_pack'"`                         // 是否随机闷包（针对闷包ID）
	MisBuyPrice            float64 `xorm:"decimal(10,5) null 'mis_buy_price'"`                    // mis报表上的进货价
	IsHaitao               int     `xorm:"tinyint(3) unsigned notnull 'is_haitao'"`               // 是否海淘商品：0，不是；1，是。
	DownPrice              float64 `xorm:"decimal(10,0) unsigned notnull 'down_price'"`           // 全网最低价
	GoodsAdjustPriceType   int     `xorm:"tinyint(1) unsigned notnull 'goods_adjust_price_type'"` // 商品调价类型1:普通调价
	State                  int     `xorm:"tinyint(1) null 'state'"`                               // 商品状态，1：首页上架，2：列表页上架，3：隐身上架，4：下架，5：仅在专区显示(弃用但保留) 6：店铺展示 7:品牌特卖
	Status                 int     `xorm:"tinyint(4) null 'status'"`                              // 操作：0:已创建；1：已保存；2：待审核；3：已通过；4：未通过;5：待运营
	Updatetime             int64   `xorm:"int(10) unsigned notnull 'updatetime'"`                 // 更新时间
	Saletime               int64   `xorm:"int(11) notnull 'saletime'"`                            // 上架时间(暂时不使用,使用start_time)
	Checktime              int64   `xorm:"int(10) null 'checktime'"`                              // 审核时间
	OutReason              string  `xorm:"varchar(200) null 'out_reason'"`                        // 下架原因
	TagId                  int64   `xorm:"int(11) null 'tag_id'"`                                 // 商品标签关联tags.tag_id
	IsFullYearVendible     int     `xorm:"tinyint(4) notnull 'is_full_year_vendible'"`            // 商品是否全年可销售：0，否；1，是。为0表示商品只在闪降期间可销售。
	IsSeckill              int     `xorm:"tinyint(4) notnull 'is_seckill'"`                       // 是否秒杀：0，否；1，是
	IsFlashDown            int     `xorm:"tinyint(1) null 'is_flash_down'"`                       // 是否分享闪降商品（0:否；1:是）
	FlashDownNum           int64   `xorm:"int(11) null 'flash_down_num'"`                         // 闪降分享人数
	IsCaveatEmptor         int     `xorm:"tinyint(1) null 'is_caveat_emptor'"`                    // 是否不可退换（0:否；1:是）
	IsCopied               int     `xorm:"tinyint(4) notnull 'is_copied'"`                        // 是否复制的商品：0，否；1，是
	IsDelete               int     `xorm:"tinyint(1) null 'is_delete'"`                           // 商品是否删除（0:未删除；1:已删除）
	SdjjCut                int     `xorm:"tinyint(1) null 'sdjj_cut'"`                            // 是否参与闪电砍价（0:不参与，1:参与）
	ForSpring              int     `xorm:"tinyint(1) null 'for_spring'"`                          // 适合春季(0:否；1：是)
	ForSummer              int     `xorm:"tinyint(1) null 'for_summer'"`                          // 适合夏季(0:否；1：是)
	ForAutumn              int     `xorm:"tinyint(1) null 'for_autumn'"`                          // 适合秋季(0:否；1：是)
	ForWinter              int     `xorm:"tinyint(1) null 'for_winter'"`                          // 适合冬季(0:否；1：是)
	ForCrowd               string  `xorm:"varchar(50) null 'for_crowd'"`                          // 适合的人群
	CatSortScore           int64   `xorm:"int(11) notnull 'cat_sort_score'"`                      // 在分类商品列表里的排序分数
	TaojjState             int     `xorm:"tinyint(1) null 'taojj_state'"`                         // 是否淘集集上架（0:未上架，1:已上架）
	SourceGoodsId          int64   `xorm:"int(11) null 'source_goods_id'"`                        // 来源商品ID（此列不为空时，标明该商品是复制商品。）
	FreightTpl             int64   `xorm:"int(11) null 'freight_tpl'"`                            // 运费模板
	RepeatCustomersSort    int64   `xorm:"int(10) null 'repeat_customers_sort'"`                  // 老用户排序
	NewCustomersSort       int64   `xorm:"int(10) null 'new_customers_sort'"`                     // 新用户排序
	GroupNumber            int     `xorm:"smallint(6) notnull 'group_number'"`                    // 拼团人数
	HotSource              string  `xorm:"varchar(200) notnull 'hot_source'"`                     // 推荐为图片
	SellingPoint           string  `xorm:"varchar(200) notnull 'selling_point'"`                  // 卖点
	PerMonthBuyUp          int64   `xorm:"int(10) null 'per_month_buy_up'"`                       // 月限购数量（0：不限制）
	IsHotCake              int     `xorm:"tinyint(2) null 'is_hot_cake'"`                         // 爆款商品 0：不是爆款  1：爆款
	ShortTitle             string  `xorm:"varchar(100) null 'short_title'"`                       // 短标题
	LimitationTimes        int64   `xorm:"int(1) null 'limitation_times'"`                        // 限购次数
	DeliveryTimeCommitment int     `xorm:"tinyint(4) null 'delivery_time_commitment'"`            // 发货时间承诺,0:无,1:24小时,2:48小时
	IsReturnChange         int     `xorm:"tinyint(2) null 'is_return_change'"`                    // 是否支持7天无理由退换货,0:否,1:是
	IsLostTen              int     `xorm:"tinyint(2) null 'is_lost_ten'"`                         // 是否支持假一赔十,0:否,1:是(淘集集)
	MaxReward              float64 `xorm:"decimal(10,2) null 'max_reward'"`                       // 最高平分奖励
	SpecialSaleStatus      int     `xorm:"tinyint(2) null 'special_sale_status'"`                 // 审核状态,0:新增(无),1:待审核,2:审核通过,3:审核不通过
	FirstPutawayTime       int64   `xorm:"int(1) null 'first_putaway_time'"`                      // 首次上架时间
	MerchantSubsidyPrice   float64 `xorm:"decimal(10,2) null 'merchant_subsidy_price'"`           // 商家补贴价
	IsShowicon             int     `xorm:"tinyint(3) unsigned null 'is_showicon'"`                // 是否显示icon 0否 1是
	IconShowStartTime      int64   `xorm:"int(11) unsigned null 'icon_show_start_time'"`          // 图标显示开始时间
	IconShowEndTime        int64   `xorm:"int(11) unsigned null 'icon_show_end_time'"`            // 图标显示结束时间
	IconGoodDetial         string  `xorm:"varchar(255) null 'icon_good_detial'"`                  // 商品页活动图标
	IconList               string  `xorm:"varchar(255) null 'icon_list'"`                         // 商品列表活动图标
	IconCart               string  `xorm:"varchar(255) null 'icon_cart'"`                         // 购物车商品活动图标
	SyncTime               int64   `xorm:"int(11) notnull 'sync_time'"`                           // 同步时间
}

func (CbdGoods) TableName() string {
	return "cbd_goods"
}


```

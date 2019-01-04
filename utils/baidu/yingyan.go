package baidu

//yingyan sdk
type YingyanParam struct {
	Ak         string `json:"ak"`
	ServiceID  string `json:"service_id"`
	EntityName string `json:"entity_name"`
	EntityDesc string `json:"entity_desc"`
	CityName   string `json:"city_name"`
	District   string `json:"district"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	PointsList string `json:"point_list"`
}

//添加鹰眼账号
func (y *YingyanParam) AddEntity(EntityName, EntityDesc, CityName, District string) *YingyanParam {
	y.Ak = Ak               //测试ak 个人未验证（单日调用量100000）
	y.ServiceID = ServiceID // 创建鹰眼服务ID （线上162000）测试：162002
	y.EntityName = EntityName
	y.EntityDesc = EntityDesc
	y.CityName = CityName
	y.District = District

	return y
}

//鹰眼批量上传坐标点
func (y *YingyanParam) AddPoints(PointList string) *YingyanParam {
	y.Ak = Ak
	y.ServiceID = ServiceID
	y.PointsList = PointList

	return y
}

//查询鹰眼纠偏里程
func (y *YingyanParam) GetDistance(StartTime, EndTime int64, EntityName string) *YingyanParam {
	y.Ak = Ak
	y.ServiceID = ServiceID
	y.StartTime = StartTime
	y.EndTime = EndTime
	y.EntityName = EntityName

	return y
}

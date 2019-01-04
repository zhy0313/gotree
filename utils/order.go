package utils

// DivisionUseType 拆分预约用车订单类型
func DivisionUseType(serviceType int8) (useType, parttenType int8) {
	switch serviceType {
	case 11, 12, 21, 22:
		parttenType = serviceType % 10
		useType = serviceType / 10
	default:
		parttenType = 0
		useType = 0
	}
	return
}

// 1：即时用车 11：接机 12：送机 21：接站 22：送站 4：日租 3：半日租 15：预约用车 7：出租车 8：导游 9：车+导
// SwitchOrderTypeToCommentType  订单类型转换为评价订单类型
func SwitchOrderTypeToCommentType(serviceType, days int8) (commentType int8) {
	switch serviceType {
	case 2:
		commentType = 8
	case 3:
		commentType = 9
	case 4:
		commentType = 3
		if days > 0 {
			commentType = 4
		}
	case 5:
		commentType = 1
	default:
		commentType = 0

	}
	return
}

// 转换车+导、导游订单状态 为车订单状态
/*
车+导
when d.OrderStatus=1 then 1 --待支付
			when d.OrderStatus=1 and d.IsPaid=1 and IsBooked in(0,2) then 2  --待确认
			when d.OrderStatus=3 then 3  --待服务
			when d.OrderStatus=8 then 3  --待服务
			when d.OrderStatus=9 then 4  -- 服务中
			when d.OrderStatus=5 then 5  -- 已完成
			when d.OrderStatus in(13,14) then 6  --已关闭
导游
case	when OrderStatus=1 then 1 --待支付
				when OrderStatus=3 and PayType=2 and IsBooked in(0,2) then 2  --待确认
			 	when OrderStatus=8 then 3  --待服务
			 	when OrderStatus=9 then 4  -- 服务中
				when OrderStatus=5 then 5  -- 已完成
				when OrderStatus in(13,14) then 6  --已关闭
*/
func SwitchGuideOrderStatusToCarStatus(orderStatus int8) (status int8) {
	switch orderStatus {
	case 2:
		status = 1
	case 3:
		status = 4
	case 4:
		status = 7
	case 5:
		status = 8
	case 6:
		status = 5
	default:
		status = 0

	}
	return
}

//CarOrderBookingServiceTypeName 获取类型名称
func CarOrderBookingServiceTypeName(typeID int8) (typeName string) {
	switch typeID {
	case 11:
		typeName = "接机"
	case 12:
		typeName = "送机"
	case 21:
		typeName = "接站"
	case 22:
		typeName = "送站"
	case 15:
		typeName = "预约用车"
	default:
	}
	return
}

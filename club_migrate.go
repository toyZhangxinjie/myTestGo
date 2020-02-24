package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ClubPostMigrate() {
	str := `电缆屏蔽带圈	334	铜带	113
电缆带圈	335	铜带	113
电缆箔圈	336	铜箔	112
电磁线圈	18	漆包线	88
铜铸件圈	115	铜合金锭	114
热交换器圈	213	电机	392
换热器圈	283	电机	392
手电筒圈	204	五金	393
密码锁圈	227	五金	393
铆钉圈	229	五金	393
螺丝圈	234	五金	393
螺母圈	235	五金	393
链条锁圈	253	五金	393
电钻圈	327	五金	393
电动螺丝刀圈	339	五金	393
电动扳手圈	340	五金	393
安全锁圈	362	五金	393
轴套圈	164	五金	393
铜套圈	118	五金	393
轴承圈	3	五金	393
水龙头圈	203	水暖	394
离心泵圈	257	水暖	394
花洒圈	286	水暖	394
电熨斗圈	332	家电	274
灯具圈	346	家电	274
电刨圈	333	低压电器	395
电锤圈	342	低压电器	395
电磁搅拌器圈	341	低压电器	395
通讯电子圈	189	电子元件	330
预焙阳极块圈	170	预焙阳极	388
电解铝圈	19	A00铝锭	373
铝合金圈	369	铝合金锭	63
低合金板卷圈	14	铝板	56
铝蜂窝板圈	248	铝板	56
铝单板圈	249	铝板	56
铝条圈	244	铝带	59
铝型材圈	68	铝型材	396
型材铝圈	175	铝型材	396
铝压铸圈	64	压铸件	398
涡轮发动机圈	185	汽配	399
涡轮起重器圈	183	汽配	399
涡轮机圈	184	汽配	399
散热片圈	208	汽配	399
离合器片圈	258	汽配	399
减震器圈	273	汽配	399
机件生铝圈	277	汽配	399
火花塞圈	280	汽配	399
缸盖圈	300	汽配	399
调压器圈	326	汽配	399
变速器圈	359	汽配	399
喷油嘴圈	223	摩配	400
铝圆片圈	242	幕墙	401
铝标牌圈	250	铝塑板	245
再生精铅圈	157	还原铅	40
电解铅圈	21	铅锭	402
铅蓄电池圈	144	电动自行车电池	403
锌压铸圈	139	压铸锌合金	145
镀锌板卷圈	24	锌板	404
镀锌结构件圈	25	镀锌	405
锡酸锌圈	138	化工	406
环保焊锡圈	284	锡丝	131
焊锡圈	390	锡丝	131
锡合金圈	178	锡合金锭	407
硫酸亚锡圈	86	焦磷酸亚锡	408
锡酸钠圈	132	锡酸钾	409
氯化亚锡圈	81	氧化亚锡	410
二氧化锡圈	320	氧化亚锡	410
二氧化硅圈	322	有机硅	411
金属铟圈	383	铟锭	154
金属镁圈	380	镓镁合金	412
金属镓圈	34	氧化镓	413
金属铋圈	376	铋锭	1
金属硒圈	377	硒锭	124
二硫化 硒圈	26	二氧化硒	414
金属碲圈	381	碲锭	345
二氧化碲圈	324	碲粉	415
金属锑圈	379	锑锭	107
硼中间合金圈	87	钕铁硼永磁材料	416
金属钴圈	375	电解钴	2
钴片圈	297	电解钴	2
三元催化器圈	209	三元前驱体	102
硝酸锂圈	136	碳酸锂	106
氧化锂圈	16	氢氧化锂	98
废电动车电池圈	29	锂电池	17
海绵钛圈	289	钛酸锂	417
氮化钛圈	347	钛酸锂	417
金属钛圈	386	钛酸锂	417
氢化锂圈	96	锰酸锂	418
镍板圈	82	镍合金	83
镍片圈	225	镍合金	83
锰合金圈	151	铜锰合金	419
鼓风机圈	296	鼓风炉	420
电感应炉圈	337	电炉	328
时效炉圈	205	吹炼炉	421
淬火炉圈	351	吹炼炉	421
电动砂轮机圈	338	球磨机	215
除氧器圈	352	真空过滤机	291
除气箱圈	353	真空过滤机	291
电阻反射炉圈	329	反射炉	279
钻探机圈	161	扒矿机	422
煤气发生炉圈	228	真空炉	423
刨床圈	224	拉光机(剥皮机)	424
开平机圈	265	矫直抛光机	425
拉伸机圈	264	拉丝机	426
分卷机圈	304	叠锭机	427
轧辊磨床圈	173	铸轧机	428
冷凝器圈	261	冷轧机	259
矫平机圈	270	矫直机	269
热剪机圈	214	分切设备	303
分条机圈	302	分切设备	303
铣床圈	177	纵剪机	162
剪切机圈	272	纵剪机	162
锯床圈	266	横切机组	287
合卷机圈	288	拉弯矫	429
牵引机圈	219	叉车	430
蒸汽清洗机圈	167	除尘器	354
电子除垢器圈	331	除尘器	354
滤清器圈	238	空气净化器	299
防护用品圈	312	空气净化器	299
防护眼镜眼罩圈	313	空气净化器	299
防护鞋圈	314	空气净化器	299
防护手套圈	315	空气净化器	299
防护面罩面具圈	316	空气净化器	299
防护口罩圈	317	空气净化器	299
防护服圈	318	空气净化器	299
防护耳塞耳罩圈	319	空气净化器	299
安全帽圈	363	空气净化器	299
安全带安全绳圈	364	空气净化器	299
力学测试仪器圈	255	XRD分析检查	431
ROHS检测仪器圈	365	XRD分析检查	431
COD检测仪器圈	367	XRD分析检查	431
离子束分析仪器圈	256	扫描电子显微镜（SEM）	432
光学仪器圈	293	扫描电子显微镜（SEM）	432
紫外分析仪圈	163	粒度分析（激光粒度仪）	254
粉尘采样仪圈	301	粒度分析（激光粒度仪）	254
色谱仪圈	46	比表面积分析（比表面仪）	433`

	arr := strings.Split(str, "\n")

	fmt.Println(arr, len(arr))

	arr1 := make([]string, 0)
	for _,item := range arr {
		flag := false
		for _, s := range arr1 {
			if s == item {
				flag = true
			}
		}
		if !flag {
			arr1 = append(arr1, item)
		}
	}
	fmt.Println(arr1, len(arr1))
	formClubIds := make([]string, 0)
	toClubIds := make([]string, 0)
	for _, item := range arr1 {
		strs := strings.Split(item, "	")
		//formClub := strs[0]
		formClubID := strs[1]
		//toClub := strs[2]
		toClubID := strs[3]
		//fmt.Println(i, "从"+formClub+"("+formClubID+") >迁移到> "+toClub+"("+toClubID+")")

		formFlag := false
		for _,id := range formClubIds {
			if id == formClubID {
				formFlag = true
			}
		}
		if !formFlag {
			formClubIds = append(formClubIds, formClubID)
		}
		toFlag := false
		for _,id := range toClubIds {
			if id == toClubID {
				toFlag = true
			}
		}
		if !toFlag {
			toClubIds = append(toClubIds, toClubID)
		}
	}

	now := time.Now()
	beginTime := strconv.FormatInt(time.Date(now.Year(),now.Month(),now.Day(),0,0,0,0,time.Local).Unix(), 10)
	endTime := strconv.FormatInt(time.Date(now.Year(),now.Month(),now.Day(),0,0,0,0,time.Local).AddDate(0,0,1).Unix(), 10)

	for _, item := range formClubIds {
		sql := "UPDATE club SET post_amount=(SELECT COUNT(1) FROM post WHERE status=1 AND club_id="+item+"),"+
			"today_post_amount=(SELECT COUNT(1) FROM post WHERE status=1 AND updated_time>="+beginTime+
			" AND updated_time<"+endTime+" AND club_id="+item+") WHERE id="+item
		fmt.Println(sql)
	}
	fmt.Println("=======================================================")
	fmt.Println("=======================================================")
	fmt.Println("=======================================================")
	for _, item := range toClubIds {
		sql := "UPDATE club SET post_amount=(SELECT COUNT(1) FROM post WHERE status=1 AND club_id="+item+"),"+
			"today_post_amount=(SELECT COUNT(1) FROM post WHERE status=1 AND updated_time>="+beginTime+
			" AND updated_time<"+endTime+" AND club_id="+item+") WHERE id="+item
		fmt.Println(sql)
	}
	fmt.Println(len(formClubIds), strings.Join(formClubIds, ","))
	fmt.Println(len(toClubIds), strings.Join(toClubIds, ","))
}
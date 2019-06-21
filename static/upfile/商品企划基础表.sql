1、belle_sh.dm_online_sale_for_billz        					   // 线上每日销售数据和流量数据表;
3、belle_sh.hg_tbl_commodity_extend_prop_clean		           // 商品编码-风格映射表(货睿);									
4、belle_sh.hg_sku_commodity_style_mapping_formodeling_v2	    //款色编码-风格映射表(商品企划模型输入数据[郭浩手工维护])；
5、belle_sh.dm_commodity_styleno_commodity_first_date 			// 分平台首次销售日期
15、belle_jw.pro_channel_prop                         			// 商品编码-跟型映射表(货睿中间表)
16、belle_jw.pro_first_info										// 首次销售日期(货睿中间表)
17、belle_sh.hg_sku_commodity_style_mapping_v3                  // 款色编码-风格映射表(用于商品企划预测结果匹配[郭浩手工维护])；
18、belle_sh.mzj_store_level_bl_18                              // 店铺的级别
19、belle_jw.dm_online_inventory_for_bill						// 线上每日库存view 和 belle_sh库同属于 create view dm_online_inventory_for_bill as select * from  bi_report.dm_online_inventory_for_bill;
20、belle_jw.pro_month_lastday                                  // 每个月的最后一天 字典表(目前维护到2022年)
21、belle_sh.dw_purchase_pms_jy                                 // belle_sh.dw_purchase_pms_jy  每日分尺码未到数据; belle_jw.dw_purchase_pms_jy 不存在；
22、belle_jw.pro_commodity_style_info          					// 商品编码-风格映射表(货睿中间表)           
-- 23、belle_sh.hg_sku_newstyle_list						        // 二级分类-风格映射表(郭浩维护;临时表)
24、belle_sh.pro_spu_pyramid_level							    // belle_sh.pro_spu_pyramid_level 金字塔预测结果表(spu级别)
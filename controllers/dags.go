package controllers

import (
	"github.com/astaxie/beego"
	"github.com/upcode/models"
	"github.com/upcode/utils"
	"log"
	"strings"
)

type DagController struct {
	beego.Controller
}

func (d *DagController) Get() {
	d.TplName = "dags/d3-dags.html"
}

func (d *DagController) Dag() {
	sql := `drop table belle_sh.hg_shanghai_19_sale_inv_last_week_temp;
create table belle_sh.hg_shanghai_19_sale_inv_last_week_temp as
select *
from belle_sh.hg_shanghai_19_sale_inv;


drop table belle_sh.lrs_st_open_store_on_UI;
create table belle_sh.lrs_st_open_store_on_UI 
(store_no string)
ROW FORMAT DELIMITED
FIELDS TERMINATED BY ','
STORED AS TEXTFILE
;
insert into belle_sh.lrs_st_open_store_on_UI
values
('I0ZHST')
,('IAJ013')
,('I038ST')
,('I022ST')
,('I032ST')
,('I023ST')
,('I045ST')
,('I01IST')
,('I010ST')
,('I044ST')
,('I051ST')
,('IAJ002')
,('I01MST')
,('I01SST')
,('IAJ009')
,('I01BST')
,('IAJ004')
,('IAJ001')
,('I01JST')
,('IAJ003')
,('IAJ011')
,('I014ST')
,('I046ST')
,('I072ST')
,('I028ST')
,('IAJ014')
,('I078ST')
,('I01FST')
,('IAJ007')
,('I064ST')
,('I063ST')
,('I070ST')
,('I067ST')
,('I035ST')
,('IAJ006')
,('I029ST')
,('I033ST')
,('I01HST')
,('I074ST')
,('I066ST')
,('IAJ008')
,('I056ST')
,('IAJ012')
,('I048ST')
,('I01WST')
,('IA02MA')
,('IAJ016') -- map to I01MST
,('IAJA04') -- map to I046ST
;


drop table belle_sh.hg_shanghai_19_sale_inv;
create table belle_sh.hg_shanghai_19_sale_inv as
select
e.*,
e.store_no_tmp as store_no
from 
(select
	a.region_name
	,a.brand_no
	,(case 
	when a.store_no='IAJ016' then 'I01MST'
	when a.store_no='IAJA04' then 'I046ST'
	else a.store_no end) as store_no_tmp
	,a.store_name
	,a.store_category_name3
	,a.product_code
	,a.product_no
	,a.product_year_name
	,a.product_season_name
	,a.inv_date
	,a.in_inv
	,a.qty
	,a.amount
	,a.tag_amount
	,b.main_mtrl_name
	,b.heel_type_name
	,b.style_name
	,b.category_name3
	,b.style_number
	,case when a.size_no in ('225','225.0') then '225'
	  when a.size_no in ('230','230.0') then '230'
	  when a.size_no in ('235','235.0') then '235'
	  when a.size_no in ('240','240.0') then '240'
	  when a.size_no in ('245','245.0') then '245'
	  else a.size_no
	  end as size_no
from ( 
	select 
	aaa.*
	from 
		 (select 
          aa.*
          ,bb.store_name,bb.normal_province_name, bb.store_category_name3 
          from belle_sh.mzj_pmt_sku_store_avg_st_19 aa
		  left join (select distinct 
		  	         store_lno, store_name, normal_province_name, store_category_name3 
		  	         from data_belle.dim_org_store_allinfo 
		  	         where brand_no = 'ST'
		  	         ) bb
		  on aa.store_no=bb.store_lno
		 ) aaa
		inner join belle_sh.lrs_st_open_store_on_UI bbb
		on aaa.store_no = bbb.store_no
	)a
left join (
	select distinct
	product_no, main_mtrl_name, heel_type_name, style_name, category_name3,style_number
	from data_belle.dim_pro_allinfo 
	where category_flag!=0 and category_name1='鞋' and brand_no='ST'
	) b
on a.product_no=b.product_no
)e
;


drop table belle_sh.sy_hg_shanghai_19_sale_before_allstore;    
create table belle_sh.sy_hg_shanghai_19_sale_before_allstore as
select a.*,
	year(to_date(from_unixtime(unix_timestamp(a.inv_date_treat,'yyyyMMdd')))) as year,
	month(to_date(from_unixtime(unix_timestamp(a.inv_date_treat,'yyyyMMdd')))) as month,
	weekofyear(to_date(from_unixtime(unix_timestamp(a.inv_date_treat,'yyyyMMdd')))) as week,
	ceiling(weekofyear(to_date(from_unixtime(unix_timestamp(a.inv_date_treat,'yyyyMMdd'))))/2) as biweek,
	amount/tag_amount as discount
from 
(
	select
	*,
	if(main_mtrl_name like '%皮%',1,0) as ind_leather,
	(case
	when inv_date in('20181230','20181231') then '20181229'
	when inv_date = '20190101' then '20181230'
	when inv_date in ('20171230','20171231') then '20171230'
	when inv_date ='20180101' then '20171231'
	else inv_date end) as inv_date_treat
	from belle_sh.hg_shanghai_19_sale_inv
) a
;



drop table belle_sh.sy_tmp1_sig1_19;
create table belle_sh.sy_tmp1_sig1_19 as
select 
store_no,
product_code,
product_no,
product_year_name,
product_season_name,
min(to_date(from_unixtime(unix_timestamp(inv_date,'yyyyMMdd')))) as first_online_date,
date_add(min(to_date(from_unixtime(unix_timestamp(inv_date,'yyyyMMdd')))),30) as first_online_date_after30D,
date_add(min(to_date(from_unixtime(unix_timestamp(inv_date,'yyyyMMdd')))),60) as first_online_date_after60D,
date_add(min(to_date(from_unixtime(unix_timestamp(inv_date,'yyyyMMdd')))),90) as first_online_date_after90D,
date_add(min(to_date(from_unixtime(unix_timestamp(inv_date,'yyyyMMdd')))),120) as first_online_date_after120D,
date_add(min(to_date(from_unixtime(unix_timestamp(inv_date,'yyyyMMdd')))),150) as first_online_date_after150D,
date_add(min(to_date(from_unixtime(unix_timestamp(inv_date,'yyyyMMdd')))),180) as first_online_date_after180D
from belle_sh.hg_shanghai_19_sale_inv
group by store_no,product_code,product_no,product_year_name,product_season_name;


drop table belle_sh.sy_sku_base_first_online_date_19_mapping;
create table belle_sh.sy_sku_base_first_online_date_19_mapping as
select 
store_no,
product_code,
product_no,
product_year_name,
product_season_name,
first_online_date,
first_online_date_after30D,
first_online_date_after60D,
first_online_date_after90D,
first_online_date_after120D,
weekofyear(to_date(from_unixtime(unix_timestamp(first_online_date),'yyyy-MM-dd'))) as first_online_week,
weekofyear(to_date(from_unixtime(unix_timestamp(first_online_date_after30D),'yyyy-MM-dd'))) as first_online_30D_week,
weekofyear(to_date(from_unixtime(unix_timestamp(first_online_date_after60D),'yyyy-MM-dd'))) as first_online_60D_week,
weekofyear(to_date(from_unixtime(unix_timestamp(first_online_date_after90D),'yyyy-MM-dd'))) as first_online_90D_week,
weekofyear(to_date(from_unixtime(unix_timestamp(first_online_date_after120D),'yyyy-MM-dd'))) as first_online_120D_week,
weekofyear(to_date(from_unixtime(unix_timestamp(first_online_date_after150D),'yyyy-MM-dd'))) as first_online_150D_week,
weekofyear(to_date(from_unixtime(unix_timestamp(first_online_date_after180D),'yyyy-MM-dd'))) as first_online_180D_week,
year(first_online_date)*100+weekofyear(first_online_date)                     as first_online_yearweek,
year(first_online_date_after30D)*100+weekofyear(first_online_date_after30D)   as first_online_30D_yearweek,
year(first_online_date_after60D)*100+weekofyear(first_online_date_after60D)   as first_online_60D_yearweek,
year(first_online_date_after90D)*100+weekofyear(first_online_date_after90D)   as first_online_90D_yearweek,
year(first_online_date_after120D)*100+weekofyear(first_online_date_after120D) as first_online_120D_yearweek,
year(first_online_date_after150D)*100+weekofyear(first_online_date_after150D) as first_online_150D_yearweek,
year(first_online_date_after180D)*100+weekofyear(first_online_date_after180D) as first_online_180D_yearweek
from belle_sh.sy_tmp1_sig1_19;


set hive.support.quoted.identifiers=None;
drop table belle_sh.sy_19_table_tmp1;
create table belle_sh.sy_19_table_tmp1 as
select
$(qty|in_inv|amount|tag_amount)?+.+$,
nvl(greatest(qty,0),0) as qty,
nvl(greatest(amount,0),0) as amount,
nvl(greatest(tag_amount,0),0) as tag_amount,
nvl(in_inv,0) as in_inv
from belle_sh.sy_hg_shanghai_19_sale_before_allstore
; 

drop table belle_sh.sy_19_store_sku_lvl_SR_monthly;
create table belle_sh.sy_19_store_sku_lvl_SR_monthly as    
select 
store_no,
product_code,
product_year_name,
product_season_name,
year,
week,
sum(nvl(greatest(qty,0),0)) as qty,
sum(greatest(amount,0)) as sum_amt,
sum(greatest(tag_amount,0)) as sum_tag_amt,
count(distinct(inv_date)) as sale_days,
sum(greatest(qty,0))/count(distinct(inv_date)) as sale_rate_daily_size_lvl_monthly
from belle_sh.sy_19_table_tmp1
group by store_no,product_code,product_year_name,product_season_name,year,week
;


drop table belle_sh.sy_19_store_sku_lvl_extend_tmp1;
create table belle_sh.sy_19_store_sku_lvl_extend_tmp1 as    
select
aaa.store_no,
aaa.product_code,
aaa.product_year_name,
aaa.product_season_name,
aaa.year,
aaa.key,
b.week
from(
select
aa.store_no,
aa.product_code,
aa.product_year_name,
aa.product_season_name,
aa.year,
1 as key
from(
select a.*,b.product_code,b.product_year_name,b.product_season_name,b.key_product
from
	(select distinct store_no,year,1 as key_store
	from belle_sh.sy_19_store_sku_lvl_SR_monthly) a
left join 
	(select distinct year,product_code,product_year_name,product_season_name,1 as key_product
	from belle_sh.sy_19_store_sku_lvl_SR_monthly) b
on a.key_store=b.key_product and a.year=b.year
) aa)aaa
left join belle_sh.week_mapping b
on aaa.key=b.key
;


drop table belle_sh.sy_19_store_sku_lvl_extend_allweek;
create table belle_sh.sy_19_store_sku_lvl_extend_allweek as    
select
a.store_no,
a.product_code,
a.product_year_name,
a.product_season_name,
a.year,
a.week,
nvl(b.in_inv_ind,0) as in_inv_ind, 
nvl(greatest(b.qty,0),0) as qty,
nvl(greatest(b.sum_amt,0),0) as sum_amt,
nvl(greatest(b.sum_tag_amt,0),0) as sum_tag_amt,
nvl(b.sale_days,0) as sale_days,
nvl(b.sale_rate_daily_size_lvl_monthly,0) as sr_size_lvl
from belle_sh.sy_19_store_sku_lvl_extend_tmp1 a left join
(select
*,
1 as in_inv_ind
from belle_sh.sy_19_store_sku_lvl_SR_monthly) b
on 
a.store_no=b.store_no and a.product_code=b.product_code 
and a.year=b.year and a.week=b.week
where (a.year*100+a.week)<=${hiveconf:yearweeknum_upperlimit} and (a.year*100+a.week)>=(${hiveconf:product_year_name}-1)*100+49
;


drop table belle_sh.sy_19_store_sku_lvl_SR_monthly_lag_sum;
create table belle_sh.sy_19_store_sku_lvl_SR_monthly_lag_sum as     
select
a.*,
(a.in_inv_ind+a.in_inv_ind_lag_1+a.in_inv_ind_lag_2+a.in_inv_ind_lag_3) as in_inv_ind_last_4_week,
(a.in_inv_ind+a.in_inv_ind_lag_1) as in_inv_ind_last_2_week,
(a.qty+a.qty_lag_1+a.qty_lag_2+a.qty_lag_3) as qty_last_4_week,
(a.sum_amt+a.sum_amt_lag_1+a.sum_amt_lag_2+a.sum_amt_lag_3) as sum_amt_last_4_week,
(a.sum_tag_amt+a.sum_tag_amt_lag_1+a.sum_tag_amt_lag_2+a.sum_tag_amt_lag_3) as sum_tag_amt_last_4_week,
(a.sale_days+a.sale_days_lag_1+a.sale_days_lag_2+a.sale_days_lag_3) as sale_days_last_4_week,
(a.qty+a.qty_lag_1) as qty_last_2_week,
(a.sum_amt+a.sum_amt_lag_1) as sum_amt_last_2_week,
(a.sum_tag_amt+a.sum_tag_amt_lag_1) as sum_tag_amt_last_2_week,
nvl((a.sale_days+a.sale_days_lag_1),0) as sale_days_last_2_week,
nvl((a.qty/a.sale_days),0) as sale_rate,
--nvl((a.qty+a.qty_lag_1)/(a.sale_days+a.sale_days_lag_1),0) as sale_rate_last_2_week,
(case when a.qty - a.qty_lag_1 > -2 and a.qty - a.qty_lag_1 < 2 then nvl((a.qty+a.qty_lag_1)/(a.sale_days+a.sale_days_lag_1),0)
   else nvl((least(nvl(a.qty,0), nvl(a.qty_lag_1,0))*2 + 2)/(a.sale_days+a.sale_days_lag_1),0)
   end ) as sale_rate_last_2_week,
nvl((a.qty+a.qty_lag_1+a.qty_lag_2+a.qty_lag_3)/(a.sale_days+a.sale_days_lag_1+a.sale_days_lag_2+a.sale_days_lag_3),0) as sale_rate_last_4_week,
nvl((a.cumu_qty_last_1_week/a.cumu_sale_days_last_1_week),0) as cumu_sale_rate
from
(
select 
*,
lag(in_inv_ind,1,0)  over(partition by store_no,product_code order by year,week asc) as in_inv_ind_lag_1,
lag(in_inv_ind,2,0)  over(partition by store_no,product_code order by year,week asc) as in_inv_ind_lag_2,
lag(in_inv_ind,3,0)  over(partition by store_no,product_code order by year,week asc) as in_inv_ind_lag_3,
lag(qty,1,0)         over(partition by store_no,product_code order by year,week asc) as qty_lag_1,
lag(qty,2,0)         over(partition by store_no,product_code order by year,week asc) as qty_lag_2,
lag(qty,3,0)         over(partition by store_no,product_code order by year,week asc) as qty_lag_3,
lag(sum_amt,1,0)     over(partition by store_no,product_code order by year,week asc) as sum_amt_lag_1,
lag(sum_amt,2,0)     over(partition by store_no,product_code order by year,week asc) as sum_amt_lag_2,
lag(sum_amt,3,0)     over(partition by store_no,product_code order by year,week asc) as sum_amt_lag_3,
lag(sum_tag_amt,1,0) over(partition by store_no,product_code order by year,week asc) as sum_tag_amt_lag_1,
lag(sum_tag_amt,2,0) over(partition by store_no,product_code order by year,week asc) as sum_tag_amt_lag_2,
lag(sum_tag_amt,3,0) over(partition by store_no,product_code order by year,week asc) as sum_tag_amt_lag_3,
lag(sale_days,1,0)   over(partition by store_no,product_code order by year,week asc) as sale_days_lag_1,
lag(sale_days,2,0)   over(partition by store_no,product_code order by year,week asc) as sale_days_lag_2,
lag(sale_days,3,0)   over(partition by store_no,product_code order by year,week asc) as sale_days_lag_3,
sum(qty)             over(partition by store_no,product_code order by year,week asc) as cumu_qty_last_1_week,
sum(sale_days)       over(partition by store_no,product_code order by year,week asc) as cumu_sale_days_last_1_week,
sum(in_inv_ind)      over(partition by store_no,product_code)                        as historical_exist_ind
from belle_sh.sy_19_store_sku_lvl_extend_allweek) a
;


drop table belle_sh.sy_19_store_sku_lvl_weekly_online;
create table belle_sh.sy_19_store_sku_lvl_weekly_online as    
select
a.store_no,
a.product_code,
a.product_year_name,
a.product_season_name,
a.year,
a.week,
a.in_inv_ind,
a.in_inv_ind_last_2_week,
b.first_online_date,
b.first_online_week,
c.first_online_date_SH,
c.first_online_week_SH
from belle_sh.sy_19_store_sku_lvl_SR_monthly_lag_sum a 
left join belle_sh.sy_sku_base_first_online_date_19_mapping b
on a.store_no=b.store_no and a.product_code=b.product_code
left join 
	(select
	store_no,
	product_code,
	min(first_online_date) as first_online_date_SH,
	min(first_online_week) as first_online_week_SH
	from belle_sh.sy_sku_base_first_online_date_19_mapping
	group by store_no,product_code) c
on a.store_no=c.store_no and a.product_code=c.product_code
;


drop table belle_sh.sy_19_sku_sunday_inv;
create table belle_sh.sy_19_sku_sunday_inv as    
select 
a.*,
lag(nvl(a.sunday_inv_end,0),2,0) over(partition by a.product_code order by a.year asc,a.week asc) as sunday_inv_end_last_2_week,
lag(nvl(a.sunday_inv_end,0),1,0) over(partition by a.product_code order by a.year asc,a.week asc) as sunday_inv_end_last_1_week
from
(select
product_code,
product_year_name,
product_season_name,
year,
week,
sum(nvl(greatest(in_inv,0),0)) as sunday_inv_end
from belle_sh.sy_19_table_tmp1
where pmod(datediff(to_date(from_unixtime(unix_timestamp(inv_date,'yyyyMMdd'))), '2017-12-31'), 7)=0
group by product_code,product_year_name,product_season_name,year,week
) a
;

drop table belle_sh.sy_19_SH_sku_lag_qty;
create table belle_sh.sy_19_SH_sku_lag_qty as   
select
a.product_code,
a.product_year_name,a.product_season_name,
a.year,
a.week,
a.SH_in_inv_ind,
a.qty,
sum(a.qty) over (partition by a.product_code order by a.year asc,a.week asc ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) as qty_last_2_week,
sum(a.qty) over (partition by a.product_code order by a.year asc,a.week asc ROWS BETWEEN 3 PRECEDING AND CURRENT ROW) as qty_last_4_week,
sum(a.qty) over (partition by a.product_code order by a.year asc,a.week asc ROWS BETWEEN CURRENT ROW AND 3 FOLLOWING) as qty_post_4_week,
sum(a.qty) over(partition by a.product_code order by a.year asc,a.week asc) as cumu_qty
from (select 
product_code,product_year_name,product_season_name,year,week,
sum(qty) as qty,
sum(in_inv_ind) as SH_in_inv_ind
from belle_sh.sy_19_store_sku_lvl_extend_allweek
group by product_code,product_year_name,product_season_name,year,week) a
;


drop table belle_sh.sy_19_SH_sku_sell_through_rate;
create table belle_sh.sy_19_SH_sku_sell_through_rate as   
select
a.product_code,
a.product_year_name,a.product_season_name,
a.year,
a.week,
a.SH_in_inv_ind,
a.qty,
a.qty_last_4_week,
a.cumu_qty,
b.sunday_inv_end,
b.sunday_inv_end_last_1_week,
b.sunday_inv_end_last_2_week,
nvl((b.sunday_inv_end+a.qty-b.sunday_inv_end_last_1_week),0) as inv_add_ind,
nvl(a.qty_last_4_week/b.sunday_inv_end,0) as STR_last_4_week,
nvl(a.qty_last_2_week/b.sunday_inv_end,0) as STR_last_2_week,
nvl(a.cumu_qty/b.sunday_inv_end,0)        as STR_cumu
from belle_sh.sy_19_SH_sku_lag_qty a
left join belle_sh.sy_19_sku_sunday_inv b
on a.product_code=b.product_code and a.year=b.year and a.week=b.week
where a.SH_in_inv_ind>0
;



drop table belle_sh.sy_19_store_sku_lvl_SR_weekly_allseason;
create table belle_sh.sy_19_store_sku_lvl_SR_weekly_allseason as    
select 
store_no,
product_code,
year,
week,
sum(in_inv_ind_last_4_week) as sku_in_inv_ind_last_4_week,
sum(qty) as qty,
sum(qty_lag_1) as qty_lag_1,
sum(qty_lag_2) as qty_lag_2,
sum(qty_lag_3) as qty_lag_3,
sum(qty_last_4_week) as qty_last_4_week,
sum(sum_amt_last_4_week) as sum_amt_last_4_week,
sum(sum_tag_amt_last_4_week) as sum_tag_amt_last_4_week,
sum(sale_rate_last_4_week) as sale_rate_last_4_week,
sum(sale_rate) as sale_rate,
sum(qty_last_2_week) as qty_last_2_week,
sum(sale_rate_last_2_week) as sale_rate_last_2_week
from belle_sh.sy_19_store_sku_lvl_SR_monthly_lag_sum
where historical_exist_ind>0
group by store_no,product_code,year,week
;


drop table belle_sh.sy_19_store_lvl_SR_weekly_rank_all_season;
create table belle_sh.sy_19_store_lvl_SR_weekly_rank_all_season as    
select 
*,
if(sku_in_inv_ind_last_4_week>0,1,0) as sku_exist_ind_last_4_week,
rank() over(partition by store_no,year,week order by sale_rate_last_4_week desc,qty_last_4_week desc,sum_amt_last_4_week desc,sum_tag_amt_last_4_week desc) as sku_sr_rank,
row_number() over(partition by store_no,year,week order by sale_rate_last_4_week desc,qty_last_4_week desc,sum_amt_last_4_week desc,sum_tag_amt_last_4_week desc) as sku_sr_rank_row_number
from belle_sh.sy_19_store_sku_lvl_SR_weekly_allseason;


drop table belle_sh.sy_19_store_lvl_top40_avg_SR_weekly_all_season;
create table belle_sh.sy_19_store_lvl_top40_avg_SR_weekly_all_season as    
select
a.*,
(a.sale_rate_last_4_week/a.cnt_distinct_sku_last_4_week) as avg_store_last_4_week_sr
from(
select
store_no,
year,
week,
sum(qty) as qty,
sum(sale_rate) as sale_rate,
sum(qty_last_4_week) as qty_last_4_week,
sum(sale_rate_last_4_week) as sale_rate_last_4_week,
sum(sku_exist_ind_last_4_week) as cnt_distinct_sku_last_4_week,
count(distinct(product_code)) as cnt_distinct_sku
from belle_sh.sy_19_store_lvl_SR_weekly_rank_all_season
where sku_sr_rank_row_number<=40
group by store_no,year,week) a;


drop table belle_sh.sy_19_store_lvl_top40_avg_SR_weekly_normed_all_season;
create table belle_sh.sy_19_store_lvl_top40_avg_SR_weekly_normed_all_season as    
select
a.store_no,a.year,a.week,
a.qty,a.sale_rate,a.qty_last_4_week,a.sale_rate_last_4_week,
a.cnt_distinct_sku_last_4_week,a.cnt_distinct_sku,
nvl(a.avg_store_last_4_week_sr,0) as avg_store_last_4_week_sr,
b.mean,
b.std,
nvl(((a.avg_store_last_4_week_sr-b.mean)/b.std),0) as normed_avg_store_mthly_sr
from belle_sh.sy_19_store_lvl_top40_avg_SR_weekly_all_season a
left join 
(select
year,
week,
avg(avg_store_last_4_week_sr) as mean,
stddev_pop(avg_store_last_4_week_sr) as std
from belle_sh.sy_19_store_lvl_top40_avg_SR_weekly_all_season
group by year,week) b 
on a.year=b.year and a.week=b.week;


drop table belle_sh.lrs_19_store_sku_lvl_SR_weekly_allseason;
create table belle_sh.lrs_19_store_sku_lvl_SR_weekly_allseason as    
select 
store_no,
product_code,
product_season_name,
year,
week,
sum(in_inv_ind_last_2_week) as sku_in_inv_ind_last_2_week,
sum(qty_last_2_week) as qty_last_2_week,
sum(sum_amt_last_2_week) as sum_amt_last_2_week,
sum(sum_tag_amt_last_2_week) as sum_tag_amt_last_2_week,
sum(sale_rate_last_2_week) as sale_rate_last_2_week
from belle_sh.sy_19_store_sku_lvl_SR_monthly_lag_sum
where historical_exist_ind>0
group by store_no,product_code,year,week,product_season_name;


drop table belle_sh.lrs_19_store_lvl_SR_weekly_rank_all_season_2week;
create table belle_sh.lrs_19_store_lvl_SR_weekly_rank_all_season_2week as    
select 
*,
if(sku_in_inv_ind_last_2_week>0,1,0) as sku_exist_ind_last_2_week,
rank() over(partition by store_no,year,week,product_season_name order by sale_rate_last_2_week asc,qty_last_2_week asc,sum_amt_last_2_week asc,sum_tag_amt_last_2_week asc) as sku_sr_rank,
row_number() over(partition by store_no,year,week,product_season_name order by sale_rate_last_2_week asc,qty_last_2_week asc,sum_amt_last_2_week asc,sum_tag_amt_last_2_week asc) as sku_sr_rank_row_number
from belle_sh.lrs_19_store_sku_lvl_SR_weekly_allseason;


drop table belle_sh.lrs_19_store_lvl_SR_weekly_rank_all_season_2week_percentile;
create table belle_sh.lrs_19_store_lvl_SR_weekly_rank_all_season_2week_percentile as
select
a.*
from belle_sh.lrs_19_store_lvl_SR_weekly_rank_all_season_2week a
inner join
(select
bb.store_no
,bb.year
,bb.week
,bb.product_season_name
,ceil(0.95*max(bb.sku_sr_rank_row_number)) as sr_percentile_9_5
from belle_sh.lrs_19_store_lvl_SR_weekly_rank_all_season_2week bb
group by store_no,year,week,product_season_name
) b
on a.store_no = b.store_no and a.year = b.year and a.week = b.week and a.product_season_name = b.product_season_name and a.sku_sr_rank_row_number = b.sr_percentile_9_5;


drop table belle_sh.sy_19_store_sku_lvl_SR_weekly_before_cap;
create table belle_sh.sy_19_store_sku_lvl_SR_weekly_before_cap as    
select 
store_no,
product_code,
product_year_name,
product_season_name,
year,
week,
sum(nvl(historical_exist_ind,0)) as sku_historical_exist_ind,
sum(in_inv_ind_last_4_week) as sku_in_inv_ind_last_4_week,
sum(in_inv_ind_last_2_week) as sku_in_inv_ind_last_2_week,
sum(qty) as qty,
sum(qty_lag_1) as qty_lag_1,
sum(qty_lag_2) as qty_lag_2,
sum(qty_lag_3) as qty_lag_3,
sum(qty_last_4_week) as qty_last_4_week,
sum(sum_amt_last_4_week) as sum_amt_last_4_week,
sum(sum_tag_amt_last_4_week) as sum_tag_amt_last_4_week,
sum(sale_rate_last_4_week) as sale_rate_last_4_week,
sum(sale_rate) as sale_rate,
sum(qty_last_2_week) as qty_last_2_week,
sum(sale_rate_last_2_week) as sale_rate_last_2_week,
sum(sum_amt) as sum_amt,
sum(sum_tag_amt) as sum_tag_amt,
(sum(sum_amt)/sum(sum_tag_amt)) as discount_week_lvl,
(sum(sum_amt)/sum(qty)) as avg_sale_amt,
(sum(sum_tag_amt)/sum(qty)) as avg_tag_amount
from belle_sh.sy_19_store_sku_lvl_SR_monthly_lag_sum
group by store_no,product_code,product_year_name,product_season_name,year,week
;


drop table belle_sh.sy_19_store_sku_lvl_SR_weekly;
create table belle_sh.sy_19_store_sku_lvl_SR_weekly as
select
a.store_no,
a.product_code,
a.product_year_name,
a.product_season_name,
a.year,
a.week,
a.sku_historical_exist_ind,
a.sku_in_inv_ind_last_4_week,
a.sku_in_inv_ind_last_2_week,
a.qty,
a.qty_lag_1,
a.qty_lag_2,
a.qty_lag_3,
a.qty_last_4_week,
a.sum_amt_last_4_week,
a.sum_tag_amt_last_4_week,
a.sale_rate_last_4_week,
a.sale_rate,
a.qty_last_2_week,
a.sale_rate_last_2_week as sale_rate_last_2_week_before_cap,
if(a.sale_rate_last_2_week >= 1, b.sale_rate_last_2_week, a.sale_rate_last_2_week) as sale_rate_last_2_week,
a.sum_amt,
a.sum_tag_amt,
a.discount_week_lvl,
a.avg_sale_amt,
a.avg_tag_amount
from belle_sh.sy_19_store_sku_lvl_SR_weekly_before_cap a
left join belle_sh.lrs_19_store_lvl_SR_weekly_rank_all_season_2week_percentile b
on a.store_no = b.store_no and a.year = b.year and a.week = b.week and a.product_season_name = b.product_season_name;


drop table belle_sh.sy_19_store_sku_lvl_SR_weekly_normed;
create table belle_sh.sy_19_store_sku_lvl_SR_weekly_normed as    
select 
a.*,
nvl(((a.sale_rate_last_2_week-b.mean_amg_store)/b.std_amg_store),0) as normed_sr_last_2_week,
(case
when a.discount_week_lvl >=0.7 then '>=7'
when a.discount_week_lvl >=0.6 and a.discount_week_lvl <0.7  then '[6,7)'
when a.discount_week_lvl >=0.5 and a.discount_week_lvl <0.6  then '[5,6)'
when a.discount_week_lvl >=0 and a.discount_week_lvl <0.5  then '<5'
else 'others' end) as discount_bin,
(case
when a.avg_sale_amt >=1000 then '>=1000'
when a.avg_sale_amt >=800 and a.avg_sale_amt <1000 then '800-1000'
when a.avg_sale_amt >=600 and a.avg_sale_amt <800 then '600-800'
when a.avg_sale_amt <600 then '<600'
else 'null' end) as avg_sale_amt_bin,
(case
when a.avg_sale_amt >=1600 then '>=1600'
when a.avg_sale_amt >=1400 and a.avg_sale_amt <1600 then '1400-1600'
when a.avg_sale_amt >=1000 and a.avg_sale_amt <1400 then '1000-1400'
when a.avg_sale_amt >=800 and a.avg_sale_amt <1000 then '800-1000'
when a.avg_sale_amt <800 then '<800'
else 'null' end) as avg_sale_amt_bin_win
from (select * from belle_sh.sy_19_store_sku_lvl_SR_weekly) a 
left join
	(
	select
	product_code,product_year_name,product_season_name,year,week,
	avg(sale_rate_last_2_week) as mean_amg_store,
	stddev_pop(sale_rate_last_2_week) as std_amg_store
	from belle_sh.sy_19_store_sku_lvl_SR_weekly
	group by product_code,product_year_name,product_season_name,year,week
	) b
on a.product_code=b.product_code and a.product_year_name=b.product_year_name and a.product_season_name=b.product_season_name and a.year=b.year and a.week=b.week
;


drop table belle_sh.sy_sku_mapping_table_19;
create table belle_sh.sy_sku_mapping_table_19 as 
select
a.*,
b.category_name2,
b.category_name3,
b.category_name4
from(
select
aa.*,
bb.heel_shape
from(
select
product_code,
product_no,
product_year_name,
product_season_name,
heel_type_name,
style_name,
main_mtrl_name,
sum(greatest(tag_amount,0))/sum(greatest(qty,0)) as avg_tag_amount_SH,
sum(greatest(amount,0))/sum(greatest(qty,0)) as avg_amt_SH,
sum(greatest(amount,0)/greatest(tag_amount,0))/sum(greatest(qty,0)) as avg_discount_SH
from belle_sh.sy_hg_shanghai_19_sale_before_allstore
group by product_code,product_no,product_year_name,product_season_name,heel_type_name,style_name,main_mtrl_name) aa 
left join belle_sh.sy_sku_heel_shape bb
on aa.product_code=bb.product_code)
a left join data_belle.dim_pro_allinfo b
on a.product_code=b.product_code;


drop table belle_sh.sy_sku_mapping_table_19_group;
create table belle_sh.sy_sku_mapping_table_19_group as 
select
*,
(case 
when category_name3='特长靴' then '热'
when category_name3='高靴' then '热'
when category_name3='中靴' then '热'
when category_name3='中高靴' then '热'
when category_name3='低靴' then '热'
when category_name3='满帮' then '满帮'
when category_name3='浅口' then '浅口'
when category_name3='凉靴' then '凉'
when category_name3='中空' then '凉'
when category_name3='前后空' then '凉'
when category_name3='前空' then '凉'
when category_name3='后空' then '凉'
when category_name3='凉拖' then '凉'
else '其他' end) as category_name3_2,
(case 
when category_name3='特长靴' then '高靴'
when category_name3='高靴' then '高靴'
when category_name3='中靴' then '中靴'
when category_name3='中高靴' then '中靴'
when category_name3='低靴' then '低靴'
when category_name3='满帮' then '满帮'
when category_name3='浅口' then '浅口'
when category_name3='凉靴' then '凉X'
when category_name3='中空' then '空'
when category_name3='前后空' then '空'
when category_name3='前空' then '空'
when category_name3='后空' then '空'
when category_name3='凉拖' then '凉X'
else '其他' end) as category_name3_3_Sum,
(case
when heel_shape='坡跟' then '坡跟'
when heel_shape like '平底%' then '平底'
when heel_shape='松糕' then '厚底'
when heel_shape='水台' then '厚底'
when heel_shape='粗跟' then '粗跟'
when heel_shape='细跟' then '细跟'
else '其他' end) as heel_shape2,
if(heel_type_name='平','低',heel_type_name) as heel_type_name2,
(case
when avg_discount_SH >=0.7 then '>=7'
when avg_discount_SH >=0.6 and avg_discount_SH <0.7  then '[6,7)'
when avg_discount_SH >=0.5 and avg_discount_SH <0.6  then '[5,6)'
when avg_discount_SH >=0 and avg_discount_SH <0.5  then '<5'
else 'others' end) as discount_bin_SH,
(case
when avg_amt_SH >=1000 then '>=1000'
when avg_amt_SH >=800 and avg_amt_SH <1000 then '800-1000'
when avg_amt_SH >=600 and avg_amt_SH <800 then '600-800'
when avg_amt_SH <600 then '<600'
else 'null' end) as avg_sale_amt_bin,
(case
when avg_amt_SH >=1600 then '>=1600'
when avg_amt_SH >=1400 and avg_amt_SH <1600 then '1400-1600'
when avg_amt_SH >=1000 and avg_amt_SH <1400 then '1000-1400'
when avg_amt_SH >=800 and avg_amt_SH <1000 then '800-1000'
when avg_amt_SH <800 then '<800'
else 'null' end) as avg_sale_amt_bin_win
from belle_sh.sy_sku_mapping_table_19;



drop table belle_sh.sy_19_sku_lvl_weekly_discount_sale_amt_mapping_SH;
create table belle_sh.sy_19_sku_lvl_weekly_discount_sale_amt_mapping_SH as    
select 
aa.*,
b.heel_type_name,
b.heel_type_name2,
b.style_name,
b.category_name3_2,
b.category_name3_3_Sum,
b.heel_shape,
b.heel_shape2
from (
select
a.*,
(case
when a.discount_week_lvl >=0.7 then '>=7'
when a.discount_week_lvl >=0.6 and a.discount_week_lvl <0.7  then '[6,7)'
when a.discount_week_lvl >=0.5 and a.discount_week_lvl <0.6  then '[5,6)'
when a.discount_week_lvl >=0 and a.discount_week_lvl <0.5  then '<5'
else 'others' end) as discount_bin_SH,
(case
when a.avg_sale_amt >=1000 then '>=1000'
when a.avg_sale_amt >=800 and a.avg_sale_amt <1000 then '800-1000'
when a.avg_sale_amt >=600 and a.avg_sale_amt <800 then '600-800'
when a.avg_sale_amt <600 then '<600'
else 'null' end) as avg_sale_amt_bin_SH,
(case
when a.avg_sale_amt >=1600 then '>=1600'
when a.avg_sale_amt >=1400 and a.avg_sale_amt <1600 then '1400-1600'
when a.avg_sale_amt >=1000 and a.avg_sale_amt <1400 then '1000-1400'
when a.avg_sale_amt >=800 and a.avg_sale_amt <1000 then '800-1000'
when a.avg_sale_amt <800 then '<800'
else 'null' end) as avg_sale_amt_bin_SH_win
from
(
select
year,
week,
product_code,
product_year_name,
product_season_name,
sum(qty) as qty,
sum(sum_amt) as sum_amt,
sum(sum_tag_amt) as sum_tag_amt,
(sum(sum_amt)/sum(sum_tag_amt)) as discount_week_lvl,
(sum(sum_amt)/sum(qty)) as avg_sale_amt,
(sum(sum_tag_amt)/sum(qty)) as avg_tag_amount
from belle_sh.sy_19_store_sku_lvl_SR_weekly_normed
group by year,week,product_code,product_year_name,product_season_name) a ) aa
left join belle_sh.sy_sku_mapping_table_19_group b
on aa.product_code=b.product_code;



drop table belle_sh.sy_19_sku_lvl_weekly_sr_include_SH;
create table belle_sh.sy_19_sku_lvl_weekly_sr_include_SH as    
select 
a.*,
(case
when a.offline_disc >=0.7 then '>=7'
when a.offline_disc >=0.6 and a.offline_disc <0.7  then '[6,7)'
when a.offline_disc >=0.5 and a.offline_disc <0.6  then '[5,6)'
when a.offline_disc >=0 and a.offline_disc <0.5  then '<5'
else 'others' end) as discount_bin_treatment_lead_1,
(case
when a.price >=1000 then '>=1000'
when a.price >=800 and a.price <1000 then '800-1000'
when a.price >=600 and a.price <800 then '600-800'
when a.price <600 and a.price>0 then '<600'
else 'null' end) as avg_amt_bin_treatment_lead_1,
(case
when a.price >=1600 then '>=1600'
when a.price >=1400 and a.price <1600 then '1400-1600'
when a.price >=1000 and a.price <1400 then '1000-1400'
when a.price >=800 and a.price <1000 then '800-1000'
when a.price <800 then '<800'
else 'null' end) as avg_amt_bin_treatment_lead_1_win,
b.heel_type_name2,
b.style_name,
b.category_name3_2,
b.category_name3_3_Sum,
b.discount_bin_SH,
b.avg_sale_amt_bin_SH,
b.avg_sale_amt_bin_SH_win
from(
select
a1.store_no,a1.product_code,a1.product_year_name,a1.product_season_name,a1.year,a1.week,
a1.sku_historical_exist_ind,
a1.sku_in_inv_ind_last_2_week,
a1.sku_in_inv_ind_last_4_week,
a1.qty_last_2_week,
a1.qty_last_4_week,
a1.normed_sr_last_2_week,
a1.sale_rate_last_2_week,
a1.sale_rate_last_4_week,
a1.discount_bin,
a1.avg_sale_amt_bin,
a1.avg_sale_amt_bin_win,
a2.offline_disc,
a2.price
from (select * from belle_sh.sy_19_store_sku_lvl_SR_weekly_normed where (year*100+week)=${hiveconf:yearweeknum_upperlimit}) a1
left join 
    (
	select * from belle_sh.yl_sku_offline_disc_spr_smr_2019
	) a2
on a1.product_code=a2.product_code and a1.year=a2.year and a1.week=a2.week) a
left join belle_sh.sy_19_sku_lvl_weekly_discount_sale_amt_mapping_SH b
on a.product_code=b.product_code and a.year=b.year and a.week=b.week
;



drop table belle_sh.sy_19_sku_lvl_weekly_all_info_tmp;
create table belle_sh.sy_19_sku_lvl_weekly_all_info_tmp as    
select 
aa.*,
row_number() over(partition by aa.store_no,aa.product_code,aa.year,aa.week order by aa.avg_store_last_4_week_sr desc,aa.normed_avg_store_mthly_sr desc) as duplicate_rank,
bb.next_week_max_month_tmp as next_week_max_month,
bb.next_week_year,
bb.week_first_date,
bb.week_last_date
from 
(select
a.store_no,
a.product_code,
a.product_year_name,
a.product_season_name,
a.year,
a.week,
a.sku_historical_exist_ind,
a.sku_in_inv_ind_last_2_week,
a.sku_in_inv_ind_last_4_week,
a.normed_sr_last_2_week,
a.sale_rate_last_2_week,
a.qty_last_2_week,
a.qty_last_4_week,
if(a.offline_disc<0.5,1,0) as discount_less_than_5,
a.price,
a.discount_bin,
a.discount_bin_SH,
a.avg_sale_amt_bin,
a.avg_sale_amt_bin_SH,
a.avg_sale_amt_bin_SH_win,
a.discount_bin_treatment_lead_1,
a.avg_amt_bin_treatment_lead_1,
a.avg_amt_bin_treatment_lead_1_win,
a.heel_type_name2,
a.style_name,
a.category_name3_2,
a.category_name3_3_Sum,
b.avg_store_last_4_week_sr,
b.normed_avg_store_mthly_sr
from belle_sh.sy_19_sku_lvl_weekly_sr_include_SH a left join belle_sh.sy_19_store_lvl_top40_avg_SR_weekly_normed_all_season b
on a.store_no=b.store_no and a.year=a.year and a.week=b.week) aa
left join 
    (select * from belle_sh.week_month_mapping
	) bb
on aa.year=bb.year and aa.week=bb.week
;

drop table belle_sh.sy_19_sku_lvl_weekly_all_info;
create table belle_sh.sy_19_sku_lvl_weekly_all_info as    
select 
*
from belle_sh.sy_19_sku_lvl_weekly_all_info_tmp
where duplicate_rank=1
;


drop table belle_sh.sy_sku_lvl_weekly_19_avg_sr_1_spr_sum_prefer_v3;
create table belle_sh.sy_sku_lvl_weekly_19_avg_sr_1_spr_sum_prefer_v3 as    
select
ccc.*,
ddd.store_lvl_mean as avg_sr_heel_type2_store_lvl_mean
from(
select
aaa.*,
bbb.avg_sr_mthly as avg_sr_heel_type2,
bbb.normed_avg_sr_mthly as normed_avg_sr_heel_type2
from(
select
cc.*,
dd.store_lvl_mean as avg_sr_categ_name3_3_sum_store_lvl_mean
from(
select 
aa.*,
bb.avg_sr_mthly as avg_sr_categ_name3_3_sum,
bb.normed_avg_sr_mthly as normed_avg_sr_categ_name3_3_sum
from 
(
select
c.*,
d.store_lvl_mean as avg_sr_categ_name3_2_store_lvl_mean
from (
select
a.*,
(case 
when a.store_no in ('I010ST','I01MST','I022ST','I023ST','I030ST','I032ST','I034ST','I038ST','I039ST','I045ST','IAJ013','I01IST','I044ST') then 'A'
when a.store_no in ('I014ST','I01BST','I01FST','I01JST','I01SST','I040ST','I046ST','I051ST','I064ST','I072ST','I078ST','IAJ001','IAJ003','IAJ004','IAJ007','IAJ009','IAJ011','IAJ002','I028ST','IAJ014') then 'B'
when a.store_no in ('I01HST','I01VST','I029ST','I033ST','I035ST','I043ST','I050ST','I065ST','I066ST','I067ST','I070ST','I074ST','IAJ006','I063ST') then 'C'
when a.store_no in ('I01UST','I01WST','I048ST','I056ST','IAJ008','IAJ012','I0ZHST') then 'D'
else 'others' end) as store_lvl,
b.avg_sr_mthly as avg_sr_categ_name3_2,
b.normed_avg_sr_mthly as normed_avg_sr_categ_name3_2
from (select * from belle_sh.sy_19_sku_lvl_weekly_all_info where product_season_name in ('春','夏')
) a left join (select * from belle_sh.sy_store_lvl_avg_SR_mthly_normed_18 where product_year_name='2018') b
on a.store_no=b.store_no and a.next_week_year=b.year+1 and a.next_week_max_month=b.month and a.product_season_name=b.product_season_name and a.category_name3_2=b.category_name3_2
) c left join (select store_lvl,year,month,product_season_name,category_name3_2,max(mean) as store_lvl_mean from belle_sh.sy_store_lvl_avg_SR_mthly_normed_18 where product_year_name='2018' group by store_lvl,year,month,product_season_name,category_name3_2) d
on c.store_lvl=d.store_lvl and c.next_week_year=d.year+1 and c.next_week_max_month=d.month and c.product_season_name=d.product_season_name and c.category_name3_2=d.category_name3_2
) aa left join (select * from belle_sh.sy_store_lvl_avg_SR_mthly_normed_sum_18 where product_year_name='2018') bb
on aa.store_no=bb.store_no and aa.next_week_year=bb.year+1 and aa.next_week_max_month=bb.month and aa.product_season_name=bb.product_season_name and aa.category_name3_3_Sum=bb.category_name3_3_Sum
) cc left join (select store_lvl,year,month,product_season_name,category_name3_3_Sum,max(mean) as store_lvl_mean from belle_sh.sy_store_lvl_avg_SR_mthly_normed_sum_18 where product_year_name='2018' group by store_lvl,year,month,product_season_name,category_name3_3_Sum) dd 
on cc.store_lvl=dd.store_lvl and cc.next_week_year=dd.year+1 and cc.next_week_max_month=dd.month and cc.product_season_name=dd.product_season_name and cc.category_name3_3_Sum=dd.category_name3_3_Sum
) aaa left join (select * from belle_sh.sy_store_lvl_heel_type_avg_SR_mthly_normed_18 where product_year_name='2018') bbb
on aaa.store_no=bbb.store_no and aaa.next_week_year=bbb.year+1 and aaa.next_week_max_month=bbb.month and aaa.product_season_name=bbb.product_season_name and aaa.heel_type_name2=bbb.heel_type_name2
) ccc left join (select store_lvl,year,month,product_season_name,heel_type_name2,max(mean) as store_lvl_mean from belle_sh.sy_store_lvl_heel_type_avg_SR_mthly_normed_18 where product_year_name='2018' group by store_lvl,year,month,product_season_name,heel_type_name2) ddd
on ccc.store_lvl=ddd.store_lvl and ccc.next_week_year=ddd.year+1 and ccc.next_week_max_month=ddd.month and ccc.product_season_name=ddd.product_season_name and ccc.heel_type_name2=ddd.heel_type_name2
;



drop table belle_sh.sy_sku_lvl_weekly_19_avg_sr_2_spr_sum_prefer_v3;
create table belle_sh.sy_sku_lvl_weekly_19_avg_sr_2_spr_sum_prefer_v3 as    
select
cccc.*,
dddd.store_lvl_mean as avg_sr_60D_store_lvl_mean
from(
select
aaaa.*,
bbbb.avg_sr_mthly as avg_sr_60D,
bbbb.normed_avg_sr_mthly as normed_avg_sr_60D
from(
select
ccc.*,
ddd.store_lvl_mean as avg_sr_30D_store_lvl_mean
from(
select
aaa.*,
bbb.avg_sr_mthly as avg_sr_30D,
bbb.normed_avg_sr_mthly as normed_avg_sr_30D
from (
select
cc.*,
dd.store_lvl_mean as avg_sr_sale_amt_bin_store_lvl_mean
from(
select
aa.*,
bb.avg_sr_mthly as avg_sr_sale_amt_bin,
bb.normed_avg_sr_mthly as normed_avg_sr_sale_amt_bin
from (
select
c.*,
d.store_lvl_mean as avg_sr_discount_bin_store_lvl_mean
from(
select
a.*,
b.avg_sr_mthly as avg_sr_discount_bin,
b.normed_avg_sr_mthly as normed_avg_sr_discount_bin
from belle_sh.sy_sku_lvl_weekly_19_avg_sr_1_spr_sum_prefer_v3 a
left join (select * from belle_sh.sy_store_lvl_discount_avg_SR_yearly_normed_spr_sum_18 where product_year_name='2018') b
on a.store_no=b.store_no and a.next_week_year=b.year+1 and a.product_season_name=b.product_season_name and a.discount_bin_treatment_lead_1=b.discount_bin
) c left join (select store_lvl,year,product_season_name,discount_bin,max(mean) as store_lvl_mean from belle_sh.sy_store_lvl_discount_avg_SR_yearly_normed_spr_sum_18  where product_year_name='2018' group by store_lvl,year,product_season_name,discount_bin) d
on c.store_lvl=d.store_lvl and c.next_week_year=d.year+1 and c.product_season_name=d.product_season_name and c.discount_bin_treatment_lead_1=d.discount_bin
) aa left join (select * from belle_sh.sy_store_lvl_sale_amt_avg_SR_mthly_normed_18 where product_year_name='2018') bb
on aa.store_no=bb.store_no and aa.next_week_year=bb.year+1 and aa.next_week_max_month=bb.month and aa.product_season_name=bb.product_season_name and aa.avg_amt_bin_treatment_lead_1=bb.avg_sale_amt_bin
) cc left join  (select store_lvl,year,month,product_season_name,avg_sale_amt_bin,max(mean) as store_lvl_mean from belle_sh.sy_store_lvl_sale_amt_avg_SR_mthly_normed_18 where product_year_name='2018' group by store_lvl,year,month,product_season_name,avg_sale_amt_bin) dd
on cc.store_lvl=dd.store_lvl and cc.next_week_year=dd.year+1 and cc.next_week_max_month=dd.month and cc.product_season_name=dd.product_season_name and cc.avg_amt_bin_treatment_lead_1=dd.avg_sale_amt_bin
) aaa left join (select * from belle_sh.sy_store_lvl_30D_avg_SR_yearly_normed_spr_sum_18 where product_year_name='2018') bbb
on aaa.store_no=bbb.store_no and aaa.next_week_year=bbb.year+1 and aaa.product_season_name=bbb.product_season_name
) ccc left join (select store_lvl,year,product_season_name,max(mean) as store_lvl_mean from belle_sh.sy_store_lvl_30D_avg_SR_yearly_normed_spr_sum_18 where product_year_name='2018' group by store_lvl,year,product_season_name) ddd
on ccc.store_lvl=ddd.store_lvl and ccc.next_week_year=ddd.year+1 and ccc.product_season_name=ddd.product_season_name
) aaaa left join (select * from belle_sh.sy_store_lvl_60D_avg_SR_yearly_normed_spr_sum_18 where product_year_name='2018') bbbb
on aaaa.store_no=bbbb.store_no and aaaa.next_week_year=bbbb.year+1 and aaaa.product_season_name=bbbb.product_season_name
) cccc left join (select store_lvl,year,product_season_name,max(mean) as store_lvl_mean from belle_sh.sy_store_lvl_60D_avg_SR_yearly_normed_spr_sum_18 where product_year_name='2018' group by store_lvl,year,product_season_name) dddd
on cccc.store_lvl=dddd.store_lvl and cccc.next_week_year=dddd.year+1 and cccc.product_season_name=dddd.product_season_name
;



drop table belle_sh.sy_sku_lvl_weekly_19_avg_sr_final_spr_sum_prefer_${hiveconf:version_tag}_${hiveconf:invsunday}_allcombine;
create table belle_sh.sy_sku_lvl_weekly_19_avg_sr_final_spr_sum_prefer_${hiveconf:version_tag}_${hiveconf:invsunday}_allcombine as    
select
(case 
when aa.use_30D_ind=1 and aa.product_season_name='夏' then ${hiveconf:wgt_sku_smr}*nvl(aa.sale_rate_last_2_week,0)+${hiveconf:wgt_store_smr}*nvl(aa.avg_store_last_4_week_sr,0)+0.06*nvl(aa.avg_sr_categ_name3_3_sum_treat,0)+0.06*nvl(aa.avg_sr_heel_type2_treat,0)+0.06*nvl(aa.avg_sr_discount_bin_treat,0)+0.06*nvl(aa.avg_sr_sale_amt_bin_treat,0)+0.06*aa.avg_sr_30d_treat
when aa.use_30D_ind=1 and aa.product_season_name='春' then ${hiveconf:wgt_sku_spr}*nvl(aa.sale_rate_last_2_week,0)+${hiveconf:wgt_store_spr}*nvl(aa.avg_store_last_4_week_sr,0)+0.06*nvl(aa.avg_sr_categ_name3_2_treat,0)    +0.06*nvl(aa.avg_sr_heel_type2_treat,0)+0.06*nvl(aa.avg_sr_discount_bin_treat,0)+0.06*nvl(aa.avg_sr_sale_amt_bin_treat,0)+0.06*aa.avg_sr_30d_treat
when aa.use_60D_ind=1 and aa.product_season_name='夏' then ${hiveconf:wgt_sku_smr}*nvl(aa.sale_rate_last_2_week,0)+${hiveconf:wgt_store_smr}*nvl(aa.avg_store_last_4_week_sr,0)+0.06*nvl(aa.avg_sr_categ_name3_3_sum_treat,0)+0.06*nvl(aa.avg_sr_heel_type2_treat,0)+0.06*nvl(aa.avg_sr_discount_bin_treat,0)+0.06*nvl(aa.avg_sr_sale_amt_bin_treat,0)+0.06*aa.avg_sr_60d_treat
when aa.use_60D_ind=1 and aa.product_season_name='春' then ${hiveconf:wgt_sku_spr}*nvl(aa.sale_rate_last_2_week,0)+${hiveconf:wgt_store_spr}*nvl(aa.avg_store_last_4_week_sr,0)+0.06*nvl(aa.avg_sr_categ_name3_2_treat,0)    +0.06*nvl(aa.avg_sr_heel_type2_treat,0)+0.06*nvl(aa.avg_sr_discount_bin_treat,0)+0.06*nvl(aa.avg_sr_sale_amt_bin_treat,0)+0.06*aa.avg_sr_60d_treat
when aa.use_30D_ind=0 and aa.use_60D_ind=0 and aa.product_season_name='夏' then ${hiveconf:wgt_sku_smr}*nvl(aa.sale_rate_last_2_week,0)+${hiveconf:wgt_store_smr}*nvl(aa.avg_store_last_4_week_sr,0)+0.06*nvl(aa.avg_sr_categ_name3_3_sum_treat,0)+0.06*nvl(aa.avg_sr_heel_type2_treat,0)+0.06*nvl(aa.avg_sr_discount_bin_treat,0)+0.06*nvl(aa.avg_sr_sale_amt_bin_treat,0)
when aa.use_30D_ind=0 and aa.use_60D_ind=0 and aa.product_season_name='春' then ${hiveconf:wgt_sku_spr}*nvl(aa.sale_rate_last_2_week,0)+${hiveconf:wgt_store_spr}*nvl(aa.avg_store_last_4_week_sr,0)+0.06*nvl(aa.avg_sr_categ_name3_2_treat,0)    +0.06*nvl(aa.avg_sr_heel_type2_treat,0)+0.06*nvl(aa.avg_sr_discount_bin_treat,0)+0.06*nvl(aa.avg_sr_sale_amt_bin_treat,0)
else 9999 end) as final_normed_avg_sr_score,
aa.*
from
(select
a.store_no,
a.store_lvl,
a.product_code,
a.heel_type_name2,
a.style_name,
a.next_week_year,
a.next_week_max_month,
a.week_first_date,
a.week_last_date,
a.category_name3_2,
a.category_name3_3_sum,
a.discount_bin_treatment_lead_1,
a.avg_amt_bin_treatment_lead_1,
a.avg_amt_bin_treatment_lead_1_win,
a.product_year_name,
a.product_season_name,
a.year,
a.week,
a.sku_historical_exist_ind,
a.sku_in_inv_ind_last_2_week,
a.sku_in_inv_ind_last_4_week,
a.qty_last_2_week,
a.qty_last_4_week,
a.normed_sr_last_2_week,
a.sale_rate_last_2_week,
a.avg_store_last_4_week_sr,
a.normed_avg_store_mthly_sr,
if(a.normed_avg_store_mthly_sr>=0,1,0) as store_sale_cpb_above_avg_ind,
if(a.store_no='I01MST',a.avg_sr_categ_name3_2_store_lvl_mean,nvl(a.avg_sr_categ_name3_2,a.avg_sr_categ_name3_2_store_lvl_mean)) as avg_sr_categ_name3_2_treat,
if(a.store_no='I01MST',a.avg_sr_categ_name3_3_sum_store_lvl_mean,nvl(a.avg_sr_categ_name3_3_sum,a.avg_sr_categ_name3_3_sum_store_lvl_mean)) as avg_sr_categ_name3_3_sum_treat,
if(a.store_no='I01MST',a.avg_sr_heel_type2_store_lvl_mean,nvl(a.avg_sr_heel_type2,a.avg_sr_heel_type2_store_lvl_mean)) as avg_sr_heel_type2_treat,
if(a.store_no='I01MST',a.avg_sr_discount_bin_store_lvl_mean,nvl(a.avg_sr_discount_bin,a.avg_sr_discount_bin_store_lvl_mean)) as avg_sr_discount_bin_treat,
if(a.store_no='I01MST',a.avg_sr_sale_amt_bin_store_lvl_mean,nvl(a.avg_sr_sale_amt_bin,a.avg_sr_sale_amt_bin_store_lvl_mean)) as avg_sr_sale_amt_bin_treat,
if(a.store_no='I01MST',a.avg_sr_30D_store_lvl_mean,nvl(a.avg_sr_30D,a.avg_sr_30D_store_lvl_mean)) as avg_sr_30D_treat,
if(a.store_no='I01MST',a.avg_sr_60D_store_lvl_mean,nvl(a.avg_sr_60D,a.avg_sr_60D_store_lvl_mean)) as avg_sr_60D_treat,
a.avg_sr_categ_name3_2,
a.avg_sr_categ_name3_3_sum,
a.avg_sr_heel_type2,
a.avg_sr_discount_bin,
a.avg_sr_sale_amt_bin,
a.avg_sr_30D,
a.avg_sr_60D,
b.first_online_30D_week,
b.first_online_60D_week,
if(a.week<=b.first_online_30D_week or b.first_online_30D_week is null,1,0) as use_30D_ind,
if(a.week<=b.first_online_60D_week and a.week>b.first_online_30D_week,1,0) as use_60D_ind,
a.discount_less_than_5,
a.price as sale_amt_lead_1
from belle_sh.sy_sku_lvl_weekly_19_avg_sr_2_spr_sum_prefer_v3 a left join belle_sh.sy_sku_base_first_online_date_19_mapping b
on a.store_no=b.store_no and a.product_code=b.product_code) aa
;

drop table belle_sh.sy_sku_lvl_weekly_19_avg_sr_final_spr_sum_prefer_${hiveconf:version_tag}_${hiveconf:invsunday}_allcombine_filter;
create table belle_sh.sy_sku_lvl_weekly_19_avg_sr_final_spr_sum_prefer_${hiveconf:version_tag}_${hiveconf:invsunday}_allcombine_filter as
select
*
from belle_sh.sy_sku_lvl_weekly_19_avg_sr_final_spr_sum_prefer_${hiveconf:version_tag}_${hiveconf:invsunday}_allcombine
where store_no <> 'I01VST' and store_no <> 'I050ST' and store_no <> 'I065ST';


drop table belle_sh.sy_sku_lvl_weekly_19_avg_sr_final_spr_sum_prefer_${hiveconf:version_tag}_before_${hiveconf:invsunday}_allcombine;
create table belle_sh.sy_sku_lvl_weekly_19_avg_sr_final_spr_sum_prefer_${hiveconf:version_tag}_before_${hiveconf:invsunday}_allcombine as    
select
*
from (
select * from belle_sh.sy_sku_lvl_weekly_19_avg_sr_final_spr_sum_prefer_${hiveconf:version_tag}_${hiveconf:invsunday}_allcombine_filter
union all
select * from belle_sh.sy_sku_lvl_weekly_19_avg_sr_final_spr_sum_prefer_${hiveconf:version_tag}_before_${hiveconf:last_invsunday}_allcombine
) test
;

`
	hql := utils.HQL{
		Hql:sql,
	}
	hql.TrimSql()
	depends_resuls, depneds_all := hql.ParseDepends()
	var obj_tables []string
	for _, item := range depends_resuls{
		obj_t, _ := item["obj_table"]
		obj_tables = append(obj_tables, obj_t.(string))
	}
	log.Println(obj_tables)
	var nodes []models.Node
	var links []models.Link
	for _, ob := range obj_tables{
		lb1 := models.Label{
			Label:ob,
		}
		nod1 := models.Node{
			Id:ob,
			Value: lb1,
		}
		nodes = append(nodes, nod1)
	}
	for _, l := range depneds_all{
		l1 := strings.Split(l, " >> ")
		u := l1[0]
		v := l1[1]
		lb1 := models.Label{
			Label:"",
		}
		lk1 := models.Link{
			U:u,
			V:v,
			Value:lb1,
		}
		links = append(links, lk1)
	}

	dag := models.DAG{
		Name:"dag",
		Nodes: nodes,
		Links:links,
	}
	d.Data["json"] = dag
	d.ServeJSON()
	return
}

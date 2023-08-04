-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.

DELETE FROM transportation_accounting_codes;
ALTER table public.transportation_accounting_codes DROP constraint transportation_accounting_codes_tac_key;

COPY public.lines_of_accounting (id, loa_sys_id, loa_dpt_id, loa_tnsfr_dpt_nm, loa_baf_id, loa_trsy_sfx_tx, loa_maj_clm_nm, loa_op_agncy_id, loa_allt_sn_id, loa_pgm_elmnt_id, loa_tsk_bdgt_sbln_tx, loa_df_agncy_alctn_rcpnt_id, loa_jb_ord_nm, loa_sbaltmt_rcpnt_id, loa_wk_cntr_rcpnt_nm, loa_maj_rmbsmt_src_id, loa_dtl_rmbsmt_src_id, loa_cust_nm, loa_obj_cls_id, loa_srv_src_id, loa_spcl_intr_id, loa_bdgt_acnt_cls_nm, loa_doc_id, loa_cls_ref_id, loa_instl_acntg_act_id, loa_lcl_instl_id, loa_fms_trnsactn_id, loa_dsc_tx, loa_bgn_dt, loa_end_dt, loa_fnct_prs_nm, loa_stat_cd, loa_hist_stat_cd, loa_hs_gds_cd, org_grp_dfas_cd, loa_uic, loa_trnsn_id, loa_sub_acnt_id, loa_bet_cd, loa_fnd_ty_fg_cd, loa_bgt_ln_itm_id, loa_scrty_coop_impl_agnc_cd, loa_scrty_coop_dsgntr_cd, loa_scrty_coop_ln_itm_id, loa_agnc_dsbr_cd, loa_agnc_acntng_cd, loa_fnd_cntr_id, loa_cst_cntr_id, loa_prj_id, loa_actvty_id, loa_cst_cd, loa_wrk_ord_id, loa_fncl_ar_id, loa_scrty_coop_cust_cd, loa_end_fy_tx, loa_bg_fy_tx, loa_bgt_rstr_cd, loa_bgt_sub_act_cd, created_at, updated_at) FROM stdin;
f24c7cea-c08d-4fc7-be18-51dcc8da55ac	10001	1	\N	1234	0000	\N	1A	123A	00000000	\N	\N	\N	\N	\N	\N	\N	\N	22NL	\N	\N	000000	HHG12345678900	\N	12345	\N	\N	PERSONAL PROPERTY - FAKE DATA DIVISION	2005-10-01	2015-10-01	\N	U	\N	HT	ZZ	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	2023-08-03 19:17:10.05037	2023-08-03 19:17:38.776652
3cbb7bb3-2d68-4d98-b0b5-b91ef7f01f96	10002	1	\N	4321	0000	\N	1A	123A	00000000	\N	\N	\N	\N	\N	\N	\N	\N	22NL	\N	\N	000000	HHG12345678900	\N	12345	\N	\N	PERSONAL PROPERTY - OBFUSCATED DATA DIVISION	2005-10-01	2015-10-01	\N	U	\N	HT	ZZ	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	2023-08-03 19:17:10.05037	2023-08-03 19:17:38.776652
06254fc3-b763-484c-b555-42855d1ad5cd	10003	1	\N	1234	0000	\N	1A	123A	00000000	\N	\N	\N	\N	\N	\N	\N	\N	22NL	\N	\N	000000	HHG12345678900	\N	12345	\N	\N	PERSONAL PROPERTY - PARANORMAL ACTIVITY DIVISION	2005-10-01	2015-10-01	\N	U	\N	HT	ZZ	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	2023-08-03 19:17:10.05037	2023-08-03 19:17:38.776652
\.

-- AAAA, BBBB, CCCC are connected to lines of accounting
-- CCCC has two entries, one of which is older than the other
-- DDDD has a loa_sys_id that does not correspond to any lines of accounting in the database
-- EEEE has null loa_sys_id
COPY public.transportation_accounting_codes (id, tac, loa_id, tac_sys_id, loa_sys_id, tac_fy_txt, tac_fn_bl_mod_cd, org_grp_dfas_cd, tac_mvt_dsg_id, tac_ty_cd, tac_use_cd, tac_maj_clmt_id, tac_bill_act_txt, tac_cost_ctr_nm, buic, tac_hist_cd, tac_stat_cd, trnsprtn_acnt_tx, trnsprtn_acnt_bgn_dt, trnsprtn_acnt_end_dt, dd_actvty_adrs_id, tac_blld_add_frst_ln_tx, tac_blld_add_scnd_ln_tx, tac_blld_add_thrd_ln_tx, tac_blld_add_frth_ln_tx, tac_fnct_poc_nm, created_at, updated_at) FROM stdin;
47170fde-4622-4058-a4ca-25e125837212	AAAA	f24c7cea-c08d-4fc7-be18-51dcc8da55ac	67891	10001	2017	W	HS	\N	O	N	12345	123456	12345	\N	\N	I	FAKE HOUSING 1	2020-10-01	2025-09-30	Z12345	COMMANDING OFFICER	FINANCE CENTER	123 ANY ST	BEVERLY HILLS CA 90210	NO POC	2023-08-03 19:17:40.126142	2023-08-03 19:17:41.411975
47170fde-4622-4058-a4ca-25e125837212	BBBB	3cbb7bb3-2d68-4d98-b0b5-b91ef7f01f96	12345	10002	2017	W	HS	\N	O	N	12345	123456	12345	\N	\N	I	FAKE HOUSING 2	2013-10-01	2020-09-30	Z12345	COMMANDING OFFICER	FINANCE CENTER	123 ANY ST	BEVERLY HILLS CA 90210	NO POC	2023-08-03 19:17:40.126142	2023-08-03 19:17:41.411975
47170fde-4622-4058-a4ca-25e125837212	CCCC	06254fc3-b763-484c-b555-42855d1ad5cd	12345	10003	2017	W	HS	\N	O	N	12345	123456	12345	\N	\N	I	FAKE HOUSING 3	2013-10-01	2020-09-30	Z12345	COMMANDING OFFICER	FINANCE CENTER	123 ANY ST	BEVERLY HILLS CA 90210	NO POC	2023-08-03 19:17:40.126142	2023-08-03 19:17:41.411975
47170fde-4622-4058-a4ca-25e125837212	CCCC	06254fc3-b763-484c-b555-42855d1ad5cd	67891	10003	2017	W	HS	\N	O	N	12345	123456	12345	\N	\N	I	FAKE HOUSING 3	2020-10-01	2025-09-30	Z12345	COMMANDING OFFICER	FINANCE CENTER	123 ANY ST	BEVERLY HILLS CA 90210	NO POC	2023-08-03 19:17:40.126142	2023-08-03 19:17:41.411975
47170fde-4622-4058-a4ca-25e125837212	DDDD	\N	\N	\N	2017	W	HS	\N	O	N	12345	123456	12345	\N	\N	I	FAKE HOUSING 4	2013-10-01	2025-09-30	Z12345	COMMANDING OFFICER	FINANCE CENTER	123 ANY ST	BEVERLY HILLS CA 90210	NO POC	2023-08-03 19:17:40.126142	2023-08-03 19:17:41.411975
47170fde-4622-4058-a4ca-25e125837212	EEEE	\N	20000	\N	2017	W	HS	\N	O	N	12345	123456	12345	\N	\N	I	FAKE HOUSING 5	2013-10-01	2025-09-30	Z12345	COMMANDING OFFICER	FINANCE CENTER	123 ANY ST	BEVERLY HILLS CA 90210	NO POC	2023-08-03 19:17:40.126142	2023-08-03 19:17:41.411975
\.

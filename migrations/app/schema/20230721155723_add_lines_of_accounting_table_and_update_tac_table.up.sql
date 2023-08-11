CREATE TABLE lines_of_accounting
(
	id uuid PRIMARY KEY,
	loa_sys_id integer,
	loa_dpt_id char(2),
	loa_tnsfr_dpt_nm varchar(4),
	loa_baf_id varchar(4),
	loa_trsy_sfx_tx varchar(4),
	loa_maj_clm_nm varchar(4),
	loa_op_agncy_id varchar(4),
	loa_allt_sn_id varchar(5),
	loa_pgm_elmnt_id varchar(12),
	loa_tsk_bdgt_sbln_tx varchar(8),
	loa_df_agncy_alctn_rcpnt_id varchar(4),
	loa_jb_ord_nm varchar(10),
	loa_sbaltmt_rcpnt_id char(1),
	loa_wk_cntr_rcpnt_nm varchar(6),
	loa_maj_rmbsmt_src_id char(1),
	loa_dtl_rmbsmt_src_id varchar(3),
	loa_cust_nm varchar(6),
	loa_obj_cls_id varchar(6),
	loa_srv_src_id char(1),
	loa_spcl_intr_id char(2),
	loa_bdgt_acnt_cls_nm varchar(8),
	loa_doc_id varchar(15),
	loa_cls_ref_id char(2),
	loa_instl_acntg_act_id varchar(6),
	loa_lcl_instl_id varchar(18),
	loa_fms_trnsactn_id varchar(12),
	loa_dsc_tx text,
	loa_bgn_dt date,
	loa_end_dt date,
	loa_fnct_prs_nm varchar(255),
	loa_stat_cd char(1),
	loa_hist_stat_cd char(1),
	loa_hs_gds_cd char(2),
	org_grp_dfas_cd char(2),
	loa_uic char(6),
	loa_trnsn_id varchar(255),
	loa_sub_acnt_id char(3),
	loa_bet_cd char(4),
	loa_fnd_ty_fg_cd char(1),
	loa_bgt_ln_itm_id varchar(8),
	loa_scrty_coop_impl_agnc_cd char(1),
	loa_scrty_coop_dsgntr_cd varchar(4),
	loa_scrty_coop_ln_itm_id varchar(3),
	loa_agnc_dsbr_cd char(6),
	loa_agnc_acntng_cd char(6),
	loa_fnd_cntr_id varchar(12),
	loa_cst_cntr_id varchar(16),
	loa_prj_id varchar(12),
	loa_actvty_id varchar(11),
	loa_cst_cd varchar(16),
	loa_wrk_ord_id varchar(16),
	loa_fncl_ar_id varchar(6),
	loa_scrty_coop_cust_cd char(2),
	loa_end_fy_tx integer,
	loa_bg_fy_tx integer,
	loa_bgt_rstr_cd char(1),
	loa_bgt_sub_act_cd char(4),
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);

-- Column Comments
COMMENT on TABLE lines_of_accounting IS 'A Line of Accounting (LOA) is the funding associated with a federal organization’s budget';
COMMENT on COLUMN lines_of_accounting.loa_sys_id IS 'Unique primary id that is referenced from rows in the Transportation Accounting spreadsheet';
COMMENT on COLUMN lines_of_accounting.loa_dpt_id IS 'Department Indicator';
COMMENT on COLUMN lines_of_accounting.loa_tnsfr_dpt_nm IS 'Transfer From Department';
COMMENT on COLUMN lines_of_accounting.loa_baf_id IS 'Basic Symbol Number';
COMMENT on COLUMN lines_of_accounting.loa_trsy_sfx_tx IS 'Subhead/Limit';
COMMENT on COLUMN lines_of_accounting.loa_maj_clm_nm IS 'Fund Code/Major Claimant';
COMMENT on COLUMN lines_of_accounting.loa_op_agncy_id IS 'Operation Agency Code/Fund Admin';
COMMENT on COLUMN lines_of_accounting.loa_allt_sn_id IS 'Allotment Serial Number';
COMMENT on COLUMN lines_of_accounting.loa_pgm_elmnt_id IS 'Program Element Code';
COMMENT on COLUMN lines_of_accounting.loa_tsk_bdgt_sbln_tx IS 'Project Task/Budget Subline';
COMMENT on COLUMN lines_of_accounting.loa_df_agncy_alctn_rcpnt_id IS 'Defense Agency Allocation Recipient';
COMMENT on COLUMN lines_of_accounting.loa_jb_ord_nm IS 'Job Order/Work Order Code';
COMMENT on COLUMN lines_of_accounting.loa_sbaltmt_rcpnt_id IS 'Sub-Allotment Recipient';
COMMENT on COLUMN lines_of_accounting.loa_wk_cntr_rcpnt_nm IS 'Work Center Recipient';
COMMENT on COLUMN lines_of_accounting.loa_maj_rmbsmt_src_id IS 'Major Reimbursement Source Code';
COMMENT on COLUMN lines_of_accounting.loa_dtl_rmbsmt_src_id IS 'Detail Reimbursement Source Code';
COMMENT on COLUMN lines_of_accounting.loa_cust_nm IS 'Customer Indicator/MPC';
COMMENT on COLUMN lines_of_accounting.loa_obj_cls_id IS 'Object Class';
COMMENT on COLUMN lines_of_accounting.loa_srv_src_id IS 'Government/Public Sector Identifier';
COMMENT on COLUMN lines_of_accounting.loa_spcl_intr_id IS 'Special Interest Code/Special Program Cost Code';
COMMENT on COLUMN lines_of_accounting.loa_bdgt_acnt_cls_nm IS 'DoD Budget & Accounting Class. Code';
COMMENT on COLUMN lines_of_accounting.loa_doc_id IS 'Document/Record Reference Number';
COMMENT on COLUMN lines_of_accounting.loa_cls_ref_id IS 'Accounting Classification Reference Number';
COMMENT on COLUMN lines_of_accounting.loa_instl_acntg_act_id IS 'Accounting Installation Number';
COMMENT on COLUMN lines_of_accounting.loa_lcl_instl_id IS 'Local Installation Data/IFS Number';
COMMENT on COLUMN lines_of_accounting.loa_fms_trnsactn_id IS 'Transaction Type';
COMMENT on COLUMN lines_of_accounting.loa_dsc_tx IS 'FMS Country Code, Implementing Agency, Case Number & Line Item Number';
COMMENT on COLUMN lines_of_accounting.loa_bgn_dt IS 'Begin Date (Fiscal Year)';
COMMENT on COLUMN lines_of_accounting.loa_end_dt IS 'End Date (Fiscal Year)';
COMMENT on COLUMN lines_of_accounting.loa_fnct_prs_nm IS 'Financial POC';
COMMENT on COLUMN lines_of_accounting.loa_stat_cd IS 'Status Code';
COMMENT on COLUMN lines_of_accounting.loa_hist_stat_cd IS 'History Status Code';
COMMENT on COLUMN lines_of_accounting.loa_hs_gds_cd IS 'Household Goods Program Code';
COMMENT on COLUMN lines_of_accounting.org_grp_dfas_cd IS 'Transportation Service Code';
COMMENT on COLUMN lines_of_accounting.loa_uic IS 'Unit Identification Code';
COMMENT on COLUMN lines_of_accounting.loa_trnsn_id IS 'Transaction ID';
COMMENT on COLUMN lines_of_accounting.loa_sub_acnt_id IS 'Sub Account';
COMMENT on COLUMN lines_of_accounting.loa_bet_cd IS 'Business Event Type Code';
COMMENT on COLUMN lines_of_accounting.loa_fnd_ty_fg_cd IS 'Reimbursable Flag';
COMMENT on COLUMN lines_of_accounting.loa_bgt_ln_itm_id IS 'Budget Line Item';
COMMENT on COLUMN lines_of_accounting.loa_scrty_coop_impl_agnc_cd IS 'Security Cooperation Implementing Agency Code';
COMMENT on COLUMN lines_of_accounting.loa_scrty_coop_dsgntr_cd IS 'Security Cooperation Case Designator';
COMMENT on COLUMN lines_of_accounting.loa_scrty_coop_ln_itm_id IS 'Security Cooperation Case Line Item Identifier';
COMMENT on COLUMN lines_of_accounting.loa_agnc_dsbr_cd IS 'Agency Disbursing Identifier Code';
COMMENT on COLUMN lines_of_accounting.loa_agnc_acntng_cd IS 'Agency Accounting Identifier';
COMMENT on COLUMN lines_of_accounting.loa_fnd_cntr_id IS 'Funding Center';
COMMENT on COLUMN lines_of_accounting.loa_cst_cntr_id IS 'Cost Center';
COMMENT on COLUMN lines_of_accounting.loa_prj_id IS 'Project Identifier';
COMMENT on COLUMN lines_of_accounting.loa_actvty_id IS 'Activity Identifier';
COMMENT on COLUMN lines_of_accounting.loa_cst_cd IS 'Cost Element Code';
COMMENT on COLUMN lines_of_accounting.loa_wrk_ord_id IS 'Work Order Number';
COMMENT on COLUMN lines_of_accounting.loa_fncl_ar_id IS 'Functional Area';
COMMENT on COLUMN lines_of_accounting.loa_scrty_coop_cust_cd IS 'Security Cooperation Customer Code';
COMMENT on COLUMN lines_of_accounting.loa_end_fy_tx IS 'End Fiscal Year';
COMMENT on COLUMN lines_of_accounting.loa_bg_fy_tx IS 'Begin Fiscal Year';
COMMENT on COLUMN lines_of_accounting.loa_bgt_rstr_cd IS 'Availability Type';
COMMENT on COLUMN lines_of_accounting.loa_bgt_sub_act_cd IS 'Sub-Allocation';

ALTER TABLE transportation_accounting_codes
ADD loa_id uuid,
ADD tac_sys_id integer,
ADD loa_sys_id integer,
ADD tac_fy_txt integer,
ADD tac_fn_bl_mod_cd char(1),
ADD org_grp_dfas_cd char(2),
ADD tac_mvt_dsg_id varchar(255),
ADD tac_ty_cd char(1),
ADD tac_use_cd varchar(2),
ADD tac_maj_clmt_id varchar(6),
ADD tac_bill_act_txt varchar(6),
ADD tac_cost_ctr_nm varchar(6),
ADD buic char(6),
ADD tac_hist_cd char(1),
ADD tac_stat_cd char(1),
ADD trnsprtn_acnt_tx text,
ADD trnsprtn_acnt_bgn_dt date,
ADD trnsprtn_acnt_end_dt date,
ADD dd_actvty_adrs_id char(6),
ADD tac_blld_add_frst_ln_tx varchar(255),
ADD tac_blld_add_scnd_ln_tx varchar(255),
ADD tac_blld_add_thrd_ln_tx varchar(255),
ADD tac_blld_add_frth_ln_tx varchar(255),
ADD tac_fnct_poc_nm varchar(255);

ALTER TABLE transportation_accounting_codes
	ADD CONSTRAINT transportation_accounting_codes_loa_id_fkey FOREIGN KEY (loa_id) REFERENCES lines_of_accounting (id);

CREATE INDEX IF NOT EXISTS transportation_accounting_codes_loa_id_idx ON transportation_accounting_codes(loa_id);

-- Column Comments
COMMENT on COLUMN transportation_accounting_codes.loa_id IS 'Associates the TAC to a Line of Accounting';
COMMENT on COLUMN transportation_accounting_codes.tac_sys_id IS 'TAC System Identifier';
COMMENT on COLUMN transportation_accounting_codes.loa_sys_id IS 'LOA System Identifier';
COMMENT on COLUMN transportation_accounting_codes.tac_fy_txt IS 'Fiscal year';
COMMENT on COLUMN transportation_accounting_codes.tac_fn_bl_mod_cd IS 'Financial Bill Mode Code';
COMMENT on COLUMN transportation_accounting_codes.org_grp_dfas_cd IS 'Transportation Service Code';
COMMENT on COLUMN transportation_accounting_codes.tac_mvt_dsg_id IS 'Movement Designator Code';
COMMENT on COLUMN transportation_accounting_codes.tac_ty_cd IS 'Type Code';
COMMENT on COLUMN transportation_accounting_codes.tac_use_cd IS 'Usage Code';
COMMENT on COLUMN transportation_accounting_codes.tac_maj_clmt_id IS 'Major Claimant';
COMMENT on COLUMN transportation_accounting_codes.tac_bill_act_txt IS 'Bill Account Code';
COMMENT on COLUMN transportation_accounting_codes.tac_cost_ctr_nm IS 'Cost Center';
COMMENT on COLUMN transportation_accounting_codes.buic IS 'BUIC';
COMMENT on COLUMN transportation_accounting_codes.tac_hist_cd IS 'History Status Code';
COMMENT on COLUMN transportation_accounting_codes.tac_stat_cd IS 'TAC Status Code';
COMMENT on COLUMN transportation_accounting_codes.trnsprtn_acnt_tx IS 'Description';
COMMENT on COLUMN transportation_accounting_codes.trnsprtn_acnt_bgn_dt IS 'Effective Begin Date';
COMMENT on COLUMN transportation_accounting_codes.trnsprtn_acnt_end_dt IS 'Effective End Date';
COMMENT on COLUMN transportation_accounting_codes.dd_actvty_adrs_id IS 'DODAAC';
COMMENT on COLUMN transportation_accounting_codes.tac_blld_add_frst_ln_tx IS 'Billed Address Line 1';
COMMENT on COLUMN transportation_accounting_codes.tac_blld_add_scnd_ln_tx IS 'Billed Address Line 2';
COMMENT on COLUMN transportation_accounting_codes.tac_blld_add_thrd_ln_tx IS 'Billed Address Line 3';
COMMENT on COLUMN transportation_accounting_codes.tac_blld_add_frth_ln_tx IS 'Billed Address Line 4';
COMMENT on COLUMN transportation_accounting_codes.tac_fnct_poc_nm IS 'TAC Functional POC';

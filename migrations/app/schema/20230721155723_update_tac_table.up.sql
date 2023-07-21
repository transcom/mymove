ALTER TABLE transportation_accounting_codes
ADD loa_id uuid NOT NULL,
ADD tac_sys_id integer,
ADD loa_sys_id integer,
ADD tac_fy_txt integer NOT NULL,
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

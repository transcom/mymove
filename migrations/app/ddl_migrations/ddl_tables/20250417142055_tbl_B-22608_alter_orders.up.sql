--added by Landan Parker on April 17th 2025
--alter orders to include pay_grade_rank_id column with fk constraint to pay_grade_ranks

ALTER TABLE orders
   ADD if not exists pay_grade_rank_id uuid DEFAULT NULL
   	CONSTRAINT fk_orders_pay_grade_rank_id REFERENCES pay_grade_ranks (id);

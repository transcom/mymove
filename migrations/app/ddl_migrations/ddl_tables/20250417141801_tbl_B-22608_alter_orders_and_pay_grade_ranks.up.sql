--added by Landan Parker on April 17th 2025
--alter pay_grade_ranks and orders to elegantly adjust constraints for re-running migrations

alter table orders drop constraint if exists fk_orders_pay_grade_rank_id;
alter table orders drop if exists pay_grade_rank_id;
drop table if exists pay_grade_ranks cascade;

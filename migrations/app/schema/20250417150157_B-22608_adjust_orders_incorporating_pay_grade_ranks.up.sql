--added by Landan Parker on April 17th 2025
--add pay_grade_rank_id to orders table
--update grade and pay_grade_rank on orders

update orders
   set grade = 'O_1',
       pay_grade_rank_id = '2cf8e36a-20fb-41fe-9268-d3d1f0219d1a'
 where grade in ('ACADEMY_CADET','O_1_ACADEMY_GRADUATE')
   and service_member_id in (select id from service_members where
affiliation = 'AIR_FORCE');

update orders
   set grade = 'O_1',
       pay_grade_rank_id = 'd447b93a-d0ae-4943-af1c-39830f5e7278'
 where grade in ('ACADEMY_CADET','O_1_ACADEMY_GRADUATE')
   and service_member_id in (select id from service_members where
affiliation = 'ARMY');

update orders
   set grade = 'O_1',
       pay_grade_rank_id = '0dc31054-0939-44ff-80c4-114b80f40895'
 where grade = 'MIDSHIPMAN'
   and service_member_id in (select id from service_members where
affiliation = 'NAVY');

--update pay_grade_rank_id in orders where grade:rank is 1:1
do
'declare i record; v_count int;
begin

		for i in (
			select a.id pay_grade_id, a.grade, c.id orders_id,
d.affiliation
				from pay_grades a, orders c,
service_members d
				where a.grade = c.grade
				  and c.service_member_id = d.id)
		loop

			select count(*) into v_count
			  from pay_grade_ranks
			 where pay_grade_id = i.pay_grade_id
			   and affiliation = i.affiliation;

			if v_count <= 1 then

				update orders o
				   set pay_grade_rank_id = p.id
				  from pay_grade_ranks p
				 where o.id = i.orders_id
				   and p.pay_grade_id = i.pay_grade_id
				   and p.affiliation = i.affiliation
				   and o.pay_grade_rank_id is null;

			end if;

		end loop;

	end'
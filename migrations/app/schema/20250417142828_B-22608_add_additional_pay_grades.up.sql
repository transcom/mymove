--added by Landan Parker on April 17th 2025
--insert missing pay_grades

INSERT INTO public.pay_grades
(id, grade, grade_description, created_at, updated_at)
VALUES('9a892c59-48d5-4eba-b5f9-193716da8827', 'O_1', 'Officer Grade O_1', now(), now())
	 on conflict (id) do nothing;

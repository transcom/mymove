-- B-23301  Jim Hawks  convert existing data to first and last name columns

select distinct
	trim(name) as name,
	case
		when position(' ' in name) > 0 then substr(name, 1, position(' ' in name) - 1)
		else name
	end as first_name,
	case
		when position(' ' in name) > 0 then substr(name, position(' ' in name) + 1, 255)
		else ''
	end as last_name
from backup_contacts
order by trim(name);

update backup_contacts
	set first_name =
			case
				when position(' ' in name) > 0 then substr(name, 1, position(' ' in name) - 1)
				else name
			end,
		last_name =
			case
				when position(' ' in name) > 0 then substr(name, position(' ' in name) + 1, 255)
				else ''
			end;

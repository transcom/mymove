--cleanup application_parameters
alter table application_parameters
drop column if exists validation_code;


Insert into application_parameters (id,  created_at, updated_at, parameter_name, parameter_value)
values ('c70ff1b6-c71e-4572-8f00-6c50d67ea07e',  now(), now(), 'src_email', 'your_email@example.com');

Insert into application_parameters (id,  created_at, updated_at, parameter_name, parameter_value)
values ('d4252473-00ff-4d2e-9842-c2016310a7e1',  now(), now(), 'transcom_distro_email', 'your_email@example.com');

Insert into application_parameters (id,  created_at, updated_at, parameter_name, parameter_value)
values ('b600af46-c295-4d85-aacd-65ccf065e941',  now(), now(), 'milmove_ops_email', 'your_email@example.com');
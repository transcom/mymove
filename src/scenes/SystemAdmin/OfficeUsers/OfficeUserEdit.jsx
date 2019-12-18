import React from 'react';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';
import {
  Edit,
  SimpleForm,
  TextInput,
  DisabledInput,
  SelectInput,
  required,
  Toolbar,
  SaveButton,
  CheckboxGroupInput,
} from 'react-admin';

const OfficeUserEditToolbar = props => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const OfficeUserEdit = props => (
  <Edit {...props}>
    <SimpleForm toolbar={<OfficeUserEditToolbar />}>
      <DisabledInput source="id" />
      <DisabledInput source="email" />
      <TextInput source="first_name" validate={required()} />
      <TextInput source="middle_initials" />
      <TextInput source="last_name" validate={required()} />
      <TextInput source="telephone" validate={phoneValidators} />
      <SelectInput
        source="active"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <CheckboxGroupInput
        source="roles"
        choices={[
          { roleType: 'customer', name: 'Customer' },
          { roleType: 'transportation_ordering_officer', name: 'Transportation Ordering Officer' },
          { roleType: 'transportation_invoicing_officer', name: 'Transportation Invoicing Officer' },
          { roleType: 'contracting_officer', name: 'Contracting Officer' },
          { roleType: 'ppm_office_users', name: 'PPM Office Users' },
        ]}
        optionValue="roleType"
      />
      <DisabledInput source="created_at" />
      <DisabledInput source="updated_at" />
    </SimpleForm>
  </Edit>
);

export default OfficeUserEdit;

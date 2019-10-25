import React from 'react';
import { Show, SimpleShowLayout, TextField, BooleanField, DateField, ReferenceField } from 'react-admin';

const OfficeUserShowTitle = ({ record }) => {
  return <span>{`${record.first_name} ${record.last_name}`}</span>;
};

const OfficeUserShow = props => {
  return (
    <Show {...props} title={<OfficeUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="email" />
        <TextField source="first_name" />
        <TextField source="middle_initials" />
        <TextField source="last_name" />
        <TextField source="telephone" />
        <BooleanField source="deactivated" />
        <ReferenceField label="Transportation Office" source="transportation_office_id" reference="offices">
          <TextField source="name" />
        </ReferenceField>
        <DateField source="created_at" showTime />
        <DateField source="updated_at" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default OfficeUserShow;

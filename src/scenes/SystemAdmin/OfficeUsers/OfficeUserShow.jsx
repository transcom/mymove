import React from 'react';
import {
  Show,
  SimpleShowLayout,
  ArrayField,
  Datagrid,
  TextField,
  BooleanField,
  DateField,
  ReferenceField,
} from 'react-admin';

const OfficeUserShowTitle = ({ record }) => {
  return <span>{`${record.firstName} ${record.lastName}`}</span>;
};

const OfficeUserShow = props => {
  return (
    <Show {...props} title={<OfficeUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="middleInitials" />
        <TextField source="lastName" />
        <TextField source="telephone" />
        <BooleanField source="active" />
        <ArrayField source="roles">
          <Datagrid>
            <TextField source="roleName" />
          </Datagrid>
        </ArrayField>
        <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" sortBy="name">
          <TextField component="pre" source="name" />
        </ReferenceField>
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default OfficeUserShow;

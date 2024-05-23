import React from 'react';
import {
  ArrayField,
  BooleanField,
  Datagrid,
  DateField,
  ReferenceField,
  Show,
  SimpleShowLayout,
  TextField,
  useRecordContext,
} from 'react-admin';

const OfficeUserShowTitle = () => {
  const record = useRecordContext();

  return <span>{`${record?.firstName} ${record?.lastName}`}</span>;
};

const OfficeUserShowRoles = () => {
  const record = useRecordContext();
  if (!record?.roles) return <p>No roles assigned to this office user.</p>;

  return (
    <ArrayField source="roles">
      <Datagrid bulkActionButtons={false}>
        <TextField source="roleName" />
      </Datagrid>
    </ArrayField>
  );
};

const OfficeUserShowPrivileges = () => {
  const record = useRecordContext();
  if (!record?.privileges) return <p>No privileges assigned to this office user.</p>;

  return (
    <ArrayField source="privileges">
      <Datagrid bulkActionButtons={false}>
        <TextField source="privilegeName" />
      </Datagrid>
    </ArrayField>
  );
};

const OfficeUserShow = () => {
  return (
    <Show title={<OfficeUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="userId" label="User Id" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="middleInitials" />
        <TextField source="lastName" />
        <TextField source="telephone" />
        <BooleanField source="active" />
        <OfficeUserShowRoles />
        <OfficeUserShowPrivileges />
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

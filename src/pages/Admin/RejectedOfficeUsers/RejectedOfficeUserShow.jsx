import React from 'react';
import {
  ArrayField,
  Datagrid,
  DateField,
  ReferenceField,
  Show,
  SimpleShowLayout,
  TextField,
  useRecordContext,
} from 'react-admin';

const RejectedOfficeUserShowTitle = () => {
  const record = useRecordContext();

  return <span>{`${record?.firstName} ${record?.lastName}`}</span>;
};

const RejectedOfficeUserShowRoles = () => {
  const record = useRecordContext();
  if (!record?.roles) return <p>This user has not rejected any roles.</p>;

  return (
    <ArrayField source="roles">
      <span>Rejected roles:</span>
      <Datagrid bulkActionButtons={false}>
        <TextField source="roleName" />
      </Datagrid>
    </ArrayField>
  );
};

const RejectedOfficeUserShow = () => {
  return (
    <Show title={<RejectedOfficeUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="userId" label="User Id" />
        <TextField source="status" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="middleInitials" />
        <TextField source="lastName" />
        <TextField source="telephone" />
        <TextField source="edipi" label="DODID#" />
        <TextField source="otherUniqueId" label="Other unique Id" />
        <RejectedOfficeUserShowRoles />
        <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" sortBy="name">
          <TextField component="pre" source="name" />
        </ReferenceField>
        <DateField label="Account rejected at" source="createdAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default RejectedOfficeUserShow;

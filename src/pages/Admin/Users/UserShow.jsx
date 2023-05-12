import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';

const UserShowTitle = () => {
  const record = useRecordContext();
  return <span>{`${record?.loginGovEmail}`}</span>;
};

const UserShow = () => {
  return (
    // eslint-disable-next-line react/jsx-props-no-spreading
    <Show title={<UserShowTitle />} data-testid="user-show-detail">
      <SimpleShowLayout>
        <TextField source="id" label="User ID" />
        <TextField source="loginGovEmail" label="User email" />
        <BooleanField source="active" />
        <TextField source="currentAdminSessionId" label="User current admin session ID" />
        <TextField source="currentOfficeSessionId" label="User current office session ID" />
        <TextField source="currentMilSessionId" label="User current mil session ID" />
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default UserShow;

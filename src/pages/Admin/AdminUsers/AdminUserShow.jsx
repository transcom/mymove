import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';
import { connect } from 'react-redux';

import { selectAdminUser } from 'store/entities/selectors';

const AdminUserShowTitle = () => {
  const record = useRecordContext();
  return <span>{`${record?.firstName} ${record?.lastName}`}</span>;
};

const AdminUserShow = ({ adminUser }) => {
  return (
    <Show title={<AdminUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="userId" label="User Id" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="lastName" />
        <TextField source="organizationId" label="Organization Id" />
        <BooleanField source="active" label="Active" />
        {adminUser?.super && <BooleanField source="super" label="Super Admin" />}
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

function mapStateToProps(state) {
  return {
    adminUser: selectAdminUser(state),
  };
}

export default connect(mapStateToProps)(AdminUserShow);

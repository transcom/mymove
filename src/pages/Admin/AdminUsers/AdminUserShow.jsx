import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';
import PropTypes from 'prop-types';

const AdminUserShowTitle = () => {
  const record = useRecordContext();
  return <span>{`${record.firstName} ${record.lastName}`}</span>;
};

AdminUserShowTitle.propTypes = {
  record: PropTypes.shape({
    firstName: PropTypes.string,
    lastName: PropTypes.string,
  }),
};

AdminUserShowTitle.defaultProps = {
  record: {
    firstName: '',
    lastName: '',
  },
};

const AdminUserShow = (props) => {
  return (
    <Show {...props} title={<AdminUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="userId" label="User Id" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="lastName" />
        <TextField source="organizationId" label="Organization Id" />
        <BooleanField source="active" addLabel label="Active" />
        <DateField source="createdAt" showTime addLabel />
        <DateField source="updatedAt" showTime addLabel />
      </SimpleShowLayout>
    </Show>
  );
};

export default AdminUserShow;

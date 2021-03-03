/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField } from 'react-admin';
import PropTypes from 'prop-types';

const AdminUserShowTitle = ({ record }) => {
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
        <TextField source="organizationId" />
        <BooleanField source="active" />
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default AdminUserShow;

import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField } from 'react-admin';
import PropTypes from 'prop-types';

const UserShowTitle = ({ record }) => {
  return <span>{`${record.loginGovEmail}`}</span>;
};

UserShowTitle.propTypes = {
  record: PropTypes.shape({
    loginGovEmail: PropTypes.string,
  }),
};

UserShowTitle.defaultProps = {
  record: {
    loginGovEmail: '',
  },
};

const UserShow = (props) => {
  return (
    // eslint-disable-next-line react/jsx-props-no-spreading
    <Show {...props} title={<UserShowTitle />} data-testid="user-show-detail">
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

import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField } from 'react-admin';
import PropTypes from 'prop-types';

const UserShowTitle = () => {
  const record = useRecordContext();
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

const UserShow = () => {
  return (
    // eslint-disable-next-line react/jsx-props-no-spreading
    <Show title={<UserShowTitle />} data-testid="user-show-detail">
      <SimpleShowLayout>
        <TextField source="id" label="User ID" />
        <TextField source="loginGovEmail" label="User email" />
        <BooleanField source="active" addLabel />
        <TextField source="currentAdminSessionId" label="User current admin session ID" />
        <TextField source="currentOfficeSessionId" label="User current office session ID" />
        <TextField source="currentMilSessionId" label="User current mil session ID" />
        <DateField source="createdAt" showTime addLabel />
        <DateField source="updatedAt" showTime addLabel />
      </SimpleShowLayout>
    </Show>
  );
};

export default UserShow;

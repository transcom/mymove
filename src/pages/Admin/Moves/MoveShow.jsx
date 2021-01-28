import React from 'react';
import { Show, SimpleShowLayout, TextField, DateField } from 'react-admin';
import PropTypes from 'prop-types';

const MoveShowTitle = ({ record }) => {
  return <span>{`${record.serviceMember.firstName} ${record.serviceMember.lastName}`}</span>;
};

MoveShowTitle.propTypes = {
  record: PropTypes.shape({
    serviceMember: PropTypes.shape({
      firstName: PropTypes.string,
      lastName: PropTypes.node,
    }),
  }),
};

MoveShowTitle.defaultProps = {
  record: {
    serviceMember: {
      firstName: '',
      lastName: '',
    },
  },
};

const MoveShow = (props) => {
  return (
    /* eslint-disable-next-line react/jsx-props-no-spreading */
    <Show {...props} title={<MoveShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="locator" />
        <TextField source="status" />
        <TextField source="show" />
        <TextField source="serviceMember.userId" label="User Id" />
        <TextField source="serviceMember.id" label="Service member Id" />
        <TextField source="serviceMember.firstName" label="Service member first name" />
        <TextField source="serviceMember.middleName" label="Service member middle name" />
        <TextField source="serviceMember.lastName" label="Service member last name" />
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default MoveShow;

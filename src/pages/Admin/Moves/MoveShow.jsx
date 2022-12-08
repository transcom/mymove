import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';
import PropTypes from 'prop-types';

const MoveShowTitle = () => {
  const record = useRecordContext();
  return <span>{`Move ID: ${record.id}`}</span>;
};

MoveShowTitle.propTypes = {
  record: PropTypes.shape({
    id: PropTypes.string,
  }),
};

MoveShowTitle.defaultProps = {
  record: {
    id: '',
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
        <BooleanField source="show" />
        <TextField source="ordersId" reference="moves" label="Order Id" />
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

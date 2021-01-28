import React from 'react';
import { Show, SimpleShowLayout, TextField, DateField } from 'react-admin';
import PropTypes from 'prop-types';

const MoveShowTitle = ({ serviceMember }) => {
  return <span>{`${serviceMember.firstName} ${serviceMember.lastName}`}</span>;
};

MoveShowTitle.propTypes = {
  serviceMember: PropTypes.shape({
    firstName: PropTypes.string,
    lastName: PropTypes.node,
  }).isRequired,
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
        <TextField source="id" label="Service Member Id" reference="moves.serviceMember" />
        <TextField source="firstname" />
        {/* <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" sortBy="name">
          <TextField component="pre" source="name" />
        </ReferenceField> */}
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default MoveShow;

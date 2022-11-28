import React from 'react';
import {
  ArrayField,
  BooleanField,
  Datagrid,
  DateField,
  ReferenceField,
  Show,
  SimpleShowLayout,
  TextField,
} from 'react-admin';
import PropTypes from 'prop-types';

const OfficeUserShowTitle = ({ record }) => {
  return <span>{`${record.firstName} ${record.lastName}`}</span>;
};

OfficeUserShowTitle.propTypes = {
  record: PropTypes.shape({
    firstName: PropTypes.string,
    lastName: PropTypes.string,
  }),
};

OfficeUserShowTitle.defaultProps = {
  record: {
    firstName: '',
    lastName: '',
  },
};

const OfficeUserShow = () => {
  return (
    <Show title={<OfficeUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="userId" label="User Id" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="middleInitials" />
        <TextField source="lastName" />
        <TextField source="telephone" />
        <BooleanField source="active" addLabel />
        <ArrayField source="roles" addLabel>
          <Datagrid>
            <TextField source="roleName" />
          </Datagrid>
        </ArrayField>
        <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" sortBy="name">
          <TextField component="pre" source="name" />
        </ReferenceField>
        <DateField source="createdAt" showTime addLabel />
        <DateField source="updatedAt" showTime addLabel />
      </SimpleShowLayout>
    </Show>
  );
};

export default OfficeUserShow;

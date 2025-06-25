import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField } from 'react-admin';
import PropTypes from 'prop-types';

const ClientCertShowTitle = ({ record }) => {
  return <span>{`${record.subject}`}</span>;
};

ClientCertShowTitle.propTypes = {
  record: PropTypes.shape({
    subject: PropTypes.string,
  }),
};

ClientCertShowTitle.defaultProps = {
  record: {
    subject: '',
  },
};

const ClientCertShow = (props) => {
  return (
    <Show {...props} title={<ClientCertShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="subject" />
        <TextField source="sha256Digest" />
        <TextField source="userId" label="User Id" />
        <BooleanField source="allowAirForceOrdersRead" />
        <BooleanField source="allowAirForceOrdersWrite" />
        <BooleanField source="allowArmyOrdersRead" />
        <BooleanField source="allowArmyOrdersWrite" />
        <BooleanField source="allowCoastGuardOrdersRead" />
        <BooleanField source="allowCoastGuardOrdersWrite" />
        <BooleanField source="allowMarineCorpsOrdersRead" />
        <BooleanField source="allowMarineCorpsOrdersWrite" />
        <BooleanField source="allowNavyOrdersRead" />
        <BooleanField source="allowNavyOrdersWrite" />
        <BooleanField source="allowPrime" />
        <BooleanField source="allowPPTAS" label="Allow PPTAS" />
        <DateField source="createdAt" showTime addLabel />
        <DateField source="updatedAt" showTime addLabel />
      </SimpleShowLayout>
    </Show>
  );
};

export default ClientCertShow;

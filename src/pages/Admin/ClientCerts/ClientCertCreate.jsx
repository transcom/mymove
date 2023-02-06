import React from 'react';
import { Create, SimpleForm, TextInput, SelectInput, required } from 'react-admin';
import { Typography } from '@material-ui/core';

const ClientCertCreate = (props) => (
  <Create {...props}>
    <SimpleForm>
      <Typography variant="h5" gutterBottom>
        Indentity
      </Typography>
      <Typography>
        This section is used create the client certificate and user relationship needed to authenticate via mutual TLS.
      </Typography>
      <TextInput source="subject" validate={required()} multiline />
      <TextInput source="sha256Digest" validate={required()} multiline />
      <TextInput source="user_id" validate={required()} multiline />
      {/* <ReferenceInput source="user_id" validate={required()} reference="users" /> */}
      <Typography variant="h5" gutterBottom>
        Roles
      </Typography>
      <Typography> This section is used to grant roles to the client certificate. </Typography>
      <SelectInput
        source="allowPrime"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowOrdersAPI"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowAirForceOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowAirForceOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowArmyOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowArmyOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowCoastGuardOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowCoastGuardOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowMarineCorpsOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowMarineCorpsOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowNavyOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowNavyOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
    </SimpleForm>
  </Create>
);

export default ClientCertCreate;

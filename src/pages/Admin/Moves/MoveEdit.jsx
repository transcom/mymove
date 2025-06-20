/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Edit, SelectInput, SimpleForm, TextInput } from 'react-admin';

import SaveToolbar from '../Shared/SaveToolbar';

const MoveEdit = () => (
  <Edit>
    <SimpleForm
      toolbar={<SaveToolbar />}
      sx={{ '& .MuiInputBase-input': { width: 232 } }}
      mode="onBlur"
      reValidateMode="onBlur"
    >
      <TextInput source="id" disabled />
      <TextInput source="locator" disabled />
      <TextInput source="status" disabled />
      <SelectInput
        source="show"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
        sx={{ width: 256 }}
      />
      <TextInput source="ordersId" reference="moves" label="Order Id" disabled />
      <TextInput source="serviceMember.userId" label="User Id" disabled />
      <TextInput source="serviceMember.id" label="Service member Id" disabled />
      <TextInput source="serviceMember.firstName" label="Service member first name" disabled />
      <TextInput source="serviceMember.middleName" label="Service member middle name" disabled />
      <TextInput source="serviceMember.lastName" label="Service member last name" disabled />
      <TextInput source="createdAt" disabled />
      <TextInput source="updatedAt" disabled />
      <TextInput source="availableToPrimeAt" label="Available to Prime at" disabled />
    </SimpleForm>
  </Edit>
);

export default MoveEdit;

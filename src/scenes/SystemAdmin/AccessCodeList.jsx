import React from 'react';
import { List, Datagrid, TextField, TextInput } from 'react-admin';
import { Filter, ReferenceInput, SelectInput } from 'react-admin';
import AdminPagination from './AdminPagination';

const AccessCodeFilter = props => (
  <Filter {...props}>
    <TextInput label="Search" source="code" reference="access_codes" alwaysOn />
    <ReferenceInput label="Move Type" source="move_type" reference="access_codes" allowEmpty>
      <SelectInput optionText="move_type" optionValue="move_type" />
    </ReferenceInput>
  </Filter>
);

const AccessCodeList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={500} filters={<AccessCodeFilter />}>
    <Datagrid>
      <TextField source="id" reference="access_codes" />
      <TextField source="code" reference="access_codes" />
      <TextField source="move_type" reference="access_codes" />
      <TextField source="locator" reference="access_codes" />
    </Datagrid>
  </List>
);

export default AccessCodeList;

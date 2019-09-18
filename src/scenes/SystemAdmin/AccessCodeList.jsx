import React from 'react';
import { Datagrid, Filter, List, SelectInput, TextField, TextInput } from 'react-admin';
import AdminPagination from './AdminPagination';
import styles from './AccessCode.module.scss';

const AccessCodeFilter = props => (
  <Filter {...props} className={styles['access-codes-filters']}>
    <TextInput label="Access Code" source="code" reference="access_codes" alwaysOn />
    <SelectInput source="move_type" choices={[{ id: 'PPM', name: 'PPM' }, { id: 'HHG', name: 'HHG' }]} />
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

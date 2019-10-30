import React from 'react';
import { Datagrid, Filter, List, SelectInput, TextField, TextInput } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import styles from 'scenes/SystemAdmin/Home.module.scss';

const defaultSort = { field: 'code', order: 'DESC' };

const AccessCodeFilter = props => (
  <Filter {...props} className={styles['system-admin-filters']}>
    <TextInput
      label="Access Code (don't include prefix)"
      source="code"
      reference="access_codes"
      fullWidth
      alwaysOn
      resettable
    />
    <SelectInput source="move_type" choices={[{ id: 'PPM', name: 'PPM' }, { id: 'HHG', name: 'HHG' }]} />
  </Filter>
);

const AccessCodeList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={25} filters={<AccessCodeFilter />} sort={defaultSort}>
    <Datagrid>
      <TextField source="id" reference="access_codes" />
      <TextField source="code" reference="access_codes" />
      <TextField source="move_type" reference="access_codes" />
      <TextField source="locator" reference="access_codes" />
    </Datagrid>
  </List>
);

export default AccessCodeList;

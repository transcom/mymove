import React from 'react';
import { Datagrid, Filter, List, TextField, TextInput, DateField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import styles from 'scenes/SystemAdmin/Home.module.scss';

const defaultSort = { field: 'locator', order: 'ASC' };

const MoveFilter = props => (
  <Filter {...props} className={styles['system-admin-filters']}>
    <TextInput label="Locator" source="locator" reference="locator" alwaysOn resettable />
  </Filter>
);

const AccessCodeList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={25} filters={<MoveFilter />} sort={defaultSort}>
    <Datagrid>
      <TextField source="id" reference="moves" />
      <TextField source="orders_id" reference="moves" />
      <TextField source="service_member_id" reference="moves" />
      <TextField source="locator" reference="moves" />
      <TextField source="status" reference="moves" />
      <TextField source="show" reference="moves" />
      <DateField source="created_at" reference="moves" showTime />
      <DateField source="updated_at" reference="moves" showTime />
    </Datagrid>
  </List>
);

export default AccessCodeList;

import React from 'react';
import { Datagrid, Filter, List, TextField, TextInput, DateField } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import styles from 'scenes/SystemAdmin/Home.module.scss';

const defaultSort = { field: 'locator', order: 'ASC' };

const MoveFilter = (props) => (
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  <Filter {...props} className={styles['system-admin-filters']}>
    <TextInput label="Locator" source="locator" reference="locator" alwaysOn resettable />
  </Filter>
);

const AccessCodeList = (props) => (
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  <List {...props} pagination={<AdminPagination />} perPage={25} filters={<MoveFilter />} sort={defaultSort}>
    <Datagrid rowClick="show">
      <TextField source="id" reference="moves" />
      <TextField source="ordersId" reference="moves" />
      <TextField source="id" label="Service Member Id" reference="moves.serviceMember" sortable={false} />
      <TextField source="locator" reference="moves" />
      <TextField source="status" reference="moves" />
      <TextField source="show" reference="moves" />
      <DateField source="createdAt" reference="moves" showTime />
      <DateField source="updatedAt" reference="moves" showTime />
    </Datagrid>
  </List>
);

export default AccessCodeList;

import React from 'react';
import { List, Datagrid, TextField } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import TitleizedField from 'scenes/SystemAdmin/shared/TitleizedField';

const defaultSort = { field: 'orders_number', order: 'DESC' };

const ElectronicOrderList = (props) => (
  <List {...props} pagination={<AdminPagination />} sort={defaultSort} perPage={25}>
    <Datagrid>
      <TextField source="id" />
      <TitleizedField source="issuer" />
      <TextField source="ordersNumber" />
      <TextField source="createdAt" />
    </Datagrid>
  </List>
);

export default ElectronicOrderList;

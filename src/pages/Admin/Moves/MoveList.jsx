/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Datagrid, DateField, Filter, List, TextField, TextInput } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'locator', order: 'ASC' };

const MoveFilter = (props) => (
  <Filter {...props}>
    <TextInput label="Locator" source="locator" reference="locator" alwaysOn resettable />
  </Filter>
);

const MoveList = () => (
  <List pagination={<AdminPagination />} perPage={25} filters={<MoveFilter />} sort={defaultSort}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" reference="moves" />
      <TextField source="ordersId" reference="moves" label="Order Id" />
      <TextField source="serviceMember.id" label="Service Member Id" sortable={false} />
      <TextField source="locator" reference="moves" />
      <TextField source="status" reference="moves" />
      <TextField source="show" reference="moves" />
      <DateField source="createdAt" reference="moves" showTime />
      <DateField source="updatedAt" reference="moves" showTime />
    </Datagrid>
  </List>
);

export default MoveList;

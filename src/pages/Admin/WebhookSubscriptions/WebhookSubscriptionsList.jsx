/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { List, Datagrid, TextField, NumberField, DateField } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'status', order: 'ASC' };

const WebhookSubscriptionList = () => (
  <List pagination={<AdminPagination />} perPage={25} sort={defaultSort} bulkActionButtons={false}>
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="eventKey" />
      <TextField source="callbackUrl" />
      <NumberField source="severity" />
      <TextField source="status" />
      <DateField source="updatedAt" showTime />
    </Datagrid>
  </List>
);

export default WebhookSubscriptionList;

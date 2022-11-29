import React from 'react';
import { List, Datagrid, TextField, Filter, TextInput, DateField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import styles from 'scenes/SystemAdmin/Home.module.scss';

const defaultSort = { field: 'service_member_id', order: 'ASC' };

const NotificationFilter = (props) => (
  <Filter {...props} className={styles['system-admin-filters']}>
    <TextInput label="Service Member ID" source="serviceMemberId" reference="notifications" alwaysOn />
  </Filter>
);

const NotificationList = (props) => (
  <List {...props} pagination={<AdminPagination />} perPage={25} sort={defaultSort} filters={<NotificationFilter />}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="serviceMemberId" />
      <TextField source="sesMessageId" />
      <TextField source="notificationType" />
      <DateField source="createdAt" showTime />
    </Datagrid>
  </List>
);

export default NotificationList;

import React from 'react';
import { List, Datagrid, TextField, Filter, TextInput, DateField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import styles from 'scenes/SystemAdmin/Home.module.scss';

const defaultSort = { field: 'service_member_id', order: 'ASC' };

const NotificationFilter = props => (
  <Filter {...props} className={styles['system-admin-filters']}>
    <TextInput label="Service Member ID" source="service_member_id" reference="notifications" alwaysOn />
  </Filter>
);

const NotificationList = props => (
  <List
    {...props}
    pagination={<AdminPagination />}
    perPage={25}
    bulkActionButtons={false}
    sort={defaultSort}
    filters={<NotificationFilter />}
  >
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="service_member_id" />
      <TextField source="ses_message_id" />
      <TextField source="notification_type" />
      <DateField source="created_at" showTime />
    </Datagrid>
  </List>
);

export default NotificationList;

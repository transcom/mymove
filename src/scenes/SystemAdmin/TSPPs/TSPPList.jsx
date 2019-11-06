import React from 'react';
import { Datagrid, List, TextField, DateField, Filter, TextInput } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import styles from 'scenes/SystemAdmin/Home.module.scss';

const defaultSort = { field: 'performance_period_start', order: 'DESC' };

const TSPPFilter = props => (
  <Filter {...props} className={styles['system-admin-filters']}>
    <TextInput
      label="Traffic distribution list id"
      source="traffic_distribution_list_id"
      reference="transportation_service_provider_performances"
    />
    <TextInput
      label="Transportation service provider id"
      source="transportation_service_provider_id"
      reference="transportation_service_provider_performances"
    />
  </Filter>
);

const TSPPList = props => (
  <List
    {...props}
    pagination={<AdminPagination />}
    perPage={25}
    filters={<TSPPFilter />}
    sort={defaultSort}
    bulkActionButtons={false}
  >
    <Datagrid rowClick="show">
      <TextField source="id" reference="transportation_service_provider_performances" />
      <TextField source="traffic_distribution_list_id" reference="transportation_service_provider_performances" />
      <TextField source="transportation_service_provider_id" reference="transportation_service_provider_performances" />
      <DateField source="performance_period_start" reference="transportation_service_provider_performances" />
      <DateField source="performance_period_end" reference="transportation_service_provider_performances" />
      <DateField source="rate_cycle_start" reference="transportation_service_provider_performances" />
      <DateField source="rate_cycle_end" reference="transportation_service_provider_performances" />
    </Datagrid>
  </List>
);

export default TSPPList;

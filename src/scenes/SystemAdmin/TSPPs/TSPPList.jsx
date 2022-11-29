import React from 'react';
import { Datagrid, List, TextField, DateField, Filter, TextInput } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import styles from 'scenes/SystemAdmin/Home.module.scss';

const defaultSort = { field: 'performance_period_start', order: 'DESC' };

const TSPPFilter = (props) => (
  <Filter {...props} className={styles['system-admin-filters']}>
    <TextInput
      label="Traffic distribution list id"
      source="trafficDistributionListId"
      reference="transportation_service_provider_performances"
    />
    <TextInput
      label="Transportation service provider id"
      source="transportationServiceProviderId"
      reference="transportation_service_provider_performances"
    />
  </Filter>
);

const TSPPList = (props) => (
  <List {...props} pagination={<AdminPagination />} perPage={25} filters={<TSPPFilter />} sort={defaultSort}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" reference="transportation_service_provider_performances" />
      <TextField source="trafficDistributionListId" reference="transportation_service_provider_performances" />
      <TextField source="transportationServiceProviderId" reference="transportation_service_provider_performances" />
      <DateField source="performancePeriodStart" reference="transportation_service_provider_performances" />
      <DateField source="performancePeriodEnd" reference="transportation_service_provider_performances" />
      <DateField source="rateCycleStart" reference="transportation_service_provider_performances" />
      <DateField source="rateCycleEnd" reference="transportation_service_provider_performances" />
    </Datagrid>
  </List>
);

export default TSPPList;

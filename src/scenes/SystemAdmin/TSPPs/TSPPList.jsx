import React from 'react';
import { Datagrid, List, TextField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const TSPPList = props => (
    <List {...props} classes={{ card: "testing" }} pagination={<AdminPagination />} perPage={25} >
      <Datagrid>
        <TextField source="id" reference="transportation_service_provider_performances" />
        <TextField source="traffic_distribution_list_id" reference="transportation_service_provider_performances" />
        <TextField source="transportation_service_provider_id" reference="transportation_service_provider_performances" />
        <TextField source="performance_period_start" reference="transportation_service_provider_performances" />
        <TextField source="performance_period_end" reference="transportation_service_provider_performances" />
        <TextField source="rate_cycle_start" reference="transportation_service_provider_performances" />
        <TextField source="rate_cycle_end" reference="transportation_service_provider_performances" />
        <TextField source="quality_band" reference="transportation_service_provider_performances" />
        <TextField source="offer_count" reference="transportation_service_provider_performances" />
        <TextField source="best_value_score" reference="transportation_service_provider_performances" />
        <TextField source="linehaul_rate" reference="transportation_service_provider_performances" />
        <TextField source="sit_rate" reference="transportation_service_provider_performances" />
      </Datagrid>
    </List>
);

export default TSPPList;

import React from 'react';
import { TextField, DateField, NumberField, Show, SimpleShowLayout } from 'react-admin';

const TSPPShowTitle = ({ record }) => {
  return <span>{`${record.id}`}</span>;
};

const TSPPShow = props => {
  return (
    <Show {...props} title={<TSPPShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" reference="transportation_service_provider_performances" />
        <TextField source="traffic_distribution_list_id" reference="transportation_service_provider_performances" />
        <TextField
          source="transportation_service_provider_id"
          reference="transportation_service_provider_performances"
        />
        <DateField source="performance_period_start" reference="transportation_service_provider_performances" />
        <DateField source="performance_period_end" reference="transportation_service_provider_performances" />
        <DateField source="rate_cycle_start" reference="transportation_service_provider_performances" />
        <DateField source="rate_cycle_end" reference="transportation_service_provider_performances" />
        <TextField source="quality_band" reference="transportation_service_provider_performances" />
        <NumberField source="offer_count" reference="transportation_service_provider_performances" />
        <NumberField source="best_value_score" reference="transportation_service_provider_performances" />
        <NumberField
          source="linehaul_rate"
          reference="transportation_service_provider_performances"
          options={{ style: 'percent', maximumFractionDigits: 20 }}
        />
        <NumberField
          source="sit_rate"
          reference="transportation_service_provider_performances"
          options={{ style: 'percent', maximumFractionDigits: 20 }}
        />
      </SimpleShowLayout>
    </Show>
  );
};

export default TSPPShow;

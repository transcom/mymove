import React from 'react';
import { TextField, DateField, NumberField, Show, SimpleShowLayout } from 'react-admin';
import { useParams } from 'react-router-dom';

const TSPPShowTitle = () => {
  const { id } = useParams();
  return <span>{`TSPP ID: ${id}`}</span>;
};

const TSPPShow = () => {
  return (
    <Show title={<TSPPShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" reference="transportation_service_provider_performances" />
        <TextField source="trafficDistributionListId" reference="transportation_service_provider_performances" />
        <TextField source="transportationServiceProviderId" reference="transportation_service_provider_performances" />
        <DateField source="performancePeriodStart" reference="transportation_service_provider_performances" />
        <DateField source="performancePeriodEnd" reference="transportation_service_provider_performances" />
        <DateField source="rateCycleStart" reference="transportation_service_provider_performances" />
        <DateField source="rateCycleEnd" reference="transportation_service_provider_performances" />
        <TextField source="qualityBand" reference="transportation_service_provider_performances" />
        <NumberField source="offerCount" reference="transportation_service_provider_performances" />
        <NumberField source="bestValueScore" reference="transportation_service_provider_performances" />
        <NumberField
          source="linehaulRate"
          reference="transportation_service_provider_performances"
          options={{ style: 'percent', maximumFractionDigits: 2 }}
        />
        <NumberField
          source="sitRate"
          reference="transportation_service_provider_performances"
          options={{ style: 'percent', maximumFractionDigits: 2 }}
        />
      </SimpleShowLayout>
    </Show>
  );
};

export default TSPPShow;

import React from 'react';

import ImportantShipmentDates from './ImportantShipmentDates';

export default {
  title: 'Office Components/ImportantShipmentDate',
};

export const Default = () => (
  <ImportantShipmentDates
    requestedPickupDate="Thursday, 26 Mar 2020"
    scheduledPickupDate="Friday, 27 Mar 2020"
    requiredDeliveryDate="Monday, 30 Mar 2020"
  />
);

export const EmptyState = () => <ImportantShipmentDates />;

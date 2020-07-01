import React from 'react';

import ImportantShipmentDates from '../components/Office/ImportantShipmentDates';

export default {
  title: 'TOO&#47;TIO Components|ImportantShipmentDate',
};

export const Default = () => (
  <ImportantShipmentDates requestedPickupDate="Thursday, 26 Mar 2020" scheduledPickupDate="Friday, 27 Mar 2020" />
);

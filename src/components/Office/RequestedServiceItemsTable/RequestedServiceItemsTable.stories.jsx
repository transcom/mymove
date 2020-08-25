import React from 'react';

import RequestedServiceItemsTable from './RequestedServiceItemsTable';

export default {
  title: 'TOO/TIO Components|RequestedServiceItemsTable',
  component: RequestedServiceItemsTable,
};

const serviceItems = [
  {
    id: 'abc-123',
    submittedAt: '2020-11-20',
    serviceItem: 'Dom. Origin 1st Day SIT',
    code: 'DOFSIT',
    details: {
      pickupPostalCode: '60612',
      reason: "here's the reason",
    },
  },
  {
    id: 'abc-1234',
    submittedAt: '2020-11-20',
    serviceItem: 'Dom. Destination 1st Day SIT',
    code: 'DDFSIT',
    details: {
      firstCustomerContact: { timeMilitary: '', firstAvailableDeliveryDate: '2020-11-22' },
      secondCustomerContact: { timeMilitary: '', firstAvailableDeliveryDate: '2020-11-23' },
    },
  },
  {
    id: 'cba-123',
    submittedAt: '2020-11-20',
    serviceItem: 'Dom. Origin Shuttle Service',
    code: 'DOSHUT',
    details: {
      reason: "Here's the reason",
    },
  },
  {
    id: 'cba-1234',
    submittedAt: '2020-11-20',
    serviceItem: 'Dom. Destination Shuttle Service',
    code: 'DDSHUT',
    details: {
      reason: "Here's the reason",
    },
  },
  {
    id: 'abc12345',
    submittedAt: '2020-11-20',
    serviceItem: 'Dom. Crating',
    code: 'DCRT',
    details: {
      description: "Here's the description",
      itemDimensions: { length: 8400, width: 2600, height: 4200 },
      crateDimensions: { length: 110000, width: 36000, height: 54000 },
      imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
    },
  },
];

export const Default = () => <RequestedServiceItemsTable serviceItems={serviceItems} />;

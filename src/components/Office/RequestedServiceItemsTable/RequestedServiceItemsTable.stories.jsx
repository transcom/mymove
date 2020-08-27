import React from 'react';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

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
    code: 'DOMSIT',
    status: 'SUBMITTED',
    details: {
      text: {
        ZIP: '60612',
        Reason: "here's the reason",
      },
      imgURL: null,
    },
  },
  {
    id: 'abc-1234',
    submittedAt: '2020-11-20',
    serviceItem: 'Dom. Destination 1st Day SIT',
    code: 'DDFSIT',
    status: 'SUBMITTED',
    details: {
      text: {
        'First available delivery date': '22 Nov 2020',
        'First customer contact': '22 Nov 2020 12:00pm',
        'Second customer contact': '22 Nov 2020 12:00pm',
      },
      imgURL: null,
    },
  },
  {
    id: 'cba-123',
    submittedAt: '2020-11-20',
    serviceItem: 'Dom. Origin Shuttle Service',
    code: 'DOSHUT',
    status: 'SUBMITTED',
    details: {
      text: {
        'Reason for request': "Here's the reason",
        'Estimated weight': '3,500lbs',
      },
      imgURL: null,
    },
  },
  {
    id: 'cba-1234',
    submittedAt: '2020-11-20',
    serviceItem: 'Dom. Destination Shuttle Service',
    code: 'DDSHUT',
    status: 'SUBMITTED',
    details: {
      text: {
        'Reason for request': "Here's the reason",
        'Estimated weight': '3,500lbs',
      },
      imgURL: null,
    },
  },
  {
    id: 'abc12345',
    submittedAt: '2020-11-20',
    serviceItem: 'Dom. Crating',
    code: 'DCRT',
    status: 'SUBMITTED',
    details: {
      text: {
        Description: "Here's the description",
        'Item dimensions': '84"x26"x42"',
        'Crate dimensions': '110"x36"x54"',
      },
      imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
    },
  },
];

const approvedServiceItems = serviceItems.map((serviceItem) => {
  return { ...serviceItem, status: 'APPROVED' };
});
const rejectedServiceItems = serviceItems.map((serviceItem) => {
  return { ...serviceItem, status: 'REJECTED' };
});

export const Default = () => (
  <RequestedServiceItemsTable serviceItems={serviceItems} statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED} />
);
export const ApprovedServiceItems = () => (
  <RequestedServiceItemsTable serviceItems={approvedServiceItems} statusForTableType={SERVICE_ITEM_STATUS.APPROVED} />
);
export const RejectedServiceItems = () => (
  <RequestedServiceItemsTable serviceItems={rejectedServiceItems} statusForTableType={SERVICE_ITEM_STATUS.REJECTED} />
);

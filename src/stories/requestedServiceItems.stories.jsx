import React from 'react';

import { storiesOf } from '@storybook/react';

import RequestedServiceItemsTable from '../components/Office/RequestedServiceItemsTable';

storiesOf('TOO/TIO Components|RequestedServiceItemsTable', module).add('default', () => {
  const serviceItems = [
    {
      id: 'abc-123',
      dateRequested: '20 Nov 2020',
      serviceItem: 'Dom. Origin 1st Day SIT',
      code: 'DOMSIT',
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
      dateRequested: '22 Nov 2020',
      serviceItem: 'Dom. Destination 1st Day SIT',
      code: 'DDFSIT',
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
      dateRequested: '22 Nov 2020',
      serviceItem: 'Dom. Origin Shuttle Service',
      code: 'DOSHUT',
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
      dateRequested: '22 Nov 2020',
      serviceItem: 'Dom. Destination Shuttle Service',
      code: 'DDSHUT',
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
      dateRequested: '22 Nov 2020',
      serviceItem: 'Dom. Crating',
      code: 'DCRT',
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

  return <RequestedServiceItemsTable serviceItems={serviceItems} />;
});

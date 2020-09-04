import React from 'react';
import { withKnobs, text, object } from '@storybook/addon-knobs';

import OrdersTable from './OrdersTable';

export default {
  title: 'TOO/TIO Components|OrdersTable',
  component: OrdersTable,
  decorators: [
    withKnobs,
    (Story) => (
      <div style={{ 'max-width': '800px' }}>
        <Story />
      </div>
    ),
  ],
};

export const Basic = () => (
  <OrdersTable
    ordersInfo={{
      currentDutyStation: object('ordersInfo.currentDutyStation', { name: 'JBSA Lackland' }),
      newDutyStation: object('ordersInfo.newDutyStation', { name: 'JB Lewis-McChord' }),
      issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
      reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
      departmentIndicator: text('ordersInfo.departmentIndicator', 'NAVY_AND_MARINES'),
      ordersNumber: text('ordersInfo.ordersNumber', '999999999'),
      ordersType: text('ordersInfo.ordersType', 'PERMANENT_CHANGE_OF_STATION'),
      ordersTypeDetail: text('ordersInfo.ordersTypeDetail', 'HHG_PERMITTED'),
      tacMDC: text('ordersInfo.tacMDC', '9999'),
      sacSDN: text('ordersInfo.sacSDN', '999 999999 999'),
    }}
  />
);

import React from 'react';

import OrdersTable from './OrdersTable';

import { ORDERS_TYPE_OPTIONS } from 'constants/orders';

export default {
  title: 'Customer Components / OrdersTable',
  component: OrdersTable,
  decorators: [
    (Story) => (
      <div style={{ padding: 40 }}>
        <Story />
      </div>
    ),
  ],
  argTypes: {
    onEditClick: { action: 'orders edit button clicked' },
  },
};

const defaultProps = {
  moveId: 'abc123',
  orderType: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION,
  issueDate: '11 June 2020',
  reportByDate: '11 Aug 2020',
  newDutyStationName: 'Fort Knox',
  hasDependents: true,
  uploads: [{ id: 1 }, { id: 2 }, { id: 3 }],
};

const OrdersTableTemplate = (args) => <OrdersTable {...args} />;
export const Basic = OrdersTableTemplate.bind({});
Basic.args = defaultProps;

import React from 'react';

import OrdersTable from './OrdersTable';

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
  orderType: 'Permanent change of station',
  issueDate: '11 June 2020',
  reportByDate: '11 Aug 2020',
  newDutyStationName: 'Fort Knox',
  hasDependents: true,
  uploads: [{ id: 1 }, { id: 2 }, { id: 3 }],
};

const OrdersTableTemplate = (args) => <OrdersTable {...args} />;
export const Basic = OrdersTableTemplate.bind({});
Basic.args = defaultProps;

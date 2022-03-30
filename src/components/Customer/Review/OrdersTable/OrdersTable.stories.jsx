import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import OrdersTable from './OrdersTable';

import { ORDERS_TYPE_OPTIONS } from 'constants/orders';

export default {
  title: 'Customer Components / OrdersTable',
  component: OrdersTable,
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
  argTypes: {
    onEditClick: { action: 'orders edit button clicked' },
  },
};

const defaultProps = {
  moveId: 'abc123',
  orderType: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION,
  issueDate: '2020-06-11',
  reportByDate: '2020-08-11',
  newDutyLocationName: 'Fort Knox',
  hasDependents: true,
  uploads: [{ id: 1 }, { id: 2 }, { id: 3 }],
};

const OrdersTableTemplate = (args) => <OrdersTable {...args} />;
export const Basic = OrdersTableTemplate.bind({});
Basic.args = defaultProps;

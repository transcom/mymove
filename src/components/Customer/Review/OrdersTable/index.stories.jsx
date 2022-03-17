/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import OrdersTable from '.';

const defaultProps = {
  orderType: 'Permanent change of station',
  issueDate: '11 June 2020',
  reportByDate: '11 Aug 2020',
  newDutyLocationName: 'Fort Knox',
  hasDependents: true,
  uploads: [1, 2, 3],
};

export default {
  title: 'Customer Components / OrdersTable',
};

export const Basic = () => (
  <div style={{ padding: 40 }}>
    <OrdersTable {...defaultProps} />
  </div>
);

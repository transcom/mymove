import React from 'react';
import { node } from 'prop-types';
import { number } from '@storybook/addon-knobs';

import EditBillableWeight from './EditBillableWeight';

export default {
  title: 'Office Components/EditBillableWeight',
  component: EditBillableWeight,
};

const Container = ({ children }) => <div style={{ width: 336, margin: '0 auto' }}>{children}</div>;

Container.propTypes = {
  children: node.isRequired,
};

export const MaxbillableWeight = () => (
  <Container>
    <EditBillableWeight
      title="Max billable weight"
      weightAllowance={number('WeightAllowance', 8000)}
      estimatedWeight={number('EstimatedWeight', 13750)}
    />
  </Container>
);

export const BillableWeight = () => (
  <Container>
    <EditBillableWeight
      title="Billable weight"
      originalWeight={number('OriginalWeight', 10000)}
      estimatedWeight={number('EstimatedWeight', 13000)}
      maxBillableWeight={number('MaxBillableWeight', 6000)}
      billableWeight={number('BillableWeight', 7000)}
      totalBillableWeight={number('TotalBillableWeight', 11000)}
    />
  </Container>
);

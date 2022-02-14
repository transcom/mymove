import React from 'react';
import { node } from 'prop-types';

import EditBillableWeight from './EditBillableWeight';

export default {
  title: 'Office Components/EditBillableWeight',
  component: EditBillableWeight,
};

const Container = ({ children }) => <div style={{ width: 336, margin: '0 auto' }}>{children}</div>;

Container.propTypes = {
  children: node.isRequired,
};
const Template = (args) => (
  <Container>
    <EditBillableWeight {...args} />
  </Container>
);

const ToggledTemplate = (args) => (
  <Container>
    <EditBillableWeight {...args} showFieldsInitial />
  </Container>
);

export const MaxBillableWeight = Template.bind({});

MaxBillableWeight.args = {
  billableWeightJustification: 'Reduced billable weight to cap at 110% of estimated.',
  estimatedWeight: 13750,
  maxBillableWeight: 13000,
  title: 'Max billable weight',
  weightAllowance: 8000,
};

export const EmptyMaxBillableWeight = ToggledTemplate.bind({});

EmptyMaxBillableWeight.args = {
  billableWeightJustification: '',
  estimatedWeight: 13750,
  title: 'Max billable weight',
  weightAllowance: 8000,
};

export const EmptyMaxBillableWeightNTSRelease = ToggledTemplate.bind({});

EmptyMaxBillableWeightNTSRelease.args = {
  billableWeightJustification: '',
  estimatedWeight: 13750,
  title: 'Max billable weight',
  weightAllowance: 8000,
  isNTSRShipment: true,
};

export const BillableWeight = Template.bind({});

BillableWeight.args = {
  billableWeight: 7000,
  billableWeightJustification: 'Reduced billable weight to cap at 110% of estimated.',
  estimatedWeight: 13000,
  maxBillableWeight: 6000,
  originalWeight: 10000,
  title: 'Billable weight',
  totalBillableWeight: 11000,
};

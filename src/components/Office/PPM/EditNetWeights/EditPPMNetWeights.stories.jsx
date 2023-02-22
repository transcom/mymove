import React from 'react';
import { node } from 'prop-types';

import EditPPMNetWeight from './EditPPMNetWeight';

export default {
  title: 'Office Components/PPM/EditPPMNetWeights',
  component: EditPPMNetWeight,
};

const Container = ({ children }) => <div style={{ width: 336, margin: '0 auto' }}>{children}</div>;

Container.propTypes = {
  children: node.isRequired,
};
const Template = (args) => (
  <Container>
    <EditPPMNetWeight {...args} />
  </Container>
);

export const EditPPMNetWeightDefault = Template.bind({});

EditPPMNetWeightDefault.args = {
  maxBillableWeight: 19500,
  originalWeight: 4500,
  weightAllowance: 18000,
  totalBillableWeight: 11000,
  estimatedWeight: 13000,
  billableWeight: 7000,
  ppmNetWeightRemarks: 'Everything seems fine',
};

export const EditPPMNetWeightExcessWeight = Template.bind({});
EditPPMNetWeightExcessWeight.args = {
  maxBillableWeight: 19500,
  originalWeight: 4500,
  weightAllowance: 18000,
  totalBillableWeight: 21000,
  estimatedWeight: 13000,
  billableWeight: 7000,
  ppmNetWeightRemarks: '',
};

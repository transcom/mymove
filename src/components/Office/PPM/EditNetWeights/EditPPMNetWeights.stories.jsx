import React from 'react';
import { node } from 'prop-types';

import EditPPMNetWeight from './EditPPMNetWeight';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';

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
  netWeightRemarks: '',
  weightAllowance: 5000,
  weightTicket: {
    fullWeight: 1200,
    emptyWeight: 200,
  },
  shipments: [
    { primeActualWeight: 1000, reweigh: null, status: 'APPROVED' },
    { primeActualWeight: 2000, reweigh: { weight: 1000 }, status: 'APPROVED' },
    {
      ppmShipment: {
        weightTickets: [
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        ],
      },
      status: 'APPROVED',
    },
  ],
};

export const EditPPMNetWeightExcessWeight = Template.bind({});
EditPPMNetWeightExcessWeight.args = {
  netWeightRemarks: '',
  weightAllowance: 5000,
  weightTicket: {
    fullWeight: 1200,
    emptyWeight: 200,
  },
  shipments: [
    { primeActualWeight: 1200, reweigh: null, status: 'APPROVED' },
    { primeActualWeight: 2000, reweigh: { weight: 3000 }, status: 'APPROVED' },
    {
      ppmShipment: {
        weightTickets: [
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        ],
      },
      status: 'APPROVED',
    },
  ],
};

export const EditPPMNetWeightReduceWeight = Template.bind({});
EditPPMNetWeightReduceWeight.args = {
  netWeightRemarks: '',
  weightAllowance: 5000,
  weightTicket: {fullWeight: 1200, emptyWeight: 200},
  shipments: [
    { primeActualWeight: 1200, reweigh: null, status: 'APPROVED' },
    { primeActualWeight: 6000, reweigh: { weight: 5000 }, status: 'APPROVED' },
    {
      ppmShipment: {
        weightTickets: [
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        ],
      },
      status: 'APPROVED',
    },
    {
      ppmShipment: {
        weightTickets: [
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        ],
      },
      status: 'APPROVED',
    },
  ],
};

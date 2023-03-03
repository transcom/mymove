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
  netWeightRemarks: '',
  weightAllowance: 8000,
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 600,
    fullWeight: 1200,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
    shipments: [
      { primeActualWeight: { weight: 1200 }, reweigh: null },
      { primeActualWeight: { weight: 4800 }, reweigh: { weight: 5000 } },
      {
        ppmShipment: {
          weightTickets: [
            {
              vehicleDescription: 'Kia Forte',
              emptyWeight: 600,
              fullWeight: 1200,
              ownsTrailer: true,
              trailerMeetsCriteria: false,
            },
            {
              vehicleDescription: 'Kia Soul',
              emptyWeight: 1200,
              fullWeight: 2000,
              ownsTrailer: false,
              trailerMeetsCriteria: false,
            },
          ],
        },
      },
    ],
  },
};

export const EditPPMNetWeightExcessWeight = Template.bind({});
EditPPMNetWeightExcessWeight.args = {
  netWeightRemarks: '',
  weightAllowance: 8000,
  weightTicket: {
    emptyWeight: 600,
    fullWeight: 18000,
  },
  shipments: [
    { primeActualWeight: { weight: 1200 }, reweigh: null },
    { primeActualWeight: { weight: 4800 }, reweigh: { weight: 5000 } },
    {
      ppmShipment: {
        weightTickets: [
          {
            emptyWeight: 600,
            fullWeight: 1200,
          },
          {
            emptyWeight: 1200,
            fullWeight: 2000,
          },
        ],
      },
    },
  ],
};

export const EditPPMNetWeightReduceWeight = Template.bind({});
EditPPMNetWeightExcessWeight.args = {
  netWeightRemarks: '',
  weightAllowance: 8000,
  weightTicket: { emptyWeight: 600, fullWeight: 10000 },
  shipments: [
    { primeActualWeight: { weight: 1200 }, reweigh: null },
    { primeActualWeight: { weight: 4800 }, reweigh: { weight: 5000 } },
    {
      ppmShipment: {
        weightTickets: [
          {
            emptyWeight: 600,
            fullWeight: 1200,
          },
          {
            emptyWeight: 1200,
            fullWeight: 2000,
          },
        ],
      },
    },
    {
      ppmShipment: {
        weightTickets: [
          {
            emptyWeight: 600,
            fullWeight: 1200,
          },
          {
            emptyWeight: 1200,
            fullWeight: 2000,
          },
        ],
      },
    },
  ],
};

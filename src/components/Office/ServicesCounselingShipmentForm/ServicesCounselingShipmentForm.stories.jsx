/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { action } from '@storybook/addon-actions';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ServicesCounselingShipmentForm from './ServicesCounselingShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

const defaultProps = {
  match: {
    isExact: false,
    path: '/counseling/moves/:moveId/shipments/:mtoShipmentId/',
    url: '',
    params: { moveCode: 'move123' },
  },
  moveTaskOrderID: 'task123',
  history: { push: () => {} },
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
    streetAddress1: '123 Main',
  },
  useCurrentResidence: false,
  mtoShipment: {
    destinationAddress: undefined,
  },
  serviceMember: {
    weightAllotment: {
      totalWeightSelf: 5000,
    },
  },
  isCreatePage: true,
  submitHandler: action('submit MTO Shipment for create or update'),
};

const mockMtoShipment = {
  id: 'mock id',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock customer remarks',
  counselorRemarks: 'mock counselor remarks',
  requestedPickupDate: '2020-03-01',
  requestedDeliveryDate: '2020-03-30',
  agents: [
    {
      firstName: 'mock receiving',
      lastName: 'agent',
      telephone: '2225551234',
      email: 'mock.delivery.agent@example.com',
      agentType: 'RECEIVING_AGENT',
    },
    {
      firstName: 'Mock Releasing',
      lastName: 'Agent Jr, PhD, MD, DDS',
      telephone: '3335551234',
      email: 'mock.pickup.agent@example.com',
      agentType: 'RELEASING_AGENT',
    },
  ],
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

const mockMtoShipmentNoCustomerRemarks = {
  ...mockMtoShipment,
  customerRemarks: '',
};

export default {
  title: 'Office Components / Forms / ServicesCounselingShipmentForm',
  component: ServicesCounselingShipmentForm,
  decorators: [
    (Story) => (
      <GridContainer className={styles.gridContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

// create shipment stories (form should not prefill customer data)
export const HHGShipment = () => (
  <ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.HHG} />
);

// edit shipment stories (form should prefill)
export const EditHHGShipment = () => (
  <ServicesCounselingShipmentForm
    {...defaultProps}
    selectedMoveType={SHIPMENT_OPTIONS.HHG}
    isCreatePage={false}
    mtoShipment={mockMtoShipment}
  />
);

// edit shipment stories, no customer remarks (form should prefill)
export const EditHHGShipmentNoCustRemarks = () => (
  <ServicesCounselingShipmentForm
    {...defaultProps}
    selectedMoveType={SHIPMENT_OPTIONS.HHG}
    isCreatePage={false}
    mtoShipment={mockMtoShipmentNoCustomerRemarks}
  />
);

export const NTSShipmentWithoutCodes = () => {
  return <ServicesCounselingShipmentForm {...defaultProps} selectedMoveType={SHIPMENT_OPTIONS.NTS} />;
};

export const NTSShipmentWithCodes = () => {
  return (
    <ServicesCounselingShipmentForm
      {...defaultProps}
      selectedMoveType={SHIPMENT_OPTIONS.NTS}
      TACs={{ HHG: '1234', NTS: '5678' }}
      SACs={{ HHG: '000012345' }}
    />
  );
};

export const NTSReleaseShipment = () => {
  return (
    <ServicesCounselingShipmentForm
      {...defaultProps}
      selectedMoveType={SHIPMENT_OPTIONS.NTSR}
      TACs={{ HHG: '1234', NTS: '5678' }}
      SACs={{ HHG: '000012345', NTS: '6789ABC' }}
    />
  );
};

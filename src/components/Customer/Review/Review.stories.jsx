import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import { Summary } from 'components/Customer/Review/Summary/Summary';
import { MOVE_STATUSES, SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';

export default {
  title: 'Customer Components / Review Shipment',
  component: Summary,
  decorators: [
    (Story) => (
      <MockProviders>
        <GridContainer>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <Story />
            </Grid>
          </Grid>
        </GridContainer>
      </MockProviders>
    ),
  ],
};

const noop = () => {};
const customerUuid = 'customerUuid';
const mtoUuid = 'mtoUuid';
const mtoLocator = 'XYZ890';

const defaultProps = {
  currentMove: {
    id: mtoUuid,
    locator: mtoLocator,
    service_member_id: customerUuid,
    status: MOVE_STATUSES.DRAFT,
  },
  currentOrders: {
    orders_type: 'PERMANENT_CHANGE_OF_STATION',
    has_dependents: false,
    issue_date: '2020-08-11',
    grade: 'E-1',
    moves: [mtoUuid],
    origin_duty_location: {
      name: 'Fort Drum',
      address: {
        postalCode: '13643',
      },
    },
    new_duty_location: {
      name: 'Naval Base San Diego',
      address: {
        postalCode: '91945',
      },
    },
    report_by_date: '2020-08-31',
    service_member_id: customerUuid,
    spouse_has_pro_gear: false,
    status: MOVE_STATUSES.DRAFT,
    uploaded_orders: {
      uploads: [],
    },
  },
  history: {
    back: noop,
    push: noop,
  },
  match: {
    url: `/moves/${mtoLocator}/review`,
    params: {
      moveId: mtoLocator,
    },
  },
  moveIsApproved: false,
  mtoShipments: [],
  serviceMember: {
    id: customerUuid,
    current_location: {
      name: 'Fort Drum',
    },
    residential_address: {
      city: 'Great Bend',
      postalCode: '13643',
      state: 'NY',
      streetAddress1: '448 Washington Blvd NE',
    },
    affiliation: 'NAVY',
    edipi: '1231231231',
    personal_email: 'test@example.com',
    first_name: 'Jason',
    last_name: 'Ash',
    rank: 'E_1',
    telephone: '323-555-7890',
  },
};

const HHGShipment = {
  id: 'hhgShipmentUuid',
  agents: [],
  customerRemarks: 'contains grandfather clock, garbage cat',
  moveTaskOrderID: mtoUuid,
  pickupAddress: {
    city: 'Great Bend',
    postalCode: '13643',
    state: 'NY',
    streetAddress1: '448 Washington Blvd NE',
  },
  requestedDeliveryDate: '2020-08-31',
  requestedPickupDate: '2020-08-31',
  shipmentType: 'HHG',
  status: MOVE_STATUSES.SUBMITTED,
  createdAt: '2020-09-01T21:00:00.000Z',
  updatedAt: '2020-09-02T21:08:38.392Z',
};

const PPMShipment = {
  id: 'ppmShipmentUuid',
  moveTaskOrderID: mtoUuid,
  status: MOVE_STATUSES.DRAFT,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  created_at: '2020-09-01T22:00:00.000Z',
  updated_at: '2020-09-01T22:30:00.000Z',
  ppmShipment: {
    expectedDepartureDate: '2021-06-23',
    pickupPostalCode: '13643',
    destinationPostalCode: '91945',
    sitExpected: false,
    estimatedWeight: 5000,
    estimatedIncentive: 1000000,
  },
};

export const WithNoShipments = () => {
  return (
    <MockProviders>
      <Summary {...defaultProps} />
    </MockProviders>
  );
};

export const WithHHGShipment = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [HHGShipment],
  };

  return (
    <MockProviders>
      <Summary {...props} />
    </MockProviders>
  );
};

export const WithPPM = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [PPMShipment],
  };
  return (
    <MockProviders>
      <Summary {...props} />
    </MockProviders>
  );
};

export const AsSubmitted = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [HHGShipment, PPMShipment],
    currentMove: {
      ...defaultProps.currentMove,
      status: MOVE_STATUSES.SUBMITTED,
    },
  };

  return (
    <MockProviders>
      <Summary {...props} />
    </MockProviders>
  );
};

export const AsApproved = () => {
  const approvedShipment = {
    ...HHGShipment,
    status: MOVE_STATUSES.SUBMITTED,
  };

  const props = {
    ...defaultProps,
    mtoShipments: [approvedShipment],
    moveIsApproved: true,
    currentMove: {
      ...defaultProps.currentMove,
      status: MOVE_STATUSES.SUBMITTED,
    },
  };

  return (
    <MockProviders>
      <Summary {...props} />
    </MockProviders>
  );
};

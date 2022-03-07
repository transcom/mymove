import React from 'react';

import { Summary } from './Summary';

import { MOVE_STATUSES } from 'shared/constants';
import { MockProviders } from 'testUtils';

export default {
  title: 'Customer Components / Review Shipment',
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
    affiliation: 'Navy',
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
  pickup_postal_code: '13643',
  destination_postal_code: '91945',
  original_move_date: '2021-06-23',
  created_at: '2020-09-01T22:00:00.000Z',
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
    currentPPM: PPMShipment,
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
    mtoShipments: [HHGShipment],
    currentMove: {
      ...defaultProps.currentMove,
      status: MOVE_STATUSES.SUBMITTED,
    },
    currentPPM: PPMShipment,
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
    currentPPM: PPMShipment,
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

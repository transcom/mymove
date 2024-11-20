import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import { Summary } from 'components/Customer/Review/Summary/Summary';
import { MOVE_STATUSES, SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { shipmentStatuses } from 'constants/shipments';
import { ORDERS_TYPE } from 'constants/orders';

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

const serviceMemberMoves = {
  currentMove: [
    {
      id: mtoUuid,
      moveCode: mtoLocator,
      mtoShipments: [],
      counselingOffice: {
        name: 'PPPO McAlester',
      },
      orders: {
        authorizedWeight: 11000,
        created_at: '2024-03-12T13:36:14.940Z',
        entitlement: {
          proGear: 2000,
          proGearSpouse: 500,
        },
        grade: 'E_7',
        has_dependents: false,
        id: 'orderId',
        issue_date: '2024-04-18',
        new_duty_location: {
          address: {
            city: 'Flagstaff',
            country: 'United States',
            id: '02df4469-90e5-4dbe-b2e0-69c7f8367912',
            postalCode: '86004',
            state: 'AZ',
            streetAddress1: 'n/a',
          },
          address_id: '02df4469-90e5-4dbe-b2e0-69c7f8367912',
          affiliation: null,
          created_at: '2024-02-27T20:40:42.164Z',
          id: '6af688f3-7be2-422e-a07a-2a26a4069ec4',
          name: 'Flagstaff, AZ 86004',
          updated_at: '2024-02-27T20:40:42.164Z',
        },
        orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
        originDutyLocationGbloc: 'HAFC',
        origin_duty_location: {
          address: {
            city: 'McAlester',
            country: 'United States',
            id: '7eccd9bc-c48b-4822-9324-7b0baa4256d1',
            postalCode: '74501',
            state: 'OK',
            streetAddress1: 'n/a',
          },
          address_id: '7eccd9bc-c48b-4822-9324-7b0baa4256d1',
          affiliation: 'ARMY',
          created_at: '2024-02-27T20:40:47.436Z',
          id: 'e017cd1a-a2b1-4e64-b887-e6e06e2d3f6c',
          name: 'McAlester Army Ammunition Plant, OK 74501',
          transportation_office: {
            address: {
              city: 'McAlester',
              country: 'United States',
              id: 'b52f5e75-620e-49c4-ac64-1330f495c956',
              postalCode: '74501',
              state: 'OK',
              streetAddress1: '1 C Tree Rd',
              streetAddress2: 'Bldg 31',
            },
            created_at: '2018-05-28T14:27:38.004Z',
            gbloc: 'HAFC',
            id: 'd1359c20-c762-4b04-9ed6-fd2b9060615b',
            name: 'PPPO McAlester - USA',
            phone_lines: [],
            updated_at: '2018-05-28T14:27:38.004Z',
          },
          transportation_office_id: 'd1359c20-c762-4b04-9ed6-fd2b9060615b',
          updated_at: '2024-02-27T20:40:47.436Z',
        },
        report_by_date: '2024-03-28',
        service_member_id: 'ab65bff1-4da3-4a51-a36f-2bdb2c4edc4d',
        spouse_has_pro_gear: false,
        status: 'DRAFT',
        updated_at: '2024-03-12T13:36:14.940Z',
        uploaded_orders: {
          id: '1d546f35-0a09-4007-9551-7934f148459d',
          service_member_id: 'ab65bff1-4da3-4a51-a36f-2bdb2c4edc4d',
          uploads: [
            {
              bytes: 1137126,
              contentType: 'image/png',
              createdAt: '2024-03-12T13:36:21.868Z',
              filename: 'Screenshot 2024-02-15 at 12.22.53 PM.png',
              id: '93374041-54a2-4ccd-9ff9-d8acb4cf440f',
              status: 'PROCESSING',
              updatedAt: '2024-03-12T13:36:21.868Z',
              url: '/storage/user/71d01de6-2181-45af-bb2b-63328fffe194/uploads/93374041-54a2-4ccd-9ff9-d8acb4cf440f?contentType=image%2Fpng',
            },
          ],
        },
      },
      status: MOVE_STATUSES.DRAFT,
      submittedAt: '0001-01-01T00:00:00.000Z',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
  previousMoes: [],
};

const defaultProps = {
  currentMove: {
    id: mtoUuid,
    locator: mtoLocator,
    service_member_id: customerUuid,
    status: MOVE_STATUSES.DRAFT,
    counseling_office: {
      name: 'PPPO McAlester',
    },
  },
  currentOrders: {
    orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    has_dependents: false,
    issue_date: '2020-08-11',
    grade: 'E_1',
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
  router: {
    location: {
      pathname: `/moves/${mtoUuid}/review`,
      search: '',
    },
    navigate: noop,
    params: { moveId: mtoUuid },
  },
  onDidMount: noop,
  setMsg: noop,
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
    telephone: '323-555-7890',
  },
  serviceMemberMoves,
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
    sitExpected: false,
    estimatedWeight: 5000,
    estimatedIncentive: 1000000,
    hasRequestedAdvance: false,
  },
};

const IncompeletePPMShipment = {
  id: 'ppmShipmentUuid',
  moveTaskOrderID: mtoUuid,
  status: MOVE_STATUSES.DRAFT,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  created_at: '2020-09-01T22:00:00.000Z',
  updated_at: '2020-09-01T22:30:00.000Z',
  ppmShipment: {
    expectedDepartureDate: '2021-06-23',
    sitExpected: false,
    hasRequestedAdvance: null,
  },
};

const serviceMemberMovesWithHhgShipment = {
  ...serviceMemberMoves,
  currentMove: [
    {
      ...serviceMemberMoves.currentMove[0],
      mtoShipments: [HHGShipment],
    },
  ],
};

const serviceMemberMovesWithPpmShipment = {
  ...serviceMemberMoves,
  currentMove: [
    {
      ...serviceMemberMoves.currentMove[0],
      mtoShipments: [PPMShipment],
    },
  ],
};

const serviceMemberMovesWithIncompletePpmShipment = {
  ...serviceMemberMoves,
  currentMove: [
    {
      ...serviceMemberMoves.currentMove[0],
      mtoShipments: [IncompeletePPMShipment],
    },
  ],
};

const serviceMemberMovesSubmitted = {
  ...serviceMemberMoves,
  currentMove: [
    {
      ...serviceMemberMoves.currentMove[0],
      mtoShipments: [HHGShipment, PPMShipment],
      status: MOVE_STATUSES.SUBMITTED,
    },
  ],
};

const serviceMemberMovesApproved = {
  ...serviceMemberMoves,
  currentMove: [
    {
      ...serviceMemberMoves.currentMove[0],
      mtoShipments: [HHGShipment, PPMShipment],
      status: MOVE_STATUSES.APPROVED,
    },
  ],
};

export const WithNoShipments = () => {
  return <Summary {...defaultProps} />;
};

export const WithHHGShipment = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [HHGShipment],
    serviceMemberMoves: serviceMemberMovesWithHhgShipment,
  };

  return <Summary {...props} />;
};

export const WithPPM = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [PPMShipment],
    serviceMemberMoves: serviceMemberMovesWithPpmShipment,
  };
  return <Summary {...props} />;
};

export const WithIncompletePPM = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [IncompeletePPMShipment],
    serviceMemberMoves: serviceMemberMovesWithIncompletePpmShipment,
  };
  return <Summary {...props} />;
};

export const AsSubmitted = () => {
  const props = {
    ...defaultProps,
    mtoShipments: [HHGShipment, PPMShipment],
    currentMove: {
      ...defaultProps.currentMove,
      status: MOVE_STATUSES.SUBMITTED,
    },
    serviceMemberMoves: serviceMemberMovesSubmitted,
  };

  return <Summary {...props} />;
};

export const AsApproved = () => {
  const approvedHHGShipment = {
    ...HHGShipment,
    status: shipmentStatuses.APPROVED,
  };

  const approvedPPMShipment = {
    ...PPMShipment,
    status: shipmentStatuses.APPROVED,
  };

  const props = {
    ...defaultProps,
    mtoShipments: [approvedHHGShipment, approvedPPMShipment],
    moveIsApproved: true,
    currentMove: {
      ...defaultProps.currentMove,
      status: MOVE_STATUSES.APPROVED,
    },
    serviceMemberMoves: serviceMemberMovesApproved,
  };

  return <Summary {...props} />;
};

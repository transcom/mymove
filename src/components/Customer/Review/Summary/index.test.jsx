import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { Summary } from './index';

import { MOVE_STATUSES } from 'shared/constants';
import { validateEntitlement } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  validateEntitlement: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const testProps = {
  serviceMember: {
    id: '666',
    current_station: {
      name: 'Test Duty Station',
    },
    residential_address: {
      city: 'New York',
      postalCode: '10001',
      state: 'NY',
      streetAddress1: '123 Main St',
    },
    affiliation: 'Navy',
    edipi: '123567890',
    personal_email: 'test@email.com',
    first_name: 'Tester',
    last_name: 'Testing',
    rank: 'RANK',
    telephone: '123-555-7890',
  },
  currentOrders: {
    orders_type: 'PERMANENT_CHANGE_OF_STATION',
    has_dependents: false,
    issue_date: '2020-08-11',
    grade: 'RANK',
    moves: ['123'],
    origin_duty_location: {
      name: 'Test Duty Station',
      address: {
        postalCode: '123456',
      },
    },
    new_duty_location: {
      name: 'New Test Duty Station',
      address: {
        postalCode: '123456',
      },
    },
    report_by_date: '2020-08-31',
    service_member_id: '666',
    spouse_has_pro_gear: false,
    status: MOVE_STATUSES.DRAFT,
    uploaded_orders: {
      uploads: [],
    },
  },
  match: { path: '', isExact: true, url: '/moves/123/review', params: { moveId: '123' } },
  history: {
    goBack: jest.fn(),
    push: jest.fn(),
  },
  currentMove: {
    id: '123',
    locator: 'CXVV3F',
    selected_move_type: 'HHG',
    service_member_id: '666',
    status: MOVE_STATUSES.DRAFT,
  },
  selectedMoveType: 'HHG',
  moveIsApproved: false,
  entitlement: {},
  mtoShipment: {
    id: 'testMtoShipment789',
    agents: [],
    customerRemarks: 'please be carefule',
    moveTaskOrderID: '123',
    pickupAddress: {
      city: 'Beverly Hills',
    },
    requestedDeliveryDate: '2020-08-31',
    requestedPickupDate: '2020-08-31',
    shipmentType: 'HHG',
    status: MOVE_STATUSES.SUBMITTED,
    updatedAt: '2020-09-02T21:08:38.392Z',
  },
  mtoShipments: [
    {
      id: 'testMtoShipment789',
      agents: [],
      customerRemarks: 'please be carefule',
      moveTaskOrderID: '123',
      pickupAddress: {
        city: 'Beverly Hills',
      },
      requestedDeliveryDate: '2020-08-31',
      requestedPickupDate: '2020-08-31',
      shipmentType: 'HHG',
      status: MOVE_STATUSES.SUBMITTED,
      updatedAt: '2020-09-02T21:08:38.392Z',
    },
  ],
  onDidMount: jest.fn(),
  showLoggedInUser: jest.fn(),
};

const testPPM = {
  advance_worksheet: {
    id: '00000000-0000-0000-0000-000000000000',
    service_member_id: '00000000-0000-0000-0000-000000000000',
    uploads: [],
  },
  created_at: '2021-04-07T16:44:03.946Z',
  destination_postal_code: '85309',
  has_additional_postal_code: false,
  has_pro_gear: 'NOT SURE',
  has_pro_gear_over_thousand: 'YES',
  has_requested_advance: false,
  has_sit: false,
  id: 'b3a8794b-0613-460d-9cac-092bbcf808bb',
  incentive_estimate_max: 2135347,
  incentive_estimate_min: 1931981,
  mileage: 757,
  move_id: '55a782e3-c4bb-4907-9f8d-b174c0a886f6',
  original_move_date: '2021-04-28',
  pickup_postal_code: '10002',
  planned_sit_max: 0,
  sit_max: 1104747,
  status: 'DRAFT',
  updated_at: '2021-04-07T17:05:15.522Z',
  weight_estimate: 20000,
};

const testPropsWithPPM = {
  ...testProps,
  currentMove: {
    ...testProps.currentMove,
    personally_procured_moves: [testPPM.id],
  },
  currentPPM: testPPM,
};

describe('Summary page', () => {
  it('does not validate the entitlement if the user does not have a PPM', () => {
    render(<Summary {...testProps} />);
    expect(validateEntitlement).not.toHaveBeenCalled();
  });

  it('validates the entitlement if the user has a PPM', () => {
    render(<Summary {...testPropsWithPPM} />);
    expect(validateEntitlement).toHaveBeenCalledWith(testProps.currentMove.id);
  });

  describe('if the user can add another shipment', () => {
    it('displays the Add Another Shipment section', () => {
      render(<Summary {...testProps} />);

      expect(screen.getByRole('link', { name: 'Add another shipment' })).toHaveAttribute(
        'href',
        '/moves/123/shipment-type',
      );
    });

    it('displays a button that opens a modal', () => {
      render(<Summary {...testProps} />);

      expect(
        screen.queryByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).not.toBeInTheDocument();

      expect(screen.getByTitle('Help with adding shipments')).toBeInTheDocument();
      userEvent.click(screen.getByTitle('Help with adding shipments'));

      expect(
        screen.getByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).toBeInTheDocument();
    });
  });

  describe('if the weight estimate is above the allotted entitlement', () => {
    it('displays the entitlement warning message', async () => {
      validateEntitlement.mockImplementation(() =>
        // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
        // eslint-disable-next-line prefer-promise-reject-errors
        Promise.reject({
          response: {
            status: 409,
            body: {
              message:
                'Your estimated weight of 20,000 lbs is above your weight entitlement of 14,000 lbs. \n You will only be paid for the weight you move up to your weight entitlement',
            },
          },
        }),
      );

      const { queryByText } = render(<Summary {...testPropsWithPPM} />);

      await waitFor(() => {
        expect(queryByText(/Your estimated weight is above your entitlement./)).toBeInTheDocument();
        expect(
          queryByText(/Your estimated weight of 20,000 lbs is above your weight entitlement of 14,000 lbs./),
        ).toBeInTheDocument();
      });
    });
  });

  afterEach(jest.clearAllMocks);
});

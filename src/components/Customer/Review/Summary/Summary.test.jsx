import React from 'react';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { Summary } from 'components/Customer/Review/Summary/Summary';
import { MOVE_STATUSES } from 'shared/constants';
import { renderWithRouterProp } from 'testUtils';
import { customerRoutes } from 'constants/routes';

const mockRouterConfig = { path: customerRoutes.MOVE_REVIEW_PATH, params: { moveId: '123' } };
const testProps = {
  serviceMember: {
    id: '666',
    current_location: {
      name: 'Test Duty Location',
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
      name: 'Test Duty Location',
      address: {
        postalCode: '123456',
      },
    },
    new_duty_location: {
      name: 'New Test Duty Location',
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
  updateShipmentList: jest.fn(),
  setMsg: jest.fn(),
};

describe('Summary page', () => {
  describe('if the user can add another shipment', () => {
    it('displays the Add Another Shipment section', () => {
      renderWithRouterProp(<Summary {...testProps} />, mockRouterConfig);

      expect(screen.getByRole('link', { name: 'Add another shipment' })).toHaveAttribute(
        'href',
        '/moves/123/shipment-type',
      );
    });

    it('displays a button that opens a modal', async () => {
      renderWithRouterProp(<Summary {...testProps} />, mockRouterConfig);

      expect(
        screen.queryByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).not.toBeInTheDocument();

      expect(screen.getByTitle('Help with adding shipments')).toBeInTheDocument();
      await userEvent.click(screen.getByTitle('Help with adding shipments'));

      expect(
        screen.getByRole('heading', { level: 3, name: 'Reasons you might need another shipment' }),
      ).toBeInTheDocument();
    });
  });
  afterEach(jest.clearAllMocks);
});

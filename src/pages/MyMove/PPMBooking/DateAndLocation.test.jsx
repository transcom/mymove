import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';

import DateAndLocation from 'pages/MyMove/PPMBooking/DateAndLocation';
import { customerRoutes, generalRoutes } from 'constants/routes';

const mockPush = jest.fn();
const mockMoveId = 'move123';
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: () => ({
    moveId: mockMoveId,
  }),
}));

const defaultProps = {
  serviceMember: {
    id: '8',
    residential_address: {
      postalCode: '20001',
    },
  },
  destinationDutyLocation: {
    address: {
      postalCode: '10002',
    },
  },
};

const fullShipmentProps = {
  mtoShipment: {
    id: '9',
    moveTaskOrderID: mockMoveId,
    ppmShipment: {
      pickupPostalCode: '20002',
      secondaryPickupPostalCode: '20003',
      destinationPostalCode: '20004',
      secondaryDestinationPostalCode: '20005',
      sitExpected: true,
      expectedDepartureDate: '2022-12-31',
    },
  },
};

describe('DateAndLocation component', () => {
  describe('creating a new PPM shipment', () => {
    it('renders the heading and empty form', () => {
      render(<DateAndLocation {...defaultProps} />);

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');
    });

    it('routes back to the new shipment type screen when back is clicked', async () => {
      const selectShipmentType = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, {
        moveId: mockMoveId,
      });

      render(<DateAndLocation {...defaultProps} />);

      const backButton = await screen.getByRole('button', { name: 'Back' });
      userEvent.click(backButton);

      expect(mockPush).toHaveBeenCalledWith(selectShipmentType);
    });
  });

  describe('editing an existing PPM shipment', () => {
    it('renders the heading and form with shipment values', async () => {
      render(<DateAndLocation {...defaultProps} {...fullShipmentProps} />);

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');

      const postalCodes = screen.getAllByLabelText('ZIP');
      const secondaryPostalCodes = screen.getAllByLabelText('Second ZIP');

      await waitFor(() => {
        expect(postalCodes[0]).toHaveValue('20002');
        expect(postalCodes[1]).toHaveValue('20004');
        expect(secondaryPostalCodes[0]).toHaveValue('20003');
        expect(secondaryPostalCodes[1]).toHaveValue('20005');
        expect(screen.getAllByLabelText('Yes')[2]).toBeChecked();
        expect(screen.getByLabelText('When do you plan to start moving your PPM?')).toHaveValue('31 Dec 2022');
      });
    });

    it('routes back to the home page screen when back is clicked', async () => {
      const selectShipmentType = generatePath(generalRoutes.HOME_PATH);

      render(<DateAndLocation {...defaultProps} {...fullShipmentProps} />);

      userEvent.click(screen.getByRole('button', { name: 'Back' }));

      expect(mockPush).toHaveBeenCalledWith(selectShipmentType);
    });
  });
});

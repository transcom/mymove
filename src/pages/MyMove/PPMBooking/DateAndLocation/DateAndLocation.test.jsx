import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';

import DateAndLocation from 'pages/MyMove/PPMBooking/DateAndLocation/DateAndLocation';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { createMTOShipment, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';

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

const mockNewShipmentId = 'newShipment123';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createMTOShipment: jest.fn(),
  patchMTOShipment: jest.fn(),
}));

jest.mock('utils/validation', () => ({
  ...jest.requireActual('utils/validation'),
  validatePostalCode: jest.fn(),
}));

const mockDispatch = jest.fn();

jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: jest.fn().mockImplementation(() => mockDispatch),
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
  postalCodeValidator: jest.fn(),
};

const fullShipmentProps = {
  ...defaultProps,
  mtoShipment: {
    id: '9',
    moveTaskOrderID: mockMoveId,
    ppmShipment: {
      id: '10',
      pickupPostalCode: '20002',
      secondaryPickupPostalCode: '20003',
      destinationPostalCode: '20004',
      secondaryDestinationPostalCode: '20005',
      sitExpected: true,
      expectedDepartureDate: '2022-12-31',
    },
    eTag: 'Za8lF',
  },
};

beforeEach(() => {
  jest.clearAllMocks();
});

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

    it('calls create shipment endpoint and formats required payload values', async () => {
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });

      render(<DateAndLocation {...defaultProps} />);

      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      userEvent.type(primaryPostalCodes[0], '10001');
      userEvent.type(primaryPostalCodes[1], '10002');

      userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '04 Jul 2022');

      userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: mockMoveId,
          shipmentType: 'PPM',
          ppmShipment: {
            pickupPostalCode: '10001',
            destinationPostalCode: '10002',
            hasSecondaryPickupPostalCode: false,
            secondaryPickupPostalCode: null,
            hasSecondaryDestinationPostalCode: false,
            secondaryDestinationPostalCode: null,
            sitExpected: false,
            expectedDepartureDate: '2022-07-04',
          },
        });

        expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: mockNewShipmentId }));
        expect(mockPush).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
            moveId: mockMoveId,
            mtoShipmentId: mockNewShipmentId,
          }),
        );
      });
    });

    it('displays an error alert when the create shipment fails', async () => {
      createMTOShipment.mockRejectedValueOnce('fatal error');

      render(<DateAndLocation {...defaultProps} />);

      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      userEvent.type(primaryPostalCodes[0], '10001');
      userEvent.type(primaryPostalCodes[1], '10002');

      userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '04 Jul 2022');

      userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: mockMoveId,
          shipmentType: 'PPM',
          ppmShipment: {
            pickupPostalCode: '10001',
            destinationPostalCode: '10002',
            hasSecondaryPickupPostalCode: false,
            secondaryPickupPostalCode: null,
            hasSecondaryDestinationPostalCode: false,
            secondaryDestinationPostalCode: null,
            sitExpected: false,
            expectedDepartureDate: '2022-07-04',
          },
        });

        expect(screen.getByText('There was an error attempting to create your shipment.')).toBeInTheDocument();
      });
    });

    it('calls create shipment endpoint and formats optional payload values', async () => {
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });

      render(<DateAndLocation {...defaultProps} />);

      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      userEvent.type(primaryPostalCodes[0], '10001');
      userEvent.type(primaryPostalCodes[1], '10002');

      const radioElements = screen.getAllByLabelText('Yes');
      userEvent.click(radioElements[0]);
      userEvent.click(radioElements[1]);

      const secondaryPostalCodes = screen.getAllByLabelText('Second ZIP');
      userEvent.type(secondaryPostalCodes[0], '10003');
      userEvent.type(secondaryPostalCodes[1], '10004');

      userEvent.click(radioElements[2]);

      userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '04 Jul 2022');

      userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: mockMoveId,
          shipmentType: 'PPM',
          ppmShipment: {
            pickupPostalCode: '10001',
            destinationPostalCode: '10002',
            hasSecondaryPickupPostalCode: true,
            secondaryPickupPostalCode: '10003',
            hasSecondaryDestinationPostalCode: true,
            secondaryDestinationPostalCode: '10004',
            sitExpected: true,
            expectedDepartureDate: '2022-07-04',
          },
        });

        expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: mockNewShipmentId }));
        expect(mockPush).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
            moveId: mockMoveId,
            mtoShipmentId: mockNewShipmentId,
          }),
        );
      });
    });
  });

  describe('editing an existing PPM shipment', () => {
    it('renders the heading and form with shipment values', async () => {
      render(<DateAndLocation {...fullShipmentProps} />);

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');

      const postalCodes = screen.getAllByLabelText('ZIP');
      const secondaryPostalCodes = screen.getAllByLabelText('Second ZIP');

      await waitFor(() => {
        expect(screen.getByLabelText('When do you plan to start moving your PPM?')).toHaveValue('31 Dec 2022');
      });

      expect(postalCodes[0]).toHaveValue('20002');
      expect(postalCodes[1]).toHaveValue('20004');
      expect(secondaryPostalCodes[0]).toHaveValue('20003');
      expect(secondaryPostalCodes[1]).toHaveValue('20005');
      expect(screen.getAllByLabelText('Yes')[2]).toBeChecked();
    });

    it('routes back to the home page screen when back is clicked', async () => {
      const selectShipmentType = generatePath(generalRoutes.HOME_PATH);

      render(<DateAndLocation {...defaultProps} {...fullShipmentProps} />);

      userEvent.click(screen.getByRole('button', { name: 'Back' }));

      expect(mockPush).toHaveBeenCalledWith(selectShipmentType);
    });

    it('displays an error alert when the update shipment fails', async () => {
      patchMTOShipment.mockRejectedValueOnce('fatal error');

      render(<DateAndLocation {...fullShipmentProps} />);

      userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalledWith(
          fullShipmentProps.mtoShipment.id,
          {
            id: fullShipmentProps.mtoShipment.id,
            moveTaskOrderID: mockMoveId,
            shipmentType: 'PPM',
            ppmShipment: {
              id: fullShipmentProps.mtoShipment.ppmShipment.id,
              pickupPostalCode: '20002',
              destinationPostalCode: '20004',
              hasSecondaryPickupPostalCode: true,
              secondaryPickupPostalCode: '20003',
              hasSecondaryDestinationPostalCode: true,
              secondaryDestinationPostalCode: '20005',
              sitExpected: true,
              expectedDepartureDate: '2022-12-31',
            },
          },
          fullShipmentProps.mtoShipment.eTag,
        );

        expect(screen.getByText('There was an error attempting to update your shipment.')).toBeInTheDocument();
      });
    });

    it('calls update shipment endpoint and formats optional payload values', async () => {
      patchMTOShipment.mockResolvedValueOnce({ id: fullShipmentProps.mtoShipment.id });

      render(<DateAndLocation {...fullShipmentProps} />);

      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      userEvent.clear(primaryPostalCodes[0]);
      userEvent.type(primaryPostalCodes[0], '10001');
      userEvent.clear(primaryPostalCodes[1]);
      userEvent.type(primaryPostalCodes[1], '10002');

      const secondaryPostalCodes = screen.getAllByLabelText('Second ZIP');
      userEvent.clear(secondaryPostalCodes[0]);
      userEvent.type(secondaryPostalCodes[0], '10003');
      userEvent.clear(secondaryPostalCodes[1]);
      userEvent.type(secondaryPostalCodes[1], '10004');

      const expectedDepartureDate = screen.getByLabelText('When do you plan to start moving your PPM?');
      userEvent.clear(expectedDepartureDate);
      userEvent.type(expectedDepartureDate, '04 Jul 2022');

      userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalledWith(
          fullShipmentProps.mtoShipment.id,
          {
            id: fullShipmentProps.mtoShipment.id,
            moveTaskOrderID: mockMoveId,
            shipmentType: 'PPM',
            ppmShipment: {
              id: fullShipmentProps.mtoShipment.ppmShipment.id,
              pickupPostalCode: '10001',
              destinationPostalCode: '10002',
              hasSecondaryPickupPostalCode: true,
              secondaryPickupPostalCode: '10003',
              hasSecondaryDestinationPostalCode: true,
              secondaryDestinationPostalCode: '10004',
              sitExpected: true,
              expectedDepartureDate: '2022-07-04',
            },
          },
          fullShipmentProps.mtoShipment.eTag,
        );

        expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: fullShipmentProps.mtoShipment.id }));
        expect(mockPush).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
            moveId: mockMoveId,
            mtoShipmentId: fullShipmentProps.mtoShipment.id,
          }),
        );
      });
    });
  });
});

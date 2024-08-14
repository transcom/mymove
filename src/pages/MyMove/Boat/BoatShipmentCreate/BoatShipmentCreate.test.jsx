import React from 'react';
import { waitFor, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';

import BoatShipmentCreate from 'pages/MyMove/Boat/BoatShipmentCreate/BoatShipmentCreate';
import { customerRoutes } from 'constants/routes';
import { createMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { renderWithRouter } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const mockNavigate = jest.fn();

const mockMoveId = 'move123';
const mockNewShipmentId = 'newShipment123';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  useParams: () => ({ moveId: mockMoveId }),
  useLocation: () => ({ search: '' }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createMTOShipment: jest.fn(),
  patchMTOShipment: jest.fn(),
  patchMove: jest.fn(),
  searchTransportationOffices: jest.fn(),
  getAllMoves: jest.fn(),
}));

jest.mock('utils/validation', () => ({
  ...jest.requireActual('utils/validation'),
  validatePostalCode: jest.fn(),
}));

const mockDispatch = jest.fn();
jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: () => mockDispatch,
}));

const serviceMember = {
  serviceMember: {
    id: '8',
    residential_address: {
      streetAddress1: '123 Any St',
      streetAddress2: '',
      city: 'Norfolk',
      state: 'VA',
      postalCode: '20001',
    },
  },
};

const defaultProps = {
  destinationDutyLocation: {
    address: {
      streetAddress1: '234 Any St',
      streetAddress2: '',
      city: 'Richmond',
      state: 'VA',
      postalCode: '10002',
    },
  },
  postalCodeValidator: jest.fn(),
  ...serviceMember,
};

beforeEach(() => {
  jest.clearAllMocks();
});

const renderBoatShipmentCreate = (props) => {
  renderWithRouter(<BoatShipmentCreate {...defaultProps} {...props} />, {
    path: customerRoutes.SHIPMENT_BOAT_CREATE_PATH,
    params: { moveId: 'move123' },
  });
};

describe('BoatShipmentCreate component', () => {
  describe('creating a new Boat shipment', () => {
    it('renders the heading and empty form', () => {
      renderBoatShipmentCreate();

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Boat details and measurements');
    });

    it('routes back to the new shipment type screen when back is clicked', async () => {
      renderBoatShipmentCreate();
      const selectShipmentType = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, {
        moveId: mockMoveId,
      });

      const backButton = await screen.getByRole('button', { name: 'Back' });
      await userEvent.click(backButton);

      expect(mockNavigate).toHaveBeenCalledWith(selectShipmentType);
    });

    it('calls create shipment endpoint and formats required payload values', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });

      renderBoatShipmentCreate();

      await act(async () => {
        await userEvent.type(screen.getByTestId('year'), '2022');
        await userEvent.type(screen.getByTestId('make'), 'Yamaha');
        await userEvent.type(screen.getByTestId('model'), 'SX210');
        await userEvent.type(screen.getByTestId('lengthFeet'), '21');
        await userEvent.type(screen.getByTestId('widthFeet'), '8');
        await userEvent.type(screen.getByTestId('heightFeet'), '7');
        await userEvent.click(screen.getByTestId('hasTrailerNo'));
      });

      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Boat Haul-Away (BHA)');

      await userEvent.click(screen.getByTestId('boatConfirmationContinue'));

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: mockMoveId,
          shipmentType: 'BOAT_HAUL_AWAY',
          boatShipment: {
            type: 'HAUL_AWAY',
            year: 2022,
            make: 'Yamaha',
            model: 'SX210',
            lengthInInches: 252,
            widthInInches: 96,
            heightInInches: 84,
            hasTrailer: false,
            isRoadworthy: null,
          },
          customerRemarks: undefined,
        });

        expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: mockNewShipmentId }));
        expect(mockNavigate).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_BOAT_LOCATION_INFO, {
            moveId: mockMoveId,
            mtoShipmentId: mockNewShipmentId,
          }),
        );
      });
    }, 10000);

    it('displays an error alert when the create shipment fails', async () => {
      createMTOShipment.mockRejectedValueOnce('fatal error');
      renderBoatShipmentCreate();

      await act(async () => {
        await userEvent.type(screen.getByTestId('year'), '2022');
        await userEvent.type(screen.getByTestId('make'), 'Yamaha');
        await userEvent.type(screen.getByTestId('model'), 'SX210');
        await userEvent.type(screen.getByTestId('lengthFeet'), '21');
        await userEvent.type(screen.getByTestId('widthFeet'), '8');
        await userEvent.type(screen.getByTestId('heightFeet'), '7');
        await userEvent.click(screen.getByTestId('hasTrailerNo'));
      });

      await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Boat Haul-Away (BHA)');

      await userEvent.click(screen.getByTestId('boatConfirmationContinue'));

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: mockMoveId,
          shipmentType: 'BOAT_HAUL_AWAY',
          boatShipment: {
            type: 'HAUL_AWAY',
            year: 2022,
            make: 'Yamaha',
            model: 'SX210',
            lengthInInches: 252,
            widthInInches: 96,
            heightInInches: 84,
            hasTrailer: false,
            isRoadworthy: null,
          },
          customerRemarks: undefined,
        });

        expect(screen.getByText('There was an error attempting to create your shipment.')).toBeInTheDocument();
      });
    });
  });
});

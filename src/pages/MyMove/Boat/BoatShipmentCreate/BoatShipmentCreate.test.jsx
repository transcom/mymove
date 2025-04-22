import React from 'react';
import { waitFor, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';

import BoatShipmentCreate from 'pages/MyMove/Boat/BoatShipmentCreate/BoatShipmentCreate';
import { customerRoutes } from 'constants/routes';
import { createMTOShipment, patchMTOShipment } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { renderWithRouter } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { MOVE_LOCKED_WARNING } from 'shared/constants';

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
  deleteMTOShipment: jest.fn(),
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
  id: '8',
  residential_address: {
    streetAddress1: '123 Any St',
    streetAddress2: '',
    city: 'Norfolk',
    state: 'VA',
    postalCode: '20001',
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
  serviceMember,
};

beforeEach(() => {
  jest.clearAllMocks();
});

const renderBoatShipmentCreate = async (props) => {
  await act(async () => {
    renderWithRouter(<BoatShipmentCreate {...defaultProps} {...props} />, {
      path: customerRoutes.SHIPMENT_BOAT_CREATE_PATH,
      params: { moveId: 'move123' },
    });
  });
};

describe('BoatShipmentCreate component', () => {
  describe('creating a new Boat shipment', () => {
    it('renders the heading and empty form', async () => {
      await renderBoatShipmentCreate();

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Boat details and measurements');
    });

    it('routes back to the new shipment type screen when back is clicked', async () => {
      await renderBoatShipmentCreate();
      const selectShipmentType = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, {
        moveId: mockMoveId,
      });

      const backButton = await screen.getByRole('button', { name: 'Back' });
      await act(async () => {
        await userEvent.click(backButton);
      });

      expect(mockNavigate).toHaveBeenCalledWith(selectShipmentType);
    });

    it('calls create shipment endpoint and formats required payload values', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });

      await renderBoatShipmentCreate();

      await act(async () => {
        await userEvent.type(screen.getByTestId('year'), '2022');
        await userEvent.type(screen.getByTestId('make'), 'Yamaha');
        await userEvent.type(screen.getByTestId('model'), 'SX210');
        await userEvent.type(screen.getByTestId('lengthFeet'), '21');
        await userEvent.type(screen.getByTestId('widthFeet'), '8');
        await userEvent.type(screen.getByTestId('heightFeet'), '7');
        await userEvent.click(screen.getByTestId('hasTrailerNo'));
        await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      });

      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Boat Haul-Away (BHA)');

      await act(async () => {
        await userEvent.click(screen.getByTestId('boatConfirmationContinue'));
      });

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

        expect(mockDispatch).toHaveBeenCalledWith(
          updateMTOShipment(expect.objectContaining({ id: mockNewShipmentId })),
        );
        expect(mockNavigate).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_BOAT_LOCATION_INFO, {
            moveId: mockMoveId,
            mtoShipmentId: mockNewShipmentId,
          }),
        );
      });
    });

    it('displays an error alert when the create shipment fails', async () => {
      createMTOShipment.mockRejectedValueOnce({
        response: { body: { invalidFields: { model: ['Some error message'] } } },
      });
      await renderBoatShipmentCreate();

      await act(async () => {
        await userEvent.type(screen.getByTestId('year'), '2022');
        await userEvent.type(screen.getByTestId('make'), 'Yamaha');
        await userEvent.type(screen.getByTestId('model'), 'SX210');
        await userEvent.type(screen.getByTestId('lengthFeet'), '21');
        await userEvent.type(screen.getByTestId('widthFeet'), '8');
        await userEvent.type(screen.getByTestId('heightFeet'), '7');
        await userEvent.click(screen.getByTestId('hasTrailerNo'));
        await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      });

      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Boat Haul-Away (BHA)');

      await act(async () => {
        await userEvent.click(screen.getByTestId('boatConfirmationContinue'));
      });

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

        expect(screen.getByText('Some error message')).toBeInTheDocument();
      });
    });

    it('disables the continue button when the move is locked by an office user', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });

      await renderBoatShipmentCreate({ isMoveLocked: true });

      await act(async () => {
        await userEvent.type(screen.getByTestId('year'), '2022');
        await userEvent.type(screen.getByTestId('make'), 'Yamaha');
        await userEvent.type(screen.getByTestId('model'), 'SX210');
        await userEvent.type(screen.getByTestId('lengthFeet'), '21');
        await userEvent.type(screen.getByTestId('widthFeet'), '8');
        await userEvent.type(screen.getByTestId('heightFeet'), '7');
        await userEvent.click(screen.getByTestId('hasTrailerNo'));
      });

      expect(screen.getByRole('button', { name: 'Continue' })).toBeDisabled();
      expect(screen.getByText(MOVE_LOCKED_WARNING)).toBeInTheDocument();
      expect(screen.getByText(MOVE_LOCKED_WARNING)).toBeVisible();
    });
  });

  describe('editing an existing Boat shipment', () => {
    const existingShipment = {
      id: 'existingShipment123',
      eTag: 'someETag',
      boatShipment: {
        id: 'boat123',
        type: 'TOW_AWAY',
        year: 2020,
        make: 'Sea Ray',
        model: 'Sundancer',
        lengthInInches: 240,
        widthInInches: 96,
        heightInInches: 84,
        hasTrailer: true,
        isRoadworthy: true,
      },
    };

    it('calls patch shipment endpoint and formats required payload values', async () => {
      patchMTOShipment.mockResolvedValueOnce({ id: existingShipment.id });

      await renderBoatShipmentCreate({ mtoShipment: existingShipment });

      await act(async () => {
        await userEvent.clear(screen.getByTestId('year'));
        await userEvent.type(screen.getByTestId('year'), '2021');
        await userEvent.clear(screen.getByTestId('make'));
        await userEvent.type(screen.getByTestId('make'), 'Bayliner');
        await userEvent.clear(screen.getByTestId('model'));
        await userEvent.type(screen.getByTestId('model'), 'Ciera');
        await userEvent.clear(screen.getByTestId('lengthFeet'));
        await userEvent.type(screen.getByTestId('lengthFeet'), '25');
        await userEvent.clear(screen.getByTestId('widthFeet'));
        await userEvent.type(screen.getByTestId('widthFeet'), '8');
        await userEvent.clear(screen.getByTestId('heightFeet'));
        await userEvent.type(screen.getByTestId('heightFeet'), '7');
        await userEvent.click(screen.getByTestId('hasTrailerYes'));
        await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      });

      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Boat Tow-Away (BTA)');

      await act(async () => {
        await userEvent.click(screen.getByTestId('boatConfirmationContinue'));
      });

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalledWith(
          existingShipment.id,
          {
            moveTaskOrderID: mockMoveId,
            shipmentType: 'BOAT_TOW_AWAY',
            boatShipment: {
              id: 'boat123',
              type: 'TOW_AWAY',
              year: 2021,
              make: 'Bayliner',
              model: 'Ciera',
              lengthInInches: 300,
              widthInInches: 96,
              heightInInches: 84,
              hasTrailer: true,
              isRoadworthy: true,
            },
            customerRemarks: undefined,
            id: 'existingShipment123',
          },
          'someETag',
        );

        expect(mockDispatch).toHaveBeenCalledWith(
          updateMTOShipment(expect.objectContaining({ id: existingShipment.id })),
        );
        expect(mockNavigate).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_BOAT_LOCATION_INFO, {
            moveId: mockMoveId,
            mtoShipmentId: existingShipment.id,
          }),
        );
      });
    });
  });
});

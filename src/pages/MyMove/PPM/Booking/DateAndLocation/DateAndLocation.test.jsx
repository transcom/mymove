import React from 'react';
import { render, waitFor, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath, MemoryRouter } from 'react-router';
import selectEvent from 'react-select-event';
import { act } from 'react-dom/test-utils';

import DateAndLocation from 'pages/MyMove/PPM/Booking/DateAndLocation/DateAndLocation';
import { customerRoutes, generalRoutes } from 'constants/routes';
import { createMTOShipment, patchMTOShipment, patchMove, searchTransportationOffices } from 'services/internalApi';
import { updateMTOShipment, updateMove } from 'store/entities/actions';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';

const mockPush = jest.fn();

const mockMoveId = 'move123';
const mockMove = {
  id: mockMoveId,
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};
const mockCloseoutId = '3210a533-19b8-4805-a564-7eb452afce10';

const mockCloseoutOffice = {
  address: {
    city: 'Test City',
    country: 'United States',
    id: 'a13806fc-0e7d-4dc3-91ca-b802d9da50f1',
    postalCode: '85309',
    state: 'AZ',
    streetAddress1: '7383 N Litchfield Rd',
    streetAddress2: 'Rm 1122',
  },
  created_at: '2018-05-28T14:27:39.198Z',
  gbloc: 'KKFA',
  id: mockCloseoutId,
  name: 'Tester',
  phone_lines: [],
  updated_at: '2018-05-28T14:27:39.198Z',
};

const mockSearchTransportationOffices = () => Promise.resolve([mockCloseoutOffice]);

jest.mock('components/LocationSearchBox/api', () => ({
  ShowAddress: jest.fn().mockImplementation(() =>
    Promise.resolve({
      city: 'Test City',
      country: 'United States',
      id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
      postalCode: '85309',
      state: 'AZ',
      streetAddress1: 'n/a',
    }),
  ),
}));

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
  patchMove: jest.fn(),
  searchTransportationOffices: jest.fn(),
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

const serviceMember = {
  serviceMember: {
    id: '8',
    residential_address: {
      postalCode: '20001',
    },
  },
};

const defaultProps = {
  destinationDutyLocation: {
    address: {
      postalCode: '10002',
    },
  },
  postalCodeValidator: jest.fn(),
  ...serviceMember,
};

const armyServiceMember = {
  ...defaultProps.serviceMember,
  affiliation: SERVICE_MEMBER_AGENCIES.ARMY,
};

const navyServiceMember = {
  ...defaultProps.serviceMember,
  affiliation: SERVICE_MEMBER_AGENCIES.NAVY,
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
      render(<DateAndLocation {...defaultProps} />, { wrapper: MemoryRouter });

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');
    });

    it('routes back to the new shipment type screen when back is clicked', async () => {
      render(<DateAndLocation {...defaultProps} />, { wrapper: MemoryRouter });

      const selectShipmentType = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, {
        moveId: mockMoveId,
      });

      const backButton = await screen.getByRole('button', { name: 'Back' });
      await userEvent.click(backButton);

      expect(mockPush).toHaveBeenCalledWith(selectShipmentType);
    });

    it('calls create shipment endpoint and formats required payload values', async () => {
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });

      render(<DateAndLocation {...defaultProps} />, { wrapper: MemoryRouter });

      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      await userEvent.type(primaryPostalCodes[0], '10001');
      await userEvent.type(primaryPostalCodes[1], '10002');

      await userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '04 Jul 2022');

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

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

      render(<DateAndLocation {...defaultProps} />, { wrapper: MemoryRouter });

      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      await userEvent.type(primaryPostalCodes[0], '10001');
      await userEvent.type(primaryPostalCodes[1], '10002');

      await userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '04 Jul 2022');

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

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

      render(<DateAndLocation {...defaultProps} />, { wrapper: MemoryRouter });

      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      await userEvent.type(primaryPostalCodes[0], '10001');
      await userEvent.type(primaryPostalCodes[1], '10002');

      const radioElements = screen.getAllByLabelText('Yes');
      await userEvent.click(radioElements[0]);
      await userEvent.click(radioElements[1]);

      const secondaryPostalCodes = screen.getAllByLabelText('Second ZIP');
      await userEvent.type(secondaryPostalCodes[0], '10003');
      await userEvent.type(secondaryPostalCodes[1], '10004');

      await userEvent.click(radioElements[2]);

      await userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '04 Jul 2022');

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

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

    it('calls patch move when there is a closeout office (Army/Air Force) and create shipment succeeds', async () => {
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });
      patchMove.mockResolvedValueOnce(mockMove);
      searchTransportationOffices.mockImplementation(mockSearchTransportationOffices);

      render(<DateAndLocation {...defaultProps} serviceMember={armyServiceMember} move={mockMove} />, {
        wrapper: MemoryRouter,
      });

      // Fill in form
      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      await userEvent.type(primaryPostalCodes[0], '10001');
      await userEvent.type(primaryPostalCodes[1], '10002');
      await userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '04 Jul 2022');

      // Set Closeout office
      const closeoutOfficeInput = await screen.getByLabelText('Which closeout office should review your PPM?');
      await fireEvent.change(closeoutOfficeInput, { target: { value: 'Tester' } });
      await act(() => selectEvent.select(closeoutOfficeInput, /Tester/));

      // Submit form
      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        // Shipment should get created
        expect(createMTOShipment).toHaveBeenCalledTimes(1);

        // Move patched with the closeout office
        expect(patchMove).toHaveBeenCalledTimes(1);
        expect(patchMove).toHaveBeenCalledWith(mockMove.id, { closeoutOfficeId: mockCloseoutId }, mockMove.eTag);

        // Redux updated with new shipment and updated move
        expect(mockDispatch).toHaveBeenCalledTimes(2);
        expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: mockNewShipmentId }));
        expect(mockDispatch).toHaveBeenCalledWith(updateMove(mockMove));

        // Finally, should get redirected to the estimated weight page
        expect(mockPush).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
            moveId: mockMoveId,
            mtoShipmentId: mockNewShipmentId,
          }),
        );
      });
    });

    it('does not call patch move when there is not a closeout office (not Army/Air Force)', async () => {
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });

      render(<DateAndLocation {...defaultProps} serviceMember={navyServiceMember} />, {
        wrapper: MemoryRouter,
      });

      // Fill in form
      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      await userEvent.type(primaryPostalCodes[0], '10001');
      await userEvent.type(primaryPostalCodes[1], '10002');
      await userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '04 Jul 2022');

      // Should not see closeout office field
      expect(screen.queryByLabelText('Which closeout office should review your PPM?')).not.toBeInTheDocument();

      // Submit form
      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        // Shipment should get created
        expect(createMTOShipment).toHaveBeenCalledTimes(1);

        // Should not try to patch the move
        expect(patchMove).toHaveBeenCalledTimes(0);

        // Redux updated with new shipment (and not a updated move)
        expect(mockDispatch).toHaveBeenCalledTimes(1);
        expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: mockNewShipmentId }));

        // Finally, should get redirected to the estimated weight page
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
      render(<DateAndLocation {...fullShipmentProps} />, { wrapper: MemoryRouter });

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
      render(<DateAndLocation {...defaultProps} {...fullShipmentProps} />, { wrapper: MemoryRouter });

      const selectShipmentType = generatePath(generalRoutes.HOME_PATH);

      await userEvent.click(screen.getByRole('button', { name: 'Back' }));

      expect(mockPush).toHaveBeenCalledWith(selectShipmentType);
    });

    it('displays an error alert when the update shipment fails', async () => {
      patchMTOShipment.mockRejectedValueOnce('fatal error');

      render(<DateAndLocation {...fullShipmentProps} />, { wrapper: MemoryRouter });

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

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

      render(<DateAndLocation {...fullShipmentProps} />, { wrapper: MemoryRouter });

      const primaryPostalCodes = screen.getAllByLabelText('ZIP');
      await userEvent.clear(primaryPostalCodes[0]);
      await userEvent.type(primaryPostalCodes[0], '10001');
      await userEvent.clear(primaryPostalCodes[1]);
      await userEvent.type(primaryPostalCodes[1], '10002');

      const secondaryPostalCodes = screen.getAllByLabelText('Second ZIP');
      await userEvent.clear(secondaryPostalCodes[0]);
      await userEvent.type(secondaryPostalCodes[0], '10003');
      await userEvent.clear(secondaryPostalCodes[1]);
      await userEvent.type(secondaryPostalCodes[1], '10004');

      const expectedDepartureDate = screen.getByLabelText('When do you plan to start moving your PPM?');
      await userEvent.clear(expectedDepartureDate);
      await userEvent.type(expectedDepartureDate, '04 Jul 2022');

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

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

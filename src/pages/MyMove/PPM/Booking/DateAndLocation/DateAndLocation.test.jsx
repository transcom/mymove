import React from 'react';
import { render, waitFor, screen, fireEvent, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';
import selectEvent from 'react-select-event';

import DateAndLocation from 'pages/MyMove/PPM/Booking/DateAndLocation/DateAndLocation';
import { customerRoutes } from 'constants/routes';
import { createMTOShipment, patchMTOShipment, patchMove, searchTransportationOffices } from 'services/internalApi';
import { updateMTOShipment, updateMove } from 'store/entities/actions';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { MockProviders } from 'testUtils';
import { PPM_TYPES } from 'shared/constants';

const mockNavigate = jest.fn();

const mockMoveId = 'move123';
const mockRoutingParams = { moveId: mockMoveId };
const mockNewShipmentId = 'newShipment123';

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

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
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
      pickupAddress: {
        streetAddress1: '234 Any St',
        streetAddress2: '',
        city: 'Richmond',
        state: 'VA',
        postalCode: '20002',
      },
      destinationAddress: {
        streetAddress1: '234 Any St',
        streetAddress2: '',
        city: 'Richmond',
        state: 'VA',
        postalCode: '20003',
      },
      secondaryPickupAddress: {
        streetAddress1: '234 Any St',
        streetAddress2: '',
        city: 'Richmond',
        state: 'VA',
        postalCode: '20004',
      },
      secondaryDestinationAddress: {
        streetAddress1: '234 Any St',
        streetAddress2: '',
        city: 'Richmond',
        state: 'VA',
        postalCode: '20005',
      },
      tertiaryPickupAddress: {
        streetAddress1: '234 Any St',
        streetAddress2: '',
        city: 'Richmond',
        state: 'VA',
        postalCode: '20006',
      },
      tertiaryDestinationAddress: {
        streetAddress1: '234 Any St',
        streetAddress2: '',
        city: 'Richmond',
        state: 'VA',
        postalCode: '20007',
      },
      sitExpected: true,
      expectedDepartureDate: '2022-12-31',
      hasTertiaryPickupAddress: true,
      hasTertiaryDestinationAddress: true,
    },
    eTag: 'Za8lF',
  },
};

beforeEach(() => {
  jest.clearAllMocks();
});

const renderDateAndLocation = (props) => {
  const mockRoutingOptions = {
    path: customerRoutes.SHIPMENT_SELECT_TYPE_PATH,
    params: mockRoutingParams,
  };
  render(
    <MockProviders {...mockRoutingOptions}>
      <DateAndLocation {...defaultProps} {...props} />
    </MockProviders>,
  );
};

describe('DateAndLocation component', () => {
  describe('creating a new PPM shipment', () => {
    it('renders the heading and empty form', () => {
      renderDateAndLocation();

      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');
    });

    it('routes back to the new shipment type screen when back is clicked', async () => {
      renderDateAndLocation();
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

      renderDateAndLocation();

      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current pickup address'));
      });

      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current delivery address'));
      });

      await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM?/), '04 Jul 2022');

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: mockMoveId,
          shipmentType: 'PPM',
          ppmShipment: {
            ppmType: PPM_TYPES.INCENTIVE_BASED,
            destinationAddress: {
              city: 'Richmond',
              postalCode: '10002',
              state: 'VA',
              streetAddress1: '234 Any St',
              streetAddress2: '',
            },
            pickupAddress: {
              city: 'Norfolk',
              postalCode: '20001',
              state: 'VA',
              streetAddress1: '123 Any St',
              streetAddress2: '',
            },
            hasSecondaryPickupAddress: false,
            hasSecondaryDestinationAddress: false,
            hasTertiaryPickupAddress: false,
            hasTertiaryDestinationAddress: false,
            sitExpected: false,
            expectedDepartureDate: '2022-07-04',
            isActualExpenseReimbursement: false,
          },
        });

        expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: mockNewShipmentId }));
        expect(mockNavigate).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
            moveId: mockMoveId,
            mtoShipmentId: mockNewShipmentId,
          }),
        );
      });
    }, 10000);

    it('displays an error alert when the create shipment fails', async () => {
      createMTOShipment.mockRejectedValueOnce('fatal error');
      renderDateAndLocation();

      // Fill in form
      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current pickup address'));
      });

      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current delivery address'));
      });

      await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM?/), '04 Jul 2022');

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: mockMoveId,
          shipmentType: 'PPM',
          ppmShipment: {
            ppmType: PPM_TYPES.INCENTIVE_BASED,
            destinationAddress: {
              city: 'Richmond',
              postalCode: '10002',
              state: 'VA',
              streetAddress1: '234 Any St',
              streetAddress2: '',
            },
            pickupAddress: {
              city: 'Norfolk',
              postalCode: '20001',
              state: 'VA',
              streetAddress1: '123 Any St',
              streetAddress2: '',
            },
            hasSecondaryPickupAddress: false,
            hasSecondaryDestinationAddress: false,
            hasTertiaryPickupAddress: false,
            hasTertiaryDestinationAddress: false,
            sitExpected: false,
            expectedDepartureDate: '2022-07-04',
            isActualExpenseReimbursement: false,
          },
        });

        expect(screen.getByText('There was an error attempting to create your shipment.')).toBeInTheDocument();
      });
    });

    it('calls create shipment endpoint and formats optional payload values', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });

      renderDateAndLocation();

      await act(async () => {
        await userEvent.click(document.querySelector('input[name="sitExpected"]'));
      });

      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current pickup address'));
      });

      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current delivery address'));
      });

      const radioElements = screen.getAllByLabelText('Yes');
      await userEvent.click(radioElements[2]);

      await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM?/), '04 Jul 2022');

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: mockMoveId,
          shipmentType: 'PPM',
          ppmShipment: {
            ppmType: PPM_TYPES.INCENTIVE_BASED,
            destinationAddress: {
              city: 'Richmond',
              postalCode: '10002',
              state: 'VA',
              streetAddress1: '234 Any St',
              streetAddress2: '',
            },
            pickupAddress: {
              city: 'Norfolk',
              postalCode: '20001',
              state: 'VA',
              streetAddress1: '123 Any St',
              streetAddress2: '',
            },
            hasSecondaryPickupAddress: false,
            hasSecondaryDestinationAddress: false,
            hasTertiaryPickupAddress: false,
            hasTertiaryDestinationAddress: false,
            sitExpected: true,
            expectedDepartureDate: '2022-07-04',
            isActualExpenseReimbursement: false,
          },
        });

        expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: mockNewShipmentId }));
        expect(mockNavigate).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
            moveId: mockMoveId,
            mtoShipmentId: mockNewShipmentId,
          }),
        );
      });
    }, 20000);

    // move and shipment successful patches are linked
    it.skip('calls patch move when there is a closeout office (Army/Air Force) and create shipment succeeds', async () => {
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });
      patchMove.mockResolvedValueOnce(mockMove);
      searchTransportationOffices.mockImplementation(mockSearchTransportationOffices);

      renderDateAndLocation({ serviceMember: armyServiceMember, move: mockMove });

      // Fill in form
      await act(async () => {
        await userEvent.type(
          document.querySelector('input[name="pickupAddress.address.streetAddress1"]'),
          '123 Any St',
        );
        await userEvent.click(screen.getByTestId('useCurrentResidence'));
      });

      await act(async () => {
        await userEvent.type(
          document.querySelector('input[name="destinationAddress.address.streetAddress1"]'),
          '123 Any St',
        );
        await userEvent.click(screen.getByTestId('useCurrentDestinationAddress'));
      });

      await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM?/), '04 Jul 2022');

      // Set Closeout office
      const closeoutOfficeInput = await screen.getByLabelText(/Which closeout office should review your PPM?/);
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
        expect(mockNavigate).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
            moveId: mockMoveId,
            mtoShipmentId: mockNewShipmentId,
          }),
        );
      });
    });

    it('does not call patch move when there is not a closeout office (not Army/Air Force)', async () => {
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });

      renderDateAndLocation({ serviceMember: navyServiceMember });

      // Fill in form
      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current pickup address'));
      });

      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current delivery address'));
      });

      await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM?/), '04 Jul 2022');

      // Should not see closeout office field
      expect(screen.queryByLabelText(/Which closeout office should review your PPM?/)).not.toBeInTheDocument();

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
        expect(mockNavigate).toHaveBeenCalledWith(
          generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
            moveId: mockMoveId,
            mtoShipmentId: mockNewShipmentId,
          }),
        );
      });
    });

    // move and shipment patches are linked
    it.skip('does not patch the move when create shipment fails', async () => {
      // createMTOShipment.mockRejectedValueOnce('fatal error');
      searchTransportationOffices.mockImplementation(mockSearchTransportationOffices);

      renderDateAndLocation({ serviceMember: armyServiceMember, move: mockMove });

      // Fill in form
      await act(async () => {
        await userEvent.type(
          document.querySelector('input[name="pickupAddress.address.streetAddress1"]'),
          '123 Any St',
        );
      });

      await act(async () => {
        await userEvent.type(document.querySelector('input[name="pickupAddress.address.city"]'), 'Norfolk');
      });

      await act(async () => {
        await userEvent.selectOptions(document.querySelector('select[name="pickupAddress.address.state"]'), 'VA');
      });

      await act(async () => {
        await userEvent.type(document.querySelector('input[name="pickupAddress.address.postalCode"]'), '10001');
      });

      await act(async () => {
        await userEvent.type(
          document.querySelector('input[name="destinationAddress.address.streetAddress1"]'),
          '123 Any St',
        );
      });

      await act(async () => {
        await userEvent.type(document.querySelector('input[name="destinationAddress.address.city"]'), 'Norfolk');
      });

      await act(async () => {
        await userEvent.selectOptions(document.querySelector('select[name="destinationAddress.address.state"]'), 'VA');
      });

      await act(async () => {
        await userEvent.type(document.querySelector('input[name="destinationAddress.address.postalCode"]'), '10002');
      });

      await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM?/), '04 Jul 2022');

      // Set Closeout office
      const closeoutOfficeInput = await screen.getByLabelText(/Which closeout office should review your PPM?/);
      await fireEvent.change(closeoutOfficeInput, { target: { value: 'Tester' } });
      await act(() => selectEvent.select(closeoutOfficeInput, /Tester/));

      // Submit form
      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        // Should have called called create shipment (set to fail above)
        expect(createMTOShipment).toHaveBeenCalledTimes(1);

        // Should not have patched the move since the create shipment failed
        expect(patchMove).not.toHaveBeenCalled();

        // Should not have done any redux updates
        expect(mockDispatch).not.toHaveBeenCalled();

        // No redirect should have happened
        expect(mockNavigate).not.toHaveBeenCalled();

        // Should show appropriate error message
        expect(screen.getByText('There was an error attempting to create your shipment.')).toBeInTheDocument();
      });
    }, 10000);

    // the shipment and move are patched at the same time so a successful shipment patch is a successful move patch
    it.skip('displays appropriate error when patch move fails after create shipment succeeds', async () => {
      createMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });
      patchMove.mockRejectedValueOnce('fatal error');
      searchTransportationOffices.mockImplementation(mockSearchTransportationOffices);

      renderDateAndLocation({ serviceMember: armyServiceMember, move: mockMove, closeoutOffice: mockCloseoutOffice });

      // Fill in form
      await act(async () => {
        await userEvent.type(
          document.querySelector('input[name="pickupAddress.address.streetAddress1"]'),
          '123 Any St',
        );
      });

      await act(async () => {
        await userEvent.type(document.querySelector('input[name="pickupAddress.address.city"]'), 'Norfolk');
      });

      await act(async () => {
        await userEvent.selectOptions(document.querySelector('select[name="pickupAddress.address.state"]'), 'VA');
      });

      await act(async () => {
        await userEvent.type(document.querySelector('input[name="pickupAddress.address.postalCode"]'), '10001');
      });

      await act(async () => {
        await userEvent.type(
          document.querySelector('input[name="destinationAddress.address.streetAddress1"]'),
          '123 Any St',
        );
      });

      await act(async () => {
        await userEvent.type(document.querySelector('input[name="destinationAddress.address.city"]'), 'Norfolk');
      });

      await act(async () => {
        await userEvent.selectOptions(document.querySelector('select[name="destinationAddress.address.state"]'), 'VA');
      });

      await act(async () => {
        await userEvent.type(document.querySelector('input[name="destinationAddress.address.postalCode"]'), '10002');
      });

      await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM?/), '04 Jul 2022');

      // Set Closeout office
      const closeoutOfficeInput = await screen.getByLabelText(/Which closeout office should review your PPM?/);
      await fireEvent.change(closeoutOfficeInput, { target: { value: 'Tester' } });
      await act(() => selectEvent.select(closeoutOfficeInput, /Tester/));

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        // Should have called both create shipment and patch move
        expect(createMTOShipment).toHaveBeenCalledTimes(1);
        expect(patchMove).toHaveBeenCalledTimes(1);

        // Should have only updated the shipment in redux
        expect(mockDispatch).toHaveBeenCalledTimes(1);
        expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: mockNewShipmentId }));

        // No redirect should have happened
        expect(mockNavigate).not.toHaveBeenCalled();

        // Should show appropriate error message
        expect(
          screen.getByText('There was an error attempting to create the move closeout office.'),
        ).toBeInTheDocument();
      });
    });
  });

  describe('editing an existing PPM shipment', () => {
    it('renders the heading and form with shipment values', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      renderDateAndLocation(fullShipmentProps);
      expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');

      const YesButtonSelectors = screen.getAllByLabelText('Yes');
      await userEvent.click(YesButtonSelectors[0]);
      await userEvent.click(YesButtonSelectors[1]);
      await userEvent.click(YesButtonSelectors[2]);
      await userEvent.click(YesButtonSelectors[3]);

      expect(screen.getAllByLabelText('Yes')[0]).toBeChecked();
      expect(screen.getAllByLabelText('Yes')[1]).toBeChecked();
      expect(screen.getAllByLabelText('Yes')[2]).toBeChecked();
      expect(screen.getAllByLabelText('Yes')[3]).toBeChecked();

      await waitFor(() => {
        expect(screen.getByLabelText(/When do you plan to start moving your PPM?/)).toHaveValue('31 Dec 2022');
      });

      expect(screen.getByText('Richmond, VA 20002 ()')).toBeInTheDocument();
      expect(screen.getByText('Richmond, VA 20004 ()')).toBeInTheDocument();
      expect(screen.getByText('Richmond, VA 20006 ()')).toBeInTheDocument();
      expect(screen.getByText('Richmond, VA 20003 ()')).toBeInTheDocument();
      expect(screen.getByText('Richmond, VA 20005 ()')).toBeInTheDocument();
      expect(screen.getByText('Richmond, VA 20007 ()')).toBeInTheDocument();
    });

    describe('editing an existing PPM shipment', () => {
      it('renders the heading and form with shipment values', async () => {
        isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
        renderDateAndLocation(fullShipmentProps);
        expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('PPM date & location');

        const YesButtonSelectors = screen.getAllByLabelText('Yes');
        await userEvent.click(YesButtonSelectors[0]);
        await userEvent.click(YesButtonSelectors[1]);
        await userEvent.click(YesButtonSelectors[2]);
        await userEvent.click(YesButtonSelectors[3]);

        expect(screen.getAllByLabelText('Yes')[0]).toBeChecked();
        expect(screen.getAllByLabelText('Yes')[1]).toBeChecked();
        expect(screen.getAllByLabelText('Yes')[2]).toBeChecked();
        expect(screen.getAllByLabelText('Yes')[3]).toBeChecked();

        await waitFor(() => {
          expect(screen.getByLabelText(/When do you plan to start moving your PPM?/)).toHaveValue('31 Dec 2022');
        });

        expect(screen.getByText('Richmond, VA 20002 ()')).toBeInTheDocument();
        expect(screen.getByText('Richmond, VA 20004 ()')).toBeInTheDocument();
        expect(screen.getByText('Richmond, VA 20006 ()')).toBeInTheDocument();
        expect(screen.getByText('Richmond, VA 20003 ()')).toBeInTheDocument();
        expect(screen.getByText('Richmond, VA 20005 ()')).toBeInTheDocument();
        expect(screen.getByText('Richmond, VA 20007 ()')).toBeInTheDocument();
      });

      it('routes back to the home page screen when back is clicked', async () => {
        isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));
        renderDateAndLocation(fullShipmentProps);

        const selectShipmentType = generatePath(customerRoutes.MOVE_HOME_PATH, { moveId: mockMoveId });

        await userEvent.click(screen.getByRole('button', { name: 'Back' }));

        expect(mockNavigate).toHaveBeenCalledWith(selectShipmentType);
      });

      it('displays an error alert when the update shipment fails', async () => {
        patchMTOShipment.mockRejectedValueOnce('fatal error');

        renderDateAndLocation(fullShipmentProps);

        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          expect(patchMTOShipment).toHaveBeenCalledWith(
            fullShipmentProps.mtoShipment.id,
            {
              id: fullShipmentProps.mtoShipment.id,
              moveTaskOrderID: mockMoveId,
              shipmentType: 'PPM',
              ppmShipment: {
                ppmType: PPM_TYPES.INCENTIVE_BASED,
                id: fullShipmentProps.mtoShipment.ppmShipment.id,
                pickupAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20002',
                },
                destinationAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20003',
                },
                secondaryPickupAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20004',
                },
                secondaryDestinationAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20005',
                },
                tertiaryPickupAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20006',
                },
                tertiaryDestinationAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20007',
                },
                hasSecondaryPickupAddress: true,
                hasSecondaryDestinationAddress: true,
                hasTertiaryPickupAddress: true,
                hasTertiaryDestinationAddress: true,
                sitExpected: true,
                expectedDepartureDate: '2022-12-31',
                isActualExpenseReimbursement: false,
              },
            },
            fullShipmentProps.mtoShipment.eTag,
          );

          expect(screen.getByText('There was an error attempting to update your shipment.')).toBeInTheDocument();
        });
      });

      it('displays an sepcific error message when the update shipment fails and return specific message', async () => {
        const mockSpecificMessage = 'The data entered no good.';
        patchMTOShipment.mockRejectedValueOnce({ response: { body: { detail: mockSpecificMessage, status: 400 } } });

        renderDateAndLocation(fullShipmentProps);

        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          expect(patchMTOShipment).toHaveBeenCalledWith(
            fullShipmentProps.mtoShipment.id,
            {
              id: fullShipmentProps.mtoShipment.id,
              moveTaskOrderID: mockMoveId,
              shipmentType: 'PPM',
              ppmShipment: {
                ppmType: PPM_TYPES.INCENTIVE_BASED,
                id: fullShipmentProps.mtoShipment.ppmShipment.id,
                pickupAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20002',
                },
                destinationAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20003',
                },
                secondaryPickupAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20004',
                },
                secondaryDestinationAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20005',
                },
                tertiaryPickupAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20006',
                },
                tertiaryDestinationAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20007',
                },
                hasSecondaryPickupAddress: true,
                hasSecondaryDestinationAddress: true,
                hasTertiaryPickupAddress: true,
                hasTertiaryDestinationAddress: true,
                sitExpected: true,
                expectedDepartureDate: '2022-12-31',
                isActualExpenseReimbursement: false,
              },
            },
            fullShipmentProps.mtoShipment.eTag,
          );

          expect(screen.getByText(mockSpecificMessage)).toBeInTheDocument();
        });
      });

      it('calls update shipment endpoint and formats optional payload values', async () => {
        patchMTOShipment.mockResolvedValueOnce({ id: fullShipmentProps.mtoShipment.id });

        renderDateAndLocation(fullShipmentProps);

        const expectedDepartureDate = screen.getByLabelText(/When do you plan to start moving your PPM?/);
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
                ppmType: PPM_TYPES.INCENTIVE_BASED,
                id: fullShipmentProps.mtoShipment.ppmShipment.id,
                pickupAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20002',
                },
                destinationAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20003',
                },
                secondaryPickupAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20004',
                },
                secondaryDestinationAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20005',
                },
                tertiaryPickupAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20006',
                },
                tertiaryDestinationAddress: {
                  streetAddress1: '234 Any St',
                  streetAddress2: '',
                  city: 'Richmond',
                  state: 'VA',
                  postalCode: '20007',
                },
                hasSecondaryPickupAddress: true,
                hasSecondaryDestinationAddress: true,
                hasTertiaryPickupAddress: true,
                hasTertiaryDestinationAddress: true,
                sitExpected: true,
                expectedDepartureDate: '2022-07-04',
                isActualExpenseReimbursement: false,
              },
            },
            fullShipmentProps.mtoShipment.eTag,
          );

          expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: fullShipmentProps.mtoShipment.id }));
          expect(mockNavigate).toHaveBeenCalledWith(
            generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
              moveId: mockMoveId,
              mtoShipmentId: fullShipmentProps.mtoShipment.id,
            }),
          );
        });
      });

      it('calls patch move when there is a closeout office (Army/Air Force) and update shipment succeeds', async () => {
        patchMTOShipment.mockResolvedValueOnce({ id: fullShipmentProps.mtoShipment.id });
        patchMove.mockResolvedValueOnce(mockMove);
        searchTransportationOffices.mockImplementation(mockSearchTransportationOffices);

        renderDateAndLocation({
          ...fullShipmentProps,
          serviceMember: armyServiceMember,
          move: {
            ...mockMove,
            closeoutOffice: mockCloseoutOffice,
          },
        });

        // Submit form
        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          // Shipment should get updated
          expect(patchMTOShipment).toHaveBeenCalledTimes(1);

          // Move patched with the closeout office
          expect(patchMove).toHaveBeenCalledTimes(1);
          expect(patchMove).toHaveBeenCalledWith(mockMove.id, { closeoutOfficeId: mockCloseoutId }, mockMove.eTag);

          // Redux updated with new shipment and updated move
          expect(mockDispatch).toHaveBeenCalledTimes(3);

          // Finally, should get redirected to the estimated weight page
          expect(mockNavigate).toHaveBeenCalledWith(
            generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
              moveId: mockMoveId,
              mtoShipmentId: fullShipmentProps.mtoShipment.id,
            }),
          );
        });
      });

      it('does not call patch move when there is not a closeout office (not Army/Air Force)', async () => {
        patchMTOShipment.mockResolvedValueOnce({ id: fullShipmentProps.mtoShipment.id });

        renderDateAndLocation({ ...fullShipmentProps, serviceMember: navyServiceMember, move: mockMove });

        // Submit form
        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          // Shipment should get updated
          expect(patchMTOShipment).toHaveBeenCalledTimes(1);

          // Should not try to patch the move
          expect(patchMove).toHaveBeenCalledTimes(0);

          // Redux updated with new shipment (and not a updated move)
          expect(mockDispatch).toHaveBeenCalledTimes(1);
          expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: fullShipmentProps.mtoShipment.id }));

          // Finally, should get redirected to the estimated weight page
          expect(mockNavigate).toHaveBeenCalledWith(
            generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, {
              moveId: mockMoveId,
              mtoShipmentId: fullShipmentProps.mtoShipment.id,
            }),
          );
        });
      });

      it('does not patch the move when patch shipment fails', async () => {
        patchMTOShipment.mockRejectedValueOnce('fatal error');

        renderDateAndLocation({
          ...fullShipmentProps,
          serviceMember: armyServiceMember,
          move: {
            ...mockMove,
            closeoutOffice: mockCloseoutOffice,
          },
        });

        // Submit form
        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          // Should have called called patch shipment (set to fail above)
          expect(patchMTOShipment).toHaveBeenCalledTimes(1);

          // Should not have patched the move since the patch shipment failed
          expect(patchMove).not.toHaveBeenCalled();

          // Should not have done any redux updates
          expect(mockDispatch).not.toHaveBeenCalled();

          // No redirect should have happened
          expect(mockNavigate).not.toHaveBeenCalled();

          // Should show appropriate error message
          expect(screen.getByText('There was an error attempting to update your shipment.')).toBeInTheDocument();
        });
      });

      it('displays appropriate error when patch move fails after patch shipment succeeds', async () => {
        patchMTOShipment.mockResolvedValueOnce({ id: mockNewShipmentId });
        patchMove.mockRejectedValueOnce('fatal error');
        searchTransportationOffices.mockImplementation(mockSearchTransportationOffices);

        renderDateAndLocation({
          ...fullShipmentProps,
          serviceMember: armyServiceMember,
          move: {
            ...mockMove,
            closeoutOffice: mockCloseoutOffice,
          },
        });

        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          // Should have called both patch shipment and patch move
          expect(patchMTOShipment).toHaveBeenCalledTimes(1);
          expect(patchMove).toHaveBeenCalledTimes(1);

          // Should have only updated the shipment in redux
          expect(mockDispatch).toHaveBeenCalled();
          expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment({ id: mockNewShipmentId }));

          // No redirect should have happened
          expect(mockNavigate).not.toHaveBeenCalled();

          // Should show appropriate error message
          expect(
            screen.getByText('There was an error attempting to update the move closeout office.'),
          ).toBeInTheDocument();
        });
      });
    });
  });
});

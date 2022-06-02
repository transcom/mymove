import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';
import { v4 as uuidv4 } from 'uuid';

import Advance from './Advance';

import { customerRoutes } from 'constants/routes';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { updateMTOShipment } from 'store/entities/actions';
import { setFlashMessage } from 'store/flash/actions';
import { selectCurrentOrders, selectMTOShipmentById } from 'store/entities/selectors';
import { MockProviders } from 'testUtils';
import { ORDERS_TYPE } from 'constants/orders';

const mockPush = jest.fn();

const mockMoveId = uuidv4();
const mockMTOShipmentId = uuidv4();

const reviewPath = generatePath(customerRoutes.MOVE_REVIEW_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

const estimatedIncentivePath = generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

const mockShipmentETag = Buffer.from(new Date()).toString('base64');
const mockPPMShipmentETag = Buffer.from(new Date()).toString('base64');

const mockMTOShipment = {
  id: mockMTOShipmentId,
  moveTaskOrderID: mockMoveId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: uuidv4(),
    pickupPostalCode: '20002',
    destinationPostalCode: '20004',
    sitExpected: false,
    expectedDepartureDate: '2022-12-31',
    eTag: mockPPMShipmentETag,
    estimatedIncentive: 1000000,
    estimatedWeight: 4000,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
    hasRequestedAdvance: null,
  },
  eTag: mockShipmentETag,
};

const mockOrders = {
  orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
};

const mockOrdersRetiree = {
  orders_type: ORDERS_TYPE.RETIREMENT,
};

const mockOrdersSeparatee = {
  orders_type: ORDERS_TYPE.SEPARATION,
};

const mockMTOShipmentWithAdvance = {
  ...mockMTOShipment,
  ppmShipment: {
    ...mockMTOShipment.ppmShipment,
    hasRequestedAdvance: true,
    advanceAmountRequested: 40000,
    eTag: mockPPMShipmentETag,
  },
  eTag: mockShipmentETag,
};

const mockDispatch = jest.fn();

jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: jest.fn().mockImplementation(() => mockDispatch),
}));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: () => ({
    moveId: mockMoveId,
    mtoShipmentId: mockMTOShipmentId,
  }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getResponseError: jest.fn(),
  patchMTOShipment: jest.fn(),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn().mockImplementation(() => mockMTOShipment),
  selectCurrentOrders: jest.fn().mockImplementation(() => mockOrders),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('Advance page', () => {
  it('renders the heading and empty form', () => {
    render(<Advance />, { wrapper: MockProviders });

    expect(screen.getByRole('heading', { level: 1 }).textContent).toMatchInlineSnapshot(`"Advances"`);

    const requestAdvanceYesInput = screen.getByRole('radio', { name: /yes/i });
    expect(requestAdvanceYesInput).toBeInstanceOf(HTMLInputElement);
    expect(requestAdvanceYesInput.checked).toBe(false);

    const requestAdvanceNoInput = screen.getByRole('radio', { name: /no/i });
    expect(requestAdvanceNoInput).toBeInstanceOf(HTMLInputElement);
    expect(requestAdvanceNoInput.checked).toBe(true);

    expect(screen.queryByLabelText('Amount requested')).not.toBeInTheDocument();

    const backButton = screen.getByRole('button', { name: /back/i });
    expect(backButton).toBeInTheDocument();
    expect(backButton).not.toBeDisabled();

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    expect(saveButton).toBeInTheDocument();
    expect(saveButton).not.toBeDisabled();
  });

  it.each([
    [mockMTOShipment, undefined],
    [undefined, mockOrders],
    [undefined, undefined],
  ])('renders the loading placeholder when mtoShipment or orders are missing', async (loadedShipment, loadedOrders) => {
    selectMTOShipmentById.mockImplementationOnce(() => loadedShipment);
    selectCurrentOrders.mockImplementationOnce(() => loadedOrders);

    render(<Advance />, { wrapper: MockProviders });

    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
  });

  it.each([[mockMTOShipment], [mockMTOShipmentWithAdvance]])(
    'renders the form with and without previously filled in amount requested for an advance',
    async (preExistingShipment) => {
      selectMTOShipmentById.mockImplementationOnce(() => preExistingShipment);

      render(<Advance />, { wrapper: MockProviders });

      const hasRequestedAdvanceYesInput = screen.getByRole('radio', { name: /yes/i });
      const hasRequestedAdvanceNoInput = screen.getByRole('radio', { name: /no/i });

      if (preExistingShipment.ppmShipment.hasRequestedAdvance) {
        expect(hasRequestedAdvanceYesInput.checked).toBe(true);
        expect(hasRequestedAdvanceNoInput.checked).toBe(false);
        await waitFor(() => {
          expect(screen.getByLabelText('Amount requested').value).toBe('400');
        });
      } else {
        expect(hasRequestedAdvanceYesInput.checked).toBe(false);
        expect(hasRequestedAdvanceNoInput.checked).toBe(true);
        expect(screen.queryByLabelText('Amount requested')).not.toBeInTheDocument();
      }
    },
  );

  it.each([[mockOrders], [mockOrdersRetiree], [mockOrdersSeparatee]])(
    'displays info alert guidance for advance to retirees and separatees',
    async (orders) => {
      selectCurrentOrders.mockImplementation(() => orders);

      render(<Advance />, { wrapper: MockProviders });

      const nonPCSOrdersAdvanceText =
        'People leaving the military may not be eligible to receive an advance, based on individual service policies. Your counselor can give you more information after you make your request.';

      if (orders.orders_type === ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION) {
        expect(screen.queryByText(nonPCSOrdersAdvanceText)).not.toBeInTheDocument();
      } else {
        expect(screen.getByText(nonPCSOrdersAdvanceText)).toBeInTheDocument();
      }
    },
  );

  it('can toggle optional fields', async () => {
    render(<Advance />, { wrapper: MockProviders });

    const hasRequestedAdvanceYesInput = screen.getByRole('radio', { name: /yes/i });
    userEvent.click(hasRequestedAdvanceYesInput);

    const advanceInput = await screen.findByLabelText('Amount requested');
    expect(advanceInput).toBeInstanceOf(HTMLInputElement);

    const hasRequestedAdvanceNoInput = screen.getByRole('radio', { name: /no/i });
    userEvent.click(hasRequestedAdvanceNoInput);

    await waitFor(() => {
      expect(screen.queryByLabelText('Amount requested')).not.toBeInTheDocument();
    });
  });

  it('routes back to the previous page when the back button is clicked', () => {
    render(<Advance />, { wrapper: MockProviders });

    const backButton = screen.getByRole('button', { name: /back/i });

    userEvent.click(backButton);

    expect(mockPush).toHaveBeenCalledWith(estimatedIncentivePath);
  });

  it('calls the patch shipment endpoint when save & continue is clicked', async () => {
    patchMTOShipment.mockResolvedValueOnce({ id: mockMTOShipment.id });

    render(<Advance />, { wrapper: MockProviders });

    const advance = 4000;
    const hasRequestedAdvanceYesInput = screen.getByRole('radio', { name: /yes/i });
    await userEvent.click(hasRequestedAdvanceYesInput);

    const advanceInput = screen.getByLabelText('Amount requested');
    userEvent.type(advanceInput, String(advance));

    const agreeToTerms = screen.getByLabelText(/I acknowledge/i);
    userEvent.click(agreeToTerms);

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    expect(saveButton).not.toBeDisabled();
    userEvent.click(saveButton);

    const expectedPayload = {
      shipmentType: mockMTOShipment.shipmentType,
      ppmShipment: {
        advance: 400000,
        advanceRequested: true,
        id: mockMTOShipment.ppmShipment.id,
      },
    };

    await waitFor(() =>
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, expectedPayload, mockMTOShipment.eTag),
    );
  });

  it('passes null for advance amount if advance requested is false', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithAdvance);
    patchMTOShipment.mockResolvedValueOnce({ id: mockMTOShipment.id });

    render(<Advance />, { wrapper: MockProviders });

    const hasRequestedAdvanceYesInput = screen.getByRole('radio', { name: /yes/i });
    const hasRequestedAdvanceNoInput = screen.getByRole('radio', { name: /no/i });

    expect(hasRequestedAdvanceYesInput.checked).toBe(true);
    expect(hasRequestedAdvanceNoInput.checked).toBe(false);

    await userEvent.click(hasRequestedAdvanceNoInput);

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    expect(saveButton).not.toBeDisabled();
    userEvent.click(saveButton);

    const expectedPayload = {
      shipmentType: mockMTOShipment.shipmentType,
      ppmShipment: {
        advance: null,
        advanceRequested: false,
        id: mockMTOShipment.ppmShipment.id,
      },
    };

    await waitFor(() =>
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, expectedPayload, mockMTOShipment.eTag),
    );
  });

  it('updates the state if shipment patch is successful', async () => {
    patchMTOShipment.mockResolvedValue(mockMTOShipment);

    render(<Advance />, { wrapper: MockProviders });

    const advance = 4000;
    const hasRequestedAdvanceYesInput = screen.getByRole('radio', { name: /yes/i });
    userEvent.click(hasRequestedAdvanceYesInput);

    const agreeToTerms = screen.getByLabelText(/I acknowledge/i);
    userEvent.click(agreeToTerms);

    const advanceInput = screen.getByLabelText('Amount requested');
    userEvent.type(advanceInput, String(advance));

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    userEvent.click(saveButton);

    await waitFor(() => expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment(mockMTOShipment)));
  });

  it('routes to the review page when the user clicks save & continue', async () => {
    patchMTOShipment.mockResolvedValue({});

    render(<Advance />, { wrapper: MockProviders });
    const hasRequestedAdvanceYesInput = screen.getByRole('radio', { name: /yes/i });
    userEvent.click(hasRequestedAdvanceYesInput);

    const agreeToTerms = screen.getByLabelText(/I acknowledge/i);
    userEvent.click(agreeToTerms);

    const advanceInput = screen.getByLabelText('Amount requested');
    userEvent.type(advanceInput, '4000');

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    userEvent.click(saveButton);

    await waitFor(() => expect(mockPush).toHaveBeenCalledWith(reviewPath));
    expect(mockDispatch).toHaveBeenCalledWith(
      setFlashMessage(
        'PPM_ONBOARDING_SUBMIT_SUCCESS',
        'success',
        'Review your info and submit your move request now, or come back and finish later.',
        'Details saved',
      ),
    );
  });

  it('displays an error message if the update fails', async () => {
    const mockErrorMsg = 'Invalid shipment ID';

    patchMTOShipment.mockRejectedValue({ id: '11' });
    getResponseError.mockReturnValue(mockErrorMsg);

    render(<Advance />, { wrapper: MockProviders });
    const hasRequestedAdvanceYesInput = screen.getByRole('radio', { name: /yes/i });
    userEvent.click(hasRequestedAdvanceYesInput);

    const advanceInput = screen.getByLabelText('Amount requested');
    userEvent.type(advanceInput, '4000');
    const agreeToTerms = screen.getByLabelText(/I acknowledge/i);
    userEvent.click(agreeToTerms);
    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    userEvent.click(saveButton);

    expect(await screen.findByText(mockErrorMsg)).toBeInTheDocument();
  });
});

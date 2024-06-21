import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 as uuidv4 } from 'uuid';

import EstimatedWeightsProGear from './EstimatedWeightsProGear';

import { customerRoutes } from 'constants/routes';
import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { updateMTOShipment } from 'store/entities/actions';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { MockProviders } from 'testUtils';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

const mockMoveId = uuidv4();
const mockMTOShipmentId = uuidv4();
const mockRoutingConfig = {
  path: customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH,
  params: {
    moveId: mockMoveId,
    mtoShipmentId: mockMTOShipmentId,
  },
};
const shipmentEditPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});
const estimatedIncentivePath = generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

const mockOrders = {
  has_dependents: false,
  authorizedWeight: 5000,
  entitlement: {
    proGear: 2000,
    proGearSpouse: 500,
  },
};

const mockServiceMember = {
  id: uuidv4(),
};

const mockMTOShipment = {
  id: mockMTOShipmentId,
  moveTaskOrderID: mockMoveId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: uuidv4(),
    sitExpected: false,
    expectedDepartureDate: '2022-12-31',
    eTag: window.btoa(new Date()),
  },
  eTag: window.btoa(new Date()),
};

const mockPreExistingShipment = {
  ...mockMTOShipment,
  ppmShipment: {
    ...mockMTOShipment.ppmShipment,
    estimatedWeight: 4000,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
    eTag: window.btoa(new Date()),
  },
  eTag: window.btoa(new Date()),
};

const mockPreExistingShipmentWithProGear = {
  ...mockPreExistingShipment,
  ppmShipment: {
    ...mockPreExistingShipment.ppmShipment,
    hasProGear: true,
    proGearWeight: 1000,
    spouseProGearWeight: 100,
    eTag: window.btoa(new Date()),
  },
  eTag: window.btoa(new Date()),
};

const mockDispatch = jest.fn();

jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: jest.fn().mockImplementation(() => mockDispatch),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getResponseError: jest.fn(),
  patchMTOShipment: jest.fn(),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectCurrentOrders: jest.fn().mockImplementation(() => mockOrders),
  selectMTOShipmentById: jest.fn().mockImplementation(() => mockMTOShipment),
  selectServiceMemberFromLoggedInUser: jest.fn().mockImplementation(() => mockServiceMember),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const renderEstimatedWeightsProGear = (props) => {
  return render(
    <MockProviders {...mockRoutingConfig}>
      <EstimatedWeightsProGear {...props} />
    </MockProviders>,
  );
};

describe('EstimatedWeightsProGear page', () => {
  it('renders the heading and empty form when weight info has not been entered', () => {
    renderEstimatedWeightsProGear();

    expect(screen.getByRole('heading', { level: 1 }).textContent).toMatchInlineSnapshot(`"Estimated weight"`);

    const estimatedWeightInput = screen.getByLabelText(/estimated weight of this ppm shipment/i);
    expect(estimatedWeightInput).toBeInTheDocument(HTMLInputElement);
    expect(estimatedWeightInput.value).toBe('');

    const hasProGearYesInput = screen.getByRole('radio', { name: /yes/i });
    expect(hasProGearYesInput).toBeInstanceOf(HTMLInputElement);
    expect(hasProGearYesInput.checked).toBe(false);

    const hasProGearNoInput = screen.getByRole('radio', { name: /no/i });
    expect(hasProGearNoInput).toBeInstanceOf(HTMLInputElement);
    expect(hasProGearNoInput.checked).toBe(true);

    expect(screen.queryByLabelText(/estimated weight of your pro-gear/i)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/estimated weight of your spouse’s pro-gear/i)).not.toBeInTheDocument();

    const backButton = screen.getByRole('button', { name: /back/i });
    expect(backButton).toBeInTheDocument();
    expect(backButton).not.toBeDisabled();

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    expect(saveButton).toBeInTheDocument();
    expect(saveButton).not.toBeDisabled();
  });

  it.each([[mockPreExistingShipment], [mockPreExistingShipmentWithProGear]])(
    'renders the form pre-filled when weight info has been entered previously',
    async (preExistingShipment) => {
      selectMTOShipmentById.mockImplementationOnce(() => preExistingShipment);

      renderEstimatedWeightsProGear();

      await waitFor(() => {
        expect(screen.getByLabelText(/estimated weight of this ppm shipment/i).value).toBe('4,000');
      });

      const hasProGearYesInput = screen.getByRole('radio', { name: /yes/i });
      const hasProGearNoInput = screen.getByRole('radio', { name: /no/i });

      if (preExistingShipment.ppmShipment.hasProGear) {
        expect(hasProGearYesInput.checked).toBe(true);
        expect(hasProGearNoInput.checked).toBe(false);

        const proGearWeightInput = screen.getByLabelText(/estimated weight of your pro-gear/i);
        expect(proGearWeightInput.value).toBe('1,000');

        const spouseProGearWeightInput = screen.getByLabelText(/estimated weight of your spouse’s pro-gear/i);
        expect(spouseProGearWeightInput.value).toBe('100');
      } else {
        expect(hasProGearYesInput.checked).toBe(false);
        expect(hasProGearNoInput.checked).toBe(true);
        expect(screen.queryByLabelText(/estimated weight of your pro-gear/i)).not.toBeInTheDocument();
        expect(screen.queryByLabelText(/estimated weight of your spouse’s pro-gear/i)).not.toBeInTheDocument();
      }
    },
  );

  it('can toggle optional fields', async () => {
    renderEstimatedWeightsProGear();

    const hasProGearYesInput = screen.getByRole('radio', { name: /yes/i });
    await userEvent.click(hasProGearYesInput);

    const proGearWeightInput = await screen.findByLabelText(/estimated weight of your pro-gear/i);
    expect(proGearWeightInput).toBeInstanceOf(HTMLInputElement);

    const spouseProGearWeightInput = screen.getByLabelText(/estimated weight of your spouse’s pro-gear/i);
    expect(spouseProGearWeightInput).toBeInstanceOf(HTMLInputElement);

    const hasProGearNoInput = screen.getByRole('radio', { name: /no/i });
    await userEvent.click(hasProGearNoInput);

    await waitFor(() => {
      expect(screen.queryByLabelText(/estimated weight of your pro-gear/i)).not.toBeInTheDocument();
    });

    expect(screen.queryByLabelText(/estimated weight of your spouse’s pro-gear/i)).not.toBeInTheDocument();
  });

  it('routes back to the previous page when the back button is clicked', async () => {
    renderEstimatedWeightsProGear();

    const backButton = screen.getByRole('button', { name: /back/i });

    await userEvent.click(backButton);

    expect(mockNavigate).toHaveBeenCalledWith(shipmentEditPath);
  });

  it('calls the patch shipment endpoint when save & continue is clicked', async () => {
    patchMTOShipment.mockResolvedValue();

    renderEstimatedWeightsProGear();

    const estimatedWeight = 4000;

    const estimatedWeightInput = screen.getByLabelText(/estimated weight of this ppm shipment/i);
    await userEvent.type(estimatedWeightInput, String(estimatedWeight));

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    expect(saveButton).not.toBeDisabled();
    await userEvent.click(saveButton);

    const expectedPayload = {
      shipmentType: mockMTOShipment.shipmentType,
      ppmShipment: {
        id: mockMTOShipment.ppmShipment.id,
        estimatedWeight,
        hasProGear: false,
        proGearWeight: null,
        spouseProGearWeight: null,
      },
    };

    await waitFor(() =>
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, expectedPayload, mockMTOShipment.eTag),
    );
  });

  it('calls the patch shipment endpoint with optional values when save & continue is clicked', async () => {
    patchMTOShipment.mockResolvedValue();

    renderEstimatedWeightsProGear();

    const estimatedWeight = 4000;

    const estimatedWeightInput = screen.getByLabelText(/estimated weight of this ppm shipment/i);
    await userEvent.type(estimatedWeightInput, String(estimatedWeight));

    const hasProGearYesInput = screen.getByRole('radio', { name: /yes/i });
    await userEvent.click(hasProGearYesInput);

    const proGearWeight = 1000;

    const proGearWeightInput = await screen.findByLabelText(/estimated weight of your pro-gear/i);
    expect(proGearWeightInput).toBeInstanceOf(HTMLInputElement);

    await userEvent.type(proGearWeightInput, String(proGearWeight));

    const spouseProGearWeight = 100;

    const spouseProGearWeightInput = screen.getByLabelText(/estimated weight of your spouse’s pro-gear/i);
    expect(spouseProGearWeightInput).toBeInstanceOf(HTMLInputElement);

    await userEvent.type(spouseProGearWeightInput, String(spouseProGearWeight));

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    expect(saveButton).not.toBeDisabled();
    await userEvent.click(saveButton);

    const expectedPayload = {
      shipmentType: mockMTOShipment.shipmentType,
      ppmShipment: {
        id: mockMTOShipment.ppmShipment.id,
        estimatedWeight,
        hasProGear: true,
        proGearWeight,
        spouseProGearWeight,
      },
    };

    await waitFor(() => {
      // screen.debug();
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, expectedPayload, mockMTOShipment.eTag);
    });
  });

  it('updates the state if shipment patch is successful', async () => {
    patchMTOShipment.mockResolvedValue(mockPreExistingShipment);

    renderEstimatedWeightsProGear();

    const estimatedWeight = 4000;

    const estimatedWeightInput = screen.getByLabelText(/estimated weight of this ppm shipment/i);
    await userEvent.type(estimatedWeightInput, String(estimatedWeight));

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    await userEvent.click(saveButton);

    await waitFor(() => expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment(mockPreExistingShipment)));
  });

  it('routes to the estimated incentive page when the user clicks save & continue', async () => {
    patchMTOShipment.mockResolvedValue({});

    renderEstimatedWeightsProGear();

    const estimatedWeightInput = screen.getByLabelText(/estimated weight of this ppm shipment/i);
    await userEvent.type(estimatedWeightInput, '4000');

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    await userEvent.click(saveButton);

    await waitFor(() => expect(mockNavigate).toHaveBeenCalledWith(estimatedIncentivePath));
  });

  it('displays an error message if the update fails', async () => {
    const mockErrorMsg = 'Invalid shipment ID';

    patchMTOShipment.mockRejectedValue({});
    getResponseError.mockReturnValue(mockErrorMsg);

    renderEstimatedWeightsProGear();

    const estimatedWeightInput = screen.getByLabelText(/estimated weight of this ppm shipment/i);
    await userEvent.type(estimatedWeightInput, '4000');

    const saveButton = screen.getByRole('button', { name: /save & continue/i });
    await userEvent.click(saveButton);

    expect(await screen.findByText(mockErrorMsg)).toBeInTheDocument();
  });
});

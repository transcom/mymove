import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 as uuidv4 } from 'uuid';

import BoatShipmentLocationInfo from './BoatShipmentLocationInfo';

import { customerRoutes } from 'constants/routes';
import { patchMTOShipment } from 'services/internalApi';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { MockProviders } from 'testUtils';

const mockNavigate = jest.fn();
const mockMoveId = uuidv4();
const mockMTOShipmentId = uuidv4();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  useLocation: () => ({ search: '' }),
  useParams: () => ({ moveId: mockMoveId, mtoShipmentId: mockMTOShipmentId }),
}));

const shipmentEditPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

const mockRoutingConfig = {
  path: customerRoutes.SHIPMENT_BOAT_LOCATION_INFO,
  params: {
    moveId: mockMoveId,
    mtoShipmentId: mockMTOShipmentId,
  },
};

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
  shipmentType: SHIPMENT_TYPES.BOAT_HAUL_AWAY,
  boatShipment: {
    id: uuidv4(),
    type: 'HAUL_AWAY',
    year: 2022,
    make: 'Yamaha',
    model: 'SX210',
    lengthInInches: 252,
    widthInInches: 96,
    heightInInches: 84,
    hasTrailer: false,
    isRoadworthy: null,
    eTag: window.btoa(new Date()),
  },
  eTag: window.btoa(new Date()),
  createdAt: '2021-06-11T18:12:11.918Z',
  customerRemarks: 'mock remarks',
  requestedPickupDate: '2021-08-01',
  requestedDeliveryDate: '2021-08-11',
  pickupAddress: {
    id: uuidv4(),
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    id: uuidv4(),
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

const mockPreExistingShipment = {
  ...mockMTOShipment,
  boatShipment: {
    ...mockMTOShipment.boatShipment,
    lengthInInches: 352,
    widthInInches: 16,
    heightInInches: 34,
    hasTrailer: true,
    isRoadworthy: false,
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

const renderBoatShipmentLocationInfo = (props) => {
  return render(
    <MockProviders {...mockRoutingConfig}>
      <BoatShipmentLocationInfo {...props} />
    </MockProviders>,
  );
};

describe('Pickup info page', () => {
  it('renders the MtoShipmentForm component', () => {
    renderBoatShipmentLocationInfo();

    expect(screen.getByRole('heading', { level: 2, name: 'Pickup info' })).toBeInTheDocument();
  });

  it.each([[mockPreExistingShipment]])(
    'renders the form pre-filled when info has been entered previously',
    async (preExistingShipment) => {
      selectMTOShipmentById.mockImplementationOnce(() => preExistingShipment);

      renderBoatShipmentLocationInfo();

      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText('Use my current address')).not.toBeChecked();
      expect(screen.getAllByLabelText(/Address 1/)[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText(/City/)[0]).toHaveValue('San Antonio');
      expect(screen.getAllByLabelText(/State/)[0]).toHaveValue('TX');
      expect(screen.getAllByLabelText(/ZIP/)[0]).toHaveValue('78234');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText(/Address 1/)[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText(/City/)[1]).toHaveValue('Tacoma');
      expect(screen.getAllByLabelText(/State/)[1]).toHaveValue('WA');
      expect(screen.getAllByLabelText(/ZIP/)[1]).toHaveValue('98421');
    },
  );

  it('routes back to the previous page when the back button is clicked', async () => {
    renderBoatShipmentLocationInfo();

    const backButton = screen.getByRole('button', { name: /back/i });

    await userEvent.click(backButton);

    expect(mockNavigate).toHaveBeenCalledWith(shipmentEditPath);
  });

  it('can submit with pickup information successfully', async () => {
    const shipmentInfo = {
      pickupAddress: {
        streetAddress1: '6622 Airport Way S',
        streetAddress2: '#1430',
        city: 'San Marcos',
        state: 'TX',
        postalCode: '78666',
      },
    };

    const expectedPayload = {
      moveTaskOrderID: mockMoveId,
      shipmentType: SHIPMENT_TYPES.BOAT_HAUL_AWAY,
      pickupAddress: { ...shipmentInfo.pickupAddress },
      customerRemarks: mockMTOShipment.customerRemarks,
      requestedPickupDate: mockMTOShipment.requestedPickupDate,
      requestedDeliveryDate: mockMTOShipment.requestedDeliveryDate,
      destinationAddress: { ...mockMTOShipment.destinationAddress, streetAddress2: '' },
      secondaryDeliveryAddress: undefined,
      hasSecondaryDeliveryAddress: false,
      secondaryPickupAddress: undefined,
      hasSecondaryPickupAddress: false,
      tertiaryDeliveryAddress: undefined,
      hasTertiaryDeliveryAddress: false,
      tertiaryPickupAddress: undefined,
      hasTertiaryPickupAddress: false,
      agents: [
        { agentType: 'RELEASING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
        { agentType: 'RECEIVING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
      ],
      counselorRemarks: undefined,
    };
    delete expectedPayload.destinationAddress.id;

    const newUpdatedAt = '2021-06-11T21:20:22.150Z';
    const expectedUpdateResponse = {
      ...mockMTOShipment,
      pickupAddress: { ...shipmentInfo.pickupAddress },
      shipmentType: SHIPMENT_OPTIONS.HHG,
      eTag: window.btoa(newUpdatedAt),
      status: 'SUBMITTED',
    };

    patchMTOShipment.mockImplementation(() => Promise.resolve(expectedUpdateResponse));

    renderBoatShipmentLocationInfo({ isCreatePage: false, mtoShipment: mockMTOShipment });

    const pickupAddress1Input = screen.getAllByLabelText(/Address 1/)[0];
    await userEvent.clear(pickupAddress1Input);
    await userEvent.type(pickupAddress1Input, shipmentInfo.pickupAddress.streetAddress1);

    const pickupAddress2Input = screen.getAllByLabelText(/Address 2/)[0];
    await userEvent.clear(pickupAddress2Input);
    await userEvent.type(pickupAddress2Input, shipmentInfo.pickupAddress.streetAddress2);

    const pickupCityInput = screen.getAllByLabelText(/City/)[0];
    await userEvent.clear(pickupCityInput);
    await userEvent.type(pickupCityInput, shipmentInfo.pickupAddress.city);

    const pickupStateInput = screen.getAllByLabelText(/State/)[0];
    await userEvent.selectOptions(pickupStateInput, shipmentInfo.pickupAddress.state);

    const pickupPostalCodeInput = screen.getAllByLabelText(/ZIP/)[0];
    await userEvent.clear(pickupPostalCodeInput);
    await userEvent.type(pickupPostalCodeInput, shipmentInfo.pickupAddress.postalCode);

    const saveButton = await screen.findByRole('button', { name: 'Save & Continue' });
    expect(saveButton).not.toBeDisabled();
    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipment.id, expectedPayload, mockMTOShipment.eTag);
    });
  });
});

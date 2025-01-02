import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 as uuidv4 } from 'uuid';

import MobileHomeShipmentLocationInfo from './MobileHomeShipmentLocationInfo';

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

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const shipmentEditPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

const mockRoutingConfig = {
  path: customerRoutes.SHIPMENT_MOBILE_HOME_LOCATION_INFO,
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
  shipmentType: SHIPMENT_TYPES.MOBILE_HOME,
  mobileHomeShipment: {
    id: uuidv4(),
    year: 2022,
    make: 'Skyline Homes',
    model: 'Crown',
    lengthInInches: 252,
    widthInInches: 96,
    heightInInches: 84,
    eTag: window.btoa(new Date()),
  },
  eTag: window.btoa(new Date()),
  createdAt: '2021-06-11T18:12:11.918Z',
  customerRemarks: 'mock remarks',
  requestedPickupDate: '2021-08-01',
  requestedDeliveryDate: '2021-08-11',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    streetAddress2: '',
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
  mobileHomeShipment: {
    ...mockMTOShipment.mobileHomeShipment,
    lengthInInches: 352,
    widthInInches: 16,
    heightInInches: 34,
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

const renderMobileHomeShipmentLocationInfo = (props) => {
  return render(
    <MockProviders {...mockRoutingConfig}>
      <MobileHomeShipmentLocationInfo {...props} />
    </MockProviders>,
  );
};

describe('Pickup info page', () => {
  it('renders the MtoShipmentForm component', () => {
    renderMobileHomeShipmentLocationInfo();

    expect(screen.getByRole('heading', { level: 2, name: 'Pickup info' })).toBeInTheDocument();
  });

  it.each([[mockPreExistingShipment]])(
    'renders the form pre-filled when info has been entered previously',
    async (preExistingShipment) => {
      selectMTOShipmentById.mockImplementationOnce(() => preExistingShipment);

      renderMobileHomeShipmentLocationInfo();

      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText('Use my current address')).not.toBeChecked();
      expect(screen.getAllByLabelText(/Address 1/)[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByTestId(/City/)[0]).toHaveTextContent('San Antonio');
      expect(screen.getAllByTestId(/State/)[0]).toHaveTextContent('TX');
      expect(screen.getAllByTestId(/ZIP/)[0]).toHaveTextContent('78234');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText(/Address 1/)[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByTestId(/City/)[1]).toHaveTextContent('Tacoma');
      expect(screen.getAllByTestId(/State/)[1]).toHaveTextContent('WA');
      expect(screen.getAllByTestId(/ZIP/)[1]).toHaveTextContent('98421');
    },
  );

  it('routes back to the previous page when the back button is clicked', async () => {
    renderMobileHomeShipmentLocationInfo();

    const backButton = screen.getByRole('button', { name: /back/i });

    await userEvent.click(backButton);

    expect(mockNavigate).toHaveBeenCalledWith(shipmentEditPath);
  });

  it('can submit with pickup information successfully', async () => {
    const expectedPayload = {
      moveTaskOrderID: mockMoveId,
      shipmentType: SHIPMENT_TYPES.MOBILE_HOME,
      pickupAddress: { ...mockMTOShipment.pickupAddress },
      customerRemarks: mockMTOShipment.customerRemarks,
      requestedPickupDate: mockMTOShipment.requestedPickupDate,
      requestedDeliveryDate: mockMTOShipment.requestedDeliveryDate,
      destinationAddress: { ...mockMTOShipment.destinationAddress, streetAddress2: '' },
      hasSecondaryDeliveryAddress: false,
      hasSecondaryPickupAddress: false,
      hasTertiaryDeliveryAddress: false,
      hasTertiaryPickupAddress: false,
      destinationType: undefined,
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
      shipmentType: SHIPMENT_OPTIONS.HHG,
      eTag: window.btoa(newUpdatedAt),
      status: 'SUBMITTED',
    };

    patchMTOShipment.mockImplementation(() => Promise.resolve(expectedUpdateResponse));

    renderMobileHomeShipmentLocationInfo({ isCreatePage: false, mtoShipment: mockMTOShipment });

    const saveButton = await screen.findByRole('button', { name: 'Save & Continue' });
    expect(saveButton).not.toBeDisabled();
    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipment.id, expectedPayload, mockMTOShipment.eTag);
    });
  });
});

/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor, within, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentForm from './ShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ORDERS_TYPE } from 'constants/orders';
import { roleTypes } from 'constants/userRoles';
import { ADDRESS_UPDATE_STATUS, ppmShipmentStatuses } from 'constants/shipments';
import { tooRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';
import { validatePostalCode } from 'utils/validation';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { dateSelectionIsWeekendHoliday } from 'services/ghcApi';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const mockMutateFunction = jest.fn();
jest.mock('@tanstack/react-query', () => ({
  ...jest.requireActual('@tanstack/react-query'),
  useMutation: () => ({ mutate: mockMutateFunction }),
}));

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  dateSelectionIsWeekendHoliday: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const mockMtoShipment = {
  id: 'shipment123',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock customer remarks',
  counselorRemarks: 'mock counselor remarks',
  requestedPickupDate: '2020-03-01',
  requestedDeliveryDate: '2020-03-30',
  // requestedPickupDate: '2021-06-07',
  // requestedDeliveryDate: '2021-06-14',
  hasSecondaryDeliveryAddress: false,
  hasSecondaryPickupAddress: false,
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  mtoAgents: [
    {
      agentType: 'RELEASING_AGENT',
      email: 'jasn@email.com',
      firstName: 'Jason',
      lastName: 'Ash',
      phone: '999-999-9999',
    },
    {
      agentType: 'RECEIVING_AGENT',
      email: 'rbaker@email.com',
      firstName: 'Riley',
      lastName: 'Baker',
      phone: '863-555-9664',
    },
  ],
  mtoServiceItems: [
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:55.858Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1NS44NTgxMjVa',
      id: '7b7e94b1-0f34-418b-866f-d052e3a1c756',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DLH',
      reServiceID: '8d600f25-1def-422d-b159-617c7d59156e',
      reServiceName: 'Domestic linehaul',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:55.912Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1NS45MTI0NDFa',
      id: 'bf3516eb-1eaa-4e71-bd94-c523a6c866d0',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'FSC',
      reServiceID: '4780b30c-e846-437a-b39a-c499a6b09872',
      reServiceName: 'Fuel surcharge',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:55.968Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1NS45Njg1Nzda',
      id: '52b087b4-8e7f-4c96-939e-772cdd406e3a',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DOP',
      reServiceID: '2bc3e5cb-adef-46b1-bde9-55570bfdd43e',
      reServiceName: 'Domestic origin price',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:56.037Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1Ni4wMzc1OTla',
      id: 'c89ec6c0-a240-4478-afa0-52c5e2466ad4',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DDP',
      reServiceID: '50f1179a-3b72-4fa1-a951-fe5bcc70bd14',
      reServiceName: 'Domestic destination price',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:56.094Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1Ni4wOTQxMjRa',
      id: 'e26c9be3-dd55-4a0c-b002-f03258c40d06',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DPK',
      reServiceID: 'bdea5a8d-f15f-47d2-85c9-bba5694802ce',
      reServiceName: 'Domestic packing',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:56.162Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1Ni4xNjIzMTla',
      id: 'aca010a5-71e5-4994-b06b-97dfe4377f18',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DUPK',
      reServiceID: '15f01bc1-0754-4341-8e0f-25c8f04d5a77',
      reServiceName: 'Domestic unpacking',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
};

const defaultProps = {
  isCreatePage: true,
  submitHandler: jest.fn(),
  dateSelectionIsWeekendHoliday: jest.fn().mockImplementation(() => Promise.resolve()),
  newDutyLocationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
    streetAddress1: '123 Main',
    streetAddress2: '',
  },
  originDutyLocationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
    streetAddress1: '123 Main',
    streetAddress2: '',
  },
  serviceMember: {
    weightAllotment: {
      totalWeightSelf: 5000,
    },
    agency: '',
  },
  moveTaskOrderID: 'mock move id',
  mtoShipments: [],
  mtoShipment: mockMtoShipment,
  userRole: roleTypes.SERVICES_COUNSELOR,
  orderType: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  isForServivcesCounseling: false,
};

const mockShipmentWithDestinationType = {
  ...mockMtoShipment,
  displayDestinationType: true,
  destinationType: 'PLACE_ENTERED_ACTIVE_DUTY',
};

const mockPPMShipment = {
  ...mockMtoShipment,
  ppmShipment: {
    id: 'ppmShipmentID',
    shipmentId: 'shipment123',
    status: ppmShipmentStatuses.NEEDS_ADVANCE_APPROVAL,
    expectedDepartureDate: '2022-04-01',
    hasSecondaryPickupAddress: true,
    hasSecondaryDestinationAddress: true,
    pickupAddress: {
      streetAddress1: '111 Test Street',
      streetAddress2: '222 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42701',
    },
    secondaryPickupAddress: {
      streetAddress1: '777 Test Street',
      streetAddress2: '888 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42702',
    },
    destinationAddress: {
      streetAddress1: '222 Test Street',
      streetAddress2: '333 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42703',
    },
    secondaryDestinationAddress: {
      streetAddress1: '444 Test Street',
      streetAddress2: '555 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42701',
    },
    sitExpected: false,
    estimatedWeight: 4999,
    hasProGear: false,
    estimatedIncentive: 1234500,
    hasRequestedAdvance: true,
    advanceAmountRequested: 487500,
    advanceStatus: 'APPROVED',
    isActualExpenseReimbursement: true,
  },
};

const mockRejectedPPMShipment = {
  ...mockMtoShipment,
  ppmShipment: {
    id: 'ppmShipmentID',
    shipmentId: 'shipment123',
    status: ppmShipmentStatuses.NEEDS_ADVANCE_APPROVAL,
    expectedDepartureDate: '2022-04-01',
    hasSecondaryPickupAddress: true,
    hasSecondaryDestinationAddress: true,
    pickupAddress: {
      streetAddress1: '111 Test Street',
      streetAddress2: '222 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42701',
    },
    secondaryPickupAddress: {
      streetAddress1: '777 Test Street',
      streetAddress2: '888 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42702',
    },
    destinationAddress: {
      streetAddress1: '222 Test Street',
      streetAddress2: '333 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42703',
    },
    secondaryDestinationAddress: {
      streetAddress1: '444 Test Street',
      streetAddress2: '555 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42701',
    },
    sitExpected: false,
    estimatedWeight: 4999,
    hasProGear: false,
    estimatedIncentive: 1234500,
    hasRequestedAdvance: true,
    advanceAmountRequested: 487500,
    advanceStatus: 'REJECTED',
  },
};

const mockDeliveryAddressUpdate = {
  deliveryAddressUpdate: {
    contractorRemarks: 'Test Contractor Remark',
    id: 'c49f7921-5a6e-46b4-bb39-022583574453',
    newAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMy0wNy0xN1QxODowODowNi42NTU5MTVa',
      id: '6b57ce91-cabd-4e3b-9f48-ed4627d4878f',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    originalAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMy0wNy0xN1QxODowODowNi42NDkyNTha',
      id: '92509013-aafc-4892-a476-2e3b97e6933d',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    shipmentID: '5c84bcf3-92f7-448f-b0e1-e5378b6806df',
    status: 'REQUESTED',
  },
};

const defaultPropsRetirement = {
  ...defaultProps,
  displayDestinationType: true,
  destinationType: 'HOME_OF_RECORD',
  orderType: ORDERS_TYPE.RETIREMENT,
};

const defaultPropsSeparation = {
  ...defaultProps,
  displayDestinationType: true,
  destinationType: 'HOME_OF_SELECTION',
  orderType: ORDERS_TYPE.SEPARATION,
};

jest.mock('utils/validation', () => ({
  ...jest.requireActual('utils/validation'),
  validatePostalCode: jest.fn(),
}));
const mockRoutingOptions = {
  path: tooRoutes.BASE_SHIPMENT_EDIT_PATH,
  params: { moveCode: 'move123', shipmentId: 'shipment123' },
};

beforeEach(() => {
  jest.clearAllMocks();
});

const renderWithRouter = (ui) => {
  render(<MockProviders {...mockRoutingOptions}>{ui}</MockProviders>);
};

describe('ShipmentForm component', () => {
  describe('when creating a new shipment', () => {
    it('does not show the delete shipment button', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      const deleteButton = screen.queryByRole('button', { name: 'Delete shipment' });
      await waitFor(() => {
        expect(deleteButton).not.toBeInTheDocument();
      });
    });
  });

  describe('when creating a new HHG shipment', () => {
    it('renders the HHG shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');

      expect(screen.getByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Pickup Address')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use pickup address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText('First name')[0]).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.getAllByLabelText('Last name')[0]).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.getAllByLabelText('Phone')[0]).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.getAllByLabelText('Email')[0]).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.getByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      const deliveryLocationSectionHeadings = screen.getAllByText('Delivery location');
      expect(deliveryLocationSectionHeadings).toHaveLength(2);
      expect(deliveryLocationSectionHeadings[0]).toBeInstanceOf(HTMLParagraphElement);
      expect(deliveryLocationSectionHeadings[1]).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText('Yes')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Yes')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[1]).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText('First name')[1]).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getAllByLabelText('Last name')[1]).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.getAllByLabelText('Phone')[1]).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.getAllByLabelText('Email')[1]).toHaveAttribute('name', 'delivery.agent.email');

      expect(screen.getByText('Customer remarks')).toBeTruthy();

      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('Service Counselor - renders date alert warnings for pickup/delivery on date picker selection', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);
      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');

      await userEvent.type(screen.getByLabelText('Requested pickup date'), '26 Mar 2024');
      await userEvent.type(screen.getByLabelText('Requested delivery date'), '30 Mar 2024');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Requested pickup date 26 Mar 2024 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Requested delivery date 30 Mar 2024 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('uses the current residence address for pickup address when checked', async () => {
      const user = userEvent.setup();
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      await act(async () => {
        await user.click(screen.getByLabelText('Use pickup address'));
      });

      expect((await screen.findAllByLabelText('Address 1'))[0]).toHaveValue(
        defaultProps.currentResidence.streetAddress1,
      );

      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue(defaultProps.currentResidence.city);
      expect(screen.getAllByLabelText('State')[0]).toHaveValue(defaultProps.currentResidence.state);
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue(defaultProps.currentResidence.postalCode);
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      await act(async () => {
        await userEvent.click(screen.getAllByLabelText('Yes')[1]);
      });

      expect((await screen.findAllByLabelText('Address 1'))[0]).toHaveAttribute(
        'name',
        'pickup.address.streetAddress1',
      );
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveAttribute('name', 'delivery.address.streetAddress1');

      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveAttribute('name', 'pickup.address.streetAddress2');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveAttribute('name', 'delivery.address.streetAddress2');

      expect(screen.getAllByLabelText('City')[0]).toHaveAttribute('name', 'pickup.address.city');
      expect(screen.getAllByLabelText('City')[1]).toHaveAttribute('name', 'delivery.address.city');

      expect(screen.getAllByLabelText('State')[0]).toHaveAttribute('name', 'pickup.address.state');
      expect(screen.getAllByLabelText('State')[1]).toHaveAttribute('name', 'delivery.address.state');

      expect(screen.getAllByLabelText('ZIP')[0]).toHaveAttribute('name', 'pickup.address.postalCode');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveAttribute('name', 'delivery.address.postalCode');
    });

    it('renders a delivery address type for retirement orders type', async () => {
      renderWithRouter(<ShipmentForm {...defaultPropsRetirement} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      await act(async () => {
        await userEvent.click(screen.getAllByLabelText('Yes')[1]);
      });

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.getAllByLabelText('Destination type')[0]).toHaveAttribute('name', 'destinationType');
    });

    it('does not render delivery address type for PCS order type', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);
      await act(async () => {
        await userEvent.click(screen.getAllByLabelText('Yes')[1]);
      });

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByLabelText('Destination type')).toBeNull();
    });

    it('renders a delivery address type for separation orders type', async () => {
      renderWithRouter(<ShipmentForm {...defaultPropsSeparation} shipmentType={SHIPMENT_OPTIONS.HHG} />);
      await act(async () => {
        await userEvent.click(screen.getAllByLabelText('Yes')[1]);
      });

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.getAllByLabelText('Destination type')[0]).toHaveAttribute('name', 'destinationType');
    });

    it('does not render an Accounting Codes section', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByRole('heading', { name: 'Accounting codes' })).not.toBeInTheDocument();
    });

    it('does not render NTS release-only sections', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByText(/Shipment weight (lbs)/)).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility info' })).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility address' })).not.toBeInTheDocument();
    });
  });

  describe('editing an already existing HHG shipment', () => {
    it('renders the HHG shipment form with pre-filled values', async () => {
      // For some reason need this mock here.
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          displayDestinationType={false}
        />,
      );

      expect(await screen.findByLabelText('Requested pickup date')).toHaveValue('01 Mar 2020');
      expect(screen.getByLabelText('Use pickup address')).not.toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue('San Antonio');
      expect(screen.getAllByLabelText('State')[0]).toHaveValue('TX');
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue('78234');
      expect(screen.getAllByLabelText('First name')[0]).toHaveValue('Jason');
      expect(screen.getAllByLabelText('Last name')[0]).toHaveValue('Ash');
      expect(screen.getAllByLabelText('Phone')[0]).toHaveValue('999-999-9999');
      expect(screen.getAllByLabelText('Email')[0]).toHaveValue('jasn@email.com');
      expect(screen.getByLabelText('Requested delivery date')).toHaveValue('30 Mar 2020');
      expect(screen.getAllByLabelText('Yes')[0]).not.toBeChecked();
      expect(screen.getAllByLabelText('Yes')[1]).toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[1]).toHaveValue('Tacoma');
      expect(screen.getAllByLabelText('State')[1]).toHaveValue('WA');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveValue('98421');
      expect(screen.getAllByLabelText('First name')[1]).toHaveValue('Riley');
      expect(screen.getAllByLabelText('Last name')[1]).toHaveValue('Baker');
      expect(screen.getAllByLabelText('Phone')[1]).toHaveValue('863-555-9664');
      expect(screen.getAllByLabelText('Email')[1]).toHaveValue('rbaker@email.com');
      expect(screen.getByText('Customer remarks')).toBeTruthy();
      expect(screen.getByText('mock customer remarks')).toBeTruthy();
      expect(screen.getByLabelText('Counselor remarks')).toHaveValue('mock counselor remarks');

      const noDestinationTypeRadioButton = await screen.getAllByLabelText('No')[1];
      await act(async () => {
        await userEvent.click(noDestinationTypeRadioButton);
      });
      expect(screen.getByText('We can use the zip of their new duty location:')).toBeTruthy();
      expect(screen.queryByLabelText('Destination type')).toBeNull();
    });
  });

  describe('editing an already existing HHG shipment for retiree/separatee', () => {
    it('renders the HHG shipment form with pre-filled values', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultPropsRetirement}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          mtoShipment={mockShipmentWithDestinationType}
          displayDestinationType
        />,
      );

      expect(await screen.findByLabelText('Requested pickup date')).toHaveValue('01 Mar 2020');
      expect(screen.getByLabelText('Use pickup address')).not.toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue('San Antonio');
      expect(screen.getAllByLabelText('State')[0]).toHaveValue('TX');
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue('78234');
      expect(screen.getAllByLabelText('First name')[0]).toHaveValue('Jason');
      expect(screen.getAllByLabelText('Last name')[0]).toHaveValue('Ash');
      expect(screen.getAllByLabelText('Phone')[0]).toHaveValue('999-999-9999');
      expect(screen.getAllByLabelText('Email')[0]).toHaveValue('jasn@email.com');
      expect(screen.getByLabelText('Requested delivery date')).toHaveValue('30 Mar 2020');
      expect(screen.getAllByLabelText('Yes')[0]).not.toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[1]).toHaveValue('Tacoma');
      expect(screen.getAllByLabelText('State')[1]).toHaveValue('WA');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveValue('98421');
      expect(screen.getAllByLabelText('First name')[1]).toHaveValue('Riley');
      expect(screen.getAllByLabelText('Last name')[1]).toHaveValue('Baker');
      expect(screen.getAllByLabelText('Phone')[1]).toHaveValue('863-555-9664');
      expect(screen.getAllByLabelText('Email')[1]).toHaveValue('rbaker@email.com');
      expect(screen.getByText('Customer remarks')).toBeTruthy();
      expect(screen.getByText('mock customer remarks')).toBeTruthy();
      expect(screen.getByLabelText('Counselor remarks')).toHaveValue('mock counselor remarks');
      expect(screen.getByLabelText('Destination type')).toHaveValue('PLACE_ENTERED_ACTIVE_DUTY');
      expect(screen.queryByTestId('alert')).not.toBeInTheDocument();

      const noDestinationTypeRadioButton = await screen.getAllByLabelText('No')[1];
      await act(async () => {
        await userEvent.click(noDestinationTypeRadioButton);
      });
      expect(screen.getByText('We can use the zip of their HOR, HOS or PLEAD:')).toBeTruthy();
      expect(screen.getByLabelText('Destination type')).toBeVisible();
    });

    const runAlertingTest = async (shipmentType) => {
      renderWithRouter(
        <ShipmentForm
          {...defaultPropsRetirement}
          isCreatePage={false}
          shipmentType={shipmentType}
          mtoShipment={{ ...mockShipmentWithDestinationType, ...mockDeliveryAddressUpdate }}
          displayDestinationType
        />,
      );

      const alerts = await screen.findAllByTestId('alert');
      expect(alerts).toHaveLength(2); // Should have 2 alerts shown due to the address update request
      expect(alerts[0]).toHaveTextContent('Request needs review. See delivery location to proceed.');
      expect(alerts[1]).toHaveTextContent(
        'Pending delivery location change request needs review. Review request to proceed.',
      );
    };

    describe('shipment address change request', () => {
      it('displays appropriate alerting when an address change is requested for HHG shipment', async () => {
        await runAlertingTest(SHIPMENT_OPTIONS.HHG);
      });

      it('displays appropriate alerting when an address change is requested for NTSr shipment', async () => {
        await runAlertingTest(SHIPMENT_OPTIONS.NTSR);
      });

      it('opens a closeable modal when Review Request is clicked', async () => {
        const user = userEvent.setup();

        const shipmentType = SHIPMENT_OPTIONS.HHG;

        renderWithRouter(
          <ShipmentForm
            {...defaultPropsRetirement}
            isCreatePage={false}
            shipmentType={shipmentType}
            mtoShipment={{ ...mockShipmentWithDestinationType, ...mockDeliveryAddressUpdate, shipmentType }}
            displayDestinationType
          />,
        );

        const queryForModal = () => screen.queryByTestId('modal');

        const reviewRequestLink = await screen.findByRole('button', { name: 'Review request' });

        // confirm the modal is not already present
        expect(queryForModal()).not.toBeInTheDocument();

        // Open the modal
        await act(async () => {
          await user.click(reviewRequestLink);
        });

        await waitFor(() => expect(queryForModal()).toBeInTheDocument());

        // Close the modal
        const modalCancel = within(queryForModal()).queryByText('Cancel');

        expect(modalCancel).toBeInTheDocument();

        await act(async () => {
          await user.click(modalCancel);
        });

        // Confirm the modal has been closed
        expect(queryForModal()).not.toBeInTheDocument();
      });

      const runShipmentAddressUpdateTest = async (shipmentType) => {
        const user = userEvent.setup();
        const eTag = '8c32882e7793d9da88e0fdfd68672e2ead2f';

        renderWithRouter(
          <ShipmentForm
            {...defaultPropsRetirement}
            isCreatePage={false}
            shipmentType={shipmentType}
            mtoShipment={{ ...mockShipmentWithDestinationType, ...mockDeliveryAddressUpdate, eTag, shipmentType }}
            displayDestinationType
          />,
        );

        const queryForModal = () => screen.queryByTestId('modal');
        const findAlerts = async () => screen.findAllByTestId('alert');

        const reviewRequestLink = await screen.findByRole('button', { name: 'Review request' });

        expect(await findAlerts()).toHaveLength(2);

        // Open the modal
        await act(async () => {
          await user.click(reviewRequestLink);
        });
        const modal = queryForModal();

        expect(modal).toBeInTheDocument();

        // Fill and submit
        const approvalQuestion = within(modal).getByRole('group', { name: 'Approve address change?' });
        const approvalYes = within(approvalQuestion).getByRole('radio', { name: 'Yes' });
        const officeRemarks = within(modal).getByLabelText('Office remarks');
        const save = within(modal).getByRole('button', { name: 'Save' });

        const officeRemarksAnswer = 'Here are my remarks from the office';
        await act(async () => {
          await user.click(approvalYes);
          await user.type(officeRemarks, officeRemarksAnswer);
          await user.click(save);
        });

        // Confirm that the request was triggered
        expect(mockMutateFunction).toHaveBeenCalledTimes(1);
        expect(mockMutateFunction).toHaveBeenCalledWith({
          shipmentID: mockShipmentWithDestinationType.id,
          ifMatchETag: eTag,
          body: {
            status: ADDRESS_UPDATE_STATUS.APPROVED,
            officeRemarks: officeRemarksAnswer,
          },
          successCallback: expect.any(Function),
        });
      };

      it('allows a shipment address update review to be submitted via the modal for an HHG shipment', async () => {
        await runShipmentAddressUpdateTest(SHIPMENT_OPTIONS.HHG);
      });

      it('allows a shipment address update review to be submitted via the modal for an NTSr shipment', async () => {
        await runShipmentAddressUpdateTest(SHIPMENT_OPTIONS.NTSR);
      });
    });
  });

  describe('creating a new NTS shipment', () => {
    it('renders the NTS shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTS} />);

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');

      expect(screen.getByLabelText('Requested pickup date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Pickup Address')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use pickup address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('City')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('State')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByLabelText('ZIP')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('First name')).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.getByLabelText('Last name')).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.getByLabelText('Phone')).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.getByLabelText('Email')).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.queryByText('Delivery location')).not.toBeInTheDocument();
      expect(screen.queryByText(/Receiving agent/)).not.toBeInTheDocument();

      expect(screen.getByText('Customer remarks')).toBeTruthy();

      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);

      expect(screen.queryByRole('heading', { level: 2, name: 'Vendor' })).not.toBeInTheDocument();
    });

    it('renders an Accounting Codes section', async () => {
      renderWithRouter(
        <ShipmentForm {...defaultProps} TACs={{ HHG: '1234', NTS: '5678' }} shipmentType={SHIPMENT_OPTIONS.NTS} />,
      );

      expect(await screen.findByText(/Accounting codes/)).toBeInTheDocument();
      expect(screen.getByLabelText('1234 (HHG)')).toBeInTheDocument();
      expect(screen.getByText('No SAC code entered.')).toBeInTheDocument();
    });

    it('does not render NTS release-only sections', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTS} />);

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');
      expect(screen.queryByText(/Shipment weight (lbs)/)).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility info' })).not.toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility address' })).not.toBeInTheDocument();
    });
  });

  describe('editing an already existing NTS shipment', () => {
    it('pre-fills the Accounting Codes section', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          mtoShipment={{
            ...mockMtoShipment,
            tacType: 'NTS',
            sacType: 'HHG',
          }}
          TACs={{ HHG: '1234', NTS: '5678' }}
          SACs={{ HHG: '000012345' }}
          shipmentType={SHIPMENT_OPTIONS.NTS}
        />,
      );

      expect(await screen.findByText(/Accounting codes/)).toBeInTheDocument();
      expect(screen.getByLabelText('1234 (HHG)')).not.toBeChecked();
      expect(screen.getByLabelText('5678 (NTS)')).toBeChecked();
      expect(screen.getByLabelText('000012345 (HHG)')).toBeChecked();
    });

    it('sends an empty string when clearing LOA types when updating a shipment', async () => {
      const mockSubmitHandler = jest.fn().mockResolvedValue(null);

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          mtoShipment={{
            ...mockMtoShipment,
            tacType: 'NTS',
            sacType: 'HHG',
          }}
          TACs={{ HHG: '1234', NTS: '5678' }}
          SACs={{ HHG: '000012345', NTS: '2222' }}
          shipmentType={SHIPMENT_OPTIONS.NTS}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );

      await act(async () => {
        await userEvent.click(screen.getByTestId('clearSelection-sacType'));
      });
      const saveButton = screen.getByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await act(async () => {
        await userEvent.click(saveButton);
      });

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({ tacType: 'NTS', sacType: '' }),
          }),
          expect.objectContaining({
            onError: expect.any(Function),
            onSuccess: expect.any(Function),
          }),
        );
      });
    });

    it('does not send undefined LOA types when creating shipment', async () => {
      const mockSubmitHandler = jest.fn().mockResolvedValue(null);

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          mtoShipment={{
            ...mockMtoShipment,
          }}
          shipmentType={SHIPMENT_OPTIONS.NTS}
          submitHandler={mockSubmitHandler}
          isCreatePage
        />,
      );

      await act(async () => {
        await userEvent.type(screen.getByLabelText('Requested pickup date'), '26 Mar 2022');
        await userEvent.click(screen.getByTestId('useCurrentResidence'));
      });

      const saveButton = screen.getByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await act(async () => {
        await userEvent.click(saveButton);
      });

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.not.objectContaining({ tacType: expect.any(String), sacType: expect.any(String) }),
          }),
          expect.objectContaining({
            onError: expect.any(Function),
            onSuccess: expect.any(Function),
          }),
        );
      });
    });
  });

  describe('creating a new NTS-release shipment', () => {
    it('renders the NTS-release shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');

      expect(screen.queryByText('Pickup Address')).not.toBeInTheDocument();
      expect(screen.queryByText(/Releasing agent/)).not.toBeInTheDocument();

      expect(screen.getByLabelText('Requested delivery date')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Delivery location')).toBeInstanceOf(HTMLLegendElement);

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('First name')).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getByLabelText('Last name')).toHaveAttribute('name', 'delivery.agent.lastName');

      expect(screen.getByText('Customer remarks')).toBeTruthy();
      expect(screen.getByLabelText('Counselor remarks')).toBeInstanceOf(HTMLTextAreaElement);

      expect(screen.queryByRole('heading', { level: 2, name: 'Vendor' })).not.toBeInTheDocument();
    });

    it('renders an Accounting Codes section', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText(/Accounting codes/)).toBeInTheDocument();
    });

    it('renders the NTS release-only sections', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTSR} />);

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');
      expect(screen.getByText(/Previously recorded weight \(lbs\)/)).toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility info' })).toBeInTheDocument();
      expect(screen.queryByRole('heading', { name: 'Storage facility address' })).toBeInTheDocument();
    });
  });

  describe('as a TOO', () => {
    it('create new - HHG: displays date alerts for pickup/delivery for weekends', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderWithRouter(
        <ShipmentForm {...defaultProps} isCreatePage shipmentType={SHIPMENT_OPTIONS.HHG} userRole={roleTypes.TOO} />,
      );
      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      await userEvent.type(screen.getByLabelText('Requested pickup date'), '26 Mar 2024');
      await userEvent.type(screen.getByLabelText('Requested delivery date'), '30 Mar 2024');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Requested pickup date 26 Mar 2024 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Requested delivery date 30 Mar 2024 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('edit-HHG: pageload displays date alerts for pickup/delivery for weekends', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          userRole={roleTypes.TOO}
        />,
      );
      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByRole('heading', { level: 2, name: 'Vendor' })).not.toBeInTheDocument();
      expect(await screen.findByLabelText('Requested pickup date')).toHaveValue('01 Mar 2020');
      expect(await screen.findByLabelText('Requested delivery date')).toHaveValue('30 Mar 2020');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Requested pickup date 01 Mar 2020 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Requested delivery date 30 Mar 2020 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('edit-HHG: pageload displays date alerts for pickup/delivery for holiday', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          userRole={roleTypes.TOO}
        />,
      );
      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Requested pickup date 01 Mar 2020 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Requested delivery date 30 Mar 2020 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('edit-HHG: pageload displays date alerts for pickup/delivery for weekend and holiday', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          userRole={roleTypes.TOO}
        />,
      );
      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Requested pickup date 01 Mar 2020 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Requested delivery date 30 Mar 2020 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('renders the HHG shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.HHG} userRole={roleTypes.TOO} />);

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');
      expect(screen.queryByRole('heading', { level: 2, name: 'Vendor' })).not.toBeInTheDocument();
      expect(screen.getByLabelText('Requested pickup date')).toBeInTheDocument();
      expect(screen.getByText('Pickup Address')).toBeInTheDocument();
      expect(screen.getByLabelText('Requested delivery date')).toBeInTheDocument();
      expect(screen.getByText(/Receiving agent/).parentElement).toBeInTheDocument();
      expect(screen.getByText('Customer remarks')).toBeInTheDocument();
      expect(screen.getByText('Counselor remarks')).toBeInTheDocument();
    });

    it('renders the NTS shipment form', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTS} userRole={roleTypes.TOO} />);

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');
      expect(screen.getByLabelText('Requested pickup date')).toBeInTheDocument();
      expect(screen.getByLabelText('Requested delivery date')).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Vendor' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Storage facility info' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2, name: 'Storage facility address' })).toBeInTheDocument();
    });

    it('create new - NTS: displays date alerts for pickup/delivery for weekends', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderWithRouter(
        <ShipmentForm {...defaultProps} isCreatePage shipmentType={SHIPMENT_OPTIONS.NTS} userRole={roleTypes.TOO} />,
      );
      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');
      await userEvent.type(screen.getByLabelText('Requested pickup date'), '26 Mar 2024');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Requested pickup date 26 Mar 2024 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('edit-NTS: pageload displays date alerts for pickup/delivery for weekend and holiday', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.NTS}
          userRole={roleTypes.TOO}
        />,
      );
      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Requested pickup date 01 Mar 2020 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Requested delivery date 30 Mar 2020 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('renders the NTS release shipment form', async () => {
      renderWithRouter(
        <ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.NTSR} userRole={roleTypes.TOO} />,
      );

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');

      expect(screen.getByRole('heading', { level: 2, name: 'Vendor' })).toBeInTheDocument();
      expect(screen.getByLabelText('Requested pickup date')).toBeInTheDocument();
      expect(screen.getByLabelText('Requested delivery date')).toBeInTheDocument();
    });
  });

  it('edit-NTSR: pageload displays date alerts for pickup/delivery for weekend and holiday', async () => {
    const expectedDateSelectionIsWeekendHolidayResponse = {
      country_code: 'US',
      country_name: 'United States',
      is_weekend: true,
      is_holiday: true,
    };
    dateSelectionIsWeekendHoliday.mockImplementation(() =>
      Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
    );
    renderWithRouter(
      <ShipmentForm
        {...defaultProps}
        isCreatePage={false}
        shipmentType={SHIPMENT_OPTIONS.NTSR}
        userRole={roleTypes.TOO}
      />,
    );
    expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');
    expect(
      screen.getByText(
        'Requested pickup date 01 Mar 2020 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
      ),
    ).toHaveClass('usa-alert__text');
    expect(
      screen.getByText(
        /Requested delivery date 30 Mar 2020 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
      ),
    ).toHaveClass('usa-alert__text');
  });

  it('create new - NTSR: displays date alerts for pickup/delivery for weekends', async () => {
    const expectedDateSelectionIsWeekendHolidayResponse = {
      country_code: 'US',
      country_name: 'United States',
      is_weekend: true,
      is_holiday: false,
    };
    dateSelectionIsWeekendHoliday.mockImplementation(() =>
      Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
    );
    renderWithRouter(
      <ShipmentForm {...defaultProps} isCreatePage shipmentType={SHIPMENT_OPTIONS.NTSR} userRole={roleTypes.TOO} />,
    );
    expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');
    await userEvent.type(screen.getByLabelText('Requested delivery date'), '01 Mar 2024');
    await waitFor(() => {
      expect(
        screen.getByText(
          /Requested delivery date 01 Mar 2024 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
        ),
      ).toHaveClass('usa-alert__text');
    });
  });

  describe('filling the form', () => {
    it('shows an error if the submitHandler returns an error', async () => {
      const mockSubmitHandler = jest.fn((payload, { onError }) => {
        // fire onError handler on form
        onError();
      });

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );

      const saveButton = screen.getByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();

      await act(async () => {
        await userEvent.click(saveButton);
      });

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalled();
      });

      expect(
        await screen.findByText('Something went wrong, and your changes were not saved. Please try again.'),
      ).toBeInTheDocument();
      expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('shows a specific error message if the submitHandler returns a specific error message', async () => {
      const mockSpecificMessage = 'The data entered no good.';
      const mockSubmitHandler = jest.fn((payload, { onError }) => {
        // fire onError handler on form
        onError({ response: { body: { message: mockSpecificMessage, status: 400 } } });
      });

      validatePostalCode.mockImplementation(() => Promise.resolve(false));

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );

      const saveButton = screen.getByRole('button', { name: 'Save and Continue' });
      expect(saveButton).not.toBeDisabled();
      await act(async () => {
        await userEvent.click(saveButton);
      });

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalled();
      });

      expect(await screen.findByText(mockSpecificMessage)).toBeInTheDocument();
      expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('shows an error if the submitHandler returns an error when editing a PPM', async () => {
      const mockSubmitHandler = jest.fn((payload, { onError }) => {
        // fire onError handler on form
        onError();
      });
      validatePostalCode.mockImplementation(() => Promise.resolve(false));

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );

      const saveButton = screen.getByRole('button', { name: 'Save and Continue' });
      expect(saveButton).not.toBeDisabled();
      await act(async () => {
        await userEvent.click(saveButton);
      });

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalled();
      });

      expect(
        await screen.findByText('Something went wrong, and your changes were not saved. Please try again.'),
      ).toBeInTheDocument();
      expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('shows an error if the submitHandler returns an error when creating a PPM', async () => {
      const mockSubmitHandler = jest.fn((payload, { onError }) => {
        // fire onError handler on form
        onError();
      });
      validatePostalCode.mockImplementation(() => Promise.resolve(false));

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
          submitHandler={mockSubmitHandler}
          isCreatePage
        />,
      );

      await act(async () => {
        await userEvent.type(screen.getByLabelText('Planned Departure Date'), '26 Mar 2022');

        await userEvent.type(screen.getAllByLabelText('Address 1')[0], 'Test Street 1');
        await userEvent.type(screen.getAllByLabelText('City')[0], 'TestOne City');
        const pickupStateInput = screen.getAllByLabelText('State')[0];
        await userEvent.selectOptions(pickupStateInput, 'CA');
        await userEvent.type(screen.getAllByLabelText('ZIP')[0], '90210');

        await userEvent.type(screen.getAllByLabelText('Address 1')[1], 'Test Street 3');
        await userEvent.type(screen.getAllByLabelText('City')[1], 'TestTwo City');
        const destinationStateInput = screen.getAllByLabelText('State')[1];
        await userEvent.selectOptions(destinationStateInput, 'CA');
        await userEvent.type(screen.getAllByLabelText('ZIP')[1], '90210');

        await userEvent.type(screen.getByLabelText('Estimated PPM weight'), '1000');

        const saveButton = screen.getByRole('button', { name: 'Save and Continue' });
        expect(saveButton).not.toBeDisabled();
        await userEvent.click(saveButton);
      });

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalled();
      });

      expect(
        await screen.findByText('Something went wrong, and your changes were not saved. Please try again.'),
      ).toBeInTheDocument();
      expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('saves the update to the counselor remarks when the save button is clicked', async () => {
      const newCounselorRemarks = 'Counselor remarks';

      const expectedPayload = {
        body: {
          customerRemarks: 'mock customer remarks',
          counselorRemarks: newCounselorRemarks,
          hasSecondaryDeliveryAddress: false,
          hasSecondaryPickupAddress: false,
          hasTertiaryDeliveryAddress: false,
          hasTertiaryPickupAddress: false,
          destinationAddress: {
            streetAddress1: '441 SW Rio de la Plata Drive',
            city: 'Tacoma',
            state: 'WA',
            postalCode: '98421',
            streetAddress2: '',
          },
          pickupAddress: {
            streetAddress1: '812 S 129th St',
            city: 'San Antonio',
            state: 'TX',
            postalCode: '78234',
            streetAddress2: '',
          },
          agents: [
            {
              agentType: 'RELEASING_AGENT',
              email: 'jasn@email.com',
              firstName: 'Jason',
              lastName: 'Ash',
              phone: '999-999-9999',
            },
            {
              agentType: 'RECEIVING_AGENT',
              email: 'rbaker@email.com',
              firstName: 'Riley',
              lastName: 'Baker',
              phone: '863-555-9664',
            },
          ],
          requestedDeliveryDate: '2020-03-30',
          requestedPickupDate: '2020-03-01',
          shipmentType: SHIPMENT_OPTIONS.HHG,
        },
        shipmentID: 'shipment123',
        moveTaskOrderID: 'mock move id',
        normalize: false,
      };

      const patchResponse = {
        ...expectedPayload,
        created_at: '2021-02-08T16:48:04.117Z',
        updated_at: '2021-02-11T16:48:04.117Z',
      };

      const mockSubmitHandler = jest.fn(() => Promise.resolve(patchResponse));

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.HHG}
          submitHandler={mockSubmitHandler}
          isCreatePage={false}
        />,
      );
      const counselorRemarks = await screen.findByLabelText('Counselor remarks');

      await act(async () => {
        await userEvent.clear(counselorRemarks);
        await userEvent.type(counselorRemarks, newCounselorRemarks);
        const saveButton = screen.getByRole('button', { name: 'Save' });
        expect(saveButton).not.toBeDisabled();
        await userEvent.click(saveButton);
      });

      await waitFor(() => {
        expect(mockSubmitHandler).toHaveBeenCalledWith(expectedPayload, {
          onSuccess: expect.any(Function),
          onError: expect.any(Function),
        });
      });
    });
  });

  describe('external vendor shipment', () => {
    it('shows the TOO an alert', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.NTSR}
          mtoShipment={{ ...mockMtoShipment, usesExternalVendor: true }}
          isCreatePage={false}
          userRole={roleTypes.TOO}
        />,
      );

      expect(
        await screen.findByText(
          'The GHC prime contractor is not handling the shipment. Information will not be automatically shared with the movers handling it.',
        ),
      ).toBeInTheDocument();
    });

    it('does not show the SC an alert', async () => {
      renderWithRouter(
        <ShipmentForm
          // SC is default role from test props
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.NTSR}
          mtoShipment={{ ...mockMtoShipment, usesExternalVendor: true }}
          isCreatePage={false}
        />,
      );

      await waitFor(() => {
        expect(
          screen.queryByText(
            'The GHC prime contractor is not handling the shipment. Information will not be automatically shared with the movers handling it.',
          ),
        ).not.toBeInTheDocument();
      });
    });
  });

  describe('creating a new PPM shipment', () => {
    it('displays PPM content', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          isCreatePage
          userRole={roleTypes.SERVICES_COUNSELOR}
        />,
      );

      expect(await screen.findByTestId('tag')).toHaveTextContent('PPM');
    });
  });

  describe('TOO editing an already existing PPM shipment', () => {
    it('renders the PPM shipment form with pre-filled values as TOO', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
          userRole={roleTypes.TOO}
        />,
      );

      expect(await screen.getByLabelText('Planned Departure Date')).toHaveValue('01 Apr 2022');

      expect(await screen.getAllByLabelText('Address 1')[0]).toHaveValue(
        mockPPMShipment.ppmShipment.pickupAddress.streetAddress1,
      );
      expect(await screen.getAllByLabelText(/Address 2/)[0]).toHaveValue(
        mockPPMShipment.ppmShipment.pickupAddress.streetAddress2,
      );
      expect(await screen.getAllByLabelText('City')[0]).toHaveValue(mockPPMShipment.ppmShipment.pickupAddress.city);
      expect(await screen.getAllByLabelText('State')[0]).toHaveValue(mockPPMShipment.ppmShipment.pickupAddress.state);
      expect(await screen.getAllByLabelText('ZIP')[0]).toHaveValue(
        mockPPMShipment.ppmShipment.pickupAddress.postalCode,
      );

      expect(await screen.getAllByLabelText('Address 1')[1]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryPickupAddress.streetAddress1,
      );
      expect(await screen.getAllByLabelText(/Address 2/)[1]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryPickupAddress.streetAddress2,
      );
      expect(await screen.getAllByLabelText('City')[1]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryPickupAddress.city,
      );
      expect(await screen.getAllByLabelText('State')[1]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryPickupAddress.state,
      );
      expect(await screen.getAllByLabelText('ZIP')[1]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryPickupAddress.postalCode,
      );

      expect(await screen.getAllByLabelText('Address 1')[2]).toHaveValue(
        mockPPMShipment.ppmShipment.destinationAddress.streetAddress1,
      );
      expect(await screen.getAllByLabelText(/Address 2/)[2]).toHaveValue(
        mockPPMShipment.ppmShipment.destinationAddress.streetAddress2,
      );
      expect(await screen.getAllByLabelText('City')[2]).toHaveValue(
        mockPPMShipment.ppmShipment.destinationAddress.city,
      );
      expect(await screen.getAllByLabelText('State')[2]).toHaveValue(
        mockPPMShipment.ppmShipment.destinationAddress.state,
      );
      expect(await screen.getAllByLabelText('ZIP')[2]).toHaveValue(
        mockPPMShipment.ppmShipment.destinationAddress.postalCode,
      );

      expect(await screen.getAllByLabelText('Address 1')[3]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryDestinationAddress.streetAddress1,
      );
      expect(await screen.getAllByLabelText(/Address 2/)[3]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryDestinationAddress.streetAddress2,
      );
      expect(await screen.getAllByLabelText('City')[3]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryDestinationAddress.city,
      );
      expect(await screen.getAllByLabelText('State')[3]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryDestinationAddress.state,
      );
      expect(await screen.getAllByLabelText('ZIP')[3]).toHaveValue(
        mockPPMShipment.ppmShipment.secondaryDestinationAddress.postalCode,
      );

      expect(screen.getAllByLabelText('Yes')[0]).toBeChecked();
      expect(screen.getAllByLabelText('No')[0]).not.toBeChecked();
      expect(screen.getByLabelText('Estimated PPM weight')).toHaveValue('4,999');
      expect(screen.getAllByLabelText('Yes')[2]).toBeChecked();
      expect(screen.getAllByLabelText('No')[2]).not.toBeChecked();
    });

    it('renders the PPM shipment form with pre-filled requested values for Advance Page for TOO', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
          userRole={roleTypes.TOO}
        />,
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      expect(await screen.getByLabelText('No')).not.toBeChecked();
      expect(screen.getByLabelText('Yes')).toBeChecked();
      expect(screen.findByText('Estimated incentive: $12,345').toBeInTheDocument);
      expect(screen.getByLabelText('Amount requested')).toHaveValue('4,875');
      expect((await screen.findByText('Maximum advance: $7,407')).toBeInTheDocument);
      expect(screen.getByLabelText('Approve')).toBeChecked();

      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Save and Continue' }));
      });

      await waitFor(() => {
        expect(defaultProps.submitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              counselorRemarks: 'mock counselor remarks',
              ppmShipment: expect.objectContaining({
                hasRequestedAdvance: true,
                advanceAmountRequested: 487500,
                advanceStatus: 'APPROVED',
              }),
            }),
          }),
          expect.objectContaining({
            onSuccess: expect.any(Function),
          }),
        );
      });
    });
    describe('editing an already existing PPM shipment', () => {
      it('renders the PPM shipment form with pre-filled values', async () => {
        isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
        renderWithRouter(
          <ShipmentForm
            {...defaultProps}
            isCreatePage={false}
            shipmentType={SHIPMENT_OPTIONS.PPM}
            mtoShipment={mockPPMShipment}
          />,
        );

        expect(await screen.getByLabelText('Planned Departure Date')).toHaveValue('01 Apr 2022');

        expect(await screen.getAllByLabelText('Address 1')[0]).toHaveValue(
          mockPPMShipment.ppmShipment.pickupAddress.streetAddress1,
        );
        expect(await screen.getAllByLabelText(/Address 2/)[0]).toHaveValue(
          mockPPMShipment.ppmShipment.pickupAddress.streetAddress2,
        );
        expect(await screen.getAllByLabelText('City')[0]).toHaveValue(mockPPMShipment.ppmShipment.pickupAddress.city);
        expect(await screen.getAllByLabelText('State')[0]).toHaveValue(mockPPMShipment.ppmShipment.pickupAddress.state);
        expect(await screen.getAllByLabelText('ZIP')[0]).toHaveValue(
          mockPPMShipment.ppmShipment.pickupAddress.postalCode,
        );

        expect(await screen.getAllByLabelText('Address 1')[1]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryPickupAddress.streetAddress1,
        );
        expect(await screen.getAllByLabelText(/Address 2/)[1]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryPickupAddress.streetAddress2,
        );
        expect(await screen.getAllByLabelText('City')[1]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryPickupAddress.city,
        );
        expect(await screen.getAllByLabelText('State')[1]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryPickupAddress.state,
        );
        expect(await screen.getAllByLabelText('ZIP')[1]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryPickupAddress.postalCode,
        );

        expect(await screen.getAllByLabelText('Address 1')[2]).toHaveValue(
          mockPPMShipment.ppmShipment.destinationAddress.streetAddress1,
        );
        expect(await screen.getAllByLabelText(/Address 2/)[2]).toHaveValue(
          mockPPMShipment.ppmShipment.destinationAddress.streetAddress2,
        );
        expect(await screen.getAllByLabelText('City')[2]).toHaveValue(
          mockPPMShipment.ppmShipment.destinationAddress.city,
        );
        expect(await screen.getAllByLabelText('State')[2]).toHaveValue(
          mockPPMShipment.ppmShipment.destinationAddress.state,
        );
        expect(await screen.getAllByLabelText('ZIP')[2]).toHaveValue(
          mockPPMShipment.ppmShipment.destinationAddress.postalCode,
        );

        expect(await screen.getAllByLabelText('Address 1')[3]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryDestinationAddress.streetAddress1,
        );
        expect(await screen.getAllByLabelText(/Address 2/)[3]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryDestinationAddress.streetAddress2,
        );
        expect(await screen.getAllByLabelText('City')[3]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryDestinationAddress.city,
        );
        expect(await screen.getAllByLabelText('State')[3]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryDestinationAddress.state,
        );
        expect(await screen.getAllByLabelText('ZIP')[3]).toHaveValue(
          mockPPMShipment.ppmShipment.secondaryDestinationAddress.postalCode,
        );

        expect(screen.getAllByLabelText('Yes')[0]).toBeChecked();
        expect(screen.getAllByLabelText('No')[0]).not.toBeChecked();
        expect(screen.getAllByLabelText('Yes')[1]).toBeChecked();
        expect(screen.getAllByLabelText('No')[1]).not.toBeChecked();
        expect(screen.getByLabelText('Estimated PPM weight')).toHaveValue('4,999');
        expect(screen.getAllByLabelText('Yes')[3]).toBeChecked();
        expect(screen.getAllByLabelText('No')[3]).not.toBeChecked();
      });
    });
    it('renders the PPM shipment form with pre-filled requested values for Advance Page', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={mockPPMShipment}
        />,
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      expect(await screen.getByLabelText('No')).not.toBeChecked();
      expect(screen.getByLabelText('Yes')).toBeChecked();
      expect(screen.findByText('Estimated incentive: $12,345').toBeInTheDocument);
      expect(screen.getByLabelText('Amount requested')).toHaveValue('4,875');
      expect((await screen.findByText('Maximum advance: $7,407')).toBeInTheDocument);
      expect(screen.getByLabelText('Approve')).toBeChecked();
      expect(screen.getByLabelText('Counselor remarks')).toHaveValue('mock counselor remarks');

      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Save and Continue' }));
      });

      await waitFor(() => {
        expect(defaultProps.submitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              counselorRemarks: 'mock counselor remarks',
              ppmShipment: expect.objectContaining({
                hasRequestedAdvance: true,
                advanceAmountRequested: 487500,
                advanceStatus: 'APPROVED',
              }),
            }),
          }),
          expect.objectContaining({
            onSuccess: expect.any(Function),
          }),
        );
      });
    });

    it('validates the Advance Page making counselor remarks required when `Advance Requested?` is changed from Yes to No', async () => {
      const ppmShipmentWithoutRemarks = {
        ...mockPPMShipment,
        counselorRemarks: '',
      };

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={ppmShipmentWithoutRemarks}
        />,
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      expect(screen.getByLabelText('No')).not.toBeChecked();
      expect(screen.getByLabelText('Yes')).toBeChecked();
      // Selecting advance not requested
      await act(async () => {
        await userEvent.click(screen.getByLabelText('No'));
      });
      await waitFor(() => {
        expect(screen.getByLabelText('No')).toBeChecked();
        expect(screen.getByLabelText('Yes')).not.toBeChecked();
      });
      const requiredAlerts = screen.getAllByRole('alert');
      expect(requiredAlerts[0]).toHaveTextContent('Required');

      expect(screen.queryByLabelText('Amount requested')).not.toBeInTheDocument();

      await act(async () => {
        await userEvent.type(screen.getByLabelText('Counselor remarks'), 'retirees are not given advances');
        await userEvent.tab();
      });

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeEnabled();
      });

      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Save and Continue' }));
      });

      await waitFor(() => {
        expect(defaultProps.submitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              counselorRemarks: 'retirees are not given advances',
              ppmShipment: expect.objectContaining({ hasRequestedAdvance: false }),
            }),
          }),
          expect.objectContaining({
            onSuccess: expect.any(Function),
          }),
        );
      });
    });

    it('validates the Advance Page making counselor remarks required when advance amount is changed', async () => {
      const ppmShipmentWithoutRemarks = {
        ...mockPPMShipment,
        counselorRemarks: '',
      };

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={ppmShipmentWithoutRemarks}
        />,
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      const advanceAmountInput = screen.getByLabelText('Amount requested');

      expect(advanceAmountInput).toHaveValue('4,875');
      // Edit a requested advance amount
      await act(async () => {
        await userEvent.clear(advanceAmountInput);
        await userEvent.type(advanceAmountInput, '2,000');
        advanceAmountInput.blur();
      });
      await waitFor(() => {
        expect(advanceAmountInput).toHaveValue('2,000');
      });

      const requiredAlerts = screen.getAllByRole('alert');

      expect(requiredAlerts[0]).toHaveTextContent('Required');
    });

    it('marks amount requested input as min of $1 expected when conditionally displayed', async () => {
      renderWithRouter(
        <ShipmentForm {...defaultProps} isCreatePage={false} isAdvancePage shipmentType={SHIPMENT_OPTIONS.PPM} />,
      );

      const inputHasRequestedAdvance = screen.getByLabelText('Yes');
      await act(async () => {
        await userEvent.click(inputHasRequestedAdvance);
      });
      const advanceAmountRequested = screen.getByLabelText('Amount requested');
      await act(async () => {
        await userEvent.type(advanceAmountRequested, '0');
      });
      expect(advanceAmountRequested).toHaveValue('0');

      await waitFor(() => {
        const requiredAlerts = screen.getAllByRole('alert');
        expect(requiredAlerts[0]).toHaveTextContent('Enter an amount $1 or more.');
      });
    });

    it('sets `Counselor Remarks` as required when an advance request is rejected', async () => {
      const ppmShipmentWithoutRemarks = {
        ...mockPPMShipment,
        counselorRemarks: '',
      };

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={ppmShipmentWithoutRemarks}
        />,
        { wrapper: MockProviders },
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      expect(screen.getByLabelText('Approve')).toBeChecked();
      expect(screen.getByLabelText('Reject')).not.toBeChecked();

      const advanceAmountInput = screen.getByLabelText('Amount requested');
      expect(advanceAmountInput).toHaveValue('4,875');

      await act(async () => {
        // Edit a requested advance amount to different number to
        // test REVERT to save on REJECT
        await userEvent.clear(advanceAmountInput);
        await userEvent.type(advanceAmountInput, '2,000');
      });

      // Rejecting advance request
      await userEvent.click(screen.getByLabelText('Reject'));
      await waitFor(() => {
        expect(screen.getByLabelText('Approve')).not.toBeChecked();
        expect(screen.getByLabelText('Reject')).toBeChecked();

        // Verify original value was reset back 2000 to 4875. This only
        // happens when REJECT is selected.
        const advanceAmountInput2 = screen.getByLabelText('Amount requested');
        expect(advanceAmountInput2).toHaveValue('4,875');
      });
      const requiredAlert = screen.getAllByRole('alert');
      expect(requiredAlert[0]).toHaveTextContent('Required');

      await act(async () => {
        await userEvent.type(
          screen.getByLabelText('Counselor remarks'),
          'I, a service counselor, have rejected your advance request',
        );
        await userEvent.tab();
      });

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeEnabled();
      });

      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Save and Continue' }));
      });

      await waitFor(() => {
        expect(defaultProps.submitHandler).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              counselorRemarks: 'I, a service counselor, have rejected your advance request',
              ppmShipment: expect.objectContaining({ advanceStatus: 'REJECTED' }),
            }),
          }),
          expect.objectContaining({
            onSuccess: expect.any(Function),
          }),
        );
      });
    });

    it('sets to ACCEPT from REJECT if advance number is changed', async () => {
      const ppmShipment = {
        ...mockRejectedPPMShipment,
        counselorRemarks: 'test',
      };

      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          isCreatePage={false}
          isAdvancePage
          shipmentType={SHIPMENT_OPTIONS.PPM}
          mtoShipment={ppmShipment}
        />,
        { wrapper: MockProviders },
      );

      expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('Incentive & advance');
      expect(screen.getByLabelText('Reject')).toBeChecked();
      expect(screen.getByLabelText('Approve')).not.toBeChecked();

      const advanceAmountInput = screen.getByLabelText('Amount requested');
      expect(advanceAmountInput).toHaveValue('4,875');

      await act(async () => {
        await userEvent.clear(advanceAmountInput);
        await userEvent.type(advanceAmountInput, '2,000');
      });

      // test REJECT is changed to ACCEPT when advance number is changed
      expect(screen.getByLabelText('Reject')).not.toBeChecked();
      expect(screen.getByLabelText('Approve')).toBeChecked();
    });
  });

  describe('creating a new PPM shipment', () => {
    it('displays PPM content', async () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          isCreatePage
          userRole={roleTypes.SERVICES_COUNSELOR}
        />,
      );

      expect(await screen.findByTestId('tag')).toHaveTextContent('PPM');
      expect(screen.getByText('Is this PPM an Actual Expense Reimbursement?')).toBeInTheDocument();
      expect(screen.getByText('What address are you moving from?')).toBeInTheDocument();
      expect(screen.getByText('Second pickup address')).toBeInTheDocument();
      expect(
        screen.getByText(
          'Will you move any belongings from a second address? (Must be near the pickup address. Subject to approval.)',
        ),
      ).toBeInTheDocument();

      expect(screen.getByText('Delivery Address')).toBeInTheDocument();
      expect(screen.getByText('Second delivery address')).toBeInTheDocument();
      expect(
        screen.getByText(
          'Will you move any belongings to a second address? (Must be near the delivery address. Subject to approval.)',
        ),
      ).toBeInTheDocument();
    });
    it('displays the third pickup address question when the Yes option for second pickup address is selected', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          isCreatePage
          userRole={roleTypes.SERVICES_COUNSELOR}
        />,
      );
      expect(screen.queryByText('Third pickup address')).not.toBeInTheDocument();
      fireEvent.click(screen.getByTestId('has-secondary-pickup'));
      expect(await screen.findByText('Third pickup address')).toBeInTheDocument();
      expect(
        await screen.findByText(
          'Will you move any belongings from a third address? (Must be near the pickup address. Subject to approval.)',
        ),
      ).toBeInTheDocument();
    });
    it('displays the third delivery address question when the Yes option for second delivery address is selected', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      renderWithRouter(
        <ShipmentForm
          {...defaultProps}
          shipmentType={SHIPMENT_OPTIONS.PPM}
          isCreatePage
          userRole={roleTypes.SERVICES_COUNSELOR}
        />,
      );
      expect(screen.queryByText('Third delivery address')).not.toBeInTheDocument();
      fireEvent.click(screen.getByTestId('has-secondary-destination'));
      expect(await screen.findByText('Third delivery address')).toBeInTheDocument();
      expect(
        await screen.findByText(
          'Will you move any belongings to a third address? (Must be near the delivery address. Subject to approval.)',
        ),
      ).toBeInTheDocument();
    });
  });

  const mockPPMShipmentWithSIT = {
    sitEstimatedCost: 123400,
    sitEstimatedWeight: 2345,
    sitLocation: 'DESTINATION',
    sitEstimatedDepartureDate: '2022-10-29',
    sitEstimatedEntryDate: '2022-08-06',
    sitExpected: true,
    pickupAddress: {
      streetAddress1: '111 Test Street',
      streetAddress2: '222 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42701',
    },
    destinationAddress: {
      streetAddress1: '222 Test Street',
      streetAddress2: '333 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42703',
    },
  };

  const defaultSITProps = {
    ...defaultProps,
    shipmentType: SHIPMENT_OPTIONS.PPM,
    isAdvancePage: true,
    mtoShipment: {
      ...mockMtoShipment,
      ppmShipment: mockPPMShipmentWithSIT,
    },
    userRole: roleTypes.SERVICES_COUNSELOR,
  };

  describe('as a SC, the SIT details block', () => {
    it('displays when SIT is expected', () => {
      renderWithRouter(<ShipmentForm {...defaultSITProps} />);
      expect(screen.getByRole('heading', { level: 2, name: /Storage in transit \(SIT\)/ })).toBeInTheDocument();
    });
    it('does not display when SIT is not expected', () => {
      renderWithRouter(
        <ShipmentForm
          {...defaultSITProps}
          mtoShipment={{
            ...mockMtoShipment,
            ppmShipment: {
              ...mockPPMShipmentWithSIT,
              sitExpected: false,
            },
          }}
        />,
      );
      expect(screen.queryByRole('heading', { level: 2, name: /Storage in transit \(SIT\)/ })).not.toBeInTheDocument();
    });
    it('does not display for TOO', () => {
      renderWithRouter(<ShipmentForm {...defaultSITProps} userRole={roleTypes.TOO} />);
      expect(screen.queryByRole('heading', { level: 2, name: /Storage in transit \(SIT\)/ })).not.toBeInTheDocument();
    });
  });

  describe('creating a new Boat shipment', () => {
    it('renders the Boat shipment form correctly', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.BOAT_HAUL_AWAY} isCreatePage />);

      expect(await screen.findByTestId('tag')).toHaveTextContent('Boat');
      expect(await screen.findByText('Boat')).toBeInTheDocument();
      expect(screen.getByLabelText('Year')).toBeInTheDocument();
      expect(screen.getByLabelText('Make')).toBeInTheDocument();
      expect(screen.getByLabelText('Model')).toBeInTheDocument();
      expect(await screen.findByText('Length')).toBeInTheDocument();
      expect(await screen.findByText('Width')).toBeInTheDocument();
      expect(await screen.findByText('Height')).toBeInTheDocument();
      expect(await screen.findByText('Does the boat have a trailer?')).toBeInTheDocument();
      expect(await screen.findByText('What is the method of shipment?')).toBeInTheDocument();
      expect(await screen.findByText('Pickup details')).toBeInTheDocument();
      expect(await screen.findByText('Delivery details')).toBeInTheDocument();
      expect(await screen.findByText('Remarks')).toBeInTheDocument();
    });

    it('validates length and width input fields to ensure they accept only numeric values', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.BOAT_HAUL_AWAY} />);

      const lengthInput = await screen.findByTestId('lengthFeet');
      const widthInput = await screen.findByTestId('widthFeet');

      await act(async () => {
        userEvent.type(lengthInput, 'abc');
        userEvent.type(widthInput, 'xyz');
      });

      await waitFor(() => {
        expect(lengthInput).toHaveValue('');
        expect(widthInput).toHaveValue('');
      });
    });

    it('validates required fields for boat shipment', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.BOAT_HAUL_AWAY} />);

      const submitButton = screen.getByRole('button', { name: 'Save' });

      await act(async () => {
        userEvent.click(submitButton);
      });

      expect(submitButton).toBeDisabled();
    });

    it('validates the year field is within the valid range', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.BOAT_HAUL_AWAY} />);

      await act(async () => {
        await userEvent.click(screen.getByTestId('year'));
        await userEvent.type(screen.getByTestId('year'), '1600');
        const submitButton = screen.getByRole('button', { name: 'Save' });
        userEvent.click(submitButton);
      });

      expect(await screen.findByText('Invalid year')).toBeInTheDocument();
    });

    it('validates dimensions - fail', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.BOAT_HAUL_AWAY} />);

      // Enter dimensions below the required minimums
      await act(async () => {
        await userEvent.click(screen.getByTestId('lengthFeet'));
        await userEvent.type(screen.getByTestId('lengthFeet'), '10');
        await userEvent.click(screen.getByTestId('widthFeet'));
        await userEvent.type(screen.getByTestId('widthFeet'), '5');
        await userEvent.click(screen.getByTestId('heightFeet'));
        await userEvent.type(screen.getByTestId('heightFeet'), '6');
        const submitButton = screen.getByRole('button', { name: 'Save' });
        userEvent.click(submitButton);
      });

      expect(
        screen.queryByText(
          'The dimensions do not meet the requirements for a boat shipment. Please cancel and select a different shipment type.',
        ),
      ).toBeInTheDocument();
    });

    it('validates dimensions - pass', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.BOAT_HAUL_AWAY} />);

      // Enter dimensions below the required minimums
      await act(async () => {
        await userEvent.click(screen.getByTestId('lengthFeet'));
        await userEvent.type(screen.getByTestId('lengthFeet'), '15');
        await userEvent.click(screen.getByTestId('widthFeet'));
        await userEvent.type(screen.getByTestId('widthFeet'), '5');
        await userEvent.click(screen.getByTestId('heightFeet'));
        await userEvent.type(screen.getByTestId('heightFeet'), '6');
        const submitButton = screen.getByRole('button', { name: 'Save' });
        userEvent.click(submitButton);
      });

      expect(
        screen.queryByText(
          'The dimensions do not meet the requirements for a boat shipment. Please cancel and select a different shipment type.',
        ),
      ).not.toBeInTheDocument();
    });
  });

  describe('creating a new Mobile Home shipment', () => {
    it('renders the Mobile Home shipment form correctly', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.MOBILE_HOME} isCreatePage />);

      expect(screen.getByLabelText('Year')).toBeInTheDocument();
      expect(screen.getByLabelText('Make')).toBeInTheDocument();
      expect(screen.getByLabelText('Model')).toBeInTheDocument();
      expect(await screen.findByText('Length')).toBeInTheDocument();
      expect(await screen.findByText('Width')).toBeInTheDocument();
      expect(await screen.findByText('Height')).toBeInTheDocument();
      expect(await screen.findByText('Remarks')).toBeInTheDocument();
    });

    it('validates length and width input fields to ensure they accept only numeric values', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.MOBILE_HOME} />);

      const lengthInput = await screen.findByTestId('lengthFeet');
      const heightInput = await screen.findByTestId('heightFeet');
      const widthInput = await screen.findByTestId('widthFeet');

      await act(async () => {
        userEvent.type(lengthInput, 'abc');
        userEvent.type(heightInput, 'xyz');
        userEvent.type(widthInput, 'zyz');
      });

      await waitFor(() => {
        expect(lengthInput).toHaveValue('');
        expect(heightInput).toHaveValue('');
        expect(widthInput).toHaveValue('');
      });
    });

    it('validates required fields for Mobile Home shipment', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.MOBILE_HOME} />);

      const submitButton = screen.getByRole('button', { name: 'Save' });

      await act(async () => {
        userEvent.click(submitButton);
      });

      expect(submitButton).toBeDisabled();
    });

    it('validates the year field is within the valid range', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.MOBILE_HOME} />);

      await act(async () => {
        await userEvent.click(screen.getByTestId('year'));
        await userEvent.type(screen.getByTestId('year'), '1600');
        const submitButton = screen.getByRole('button', { name: 'Save' });
        userEvent.click(submitButton);
      });

      expect(await screen.findByText('Invalid year')).toBeInTheDocument();
    });

    it('validates dimensions - pass', async () => {
      renderWithRouter(<ShipmentForm {...defaultProps} shipmentType={SHIPMENT_OPTIONS.MOBILE_HOME} />);

      // Enter dimensions below the required minimums
      await act(async () => {
        await userEvent.click(screen.getByTestId('lengthFeet'));
        await userEvent.type(screen.getByTestId('lengthFeet'), '15');
        await userEvent.click(screen.getByTestId('widthFeet'));
        await userEvent.type(screen.getByTestId('widthFeet'), '5');
        await userEvent.click(screen.getByTestId('heightFeet'));
        await userEvent.type(screen.getByTestId('heightFeet'), '6');
        const submitButton = screen.getByRole('button', { name: 'Save' });
        userEvent.click(submitButton);
      });

      expect(screen.queryByText('Where and when should the movers deliver your mobile home?')).not.toBeInTheDocument();
    });
  });
});

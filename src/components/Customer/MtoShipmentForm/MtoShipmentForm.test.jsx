/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { generatePath } from 'react-router-dom';
import { waitFor, screen, fireEvent, act, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';
import { v4 as uuidv4 } from 'uuid';
import moment from 'moment';

import MtoShipmentForm from './MtoShipmentForm';

import { configureStore } from 'shared/store';
import { customerRoutes } from 'constants/routes';
import {
  createMTOShipment,
  getResponseError,
  patchMTOShipment,
  dateSelectionIsWeekendHoliday,
} from 'services/internalApi';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES, MOVE_LOCKED_WARNING } from 'shared/constants';
import { renderWithRouter } from 'testUtils';
import { ORDERS_TYPE } from 'constants/orders';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { formatDateWithUTC, formatDateForDatePicker } from 'shared/dates';
import { boatShipmentTypes } from 'constants/shipments';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  createMTOShipment: jest.fn(),
  getResponseError: jest.fn(),
  patchMTOShipment: jest.fn(),
  dateSelectionIsWeekendHoliday: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const moveId = uuidv4();

const defaultProps = {
  isCreatePage: true,
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  showLoggedInUser: jest.fn(),
  createMTOShipment: jest.fn(),
  updateMTOShipment: jest.fn(),
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
  orders: {
    orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    has_dependents: false,
    authorizedWeight: 5000,
  },
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

const ubProps = {
  isCreatePage: true,
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  showLoggedInUser: jest.fn(),
  createMTOShipment: jest.fn(),
  updateMTOShipment: jest.fn(),
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
  orders: {
    orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    has_dependents: false,
    entitlement: {
      ub_allowance: 600,
    },
  },
  shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
};

const updatedAt = '2021-06-11T18:12:11.918Z';
const tomorrow = formatDateWithUTC(moment().add(1, 'days').toDate(), 'YYYY-MM-DD');
const tomorrowDatePicker = formatDateForDatePicker(formatDateWithUTC(moment().add(1, 'days').toDate()));

const mockMtoShipmentHHGWithDest = {
  id: uuidv4(),
  eTag: window.btoa(updatedAt),
  createdAt: '2021-06-11T18:12:11.918Z',
  updatedAt,
  moveTaskOrderId: moveId,
  customerRemarks: 'mock remarks',
  requestedPickupDate: tomorrow,
  requestedDeliveryDate: '2021-06-14',
  newDutyLocationAddress: {
    id: uuidv4(),
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
  },
  pickupAddress: {
    id: uuidv4(),
    streetAddress1: '812 S 129th St',
    streetAddress2: '#123',
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
  secondaryDeliveryAddress: {
    id: uuidv4(),
    streetAddress1: '123 N Main',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  shipmentType: SHIPMENT_OPTIONS.HHG,
  hasSecondaryPickupAddress: false,
  hasSecondaryDeliveryAddress: true,
  hasTertiaryPickupAddress: false,
  hasTertiaryDeliveryAddress: false,
};

const mockMtoShipmentUB = {
  id: uuidv4(),
  eTag: window.btoa(updatedAt),
  createdAt: '2021-06-11T18:12:11.918Z',
  updatedAt,
  moveTaskOrderId: moveId,
  customerRemarks: 'mock remarks',
  requestedPickupDate: tomorrow,
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
    city: 'Auburn',
    state: 'WA',
    postalCode: '98001',
  },
  shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
};

const mockMtoShipmentSecondaryAddress = {
  id: uuidv4(),
  eTag: window.btoa(updatedAt),
  createdAt: '2021-06-11T18:12:11.918Z',
  updatedAt,
  moveTaskOrderId: moveId,
  customerRemarks: 'mock remarks',
  requestedPickupDate: tomorrow,
  requestedDeliveryDate: '2021-08-11',
  hasSecondaryPickupAddress: true,
  hasSecondaryDeliveryAddress: true,
  secondaryPickupAddress: {
    id: uuidv4(),
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    id: uuidv4(),
    streetAddress1: '123 SW Main',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  secondaryDeliveryAddress: {
    id: uuidv4(),
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
  shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
};

const reviewPath = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId });

beforeEach(() => {
  isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
});

afterEach(() => {
  jest.clearAllMocks();
});

const mockStore = configureStore({});

const renderMtoShipmentForm = (props) => {
  return renderWithRouter(
    <Provider store={mockStore.store}>
      <MtoShipmentForm {...defaultProps} {...props} />
    </Provider>,
    {
      path: customerRoutes.SHIPMENT_CREATE_PATH,
      params: { moveId },
    },
  );
};

const renderUBShipmentForm = (props) => {
  return renderWithRouter(
    <Provider store={mockStore.store}>
      <MtoShipmentForm {...ubProps} {...props} />
    </Provider>,
    {
      path: customerRoutes.SHIPMENT_CREATE_PATH,
      params: { moveId },
    },
  );
};

describe('MtoShipmentForm component', () => {
  describe('when creating a new HHG shipment and asterisks for required fields', () => {
    it('renders the HHG shipment form', async () => {
      renderMtoShipmentForm();

      expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');

      expect(screen.getByText(/5,000 lbs/)).toHaveClass('usa-alert__text');

      expect(screen.getAllByText('Date')[0]).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText(/Preferred pickup date/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Pickup info' })).toBeInTheDocument();
      expect(screen.getByTestId('pickupDateHint')).toHaveTextContent(
        'This is the day movers would put this shipment on their truck. Packing starts earlier. Dates will be finalized when you talk to your Customer Care Representative. Your requested pickup/load date should be your latest preferred pickup/load date, or the date you need to be out of your origin residence.',
      );
      expect(screen.getByText('Pickup Address')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 1 */)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 1 */)).toBeRequired();
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('City')).toBeInstanceOf(HTMLLabelElement);
      expect(screen.getByTestId('State')).toBeInstanceOf(HTMLLabelElement);
      expect(screen.getByTestId('ZIP')).toBeInstanceOf(HTMLLabelElement);
      expect(screen.getByLabelText(/Location Lookup */)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('Yes, I have a second pickup address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not have a second pickup address')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText(/First name/)[0]).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.getAllByLabelText(/Last name/)[0]).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.getAllByLabelText(/Phone/)[0]).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.getAllByLabelText(/Email/)[0]).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.getAllByText('Date')[1]).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText(/Preferred delivery date/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Delivery Address')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not know my delivery address')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByRole('heading', { level: 4, name: 'Second Delivery Address' })).not.toBeInTheDocument();
      expect(screen.queryByTitle('Yes, I have a second delivery address')).not.toBeInTheDocument();
      expect(screen.queryByTitle('No, I do not have a second delivery address')).not.toBeInTheDocument();

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText(/First name/)[1]).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getAllByLabelText(/Last name/)[1]).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.getAllByLabelText(/Phone/)[1]).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.getAllByLabelText(/Email/)[1]).toHaveAttribute('name', 'delivery.agent.email');

      expect(
        screen.queryByText(
          'Details about the facility where your things are now, including the name or address (if you know them)',
        ),
      ).not.toBeInTheDocument();

      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('renders the correct weight allowance when there are dependents', async () => {
      renderMtoShipmentForm({ orders: { has_dependents: true, authorizedWeight: 8000 } });

      expect(await screen.findByText('HHG')).toHaveClass('usa-tag');

      expect(screen.getByText(/8,000 lbs/)).toHaveClass('usa-alert__text');
    });

    it('renders the correct helper text for Delivery Address when orders type is RETIREMENT', async () => {
      renderMtoShipmentForm({ orders: { orders_type: ORDERS_TYPE.RETIREMENT } });
      await waitFor(() =>
        expect(
          screen.getByText('We can use the zip of the HOR, PLEAD or HOS you entered with your orders.')
            .toBeInTheDocument,
        ),
      );
    });

    it('renders the correct helper text for Delivery Address when orders type is SEPARATION', async () => {
      renderMtoShipmentForm({ orders: { orders_type: ORDERS_TYPE.SEPARATION } });
      await waitFor(() =>
        expect(
          screen.getByText('We can use the zip of the HOR, PLEAD or HOS you entered with your orders.')
            .toBeInTheDocument,
        ),
      );
    });

    it('renders the correct helper text for Delivery Address when orders type is PERMANENT_CHANGE_OF_STATION', async () => {
      renderMtoShipmentForm({ orders: { orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION } });
      await waitFor(() => expect(screen.getByText(/We can use the zip of your new duty location./).toBeInTheDocument));
    });

    it('renders the correct helper text for Delivery Address when orders type is LOCAL_MOVE', async () => {
      renderMtoShipmentForm({ orders: { orders_type: ORDERS_TYPE.LOCAL_MOVE } });
      await waitFor(() => expect(screen.getByText(/We can use the zip of your new duty location./).toBeInTheDocument));
    });

    it('renders the correct helper text for Delivery Address when orders type is TEMPORARY_DUTY', async () => {
      renderMtoShipmentForm({ orders: { orders_type: ORDERS_TYPE.TEMPORARY_DUTY } });
      await waitFor(() => expect(screen.getByText(/We can use the zip of your new duty location./).toBeInTheDocument));
    });

    it('does not render special NTS What to expect section', async () => {
      const { queryByTestId } = renderMtoShipmentForm();

      await waitFor(() => {
        expect(queryByTestId('nts-what-to-expect')).not.toBeInTheDocument();
      });
    });

    it('uses the current residence address for pickup address when checked', async () => {
      const { queryByLabelText, queryAllByLabelText, getAllByTestId } = renderMtoShipmentForm();

      await userEvent.click(queryByLabelText('Use my current address'));

      await waitFor(() => {
        expect(queryAllByLabelText(/Address 1/)[0]).toHaveValue(defaultProps.currentResidence.streetAddress1);
        expect(queryAllByLabelText(/Address 2/)[0]).toHaveValue('');
        expect(getAllByTestId('City')[0]).toHaveTextContent(defaultProps.currentResidence.city);
        expect(getAllByTestId(/State/)[0]).toHaveTextContent(defaultProps.currentResidence.state);
        expect(getAllByTestId(/ZIP/)[0]).toHaveTextContent(defaultProps.currentResidence.postalCode);
        expect(
          screen.getAllByText(
            `${defaultProps.currentResidence.city}, ${defaultProps.currentResidence.state} ${defaultProps.currentResidence.postalCode} ()`,
          ),
        );
      });
    });

    it('renders a second address fieldset when the user has a second pickup address', async () => {
      renderMtoShipmentForm();

      fireEvent.click(screen.getByTestId('has-secondary-pickup'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1.length).toBe(2);
      expect(streetAddress1[1]).toHaveAttribute('name', 'secondaryPickup.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[1]).toHaveAttribute('name', 'secondaryPickup.address.streetAddress2');

      const city = screen.getAllByTestId('City');
      expect(city[0]).toHaveAttribute('aria-label', 'pickup.address.city');

      const state = screen.getAllByTestId(/State/);
      expect(state[0]).toHaveAttribute('aria-label', 'pickup.address.state');

      const zip = screen.getAllByTestId(/ZIP/);
      expect(zip[0]).toHaveAttribute('aria-label', 'pickup.address.postalCode');
      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(2);
    });

    it('renders a third address fieldset when the user has a third pickup address', async () => {
      const { queryByLabelText } = renderMtoShipmentForm({ mtoShipment: mockMtoShipmentSecondaryAddress });

      await userEvent.click(queryByLabelText('Use my current address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));
      await userEvent.click(screen.getByTitle('Yes, I have a third pickup address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1[2]).toHaveAttribute('name', 'tertiaryPickup.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[2]).toHaveAttribute('name', 'tertiaryPickup.address.streetAddress2');

      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(5);
    });

    it('renders an address fieldset when the user has a delivery address', async () => {
      renderMtoShipmentForm();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1[1]).toHaveAttribute('name', 'delivery.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[1]).toHaveAttribute('name', 'delivery.address.streetAddress2');

      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(2);
    });

    it('renders the secondary delivery address question once a user says they have a primary delivery address', async () => {
      renderMtoShipmentForm();

      expect(screen.queryByRole('heading', { level: 4, name: 'Second Delivery Address' })).not.toBeInTheDocument();
      expect(screen.queryByTitle('Yes, I have a second delivery address')).not.toBeInTheDocument();
      expect(screen.queryByTitle('No, I do not have a second delivery address')).not.toBeInTheDocument();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));

      expect(screen.getByTitle('Yes, I have a second delivery address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not have a second delivery address')).toBeInstanceOf(HTMLInputElement);
    });

    it('renders a second address fieldset when the user has a second delivery address', async () => {
      renderMtoShipmentForm({ mtoShipment: mockMtoShipmentHHGWithDest });

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1[2]).toHaveAttribute('name', 'secondaryDelivery.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[2]).toHaveAttribute('name', 'secondaryDelivery.address.streetAddress2');

      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(3);
    });

    it('renders the third delivery address question once a user says they have a second delivery address', async () => {
      renderMtoShipmentForm({ mtoShipment: mockMtoShipmentHHGWithDest });

      expect(screen.queryByRole('heading', { level: 4, name: 'Third Delivery Address' })).not.toBeInTheDocument();
      expect(screen.queryByTitle('Yes, I have a third delivery address')).not.toBeInTheDocument();
      expect(screen.queryByTitle('No, I do not have a third delivery address')).not.toBeInTheDocument();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      expect(screen.getByTitle('Yes, I have a third delivery address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not have a third delivery address')).toBeInstanceOf(HTMLInputElement);
    });

    it('goes back when the back button is clicked', async () => {
      renderMtoShipmentForm();

      const backButton = await screen.findByRole('button', { name: 'Back' });
      await userEvent.click(backButton);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith(-1);
      });
    });

    it('renders a third address fieldset when the user has a third delivery address', async () => {
      renderMtoShipmentForm({ mtoShipment: mockMtoShipmentHHGWithDest });

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a third delivery address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1[3]).toHaveAttribute('name', 'tertiaryDelivery.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[3]).toHaveAttribute('name', 'tertiaryDelivery.address.streetAddress2');

      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(4);
    });

    it('can submit a new HHG shipment successfully', async () => {
      const mockMtoShipmentHHG = {
        id: uuidv4(),
        eTag: window.btoa(updatedAt),
        createdAt: '2021-06-11T18:12:11.918Z',
        updatedAt,
        moveTaskOrderId: moveId,
        requestedPickupDate: tomorrow,
        requestedDeliveryDate: '2021-06-14',
        pickupAddress: {
          streetAddress1: '812 S 129th St',
          streetAddress2: '#123',
          city: 'San Antonio',
          state: 'TX',
          postalCode: '78234',
        },
        shipmentType: SHIPMENT_OPTIONS.HHG,
        hasSecondaryPickupAddress: false,
        hasSecondaryDeliveryAddress: false,
        hasTertiaryPickupAddress: false,
        hasTertiaryDeliveryAddress: false,
      };

      const expectedPayload = {
        agents: [
          { agentType: 'RELEASING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
          { agentType: 'RECEIVING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
        ],
        moveTaskOrderID: moveId,
        shipmentType: SHIPMENT_OPTIONS.HHG,
        customerRemarks: '',
        requestedPickupDate: tomorrow,
        pickupAddress: { ...mockMtoShipmentHHG.pickupAddress },
        requestedDeliveryDate: '2021-06-14',
        hasSecondaryPickupAddress: false,
        hasSecondaryDeliveryAddress: false,
        hasTertiaryPickupAddress: false,
        hasTertiaryDeliveryAddress: false,
      };

      const expectedCreateResponse = {
        createdAt: '2021-06-11T18:12:11.918Z',
        customerRemarks: '',
        eTag: window.btoa(updatedAt),
        id: uuidv4(),
        moveTaskOrderID: moveId,
        pickupAddress: { ...mockMtoShipmentHHG.pickupAddress, id: uuidv4() },
        requestedDeliveryDate: expectedPayload.requestedDeliveryDate,
        requestedPickupDate: expectedPayload.requestedPickupDate,
        shipmentType: SHIPMENT_OPTIONS.HHG,
        status: 'SUBMITTED',
        updatedAt,
      };

      createMTOShipment.mockImplementation(() => Promise.resolve(expectedCreateResponse));
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ mtoShipment: mockMtoShipmentHHG });

      const nextButton = await screen.findByRole('button', { name: 'Next' });
      expect(nextButton).not.toBeDisabled();
      await userEvent.click(nextButton);

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith(expectedPayload);
      });

      expect(defaultProps.updateMTOShipment).toHaveBeenCalledWith(expectedCreateResponse);

      expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
    });

    it('disables the continue button when the move is locked by an office user', async () => {
      const mockMtoShipmentHHG = {
        id: uuidv4(),
        eTag: window.btoa(updatedAt),
        createdAt: '2021-06-11T18:12:11.918Z',
        updatedAt,
        moveTaskOrderId: moveId,
        requestedPickupDate: '2021-06-07',
        requestedDeliveryDate: '2021-06-14',
        pickupAddress: {
          streetAddress1: '812 S 129th St',
          streetAddress2: '#123',
          city: 'San Antonio',
          state: 'TX',
          postalCode: '78234',
        },
        shipmentType: SHIPMENT_OPTIONS.HHG,
        hasSecondaryPickupAddress: false,
        hasSecondaryDeliveryAddress: false,
        hasTertiaryPickupAddress: false,
        hasTertiaryDeliveryAddress: false,
      };

      renderMtoShipmentForm({ mtoShipment: mockMtoShipmentHHG, isMoveLocked: true });

      const nextButton = await screen.findByRole('button', { name: 'Next' });
      expect(nextButton).toBeDisabled();
    });

    it('shows an error when there is an error with the submission', async () => {
      const shipmentInfo = {
        requestedPickupDate: tomorrow,
        pickupAddress: {
          streetAddress1: '812 S 129th St',
          streetAddress2: '#123',
          city: 'San Antonio',
          state: 'TX',
          postalCode: '78234',
        },
        requestedDeliveryDate: '14 Jun 2021',
      };

      const errorMessage = 'Something broke!';
      const errorResponse = { response: { errorMessage } };
      createMTOShipment.mockImplementation(() => Promise.reject(errorResponse));
      getResponseError.mockImplementation(() => errorMessage);
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ mtoShipment: shipmentInfo });

      const nextButton = await screen.findByRole('button', { name: 'Next' });
      expect(nextButton).not.toBeDisabled();
      await userEvent.click(nextButton);

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalled();
      });

      expect(getResponseError).toHaveBeenCalledWith(
        errorResponse.response,
        'failed to create MTO shipment due to server error',
      );

      expect(await screen.findByText(errorMessage)).toBeInTheDocument();
    });
  });

  describe('editing an already existing HHG shipment', () => {
    const mockMtoShipment = {
      id: uuidv4(),
      eTag: window.btoa(updatedAt),
      createdAt: '2021-06-11T18:12:11.918Z',
      updatedAt,
      moveTaskOrderId: moveId,
      customerRemarks: 'mock remarks',
      requestedPickupDate: tomorrow,
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

    it('renders the HHG shipment form with pre-filled values', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText('Use my current address')).not.toBeChecked();
      expect(screen.getAllByLabelText(/Address 1 */)[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getByText('San Antonio, TX 78234 ()'));
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText(/Address 1 */)[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getByText('Tacoma, WA 98421 ()'));
      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');

      expect(
        screen.getByText(
          `Preferred pickup date ${tomorrowDatePicker} is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
        ),
      ).toHaveClass('usa-alert__text');
      expect(
        screen.getAllByText(
          `Preferred pickup date ${tomorrowDatePicker} is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
        ),
      ).toHaveLength(1);
    });

    it('renders the HHG shipment form with pre-filled secondary addresses', async () => {
      const shipment = {
        ...mockMtoShipment,
        secondaryPickupAddress: {
          streetAddress1: '142 E Barrel Hoop Circle',
          streetAddress2: '#4A',
          city: 'Corpus Christi',
          state: 'TX',
          postalCode: '78412',
        },
        secondaryDeliveryAddress: {
          streetAddress1: '3373 NW Martin Luther King Jr Blvd',
          streetAddress2: '',
          city: 'Auburn',
          state: 'WA',
          postalCode: '98002',
        },
      };
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: shipment });

      expect(await screen.findByTitle('Yes, I have a second pickup address')).toBeChecked();
      expect(await screen.findByTitle('Yes, I have a second delivery address')).toBeChecked();

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1.length).toBe(4);

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2.length).toBe(4);

      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(4);

      // Secondary pickup address should be the 2nd address
      expect(streetAddress1[1]).toHaveValue('142 E Barrel Hoop Circle');
      expect(streetAddress2[1]).toHaveValue('#4A');
      expect(screen.getByText('Corpus Christi, TX 78412 ()'));

      // Secondary delivery address should be the 4th address
      expect(streetAddress1[3]).toHaveValue('3373 NW Martin Luther King Jr Blvd');
      expect(streetAddress2[3]).toHaveValue('');
      expect(screen.getByText('Auburn, WA 98002 ()'));
    });

    it('does not allow the user to save the form if the address fields on a secondary addreess is the only one filled out', async () => {
      const shipment = {
        ...mockMtoShipment,
        secondaryPickupAddress: {
          streetAddress1: '142 E Barrel Hoop Circle',
          streetAddress2: '#4A',
          city: 'Corpus Christi',
          state: 'TX',
          postalCode: '78412',
        },
        secondaryDeliveryAddress: {
          streetAddress1: '3373 NW Martin Luther King Jr Blvd',
          streetAddress2: '',
          city: mockMtoShipment.destinationAddress.city,
          state: mockMtoShipment.destinationAddress.state,
          postalCode: mockMtoShipment.destinationAddress.postalCode,
        },
      };
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: shipment });

      // Verify that the form is good to submit by checking that the save button is not disabled.
      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      const address = await screen.findAllByLabelText(/Address 1 */);

      // The second instance of a field is the secondary pickup
      await userEvent.type(address[1], '6622 Airport Way S');
      await waitFor(() => {
        expect(saveButton).not.toBeDisabled();
      });

      // Clear the field so that the secondary delivery address can be checked
      await userEvent.clear(address[1]);
      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });
    });

    it('disables the save button if the move has been locked by an office user', async () => {
      const shipment = {
        ...mockMtoShipment,
        secondaryPickupAddress: {
          streetAddress1: '142 E Barrel Hoop Circle',
          streetAddress2: '#4A',
          city: 'Corpus Christi',
          state: 'TX',
          postalCode: '78412',
        },
        secondaryDeliveryAddress: {
          streetAddress1: '3373 NW Martin Luther King Jr Blvd',
          streetAddress2: '',
          city: mockMtoShipment.destinationAddress.city,
          state: mockMtoShipment.destinationAddress.state,
          postalCode: mockMtoShipment.destinationAddress.postalCode,
        },
      };

      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: shipment, isMoveLocked: true });

      // Verify that the form is good to submit by checking that the save button is not disabled.
      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).toBeDisabled();

      // test to verify that the MOVE_LOCKED_WARNING text is showing:
      expect(screen.getByText(MOVE_LOCKED_WARNING)).toBeInTheDocument();
      expect(screen.getByText(MOVE_LOCKED_WARNING)).toBeVisible();
    });

    it('allow the user to save the form if the secondary address1 field is cleared but the toggle is switched to No', async () => {
      const shipment = {
        ...mockMtoShipment,
        secondaryPickupAddress: {
          streetAddress1: '142 E Barrel Hoop Circle',
          streetAddress2: '#4A',
          city: 'Corpus Christi',
          state: 'TX',
          postalCode: '78412',
        },
        secondaryDeliveryAddress: {
          streetAddress1: '3373 NW Martin Luther King Jr Blvd',
          streetAddress2: '',
          city: mockMtoShipment.destinationAddress.city,
          state: mockMtoShipment.destinationAddress.state,
          postalCode: mockMtoShipment.destinationAddress.postalCode,
        },
      };

      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: shipment });

      // Verify that the form is good to submit by checking that the save button is not disabled.
      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));
      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      const address = await screen.findAllByLabelText(/Address 1 */);

      // Clear the second pickup address1 field so that it triggers required validation, disables Save
      await userEvent.clear(address[1]);
      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });

      //  Click No to second pickup address, should be able to save
      await userEvent.click(screen.getByTitle('No, I do not have a second pickup address'));
      await waitFor(() => {
        expect(saveButton).not.toBeDisabled();
      });

      // get new address1 pull since disabled above
      const newAddress = await screen.findAllByLabelText(/Address 1 */);
      // Clear the second delivery address1 field so that it triggers required validation, disables Save
      await userEvent.clear(newAddress[2]);
      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });

      //  Click No to second delivery address, should be able to save
      await userEvent.click(screen.getByTitle('No, I do not have a second delivery address'));
      await waitFor(() => {
        expect(saveButton).not.toBeDisabled();
      });
    });

    it('allow the user to save the form if the tertiary address1 field is cleared but the toggle is switched to No', async () => {
      const shipment = {
        ...mockMtoShipment,
        secondaryPickupAddress: {
          streetAddress1: '142 E Barrel Hoop Circle',
          streetAddress2: '#4A',
          city: 'Corpus Christi',
          state: 'TX',
          postalCode: '78412',
        },
        tertiaryPickupAddress: {
          streetAddress1: '789 S Elm',
          city: 'Corpus Christi',
          state: 'TX',
          postalCode: '78412',
        },
        secondaryDeliveryAddress: {
          streetAddress1: '3373 NW Martin Luther King Jr Blvd',
          streetAddress2: '',
          city: mockMtoShipment.destinationAddress.city,
          state: mockMtoShipment.destinationAddress.state,
          postalCode: mockMtoShipment.destinationAddress.postalCode,
        },
        tertiaryDeliveryAddress: {
          streetAddress1: '453 N Main Blvd',
          city: mockMtoShipment.destinationAddress.city,
          state: mockMtoShipment.destinationAddress.state,
          postalCode: mockMtoShipment.destinationAddress.postalCode,
        },
      };

      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: shipment });

      // Verify that the form is good to submit by checking that the save button is not disabled.
      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));
      await userEvent.click(screen.getByTitle('Yes, I have a third pickup address'));
      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a third delivery address'));

      const address = await screen.findAllByLabelText(/Address 1 */);

      // Clear the third pickup address1 field so that it triggers required validation, disables Save
      await userEvent.clear(address[2]);
      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });

      //  Click No to third pickup address, should be able to save
      await userEvent.click(screen.getByTitle('No, I do not have a third pickup address'));
      await waitFor(() => {
        expect(saveButton).not.toBeDisabled();
      });

      // get new address1 pull since disabled above
      const newAddress = await screen.findAllByLabelText(/Address 1 */);
      // Clear the third delivery address1 field so that it triggers required validation, disables Save
      await userEvent.clear(newAddress[4]);
      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });

      //  Click No to third delivery address, should be able to save
      await userEvent.click(screen.getByTitle('No, I do not have a third delivery address'));
      await waitFor(() => {
        expect(saveButton).not.toBeDisabled();
      });
    });

    it('goes back when the cancel button is clicked', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

      const cancelButton = await screen.findByRole('button', { name: 'Cancel' });
      await userEvent.click(cancelButton);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith(-1);
      });
    });

    it('can submit edits to an HHG shipment successfully', async () => {
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
        moveTaskOrderID: moveId,
        shipmentType: SHIPMENT_OPTIONS.HHG,
        pickupAddress: { ...shipmentInfo.pickupAddress, city: 'San Antonio', state: 'TX', postalCode: '78234' },
        customerRemarks: mockMtoShipmentUB.customerRemarks,
        requestedPickupDate: mockMtoShipmentUB.requestedPickupDate,
        requestedDeliveryDate: mockMtoShipmentUB.requestedDeliveryDate,
        destinationAddress: {
          ...mockMtoShipmentUB.destinationAddress,
          city: 'Tacoma',
          state: 'WA',
          postalCode: '98421',
          streetAddress2: '',
        },
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
        ...mockMtoShipment,
        pickupAddress: { ...shipmentInfo.pickupAddress },
        shipmentType: SHIPMENT_OPTIONS.HHG,
        eTag: window.btoa(newUpdatedAt),
        status: 'SUBMITTED',
      };

      patchMTOShipment.mockImplementation(() => Promise.resolve(expectedUpdateResponse));
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

      const pickupAddress1Input = screen.getAllByLabelText(/Address 1 */)[0];
      await userEvent.clear(pickupAddress1Input);
      await userEvent.type(pickupAddress1Input, shipmentInfo.pickupAddress.streetAddress1);

      const pickupAddress2Input = screen.getAllByLabelText(/Address 2/)[0];
      await userEvent.clear(pickupAddress2Input);
      await userEvent.type(pickupAddress2Input, shipmentInfo.pickupAddress.streetAddress2);

      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalledWith(mockMtoShipment.id, expectedPayload, mockMtoShipment.eTag);
      });

      expect(defaultProps.updateMTOShipment).toHaveBeenCalledWith(expectedUpdateResponse);

      expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
    });

    it('collapses delivery address fieldset when the user says No to having a delivery address', async () => {
      const blankSecondaryDelivery = {
        secondaryDeliveryAddress: {
          streetAddress1: '',
          streetAddress2: '',
          city: '',
          state: '',
          postalCode: '',
        },
      };

      const newUpdatedAt = '2021-06-11T21:20:22.150Z';
      const expectedUpdateResponse = {
        ...mockMtoShipmentHHGWithDest,
        destinationAddress: { ...mockMtoShipmentHHGWithDest.newDutyLocationAddress, streetAddress1: 'N/A' },
        secondaryDeliveryAddress: blankSecondaryDelivery,
        shipmentType: SHIPMENT_OPTIONS.HHG,
        eTag: window.btoa(newUpdatedAt),
        status: 'SUBMITTED',
      };

      patchMTOShipment.mockImplementation(() => Promise.resolve(expectedUpdateResponse));

      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentHHGWithDest });

      // verify second delivery has values
      expect(await screen.getAllByLabelText(/Address 1 */)[2]).toHaveValue(
        mockMtoShipmentHHGWithDest.secondaryDeliveryAddress.streetAddress1,
      );
      expect(
        screen.getAllByText(
          `${mockMtoShipmentHHGWithDest.secondaryDeliveryAddress.city}, ${mockMtoShipmentHHGWithDest.secondaryDeliveryAddress.state} ${mockMtoShipmentHHGWithDest.secondaryDeliveryAddress.postalCode} ()`,
        ),
      );

      await userEvent.click(screen.getByTitle('No, I do not know my delivery address'));

      // No to delivery should also hide second delivery address fields, only pickup address left
      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1.length).toBe(1);

      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(1);

      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      expect(defaultProps.updateMTOShipment).toHaveBeenCalledWith(expectedUpdateResponse);

      expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
    });

    it('shows an error when there is an error with the submission', async () => {
      const errorMessage = 'Something broke!';
      const errorResponse = { response: { errorMessage } };
      patchMTOShipment.mockImplementation(() => Promise.reject(errorResponse));
      getResponseError.mockImplementation(() => errorMessage);
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalled();
      });

      expect(getResponseError).toHaveBeenCalledWith(
        errorResponse.response,
        'failed to update MTO shipment due to server error',
      );

      expect(await screen.findByText(errorMessage)).toBeInTheDocument();
    });

    it('renders the HHG shipment form with pre-filled values', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText('Use my current address')).not.toBeChecked();
      expect(screen.getAllByLabelText(/Address 1 */)[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByTestId('City')[0]).toHaveTextContent('San Antonio');
      expect(screen.getAllByTestId('State')[0]).toHaveTextContent('TX');
      expect(screen.getAllByTestId('ZIP')[0]).toHaveTextContent('78234');
      expect(screen.getByText('San Antonio, TX 78234 ()'));
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText(/Address 1 */)[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByTestId('City')[1]).toHaveTextContent('Tacoma');
      expect(screen.getAllByTestId('State')[1]).toHaveTextContent('WA');
      expect(screen.getAllByTestId('ZIP')[1]).toHaveTextContent('98421');
      expect(screen.getByText('Tacoma, WA 98421 ()'));
      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');
    });

    it('renders the HHG shipment with date validaton alerts for weekend and holiday', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United of States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a holiday and weekend in the United of States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Preferred delivery date 11 Aug 2021 is on a holiday and weekend in the United of States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('renders the HHG shipment with date validaton alerts for weekend', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Preferred delivery date 11 Aug 2021 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('renders the HHG shipment with date validaton alerts for holiday', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Preferred delivery date 11 Aug 2021 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('renders the HHG shipment with no date validaton alerts for pickup/delivery', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');

      await waitFor(() => {
        expect(
          screen.queryAllByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveLength(0);
        expect(
          screen.queryAllByText(
            'Preferred delivery date 11 Aug 2021 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
          ),
        ).toHaveLength(0);
      });
    });
  });

  describe('when creating a new UB shipment', () => {
    it('renders the UB shipment form', async () => {
      renderUBShipmentForm();

      expect(await screen.findByText('UB')).toHaveClass('usa-tag');

      expect(
        screen.queryByText(
          'Remember: You can move up to 600 lbs for this UB shipment. The weight of your UB is part of your authorized weight allowance. You’ll be billed for any excess weight you move.',
        ),
      ).toBeInTheDocument();

      expect(screen.getAllByText('Date')[0]).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText(/Preferred pickup date/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Pickup info' })).toBeInTheDocument();
      expect(screen.getByTestId('pickupDateHint')).toHaveTextContent(
        'This is the day movers would put this shipment on their truck. Packing starts earlier. Dates will be finalized when you talk to your Customer Care Representative. Your requested pickup/load date should be your latest preferred pickup/load date, or the date you need to be out of your origin residence.',
      );
      expect(screen.getByText('Pickup Address')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 1 */)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Location Lookup */)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByTitle('Yes, I have a second pickup address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not have a second pickup address')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText(/First name/)[0]).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.getAllByLabelText(/Last name/)[0]).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.getAllByLabelText(/Phone/)[0]).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.getAllByLabelText(/Email/)[0]).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.getAllByText('Date')[1]).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText(/Preferred delivery date/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Delivery Address')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not know my delivery address')).toBeInstanceOf(HTMLInputElement);

      expect(screen.queryByRole('heading', { level: 4, name: 'Second Delivery Address' })).not.toBeInTheDocument();
      expect(screen.queryByTitle('Yes, I have a second delivery address')).not.toBeInTheDocument();
      expect(screen.queryByTitle('No, I do not have a second delivery address')).not.toBeInTheDocument();

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getAllByLabelText(/First name/)[1]).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getAllByLabelText(/Last name/)[1]).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.getAllByLabelText(/Phone/)[1]).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.getAllByLabelText(/Email/)[1]).toHaveAttribute('name', 'delivery.agent.email');

      expect(
        screen.queryByText(
          'Details about the facility where your things are now, including the name or address (if you know them)',
        ),
      ).not.toBeInTheDocument();

      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('renders the correct helper text when the UB allowance is null', async () => {
      renderUBShipmentForm({ orders: { entitlement: { ub_allowance: null } } });
      expect(
        screen.queryByText(
          'Remember: You can move up to your UB allowance for this UB shipment. The weight of your UB is part of your authorized weight allowance. You’ll be billed for any excess weight you move.',
        ),
      ).toBeInTheDocument();
    });

    it('renders the correct helper text for Delivery Address when orders type is RETIREMENT', async () => {
      renderUBShipmentForm({ orders: { orders_type: ORDERS_TYPE.RETIREMENT } });
      await waitFor(() =>
        expect(
          screen.getByText('We can use the zip of the HOR, PLEAD or HOS you entered with your orders.')
            .toBeInTheDocument,
        ),
      );
    });

    it('renders the correct helper text for Delivery Address when orders type is SEPARATION', async () => {
      renderUBShipmentForm({ orders: { orders_type: ORDERS_TYPE.SEPARATION } });
      await waitFor(() =>
        expect(
          screen.getByText('We can use the zip of the HOR, PLEAD or HOS you entered with your orders.')
            .toBeInTheDocument,
        ),
      );
    });

    it('renders the correct helper text for Delivery Address when orders type is PERMANENT_CHANGE_OF_STATION', async () => {
      renderUBShipmentForm({ orders: { orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION } });
      await waitFor(() => expect(screen.getByText(/We can use the zip of your new duty location./).toBeInTheDocument));
    });

    it('renders the correct helper text for Delivery Address when orders type is LOCAL_MOVE', async () => {
      renderUBShipmentForm({ orders: { orders_type: ORDERS_TYPE.LOCAL_MOVE } });
      await waitFor(() => expect(screen.getByText(/We can use the zip of your new duty location./).toBeInTheDocument));
    });

    it('renders the correct helper text for Delivery Address when orders type is TEMPORARY_DUTY', async () => {
      renderUBShipmentForm({ orders: { orders_type: ORDERS_TYPE.TEMPORARY_DUTY } });
      await waitFor(() => expect(screen.getByText(/We can use the zip of your new duty location./).toBeInTheDocument));
    });

    it('does not render special NTS What to expect section', async () => {
      const { queryByTestId } = renderUBShipmentForm();

      await waitFor(() => {
        expect(queryByTestId('nts-what-to-expect')).not.toBeInTheDocument();
      });
    });

    it('uses the current residence address for pickup address when checked', async () => {
      const { queryByLabelText, queryAllByLabelText, getAllByTestId } = renderUBShipmentForm();

      await userEvent.click(queryByLabelText('Use my current address'));

      await waitFor(() => {
        expect(queryAllByLabelText(/Address 1 */)[0]).toHaveValue(defaultProps.currentResidence.streetAddress1);
        expect(queryAllByLabelText(/Address 2/)[0]).toHaveValue('');
        expect(getAllByTestId('City')[0]).toHaveTextContent(defaultProps.currentResidence.city);
        expect(getAllByTestId('State')[0]).toHaveTextContent(defaultProps.currentResidence.state);
        expect(getAllByTestId('ZIP')[0]).toHaveTextContent(defaultProps.currentResidence.postalCode);
        expect(
          screen.getAllByText(
            `${defaultProps.currentResidence.city}, ${defaultProps.currentResidence.state} ${defaultProps.currentResidence.postalCode} ()`,
          ),
        );
      });
    });

    it('renders a second address fieldset when the user has a pickup address', async () => {
      renderUBShipmentForm();

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1[0]).toHaveAttribute('name', 'pickup.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[0]).toHaveAttribute('name', 'pickup.address.streetAddress2');

      const city = screen.getAllByTestId('City');
      expect(city[0]).toHaveAttribute('aria-label', 'pickup.address.city');

      const state = screen.getAllByTestId('State');
      expect(state[0]).toHaveAttribute('aria-label', 'pickup.address.state');

      const zip = screen.getAllByTestId('ZIP');
      expect(zip[0]).toHaveAttribute('aria-label', 'pickup.address.postalCode');

      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(1);
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      renderUBShipmentForm();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1[0]).toHaveAttribute('name', 'pickup.address.streetAddress1');
      expect(streetAddress1[1]).toHaveAttribute('name', 'delivery.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[0]).toHaveAttribute('name', 'pickup.address.streetAddress2');
      expect(streetAddress2[1]).toHaveAttribute('name', 'delivery.address.streetAddress2');

      const city = screen.getAllByTestId('City');
      expect(city[0]).toHaveAttribute('aria-label', 'pickup.address.city');
      expect(city[1]).toHaveAttribute('aria-label', 'delivery.address.city');

      const state = screen.getAllByTestId('State');
      expect(state[0]).toHaveAttribute('aria-label', 'pickup.address.state');
      expect(state[1]).toHaveAttribute('aria-label', 'delivery.address.state');

      const zip = screen.getAllByTestId('ZIP');
      expect(zip[0]).toHaveAttribute('aria-label', 'pickup.address.postalCode');
      expect(zip[1]).toHaveAttribute('aria-label', 'delivery.address.postalCode');

      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(2);
    });

    it('renders the secondary delivery address question once a user says they have a primary delivery address', async () => {
      renderUBShipmentForm();

      expect(screen.queryByRole('heading', { level: 4, name: 'Second Delivery Address' })).not.toBeInTheDocument();
      expect(screen.queryByTitle('Yes, I have a second delivery address')).not.toBeInTheDocument();
      expect(screen.queryByTitle('No, I do not have a second delivery address')).not.toBeInTheDocument();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));
      expect(screen.getByTitle('Yes, I have a second delivery address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not have a second delivery address')).toBeInstanceOf(HTMLInputElement);
    });

    it('renders another address fieldset when the user has a second delivery address', async () => {
      renderUBShipmentForm({ mtoShipment: mockMtoShipmentUB });

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1.length).toBe(3);
      expect(streetAddress1[2]).toHaveAttribute('name', 'secondaryDelivery.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2.length).toBe(3);
      expect(streetAddress2[2]).toHaveAttribute('name', 'secondaryDelivery.address.streetAddress2');

      const city = screen.getAllByTestId('City');
      expect(city.length).toBe(3);
      expect(city[2]).toHaveAttribute('aria-label', 'secondaryDelivery.address.city');

      const state = screen.getAllByTestId('State');
      expect(state.length).toBe(3);
      expect(state[2]).toHaveAttribute('aria-label', 'secondaryDelivery.address.state');

      const zip = screen.getAllByTestId('ZIP');
      expect(zip.length).toBe(3);
      expect(zip[2]).toHaveAttribute('aria-label', 'secondaryDelivery.address.postalCode');
      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(3);
    });

    it('goes back when the back button is clicked', async () => {
      renderUBShipmentForm();

      const backButton = await screen.findByRole('button', { name: 'Back' });
      await userEvent.click(backButton);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith(-1);
      });
    });

    it('can submit a new UB shipment successfully', async () => {
      const expectedPayload = {
        moveTaskOrderID: moveId,
        shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
        pickupAddress: { ...mockMtoShipmentUB.pickupAddress, streetAddress2: '' },
        customerRemarks: mockMtoShipmentUB.customerRemarks,
        requestedPickupDate: mockMtoShipmentUB.requestedPickupDate,
        requestedDeliveryDate: mockMtoShipmentUB.requestedDeliveryDate,
        destinationAddress: { ...mockMtoShipmentUB.destinationAddress, streetAddress2: '' },
        hasSecondaryDeliveryAddress: false,
        hasSecondaryPickupAddress: false,
        hasTertiaryDeliveryAddress: false,
        hasTertiaryPickupAddress: false,
        destinationType: undefined,
        sacType: undefined,
        tacType: undefined,
        agents: [
          { agentType: 'RELEASING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
          { agentType: 'RECEIVING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
        ],
        counselorRemarks: undefined,
      };
      delete expectedPayload.destinationAddress.id;
      delete expectedPayload.pickupAddress.id;

      const expectedCreateResponse = {
        createdAt: '2021-06-11T18:12:11.918Z',
        customerRemarks: '',
        eTag: window.btoa(updatedAt),
        id: uuidv4(),
        moveTaskOrderID: moveId,
        pickupAddress: { ...mockMtoShipmentUB.pickupAddress, id: uuidv4() },
        requestedDeliveryDate: expectedPayload.requestedDeliveryDate,
        requestedPickupDate: expectedPayload.requestedPickupDate,
        shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
        status: 'SUBMITTED',
        updatedAt,
      };

      createMTOShipment.mockImplementation(() => Promise.resolve(expectedCreateResponse));
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ mtoShipment: mockMtoShipmentUB });

      const nextButton = await screen.findByRole('button', { name: 'Next' });
      expect(nextButton).not.toBeDisabled();
      await userEvent.click(nextButton);

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith(expectedPayload);
      });

      expect(ubProps.updateMTOShipment).toHaveBeenCalledWith(expectedCreateResponse);

      expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
    });

    it('shows an error when there is an error with the submission', async () => {
      const errorMessage = 'Something broke!';
      const errorResponse = { response: { errorMessage } };
      createMTOShipment.mockImplementation(() => Promise.reject(errorResponse));
      getResponseError.mockImplementation(() => errorMessage);
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ mtoShipment: mockMtoShipmentUB });

      const nextButton = await screen.findByRole('button', { name: 'Next' });
      expect(nextButton).not.toBeDisabled();
      await userEvent.click(nextButton);

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalled();
      });

      expect(getResponseError).toHaveBeenCalledWith(
        errorResponse.response,
        'failed to create MTO shipment due to server error',
      );

      expect(await screen.findByText(errorMessage)).toBeInTheDocument();
    });
  });

  describe('editing an already existing UB shipment', () => {
    it('renders the UB shipment form with pre-filled values', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentUB });

      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText('Use my current address')).not.toBeChecked();
      expect(screen.getAllByLabelText(/Address 1 */)[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByTestId('City')[0]).toHaveTextContent('San Antonio');
      expect(screen.getAllByTestId('State')[0]).toHaveTextContent('TX');
      expect(screen.getAllByTestId('ZIP')[0]).toHaveTextContent('78234');
      expect(
        screen.getAllByText(
          `${mockMtoShipmentUB.pickupAddress.city}, ${mockMtoShipmentUB.pickupAddress.state} ${mockMtoShipmentUB.pickupAddress.postalCode} ()`,
        ),
      );
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText(/Address 1 */)[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByTestId('City')[1]).toHaveTextContent('Auburn');
      expect(screen.getAllByTestId('State')[1]).toHaveTextContent('WA');
      expect(screen.getAllByTestId('ZIP')[1]).toHaveTextContent('98001');
      expect(
        screen.getAllByText(
          `${mockMtoShipmentUB.destinationAddress.city}, ${mockMtoShipmentUB.destinationAddress.state} ${mockMtoShipmentUB.destinationAddress.postalCode} ()`,
        ),
      );
      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');

      expect(
        screen.getByText(
          `Preferred pickup date ${tomorrowDatePicker} is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
        ),
      ).toHaveClass('usa-alert__text');
      expect(
        screen.getAllByText(
          `Preferred pickup date ${tomorrowDatePicker} is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
        ),
      ).toHaveLength(1);
    });

    it('renders the UB shipment form with pre-filled secondary addresses', async () => {
      const shipment = {
        ...mockMtoShipmentUB,
        secondaryPickupAddress: {
          streetAddress1: '142 E Barrel Hoop Circle',
          streetAddress2: '#4A',
          city: 'Corpus Christi',
          state: 'TX',
          postalCode: '78412',
        },
        secondaryDeliveryAddress: {
          streetAddress1: '3373 NW Martin Luther King Jr Blvd',
          streetAddress2: '',
          city: 'Auburn',
          state: 'WA',
          postalCode: '98002',
        },
      };
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: shipment });

      expect(await screen.findByTitle('Yes, I have a second pickup address')).toBeChecked();
      expect(await screen.findByTitle('Yes, I have a second delivery address')).toBeChecked();

      const streetAddress1 = await screen.findAllByLabelText(/Address 1 */);
      expect(streetAddress1.length).toBe(4);

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2.length).toBe(4);

      const city = screen.getAllByTestId('City');
      expect(city.length).toBe(4);

      const state = screen.getAllByTestId('State');
      expect(state.length).toBe(4);

      const zip = screen.getAllByTestId('ZIP');
      expect(zip.length).toBe(4);
      expect(screen.getAllByLabelText(/Location Lookup */).length).toBe(4);

      // Secondary pickup address should be the 2nd address
      expect(streetAddress1[1]).toHaveValue('142 E Barrel Hoop Circle');
      expect(streetAddress2[1]).toHaveValue('#4A');
      expect(city[1]).toHaveTextContent('Corpus Christi');
      expect(state[1]).toHaveTextContent('TX');
      expect(zip[1]).toHaveTextContent('78412');
      expect(screen.getByText('Corpus Christi, TX 78412 ()'));

      // Secondary delivery address should be the 4th address
      expect(streetAddress1[3]).toHaveValue('3373 NW Martin Luther King Jr Blvd');
      expect(streetAddress2[3]).toHaveValue('');
      expect(city[3]).toHaveTextContent('Auburn');
      expect(state[3]).toHaveTextContent('WA');
      expect(zip[3]).toHaveTextContent('98002');
      expect(screen.getByText('Auburn, WA 98002 ()'));
    });

    it('does not allow the user to save the form if the secondary addreess is the only one filled out', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentSecondaryAddress });

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));

      // Verify that the form cannot submit by checking that the save button is disabled.
      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).toBeDisabled();
    });

    it('goes back when the cancel button is clicked', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentUB });

      const cancelButton = await screen.findByRole('button', { name: 'Cancel' });
      await userEvent.click(cancelButton);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith(-1);
      });
    });

    it('can submit edits to a UB shipment successfully', async () => {
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
        moveTaskOrderID: moveId,
        shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
        pickupAddress: { ...shipmentInfo.pickupAddress, city: 'San Antonio', state: 'TX', postalCode: '78234' },
        customerRemarks: mockMtoShipmentUB.customerRemarks,
        requestedPickupDate: mockMtoShipmentUB.requestedPickupDate,
        requestedDeliveryDate: mockMtoShipmentUB.requestedDeliveryDate,
        destinationAddress: { ...mockMtoShipmentUB.destinationAddress, streetAddress2: '' },
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
        ...mockMtoShipmentUB,
        pickupAddress: { ...shipmentInfo.pickupAddress },
        shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
        eTag: window.btoa(newUpdatedAt),
        status: 'SUBMITTED',
      };

      patchMTOShipment.mockImplementation(() => Promise.resolve(expectedUpdateResponse));
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentUB });

      const pickupAddress1Input = screen.getAllByLabelText(/Address 1 */)[0];
      await userEvent.clear(pickupAddress1Input);
      await userEvent.type(pickupAddress1Input, shipmentInfo.pickupAddress.streetAddress1);

      const pickupAddress2Input = screen.getAllByLabelText(/Address 2/)[0];
      await userEvent.clear(pickupAddress2Input);
      await userEvent.type(pickupAddress2Input, shipmentInfo.pickupAddress.streetAddress2);

      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalledWith(mockMtoShipmentUB.id, expectedPayload, mockMtoShipmentUB.eTag);
      });

      expect(ubProps.updateMTOShipment).toHaveBeenCalledWith(expectedUpdateResponse);

      expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
    });

    it('shows an error when there is an error with the submission', async () => {
      const shipmentInfo = {
        pickupAddress: {
          streetAddress1: '6622 Airport Way S',
          streetAddress2: '#1430',
          city: 'San Marcos',
          state: 'TX',
          postalCode: '78666',
        },
      };

      const errorMessage = 'Something broke!';
      const errorResponse = { response: { body: { detail: errorMessage } } };

      patchMTOShipment.mockImplementation(() => Promise.reject(errorResponse));
      getResponseError.mockImplementation((response, defaultMessage) => {
        return response.body?.detail || defaultMessage;
      });
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentUB });

      const pickupAddress1Input = screen.getAllByLabelText(/Address 1 */)[0];
      await userEvent.clear(pickupAddress1Input);
      await userEvent.type(pickupAddress1Input, shipmentInfo.pickupAddress.streetAddress1);

      const pickupAddress2Input = screen.getAllByLabelText(/Address 2/)[0];
      await userEvent.clear(pickupAddress2Input);
      await userEvent.type(pickupAddress2Input, shipmentInfo.pickupAddress.streetAddress2);

      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalled();
      });

      expect(getResponseError).toHaveBeenCalledWith(
        errorResponse.response,
        'failed to update MTO shipment due to server error',
      );

      expect(await screen.findByText(errorMessage)).toBeInTheDocument();
    });

    it('renders the UB shipment form with pre-filled values', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentUB });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText('Use my current address')).not.toBeChecked();
      expect(screen.getAllByLabelText(/Address 1 */)[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByTestId('City')[0]).toHaveTextContent('San Antonio');
      expect(screen.getAllByTestId(/State/)[0]).toHaveTextContent('TX');
      expect(screen.getAllByTestId(/ZIP/)[0]).toHaveTextContent('78234');
      expect(screen.getByText('San Antonio, TX 78234 ()'));
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText(/Address 1 */)[1]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByTestId('City')[1]).toHaveTextContent('Auburn');
      expect(screen.getAllByTestId(/State/)[1]).toHaveTextContent('WA');
      expect(screen.getAllByTestId(/ZIP/)[1]).toHaveTextContent('98001');
      expect(screen.getByText('Auburn, WA 98001 ()'));

      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');
    });

    it('renders the UB shipment with date validaton alerts for weekend and holiday', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United of States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentUB });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a holiday and weekend in the United of States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Preferred delivery date 11 Aug 2021 is on a holiday and weekend in the United of States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('renders the UB shipment with date validaton alerts for weekend', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentUB });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Preferred delivery date 11 Aug 2021 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('renders the UB shipment with date validaton alerts for holiday', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentUB });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.getByText(
            /Preferred delivery date 11 Aug 2021 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });

    it('renders the UB shipment with no date validaton alerts for pickup/delivery', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipmentUB });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');

      await waitFor(() => {
        expect(
          screen.queryAllByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveLength(0);
        expect(
          screen.queryAllByText(
            'Preferred delivery date 11 Aug 2021 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
          ),
        ).toHaveLength(0);
      });
    });
  });

  describe('creating a new NTS shipment', () => {
    it('renders the NTS shipment form', async () => {
      renderMtoShipmentForm({ shipmentType: SHIPMENT_OPTIONS.NTS });

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');

      expect(screen.getByText(/5,000 lbs/)).toHaveClass('usa-alert__text');

      expect(screen.getByTestId('pickupDateHint')).toHaveTextContent(
        'This is the day movers would put this shipment on their truck. Packing starts earlier. Dates will be finalized when you talk to your Customer Care Representative. Your requested pickup/load date should be your latest preferred pickup/load date, or the date you need to be out of your origin residence.',
      );
      expect(screen.getByText('Date')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText(/Preferred pickup date/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Pickup Address')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 1 */)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('City')).toBeInstanceOf(HTMLLabelElement);
      expect(screen.getByTestId('State')).toBeInstanceOf(HTMLLabelElement);
      expect(screen.getByTestId('ZIP')).toBeInstanceOf(HTMLLabelElement);
      expect(screen.getByLabelText('Location Lookup *', { exact: false })).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Releasing agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText(/First name/)).toHaveAttribute('name', 'pickup.agent.firstName');
      expect(screen.getByLabelText(/Last name/)).toHaveAttribute('name', 'pickup.agent.lastName');
      expect(screen.getByLabelText(/Phone/)).toHaveAttribute('name', 'pickup.agent.phone');
      expect(screen.getByLabelText(/Email/)).toHaveAttribute('name', 'pickup.agent.email');

      expect(screen.getAllByText('Date')).toHaveLength(1);
      expect(screen.getAllByText('Pickup Address')).toHaveLength(1);
      expect(screen.queryByText(/Receiving agent/)).not.toBeInTheDocument();
      expect(
        screen.queryByText(
          'Details about the facility where your things are now, including the name or address (if you know them)',
        ),
      ).not.toBeInTheDocument();

      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('renders the correct weight allowance when there are dependents', async () => {
      renderMtoShipmentForm({
        shipmentType: SHIPMENT_OPTIONS.NTS,
        orders: { has_dependents: true, authorizedWeight: 8000 },
      });

      expect(await screen.findByText('NTS')).toHaveClass('usa-tag');

      expect(screen.getByText(/8,000 lbs/)).toHaveClass('usa-alert__text');
    });

    it('renders special NTS What to expect section', async () => {
      const { queryByTestId } = renderMtoShipmentForm({ shipmentType: SHIPMENT_OPTIONS.NTS });

      await waitFor(() => {
        expect(queryByTestId('nts-what-to-expect')).toBeInTheDocument();
      });
    });

    it('renders NTS with preferred pickup date alert for holiday and weekend', async () => {
      const mockMtoShipment = {
        id: uuidv4(),
        eTag: window.btoa(updatedAt),
        createdAt: '2021-06-11T18:12:11.918Z',
        updatedAt,
        moveTaskOrderId: moveId,
        customerRemarks: 'mock remarks',
        requestedPickupDate: tomorrow,
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
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, shipmentType: SHIPMENT_OPTIONS.NTS, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
      await waitFor(() => {
        // only pickup date is available. delivery alert will never be present.
        expect(
          screen.getByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveClass('usa-alert__text');
        expect(
          screen.queryAllByText(
            'Preferred delivery date 11 Aug 2021 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
          ),
        ).toHaveLength(0);
      });
    });
  });

  describe('creating a new NTS-release shipment', () => {
    it('renders the NTS-release shipment form', async () => {
      renderMtoShipmentForm({ shipmentType: SHIPMENT_OPTIONS.NTSR });

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');

      expect(screen.getByText(/5,000 lbs/)).toHaveClass('usa-alert__text');

      expect(screen.queryByLabelText(/Preferred pickup date/)).not.toBeInTheDocument();
      expect(screen.queryByText('Pickup Info')).not.toBeInTheDocument();
      expect(screen.queryByText(/Releasing agent/)).not.toBeInTheDocument();

      expect(screen.getAllByText('Date')).toHaveLength(1);
      expect(screen.getAllByText('Delivery Address')).toHaveLength(1);

      expect(screen.getByText('Date')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText(/Preferred delivery date/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Delivery Address')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText(/Receiving agent/).parentElement).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText(/First name/)).toHaveAttribute('name', 'delivery.agent.firstName');
      expect(screen.getByLabelText(/Last name/)).toHaveAttribute('name', 'delivery.agent.lastName');
      expect(screen.getByLabelText(/Phone/)).toHaveAttribute('name', 'delivery.agent.phone');
      expect(screen.getByLabelText(/Email/)).toHaveAttribute('name', 'delivery.agent.email');

      expect(
        screen.queryByText(
          'Details about the facility where your things are now, including the name or address (if you know them)',
        ),
      ).toBeInTheDocument();

      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toBeInstanceOf(HTMLTextAreaElement);
    });

    it('renders the correct weight allowance when there are dependents', async () => {
      renderMtoShipmentForm({
        shipmentType: SHIPMENT_OPTIONS.NTSR,
        orders: { has_dependents: true, authorizedWeight: 8000 },
      });

      expect(await screen.findByText('NTS-release')).toHaveClass('usa-tag');

      expect(screen.getByText(/8,000 lbs/)).toHaveClass('usa-alert__text');
    });

    it('does not render special NTS What to expect section', async () => {
      const { queryByTestId } = renderMtoShipmentForm({ shipmentType: SHIPMENT_OPTIONS.NTSR });

      await waitFor(() => {
        expect(queryByTestId('nts-what-to-expect')).not.toBeInTheDocument();
      });
    });

    it('renders NTSR with preferred delivery date alert for holiday and weekend', async () => {
      const mockMtoShipment = {
        id: uuidv4(),
        eTag: window.btoa(updatedAt),
        createdAt: '2021-06-11T18:12:11.918Z',
        updatedAt,
        moveTaskOrderId: moveId,
        customerRemarks: 'mock remarks',
        requestedPickupDate: tomorrow,
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
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderMtoShipmentForm({ isCreatePage: false, shipmentType: SHIPMENT_OPTIONS.NTSR, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        // only delivery date is available. pickup alert will never be present.
        expect(
          screen.queryAllByText(
            `Preferred pickup date ${tomorrowDatePicker} is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.`,
          ),
        ).toHaveLength(0);
        expect(
          screen.getByText(
            /Preferred delivery date 11 Aug 2021 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
          ),
        ).toHaveClass('usa-alert__text');
      });
    });
  });

  describe('requestedPickupDate validation when creating and editing non-PPM shipments', () => {
    const mockHHGShipment = {
      ...mockMtoShipmentHHGWithDest,
      requestedPickupDate: '2021-06-07',
    };
    const mockNtsShipment = {
      ...mockHHGShipment,
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      storageFacility: {
        facilityName: 'Storage Facility',
        address: {
          city: 'Anytown',
          country: 'USA',
          postalCode: '90210',
          state: 'OK',
          streetAddress1: '555 Main Ave',
          streetAddress2: 'Apartment 900',
        },
      },
      tacType: 'HHG',
      sacType: 'NTS',
      tac: '123',
      sac: '456',
      serviceOrderNumber: '12341234',
    };

    const mockBoatShipment = (boatShipmentType) => ({
      ...mockHHGShipment,
      boatShipment: {
        type: boatShipmentType,
        year: 2020,
        make: 'Yamaha',
        model: '242X E-Series',
        lengthInInches: 276,
        widthInInches: 102,
        heightInInches: 120,
        hasTrailer: true,
        isRoadworthy: true,
      },
    });

    const mockMobileHomeShipment = {
      ...mockHHGShipment,
      mobileHomeShipment: {
        year: 2020,
        make: 'Yamaha',
        model: '242X E-Series',
        lengthInInches: 276,
        widthInInches: 102,
        heightInInches: 120,
      },
    };

    const mockedUB = {
      ...mockMtoShipmentUB,
      requestedPickupDate: '2021-06-07',
    };

    const shipmentTypesSource = [
      [SHIPMENT_TYPES.HHG, mockHHGShipment],
      [SHIPMENT_TYPES.NTS, mockNtsShipment],
      [SHIPMENT_TYPES.BOAT_HAUL_AWAY, mockBoatShipment(boatShipmentTypes.HAUL_AWAY)],
      [SHIPMENT_TYPES.BOAT_TOW_AWAY, mockBoatShipment(boatShipmentTypes.TOW_AWAY)],
      [SHIPMENT_TYPES.MOBILE_HOME, mockMobileHomeShipment],
      [SHIPMENT_TYPES.UNACCOMPANIED_BAGGAGE, mockedUB],
    ];

    const shipmentTypesToTest = [
      ...shipmentTypesSource.map((v) => [true].concat(v)),
      ...shipmentTypesSource.map((v) => [false].concat(v)),
    ];

    it.each(shipmentTypesToTest)(
      'preferredPickupDate (isCreate: %s | %s) - validation errors show',
      async (isCreate, shipmentType, mockShipment) => {
        const expectedDateSelectionIsWeekendHolidayResponse = {
          country_code: 'US',
          country_name: 'United States',
          is_weekend: true,
          is_holiday: true,
        };
        dateSelectionIsWeekendHoliday.mockImplementation(() =>
          Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
        );

        renderMtoShipmentForm({
          shipmentType,
          isCreatePage: isCreate,
          mtoShipment: isCreate ? undefined : mockShipment,
        });

        if (isCreate) {
          // Error doesn't show unless touched
          expect(
            within(await screen.findByTestId('preferredPickupDateFieldSet')).queryByTestId('errorMessage'),
          ).not.toBeInTheDocument();

          // Trigger error with empty date, field touched
          await act(async () => {
            const node = screen.getByLabelText(/Preferred pickup date/);
            await userEvent.clear(node);
            node.blur();
          });
          const dateRequiredParent = within(await screen.findByTestId('preferredPickupDateFieldSet')).queryByTestId(
            'formGroup',
          );
          await waitFor(() => {
            expect(within(dateRequiredParent).queryByTestId('errorMessage')).toHaveTextContent('Required');
          });
        } else {
          // Trigger error with invalid date, field touched
          await act(async () => {
            const node = screen.getByLabelText(/Preferred pickup date/);
            await userEvent.clear(node);
            await userEvent.paste('22 Mar 2022');
            node.blur();
          });
          const dateRequiredParent = within(await screen.findByTestId('preferredPickupDateFieldSet')).queryByTestId(
            'formGroup',
          );
          await waitFor(() => {
            expect(within(dateRequiredParent).queryByTestId('errorMessage')).toHaveTextContent(
              'Preferred pickup date must be in the future.',
            );
          });
        }

        // Trigger invalid date error - cannot be today
        const now = formatDateForDatePicker(formatDateWithUTC(new Date()));
        await act(async () => {
          const node = screen.getByLabelText(/Preferred pickup date/);
          await userEvent.clear(node);
          await userEvent.paste(now);
          node.blur();
        });
        expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(now);
        const dateRequiredParent = within(await screen.findByTestId('preferredPickupDateFieldSet')).queryByTestId(
          'formGroup',
        );
        await waitFor(() => {
          expect(within(dateRequiredParent).queryByTestId('errorMessage')).toHaveTextContent(
            'Preferred pickup date must be in the future.',
          );
        });
      },
    );

    it.each(shipmentTypesToTest)(
      'preferredPickupDate (isCreate: %s | %s) - validation errors hide when valid and holiday alert shows',
      async (isCreate, shipmentType, mockShipment) => {
        const expectedDateSelectionIsWeekendHolidayResponse = {
          country_code: 'US',
          country_name: 'United States',
          is_weekend: true,
          is_holiday: true,
        };
        dateSelectionIsWeekendHoliday.mockImplementation(() =>
          Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
        );
        renderMtoShipmentForm({
          shipmentType,
          isCreatePage: isCreate,
          mtoShipment: isCreate ? {} : mockShipment,
        });

        // Trigger invalid date error - must be in the future
        await act(async () => {
          const node = screen.getByLabelText(/Preferred pickup date/);
          await userEvent.clear(node);
          await userEvent.paste('26 Mar 2022');
          node.blur();
        });
        expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('26 Mar 2022');
        const dateRequiredParent = within(await screen.findByTestId('preferredPickupDateFieldSet')).queryByTestId(
          'formGroup',
        );
        await waitFor(() => {
          expect(within(dateRequiredParent).queryByTestId('errorMessage')).toHaveTextContent(
            'Preferred pickup date must be in the future.',
          );
        });
        // should hide holiday alert
        expect(screen.queryByTestId('preferredPickupDateAlert')).not.toBeInTheDocument();

        // Valid date, hides errors
        await act(async () => {
          const node = screen.getByLabelText(/Preferred pickup date/);
          await userEvent.clear(node);
          await userEvent.paste(tomorrowDatePicker);
          node.blur();
        });

        expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue(tomorrowDatePicker);
        await waitFor(() => {
          expect(
            within(screen.getByTestId('preferredPickupDateFieldSet')).queryByTestId('errorMessage'),
          ).not.toBeInTheDocument();
          expect(screen.getByTestId('preferredPickupDateAlert')).toBeVisible();
        });
      },
    );
  });
});

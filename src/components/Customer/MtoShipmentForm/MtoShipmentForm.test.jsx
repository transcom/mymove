/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { generatePath } from 'react-router-dom';
import { waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { v4 as uuidv4 } from 'uuid';

import MtoShipmentForm from './MtoShipmentForm';

import { customerRoutes } from 'constants/routes';
import {
  createMTOShipment,
  getResponseError,
  patchMTOShipment,
  dateSelectionIsWeekendHoliday,
} from 'services/internalApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { renderWithRouter } from 'testUtils';
import { ORDERS_TYPE } from 'constants/orders';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

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

const renderMtoShipmentForm = (props) => {
  return renderWithRouter(<MtoShipmentForm {...defaultProps} {...props} />, {
    path: customerRoutes.SHIPMENT_CREATE_PATH,
    params: { moveId },
  });
};

const renderUBShipmentForm = (props) => {
  return renderWithRouter(<MtoShipmentForm {...ubProps} {...props} />, {
    path: customerRoutes.SHIPMENT_CREATE_PATH,
    params: { moveId },
  });
};

describe('MtoShipmentForm component', () => {
  describe('when creating a new HHG shipment', () => {
    it('renders the HHG shipment form', async () => {
      renderMtoShipmentForm();

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
      expect(screen.getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/City/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/State/)).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByLabelText(/ZIP/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 4, name: 'Second Pickup Address' })).toBeInTheDocument();
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
      const { queryByLabelText, queryAllByLabelText } = renderMtoShipmentForm();

      await userEvent.click(queryByLabelText('Use my current address'));

      await waitFor(() => {
        expect(queryAllByLabelText(/Address 1/)[0]).toHaveValue(defaultProps.currentResidence.streetAddress1);
        expect(queryAllByLabelText(/Address 2/)[0]).toHaveValue('');
        expect(queryAllByLabelText(/City/)[0]).toHaveValue(defaultProps.currentResidence.city);
        expect(queryAllByLabelText(/State/)[0]).toHaveValue(defaultProps.currentResidence.state);
        expect(queryAllByLabelText(/ZIP/)[0]).toHaveValue(defaultProps.currentResidence.postalCode);
      });
    });

    it('renders a second address fieldset when the user has a second pickup address', async () => {
      renderMtoShipmentForm();

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1/);
      expect(streetAddress1[1]).toHaveAttribute('name', 'secondaryPickup.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[1]).toHaveAttribute('name', 'secondaryPickup.address.streetAddress2');

      const city = await screen.findAllByLabelText(/City/);
      expect(city[1]).toHaveAttribute('name', 'secondaryPickup.address.city');

      const state = await screen.findAllByLabelText(/State/);
      expect(state[1]).toHaveAttribute('name', 'secondaryPickup.address.state');

      const zip = await screen.findAllByLabelText(/ZIP/);
      expect(zip[1]).toHaveAttribute('name', 'secondaryPickup.address.postalCode');
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      renderMtoShipmentForm();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1/);
      expect(streetAddress1[0]).toHaveAttribute('name', 'pickup.address.streetAddress1');
      expect(streetAddress1[1]).toHaveAttribute('name', 'delivery.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[0]).toHaveAttribute('name', 'pickup.address.streetAddress2');
      expect(streetAddress2[1]).toHaveAttribute('name', 'delivery.address.streetAddress2');

      const city = await screen.findAllByLabelText(/City/);
      expect(city[0]).toHaveAttribute('name', 'pickup.address.city');
      expect(city[1]).toHaveAttribute('name', 'delivery.address.city');

      const state = await screen.findAllByLabelText(/State/);
      expect(state[0]).toHaveAttribute('name', 'pickup.address.state');
      expect(state[1]).toHaveAttribute('name', 'delivery.address.state');

      const zip = await screen.findAllByLabelText(/ZIP/);
      expect(zip[0]).toHaveAttribute('name', 'pickup.address.postalCode');
      expect(zip[1]).toHaveAttribute('name', 'delivery.address.postalCode');
    });

    it('renders the secondary delivery address question once a user says they have a primary delivery address', async () => {
      renderMtoShipmentForm();

      expect(screen.queryByRole('heading', { level: 4, name: 'Second Delivery Address' })).not.toBeInTheDocument();
      expect(screen.queryByTitle('Yes, I have a second delivery address')).not.toBeInTheDocument();
      expect(screen.queryByTitle('No, I do not have a second delivery address')).not.toBeInTheDocument();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));

      expect(await screen.findByRole('heading', { level: 4, name: 'Second Delivery Address' })).toBeInTheDocument();
      expect(screen.getByTitle('Yes, I have a second delivery address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not have a second delivery address')).toBeInstanceOf(HTMLInputElement);
    });

    it('renders another address fieldset when the user has a second delivery address', async () => {
      renderMtoShipmentForm();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1/);
      expect(streetAddress1.length).toBe(3);
      expect(streetAddress1[2]).toHaveAttribute('name', 'secondaryDelivery.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2.length).toBe(3);
      expect(streetAddress2[2]).toHaveAttribute('name', 'secondaryDelivery.address.streetAddress2');

      const city = await screen.findAllByLabelText(/City/);
      expect(city.length).toBe(3);
      expect(city[2]).toHaveAttribute('name', 'secondaryDelivery.address.city');

      const state = await screen.findAllByLabelText(/State/);
      expect(state.length).toBe(3);
      expect(state[2]).toHaveAttribute('name', 'secondaryDelivery.address.state');

      const zip = await screen.findAllByLabelText(/ZIP/);
      expect(zip.length).toBe(3);
      expect(zip[2]).toHaveAttribute('name', 'secondaryDelivery.address.postalCode');
    });

    it('goes back when the back button is clicked', async () => {
      renderMtoShipmentForm();

      const backButton = await screen.findByRole('button', { name: 'Back' });
      await userEvent.click(backButton);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith(-1);
      });
    });

    it('can submit a new HHG shipment successfully', async () => {
      const shipmentInfo = {
        requestedPickupDate: '07 Jun 2021',
        pickupAddress: {
          streetAddress1: '812 S 129th St',
          streetAddress2: '#123',
          city: 'San Antonio',
          state: 'TX',
          postalCode: '78234',
        },
        requestedDeliveryDate: '14 Jun 2021',
      };

      const expectedPayload = {
        agents: [
          { agentType: 'RELEASING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
          { agentType: 'RECEIVING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
        ],
        moveTaskOrderID: moveId,
        shipmentType: SHIPMENT_OPTIONS.HHG,
        customerRemarks: '',
        requestedPickupDate: '2021-06-07',
        pickupAddress: { ...shipmentInfo.pickupAddress },
        requestedDeliveryDate: '2021-06-14',
        hasSecondaryPickupAddress: false,
        hasSecondaryDeliveryAddress: false,
        hasTertiaryPickupAddress: false,
        hasTertiaryDeliveryAddress: false,
      };

      const updatedAt = '2021-06-11T18:12:11.918Z';
      const expectedCreateResponse = {
        createdAt: '2021-06-11T18:12:11.918Z',
        customerRemarks: '',
        eTag: window.btoa(updatedAt),
        id: uuidv4(),
        moveTaskOrderID: moveId,
        pickupAddress: { ...shipmentInfo.pickupAddress, id: uuidv4() },
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
      renderMtoShipmentForm();

      const pickupDateInput = await screen.findByLabelText(/Preferred pickup date/);
      await userEvent.type(pickupDateInput, shipmentInfo.requestedPickupDate);

      const pickupAddress1Input = screen.getByLabelText(/Address 1/);
      await userEvent.type(pickupAddress1Input, shipmentInfo.pickupAddress.streetAddress1);

      const pickupAddress2Input = screen.getByLabelText(/Address 2/);
      await userEvent.type(pickupAddress2Input, shipmentInfo.pickupAddress.streetAddress2);

      const pickupCityInput = screen.getByLabelText(/City/);
      await userEvent.type(pickupCityInput, shipmentInfo.pickupAddress.city);

      const pickupStateInput = screen.getByLabelText(/State/);
      await userEvent.selectOptions(pickupStateInput, shipmentInfo.pickupAddress.state);

      const pickupPostalCodeInput = screen.getByLabelText(/ZIP/);
      await userEvent.type(pickupPostalCodeInput, shipmentInfo.pickupAddress.postalCode);

      const deliveryDateInput = await screen.findByLabelText(/Preferred delivery date/);
      await userEvent.type(deliveryDateInput, shipmentInfo.requestedDeliveryDate);

      const nextButton = await screen.findByRole('button', { name: 'Next' });
      expect(nextButton).not.toBeDisabled();
      await userEvent.click(nextButton);

      await waitFor(() => {
        expect(createMTOShipment).toHaveBeenCalledWith(expectedPayload);
      });

      expect(defaultProps.updateMTOShipment).toHaveBeenCalledWith(expectedCreateResponse);

      expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
    });

    it('shows an error when there is an error with the submission', async () => {
      const shipmentInfo = {
        requestedPickupDate: '07 Jun 2021',
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
      renderMtoShipmentForm();

      const pickupDateInput = await screen.findByLabelText(/Preferred pickup date/);
      await userEvent.type(pickupDateInput, shipmentInfo.requestedPickupDate);

      const pickupAddress1Input = screen.getByLabelText(/Address 1/);
      await userEvent.type(pickupAddress1Input, shipmentInfo.pickupAddress.streetAddress1);

      const pickupAddress2Input = screen.getByLabelText(/Address 2/);
      await userEvent.type(pickupAddress2Input, shipmentInfo.pickupAddress.streetAddress2);

      const pickupCityInput = screen.getByLabelText(/City/);
      await userEvent.type(pickupCityInput, shipmentInfo.pickupAddress.city);

      const pickupStateInput = screen.getByLabelText(/State/);
      await userEvent.selectOptions(pickupStateInput, shipmentInfo.pickupAddress.state);

      const pickupPostalCodeInput = screen.getByLabelText(/ZIP/);
      await userEvent.type(pickupPostalCodeInput, shipmentInfo.pickupAddress.postalCode);

      const deliveryDateInput = await screen.findByLabelText(/Preferred delivery date/);
      await userEvent.type(deliveryDateInput, shipmentInfo.requestedDeliveryDate);

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
    const updatedAt = '2021-06-11T18:12:11.918Z';

    const mockMtoShipment = {
      id: uuidv4(),
      eTag: window.btoa(updatedAt),
      createdAt: '2021-06-11T18:12:11.918Z',
      updatedAt,
      moveTaskOrderId: moveId,
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
      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');

      expect(
        screen.getByText(
          /Preferred pickup date 01 Aug 2021 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
        ),
      ).toHaveClass('usa-alert__text');
      expect(
        screen.getAllByText(
          'Preferred pickup date 01 Aug 2021 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
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
          city: mockMtoShipment.destinationAddress.city,
          state: mockMtoShipment.destinationAddress.state,
          postalCode: mockMtoShipment.destinationAddress.postalCode,
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

      const streetAddress1 = await screen.findAllByLabelText(/Address 1/);
      expect(streetAddress1.length).toBe(4);

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2.length).toBe(4);

      const city = await screen.findAllByLabelText(/City/);
      expect(city.length).toBe(4);

      const state = await screen.findAllByLabelText(/State/);
      expect(state.length).toBe(4);

      const zip = await screen.findAllByLabelText(/ZIP/);
      expect(zip.length).toBe(4);

      // Secondary pickup address should be the 2nd address
      expect(streetAddress1[1]).toHaveValue('142 E Barrel Hoop Circle');
      expect(streetAddress2[1]).toHaveValue('#4A');
      expect(city[1]).toHaveValue('Corpus Christi');
      expect(state[1]).toHaveValue('TX');
      expect(zip[1]).toHaveValue('78412');

      // Secondary delivery address should be the 4th address
      expect(streetAddress1[3]).toHaveValue('3373 NW Martin Luther King Jr Blvd');
      expect(streetAddress2[3]).toHaveValue('');
      expect(city[3]).toHaveValue(mockMtoShipment.destinationAddress.city);
      expect(state[3]).toHaveValue(mockMtoShipment.destinationAddress.state);
      expect(zip[3]).toHaveValue(mockMtoShipment.destinationAddress.postalCode);
    });

    it.each([
      [/Address 1/, 'Some Address'],
      [/Address 2/, '123'],
      [/City/, 'Some City'],
      [/ZIP/, '92131'],
    ])(
      'does not allow the user to save the form if the %s field on a secondary addreess is the only one filled out',
      async (fieldName, text) => {
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

        // Verify that the form is good to submit by checking that the save button is not disabled.
        const saveButton = await screen.findByRole('button', { name: 'Save' });
        expect(saveButton).not.toBeDisabled();

        await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));
        await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

        const address = await screen.findAllByLabelText(fieldName);
        // The second instance of a field is the secondary pickup
        await userEvent.type(address[1], text);
        await waitFor(() => {
          expect(saveButton).toBeDisabled();
        });

        // Clear the field so that the secondary delivery address can be checked
        await userEvent.clear(address[1]);
        await waitFor(() => {
          expect(saveButton).not.toBeDisabled();
        });

        // The fourth instance found is the secondary delivery
        await userEvent.type(address[3], text);
        await waitFor(() => {
          expect(saveButton).toBeDisabled();
        });

        await userEvent.clear(address[3]);
        await waitFor(() => {
          expect(saveButton).not.toBeDisabled();
        });
      },
    );

    // Similar test as above, but with the state input.
    // Extracted out since the state field is not a text input.
    it('does not allow the user to save the form if the state field on a secondary addreess is the only one filled out', async () => {
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

      // Verify that the form is good to submit by checking that the save button is not disabled.
      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      const state = await screen.findAllByLabelText(/State/);
      // The second instance of a field is the secondary pickup
      await userEvent.selectOptions(state[1], 'CA');
      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });

      // Change the selection to blank so that the secondary delivery address can be checked
      await userEvent.selectOptions(state[1], '');
      await waitFor(() => {
        expect(saveButton).not.toBeDisabled();
      });

      // The fourth instance found is the secondary delivery
      await userEvent.selectOptions(state[3], 'CA');
      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });

      await userEvent.selectOptions(state[3], '');
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
        pickupAddress: { ...shipmentInfo.pickupAddress },
        customerRemarks: mockMtoShipment.customerRemarks,
        requestedPickupDate: mockMtoShipment.requestedPickupDate,
        requestedDeliveryDate: mockMtoShipment.requestedDeliveryDate,
        destinationAddress: { ...mockMtoShipment.destinationAddress, streetAddress2: '' },
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

      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalledWith(mockMtoShipment.id, expectedPayload, mockMtoShipment.eTag);
      });

      expect(defaultProps.updateMTOShipment).toHaveBeenCalledWith(expectedUpdateResponse);

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
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Preferred pickup date 01 Aug 2021 is on a holiday and weekend in the United of States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
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
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Preferred pickup date 01 Aug 2021 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
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
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Preferred pickup date 01 Aug 2021 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
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
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');

      await waitFor(() => {
        expect(
          screen.queryAllByText(
            'Preferred pickup date 01 Aug 2021 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
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

      expect(screen.getAllByText('Date')[0]).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText(/Preferred pickup date/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Pickup info' })).toBeInTheDocument();
      expect(screen.getByTestId('pickupDateHint')).toHaveTextContent(
        'This is the day movers would put this shipment on their truck. Packing starts earlier. Dates will be finalized when you talk to your Customer Care Representative. Your requested pickup/load date should be your latest preferred pickup/load date, or the date you need to be out of your origin residence.',
      );
      expect(screen.getByText('Pickup Address')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Use my current address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/City/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/State/)).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByLabelText(/ZIP/)).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 4, name: 'Second Pickup Address' })).toBeInTheDocument();
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
      const { queryByLabelText, queryAllByLabelText } = renderUBShipmentForm();

      await userEvent.click(queryByLabelText('Use my current address'));

      await waitFor(() => {
        expect(queryAllByLabelText(/Address 1/)[0]).toHaveValue(defaultProps.currentResidence.streetAddress1);
        expect(queryAllByLabelText(/Address 2/)[0]).toHaveValue('');
        expect(queryAllByLabelText(/City/)[0]).toHaveValue(defaultProps.currentResidence.city);
        expect(queryAllByLabelText(/State/)[0]).toHaveValue(defaultProps.currentResidence.state);
        expect(queryAllByLabelText(/ZIP/)[0]).toHaveValue(defaultProps.currentResidence.postalCode);
      });
    });

    it('renders a second address fieldset when the user has a second pickup address', async () => {
      renderUBShipmentForm();

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1/);
      expect(streetAddress1[1]).toHaveAttribute('name', 'secondaryPickup.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[1]).toHaveAttribute('name', 'secondaryPickup.address.streetAddress2');

      const city = await screen.findAllByLabelText(/City/);
      expect(city[1]).toHaveAttribute('name', 'secondaryPickup.address.city');

      const state = await screen.findAllByLabelText(/State/);
      expect(state[1]).toHaveAttribute('name', 'secondaryPickup.address.state');

      const zip = await screen.findAllByLabelText(/ZIP/);
      expect(zip[1]).toHaveAttribute('name', 'secondaryPickup.address.postalCode');
    });

    it('renders a second address fieldset when the user has a delivery address', async () => {
      renderUBShipmentForm();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1/);
      expect(streetAddress1[0]).toHaveAttribute('name', 'pickup.address.streetAddress1');
      expect(streetAddress1[1]).toHaveAttribute('name', 'delivery.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2[0]).toHaveAttribute('name', 'pickup.address.streetAddress2');
      expect(streetAddress2[1]).toHaveAttribute('name', 'delivery.address.streetAddress2');

      const city = await screen.findAllByLabelText(/City/);
      expect(city[0]).toHaveAttribute('name', 'pickup.address.city');
      expect(city[1]).toHaveAttribute('name', 'delivery.address.city');

      const state = await screen.findAllByLabelText(/State/);
      expect(state[0]).toHaveAttribute('name', 'pickup.address.state');
      expect(state[1]).toHaveAttribute('name', 'delivery.address.state');

      const zip = await screen.findAllByLabelText(/ZIP/);
      expect(zip[0]).toHaveAttribute('name', 'pickup.address.postalCode');
      expect(zip[1]).toHaveAttribute('name', 'delivery.address.postalCode');
    });

    it('renders the secondary delivery address question once a user says they have a primary delivery address', async () => {
      renderUBShipmentForm();

      expect(screen.queryByRole('heading', { level: 4, name: 'Second Delivery Address' })).not.toBeInTheDocument();
      expect(screen.queryByTitle('Yes, I have a second delivery address')).not.toBeInTheDocument();
      expect(screen.queryByTitle('No, I do not have a second delivery address')).not.toBeInTheDocument();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));

      expect(await screen.findByRole('heading', { level: 4, name: 'Second Delivery Address' })).toBeInTheDocument();
      expect(screen.getByTitle('Yes, I have a second delivery address')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTitle('No, I do not have a second delivery address')).toBeInstanceOf(HTMLInputElement);
    });

    it('renders another address fieldset when the user has a second delivery address', async () => {
      renderUBShipmentForm();

      await userEvent.click(screen.getByTitle('Yes, I know my delivery address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      const streetAddress1 = await screen.findAllByLabelText(/Address 1/);
      expect(streetAddress1.length).toBe(3);
      expect(streetAddress1[2]).toHaveAttribute('name', 'secondaryDelivery.address.streetAddress1');

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2.length).toBe(3);
      expect(streetAddress2[2]).toHaveAttribute('name', 'secondaryDelivery.address.streetAddress2');

      const city = await screen.findAllByLabelText(/City/);
      expect(city.length).toBe(3);
      expect(city[2]).toHaveAttribute('name', 'secondaryDelivery.address.city');

      const state = await screen.findAllByLabelText(/State/);
      expect(state.length).toBe(3);
      expect(state[2]).toHaveAttribute('name', 'secondaryDelivery.address.state');

      const zip = await screen.findAllByLabelText(/ZIP/);
      expect(zip.length).toBe(3);
      expect(zip[2]).toHaveAttribute('name', 'secondaryDelivery.address.postalCode');
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
      const shipmentInfo = {
        requestedPickupDate: '07 Jun 2021',
        pickupAddress: {
          streetAddress1: '812 S 129th St',
          streetAddress2: '#123',
          city: 'San Antonio',
          state: 'TX',
          postalCode: '78234',
        },
        requestedDeliveryDate: '14 Jun 2021',
        shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
      };

      const expectedPayload = {
        agents: [
          { agentType: 'RELEASING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
          { agentType: 'RECEIVING_AGENT', email: '', firstName: '', lastName: '', phone: '' },
        ],
        moveTaskOrderID: moveId,
        shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
        customerRemarks: '',
        requestedPickupDate: '2021-06-07',
        pickupAddress: { ...shipmentInfo.pickupAddress },
        requestedDeliveryDate: '2021-06-14',
        hasSecondaryPickupAddress: false,
        hasSecondaryDeliveryAddress: false,
        hasTertiaryPickupAddress: false,
        hasTertiaryDeliveryAddress: false,
      };

      const updatedAt = '2021-06-11T18:12:11.918Z';
      const expectedCreateResponse = {
        createdAt: '2021-06-11T18:12:11.918Z',
        customerRemarks: '',
        eTag: window.btoa(updatedAt),
        id: uuidv4(),
        moveTaskOrderID: moveId,
        pickupAddress: { ...shipmentInfo.pickupAddress, id: uuidv4() },
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
      renderUBShipmentForm();

      const pickupDateInput = await screen.findByLabelText(/Preferred pickup date/);
      await userEvent.type(pickupDateInput, shipmentInfo.requestedPickupDate);

      const pickupAddress1Input = screen.getByLabelText(/Address 1/);
      await userEvent.type(pickupAddress1Input, shipmentInfo.pickupAddress.streetAddress1);

      const pickupAddress2Input = screen.getByLabelText(/Address 2/);
      await userEvent.type(pickupAddress2Input, shipmentInfo.pickupAddress.streetAddress2);

      const pickupCityInput = screen.getByLabelText(/City/);
      await userEvent.type(pickupCityInput, shipmentInfo.pickupAddress.city);

      const pickupStateInput = screen.getByLabelText(/State/);
      await userEvent.selectOptions(pickupStateInput, shipmentInfo.pickupAddress.state);

      const pickupPostalCodeInput = screen.getByLabelText(/ZIP/);
      await userEvent.type(pickupPostalCodeInput, shipmentInfo.pickupAddress.postalCode);

      const deliveryDateInput = await screen.findByLabelText(/Preferred delivery date/);
      await userEvent.type(deliveryDateInput, shipmentInfo.requestedDeliveryDate);

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
      const shipmentInfo = {
        requestedPickupDate: '07 Jun 2021',
        pickupAddress: {
          streetAddress1: '812 S 129th St',
          streetAddress2: '#123',
          city: 'San Antonio',
          state: 'TX',
          postalCode: '78234',
        },
        requestedDeliveryDate: '14 Jun 2021',
        shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
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
      renderUBShipmentForm();

      const pickupDateInput = await screen.findByLabelText(/Preferred pickup date/);
      await userEvent.type(pickupDateInput, shipmentInfo.requestedPickupDate);

      const pickupAddress1Input = screen.getByLabelText(/Address 1/);
      await userEvent.type(pickupAddress1Input, shipmentInfo.pickupAddress.streetAddress1);

      const pickupAddress2Input = screen.getByLabelText(/Address 2/);
      await userEvent.type(pickupAddress2Input, shipmentInfo.pickupAddress.streetAddress2);

      const pickupCityInput = screen.getByLabelText(/City/);
      await userEvent.type(pickupCityInput, shipmentInfo.pickupAddress.city);

      const pickupStateInput = screen.getByLabelText(/State/);
      await userEvent.selectOptions(pickupStateInput, shipmentInfo.pickupAddress.state);

      const pickupPostalCodeInput = screen.getByLabelText(/ZIP/);
      await userEvent.type(pickupPostalCodeInput, shipmentInfo.pickupAddress.postalCode);

      const deliveryDateInput = await screen.findByLabelText(/Preferred delivery date/);
      await userEvent.type(deliveryDateInput, shipmentInfo.requestedDeliveryDate);

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
    const updatedAt = '2021-06-11T18:12:11.918Z';

    const mockMtoShipment = {
      id: uuidv4(),
      eTag: window.btoa(updatedAt),
      createdAt: '2021-06-11T18:12:11.918Z',
      updatedAt,
      moveTaskOrderId: moveId,
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
      shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
    };

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
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

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
      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');

      expect(
        screen.getByText(
          /Preferred pickup date 01 Aug 2021 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
        ),
      ).toHaveClass('usa-alert__text');
      expect(
        screen.getAllByText(
          'Preferred pickup date 01 Aug 2021 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
        ),
      ).toHaveLength(1);
    });

    it('renders the UB shipment form with pre-filled secondary addresses', async () => {
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
        is_weekend: true,
        is_holiday: true,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: shipment });

      expect(await screen.findByTitle('Yes, I have a second pickup address')).toBeChecked();
      expect(await screen.findByTitle('Yes, I have a second delivery address')).toBeChecked();

      const streetAddress1 = await screen.findAllByLabelText(/Address 1/);
      expect(streetAddress1.length).toBe(4);

      const streetAddress2 = await screen.findAllByLabelText(/Address 2/);
      expect(streetAddress2.length).toBe(4);

      const city = await screen.findAllByLabelText(/City/);
      expect(city.length).toBe(4);

      const state = await screen.findAllByLabelText(/State/);
      expect(state.length).toBe(4);

      const zip = await screen.findAllByLabelText(/ZIP/);
      expect(zip.length).toBe(4);

      // Secondary pickup address should be the 2nd address
      expect(streetAddress1[1]).toHaveValue('142 E Barrel Hoop Circle');
      expect(streetAddress2[1]).toHaveValue('#4A');
      expect(city[1]).toHaveValue('Corpus Christi');
      expect(state[1]).toHaveValue('TX');
      expect(zip[1]).toHaveValue('78412');

      // Secondary delivery address should be the 4th address
      expect(streetAddress1[3]).toHaveValue('3373 NW Martin Luther King Jr Blvd');
      expect(streetAddress2[3]).toHaveValue('');
      expect(city[3]).toHaveValue(mockMtoShipment.destinationAddress.city);
      expect(state[3]).toHaveValue(mockMtoShipment.destinationAddress.state);
      expect(zip[3]).toHaveValue(mockMtoShipment.destinationAddress.postalCode);
    });

    it.each([
      [/Address 1/, 'Some Address'],
      [/Address 2/, '123'],
      [/City/, 'Some City'],
      [/ZIP/, '92131'],
    ])(
      'does not allow the user to save the form if the %s field on a secondary addreess is the only one filled out',
      async (fieldName, text) => {
        const expectedDateSelectionIsWeekendHolidayResponse = {
          country_code: 'US',
          country_name: 'United States',
          is_weekend: false,
          is_holiday: false,
        };
        dateSelectionIsWeekendHoliday.mockImplementation(() =>
          Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
        );
        renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

        // Verify that the form is good to submit by checking that the save button is not disabled.
        const saveButton = await screen.findByRole('button', { name: 'Save' });
        expect(saveButton).not.toBeDisabled();

        await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));
        await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

        const address = await screen.findAllByLabelText(fieldName);
        // The second instance of a field is the secondary pickup
        await userEvent.type(address[1], text);
        await waitFor(() => {
          expect(saveButton).toBeDisabled();
        });

        // Clear the field so that the secondary delivery address can be checked
        await userEvent.clear(address[1]);
        await waitFor(() => {
          expect(saveButton).not.toBeDisabled();
        });

        // The fourth instance found is the secondary delivery
        await userEvent.type(address[3], text);
        await waitFor(() => {
          expect(saveButton).toBeDisabled();
        });

        await userEvent.clear(address[3]);
        await waitFor(() => {
          expect(saveButton).not.toBeDisabled();
        });
      },
    );

    // Similar test as above, but with the state input.
    // Extracted out since the state field is not a text input.
    it('does not allow the user to save the form if the state field on a secondary addreess is the only one filled out', async () => {
      const expectedDateSelectionIsWeekendHolidayResponse = {
        country_code: 'US',
        country_name: 'United States',
        is_weekend: false,
        is_holiday: false,
      };
      dateSelectionIsWeekendHoliday.mockImplementation(() =>
        Promise.resolve({ data: JSON.stringify(expectedDateSelectionIsWeekendHolidayResponse) }),
      );
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

      // Verify that the form is good to submit by checking that the save button is not disabled.
      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      const state = await screen.findAllByLabelText(/State/);
      // The second instance of a field is the secondary pickup
      await userEvent.selectOptions(state[1], 'CA');
      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });

      // Change the selection to blank so that the secondary delivery address can be checked
      await userEvent.selectOptions(state[1], '');
      await waitFor(() => {
        expect(saveButton).not.toBeDisabled();
      });

      // The fourth instance found is the secondary delivery
      await userEvent.selectOptions(state[3], 'CA');
      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });

      await userEvent.selectOptions(state[3], '');
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
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

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
        pickupAddress: { ...shipmentInfo.pickupAddress },
        customerRemarks: mockMtoShipment.customerRemarks,
        requestedPickupDate: mockMtoShipment.requestedPickupDate,
        requestedDeliveryDate: mockMtoShipment.requestedDeliveryDate,
        destinationAddress: { ...mockMtoShipment.destinationAddress, streetAddress2: '' },
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
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

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

      const saveButton = await screen.findByRole('button', { name: 'Save' });
      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(patchMTOShipment).toHaveBeenCalledWith(mockMtoShipment.id, expectedPayload, mockMtoShipment.eTag);
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
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });

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
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
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
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Preferred pickup date 01 Aug 2021 is on a holiday and weekend in the United of States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
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
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Preferred pickup date 01 Aug 2021 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
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
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      await waitFor(() => {
        expect(
          screen.getByText(
            /Preferred pickup date 01 Aug 2021 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
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
      renderUBShipmentForm({ isCreatePage: false, mtoShipment: mockMtoShipment });
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      expect(screen.getByLabelText(/Preferred delivery date/)).toHaveValue('11 Aug 2021');
      expect(
        screen.getByLabelText(
          'Are there things about this shipment that your counselor or movers should discuss with you?',
        ),
      ).toHaveValue('mock remarks');

      await waitFor(() => {
        expect(
          screen.queryAllByText(
            'Preferred pickup date 01 Aug 2021 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
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
      expect(screen.getByLabelText(/Address 1/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/Address 2/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/City/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText(/State/)).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByLabelText(/ZIP/)).toBeInstanceOf(HTMLInputElement);

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
      const updatedAt = '2021-06-11T18:12:11.918Z';
      const mockMtoShipment = {
        id: uuidv4(),
        eTag: window.btoa(updatedAt),
        createdAt: '2021-06-11T18:12:11.918Z',
        updatedAt,
        moveTaskOrderId: moveId,
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
      expect(await screen.findByLabelText(/Preferred pickup date/)).toHaveValue('01 Aug 2021');
      await waitFor(() => {
        // only pickup date is available. delivery alert will never be present.
        expect(
          screen.getByText(
            /Preferred pickup date 01 Aug 2021 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date./,
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
      expect(screen.getAllByText('Second Delivery Address')).toHaveLength(1);

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
      const updatedAt = '2021-06-11T18:12:11.918Z';
      const mockMtoShipment = {
        id: uuidv4(),
        eTag: window.btoa(updatedAt),
        createdAt: '2021-06-11T18:12:11.918Z',
        updatedAt,
        moveTaskOrderId: moveId,
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
            'Preferred pickup date 01 Aug 2021 is on a holiday and weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
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
});

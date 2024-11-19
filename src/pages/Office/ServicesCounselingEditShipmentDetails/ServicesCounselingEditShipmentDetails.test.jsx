/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen, waitFor, within, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingEditShipmentDetails from './ServicesCounselingEditShipmentDetails';

import { updateMTOShipment, updateMoveCloseoutOffice, searchTransportationOffices } from 'services/ghcApi';
import { validatePostalCode } from 'utils/validation';
import { useEditShipmentQueries } from 'hooks/queries';
import { MOVE_STATUSES, SHIPMENT_OPTIONS } from 'shared/constants';
import { servicesCounselingRoutes } from 'constants/routes';
import { renderWithProviders } from 'testUtils';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const mockRoutingParams = { moveCode: 'move123', shipmentId: 'shipment123' };
const mockRoutingConfig = { path: servicesCounselingRoutes.BASE_SHIPMENT_EDIT_PATH, params: mockRoutingParams };
const mockTransportationOffice = [
  {
    address: {
      city: '',
      id: '00000000-0000-0000-0000-000000000000',
      postalCode: '',
      state: '',
      streetAddress1: '',
      county: '',
    },
    address_id: '46c4640b-c35e-4293-a2f1-36c7b629f903',
    affiliation: 'AIR_FORCE',
    created_at: '2021-02-11T16:48:04.117Z',
    id: '93f0755f-6f35-478b-9a75-35a69211da1c',
    name: 'Altus AFB',
    updated_at: '2021-02-11T16:48:04.117Z',
  },
];

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  updateMTOShipment: jest.fn(),
  updateMoveCloseoutOffice: jest.fn(),
  searchTransportationOffices: jest.fn().mockImplementation(() => Promise.resolve(mockTransportationOffice)),
}));

jest.mock('hooks/queries', () => ({
  ...jest.requireActual('@tanstack/react-query'),
  useEditShipmentQueries: jest.fn(),
}));

jest.mock('utils/validation', () => ({
  ...jest.requireActual('utils/validation'),
  validatePostalCode: jest.fn(),
}));

jest.mock('components/LocationSearchBox/api', () => ({
  ShowAddress: jest.fn().mockImplementation(() =>
    Promise.resolve({
      city: 'Glendale Luke AFB',
      country: 'United States',
      id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
      postalCode: '85309',
      state: 'AZ',
      streetAddress1: 'n/a',
      county: 'MARICOPA',
    }),
  ),
}));

const useEditShipmentQueriesReturnValue = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: 'NEEDS SERVICE COUNSELING',
  },
  order: {
    id: '1',
    originDutyLocation: {
      address: {
        streetAddress1: '',
        city: 'Fort Knox',
        state: 'KY',
        postalCode: '40121',
        county: 'HARDIN',
      },
    },
    destinationDutyLocation: {
      address: {
        streetAddress1: '',
        city: 'Fort Irwin',
        state: 'CA',
        postalCode: '92310',
        county: 'SAN BERNARDINO',
      },
    },
    customer: {
      agency: 'ARMY',
      backup_contact: {
        email: 'email@example.com',
        name: 'name',
        phone: '555-555-5555',
      },
      current_address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41Mzg0Njha',
        id: '3a5f7cf2-6193-4eb3-a244-14d21ca05d7b',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
        county: 'LOS ANGELES',
      },
      dodID: '6833908165',
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NjAzNTJa',
      email: 'combo@ppm.hhg',
      first_name: 'Submitted',
      id: 'f6bd793f-7042-4523-aa30-34946e7339c9',
      last_name: 'Ppmhhg',
      phone: '555-555-5555',
    },
    entitlement: {
      authorizedWeight: 8000,
      dependentsAuthorized: true,
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NzgwMzda',
      id: 'e0fefe58-0710-40db-917b-5b96567bc2a8',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 8000,
    },
    order_number: 'ORDER3',
    order_type: 'PERMANENT_CHANGE_OF_STATION',
    order_type_detail: 'HHG_PERMITTED',
    tac: '9999',
  },
  mtoShipments: [
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
        county: 'SOLANO',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'shipment123',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
        county: 'LOS ANGELES',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      requestedDeliveryDate: '2018-04-15',
      scheduledDeliveryDate: '2014-04-16',
      shipmentType: SHIPMENT_OPTIONS.HHG,
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const ppmShipment = {
  id: 'shipment123',
  shipmentType: SHIPMENT_OPTIONS.PPM,
  status: MOVE_STATUSES.SUBMITTED,
  updatedAt: '2020-09-02T21:08:38.392Z',
  ppmShipment: {
    expectedDepartureDate: '2022-06-28',
    actualMoveDate: '2022-05-11',
    hasSecondaryPickupAddress: true,
    hasSecondaryDestinationAddress: true,
    pickupAddress: {
      streetAddress1: '111 Test Street',
      streetAddress2: '222 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42701',
      county: 'HARDIN',
    },
    secondaryPickupAddress: {
      streetAddress1: '777 Test Street',
      streetAddress2: '888 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42702',
      county: 'HARDIN',
    },
    destinationAddress: {
      streetAddress1: '222 Test Street',
      streetAddress2: '333 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42703',
      county: 'HARDIN',
    },
    secondaryDestinationAddress: {
      streetAddress1: '444 Test Street',
      streetAddress2: '555 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      state: 'KY',
      postalCode: '42701',
      county: 'HARDIN',
    },
    sitExpected: false,
    estimatedWeight: 1111,
    hasProGear: false,
  },
};

const ppmShipmentWithSIT = {
  ...ppmShipment,
  ppmShipment: {
    ...ppmShipment.ppmShipment,
    sitExpected: true,
    sitEstimatedWeight: 999,
    sitEstimatedDepartureDate: '2022-07-13',
    sitEstimatedEntryDate: '2022-07-05',
  },
};

const loadingReturnValue = {
  ...useEditShipmentQueriesReturnValue,
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  ...useEditShipmentQueriesReturnValue,
  isLoading: false,
  isError: true,
  isSuccess: false,
};

const props = {
  onUpdate: () => {},
};

afterEach(() => {
  jest.resetAllMocks();
});

describe('ServicesCounselingEditShipmentDetails component', () => {
  describe('check different component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useEditShipmentQueries.mockReturnValue(loadingReturnValue);

      renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useEditShipmentQueries.mockReturnValue(errorReturnValue);

      renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  it('renders the Services Counseling Shipment Form', async () => {
    useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);

    renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

    const h1 = await screen.getByRole('heading', { name: 'Edit shipment details', level: 1 });
    await waitFor(() => {
      expect(h1).toBeInTheDocument();
    });
  });

  it('calls props.onUpdate with success and routes to move details when the save button is clicked and the shipment update is successful', async () => {
    updateMTOShipment.mockImplementation(() => Promise.resolve({}));
    useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
    const onUpdateMock = jest.fn();

    renderWithProviders(
      <ServicesCounselingEditShipmentDetails {...props} onUpdate={onUpdateMock} />,
      mockRoutingConfig,
    );

    const saveButton = screen.getByRole('button', { name: 'Save' });

    expect(saveButton).not.toBeDisabled();

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/counseling/moves/move123/details');
      expect(onUpdateMock).toHaveBeenCalledWith('success');
    });
  });

  it('stays on edit shipment form and displays error when the save button is clicked and the shipment update is unsuccessful', async () => {
    jest.spyOn(console, 'error').mockImplementation(() => {});
    updateMTOShipment.mockImplementation(() => Promise.reject(new Error('something went wrong')));
    useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);

    renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    expect(saveButton).not.toBeDisabled();

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(
        screen.getByText('Something went wrong, and your changes were not saved. Please try again.'),
      ).toBeVisible();
    }, 10000);
  });

  it('routes to the move details page when the cancel button is clicked', async () => {
    useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
    renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

    const cancelButton = screen.getByRole('button', { name: 'Cancel' });

    expect(cancelButton).not.toBeDisabled();

    await userEvent.click(cancelButton);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/counseling/moves/move123/details');
    });
  });

  describe('editing PPMs', () => {
    const ppmUseEditShipmentQueriesReturnValue = {
      ...useEditShipmentQueriesReturnValue,
      mtoShipments: [{ ...ppmShipment }],
    };

    const ppmWithSITUseEditShipmentQueriesReturnValue = {
      ...useEditShipmentQueriesReturnValue,
      mtoShipments: [{ ...ppmShipmentWithSIT }],
    };

    it('renders the first page of the edit ppm Shipment Form with prefilled values', async () => {
      useEditShipmentQueries.mockReturnValue(ppmUseEditShipmentQueriesReturnValue);
      renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

      expect(await screen.findByTestId('tag')).toHaveTextContent('PPM');
      expect(screen.getByRole('textbox', { name: 'Planned Departure Date' })).toHaveValue('28 Jun 2022');

      expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue(
        ppmShipment.ppmShipment.pickupAddress.streetAddress1,
      );
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue(
        ppmShipment.ppmShipment.pickupAddress.streetAddress2,
      );

      expect(screen.getAllByTestId('City')[0]).toHaveTextContent(ppmShipment.ppmShipment.pickupAddress.city);
      expect(screen.getAllByTestId('State')[0]).toHaveTextContent(ppmShipment.ppmShipment.pickupAddress.state);
      expect(screen.getAllByTestId('ZIP')[0]).toHaveTextContent(ppmShipment.ppmShipment.pickupAddress.postalCode);

      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue(
        ppmShipment.ppmShipment.secondaryPickupAddress.streetAddress1,
      );
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue(
        ppmShipment.ppmShipment.secondaryPickupAddress.streetAddress2,
      );
      expect(screen.getAllByTestId('City')[1]).toHaveTextContent(ppmShipment.ppmShipment.secondaryPickupAddress.city);
      expect(screen.getAllByTestId('State')[1]).toHaveTextContent(ppmShipment.ppmShipment.secondaryPickupAddress.state);
      expect(screen.getAllByTestId('ZIP')[1]).toHaveTextContent(
        ppmShipment.ppmShipment.secondaryPickupAddress.postalCode,
      );

      expect(screen.getAllByLabelText(/Address 1/)[2]).toHaveValue(
        ppmShipment.ppmShipment.destinationAddress.streetAddress1,
      );
      expect(screen.getAllByLabelText(/Address 2/)[2]).toHaveValue(
        ppmShipment.ppmShipment.destinationAddress.streetAddress2,
      );
      expect(screen.getAllByTestId('City')[2]).toHaveTextContent(ppmShipment.ppmShipment.destinationAddress.city);
      expect(screen.getAllByTestId('State')[2]).toHaveTextContent(ppmShipment.ppmShipment.destinationAddress.state);
      expect(screen.getAllByTestId(/ZIP/)[2]).toHaveTextContent(ppmShipment.ppmShipment.destinationAddress.postalCode);

      expect(screen.getAllByLabelText(/Address 1/)[3]).toHaveValue(
        ppmShipment.ppmShipment.secondaryDestinationAddress.streetAddress1,
      );
      expect(screen.getAllByLabelText(/Address 2/)[3]).toHaveValue(
        ppmShipment.ppmShipment.secondaryDestinationAddress.streetAddress2,
      );
      expect(screen.getAllByTestId(/City/)[3]).toHaveTextContent(
        ppmShipment.ppmShipment.secondaryDestinationAddress.city,
      );
      expect(screen.getAllByTestId('State')[3]).toHaveTextContent(
        ppmShipment.ppmShipment.secondaryDestinationAddress.state,
      );
      expect(screen.getAllByTestId(/ZIP/)[3]).toHaveTextContent(
        ppmShipment.ppmShipment.secondaryDestinationAddress.postalCode,
      );

      expect(screen.queryByRole('textbox', { name: 'Estimated SIT weight' })).not.toBeInTheDocument();
      expect(screen.queryByRole('textbox', { name: 'Estimated storage start' })).not.toBeInTheDocument();
      expect(screen.queryByRole('textbox', { name: 'Estimated storage end' })).not.toBeInTheDocument();
      expect(await screen.findByRole('textbox', { name: 'Estimated PPM weight' })).toHaveValue('1,111');
      expect(screen.queryByRole('textbox', { name: 'Estimated pro-gear weight' })).not.toBeInTheDocument();
      expect(screen.queryByRole('textbox', { name: 'Estimated spouse pro-gear weight' })).not.toBeInTheDocument();
      expect(await screen.findByRole('button', { name: 'Save and Continue' })).toBeInTheDocument();
    });

    it('verify toggling from Yes to No to Yes restores PPM SIT prefilled values', async () => {
      useEditShipmentQueries.mockReturnValue(ppmWithSITUseEditShipmentQueriesReturnValue);
      searchTransportationOffices.mockImplementation(() => Promise.resolve(mockTransportationOffice));
      renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

      expect(await screen.findByTestId('tag')).toHaveTextContent('PPM');

      expect(await screen.queryByRole('textbox', { name: 'Estimated SIT weight' })).toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage start' })).toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage end' })).toBeInTheDocument();
      expect(await screen.findByRole('button', { name: 'Save and Continue' })).toBeInTheDocument();

      expect(await screen.findByRole('textbox', { name: 'Estimated SIT weight' })).toHaveValue('999');
      expect(await screen.findByRole('textbox', { name: 'Estimated storage start' })).toHaveValue('05 Jul 2022');
      expect(await screen.findByRole('textbox', { name: 'Estimated storage end' })).toHaveValue('13 Jul 2022');

      act(() => {
        const closeoutField = screen
          .getAllByRole('combobox')
          .find((comboBox) => comboBox.getAttribute('id') === 'closeoutOffice-input');

        userEvent.click(closeoutField);
        userEvent.keyboard('Altus{enter}');
      });

      await waitFor(() => {
        expect(screen.queryByRole('alert')).not.toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeDisabled();
      });

      // Input invalid date format will cause form to be invalid. save must be disabled.
      await userEvent.type(screen.getByLabelText('Estimated storage start'), 'FOOBAR');
      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeDisabled();
      });

      // Schema validation is fail state thus Save button is disabled. click No to hide
      // SIT related widget. Hiding SIT widget must reset schema because previous SIT related
      // schema failure is nolonger applicable.
      const sitExpected = document.getElementById('sitExpectedNo').parentElement;
      const sitExpectedNo = within(sitExpected).getByRole('radio', { name: 'No' });
      await userEvent.click(sitExpectedNo);

      // Verify No is really hiding SIT related inputs
      expect(await screen.queryByRole('textbox', { name: 'Estimated SIT weight' })).not.toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage start' })).not.toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage end' })).not.toBeInTheDocument();

      // Verify clicking Yes again will restore persisted data for each SIT related control.
      const sitExpected2 = document.getElementById('sitExpectedYes').parentElement;
      const sitExpectedYes = within(sitExpected2).getByRole('radio', { name: 'Yes' });
      await userEvent.click(sitExpectedYes);

      // Verify persisted values are restored to expected values.
      expect(await screen.findByRole('textbox', { name: 'Estimated SIT weight' })).toHaveValue('999');
      expect(await screen.findByRole('textbox', { name: 'Estimated storage start' })).toHaveValue('05 Jul 2022');
      expect(await screen.findByRole('textbox', { name: 'Estimated storage end' })).toHaveValue('13 Jul 2022');
      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeDisabled();
      });
    }, 10000);

    describe('Check SIT field validations', () => {
      it.each([
        [
          'sitEstimatedWeight',
          {
            sitEstimatedWeight: '-1',
            sitEstimatedEntryDate: '15 Jun 2022',
            sitEstimatedDepartureDate: '25 Jul 2022',
          },
          'Enter a weight greater than 0 lbs',
        ],
        [
          'sitEstimatedWeight',
          {
            sitEstimatedWeight: '0',
            sitEstimatedEntryDate: '15 Jun 2022',
            sitEstimatedDepartureDate: '25 Jul 2022',
          },
          'Enter a weight greater than 0 lbs',
        ],
        [
          'sitEstimatedWeight',
          {
            sitEstimatedWeight: '{Tab}',
            sitEstimatedEntryDate: '15 Jun 2022',
            sitEstimatedDepartureDate: '25 Jul 2022',
          },
          'Required',
        ],
        [
          'sitEstimatedEntryDate',
          { sitEstimatedWeight: '1050', sitEstimatedEntryDate: 'asdf', sitEstimatedDepartureDate: '25 Jul 2022' },
          'Enter a complete date in DD MMM YYYY format (day, month, year).',
        ],
        [
          'sitEstimatedDepartureDate',
          { sitEstimatedWeight: '1025', sitEstimatedEntryDate: '15 Jun 2022', sitEstimatedDepartureDate: 'asdf' },
          'Enter a complete date in DD MMM YYYY format (day, month, year).',
        ],
      ])(
        'Verify invalid %s field shows validation error',
        async (field, data, expectedError) => {
          useEditShipmentQueries.mockReturnValue(ppmUseEditShipmentQueriesReturnValue);
          renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

          const sitExpected = document.getElementById('sitExpectedYes').parentElement;
          const sitExpectedYes = within(sitExpected).getByRole('radio', { name: 'Yes' });
          await userEvent.click(sitExpectedYes);

          // The test is dependent on the ordering of these three lines, and I'm not sure why.
          // If either of the estimated storage dates is entered last, the test that puts an invalid value
          // in that field will fail. But if the estimated SIT weight comes last, everything works fine.
          await userEvent.type(screen.getByLabelText('Estimated storage start'), data.sitEstimatedEntryDate);
          await userEvent.type(screen.getByLabelText('Estimated storage end'), data.sitEstimatedDepartureDate);
          await userEvent.type(screen.getByLabelText('Estimated SIT weight'), data.sitEstimatedWeight);
          await userEvent.tab();

          await waitFor(
            () => {
              const alerts = screen.getAllByRole('alert');
              expect(alerts).toHaveLength(1);
              expect(alerts[0]).toHaveTextContent(expectedError);
            },
            { timeout: 10000 },
          );

          expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeDisabled();
          expect(screen.getByRole('alert').nextElementSibling.firstElementChild).toHaveAttribute('name', field);
        },
        20000,
      );
    });

    it('Enables Save and Continue button when sit required fields are filled in', async () => {
      useEditShipmentQueries.mockReturnValue(ppmUseEditShipmentQueriesReturnValue);
      searchTransportationOffices.mockImplementation(() => Promise.resolve(mockTransportationOffice));
      renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

      const sitExpected = document.getElementById('sitExpectedYes').parentElement;
      const sitExpectedYes = within(sitExpected).getByRole('radio', { name: 'Yes' });
      await userEvent.click(sitExpectedYes);
      await userEvent.type(screen.getByLabelText('Estimated SIT weight'), '1050');
      await userEvent.type(screen.getByLabelText('Estimated storage start'), '15 Jun 2022');
      await userEvent.type(screen.getByLabelText('Estimated storage end'), '25 Jun 2022');
      await userEvent.tab();
      await userEvent.type(screen.getByLabelText(/Closeout location/), 'Altus');
      await userEvent.click(await screen.findByText('Altus'));

      await waitFor(() => {
        expect(screen.queryByRole('alert')).not.toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save and Continue' })).not.toBeDisabled();
      });
    }, 10000);

    it('verify toggling from Yes to No to Yes restores PPM SIT prefilled values', async () => {
      useEditShipmentQueries.mockReturnValue(ppmWithSITUseEditShipmentQueriesReturnValue);
      searchTransportationOffices.mockImplementation(() => Promise.resolve(mockTransportationOffice));
      renderWithProviders(<ServicesCounselingEditShipmentDetails {...props} />, mockRoutingConfig);

      expect(await screen.findByTestId('tag')).toHaveTextContent('PPM');

      expect(await screen.queryByRole('textbox', { name: 'Estimated SIT weight' })).toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage start' })).toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage end' })).toBeInTheDocument();
      expect(await screen.findByRole('button', { name: 'Save and Continue' })).toBeInTheDocument();

      expect(await screen.findByRole('textbox', { name: 'Estimated SIT weight' })).toHaveValue('999');
      expect(await screen.findByRole('textbox', { name: 'Estimated storage start' })).toHaveValue('05 Jul 2022');
      expect(await screen.findByRole('textbox', { name: 'Estimated storage end' })).toHaveValue('13 Jul 2022');

      act(() => {
        const closeoutField = screen
          .getAllByRole('combobox')
          .find((comboBox) => comboBox.getAttribute('id') === 'closeoutOffice-input');

        userEvent.click(closeoutField);
        userEvent.keyboard('Altus{enter}');
      });

      await waitFor(() => {
        expect(screen.queryByRole('alert')).not.toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeDisabled();
      });

      // Input invalid date format will cause form to be invalid. save must be disabled.
      await userEvent.type(screen.getByLabelText('Estimated storage start'), 'FOOBAR');
      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeDisabled();
      });

      // Schema validation is fail state thus Save button is disabled. click No to hide
      // SIT related widget. Hiding SIT widget must reset schema because previous SIT related
      // schema failure is nolonger applicable.
      const sitExpected = document.getElementById('sitExpectedNo').parentElement;
      const sitExpectedNo = within(sitExpected).getByRole('radio', { name: 'No' });
      await userEvent.click(sitExpectedNo);

      // Verify No is really hiding SIT related inputs
      expect(await screen.queryByRole('textbox', { name: 'Estimated SIT weight' })).not.toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage start' })).not.toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage end' })).not.toBeInTheDocument();

      // Verify clicking Yes again will restore persisted data for each SIT related control.
      const sitExpected2 = document.getElementById('sitExpectedYes').parentElement;
      const sitExpectedYes = within(sitExpected2).getByRole('radio', { name: 'Yes' });
      await userEvent.click(sitExpectedYes);

      // Verify persisted values are restored to expected values.
      expect(await screen.findByRole('textbox', { name: 'Estimated SIT weight' })).toHaveValue('999');
      expect(await screen.findByRole('textbox', { name: 'Estimated storage start' })).toHaveValue('05 Jul 2022');
      expect(await screen.findByRole('textbox', { name: 'Estimated storage end' })).toHaveValue('13 Jul 2022');
      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeDisabled();
      });
    }, 10000);

    it('calls props.onUpdate with success and routes to Advance page when the save button is clicked and the shipment update is successful', async () => {
      useEditShipmentQueries.mockReturnValue(ppmUseEditShipmentQueriesReturnValue);
      updateMTOShipment.mockImplementation(() => Promise.resolve({}));
      updateMoveCloseoutOffice.mockImplementation(() => Promise.resolve({}));
      searchTransportationOffices.mockImplementation(() => Promise.resolve(mockTransportationOffice));
      validatePostalCode.mockImplementation(() => Promise.resolve(false));
      const onUpdateMock = jest.fn();

      renderWithProviders(
        <ServicesCounselingEditShipmentDetails {...props} onUpdate={onUpdateMock} />,
        mockRoutingConfig,
      );

      await waitFor(() => {
        expect(screen.getByLabelText('Estimated PPM weight')).toHaveValue('1,111');
      });
      await userEvent.type(screen.getByLabelText(/Closeout location/), 'Altus');
      await userEvent.click(await screen.findByText('Altus'));

      const saveButton = screen.getByRole('button', { name: 'Save and Continue' });
      expect(saveButton).not.toBeDisabled();

      await userEvent.click(saveButton);
      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/counseling/moves/move123/shipments/shipment123/advance');
        expect(onUpdateMock).toHaveBeenCalledWith('success');
      });
    });

    it('displays error when the save button is clicked and the closeout office update is unsuccessful', async () => {
      // don't freak out when we get a console.error
      jest.spyOn(console, 'error').mockImplementation(() => {});

      useEditShipmentQueries.mockReturnValue(ppmUseEditShipmentQueriesReturnValue);
      updateMTOShipment.mockImplementation(() => Promise.resolve({}));
      searchTransportationOffices.mockImplementation(() => Promise.resolve(mockTransportationOffice));
      updateMoveCloseoutOffice.mockImplementation(() => Promise.reject(new Error('something went wrong')));
      validatePostalCode.mockImplementation(() => Promise.resolve(false));
      const onUpdateMock = jest.fn();
      renderWithProviders(
        <ServicesCounselingEditShipmentDetails {...props} onUpdate={onUpdateMock} />,
        mockRoutingConfig,
      );

      await waitFor(() => {
        expect(screen.getByLabelText('Estimated PPM weight')).toHaveValue('1,111');
      });
      await userEvent.type(screen.getByLabelText(/Closeout location/), 'Altus');
      await userEvent.click(await screen.findByText('Altus'));

      const saveButton = screen.getByRole('button', { name: 'Save and Continue' });
      expect(saveButton).not.toBeDisabled();

      await userEvent.click(saveButton);
      await waitFor(() => {
        expect(
          screen.getByText('Something went wrong, and your changes were not saved. Please try again.'),
        ).toBeVisible();
      });
    });
  });
});

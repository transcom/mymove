import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { Formik } from 'formik';

import { primeSimulatorRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';
import { createPrimeMTOShipmentV3 } from 'services/primeApi';
import PrimeUIShipmentCreate from 'pages/PrimeUI/Shipment/PrimeUIShipmentCreate';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const moveCode = 'LR4T8V';
const moveId = '9c7b255c-2981-4bf8-839f-61c7458e2b4d';
const shipmentId = 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee';
const routingParams = {
  moveCode,
  moveCodeOrID: moveId,
  shipmentId,
};

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  createPrimeMTOShipmentV3: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(true)),
}));

const moveDetailsURL = generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID: moveId });

const initialValues = {
  shipmentType: '',

  // PPM
  counselorRemarks: '',
  ppmShipment: {
    expectedDepartureDate: '',
    sitExpected: false,
    sitLocation: '',
    sitEstimatedWeight: '',
    sitEstimatedEntryDate: '',
    sitEstimatedDepartureDate: '',
    estimatedWeight: '',
    hasProGear: false,
    proGearWeight: '',
    spouseProGearWeight: '',
    pickupAddress: {
      city: '',
      postalCode: '',
      state: '',
      streetAddress1: '',
    },
    destinationAddress: {
      city: '',
      postalCode: '',
      state: '',
      streetAddress1: '',
    },
    secondaryDeliveryAddress: {
      city: '',
      postalCode: '',
      state: '',
      streetAddress1: '',
    },
    secondaryPickupAddress: {
      city: '',
      postalCode: '',
      state: '',
      streetAddress1: '',
    },
    tertiaryDeliveryAddress: {
      city: '',
      postalCode: '',
      state: '',
      streetAddress1: '',
    },
    tertiaryPickupAddress: {
      city: '',
      postalCode: '',
      state: '',
      streetAddress1: '',
    },
    hasSecondaryPickupAddress: 'false',
    hasSecondaryDestinationAddress: 'false',
    hasTertiaryPickupAddress: 'false',
    hasTertiaryDestinationAddress: 'false',
  },

  // Boat Shipment
  boatShipment: {
    make: 'make',
    model: 'model',
    year: 1999,
    hasTrailer: true,
    isRoadworthy: true,
    lengthInFeet: 16,
    lengthInInches: 0,
    widthInFeet: 1,
    widthInInches: 1,
    heightInFeet: 1,
    heightInInches: 1,
  },

  // Mobile Home Shipment
  mobileHomeShipment: {
    make: 'mobile make',
    model: 'mobile model',
    year: 1999,
    lengthInFeet: 16,
    lengthInInches: 0,
    widthInFeet: 1,
    widthInInches: 1,
    heightInFeet: 1,
    heightInInches: 1,
  },

  // Other shipment types
  requestedPickupDate: '',
  estimatedWeight: '',
  pickupAddress: {},
  destinationAddress: {},
  secondaryDeliveryAddress: {
    city: '',
    postalCode: '',
    state: '',
    streetAddress1: '',
  },
  secondaryPickupAddress: {
    city: '',
    postalCode: '',
    state: '',
    streetAddress1: '',
  },
  tertiaryDeliveryAddress: {
    city: '',
    postalCode: '',
    state: '',
    streetAddress1: '',
  },
  tertiaryPickupAddress: {
    city: '',
    postalCode: '',
    state: '',
    streetAddress1: '',
  },
  hasSecondaryPickupAddress: 'false',
  hasSecondaryDestinationAddress: 'false',
  hasTertiaryPickupAddress: 'false',
  hasTertiaryDestinationAddress: 'false',
  diversion: '',
  divertedFromShipmentId: '',
};

const mockedComponent = (
  <MockProviders path={primeSimulatorRoutes.CREATE_SHIPMENT_PATH} params={routingParams}>
    <Formik initialValues={initialValues}>
      <form>
        <PrimeUIShipmentCreate setFlashMessage={jest.fn()} />
      </form>
    </Formik>
  </MockProviders>
);

describe('Create Shipment Page', () => {
  it('renders the page without errors', async () => {
    render(mockedComponent);

    expect(await screen.findByText('Shipment Type')).toBeInTheDocument();
  });

  it('navigates the user to the home page when the cancel button is clicked', async () => {
    render(mockedComponent);

    expect(await screen.findByText('Shipment Type')).toBeInTheDocument();

    const cancel = screen.getByRole('button', { name: 'Cancel' });
    await userEvent.click(cancel);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(moveDetailsURL);
    });
  });
});

describe('successful submission of form', () => {
  it('calls history router back to move details', async () => {
    createPrimeMTOShipmentV3.mockReturnValue({});

    render(mockedComponent);

    await userEvent.selectOptions(screen.getByLabelText('Shipment type'), 'HHG');

    const saveButton = await screen.getByRole('button', { name: 'Save' });

    expect(saveButton).not.toBeDisabled();
    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(moveDetailsURL);
    });
  });
});

describe('Error when submitting', () => {
  it('Correctly displays the unexpected server error window when an unusuable api error response is returned', async () => {
    createPrimeMTOShipmentV3.mockRejectedValue('malformed api error response');
    render(mockedComponent);

    waitFor(async () => {
      await userEvent.selectOptions(screen.getByLabelText('Shipment type'), 'HHG');

      const saveButton = await screen.getByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);
      expect(screen.getByText('Unexpected error')).toBeInTheDocument();
      expect(
        screen.getByText('An unknown error has occurred, please check the address values used'),
      ).toBeInTheDocument();
    });
  });

  it('Correctly displays the invalid fields in the error window when an api error response is returned', async () => {
    createPrimeMTOShipmentV3.mockRejectedValue({ body: { title: 'Error', invalidFields: { someField: true } } });
    render(mockedComponent);

    waitFor(async () => {
      await userEvent.selectOptions(screen.getByLabelText('Shipment type'), 'HHG');

      const saveButton = await screen.getByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);
      expect(screen.getByText('Prime API: Error')).toBeInTheDocument();
      expect(
        screen.getByText('An unknown error has occurred, please check the address values used'),
      ).toBeInTheDocument();
    });
  });

  it('Correctly displays a specific error message when an error response is returned', async () => {
    createPrimeMTOShipmentV3.mockRejectedValue({ body: { title: 'Error', detail: 'The data entered no good.' } });
    render(mockedComponent);

    waitFor(async () => {
      await userEvent.selectOptions(screen.getByLabelText('Shipment type'), 'HHG');

      const saveButton = await screen.getByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);
      expect(screen.getByText('Prime API: Error')).toBeInTheDocument();
      expect(screen.getByText('The data entered no good.')).toBeInTheDocument();
    });
  });
});

describe('Create PPM', () => {
  it('test destination address street 1 is OPTIONAL', async () => {
    createPrimeMTOShipmentV3.mockReturnValue({});

    render(mockedComponent);

    waitFor(async () => {
      await userEvent.selectOptions(screen.getByLabelText('Shipment type'), 'PPM');

      // Start controlled test case to verify everything is working.
      let input = await document.querySelector('input[name=ppmShipment.pickupAddress.streetAddress1]');
      expect(input).toBeInTheDocument();
      // enter required street 1 for pickup
      await userEvent.type(input, '123 Street');
      // clear
      await userEvent.clear(input);
      await userEvent.tab();
      // verify Required alert is displayed
      const requiredAlerts = screen.getByRole('alert');
      expect(requiredAlerts).toHaveTextContent('Required');
      // make valid again to clear alert
      await userEvent.type(input, '123 Street');

      // Verify destination address street 1 is OPTIONAL.
      input = await document.querySelector('input[name=ppmShipment.destinationAddress.streetAddress1]');
      expect(input).toBeInTheDocument();
      // enter something
      await userEvent.type(input, '123 Street');
      // clear
      await userEvent.clear(input);
      await userEvent.tab();
      // verify no validation is displayed after clearing destination address street 1 because it's OPTIONAL
      expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    });
  });
});

describe('Create Mobile Home', () => {
  it.each(['MOBILE_HOME', 'BOAT_TOW_AWAY', 'BOAT_HAUL_AWAY'])(
    'resets secondary and tertiary addresses when flags are not true for shipment type %s',
    async (shipmentType) => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      createPrimeMTOShipmentV3.mockReturnValue({});

      // Render the component
      render(mockedComponent);

      // Wait for the component to load
      waitFor(async () => {
        expect(screen.getByLabelText('Shipment type')).toBeInTheDocument();

        // Select shipment type
        await userEvent.selectOptions(screen.getByLabelText('Shipment type'), shipmentType);

        await userEvent.type(screen.getByLabelText('Requested pickup'), '01 Nov 2020');

        // Fill in required pickup and destination addresses
        let input = document.querySelector('input[name=pickupAddress.streetAddress1]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '123 Street');
        input = document.querySelector('input[name=pickupAddress.city]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, 'Folsom');
        input = document.querySelector('input[name=pickupAddress.state]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, 'CA');
        input = document.querySelector('input[name=pickupAddress.postalCode]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '95630');

        input = document.querySelector('input[name=destinationAddress.streetAddress1]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '456 Destination St');
        input = document.querySelector('input[name=destinationAddress.city]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, 'Bevy Hills');
        input = document.querySelector('input[name=destinationAddress.state]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, 'CA');
        input = document.querySelector('input[name=destinationAddress.postalCode]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '90210');

        // Enable and disable secondary and tertiary toggles
        const secondAddressToggle = document.querySelector('[data-testid="has-secondary-pickup"]');
        expect(secondAddressToggle).toBeInTheDocument();
        await userEvent.click(secondAddressToggle);

        input = await document.querySelector('input[name="secondaryPickupAddress.streetAddress1"]');
        expect(input).toBeInTheDocument();
        // enter required street 1 for pickup 2
        await userEvent.type(input, '123 Pickup Street 2');

        const thirdAddressToggle = document.querySelector('[data-testid="has-tertiary-pickup"]');
        expect(thirdAddressToggle).toBeInTheDocument();
        await userEvent.click(thirdAddressToggle);

        input = await document.querySelector('input[name="tertiaryPickupAddress.streetAddress1"]');
        expect(input).toBeInTheDocument();
        // enter required street 1 for pickup 2
        await userEvent.type(input, '123 Pickup Street 3');

        const disable2ndAddressToggle = document.querySelector('[data-testid="no-secondary-pickup"]');
        expect(disable2ndAddressToggle).toBeInTheDocument();
        await userEvent.click(disable2ndAddressToggle);

        // input boat/mobile home model info
        input = document.createElement('input[label="Year"]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '2023');

        input = document.createElement('input[label="Make"]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, 'Genesis');

        input = document.createElement('input[label="Model"]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, 'G70');

        // input boat/mobile home dimensions
        input = document.createElement('input[label="Length (Feet)"]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '10');

        input = document.createElement('input[label="Length (Inches)"]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '10');

        input = document.createElement('input[label="Width (Feet)"]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '10');

        input = document.createElement('input[label="Width (Inches)"]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '10');

        input = document.createElement('input[label="Height (Feet)"]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '10');

        input = document.createElement('input[label="Height (Inches)"]');
        expect(input).toBeInTheDocument();
        await userEvent.type(input, '10');

        // Submit the form
        const saveButton = screen.getByRole('button', { name: 'Save' });
        expect(saveButton).not.toBeDisabled();
        await userEvent.click(saveButton);

        // Verify that API call resets addresses when flags are not 'true'
        expect(createPrimeMTOShipmentV3).toHaveBeenCalledWith({
          body: expect.objectContaining({
            destinationAddress: null,
            diversion: null,
            divertedFromShipmentId: null,
            hasSecondaryDestinationAddress: false,
            hasSecondaryPickupAddress: false,
            hasTertiaryDestinationAddress: false,
            hasTertiaryPickupAddress: false,
            secondaryDestinationAddress: {},
            secondaryPickupAddress: {},
            tertiaryDestinationAddress: {},
            tertiaryPickupAddress: {},
            moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
            pickupAddress: null,
            primeEstimatedWeight: null,
            requestedPickupDate: null,
          }),
        });
      });
    },
  );
});

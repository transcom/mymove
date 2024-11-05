import React from 'react';
import { render, screen, waitFor, act } from '@testing-library/react';
import { Formik } from 'formik';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import { primeSimulatorRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import PrimeUIShipmentCreateForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentCreateForm';

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

const testPickupAddress = {
  city: 'Beverly Hills',
  country: 'US',
  eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTMyNDha',
  id: '14b1d10d-b34b-4ec5-80e6-69d885206a2a',
  postalCode: '90210',
  state: 'CA',
  streetAddress1: '123 Any Street',
  streetAddress2: 'P.O. Box 12345',
  streetAddress3: 'c/o Some Person',
};

const testDestinationAddress = {
  city: 'Venice',
  country: 'US',
  eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
  id: '672ff379-f6e3-48b4-a87d-796713f8f997',
  postalCode: '90292',
  state: 'CA',
  streetAddress1: '987 Any Avenue',
  streetAddress2: 'P.O. Box 9876',
  streetAddress3: 'c/o Some Person',
};

const moveCode = 'LR4T8V';
const moveId = '9c7b255c-2981-4bf8-839f-61c7458e2b4d';
const shipmentId = 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee';
const routingParams = {
  moveCode,
  moveCodeOrID: moveId,
  shipmentId,
};

const moveDetailsURL = generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID: moveId });

function renderShipmentCreateForm(props) {
  render(
    <MockProviders path={primeSimulatorRoutes.CREATE_SHIPMENT_PATH} params={routingParams}>
      <Formik initialValues={initialValues}>
        <form>
          <PrimeUIShipmentCreateForm {...props} />
        </form>
      </Formik>
    </MockProviders>,
  );
}

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

beforeEach(() => {
  renderShipmentCreateForm();
});

describe('PrimeUIShipmentCreateForm', () => {
  it('renders the initial form', async () => {
    isBooleanFlagEnabled.mockResolvedValue(false);
    expect(await screen.queryByText('BOAT_HAUL_AWAY')).not.toBeInTheDocument();
    expect(await screen.queryByText('BOAT_TOW_AWAY')).not.toBeInTheDocument();
    expect(await screen.queryByText('MOBILE_HOME')).not.toBeInTheDocument();
    expect(await screen.findByLabelText('Shipment type')).toBeInTheDocument();
    expect(await screen.queryByText('MOBILE_HOME')).not.toBeInTheDocument();
  });

  it('renders the initial form, selecting PPM and checkboxes', async () => {
    isBooleanFlagEnabled.mockResolvedValue(false);
    expect(await screen.queryByText('BOAT_HAUL_AWAY')).not.toBeInTheDocument();
    expect(await screen.queryByText('BOAT_TOW_AWAY')).not.toBeInTheDocument();
    expect(await screen.queryByText('MOBILE_HOME')).not.toBeInTheDocument();
    const shipmentTypeInput = await screen.findByLabelText('Shipment type');
    expect(shipmentTypeInput).toBeInTheDocument();

    // Make it a PPM.
    await userEvent.selectOptions(shipmentTypeInput, ['PPM']);

    // Make sure than an HHG-specific field is not visible.
    expect(await screen.queryByLabelText('Requested pickup')).not.toBeInTheDocument();

    expect(await screen.findByText('Dates')).toBeInTheDocument();
    expect(await screen.findByLabelText('Expected Departure Date')).toHaveValue(
      initialValues.ppmShipment.expectedDepartureDate,
    );

    expect(await screen.getAllByLabelText('Address 1')[0]).toHaveValue(
      initialValues.ppmShipment.pickupAddress.streetAddress1,
    );

    expect(await screen.getAllByLabelText('City')[0]).toHaveValue(initialValues.ppmShipment.pickupAddress.city);
    expect(await screen.getAllByLabelText('State')[0]).toHaveValue(initialValues.ppmShipment.pickupAddress.state);
    expect(await screen.getAllByLabelText('ZIP')[0]).toHaveValue(initialValues.ppmShipment.pickupAddress.postalCode);

    expect(await screen.getAllByLabelText('Address 1')[1]).toHaveValue(
      initialValues.ppmShipment.secondaryPickupAddress.streetAddress1,
    );
    expect(await screen.getAllByLabelText('City')[1]).toHaveValue(
      initialValues.ppmShipment.secondaryPickupAddress.city,
    );
    expect(await screen.getAllByLabelText('State')[1]).toHaveValue(
      initialValues.ppmShipment.secondaryPickupAddress.state,
    );
    expect(await screen.getAllByLabelText('ZIP')[1]).toHaveValue(
      initialValues.ppmShipment.secondaryPickupAddress.postalCode,
    );

    expect(await screen.findByText('Storage In Transit (SIT)')).toBeInTheDocument();
    const sitExpectedInput = await screen.findByLabelText('SIT Expected');
    expect(sitExpectedInput).not.toBeChecked();

    expect(await screen.findByText('Weights')).toBeInTheDocument();
    expect(await screen.findByLabelText('Estimated Weight (lbs)')).toHaveValue(
      initialValues.ppmShipment.estimatedWeight,
    );

    const hasProGearInput = await screen.findByLabelText('Has Pro Gear');
    expect(hasProGearInput).not.toBeChecked();

    expect(await screen.findByText('Remarks')).toBeInTheDocument();
    expect(await screen.findByLabelText('Counselor Remarks')).toHaveValue(initialValues.counselorRemarks);

    // Turn on SIT.
    await userEvent.click(sitExpectedInput);

    expect(await screen.findByLabelText('SIT Location')).toHaveValue(initialValues.ppmShipment.sitLocation);
    expect(await screen.findByLabelText('SIT Estimated Weight (lbs)')).toHaveValue(
      initialValues.ppmShipment.sitEstimatedWeight,
    );
    expect(await screen.findByLabelText('SIT Estimated Entry Date')).toHaveValue(
      initialValues.ppmShipment.sitEstimatedEntryDate,
    );
    expect(await screen.findByLabelText('SIT Estimated Departure Date')).toHaveValue(
      initialValues.ppmShipment.sitEstimatedDepartureDate,
    );

    // Turn on pro gear.
    await userEvent.click(hasProGearInput);

    expect(await screen.findByLabelText('Pro Gear Weight (lbs)')).toHaveValue(initialValues.ppmShipment.proGearWeight);
    expect(await screen.findByLabelText('Spouse Pro Gear Weight (lbs)')).toHaveValue(
      initialValues.ppmShipment.spouseProGearWeight,
    );
  });

  it.each(
    ['BOAT_HAUL_AWAY', 'BOAT_TOW_AWAY', 'MOBILE_HOME'],
    'renders the initial form, selects a Boat or Mobile Home shipment type, and shows correct fields',
    async (shipmentType) => {
      isBooleanFlagEnabled.mockResolvedValue(true); // Allow for testing of boats and mobile homes
      const shipmentTypeInput = await screen.findByLabelText('Shipment type');
      expect(shipmentTypeInput).toBeInTheDocument();

      // Select the boat or mobile home shipment type
      await userEvent.selectOptions(shipmentTypeInput, [shipmentType]);

      // Make sure that a PPM-specific field is not visible.
      expect(await screen.queryByLabelText('Expected Departure Date')).not.toBeInTheDocument();

      // Check for usual HHG fields
      expect(await screen.findByRole('heading', { name: 'Diversion', level: 2 })).toBeInTheDocument();
      expect(await screen.findByLabelText('Diversion')).not.toBeChecked();

      // Checking to make sure the text box isn't shown prior to clicking the box
      expect(screen.queryByTestId('divertedFromShipmentIdInput')).toBeNull();

      // Check the diversion box
      const diversionCheckbox = await screen.findByLabelText('Diversion');
      await userEvent.click(diversionCheckbox);

      // now the text input should be visible
      expect(await screen.findByTestId('divertedFromShipmentIdInput')).toBeInTheDocument();

      // Now check for a boat and mobile home shipment specific field
      expect(await screen.findByLabelText('Length (Feet)')).toBeVisible();
    },
  );

  it.each(
    ['BOAT_HAUL_AWAY', 'BOAT_TOW_AWAY'],
    'correct identifies if a boat shipment qualifies as a separate shipment via its dimensions',
    async (shipmentType) => {
      isBooleanFlagEnabled.mockResolvedValue(true);
      const shipmentTypeInput = await screen.findByLabelText('Shipment type');
      expect(shipmentTypeInput).toBeInTheDocument();

      // Select the boat shipment type
      await userEvent.selectOptions(shipmentTypeInput, [shipmentType]);

      // Fill in form so that we can check that the form correctly identifies eligible shipments via their dimensions
      act(() => {
        initialValues.pickupAddress = testPickupAddress;
        initialValues.destinationAddress = testDestinationAddress;
        initialValues.estimatedWeight = '2000';
        initialValues.diversion = false;
        initialValues.requestedPickupDate = '2024-01-01';
      });

      // Now submit and check that modal confirms that the shipment is eligible as a boat shipment
      const saveButton = await screen.findByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      // Ensure that the shipment was created successfully
      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith(moveDetailsURL);
      });
    },
  );

  it.each(
    ['BOAT_HAUL_AWAY', 'BOAT_TOW_AWAY'],
    'correct identifies if a boat shipment is too small to ship separately, and should instead be shipped with HHG',
    async (shipmentType) => {
      isBooleanFlagEnabled.mockResolvedValue(true);

      const shipmentTypeInput = await screen.findByLabelText('Shipment type');
      expect(shipmentTypeInput).toBeInTheDocument();

      // Select the boat shipment type
      await userEvent.selectOptions(shipmentTypeInput, [shipmentType]);

      // Fill in form so that we can check that the form correctly identifies eligible shipments via their dimensions
      act(() => {
        initialValues.pickupAddress = testPickupAddress;
        initialValues.destinationAddress = testDestinationAddress;
        initialValues.estimatedWeight = '2000';
        initialValues.diversion = false;
        initialValues.requestedPickupDate = '2024-01-01';

        initialValues.boatShipment.lengthInInches = 1; // Set length to be smaller than minimum for boat shipment
      });

      // Now submit and check that modal confirms that the shipment is eligible as a boat shipment
      const saveButton = await screen.findByRole('button', { name: 'Save' });

      expect(saveButton).not.toBeDisabled();
      await userEvent.click(saveButton);

      // Now check that the validation failed
      expect(
        await screen.findByText(
          /One of these criteria must be met for it to be a boat shipment: lengthInInches > 168, widthInInches > 82, or heightInInches > 77/,
        ),
      ).toBeVisible();
    },
  );

  it.each(['HHG', 'HHG_INTO_NTS_DOMESTIC', 'HHG_OUTOF_NTS_DOMESTIC'])(
    'renders the initial form, selecting %s',
    async (shipmentType) => {
      isBooleanFlagEnabled.mockResolvedValue(false);
      expect(await screen.queryByText('BOAT_HAUL_AWAY')).not.toBeInTheDocument();
      expect(await screen.queryByText('BOAT_TOW_AWAY')).not.toBeInTheDocument();
      const shipmentTypeInput = await screen.findByLabelText('Shipment type');
      expect(shipmentTypeInput).toBeInTheDocument();

      // Select the shipment type
      await userEvent.selectOptions(shipmentTypeInput, [shipmentType]);

      // Make sure than a PPM-specific field is not visible.
      expect(await screen.queryByLabelText('Expected Departure Date')).not.toBeInTheDocument();

      expect(await screen.findByText('Shipment Dates')).toBeInTheDocument();
      expect(await screen.findByLabelText('Requested pickup')).toHaveValue(initialValues.requestedPickupDate);

      expect(await screen.findByRole('heading', { name: 'Diversion', level: 2 })).toBeInTheDocument();
      expect(await screen.findByLabelText('Diversion')).not.toBeChecked();

      expect(await screen.findByText('Shipment Weights')).toBeInTheDocument();
      expect(await screen.findByLabelText('Estimated weight (lbs)')).toHaveValue(initialValues.estimatedWeight);

      expect(await screen.findByText('Shipment Addresses')).toBeInTheDocument();
      expect(await screen.findByText('Pickup Address')).toBeInTheDocument();
      expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue('');

      expect(await screen.findByText('Delivery Address')).toBeInTheDocument();
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('');
    },
  );

  it('renders secondary/tertiary address', async () => {
    renderShipmentCreateForm();

    const shipmentTypeInput = await screen.findByLabelText('Shipment type');
    expect(shipmentTypeInput).toBeInTheDocument();

    // Select the shipment type
    await userEvent.selectOptions(shipmentTypeInput, 'HHG');

    // Make sure than a PPM-specific field is not visible.
    expect(await screen.queryByLabelText('Expected Departure Date')).not.toBeInTheDocument();

    expect(await screen.findByText('Shipment Dates')).toBeInTheDocument();
    expect(await screen.findByLabelText('Requested pickup')).toHaveValue(initialValues.requestedPickupDate);

    expect(await screen.findByRole('heading', { name: 'Diversion', level: 2 })).toBeInTheDocument();
    expect(await screen.findByLabelText('Diversion')).not.toBeChecked();

    expect(await screen.findByText('Shipment Weights')).toBeInTheDocument();
    expect(await screen.findByLabelText('Estimated weight (lbs)')).toHaveValue(initialValues.estimatedWeight);

    expect(await screen.findByText('Shipment Addresses')).toBeInTheDocument();
    expect(await screen.findByText('Pickup Address')).toBeInTheDocument();
    expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue('');

    const hasSecondaryPickup = await screen.findByTestId('has-secondary-pickup');
    await userEvent.click(hasSecondaryPickup);
    expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('');

    const hasTertiaryPickup = await screen.findByTestId('has-tertiary-pickup');
    await userEvent.click(hasTertiaryPickup);
    expect(screen.getAllByLabelText('Address 1')[2]).toHaveValue('');

    expect(await screen.findByText('Delivery Address')).toBeInTheDocument();
    expect(screen.getAllByLabelText('Address 1')[3]).toHaveValue('');

    const hasSecondaryDestination = await screen.findByTestId('has-secondary-destination');
    await userEvent.click(hasSecondaryDestination);
    expect(screen.getAllByLabelText('Address 1')[4]).toHaveValue('');

    const hasTertiaryDestination = await screen.findByTestId('has-tertiary-destination');
    await userEvent.click(hasTertiaryDestination);
    expect(screen.getAllByLabelText('Address 1')[5]).toHaveValue('');
  });

  it('renders the HHG form and displays the shipment id text input when diversion box is checked', async () => {
    isBooleanFlagEnabled.mockResolvedValue(false);
    expect(await screen.queryByText('BOAT_HAUL_AWAY')).not.toBeInTheDocument();
    expect(await screen.queryByText('BOAT_TOW_AWAY')).not.toBeInTheDocument();
    const shipmentTypeInput = await screen.findByLabelText('Shipment type');
    expect(shipmentTypeInput).toBeInTheDocument();

    // Make it a HHG move
    await userEvent.selectOptions(shipmentTypeInput, ['HHG']);

    expect(await screen.findByRole('heading', { name: 'Diversion', level: 2 })).toBeInTheDocument();
    expect(await screen.findByLabelText('Diversion')).not.toBeChecked();

    // Checking to make sure the text box isn't shown prior to clicking the box
    expect(screen.queryByTestId('divertedFromShipmentIdInput')).toBeNull();

    // Check the diversion box
    const diversionCheckbox = await screen.findByLabelText('Diversion');
    await userEvent.click(diversionCheckbox);

    // now the text input should be visible
    expect(await screen.findByTestId('divertedFromShipmentIdInput')).toBeInTheDocument();

    // Uncheck
    await userEvent.click(diversionCheckbox);

    // now the text input should be invisible
    expect(await screen.queryByTestId('divertedFromShipmentIdInput')).toBeNull();
  });
});

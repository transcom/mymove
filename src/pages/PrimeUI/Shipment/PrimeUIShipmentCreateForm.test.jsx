import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';
import userEvent from '@testing-library/user-event';

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
    hasSecondaryPickupAddress: 'false',
    hasSecondaryDestinationAddress: 'false',
  },

  // Other shipment types
  requestedPickupDate: '',
  estimatedWeight: '',
  pickupAddress: {},
  destinationAddress: {},
  diversion: '',
  divertedFromShipmentId: '',
};

function renderShipmentCreateForm(props) {
  render(
    <Formik initialValues={initialValues}>
      <form>
        <PrimeUIShipmentCreateForm {...props} />
      </form>
    </Formik>,
  );
}

describe('PrimeUIShipmentCreateForm', () => {
  it('renders the initial form', async () => {
    renderShipmentCreateForm();

    expect(await screen.findByLabelText('Shipment type')).toBeInTheDocument();
  });

  it('renders the initial form, selecting PPM and checkboxes', async () => {
    renderShipmentCreateForm();

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

  it.each(['HHG', 'HHG_INTO_NTS_DOMESTIC', 'HHG_OUTOF_NTS_DOMESTIC'])(
    'renders the initial form, selecting %s',
    async (shipmentType) => {
      renderShipmentCreateForm();

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

      expect(await screen.findByText('Destination Address')).toBeInTheDocument();
      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('');
    },
  );

  it('renders the HHG form and displays the shipment id text input when diversion box is checked', async () => {
    renderShipmentCreateForm();

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
  });
});

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
    pickupPostalCode: '',
    secondaryPickupPostalCode: '',
    destinationPostalCode: '',
    secondaryDestinationPostalCode: '',
    sitExpected: false,
    sitLocation: '',
    sitEstimatedWeight: '',
    sitEstimatedEntryDate: '',
    sitEstimatedDepartureDate: '',
    estimatedWeight: '',
    hasProGear: false,
    proGearWeight: '',
    spouseProGearWeight: '',
  },

  // Other shipment types
  requestedPickupDate: '',
  estimatedWeight: '',
  pickupAddress: {},
  destinationAddress: {},
  diversion: '',
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

    expect(await screen.findByText('Origin Info')).toBeInTheDocument();
    expect(await screen.findByLabelText('Pickup Postal Code')).toHaveValue(initialValues.ppmShipment.pickupPostalCode);
    expect(await screen.findByLabelText('Secondary Pickup Postal Code')).toHaveValue(
      initialValues.ppmShipment.secondaryPickupPostalCode,
    );

    expect(await screen.findByText('Destination Info')).toBeInTheDocument();
    expect(await screen.findByLabelText('Destination Postal Code')).toHaveValue(
      initialValues.ppmShipment.destinationPostalCode,
    );
    expect(await screen.findByLabelText('Secondary Destination Postal Code')).toHaveValue(
      initialValues.ppmShipment.secondaryDestinationPostalCode,
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

  it('renders the initial form, selecting HHG', async () => {
    renderShipmentCreateForm();

    const shipmentTypeInput = await screen.findByLabelText('Shipment type');
    expect(shipmentTypeInput).toBeInTheDocument();

    // Make it an HHG.
    await userEvent.selectOptions(shipmentTypeInput, ['HHG']);

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
  });
});

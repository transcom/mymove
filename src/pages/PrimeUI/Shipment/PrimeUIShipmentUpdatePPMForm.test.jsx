import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import PrimeUIShipmentUpdatePPMForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentUpdatePPMForm';
import { formatCustomerDate } from 'utils/formatters';

const shipment = {
  actualPickupDate: null,
  approvedDate: null,
  counselorRemarks: 'These are counselor remarks for a PPM.',
  createdAt: '2022-07-01T13:41:33.261Z',
  destinationAddress: {
    city: null,
    postalCode: null,
    state: null,
    streetAddress1: null,
  },
  eTag: 'MjAyMi0wNy0wMVQxNDoyMzoxOS43MzgzODla',
  firstAvailableDeliveryDate: null,
  id: '1b695b60-c3ed-401b-b2e3-808d095eb8cc',
  moveTaskOrderID: '7024c8c5-52ca-4639-bf69-dd8238308c98',
  pickupAddress: {
    city: null,
    postalCode: null,
    state: null,
    streetAddress1: null,
  },
  ppmShipment: {
    actualDestinationPostalCode: '30814',
    actualMoveDate: '2022-07-13',
    actualPickupPostalCode: '90212',
    advanceAmountReceived: 598600,
    advanceAmountRequested: 598700,
    approvedAt: '2022-07-03T14:20:21.620Z',
    createdAt: '2022-06-30T13:41:33.265Z',
    destinationPostalCode: '30813',
    eTag: 'MjAyMi0wNy0wMVQxNDoyMzoxOS43ODA1Mlo=',
    estimatedIncentive: 1000000,
    estimatedWeight: 4000,
    expectedDepartureDate: '2020-03-15',
    hasProGear: true,
    hasReceivedAdvance: true,
    hasRequestedAdvance: true,
    id: 'd733fe2f-b08d-434a-ad8d-551f4d597b03',
    netWeight: 3900,
    pickupPostalCode: '90210',
    proGearWeight: 1987,
    reviewedAt: '2022-07-02T14:20:14.636Z',
    secondaryDestinationPostalCode: '30814',
    secondaryPickupPostalCode: '90211',
    shipmentId: '1b695b60-c3ed-401b-b2e3-808d095eb8cc',
    sitEstimatedCost: 123456,
    sitEstimatedDepartureDate: '2022-07-13',
    sitEstimatedEntryDate: '2022-07-05',
    sitEstimatedWeight: 1100,
    sitExpected: true,
    sitLocation: 'DESTINATION',
    spouseProGearWeight: 498,
    status: 'SUBMITTED',
    submittedAt: '2022-07-01T13:41:33.252Z',
    updatedAt: '2022-07-01T14:23:19.780Z',
  },
  primeEstimatedWeightRecordedDate: null,
  requestedPickupDate: null,
  requiredDeliveryDate: null,
  scheduledPickupDate: null,
  secondaryDeliveryAddress: {
    city: null,
    postalCode: null,
    state: null,
    streetAddress1: null,
  },
  secondaryPickupAddress: {
    city: null,
    postalCode: null,
    state: null,
    streetAddress1: null,
  },
  shipmentType: 'PPM',
  status: 'APPROVED',
  updatedAt: '2022-07-01T14:23:19.738Z',
  mtoServiceItems: [],
};

const counselorRemarks = 'These are counselor remarks.';

const initialValues = {
  ppmShipment: {
    ...shipment.ppmShipment,
    sitEstimatedWeight: shipment.ppmShipment.sitEstimatedWeight?.toLocaleString(),
    estimatedWeight: shipment.ppmShipment.estimatedWeight?.toLocaleString(),
    netWeight: shipment.ppmShipment.netWeight?.toLocaleString(),
    proGearWeight: shipment.ppmShipment.proGearWeight?.toLocaleString(),
    spouseProGearWeight: shipment.ppmShipment.spouseProGearWeight?.toLocaleString(),
  },
  counselorRemarks,
};

function renderShipmentUpdatePPMForm(props) {
  render(
    <Formik initialValues={initialValues}>
      <form>
        <PrimeUIShipmentUpdatePPMForm {...props} />
      </form>
    </Formik>,
  );
}

describe('PrimeUIShipmentUpdatePPMForm', () => {
  it('renders the form', async () => {
    renderShipmentUpdatePPMForm();

    expect(await screen.findByText('Dates')).toBeInTheDocument();
    expect(await screen.findByLabelText('Expected Departure Date')).toHaveValue(
      formatCustomerDate(initialValues.ppmShipment.expectedDepartureDate),
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
    expect(await screen.findByLabelText('SIT Expected')).toBeChecked();
    expect(await screen.findByLabelText('SIT Location')).toHaveValue(initialValues.ppmShipment.sitLocation);
    expect(await screen.findByLabelText('SIT Estimated Weight (lbs)')).toHaveValue(
      initialValues.ppmShipment.sitEstimatedWeight,
    );
    expect(await screen.findByLabelText('SIT Estimated Entry Date')).toHaveValue(
      formatCustomerDate(initialValues.ppmShipment.sitEstimatedEntryDate),
    );
    expect(await screen.findByLabelText('SIT Estimated Departure Date')).toHaveValue(
      formatCustomerDate(initialValues.ppmShipment.sitEstimatedDepartureDate),
    );

    expect(await screen.findByText('Weights')).toBeInTheDocument();
    expect(await screen.findByLabelText('Estimated Weight (lbs)')).toHaveValue(
      initialValues.ppmShipment.estimatedWeight,
    );
    expect(await screen.findByLabelText('Net Weight (lbs)')).toHaveValue(initialValues.ppmShipment.netWeight);
    expect(await screen.findByLabelText('Has Pro Gear')).toBeChecked();
    expect(await screen.findByLabelText('Pro Gear Weight (lbs)')).toHaveValue(initialValues.ppmShipment.proGearWeight);
    expect(await screen.findByLabelText('Spouse Pro Gear Weight (lbs)')).toHaveValue(
      initialValues.ppmShipment.spouseProGearWeight,
    );

    expect(await screen.findByText('Remarks')).toBeInTheDocument();
    expect(await screen.findByLabelText('Counselor Remarks')).toHaveValue(initialValues.counselorRemarks);
  });
});

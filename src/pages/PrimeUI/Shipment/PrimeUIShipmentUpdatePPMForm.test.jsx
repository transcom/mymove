import React from 'react';
import { render, screen } from '@testing-library/react';
import { Formik } from 'formik';

import PrimeUIShipmentUpdatePPMForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentUpdatePPMForm';
import { formatCustomerDate } from 'utils/formatters';
import { MockProviders } from 'testUtils';

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
    eTag: 'MjAyMi0wNy0wMVQxNDoyMzoxOS43ODA1Mlo=',
    estimatedIncentive: 1000000,
    estimatedWeight: 4000,
    expectedDepartureDate: '2020-03-15',
    hasProGear: true,
    hasReceivedAdvance: true,
    hasRequestedAdvance: true,
    id: 'd733fe2f-b08d-434a-ad8d-551f4d597b03',
    proGearWeight: 1987,
    reviewedAt: '2022-07-02T14:20:14.636Z',
    shipmentId: '1b695b60-c3ed-401b-b2e3-808d095eb8cc',
    sitEstimatedCost: 123456,
    sitEstimatedDepartureDate: '2022-07-13',
    sitEstimatedEntryDate: '2022-07-05',
    sitEstimatedWeight: 1100,
    sitExpected: true,
    sitLocation: 'DESTINATION',
    spouseProGearWeight: 498,
    status: 'SUBMITTED',
    pickupAddress: {
      streetAddress1: '111 Test Street',
      streetAddress2: '222 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City 1',
      county: 'Murray 1',
      state: 'KY 1',
      postalCode: '42701',
    },
    secondaryPickupAddress: {
      streetAddress1: '777 Test Street',
      streetAddress2: '888 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City 2',
      county: 'Murray 2',
      state: 'KY 2',
      postalCode: '42702',
    },
    tertiaryPickupAddress: {
      streetAddress1: '123 Test Lane',
      streetAddress2: '234 Test Lane',
      streetAddress3: 'Test Woman',
      city: 'Missoula',
      county: 'Missoula',
      state: 'MT',
      postalCode: '59801',
    },
    destinationAddress: {
      streetAddress1: '222 Test Street',
      streetAddress2: '333 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City 3',
      county: 'Murray 3',
      state: 'KY 3',
      postalCode: '42703',
    },
    secondaryDestinationAddress: {
      streetAddress1: '444 Test Street',
      streetAddress2: '555 Test Street',
      streetAddress3: 'Test Man',
      city: 'Test City',
      county: 'Murray 4',
      state: 'KY 4',
      postalCode: '42704',
    },
    tertiaryDestinationAddress: {
      streetAddress1: '321 Test Lane',
      streetAddress2: '432 Test Lane',
      streetAddress3: 'Test Woman',
      city: 'Silver Spring',
      county: 'Montgomery',
      state: 'MD',
      postalCode: '20906',
    },
    hasSecondaryPickupAddress: 'true',
    hasSecondaryDestinationAddress: 'true',
    hasTertiaryPickupAddress: 'true',
    hasTertiaryDestinationAddress: 'true',
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
  tertiaryDeliveryAddress: {
    city: null,
    postalCode: null,
    state: null,
    streetAddress1: null,
  },
  tertiaryPickupAddress: {
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
    proGearWeight: shipment.ppmShipment.proGearWeight?.toLocaleString(),
    spouseProGearWeight: shipment.ppmShipment.spouseProGearWeight?.toLocaleString(),
  },
  counselorRemarks,
};

function renderShipmentUpdatePPMForm(props) {
  render(
    <MockProviders>
      <Formik initialValues={initialValues}>
        <form>
          <PrimeUIShipmentUpdatePPMForm {...props} />
        </form>
      </Formik>
    </MockProviders>,
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

    expect(await screen.getAllByLabelText('Address 1')[0]).toHaveValue(
      initialValues.ppmShipment.pickupAddress.streetAddress1,
    );
    expect(await screen.getAllByLabelText(/Address 2/)[0]).toHaveValue(
      initialValues.ppmShipment.pickupAddress.streetAddress2,
    );
    expect(await screen.getAllByLabelText(/Address 3/)[0]).toHaveValue(
      initialValues.ppmShipment.pickupAddress.streetAddress3,
    );
    expect(
      screen.getAllByText(
        `${initialValues.ppmShipment.pickupAddress.city}, ${initialValues.ppmShipment.pickupAddress.state} ${initialValues.ppmShipment.pickupAddress.postalCode} (${initialValues.ppmShipment.pickupAddress.county})`,
      ),
    );

    expect(await screen.getAllByLabelText('Address 1')[1]).toHaveValue(
      initialValues.ppmShipment.secondaryPickupAddress.streetAddress1,
    );
    expect(await screen.getAllByLabelText(/Address 2/)[1]).toHaveValue(
      initialValues.ppmShipment.secondaryPickupAddress.streetAddress2,
    );
    expect(await screen.getAllByLabelText(/Address 3/)[1]).toHaveValue(
      initialValues.ppmShipment.secondaryPickupAddress.streetAddress3,
    );
    expect(
      screen.getAllByText(
        `${initialValues.ppmShipment.secondaryPickupAddress.city}, ${initialValues.ppmShipment.secondaryPickupAddress.state} ${initialValues.ppmShipment.secondaryPickupAddress.postalCode} (${initialValues.ppmShipment.secondaryPickupAddress.county})`,
      ),
    );

    expect(await screen.getAllByLabelText('Address 1')[2]).toHaveValue(
      initialValues.ppmShipment.tertiaryPickupAddress.streetAddress1,
    );
    expect(await screen.getAllByLabelText(/Address 2/)[2]).toHaveValue(
      initialValues.ppmShipment.tertiaryPickupAddress.streetAddress2,
    );
    expect(await screen.getAllByLabelText(/Address 3/)[2]).toHaveValue(
      initialValues.ppmShipment.tertiaryPickupAddress.streetAddress3,
    );
    expect(
      screen.getAllByText(
        `${initialValues.ppmShipment.tertiaryPickupAddress.city}, ${initialValues.ppmShipment.tertiaryPickupAddress.state} ${initialValues.ppmShipment.tertiaryPickupAddress.postalCode} (${initialValues.ppmShipment.tertiaryPickupAddress.county})`,
      ),
    );

    expect(await screen.findByText('Destination Info')).toBeInTheDocument();

    expect(await screen.getAllByLabelText(/Address 1/)[3]).toHaveValue(
      initialValues.ppmShipment.destinationAddress.streetAddress1,
    );
    expect(await screen.getAllByLabelText(/Address 2/)[3]).toHaveValue(
      initialValues.ppmShipment.destinationAddress.streetAddress2,
    );
    expect(await screen.getAllByLabelText(/Address 3/)[3]).toHaveValue(
      initialValues.ppmShipment.destinationAddress.streetAddress3,
    );
    expect(
      screen.getAllByText(
        `${initialValues.ppmShipment.destinationAddress.city}, ${initialValues.ppmShipment.destinationAddress.state} ${initialValues.ppmShipment.destinationAddress.postalCode} (${initialValues.ppmShipment.destinationAddress.county})`,
      ),
    );

    expect(await screen.getAllByLabelText(/Address 1/)[4]).toHaveValue(
      initialValues.ppmShipment.secondaryDestinationAddress.streetAddress1,
    );
    expect(await screen.getAllByLabelText(/Address 2/)[4]).toHaveValue(
      initialValues.ppmShipment.secondaryDestinationAddress.streetAddress2,
    );
    expect(await screen.getAllByLabelText(/Address 3/)[4]).toHaveValue(
      initialValues.ppmShipment.secondaryDestinationAddress.streetAddress3,
    );

    expect(
      screen.getByText('Will the movers deliver any belongings to a third address?', {
        exact: false,
      }),
    ).toBeInTheDocument();

    expect(
      screen.getAllByText(
        `${initialValues.ppmShipment.secondaryDestinationAddress.city}, ${initialValues.ppmShipment.secondaryDestinationAddress.state} ${initialValues.ppmShipment.secondaryDestinationAddress.postalCode} (${initialValues.ppmShipment.secondaryDestinationAddress.county})`,
      ),
    );

    expect(await screen.getAllByLabelText(/Address 1/)[5]).toHaveValue(
      initialValues.ppmShipment.tertiaryDestinationAddress.streetAddress1,
    );
    expect(await screen.getAllByLabelText(/Address 2/)[5]).toHaveValue(
      initialValues.ppmShipment.tertiaryDestinationAddress.streetAddress2,
    );
    expect(await screen.getAllByLabelText(/Address 3/)[5]).toHaveValue(
      initialValues.ppmShipment.tertiaryDestinationAddress.streetAddress3,
    );
    expect(
      screen.getAllByText(
        `${initialValues.ppmShipment.tertiaryDestinationAddress.city}, ${initialValues.ppmShipment.tertiaryDestinationAddress.state} ${initialValues.ppmShipment.tertiaryDestinationAddress.postalCode} (${initialValues.ppmShipment.tertiaryDestinationAddress.county})`,
      ),
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
    expect(await screen.findByLabelText('Has Pro Gear')).toBeChecked();
    expect(await screen.findByLabelText('Pro Gear Weight (lbs)')).toHaveValue(initialValues.ppmShipment.proGearWeight);
    expect(await screen.findByLabelText('Spouse Pro Gear Weight (lbs)')).toHaveValue(
      initialValues.ppmShipment.spouseProGearWeight,
    );

    expect(await screen.findByText('Remarks')).toBeInTheDocument();
    expect(await screen.findByLabelText('Counselor Remarks')).toHaveValue(initialValues.counselorRemarks);
  });
});

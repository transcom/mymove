import { render, screen } from '@testing-library/react';
import React from 'react';

import { useViewEvaluationReportQueries } from '../../../hooks/queries';

import EvaluationReportContainer from './EvaluationReportContainer';

jest.mock('hooks/queries', () => ({
  useViewEvaluationReportQueries: jest.fn(),
}));

const setIsModalVisible = jest.fn();

const mtoShipments = [
  {
    actualPickupDate: '2020-03-16',
    createdAt: '2022-08-15T16:11:26.527Z',
    customerRemarks: 'Please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41MjI2MjRa',
      id: 'a1d0f3b7-87ad-45b1-b3a2-613b339aa8c1',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41MjcyMTha',
    id: '765a89c7-e59b-426a-a048-cf0ac91735bc',
    moveTaskOrderID: '01912827-b4e5-46cb-a800-4273830956cd',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41MTg2ODha',
      id: 'c717b236-3b3e-4b5d-b743-0811aa51ff1f',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    primeActualWeight: 980,
    requestedDeliveryDate: '2020-03-15',
    requestedPickupDate: '2020-03-15',
    scheduledPickupDate: '2020-03-16',
    shipmentType: 'HHG',
    status: 'SUBMITTED',
    updatedAt: '2022-08-15T16:11:26.527Z',
  },
  {
    createdAt: '2022-08-15T16:11:26.538Z',
    eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41MzgzMzda',
    id: '22ee7ad8-7147-4ee0-bbe6-264aef36d7cc',
    moveTaskOrderID: '01912827-b4e5-46cb-a800-4273830956cd',
    ppmShipment: {
      actualDestinationPostalCode: null,
      actualMoveDate: null,
      actualPickupPostalCode: null,
      advanceAmountReceived: null,
      advanceAmountRequested: 598700,
      approvedAt: null,
      createdAt: '2022-08-15T16:11:26.545Z',
      destinationPostalCode: '30813',
      eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41NDU2NTha',
      estimatedIncentive: 1000000,
      estimatedWeight: 4000,
      expectedDepartureDate: '2020-03-15',
      hasProGear: true,
      hasReceivedAdvance: null,
      hasRequestedAdvance: true,
      id: '62f0e96e-e562-4d22-91c8-a17237b9c609',
      netWeight: null,
      pickupPostalCode: '90210',
      proGearWeight: 1987,
      reviewedAt: null,
      secondaryDestinationPostalCode: '30814',
      secondaryPickupPostalCode: '90211',
      shipmentId: '22ee7ad8-7147-4ee0-bbe6-264aef36d7cc',
      sitEstimatedCost: null,
      sitEstimatedDepartureDate: null,
      sitEstimatedEntryDate: null,
      sitEstimatedWeight: null,
      sitExpected: false,
      spouseProGearWeight: 498,
      status: 'SUBMITTED',
      submittedAt: '2022-08-15T16:11:26.532Z',
      updatedAt: '2022-08-15T16:11:26.545Z',
      weightTickets: null,
    },
    requestedDeliveryDate: '0001-01-01',
    requestedPickupDate: '0001-01-01',
    shipmentType: 'PPM',
    status: 'SUBMITTED',
    updatedAt: '2022-08-15T16:11:26.538Z',
  },
  {
    createdAt: '2022-08-15T16:11:26.567Z',
    customerRemarks: 'Please treat gently',
    eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41Njc3MDFa',
    id: '34246c77-97f6-4a92-a682-52201ed3fbc4',
    moveTaskOrderID: '01912827-b4e5-46cb-a800-4273830956cd',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41NTg1MjNa',
      id: '298f4e72-b069-48bd-bb24-fbd6e6518896',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    requestedDeliveryDate: '0001-01-01',
    requestedPickupDate: '2020-03-15',
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41NjI1NjFa',
      id: 'de426328-3af5-4e3b-8734-c5ee0145f658',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: 'HHG_INTO_NTS_DOMESTIC',
    status: 'DRAFT',
    storageFacility: {
      address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41NDk1MTla',
        id: '375ff152-8f5a-4d92-a338-0832a9b4cece',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41NTQwMjRa',
      email: 'storage@email.com',
      facilityName: 'Storage R Us',
      id: '016e526f-c9c7-47a3-b742-d54ab546cadc',
      lotNumber: '1234',
      phone: '5555555555',
    },
    updatedAt: '2022-08-15T16:11:26.567Z',
  },
  {
    createdAt: '2022-08-15T16:11:26.583Z',
    customerRemarks: 'Please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41NzQ5OTVa',
      id: '2a9d2d70-3e29-40f4-a93d-838df2a1020d',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41ODM5NTha',
    id: '99371bdd-9cc6-4668-9c34-91a63bb8511e',
    moveTaskOrderID: '01912827-b4e5-46cb-a800-4273830956cd',
    requestedDeliveryDate: '2020-03-15',
    requestedPickupDate: '0001-01-01',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41Nzk0Mjda',
      id: '93297088-76df-41b9-88ab-ef94a3ee7563',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
    status: 'DRAFT',
    storageFacility: {
      address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41NDk1MTla',
        id: '375ff152-8f5a-4d92-a338-0832a9b4cece',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41NTQwMjRa',
      email: 'storage@email.com',
      facilityName: 'Storage R Us',
      id: '016e526f-c9c7-47a3-b742-d54ab546cadc',
      lotNumber: '1234',
      phone: '5555555555',
    },
    updatedAt: '2022-08-15T16:11:26.583Z',
  },
];
const evaluationReport = {
  createdAt: '2022-08-15T16:11:26.590Z',
  eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi41OTAyMjla',
  evaluationLengthMinutes: 45,
  id: '7d7180cd-f286-495d-a518-0d503e4cdc1f',
  inspectionDate: '2022-08-15',
  inspectionType: 'DATA_REVIEW',
  location: 'ORIGIN',
  moveID: '01912827-b4e5-46cb-a800-4273830956cd',
  moveReferenceID: '1018-3234',
  officeUser: {
    email: 'qae_csr_role@office.mil',
    firstName: 'Leo',
    id: 'ef4f6d1f-4ac3-4159-a364-5403e7d958ff',
    lastName: 'Spaceman',
    phone: '415-555-1212',
  },
  remarks: 'this is a submitted counseling report',
  submittedAt: '2022-08-15T16:11:26.589Z',
  type: 'COUNSELING',
  violationsObserved: false,
};

const customerInfo = {
  agency: 'ARMY',
  backup_contact: { email: 'email@example.com', name: 'name', phone: '555-555-5555' },
  current_address: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zMzIwOTFa',
    id: '28f11990-7ced-4d01-87ad-b16f2c86ea83',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
    streetAddress2: 'P.O. Box 12345',
    streetAddress3: 'c/o Some Person',
  },
  dodID: '5052247544',
  eTag: 'MjAyMi0wOC0xNVQxNjoxMToyNi4zNTkzNFo=',
  email: 'leo_spaceman_sm@example.com',
  first_name: 'Leo',
  id: 'ea557b1f-2660-4d6b-89a0-fb1b5efd2113',
  last_name: 'Spacemen',
  phone: '555-555-5555',
  userID: 'f4bbfcdf-ef66-4ce7-92f8-4c1bf507d596',
};

const viewReportReturn = { evaluationReport, mtoShipments };

describe('Evaluation Report Container', () => {
  it('renders the sample text', async () => {
    useViewEvaluationReportQueries.mockReturnValue(viewReportReturn);

    render(
      <EvaluationReportContainer
        evaluationReportId={evaluationReport.id}
        grade="E_4"
        setIsModalVisible={setIsModalVisible}
        moveCode="TEST123"
        customerInfo={customerInfo}
      />,
    );

    const evaluationReportContainer = await screen.findByTestId('EvaluationReportContainer');
    expect(evaluationReportContainer).toBeInTheDocument();
    // shipment type rendered
    expect(await screen.findByText('HHG')).toBeInTheDocument();
    // Title
    expect(await screen.findByText('Evaluation report')).toBeInTheDocument();
  });
});

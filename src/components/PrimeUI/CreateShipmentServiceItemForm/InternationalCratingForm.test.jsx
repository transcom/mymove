import React from 'react';
import { render, screen } from '@testing-library/react';

import InternationalCratingForm from './InternationalCratingForm';

const moveId = '9c7b255c-2981-4bf8-839f-61c7458e2b4d';

const approvedMoveTaskOrder = {
  moveTaskOrder: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    moveCode: 'LR4T8V',
    mtoShipments: [
      {
        actualPickupDate: '2020-03-17',
        agents: [],
        approvedDate: '2021-10-20',
        createdAt: '2021-10-21',
        customerRemarks: 'Please treat gently',
        destinationAddress: {
          city: 'Fairfield',
          id: 'bfe61147-5fd7-426e-b473-54ccf77bde35',
          postalCode: '94535',
          state: 'CA',
          streetAddress1: '987 Any Avenue',
          streetAddress2: 'P.O. Box 9876',
          streetAddress3: 'c/o Some Person',
        },
        eTag: 'MjAyMS0xMC0xOFQxODoyNDo0MS4zNzc5Nzha',
        firstAvailableDeliveryDate: null,
        id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        pickupAddress: {
          city: 'Beverly Hills',
          id: 'cf159eca-162c-4131-84a0-795e684416a6',
          postalCode: '90210',
          state: 'CA',
          streetAddress1: '123 Any Street',
          streetAddress2: 'P.O. Box 12345',
          streetAddress3: 'c/o Some Person',
        },
        primeActualWeight: 2000,
        primeEstimatedWeight: 1400,
        primeEstimatedWeightRecordedDate: null,
        requestedPickupDate: '2020-03-15',
        requiredDeliveryDate: null,
        scheduledPickupDate: '2020-03-16',
        secondaryDeliveryAddress: {
          city: null,
          postalCode: null,
          state: null,
          streetAddress1: null,
        },
        shipmentType: 'HHG',
        status: 'APPROVED',
        updatedAt: '2021-10-22',
        mtoServiceItems: null,
        reweigh: {
          id: '1234',
          weight: 9000,
          requestedAt: '2021-10-23',
        },
      },
    ],
  },
};

describe('InternationalCratingForm component', () => {
  it.each([
    ['Service item code', 'reServiceCode'],
    ['Item length (thousandths of an inch)', 'itemLength'],
    ['Item width (thousandths of an inch)', 'itemWidth'],
    ['Item height (thousandths of an inch)', 'itemHeight'],
    ['Crate length (thousandths of an inch)', 'crateLength'],
    ['Crate width (thousandths of an inch)', 'crateWidth'],
    ['Crate height (thousandths of an inch)', 'crateHeight'],
    ['Description', 'description'],
    ['Reason', 'reason'],
  ])('renders field %s in form', (labelName, inputName) => {
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];

    render(<InternationalCratingForm shipment={shipment} moveId={moveId} submission={jest.fn()} />);

    // shipment text values
    const field = screen.getByText(labelName);
    expect(field).toBeInTheDocument();
    expect(field.closest('div').nextElementSibling.name).toBe(inputName);
  });
});

import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DestinationSITServiceItemForm from './DestinationSITServiceItemForm';

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

describe('DestinationSITServiceItemForm component', () => {
  it.each([
    ['Reason', 'reason'],
    ['First available delivery date', 'firstAvailableDeliveryDate1'],
    ['First date of attempted contact', 'dateOfContact1'],
    ['First time of attempted contact', 'timeMilitary1'],
    ['Second available delivery date', 'firstAvailableDeliveryDate2'],
    ['Second date of attempted contact', 'dateOfContact2'],
    ['Second time of attempted contact', 'timeMilitary2'],
    ['SIT entry date', 'sitEntryDate'],
    ['SIT departure date', 'sitDepartureDate'],
  ])('renders field %s in form', (labelName) => {
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];

    render(<DestinationSITServiceItemForm shipment={shipment} submission={jest.fn()} />);

    const field = screen.getByText(labelName);
    expect(field).toBeInTheDocument();
  });

  it('renders hint component at bottom of page', async () => {
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];

    render(<DestinationSITServiceItemForm shipment={shipment} submission={jest.fn()} />);

    const hintInfo = screen.getByTestId('destinationSitInfo');
    expect(hintInfo).toBeInTheDocument();

    expect(hintInfo).toHaveTextContent(
      'The following service items will be created for domestic SIT: DDFSIT (Domestic Destination 1st day SIT) DDASIT (Domestic Destination additional days SIT) DDDSIT (Domestic Destination SIT delivery) DDSFSC (Domestic Destination SIT fuel surcharge) The following service items will be created for international SIT: IDFSIT (International Destination 1st day SIT) IDASIT (International Destination additional days SIT) IDDSIT (International Destination SIT delivery) IDSFSC (International Destination SIT fuel surcharge) NOTE: The above service items will use the current destination address of the shipment as their final destination address. Ensure the shipment address is accurate before creating these service items.',
    );
  });

  it('renders the Create Service Item button', async () => {
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];

    render(<DestinationSITServiceItemForm shipment={shipment} submission={jest.fn()} />);

    // Check if the button renders
    const createBtn = screen.getByRole('button', { name: 'Create service item' });
    expect(createBtn).toBeInTheDocument();
  });

  it('submits values when create service item button is clicked', async () => {
    const shipment = approvedMoveTaskOrder.moveTaskOrder.mtoShipments[0];
    const submissionMock = jest.fn();

    render(<DestinationSITServiceItemForm shipment={shipment} submission={submissionMock} />);

    await userEvent.type(screen.getByLabelText('Reason'), 'Testing');
    await userEvent.type(screen.getByLabelText('First available delivery date'), '01 Feb 2024');
    await userEvent.type(screen.getByLabelText('First date of attempted contact'), '28 Dec 2023');
    await userEvent.type(screen.getByLabelText('First time of attempted contact'), '1400Z');
    await userEvent.type(screen.getByLabelText('Second available delivery date'), '05 Feb 2024');
    await userEvent.type(screen.getByLabelText('Second date of attempted contact'), '05 Jan 2024');
    await userEvent.type(screen.getByLabelText('Second time of attempted contact'), '1400Z');
    await userEvent.type(screen.getByLabelText('SIT entry date'), '10 Jan 2024');
    await userEvent.type(screen.getByLabelText('SIT departure date'), '24 Jan 2024');

    // Submit form
    await userEvent.click(screen.getByRole('button', { name: 'Create service item' }));
    expect(submissionMock).toHaveBeenCalledTimes(1);
    expect(submissionMock).toHaveBeenCalledWith({
      body: {
        reason: 'Testing',
        dateOfContact1: '2023-12-28',
        dateOfContact2: '2024-01-05',
        firstAvailableDeliveryDate1: '2024-02-01',
        firstAvailableDeliveryDate2: '2024-02-05',
        modelType: 'MTOServiceItemDestSIT',
        moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
        mtoShipmentID: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
        reServiceCode: 'DDFSIT',
        sitDepartureDate: '2024-01-24',
        sitDestinationFinalAddress: null,
        sitEntryDate: '2024-01-10',
        timeMilitary1: '1400Z',
        timeMilitary2: '1400Z',
      },
    });
  });
});

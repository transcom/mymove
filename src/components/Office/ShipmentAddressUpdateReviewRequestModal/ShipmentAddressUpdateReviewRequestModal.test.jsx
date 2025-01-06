import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { ShipmentAddressUpdateReviewRequestModal } from './ShipmentAddressUpdateReviewRequestModal';

import { ADDRESS_UPDATE_STATUS, shipmentTypes } from 'constants/shipments';

const mockDeliveryAddressUpdate = {
  contractorRemarks: 'Test Contractor Remark',
  id: 'c49f7921-5a6e-46b4-bb39-022583574453',
  newAddress: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMy0wNy0xN1QxODowODowNi42NTU5MTVa',
    id: '6b57ce91-cabd-4e3b-9f48-ed4627d4878f',
    postalCode: '90210',
    state: 'CA',
    streetAddress1: '123 Any Street',
  },
  originalAddress: {
    city: 'Fairfield',
    country: 'US',
    id: '92509013-aafc-4892-a476-2e3b97e6933d',
    postalCode: '94535',
    state: 'CA',
    streetAddress1: '987 Any Avenue',
  },
  shipmentID: '5c84bcf3-92f7-448f-b0e1-e5378b6806df',
  status: 'REQUESTED',
};

const mockShipment = {
  id: '5c84bcf3-92f7-448f-b0e1-e5378b6806df',
  shipmentType: shipmentTypes.HHG,
  deliveryAddressUpdate: mockDeliveryAddressUpdate,
  eTag: 'eTag',
  mtoServiceItems: [
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:55.858Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1NS44NTgxMjVa',
      id: '7b7e94b1-0f34-418b-866f-d052e3a1c756',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DLH',
      reServiceID: '8d600f25-1def-422d-b159-617c7d59156e',
      reServiceName: 'Domestic linehaul',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:55.912Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1NS45MTI0NDFa',
      id: 'bf3516eb-1eaa-4e71-bd94-c523a6c866d0',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'FSC',
      reServiceID: '4780b30c-e846-437a-b39a-c499a6b09872',
      reServiceName: 'Fuel surcharge',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:55.968Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1NS45Njg1Nzda',
      id: '52b087b4-8e7f-4c96-939e-772cdd406e3a',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DOP',
      reServiceID: '2bc3e5cb-adef-46b1-bde9-55570bfdd43e',
      reServiceName: 'Domestic origin price',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:56.037Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1Ni4wMzc1OTla',
      id: 'c89ec6c0-a240-4478-afa0-52c5e2466ad4',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DDP',
      reServiceID: '50f1179a-3b72-4fa1-a951-fe5bcc70bd14',
      reServiceName: 'Domestic destination price',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:56.094Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1Ni4wOTQxMjRa',
      id: 'e26c9be3-dd55-4a0c-b002-f03258c40d06',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DPK',
      reServiceID: 'bdea5a8d-f15f-47d2-85c9-bba5694802ce',
      reServiceName: 'Domestic packing',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
    {
      approvedAt: '2023-12-14T19:10:55.840Z',
      createdAt: '2023-12-14T19:10:56.162Z',
      deletedAt: '0001-01-01',
      eTag: 'MjAyMy0xMi0xNFQxOToxMDo1Ni4xNjIzMTla',
      id: 'aca010a5-71e5-4994-b06b-97dfe4377f18',
      moveTaskOrderID: 'be44a6c6-55a2-4a36-8d8d-97e89a3b2043',
      mtoShipmentID: '3b4ecb78-0643-406f-ad74-8c1587bbba02',
      reServiceCode: 'DUPK',
      reServiceID: '15f01bc1-0754-4341-8e0f-25c8f04d5a77',
      reServiceName: 'Domestic unpacking',
      status: 'APPROVED',
      submittedAt: '0001-01-01',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
};

afterEach(() => {
  jest.resetAllMocks();
});

describe('ShipmentAddressUpdateReviewRequestModal', () => {
  it('renders the modal', async () => {
    render(
      <ShipmentAddressUpdateReviewRequestModal shipment={mockShipment} onSubmit={jest.fn()} onClose={jest.fn()} />,
    );

    await waitFor(() => {
      // Shipment type flag
      expect(screen.getByTestId('tag')).toHaveTextContent('HHG');

      // Heading
      expect(screen.getByRole('heading', { level: 2, name: 'Review request' })).toBeInTheDocument();

      // Address update preview component
      expect(screen.getByRole('heading', { level: 3, name: 'Delivery Address' })).toBeInTheDocument();

      // Form
      expect(screen.getByRole('heading', { level: 4, name: 'Review Request' })).toBeInTheDocument();

      // Form fields
      const approvalQuestion = screen.getByRole('group', { name: 'Approve address change?' });
      expect(approvalQuestion).toBeInTheDocument();
      const approvalYes = within(approvalQuestion).getByRole('radio', { name: 'Yes' });
      const approvalNo = within(approvalQuestion).getByRole('radio', { name: 'No' });
      expect(approvalYes).toBeInTheDocument();
      expect(approvalNo).toBeInTheDocument();

      expect(screen.getByLabelText('Office remarks')).toBeInTheDocument();
      expect(screen.getByText('Office remarks will be sent to the contractor.')).toBeInTheDocument();
      expect(screen.getByTestId('officeRemarks')).toBeInTheDocument();

      // Form buttons (Save and Cancel)
      expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    });
  });

  it('displays an errorMessage', () => {
    const errorText = 'An error!!';

    render(
      <ShipmentAddressUpdateReviewRequestModal
        shipment={mockShipment}
        errorMessage={errorText}
        onSubmit={jest.fn()}
        onClose={jest.fn()}
      />,
    );

    const errorMessage = screen.getByRole('alert');

    expect(errorMessage).toBeInTheDocument();
    expect(errorMessage).toHaveTextContent(errorText);
  });

  it('Runs an onClose callback on cancel', async () => {
    const user = userEvent.setup();

    const onClose = jest.fn();

    render(<ShipmentAddressUpdateReviewRequestModal shipment={mockShipment} onSubmit={jest.fn()} onClose={onClose} />);

    const cancel = screen.getByRole('button', { name: 'Cancel' });

    expect(cancel).toBeInTheDocument();

    await user.click(cancel);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('Runs an onClose callback on close', async () => {
    const user = userEvent.setup();

    const onClose = jest.fn();

    render(<ShipmentAddressUpdateReviewRequestModal shipment={mockShipment} onSubmit={jest.fn()} onClose={onClose} />);

    const close = screen.getByTestId('modalCloseButton');

    expect(close).toBeInTheDocument();

    await user.click(close);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('Runs an onSubmit callback on save', async () => {
    const user = userEvent.setup();

    const onSubmit = jest.fn();
    const officeRemarksAnswer = 'Here are my remarks from the office';

    render(<ShipmentAddressUpdateReviewRequestModal shipment={mockShipment} onSubmit={onSubmit} onClose={jest.fn()} />);

    const approvalQuestion = screen.getByRole('group', { name: 'Approve address change?' });
    const approvalYes = within(approvalQuestion).getByRole('radio', { name: 'Yes' });
    const officeRemarks = screen.getByLabelText('Office remarks');
    const save = screen.getByRole('button', { name: 'Save' });

    await user.click(approvalYes);
    await user.type(officeRemarks, officeRemarksAnswer);

    expect(approvalYes).toBeChecked();
    expect(officeRemarks).toHaveValue(officeRemarksAnswer);

    await user.click(save);

    expect(onSubmit).toHaveBeenCalledTimes(1);
    expect(onSubmit).toHaveBeenCalledWith(
      mockShipment.id,
      mockShipment.eTag,
      ADDRESS_UPDATE_STATUS.APPROVED,
      officeRemarksAnswer,
    );
  });
});

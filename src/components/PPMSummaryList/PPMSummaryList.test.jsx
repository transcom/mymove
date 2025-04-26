import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PPMSummaryList from './PPMSummaryList';

import { MockProviders, renderWithProviders } from 'testUtils';
import { downloadPPMPaymentPacket } from 'services/internalApi';
import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  downloadPPMPaymentPacket: jest.fn(),
}));

afterEach(() => {
  jest.resetAllMocks();
});

const shipments = [
  {
    id: '1',
    status: shipmentStatuses.SUBMITTED,
    ppmShipment: {
      id: '11',
      status: ppmShipmentStatuses.SUBMITTED,
      hasRequestedAdvance: true,
      advanceAmountRequested: 10000,
      pickupAddress: {
        streetAddress1: '1 Test Street',
        streetAddress2: '2 Test Street',
        streetAddress3: '3 Test Street',
        city: 'Pickup Test City',
        state: 'NY',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '1 Test Street',
        streetAddress2: '2 Test Street',
        streetAddress3: '3 Test Street',
        city: 'Destination Test City',
        state: 'NY',
        postalCode: '11111',
      },
    },
  },
  {
    id: '2',
    status: shipmentStatuses.APPROVED,
    ppmShipment: {
      id: '22',
      status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
      approvedAt: '2022-04-15T15:38:07.103Z',
      hasRequestedAdvance: true,
      advanceAmountRequested: 10000,
      pickupAddress: {
        streetAddress1: '1 Test Street',
        streetAddress2: '2 Test Street',
        streetAddress3: '3 Test Street',
        city: 'Pickup Test City',
        state: 'NY',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '1 Test Street',
        streetAddress2: '2 Test Street',
        streetAddress3: '3 Test Street',
        city: 'Destination Test City',
        state: 'NY',
        postalCode: '11111',
      },
    },
  },
  {
    id: '3',
    status: shipmentStatuses.APPROVED,
    ppmShipment: {
      id: '33',
      status: ppmShipmentStatuses.NEEDS_CLOSEOUT,
      approvedAt: '2022-04-15T15:38:07.103Z',
      submittedAt: '2022-04-19T15:38:07.103Z',
      hasRequestedAdvance: true,
      advanceAmountRequested: 10000,
      pickupAddress: {
        streetAddress1: '1 Test Street',
        streetAddress2: '2 Test Street',
        streetAddress3: '3 Test Street',
        city: 'Pickup Test City',
        state: 'NY',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '1 Test Street',
        streetAddress2: '2 Test Street',
        streetAddress3: '3 Test Street',
        city: 'Destination Test City',
        state: 'NY',
        postalCode: '11111',
      },
      movingExpenses: [
        {
          status: 'APPROVED',
        },
      ],
    },
  },
  {
    id: '4',
    status: shipmentStatuses.APPROVED,
    ppmShipment: {
      id: '44',
      status: ppmShipmentStatuses.CLOSEOUT_COMPLETE,
      approvedAt: '2022-04-15T15:38:07.103Z',
      submittedAt: '2022-04-19T15:38:07.103Z',
      reviewedAt: '2022-04-23T15:38:07.103Z',
      hasRequestedAdvance: true,
      advanceAmountRequested: 10000,
      pickupAddress: {
        streetAddress1: '1 Test Street',
        streetAddress2: '2 Test Street',
        streetAddress3: '3 Test Street',
        city: 'Pickup Test City',
        state: 'NY',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '1 Test Street',
        streetAddress2: '2 Test Street',
        streetAddress3: '3 Test Street',
        city: 'Destination Test City',
        state: 'NY',
        postalCode: '11111',
      },
      movingExpenses: [
        {
          status: 'REJECTED',
        },
      ],
    },
  },
];
const onUploadClick = jest.fn();
const onDownloadError = jest.fn();

const defaultProps = {
  shipments,
  onDownloadError,
  onUploadClick,
};

describe('PPMSummaryList component', () => {
  describe('pending approval from counselor', () => {
    const props = { ...defaultProps, shipments: shipments.slice(0, 1) };
    it('should have button disabled', () => {
      render(<PPMSummaryList {...props} />);
      expect(screen.getByRole('button', { name: 'Upload PPM Documents' })).toBeDisabled();
    });
  });
  describe('approved by counselor', () => {
    const props = { ...defaultProps, shipments: shipments.slice(1, 2) };
    it('should have button enabled', async () => {
      render(<PPMSummaryList {...props} />);
      const uploadButton = screen.getByRole('button', { name: 'Upload PPM Documents' });
      expect(uploadButton).toBeEnabled();
      await userEvent.click(uploadButton);
      expect(onUploadClick).toHaveBeenCalledWith(props.shipments[0].id); // called with mtoShipmentId
      expect(onUploadClick).toHaveBeenCalledTimes(1);
      expect(screen.getByText(/From:/, { selector: 'span' })).toBeInTheDocument();
      expect(screen.getByText(/Pickup Test City, NY 10001/, { selector: 'p' })).toBeInTheDocument();
      expect(screen.getByText(/To:/, { selector: 'span' })).toBeInTheDocument();
      expect(screen.getByText(/Destination Test City, NY 11111/, { selector: 'p' })).toBeInTheDocument();
    });
    it('should contain approved date', () => {
      render(<PPMSummaryList {...props} />);
      expect(screen.queryByText(`PPM approved: 15 Apr 2022.`)).toBeInTheDocument();

      expect(screen.queryByText(`PPM`)).toBeInTheDocument();
    });
  });

  describe('move is locked by office user', () => {
    const props = { ...defaultProps, shipments: shipments.slice(0, 1) };
    it('should disable any edit/delete buttons', () => {
      render(<PPMSummaryList {...props} isMoveLocked />);
      const uploadButton = screen.getByRole('button', { name: 'Upload PPM Documents' });
      expect(uploadButton).toBeDisabled();
    });
  });

  describe('payment docs submitted for closeout review', () => {
    it('should display submitted date and disabled button with copy', () => {
      render(<PPMSummaryList shipments={[shipments[2]]} />);
      expect(screen.getByRole('button', { name: 'Download Payment Packet' })).toBeDisabled();

      expect(screen.getByText(/From:/, { selector: 'span' })).toBeInTheDocument();
      expect(screen.getByText(/Pickup Test City, NY 10001/, { selector: 'p' })).toBeInTheDocument();
      expect(screen.getByText(/To:/, { selector: 'span' })).toBeInTheDocument();
      expect(screen.getByText(/Destination Test City, NY 11111/, { selector: 'p' })).toBeInTheDocument();

      expect(screen.queryByText(`PPM approved: 15 Apr 2022`)).toBeInTheDocument();
      expect(screen.queryByText(`PPM documentation submitted: 19 Apr 2022`)).toBeInTheDocument();

      expect(
        screen.queryByText(
          "A counselor will review your documentation. When it's verified, you can visit MilMove to download the incentive packet that you'll need to give to Finance.",
        ),
      ).toBeInTheDocument();
    });
  });

  describe('payment docs reviewed', () => {
    it('should display reviewed date and enabled button with copy', () => {
      renderWithProviders(<PPMSummaryList shipments={[shipments[3]]} />);
      expect(screen.getByRole('button', { name: 'Download Payment Packet' })).toBeEnabled();

      expect(screen.queryByText(`PPM approved: 15 Apr 2022`)).toBeInTheDocument();
      expect(screen.queryByText(`PPM documentation submitted: 19 Apr 2022`)).toBeInTheDocument();
      expect(screen.queryByText(`PPM closeout completed: 23 Apr 2022`)).toBeInTheDocument();

      expect(screen.getByText(/From:/, { selector: 'span' })).toBeInTheDocument();
      expect(screen.getByText(/Pickup Test City, NY 10001/, { selector: 'p' })).toBeInTheDocument();
      expect(screen.getByText(/To:/, { selector: 'span' })).toBeInTheDocument();
      expect(screen.getByText(/Destination Test City, NY 11111/, { selector: 'p' })).toBeInTheDocument();

      expect(
        screen.queryByText(
          'You can now download your incentive packet and submit it to Finance to request payment. You will also need to include a completed DD-1351-2, and any other paperwork required by your service.',
        ),
      ).toBeInTheDocument();
    });

    it('should display button for feedback if any document is not approved', () => {
      renderWithProviders(<PPMSummaryList shipments={[shipments[3]]} />);

      expect(screen.getByRole('button', { name: 'View Closeout Feedback' })).toBeEnabled();
    });
    it('should not display button for feedback if all documents are approved', () => {
      render(<PPMSummaryList shipments={[shipments[2]]} />);

      expect(screen.queryByRole('button', { name: 'View Closeout Feedback' })).not.toBeInTheDocument();
    });
  });

  describe('there is only one shipment', () => {
    it('should not render numbers next to PPM', () => {
      const props = { ...defaultProps, shipments: shipments.slice(0, 1) };
      render(<PPMSummaryList {...props} />);
      expect(screen.queryByText('PPM 1')).not.toBeInTheDocument();
    });
  });
  describe('there are multiple shipments', () => {
    it('should render numbers next to PPM', () => {
      renderWithProviders(<PPMSummaryList {...defaultProps} />);
      expect(screen.queryByText('PPM 1')).toBeInTheDocument();
      expect(screen.queryByText('PPM 2')).toBeInTheDocument();
    });
  });

  it('PPM Download Payment Packet - success', async () => {
    const mockResponse = {
      ok: true,
      headers: {
        'content-disposition': 'filename="test.pdf"',
      },
      status: 200,
      data: null,
    };
    downloadPPMPaymentPacket.mockImplementation(() => Promise.resolve(mockResponse));
    render(
      <MockProviders>
        <PPMSummaryList onDownloadError={onDownloadError} shipments={[shipments[3]]} />
      </MockProviders>,
    );

    expect(screen.getByText('Download Payment Packet', { exact: false })).toBeInTheDocument();

    const downloadPaymentButton = screen.getByText('Download Payment Packet');
    expect(downloadPaymentButton).toBeInTheDocument();

    await userEvent.click(downloadPaymentButton);

    await waitFor(() => {
      expect(downloadPPMPaymentPacket).toHaveBeenCalledTimes(1);
    });
  });

  it('PPM Download Payment Packet - failure', async () => {
    downloadPPMPaymentPacket.mockRejectedValue({
      response: { body: { title: 'Error title', detail: 'Error detail' } },
    });

    const shipment = {
      ppmShipment: {
        status: ppmShipmentStatuses.CLOSEOUT_COMPLETE,
        pickupAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Pickup Test City',
          state: 'NY',
          postalCode: '10001',
        },
        destinationAddress: {
          streetAddress1: '1 Test Street',
          streetAddress2: '2 Test Street',
          streetAddress3: '3 Test Street',
          city: 'Destination Test City',
          state: 'NY',
          postalCode: '11111',
        },
      },
    };

    render(
      <MockProviders>
        <PPMSummaryList shipments={[shipment]} onDownloadError={onDownloadError} />
      </MockProviders>,
    );

    expect(screen.getByText('Download Payment Packet')).toBeInTheDocument();

    const downloadPaymentButton = screen.getByText('Download Payment Packet');
    expect(downloadPaymentButton).toBeInTheDocument();
    await userEvent.click(downloadPaymentButton);

    await waitFor(() => {
      expect(downloadPPMPaymentPacket).toHaveBeenCalledTimes(1);
      expect(onDownloadError).toHaveBeenCalledTimes(1);
    });
  });
});

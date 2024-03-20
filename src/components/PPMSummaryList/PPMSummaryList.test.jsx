import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PPMSummaryList from './PPMSummaryList';

import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';

const shipments = [
  {
    id: '1',
    status: shipmentStatuses.SUBMITTED,
    ppmShipment: {
      id: '11',
      status: ppmShipmentStatuses.SUBMITTED,
      hasRequestedAdvance: true,
      advanceAmountRequested: 10000,
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
    },
  },
  {
    id: '3',
    status: shipmentStatuses.APPROVED,
    ppmShipment: {
      id: '33',
      status: ppmShipmentStatuses.NEEDS_PAYMENT_APPROVAL,
      approvedAt: '2022-04-15T15:38:07.103Z',
      submittedAt: '2022-04-19T15:38:07.103Z',
      hasRequestedAdvance: true,
      advanceAmountRequested: 10000,
    },
  },
  {
    id: '4',
    status: shipmentStatuses.APPROVED,
    ppmShipment: {
      id: '44',
      status: ppmShipmentStatuses.PAYMENT_APPROVED,
      approvedAt: '2022-04-15T15:38:07.103Z',
      submittedAt: '2022-04-19T15:38:07.103Z',
      reviewedAt: '2022-04-23T15:38:07.103Z',
      hasRequestedAdvance: true,
      advanceAmountRequested: 10000,
    },
  },
];
const onUploadClick = jest.fn();

const defaultProps = {
  shipments,
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
    });
    it('should contain approved date', () => {
      render(<PPMSummaryList {...props} />);
      expect(screen.queryByText(`PPM approved: 15 Apr 2022.`)).toBeInTheDocument();
    });
  });

  describe('payment docs submitted for closeout review', () => {
    it('should display submitted date and disabled button with copy', () => {
      render(<PPMSummaryList shipments={[shipments[2]]} />);
      expect(screen.getByRole('button', { name: 'Download Incentive Packet' })).toBeDisabled();

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
      render(<PPMSummaryList shipments={[shipments[3]]} />);
      expect(screen.getByRole('button', { name: 'Download Incentive Packet' })).toBeEnabled();

      expect(screen.queryByText(`PPM approved: 15 Apr 2022`)).toBeInTheDocument();
      expect(screen.queryByText(`PPM documentation submitted: 19 Apr 2022`)).toBeInTheDocument();
      expect(screen.queryByText(`Documentation accepted and verified: 23 Apr 2022`)).toBeInTheDocument();

      expect(
        screen.queryByText(
          'You can now download your incentive packet and submit it to Finance to request payment. You will also need to include a completed DD-1351-2, and any other paperwork required by your service.',
        ),
      ).toBeInTheDocument();
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
      render(<PPMSummaryList {...defaultProps} />);
      expect(screen.queryByText('PPM 1')).toBeInTheDocument();
      expect(screen.queryByText('PPM 2')).toBeInTheDocument();
    });
  });
});

import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PPMSummaryList from './PPMSummaryList';

import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';

const shipments = [
  {
    id: '1',
    status: shipmentStatuses.SUBMITTED,
    ppmShipment: { id: '11', status: ppmShipmentStatuses.SUBMITTED, advanceRequested: true, advance: 10000 },
  },
  {
    id: '2',
    status: shipmentStatuses.APPROVED,
    ppmShipment: {
      id: '22',
      status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
      approvedAt: '2022-04-15T15:38:07.103Z',
      advanceRequested: true,
      advance: 10000,
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
    it('should have button enabled', () => {
      render(<PPMSummaryList {...props} />);
      const uploadButton = screen.getByRole('button', { name: 'Upload PPM Documents' });
      expect(uploadButton).not.toBeDisabled();
      userEvent.click(uploadButton);
      expect(onUploadClick).toHaveBeenCalledWith(props.shipments[0].id); // called with mtoShipmentId
      expect(onUploadClick).toHaveBeenCalledTimes(1);
    });
    it('should contain approved date', () => {
      render(<PPMSummaryList {...props} />);
      expect(screen.queryByText(`PPM approved: 15 Apr 2022.`)).toBeInTheDocument();
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

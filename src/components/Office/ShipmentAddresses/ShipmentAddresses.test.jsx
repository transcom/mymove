import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentAddresses from './ShipmentAddresses';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

const testProps = {
  pickupAddress: {
    city: 'Fairfax',
    state: 'VA',
    postalCode: '12345',
    streetAddress1: '123 Fake Street',
    streetAddress2: '',
    streetAddress3: '',
    country: 'USA',
  },
  destinationAddress: {
    city: 'Boston',
    state: 'MA',
    postalCode: '01054',
    streetAddress1: '5 Main Street',
    streetAddress2: '',
    streetAddress3: '',
    country: 'USA',
  },
  destinationDutyLocation: {
    streetAddress1: '',
    city: 'Fort Irwin',
    state: 'CA',
    postalCode: '92310',
  },
  handleDivertShipment: jest.fn(),
  shipmentInfo: {
    id: '456',
    eTag: 'abc123',
    status: 'APPROVED',
    shipmentType: SHIPMENT_OPTIONS.HHG,
  },
};

const cancelledShipment = {
  pickupAddress: {
    city: 'Fairfax',
    state: 'VA',
    postalCode: '12345',
    streetAddress1: '123 Fake Street',
    streetAddress2: '',
    streetAddress3: '',
    country: 'USA',
  },
  destinationAddress: {
    city: 'Boston',
    state: 'MA',
    postalCode: '01054',
    streetAddress1: '5 Main Street',
    streetAddress2: '',
    streetAddress3: '',
    country: 'USA',
  },
  destinationDutyLocation: {
    streetAddress1: '',
    city: 'Fort Irwin',
    state: 'CA',
    postalCode: '92310',
  },
  handleDivertShipment: jest.fn(),
  shipmentInfo: {
    id: '456',
    eTag: 'abc123',
    status: 'CANCELED',
    shipmentType: SHIPMENT_OPTIONS.HHG,
  },
};

describe('ShipmentAddresses', () => {
  it('calls props.handleDivertShipment on request diversion button click', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createShipmentDiversionRequest]}>
        <ShipmentAddresses {...testProps} />
      </MockProviders>,
    );
    const requestDiversionBtn = screen.getByRole('button', { name: 'Request diversion' });

    userEvent.click(requestDiversionBtn);
    await waitFor(() => {
      expect(testProps.handleDivertShipment).toHaveBeenCalled();
      expect(testProps.handleDivertShipment).toHaveBeenCalledWith(
        testProps.shipmentInfo.id,
        testProps.shipmentInfo.eTag,
      );
    });
  });

  it('hides the request diversion button for a cancelled shipment', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createShipmentDiversionRequest]}>
        <ShipmentAddresses {...cancelledShipment} />
      </MockProviders>,
    );
    const requestDiversionBtn = screen.queryByRole('button', { name: 'Request diversion' });

    await waitFor(() => {
      expect(requestDiversionBtn).toBeNull();
    });
  });

  it('hides the request diversion button when user does not have permissions', async () => {
    render(<ShipmentAddresses {...cancelledShipment} />);
    const requestDiversionBtn = screen.queryByRole('button', { name: 'Request diversion' });

    await waitFor(() => {
      expect(requestDiversionBtn).toBeNull();
    });
  });

  it('shows correct headings for HHG', () => {
    render(<ShipmentAddresses {...testProps} />);
    expect(screen.getByText("Customer's addresses")).toBeInTheDocument();
  });

  it('shows correct headings for NTS', () => {
    const NTSProps = {
      ...testProps,
      shipmentInfo: { ...testProps.shipmentInfo, shipmentType: SHIPMENT_OPTIONS.NTS },
    };
    render(<ShipmentAddresses {...NTSProps} />);
    expect(screen.getByText('Pickup address')).toBeInTheDocument();
    expect(screen.getByText('Facility address')).toBeInTheDocument();
  });

  it('shows correct headings for NTSR', () => {
    const NTSRProps = {
      ...testProps,
      shipmentInfo: { ...testProps.shipmentInfo, shipmentType: SHIPMENT_OPTIONS.NTSR },
    };
    render(<ShipmentAddresses {...NTSRProps} />);
    expect(screen.getByText('Facility address')).toBeInTheDocument();
    expect(screen.getByText('Delivery address')).toBeInTheDocument();
  });
});

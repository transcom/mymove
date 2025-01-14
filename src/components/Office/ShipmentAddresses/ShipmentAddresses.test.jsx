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
    streetAddress1: 'Street Address',
    city: 'Fort Irwin',
    state: 'CA',
    postalCode: '92310',
  },
  handleShowDiversionModal: jest.fn(),
  shipmentInfo: {
    id: '456',
    eTag: 'abc123',
    status: 'APPROVED',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    shipmentLocator: 'ABCDEF-01',
  },
  diversionReason: '',
};

const ppmShipment = {
  pickupAddress: {
    city: 'Tampa',
    state: 'FL',
    postalCode: '33621',
    streetAddress1: '123 Fake Street',
    streetAddress2: '',
    streetAddress3: '',
    country: 'USA',
  },
  destinationAddress: {
    city: 'Chicago',
    state: 'IL',
    postalCode: '01054',
    streetAddress1: '5 Main Street',
    streetAddress2: '',
    streetAddress3: '',
    country: 'USA',
  },
  shipmentInfo: {
    id: '1234',
    eTag: 'abc123',
    status: 'APPROVED',
    shipmentType: SHIPMENT_OPTIONS.PPM,
  },
};

const canceledShipment = {
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
  handleShowDiversionModal: jest.fn(),
  shipmentInfo: {
    id: '456',
    eTag: 'abc123',
    status: 'CANCELED',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    shipmentLocator: 'ABCDEF-01',
  },
};

const cancellationRequestedShipment = {
  ...canceledShipment,
  shipmentInfo: {
    id: '456',
    eTag: 'abc123',
    status: 'CANCELLATION_REQUESTED',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    shipmentLocator: 'ABCDEF-01',
  },
};

describe('ShipmentAddresses', () => {
  it('calls props.handleShowDiversionModal on request diversion button click', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createShipmentDiversionRequest, permissionTypes.updateMTOPage]}>
        <ShipmentAddresses {...testProps} />
      </MockProviders>,
    );
    const requestDiversionBtn = screen.getByRole('button', { name: 'Request Diversion' });

    await userEvent.click(requestDiversionBtn);
    await waitFor(() => {
      expect(testProps.handleShowDiversionModal).toHaveBeenCalled();
      expect(testProps.handleShowDiversionModal).toHaveBeenCalledWith(testProps.shipmentInfo);
    });
  });

  it('hides the request diversion button for a cancelled shipment', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createShipmentDiversionRequest, permissionTypes.updateMTOPage]}>
        <ShipmentAddresses {...canceledShipment} />
      </MockProviders>,
    );
    const requestDiversionBtn = screen.queryByRole('button', { name: 'Request Diversion' });

    await waitFor(() => {
      expect(requestDiversionBtn).toBeNull();
    });
  });

  it('hides the request diversion button for a cancelation requested shipment', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createShipmentDiversionRequest, permissionTypes.updateMTOPage]}>
        <ShipmentAddresses {...cancellationRequestedShipment} />
      </MockProviders>,
    );
    const requestDiversionBtn = screen.queryByRole('button', { name: 'Request Diversion' });

    await waitFor(() => {
      expect(requestDiversionBtn).toBeNull();
    });
  });

  it('hides the request diversion button when user does not have permissions', async () => {
    render(<ShipmentAddresses {...canceledShipment} />);
    const requestDiversionBtn = screen.queryByRole('button', { name: 'Request Diversion' });

    await waitFor(() => {
      expect(requestDiversionBtn).toBeNull();
    });
  });

  it('hides the request diversion button when user does not have updateMTOPage permissions', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createShipmentDiversionRequest]}>
        <ShipmentAddresses {...canceledShipment} />
      </MockProviders>,
    );
    const requestDiversionBtn = screen.queryByRole('button', { name: 'Request Diversion' });

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
    expect(screen.getByText('Pickup Address')).toBeInTheDocument();
    expect(screen.getByText('Facility address')).toBeInTheDocument();
  });

  it('shows correct headings for NTSR', () => {
    const NTSRProps = {
      ...testProps,
      shipmentInfo: { ...testProps.shipmentInfo, shipmentType: SHIPMENT_OPTIONS.NTSR },
    };
    render(<ShipmentAddresses {...NTSRProps} />);
    expect(screen.getByText('Facility address')).toBeInTheDocument();
    expect(screen.getByText('Delivery Address')).toBeInTheDocument();
  });

  it('shows correct headings for PPM', () => {
    render(<ShipmentAddresses {...ppmShipment} />);

    expect(screen.getByText("Customer's addresses")).toBeInTheDocument();
    expect(screen.getByText('Authorized addresses')).toBeInTheDocument();
  });

  it('does not show request diversion for PPM', () => {
    render(
      <MockProviders permissions={[permissionTypes.createShipmentDiversionRequest, permissionTypes.updateMTOPage]}>
        <ShipmentAddresses {...ppmShipment} />
      </MockProviders>,
    );

    expect(screen.queryByText('Request Diversion')).not.toBeInTheDocument();
  });

  it('renders with disabled request diversion button', async () => {
    const isMoveLocked = true;
    render(
      <MockProviders permissions={[permissionTypes.createShipmentDiversionRequest, permissionTypes.updateMTOPage]}>
        <ShipmentAddresses {...testProps} isMoveLocked={isMoveLocked} />
      </MockProviders>,
    );
    const requestDiversionBtn = screen.getByRole('button', { name: 'Request Diversion' });
    expect(requestDiversionBtn).toBeDisabled();
  });
});

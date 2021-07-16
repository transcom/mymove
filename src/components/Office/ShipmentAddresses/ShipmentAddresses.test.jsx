import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentAddresses from './ShipmentAddresses';

const testProps = {
  pickupAddress: {
    city: 'Fairfax',
    state: 'VA',
    postal_code: '12345',
    street_address_1: '123 Fake Street',
    street_address_2: '',
    street_address_3: '',
    country: 'USA',
  },
  destinationAddress: {
    city: 'Boston',
    state: 'MA',
    postal_code: '01054',
    street_address_1: '5 Main Street',
    street_address_2: '',
    street_address_3: '',
    country: 'USA',
  },
  destinationDutyStation: {
    street_address_1: '',
    city: 'Fort Irwin',
    state: 'CA',
    postal_code: '92310',
  },
  handleDivertShipment: jest.fn(),
  shipmentInfo: {
    shipmentID: '456',
    ifMatchEtag: 'abc123',
    shipmentStatus: 'APPROVED',
  },
};

const cancelledShipment = {
  pickupAddress: {
    city: 'Fairfax',
    state: 'VA',
    postal_code: '12345',
    street_address_1: '123 Fake Street',
    street_address_2: '',
    street_address_3: '',
    country: 'USA',
  },
  destinationAddress: {
    city: 'Boston',
    state: 'MA',
    postal_code: '01054',
    street_address_1: '5 Main Street',
    street_address_2: '',
    street_address_3: '',
    country: 'USA',
  },
  destinationDutyStation: {
    street_address_1: '',
    city: 'Fort Irwin',
    state: 'CA',
    postal_code: '92310',
  },
  handleDivertShipment: jest.fn(),
  shipmentInfo: {
    shipmentID: '456',
    ifMatchEtag: 'abc123',
    shipmentStatus: 'CANCELED',
  },
};

describe('ShipmentAddresses', () => {
  it('calls props.handleDivertShipment on request diversion button click', async () => {
    render(<ShipmentAddresses {...testProps} />);
    const requestDiversionBtn = screen.getByRole('button', { name: 'Request diversion' });

    userEvent.click(requestDiversionBtn);
    await waitFor(() => {
      expect(testProps.handleDivertShipment).toHaveBeenCalled();
      expect(testProps.handleDivertShipment).toHaveBeenCalledWith(
        testProps.shipmentInfo.shipmentID,
        testProps.shipmentInfo.ifMatchEtag,
      );
    });
  });

  it('hides the request diversion button for a cancelled shipment', async () => {
    render(<ShipmentAddresses {...cancelledShipment} />);
    const requestDiversionBtn = screen.queryByRole('button', { name: 'Request diversion' });

    await waitFor(() => {
      expect(requestDiversionBtn).toBeNull();
    });
  });
});

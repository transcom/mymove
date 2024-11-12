import React from 'react';
import { render, screen, within, act } from '@testing-library/react';

import BoatShipmentInfoList from './BoatShipmentInfoList';

import { boatShipmentTypes } from 'constants/shipments';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const shipment = {
  boatShipment: {
    type: boatShipmentTypes.TOW_AWAY,
    year: 2020,
    make: 'Yamaha',
    model: '242X E-Series',
    lengthInInches: 276,
    widthInInches: 102,
    heightInInches: 120,
    hasTrailer: true,
    isRoadworthy: true,
  },
  requestedPickupDate: '2020-03-26',
  pickupAddress: {
    streetAddress1: '123 Harbor Dr',
    city: 'Miami',
    state: 'FL',
    postalCode: '33101',
  },
  destinationAddress: {
    streetAddress1: '456 Marina Blvd',
    city: 'Key West',
    state: 'FL',
    postalCode: '33040',
  },
  mtoAgents: [
    {
      agentType: 'RELEASING_AGENT',
      firstName: 'John',
      lastName: 'Doe',
      phone: '123-456-7890',
      email: 'john.doe@example.com',
    },
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Jane',
      lastName: 'Smith',
      phone: '987-654-3210',
      email: 'jane.smith@example.com',
    },
  ],
  counselorRemarks: 'Handle with care.',
  customerRemarks: 'Please avoid scratches.',
};

const labels = {
  requestedPickupDate: 'Requested pickup date',
  pickupAddress: 'Pickup Address',
  destinationAddress: 'Delivery Address',
  mtoAgents: ['Releasing agent', 'Receiving agent'],
  counselorRemarks: 'Counselor remarks',
  customerRemarks: 'Customer remarks',
  dimensions: 'Dimensions',
  hasTrailer: 'Trailer',
  isRoadworthy: 'Is trailer roadworthy',
};

describe('Shipment Info List - Boat Shipment', () => {
  it('renders all boat shipment fields when provided and expanded', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    await act(async () => {
      render(<BoatShipmentInfoList isExpanded shipment={shipment} />);
    });

    const requestedPickupDate = screen.getByText(labels.requestedPickupDate);
    expect(within(requestedPickupDate.parentElement).getByText('26 Mar 2020')).toBeInTheDocument();

    const pickupAddress = screen.getByText(labels.pickupAddress);
    expect(
      within(pickupAddress.parentElement).getByText(shipment.pickupAddress.streetAddress1, { exact: false }),
    ).toBeInTheDocument();

    const destinationAddress = screen.getByText(labels.destinationAddress);
    expect(
      within(destinationAddress.parentElement).getByText(shipment.destinationAddress.streetAddress1, {
        exact: false,
      }),
    ).toBeInTheDocument();

    const releasingAgent = screen.getByText(labels.mtoAgents[0]);
    expect(
      within(releasingAgent.parentElement).getByText(shipment.mtoAgents[0].email, { exact: false }),
    ).toBeInTheDocument();

    const receivingAgent = screen.getByText(labels.mtoAgents[1]);
    expect(
      within(receivingAgent.parentElement).getByText(shipment.mtoAgents[1].email, { exact: false }),
    ).toBeInTheDocument();

    const counselorRemarks = screen.getByText(labels.counselorRemarks);
    expect(within(counselorRemarks.parentElement).getByText(shipment.counselorRemarks)).toBeInTheDocument();

    const customerRemarks = screen.getByText(labels.customerRemarks);
    expect(within(customerRemarks.parentElement).getByText(shipment.customerRemarks)).toBeInTheDocument();

    const dimensions = screen.getByText(labels.dimensions);
    expect(
      within(dimensions.parentElement).getByText("23' L x 8' 6\" W x 10' H", { exact: false }),
    ).toBeInTheDocument();

    const hasTrailer = screen.getByText(labels.hasTrailer);
    expect(within(hasTrailer.parentElement).getByText('Yes')).toBeInTheDocument();

    const isRoadworthy = screen.getByText(labels.isRoadworthy);
    expect(within(isRoadworthy.parentElement).getByText('Yes')).toBeInTheDocument();
  });

  it('does not render mtoAgents when not provided', async () => {
    await act(async () => {
      render(
        <BoatShipmentInfoList
          shipment={{
            ...shipment,
            mtoAgents: [],
          }}
        />,
      );
    });

    expect(screen.queryByText(labels.mtoAgents[0])).not.toBeInTheDocument();
    expect(screen.queryByText(labels.mtoAgents[1])).not.toBeInTheDocument();
  });
});

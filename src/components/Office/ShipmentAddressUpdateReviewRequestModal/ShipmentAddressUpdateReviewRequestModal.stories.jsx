import React from 'react';

import { ShipmentAddressUpdateReviewRequestModal } from './ShipmentAddressUpdateReviewRequestModal';

import { shipmentTypes } from 'constants/shipments';

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
};

export default {
  title: 'Office Components/ShipmentAddressUpdateReviewRequestModal',
  component: ShipmentAddressUpdateReviewRequestModal,
};

export const ReviewModal = {
  render: () => <ShipmentAddressUpdateReviewRequestModal shipment={mockShipment} onClose={() => {}} />,
};

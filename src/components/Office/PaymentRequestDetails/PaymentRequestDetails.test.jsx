import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import PaymentRequestDetails from './PaymentRequestDetails';

import { PAYMENT_SERVICE_ITEM_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';
import { shipmentModificationTypes } from 'constants/shipments';
import { MockProviders } from 'testUtils';
import PAYMENT_REQUEST_STATUSES from 'constants/paymentRequestStatus';

const basicServiceItems = [
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:00:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: null,
    mtoServiceItemName: 'Move management',
  },
  {
    id: '39474c6a-69b6-4501-8e08-670a12512a5e',
    createdAt: '2020-12-01T00:00:00.000Z',
    mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    priceCents: 4000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
    rejectionReason: 'duplicate charge',
    mtoShipmentType: null,
    mtoServiceItemName: 'Counseling',
  },
];

const oneBasicServiceItem = [
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5e',
    createdAt: '2020-12-01T00:00:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    priceCents: 2000001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoServiceItemName: 'Move management',
  },
];

const hhgServiceItems = [
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5a',
    createdAt: '2020-12-01T00:04:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 100001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic linehaul',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5b',
    createdAt: '2020-12-01T00:05:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbb',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 200001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Fuel surcharge',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5c',
    createdAt: '2020-12-01T00:06:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 300001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic origin price',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5d',
    createdAt: '2020-12-01T00:07:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbd',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 400001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic destination price',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5e',
    createdAt: '2020-12-01T00:08:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbe',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 500001,
    status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic packing',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:09:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbf',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 600001,
    status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
    mtoShipmentType: SHIPMENT_OPTIONS.HHG,
    mtoServiceItemName: 'Domestic unpacking',
  },
];

const ntsrServiceItems = [
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5a',
    createdAt: '2020-12-01T00:04:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 100001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Domestic linehaul',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5b',
    createdAt: '2020-12-01T00:05:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbb',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 200001,
    status: PAYMENT_SERVICE_ITEM_STATUS.REQUESTED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Fuel surcharge',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5c',
    createdAt: '2020-12-01T00:06:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 300001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Domestic origin price',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5d',
    createdAt: '2020-12-01T00:07:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbd',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 400001,
    status: PAYMENT_SERVICE_ITEM_STATUS.APPROVED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Domestic destination price',
  },
  {
    id: '09474c6a-69b6-4501-8e08-670a12512a5f',
    createdAt: '2020-12-01T00:09:00.000Z',
    mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbf',
    mtoShipmentID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dba',
    priceCents: 600001,
    status: PAYMENT_SERVICE_ITEM_STATUS.DENIED,
    mtoShipmentType: SHIPMENT_OPTIONS.NTSR,
    mtoServiceItemName: 'Domestic unpacking',
  },
];

const hhgShipment = {
  address: 'Beverly Hills, CA 90210 to Fairfield, CA 94535',
  departureDate: '2020-12-01T00:00:00.000Z',
  tacType: 'HHG',
  sacType: 'HHG',
};

const hhgShipmentCanceled = {
  address: 'Beverly Hills, CA 90210 to Fairfield, CA 94535',
  departureDate: '2020-12-01T00:00:00.000Z',
  modificationType: shipmentModificationTypes.CANCELED,
};

const hhgShipmentDiversion = {
  address: 'Beverly Hills, CA 90210 to Fairfield, CA 94535',
  departureDate: '2020-12-01T00:00:00.000Z',
  modificationType: shipmentModificationTypes.DIVERSION,
};

const basicShipment = {
  address: '',
  departureDate: '',
};

const ntsShipment = {
  address: 'Boston, MA 02101 to Princeton, NJ 08540',
  departureDate: '020-12-01T00:00:00.000Z',
  tacType: 'NTS',
  sacType: 'HHG',
};

const testMoveLocator = 'AF7K1P';

describe('PaymentRequestDetails', () => {
  describe('When given basic service items', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails
          serviceItems={basicServiceItems}
          shipment={basicShipment}
          paymentRequestStatus={PAYMENT_REQUEST_STATUSES.REVIEWED}
        />
      </MockProviders>,
    );

    it('renders the service items', async () => {
      expect(wrapper.find('td')).toBeTruthy();
    });

    it('renders the expected table title', () => {
      expect(wrapper.text().includes('Basic service items (2 items)')).toBeTruthy();
    });

    it('renders the service item names', () => {
      const serviceItemNames = wrapper.find({ 'data-testid': 'serviceItemName' });
      expect(serviceItemNames.at(0).text()).toEqual('Move management');
      expect(serviceItemNames.at(1).text()).toEqual('Counseling');
    });

    it('renders the service item amounts', () => {
      const serviceItemAmounts = wrapper.find({ 'data-testid': 'serviceItemAmount' });
      expect(serviceItemAmounts.at(0).text()).toEqual('$20,000.01');
      expect(serviceItemAmounts.at(1).text()).toEqual('$40,000.01');
    });

    it('renders the service item statuses', () => {
      const serviceItemStatuses = wrapper.find({ 'data-testid': 'serviceItemStatus' });
      expect(serviceItemStatuses.at(0).text().includes('Accepted')).toBeTruthy();
      expect(serviceItemStatuses.at(1).text().includes('Rejected')).toBeTruthy();
    });

    it('does not render the Departure Date, Pickup Address, and Destination Address', async () => {
      expect(wrapper.find({ 'data-testid': 'pickup-to-destination' }).length).toBe(0);
      expect(wrapper.find({ 'data-testid': 'departure-date' }).length).toBe(0);
    });
  });

  describe('When given a single basic service item', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails
          serviceItems={oneBasicServiceItem}
          paymentRequestStatus={PAYMENT_REQUEST_STATUSES.PENDING}
        />
      </MockProviders>,
    );

    it('renders the expected table title', () => {
      expect(wrapper.text().includes('Basic service items (1 item)')).toBeTruthy();
    });

    it('does not render the Departure Date, Pickup Address, and Destination Address', async () => {
      expect(wrapper.find({ 'data-testid': 'pickup-to-destination' }).length).toBe(0);
      expect(wrapper.find({ 'data-testid': 'departure-date' }).length).toBe(0);
    });
  });

  describe('When given a hhg shipment service items', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails
          serviceItems={hhgServiceItems}
          shipment={hhgShipment}
          paymentRequestStatus={PAYMENT_REQUEST_STATUSES.PENDING}
        />
      </MockProviders>,
    );

    it('renders the expected table title', () => {
      expect(wrapper.text().includes('HHG (6 items)')).toBeTruthy();
    });

    it('does renders the Departure Date, Pickup Address, and Destination Address', async () => {
      expect(wrapper.find({ 'data-testid': 'pickup-to-destination' })).toBeTruthy();
      expect(
        wrapper.find({ 'data-testid': 'pickup-to-destination' }).at(0).text().includes('Fairfield, CA 94535'),
      ).toBeTruthy();
      expect(wrapper.find({ 'data-testid': 'departure-date' }).text().includes('Departed')).toBeTruthy();
    });

    it('renders the service item names', () => {
      const serviceItemNames = wrapper.find({ 'data-testid': 'serviceItemName' });
      expect(serviceItemNames.at(0).text()).toEqual('Domestic linehaul');
      expect(serviceItemNames.at(1).text()).toEqual('Fuel surcharge');
      expect(serviceItemNames.at(2).text()).toEqual('Domestic origin price');
      expect(serviceItemNames.at(3).text()).toEqual('Domestic destination price');
      expect(serviceItemNames.at(4).text()).toEqual('Domestic packing');
      expect(serviceItemNames.at(5).text()).toEqual('Domestic unpacking');
    });

    it('renders the service item amounts', () => {
      const serviceItemAmounts = wrapper.find({ 'data-testid': 'serviceItemAmount' });
      expect(serviceItemAmounts.at(0).text()).toEqual('$1,000.01');
      expect(serviceItemAmounts.at(1).text()).toEqual('$2,000.01');
      expect(serviceItemAmounts.at(2).text()).toEqual('$3,000.01');
      expect(serviceItemAmounts.at(3).text()).toEqual('$4,000.01');
      expect(serviceItemAmounts.at(4).text()).toEqual('$5,000.01');
      expect(serviceItemAmounts.at(5).text()).toEqual('$6,000.01');
    });

    it('renders the service item statuses', () => {
      const serviceItemStatuses = wrapper.find({ 'data-testid': 'serviceItemStatus' });
      expect(serviceItemStatuses.at(0).text().includes('Needs review')).toBeTruthy();
      expect(serviceItemStatuses.at(1).text().includes('Needs review')).toBeTruthy();
      expect(serviceItemStatuses.at(2).text().includes('Accepted')).toBeTruthy();
      expect(serviceItemStatuses.at(3).text().includes('Accepted')).toBeTruthy();
      expect(serviceItemStatuses.at(4).text().includes('Rejected')).toBeTruthy();
      expect(serviceItemStatuses.at(5).text().includes('Rejected')).toBeTruthy();
    });

    it('renders the TAC/SAC codes', () => {
      render(
        <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
          <PaymentRequestDetails
            serviceItems={hhgServiceItems}
            shipment={hhgShipment}
            paymentRequestStatus={PAYMENT_REQUEST_STATUSES.PENDING}
            tacs={{ HHG: '1234' }}
            sacs={{ HHG: 'AB12' }}
          />
        </MockProviders>,
      );

      expect(screen.getByText(/1234 \(HHG\)/)).toBeInTheDocument();
      expect(screen.getByText(/AB12 \(HHG\)/)).toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
    });
  });

  describe('When given a ntsr shipment service items', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails
          serviceItems={ntsrServiceItems}
          shipment={ntsShipment}
          paymentRequestStatus={PAYMENT_REQUEST_STATUSES.PENDING}
        />
      </MockProviders>,
    );

    it('renders the expected table title', () => {
      expect(wrapper.text().includes('Non-temp storage release (5 items)')).toBeTruthy();
    });

    it('does renders the Departure Date, Pickup Address, and Destination Address', async () => {
      expect(wrapper.find({ 'data-testid': 'pickup-to-destination' })).toBeTruthy();
      expect(
        wrapper.find({ 'data-testid': 'pickup-to-destination' }).at(0).text().includes('Princeton, NJ 08540'),
      ).toBeTruthy();
      expect(wrapper.find({ 'data-testid': 'departure-date' }).text().includes('Departed')).toBeTruthy();
    });

    it('renders the service item names', () => {
      const serviceItemNames = wrapper.find({ 'data-testid': 'serviceItemName' });
      expect(serviceItemNames.at(0).text()).toEqual('Domestic linehaul');
      expect(serviceItemNames.at(1).text()).toEqual('Fuel surcharge');
      expect(serviceItemNames.at(2).text()).toEqual('Domestic origin price');
      expect(serviceItemNames.at(3).text()).toEqual('Domestic destination price');
      expect(serviceItemNames.at(4).text()).toEqual('Domestic unpacking');
    });

    it('renders the service item amounts', () => {
      const serviceItemAmounts = wrapper.find({ 'data-testid': 'serviceItemAmount' });
      expect(serviceItemAmounts.at(0).text()).toEqual('$1,000.01');
      expect(serviceItemAmounts.at(1).text()).toEqual('$2,000.01');
      expect(serviceItemAmounts.at(2).text()).toEqual('$3,000.01');
      expect(serviceItemAmounts.at(3).text()).toEqual('$4,000.01');
      expect(serviceItemAmounts.at(4).text()).toEqual('$6,000.01');
    });

    it('renders the service item statuses', () => {
      const serviceItemStatuses = wrapper.find({ 'data-testid': 'serviceItemStatus' });
      expect(serviceItemStatuses.at(0).text().includes('Needs review')).toBeTruthy();
      expect(serviceItemStatuses.at(1).text().includes('Needs review')).toBeTruthy();
      expect(serviceItemStatuses.at(2).text().includes('Accepted')).toBeTruthy();
      expect(serviceItemStatuses.at(3).text().includes('Accepted')).toBeTruthy();
      expect(serviceItemStatuses.at(4).text().includes('Rejected')).toBeTruthy();
    });

    it('renders the TAC/SAC codes', () => {
      render(
        <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
          <PaymentRequestDetails
            serviceItems={ntsrServiceItems}
            shipment={ntsShipment}
            paymentRequestStatus={PAYMENT_REQUEST_STATUSES.PENDING}
            tacs={{ HHG: '1234', NTS: '5678' }}
            sacs={{ HHG: 'AB12', NTS: 'CD34' }}
          />
        </MockProviders>,
      );

      expect(screen.getByText(/5678 \(NTS\)/)).toBeInTheDocument();
      expect(screen.getByText(/AB12 \(HHG\)/)).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
    });
  });

  describe('When a payment request is in the pending status', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails
          serviceItems={hhgServiceItems}
          shipment={hhgShipment}
          paymentRequestStatus={PAYMENT_REQUEST_STATUSES.PENDING}
        />
      </MockProviders>,
    );
    it('disables expanding the service item pricer calculations', () => {
      expect(wrapper.find('ExpandableServiceItemRow').at(0).prop('disableExpansion')).toBe(true);
    });
  });

  describe('When a payment request is in a reviewed status', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails
          serviceItems={hhgServiceItems}
          shipment={hhgShipment}
          paymentRequestStatus={PAYMENT_REQUEST_STATUSES.REVIEWED}
        />
      </MockProviders>,
    );
    it('disables expanding the service item pricer calculations', () => {
      expect(wrapper.find('ExpandableServiceItemRow').at(0).prop('disableExpansion')).toBe(false);
    });
  });

  describe('When a payment request has a shipment that was canceled ', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails
          serviceItems={hhgServiceItems}
          shipment={hhgShipmentCanceled}
          paymentRequestStatus={PAYMENT_REQUEST_STATUSES.PENDING}
        />
      </MockProviders>,
    );
    it('there is a canceled tag displayed', () => {
      expect(wrapper.find('ShipmentModificationTag').text()).toBe(shipmentModificationTypes.CANCELED);
    });
  });

  describe('When a payment request has a shipment that was diverted ', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails
          serviceItems={hhgServiceItems}
          shipment={hhgShipmentDiversion}
          paymentRequestStatus={PAYMENT_REQUEST_STATUSES.PENDING}
        />
      </MockProviders>,
    );
    it('there is a diversion tag displayed', () => {
      expect(wrapper.find('ShipmentModificationTag').text()).toBe(shipmentModificationTypes.DIVERSION);
    });
  });
});

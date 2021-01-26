import React from 'react';
import { mount } from 'enzyme';

import PaymentRequestDetails from './PaymentRequestDetails';

import { PAYMENT_SERVICE_ITEM_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';

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

const shipmentAddressBasic = '';
const shipmentAddressHHG = 'Beverly Hills, CA 90210 to Fairfield, CA 94535';
const shipmentAddressNTS = 'Boston, MA 02101 to Princeton, NJ 08540';

const testMoveLocator = 'AF7K1P';

describe('PaymentRequestDetails', () => {
  describe('When given basic service items', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails serviceItems={basicServiceItems} shipmentAddress={shipmentAddressBasic} />
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

    it('does not render the Pickup Address and Destination Address', async () => {
      expect(wrapper.find({ 'data-testid': 'pickup-to-destination' }).length).toBe(0);
    });
  });

  describe('When given a single basic service item', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails serviceItems={oneBasicServiceItem} />
      </MockProviders>,
    );

    it('renders the expected table title', () => {
      expect(wrapper.text().includes('Basic service items (1 item)')).toBeTruthy();
    });

    it('does not render the Pickup Address and Destination Address', async () => {
      expect(wrapper.find({ 'data-testid': 'pickup-to-destination' }).length).toBe(0);
    });
  });

  describe('When given a hhg shipment service items', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails serviceItems={hhgServiceItems} shipmentAddress={shipmentAddressHHG} />
      </MockProviders>,
    );

    it('renders the expected table title', () => {
      expect(wrapper.text().includes('Household goods (6 items)')).toBeTruthy();
    });

    it('does renders the Pickup Address and Destination Address', async () => {
      expect(wrapper.find({ 'data-testid': 'pickup-to-destination' })).toBeTruthy();
      expect(
        wrapper.find({ 'data-testid': 'pickup-to-destination' }).at(0).text().includes('Fairfield, CA 94535'),
      ).toBeTruthy();
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
  });

  describe('When given a ntsr shipment service items', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${testMoveLocator}/payment-requests`]}>
        <PaymentRequestDetails serviceItems={ntsrServiceItems} shipmentAddress={shipmentAddressNTS} />
      </MockProviders>,
    );

    it('renders the expected table title', () => {
      expect(wrapper.text().includes('NTS release (5 items)')).toBeTruthy();
    });

    it('does renders the Pickup Address and Destination Address', async () => {
      expect(wrapper.find({ 'data-testid': 'pickup-to-destination' })).toBeTruthy();
      expect(
        wrapper.find({ 'data-testid': 'pickup-to-destination' }).at(0).text().includes('Princeton, NJ 08540'),
      ).toBeTruthy();
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
  });
});

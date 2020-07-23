import sortServiceItemsByGroup from './serviceItems';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('serviceItems utils', () => {
  describe('sortServiceItemsByGroup', () => {
    describe('when there are service items without a shipment', () => {
      it('sorts basic service items together', () => {
        const serviceItemCards = [
          {
            id: '1',
            serviceItemName: 'Counseling services',
            amount: 0.01,
            createdAt: '2020-01-01T00:09:00.999Z',
          },
          {
            id: '2',
            serviceItemName: 'Move management',
            amount: 1234.0,
            createdAt: '2020-01-01T00:06:00.999Z',
          },
          {
            id: '3',
            shipmentId: '20',
            shipmentType: SHIPMENT_OPTIONS.HHG,
            serviceItemName: 'Domestic linehaul',
            amount: 5678.05,
            createdAt: '2020-01-01T00:08:00.999Z',
          },
        ];
        expect(sortServiceItemsByGroup(serviceItemCards)).toEqual([
          {
            id: '2',
            serviceItemName: 'Move management',
            amount: 1234.0,
            createdAt: '2020-01-01T00:06:00.999Z',
          },
          {
            id: '1',
            serviceItemName: 'Counseling services',
            amount: 0.01,
            createdAt: '2020-01-01T00:09:00.999Z',
          },
          {
            id: '3',
            shipmentId: '20',
            shipmentType: SHIPMENT_OPTIONS.HHG,
            serviceItemName: 'Domestic linehaul',
            amount: 5678.05,
            createdAt: '2020-01-01T00:08:00.999Z',
          },
        ]);
      });
      describe('when there are multiple service items per shipment', () => {
        it('sorts basic service items together', () => {
          const serviceItemCards = [
            {
              id: '1',
              serviceItemName: 'Counseling services',
              amount: 0.01,
              createdAt: '2020-01-01T00:09:00.999Z',
            },
            {
              id: '2',
              serviceItemName: 'Move management',
              amount: 1234.0,
              createdAt: '2020-01-01T00:06:00.999Z',
            },
            {
              id: '3',
              shipmentId: '20',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Domestic linehaul',
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:10.999Z',
            },
            {
              id: '4',
              shipmentId: '20',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Fuel Surcharge',
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:00.999Z',
            },
          ];
          expect(sortServiceItemsByGroup(serviceItemCards)).toEqual([
            {
              id: '2',
              serviceItemName: 'Move management',
              amount: 1234.0,
              createdAt: '2020-01-01T00:06:00.999Z',
            },
            {
              id: '1',
              serviceItemName: 'Counseling services',
              amount: 0.01,
              createdAt: '2020-01-01T00:09:00.999Z',
            },
            {
              id: '4',
              shipmentId: '20',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Fuel Surcharge',
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:00.999Z',
            },
            {
              id: '3',
              shipmentId: '20',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Domestic linehaul',
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:10.999Z',
            },
          ]);
        });
      });
      describe('when there are multiple shipments of the same type', () => {
        it('sorts basic service items together', () => {
          const serviceItemCards = [
            {
              id: '1',
              serviceItemName: 'Counseling services',
              amount: 0.01,
              createdAt: '2020-01-01T00:09:00.999Z',
            },
            {
              id: '2',
              serviceItemName: 'Move management',
              amount: 1234.0,
              createdAt: '2020-01-01T00:06:00.999Z',
            },
            {
              id: '3',
              shipmentId: '20',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Domestic linehaul',
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:10.999Z',
            },
            {
              id: '4',
              shipmentId: '20',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Fuel Surcharge',
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:00.999Z',
            },
            {
              id: '5',
              shipmentId: '30',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Domestic linehaul',
              amount: 5678.05,
              createdAt: '2020-01-01T00:03:10.999Z',
            },
            {
              id: '6',
              shipmentId: '30',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Fuel Surcharge',
              amount: 5678.05,
              createdAt: '2020-01-01T00:03:00.999Z',
            },
          ];
          expect(sortServiceItemsByGroup(serviceItemCards)).toEqual([
            {
              id: '6',
              shipmentId: '30',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Fuel Surcharge',
              amount: 5678.05,
              createdAt: '2020-01-01T00:03:00.999Z',
            },
            {
              id: '5',
              shipmentId: '30',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Domestic linehaul',
              amount: 5678.05,
              createdAt: '2020-01-01T00:03:10.999Z',
            },
            {
              id: '2',
              serviceItemName: 'Move management',
              amount: 1234.0,
              createdAt: '2020-01-01T00:06:00.999Z',
            },
            {
              id: '1',
              serviceItemName: 'Counseling services',
              amount: 0.01,
              createdAt: '2020-01-01T00:09:00.999Z',
            },
            {
              id: '4',
              shipmentId: '20',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Fuel Surcharge',
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:00.999Z',
            },
            {
              id: '3',
              shipmentId: '20',
              shipmentType: SHIPMENT_OPTIONS.HHG,
              serviceItemName: 'Domestic linehaul',
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:10.999Z',
            },
          ]);
        });
      });
    });
  });
});

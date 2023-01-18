import { sortServiceItemsByGroup, formatDimensions, hasCounseling, hasMoveManagement } from './serviceItems';

import { formatToThousandthInches } from 'utils/formatters';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { serviceItemCodes } from 'content/serviceItems';

describe('serviceItems utils', () => {
  describe('formatDimensions', () => {
    describe('default conversion from thousands of inch', () => {
      it('converts to inches and adds inches symbol', () => {
        expect(formatDimensions({ length: 10000, width: 2500, height: 50000 })).toBe('10"x2.5"x50"');
      });
    });
    describe('conversion from inches to thousands of an inch', () => {
      it('converts to inches and adds inches symbol', () => {
        expect(formatDimensions({ length: 10, width: 2.5, height: 50 }, formatToThousandthInches, '')).toBe(
          '10000x2500x50000',
        );
      });
    });
  });
  describe('sortServiceItemsByGroup', () => {
    describe('when there are service items without a shipment', () => {
      it('sorts basic service items together', () => {
        const serviceItemCards = [
          {
            id: '1',
            mtoServiceItemName: serviceItemCodes.CS,
            amount: 0.01,
            createdAt: '2020-01-01T00:09:00.999Z',
          },
          {
            id: '2',
            mtoServiceItemName: serviceItemCodes.MS,
            amount: 1234.0,
            createdAt: '2020-01-01T00:06:00.999Z',
          },
          {
            id: '3',
            mtoShipmentID: '20',
            mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
            mtoServiceItemName: serviceItemCodes.DLH,
            amount: 5678.05,
            createdAt: '2020-01-01T00:08:00.999Z',
          },
        ];
        expect(sortServiceItemsByGroup(serviceItemCards)).toEqual([
          {
            id: '2',
            mtoServiceItemName: serviceItemCodes.MS,
            amount: 1234.0,
            createdAt: '2020-01-01T00:06:00.999Z',
          },
          {
            id: '1',
            mtoServiceItemName: serviceItemCodes.CS,
            amount: 0.01,
            createdAt: '2020-01-01T00:09:00.999Z',
          },
          {
            id: '3',
            mtoShipmentID: '20',
            mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
            mtoServiceItemName: serviceItemCodes.DLH,
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
              mtoServiceItemName: serviceItemCodes.CS,
              amount: 0.01,
              createdAt: '2020-01-01T00:09:00.999Z',
            },
            {
              id: '2',
              mtoServiceItemName: serviceItemCodes.MS,
              amount: 1234.0,
              createdAt: '2020-01-01T00:06:00.999Z',
            },
            {
              id: '3',
              mtoShipmentID: '20',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.DLH,
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:10.999Z',
            },
            {
              id: '4',
              mtoShipmentID: '20',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.FSC,
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:00.999Z',
            },
          ];
          expect(sortServiceItemsByGroup(serviceItemCards)).toEqual([
            {
              id: '2',
              mtoServiceItemName: serviceItemCodes.MS,
              amount: 1234.0,
              createdAt: '2020-01-01T00:06:00.999Z',
            },
            {
              id: '1',
              mtoServiceItemName: serviceItemCodes.CS,
              amount: 0.01,
              createdAt: '2020-01-01T00:09:00.999Z',
            },
            {
              id: '4',
              mtoShipmentID: '20',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.FSC,
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:00.999Z',
            },
            {
              id: '3',
              mtoShipmentID: '20',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.DLH,
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
              mtoServiceItemName: serviceItemCodes.CS,
              amount: 0.01,
              createdAt: '2020-01-01T00:09:00.999Z',
            },
            {
              id: '2',
              mtoServiceItemName: serviceItemCodes.MS,
              amount: 1234.0,
              createdAt: '2020-01-01T00:06:00.999Z',
            },
            {
              id: '3',
              mtoShipmentID: '20',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.DLH,
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:10.999Z',
            },
            {
              id: '4',
              mtoShipmentID: '20',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.FSC,
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:00.999Z',
            },
            {
              id: '5',
              mtoShipmentID: '30',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.DLH,
              amount: 5678.05,
              createdAt: '2020-01-01T00:03:10.999Z',
            },
            {
              id: '6',
              mtoShipmentID: '30',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.FSC,
              amount: 5678.05,
              createdAt: '2020-01-01T00:03:00.999Z',
            },
          ];
          expect(sortServiceItemsByGroup(serviceItemCards)).toEqual([
            {
              id: '6',
              mtoShipmentID: '30',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.FSC,
              amount: 5678.05,
              createdAt: '2020-01-01T00:03:00.999Z',
            },
            {
              id: '5',
              mtoShipmentID: '30',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.DLH,
              amount: 5678.05,
              createdAt: '2020-01-01T00:03:10.999Z',
            },
            {
              id: '2',
              mtoServiceItemName: serviceItemCodes.MS,
              amount: 1234.0,
              createdAt: '2020-01-01T00:06:00.999Z',
            },
            {
              id: '1',
              mtoServiceItemName: serviceItemCodes.CS,
              amount: 0.01,
              createdAt: '2020-01-01T00:09:00.999Z',
            },
            {
              id: '4',
              mtoShipmentID: '20',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.FSC,
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:00.999Z',
            },
            {
              id: '3',
              mtoShipmentID: '20',
              mtoShipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
              mtoServiceItemName: serviceItemCodes.DLH,
              amount: 5678.05,
              createdAt: '2020-01-01T00:08:10.999Z',
            },
          ]);
        });
      });
    });
  });
  describe('service item lookup functions', () => {
    const serviceItems = {
      MS: {
        reServiceName: 'Move management',
        approvedAt: '2020-01-01',
        id: '76055c99-0990-410c-a7c9-69373b0b53eb',
        status: 'APPROVED',
        reServiceCode: 'MS',
      },
      CS: {
        reServiceName: 'Counseling fee',
        id: '76055c99-0990-410c-a7c9-69373b0b5322',
        status: 'APPROVED',
        reServiceCode: 'CS',
        approvedAt: '2020-01-01',
      },
      DLH: {
        reServiceName: 'Domestic linehaul',
        id: '76055c99-0990-410c-a7c9-69373b0b5322',
        status: 'APPROVED',
        reServiceCode: 'DLH',
        approvedAt: '2020-01-01',
      },
    };
    describe('has counseling function', () => {
      it('returns false when given an empty array of service items', () => {
        expect(hasCounseling([])).toBe(false);
      });
      it('returns false when given an array of service items without counseling', () => {
        expect(hasCounseling([serviceItems.MS, serviceItems.DLH])).toBe(false);
      });
      it('returns true when given an array of service items with counseling', () => {
        expect(hasCounseling([serviceItems.MS, serviceItems.DLH, serviceItems.CS])).toBe(true);
      });
      it('returns true when given an array of service items with only counseling', () => {
        expect(hasCounseling([serviceItems.CS])).toBe(true);
      });
    });
    describe('has move management function', () => {
      it('returns false when given an empty array of service items', () => {
        expect(hasMoveManagement([])).toBe(false);
      });
      it('returns false when given an array of service items without counseling', () => {
        expect(hasMoveManagement([serviceItems.CS, serviceItems.DLH])).toBe(false);
      });
      it('returns true when given an array of service items with counseling', () => {
        expect(hasMoveManagement([serviceItems.MS, serviceItems.DLH, serviceItems.CS])).toBe(true);
      });
      it('returns true when given an array of service items with only counseling', () => {
        expect(hasMoveManagement([serviceItems.MS])).toBe(true);
      });
    });
  });
});

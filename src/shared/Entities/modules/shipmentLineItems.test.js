import {
  selectUnbilledShipmentLineItems,
  selectSortedPreApprovalShipmentLineItems,
  selectInvoiceShipmentLineItems,
} from 'shared/Entities/modules/shipmentLineItems';

describe('shipment line items tests', () => {
  describe('When a state with un-billed line items is passed', () => {
    it('returns an empty array of items when no shipment id is found', () => {
      const state = {
        entities: {
          shipmentLineItems: [],
        },
      };
      expect(selectUnbilledShipmentLineItems(state, 'aaabbbccc').length).toEqual(0);
    });
  });
  describe('When a state with un-billed line items is passed', () => {
    let state;
    beforeEach(() => {
      state = {
        entities: {
          tariff400ngItems: {
            'deb28967-d52c-4f04-8a0b-a264c9d80457': {
              code: '105B',
              created_at: '2018-11-05T16:13:22.946Z',
              discount_type: 'HHG',
              id: 'deb28967-d52c-4f04-8a0b-a264c9d80457',
              item: 'Pack Reg Crate',
              location: 'ORIGIN',
              ref_code: 'NONE',
              requires_pre_approval: true,
              uom_1: 'CF',
              uom_2: 'NONE',
              updated_at: '2018-11-05T16:13:22.946Z',
            },
            'c1c2caf0-90e0-445a-9069-328c3201d8b7': {
              code: '105D',
              created_at: '2018-11-05T16:13:22.946Z',
              discount_type: 'HHG',
              id: 'c1c2caf0-90e0-445a-9069-328c3201d8b7',
              item: 'Debris Removal within 30 days',
              location: 'EITHER',
              ref_code: 'NONE',
              requires_pre_approval: true,
              uom_1: 'EA',
              uom_2: 'NONE',
              updated_at: '2018-11-05T16:13:22.946Z',
            },
            'c6a865dd-324a-48a5-9b03-5db8dcd044d1': {
              code: '35A',
              created_at: '2018-11-05T16:13:22.946Z',
              discount_type: 'HHG',
              id: 'c6a865dd-324a-48a5-9b03-5db8dcd044d1',
              item: 'Third Party Service',
              location: 'ORIGIN',
              ref_code: 'NONE',
              requires_pre_approval: true,
              uom_1: 'EA',
              uom_2: 'NONE',
              updated_at: '2018-11-05T16:13:22.946Z',
            },
          },
          shipmentLineItems: {
            'e2a787d5-ff90-4331-8caa-c4c11d5002a1': {
              approved_date: '0001-01-01T00:00:00.000Z',
              created_at: '0001-01-01T00:00:00.000Z',
              id: 'e2a787d5-ff90-4331-8caa-c4c11d5002a1',
              location: 'DESTINATION',
              notes: 'this is a test request',
              quantity_1: 3000,
              quantity_2: 0,
              shipment_id: '4612cfed-acbd-47ca-840a-7b7de190d6d2',
              status: 'APPROVED',
              submitted_date: '2018-11-06T10:22:34.370Z',
              tariff400ng_item: {
                code: '105B',
                created_at: '2018-11-05T16:13:22.946Z',
                discount_type: 'HHG',
                id: 'deb28967-d52c-4f04-8a0b-a264c9d80457',
                item: 'Pack Reg Crate',
                location: 'ORIGIN',
                ref_code: 'NONE',
                requires_pre_approval: true,
                uom_1: 'CF',
                uom_2: 'NONE',
                updated_at: '2018-11-05T16:13:22.946Z',
              },
              tariff400ng_item_id: 'deb28967-d52c-4f04-8a0b-a264c9d80457',
              updated_at: '0001-01-01T00:00:00.000Z',
            },
            'e2a787d5-ff90-4331-8caa-c4c11d5002a2': {
              approved_date: '0001-01-01T00:00:00.000Z',
              created_at: '0001-01-01T00:00:00.000Z',
              id: 'e2a787d5-ff90-4331-8caa-c4c11d5002a2',
              location: 'ORIGIN',
              notes: 'this is a test request',
              quantity_1: 20000,
              quantity_2: 0,
              shipment_id: '4612cfed-acbd-47ca-840a-7b7de190d6723',
              status: 'APPROVED',
              submitted_date: '2018-11-06T10:22:34.370Z',
              tariff400ng_item: {
                code: '105B',
                created_at: '2018-11-05T16:13:22.946Z',
                discount_type: 'HHG',
                id: 'deb28967-d52c-4f04-8a0b-a264c9d80457',
                item: 'Pack Reg Crate',
                location: 'ORIGIN',
                ref_code: 'NONE',
                requires_pre_approval: true,
                uom_1: 'CF',
                uom_2: 'NONE',
                updated_at: '2018-11-05T16:13:22.946Z',
              },
              tariff400ng_item_id: 'deb28967-d52c-4f04-8a0b-a264c9d80457',
              updated_at: '0001-01-01T00:00:00.000Z',
              invoice_id: '683b4683-256c-4108-8400-300f3d87012f',
            },
            'e2a787d5-ff90-4331-8caa-c4c11d5002a3': {
              approved_date: '0001-01-01T00:00:00.000Z',
              created_at: '0001-01-01T00:00:00.000Z',
              id: 'e2a787d5-ff90-4331-8caa-c4c11d5002a3',
              location: 'ORIGIN',
              notes: 'this is a test request',
              quantity_1: 10000,
              quantity_2: 0,
              shipment_id: '4612cfed-acbd-47ca-840a-7b7de190d6723',
              status: 'APPROVED',
              submitted_date: '2018-11-06T10:22:34.370Z',
              tariff400ng_item: {
                code: '105B',
                created_at: '2018-11-05T16:13:22.946Z',
                discount_type: 'HHG',
                id: 'deb28967-d52c-4f04-8a0b-a264c9d80457',
                item: 'Pack Reg Crate',
                location: 'ORIGIN',
                ref_code: 'NONE',
                requires_pre_approval: true,
                uom_1: 'CF',
                uom_2: 'NONE',
                updated_at: '2018-11-05T16:13:22.946Z',
              },
              tariff400ng_item_id: 'deb28967-d52c-4f04-8a0b-a264c9d80457',
              updated_at: '0001-01-01T00:00:00.000Z',
            },
            'e2a787d5-ff90-4331-8caa-c4c11d50ghsdgha4': {
              approved_date: '0001-01-01T00:00:00.000Z',
              created_at: '0001-01-01T00:00:00.000Z',
              id: 'e2a787d5-ff90-4331-8caa-c4c11d50ghsdgha4',
              location: 'ORIGIN',
              notes: 'this is a test request',
              quantity_1: 0,
              quantity_2: 0,
              shipment_id: '4612cfed-acbd-47ca-840a-7b7de190d6d2',
              status: 'APPROVED',
              submitted_date: '2018-11-06T10:22:34.370Z',
              actual_amount_cents: 50,
              estimated_amount_cents: 50,
              tariff400ng_item: {
                code: '35A',
                created_at: '2018-11-05T16:13:22.946Z',
                discount_type: 'HHG',
                id: 'c6a865dd-324a-48a5-9b03-5db8dcd044d1',
                item: 'Third Party Service',
                location: 'ORIGIN',
                ref_code: 'NONE',
                requires_pre_approval: true,
                uom_1: 'CF',
                uom_2: 'NONE',
                updated_at: '2018-11-05T16:13:22.946Z',
              },
              tariff400ng_item_id: 'c6a865dd-324a-48a5-9b03-5db8dcd044d1',
              updated_at: '0001-01-01T00:00:00.000Z',
              invoice_id: '683b4683-256c-4108-8400-300f3d87012f',
            },
            'e2a787d5-ff90-4331-8caa-c4c11d50ghsdgha5': {
              approved_date: '0001-01-01T00:00:00.000Z',
              created_at: '0001-01-01T00:00:00.000Z',
              id: 'e2a787d5-ff90-4331-8caa-c4c11d50ghsdgha5',
              location: 'ORIGIN',
              notes: 'this is a test request',
              quantity_1: 0,
              quantity_2: 0,
              shipment_id: '4612cfed-acbd-47ca-840a-7b7de190d6d2',
              status: 'APPROVED',
              submitted_date: '2018-11-06T10:22:34.370Z',
              estimated_amount_cents: 50,
              tariff400ng_item: {
                code: '35A',
                created_at: '2018-11-05T16:13:22.946Z',
                discount_type: 'HHG',
                id: 'c6a865dd-324a-48a5-9b03-5db8dcd044d1',
                item: 'Third Party Service',
                location: 'ORIGIN',
                ref_code: 'NONE',
                requires_pre_approval: true,
                uom_1: 'CF',
                uom_2: 'NONE',
                updated_at: '2018-11-05T16:13:22.946Z',
              },
              tariff400ng_item_id: 'c6a865dd-324a-48a5-9b03-5db8dcd044d1',
              updated_at: '0001-01-01T00:00:00.000Z',
            },
          },
        },
      };
    });
    it('selectUnbilledShipmentLineItems returns an array of items when shipment id is found', () => {
      expect(selectUnbilledShipmentLineItems(state, '4612cfed-acbd-47ca-840a-7b7de190d6d2').length).toEqual(2);
    });
    it('selectSortedPreApprovalShipmentLineItems returns pre-approval line items that are filtered by a shipmentId', () => {
      expect(selectSortedPreApprovalShipmentLineItems(state, '4612cfed-acbd-47ca-840a-7b7de190d6723').length).toEqual(
        2,
      );
    });
    it('selectSortedPreApprovalShipmentLineItems returns all line items if no shipmentId is passed', () => {
      expect(selectSortedPreApprovalShipmentLineItems(state).length).toEqual(5);
    });
    it('selectInvoiceShipmentLineItems returns all line items if no shipmentId is passed', () => {
      expect(selectInvoiceShipmentLineItems(state, '683b4683-256c-4108-8400-300f3d87012f').length).toEqual(2);
    });
  });
});

import * as utils from './utils';

describe('utils', () => {
  describe('upsert', () => {
    const item = { id: 'foo', name: 'something' };
    describe('when upserting a new item to an array', () => {
      const arr = [
        { id: 'bar', name: 'foo' },
        { id: 'baz', name: 'baz' },
      ];
      utils.upsert(arr, item);
      it('should be appended to the array', () => {
        expect(arr).toEqual([{ id: 'bar', name: 'foo' }, { id: 'baz', name: 'baz' }, item]);
      });
    });
    describe('when upserting an update to an array', () => {
      const arr = [
        { id: 'foo', name: 'foo' },
        { id: 'baz', name: 'baz' },
      ];
      utils.upsert(arr, item);
      it('should be appended to the array', () => {
        expect(arr).toEqual([
          { id: 'foo', name: 'something' },
          { id: 'baz', name: 'baz' },
        ]);
      });
    });
  });

  describe('fetch Active', () => {
    describe('when there are no foos', () => {
      const foos = null;
      const res = utils.fetchActive(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
    describe('when there are some active and some inactive foos', () => {
      const foos = [
        { status: 'CANCELED', id: 'foo0' },
        { status: 'DRAFT', id: 'foo1' },
        { status: 'SUBMITTED', id: 'foo2' },
      ];
      const res = utils.fetchActive(foos);
      it('should return the first active foo', () => {
        expect(res.id).toEqual('foo1');
      });
    });
    describe('when there are only inactive foos', () => {
      const foos = [
        { status: 'CANCELED', id: 'foo' },
        { status: 'COMPLETED', id: 'foo' },
      ];
      const res = utils.fetchActive(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
  });

  describe('fetch Active Shipment', () => {
    describe('when there are no foos', () => {
      const foos = null;
      const res = utils.fetchActiveShipment(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
    describe('when there are some active and some inactive foos', () => {
      const foos = [
        { status: 'CANCELED', id: 'foo0' },
        { status: 'DRAFT', id: 'foo1' },
        { status: 'SUBMITTED', id: 'foo2' },
      ];
      const res = utils.fetchActiveShipment(foos);
      it('should return the first active foo', () => {
        expect(res.id).toEqual('foo1');
      });
    });
    describe('when there are only inactive foos', () => {
      const foos = [{ status: 'CANCELED', id: 'foo' }];
      const res = utils.fetchActiveShipment(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
  });

  it('check if 2nd and 3rd addresses should be cleared from prime shipment create payload', () => {
    const ppmValues = {
      shipmentType: 'PPM',
      ppmShipment: {
        hasSecondaryPickupAddress: 'false',
        hasTertiaryPickupAddress: 'false',
        hasSecondaryDestinationAddress: 'false',
        hasTertiaryDestinationAddress: 'false',
        secondaryPickupAddress: '',
        tertiaryPickupAddress: '',
        secondaryDestinationAddress: '',
        tertiaryDestinationAddress: '',
      },
    };
    const hhgValues = {
      shipmentType: 'HHG',
      hasSecondaryPickupAddress: 'false',
      hasTertiaryPickupAddress: 'false',
      hasSecondaryDestinationAddress: 'false',
      hasTertiaryDestinationAddress: 'false',
      secondaryPickupAddress: '',
      tertiaryPickupAddress: '',
      secondaryDestinationAddress: '',
      tertiaryDestinationAddress: '',
    };

    const updatedPPMValues = utils.checkAddressTogglesToClearAddresses(ppmValues);
    expect(updatedPPMValues).toEqual({
      shipmentType: 'PPM',
      ppmShipment: {
        hasSecondaryPickupAddress: 'false',
        hasTertiaryPickupAddress: 'false',
        hasSecondaryDestinationAddress: 'false',
        hasTertiaryDestinationAddress: 'false',
        secondaryPickupAddress: {},
        tertiaryPickupAddress: {},
        secondaryDestinationAddress: {},
        tertiaryDestinationAddress: {},
      },
    });

    const updatedHHGValues = utils.checkAddressTogglesToClearAddresses(hhgValues);
    expect(updatedHHGValues).toEqual({
      shipmentType: 'HHG',
      hasSecondaryPickupAddress: 'false',
      hasTertiaryPickupAddress: 'false',
      hasSecondaryDestinationAddress: 'false',
      hasTertiaryDestinationAddress: 'false',
      secondaryPickupAddress: {},
      tertiaryPickupAddress: {},
      secondaryDestinationAddress: {},
      tertiaryDestinationAddress: {},
    });
  });
});

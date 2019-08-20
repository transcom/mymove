import { sitTotalDaysUsed, sitDaysUsed } from './calculator';
import MockDate from 'mockdate';

const requestedStorageInTransit = {
  status: 'REQUESTED',
  location: 'DESTINATION',
  estimated_start_date: '2019-05-15',
};

const inSitStorageInTransit = {
  status: 'IN_SIT',
  location: 'DESTINATION',
  estimated_start_date: '2019-05-10',
  authorized_start_date: '2019-05-10',
  actual_start_date: '2019-05-10',
};

const releasedStorageInTransit = {
  status: 'RELEASED',
  location: 'ORIGIN',
  estimated_start_date: '2019-05-07',
  authorized_start_date: '2019-05-07',
  actual_start_date: '2019-05-07',
  out_date: '2019-05-11',
};

const deliveredStorageInTransit = {
  status: 'DELIVERED',
  location: 'ORIGIN',
  estimated_start_date: '2019-05-02',
  authorized_start_date: '2019-05-03',
  actual_start_date: '2019-05-03',
  out_date: '2019-05-14',
};

const inSitFutureDateStorageInTransit = {
  status: 'IN_SIT',
  location: 'DESTINATION',
  estimated_start_date: '2019-05-10',
  authorized_start_date: '2019-05-10',
  actual_start_date: '2019-05-20',
};

describe('SIT calculator', () => {
  beforeAll(() => {
    MockDate.set('2019-05-12');
  });

  afterAll(() => MockDate.reset());

  describe('calculate SIT days used', () => {
    it('requested SIT', () => {
      expect(sitDaysUsed(requestedStorageInTransit)).toEqual(0);
    });

    it('in SIT', () => {
      expect(sitDaysUsed(inSitStorageInTransit)).toEqual(3);
    });

    it('released SIT', () => {
      expect(sitDaysUsed(releasedStorageInTransit)).toEqual(5);
    });

    it('delivered SIT', () => {
      expect(sitDaysUsed(deliveredStorageInTransit)).toEqual(12);
    });

    it('start date after today', () => {
      expect(sitDaysUsed(inSitFutureDateStorageInTransit)).toEqual(0);
    });
  });

  describe('calculate total SIT days used', () => {
    it('all SIT records', () => {
      const all = [
        requestedStorageInTransit,
        inSitStorageInTransit,
        releasedStorageInTransit,
        deliveredStorageInTransit,
      ];
      expect(sitTotalDaysUsed(all)).toEqual(20);
    });

    it('empty array', () => {
      expect(sitTotalDaysUsed([])).toEqual(0);
    });
  });
});

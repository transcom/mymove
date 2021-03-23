import { makeCalculations } from './helpers';
import testParams from './serviceItemTestParams';

describe('makeCalculations', () => {
  it('returns correct data for DomesticLongHaul', () => {
    const result = makeCalculations('DLH', 99999, testParams.DomesticLongHaul);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '210',
        label: 'Mileage',
        details: ['Zip 210 to Zip 910'],
      },
      {
        value: '1.033',
        label: 'Baseline linehaul price',
        details: ['Domestic non-peak', 'Origin service area: 176', 'Pickup date: 11 Mar 2020'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [''],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticShortHaul', () => {
    const result = makeCalculations('DSH', 99999, testParams.DomesticShortHaul);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '210',
        label: 'Mileage',
        details: ['Zip 32210 to Zip 91910'],
      },
      {
        value: '1.033',
        label: 'Baseline linehaul price',
        details: ['Domestic non-peak', 'Origin service area: 176', 'Pickup date: 11 Mar 2020'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [''],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticOrignPrice', () => {
    const result = makeCalculations('DOP', 99999, testParams.DomesticOrignPrice);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticDestinationPrice', () => {
    const result = makeCalculations('DDP', 99999, testParams.DomesticDestinationPrice);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticOrigin1stSIT', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticOrigin1stSIT);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticDestination1stSIT', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticDestination1stSIT);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticOriginAdditionalSIT', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticOriginAdditionalSIT);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticDestinationAdditionalSIT', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticDestinationAdditionalSIT);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticOriginSITDelivery', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticOriginSITDelivery);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticDestinationSITDelivery', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticDestinationSITDelivery);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticPacking', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticPacking);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticUnpacking', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticUnpacking);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticCrating', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticCrating);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticCratingStandalone', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticCratingStandalone);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticUncrating', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticUncrating);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticOriginShuttleService', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticOriginShuttleService);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticDestinationShuttleService', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticDestinationShuttleService);
    expect(result).toEqual([]);
  });

  it('returns correct data for NonStandardHHG', () => {
    const result = makeCalculations('?', 99999, testParams.NonStandardHHG);
    expect(result).toEqual([]);
  });

  it('returns correct data for NonStandardUB', () => {
    const result = makeCalculations('?', 99999, testParams.NonStandardUB);
    expect(result).toEqual([]);
  });

  it('returns correct data for FuelSurchage', () => {
    const result = makeCalculations('?', 99999, testParams.FuelSurchage);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticMobileHomeFactor', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticMobileHomeFactor);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticTowAwayBoatFactor', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticTowAwayBoatFactor);
    expect(result).toEqual([]);
  });

  it('returns correct data for DomesticNTSPackingFactor', () => {
    const result = makeCalculations('?', 99999, testParams.DomesticNTSPackingFactor);
    expect(result).toEqual([]);
  });
});

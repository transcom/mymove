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
        details: ['ZIP 210 to ZIP 910'],
      },
      {
        value: '1.71',
        label: 'Baseline linehaul price',
        details: ['Domestic non-peak', 'Origin service area: 176', 'Requested pickup: 09 Mar 2020'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
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
        value: '32210',
        label: 'Mileage',
        details: ['ZIP 32210 to ZIP 91910'],
      },
      {
        value: '1.71',
        label: 'Baseline shorthaul price',
        details: ['Domestic non-peak', 'Origin service area: 176', 'Requested pickup: 09 Mar 2020'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticOrignPrice', () => {
    const resultDOP = makeCalculations('DOP', 99998, testParams.DomesticOriginPrice);
    expect(resultDOP).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '1.71',
        label: 'Origin price',
        details: ['Origin service area: 176', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.98',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticDestinationPrice', () => {
    const result = makeCalculations('DDP', 99999, testParams.DomesticDestinationPrice);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '1.71',
        label: 'Destination price',
        details: ['Destination service area: 080', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticOrigin1stSIT', () => {
    const result = makeCalculations('DOFSIT', 99999, testParams.DomesticOrigin1stSIT);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '1.71',
        label: 'Origin price',
        details: ['Origin service area: 176', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticDestination1stSIT', () => {
    const result = makeCalculations('DDFSIT', 99999, testParams.DomesticDestination1stSIT);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '1.71',
        label: 'Destination price',
        details: ['Destination service area: 080', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticOriginAdditionalSIT', () => {
    const result = makeCalculations('DOASIT', 99999, testParams.DomesticOriginAdditionalSIT);
    expect(result).toEqual([
      {
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
        label: 'Billable weight (cwt)',
        value: '85 cwt',
      },
      {
        details: [],
        label: 'Days in SIT',
        value: '2',
      },
      {
        details: ['Origin service area: 176', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
        label: 'Additional day SIT price',
        value: '1.71',
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticDestinationAdditionalSIT', () => {
    const result = makeCalculations('DDASIT', 99999, testParams.DomesticDestinationAdditionalSIT);
    expect(result).toEqual([
      {
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
        label: 'Billable weight (cwt)',
        value: '85 cwt',
      },
      {
        details: [],
        label: 'Days in SIT',
        value: '2',
      },
      {
        details: ['Destination service area: 080', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
        label: 'Additional day SIT price',
        value: '1.71',
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticOriginSITPickup', () => {
    const result = makeCalculations('DOPSIT', 99999, testParams.DomesticOriginSITPickup);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '29',
        label: 'Mileage',
        details: ['ZIP 90210 to ZIP 90211'],
      },
      {
        value: '1.71',
        label: 'SIT pickup price',
        details: ['Origin SIT schedule: 3', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticDestinationSITDelivery', () => {
    const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDelivery);
    expect(result).toEqual([
      {
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
        label: 'Billable weight (cwt)',
        value: '85 cwt',
      },
      {
        value: '29',
        label: 'Mileage',
        details: ['ZIP 91910 to ZIP 94535'],
      },
      {
        details: ['Destination SIT schedule: 3', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
        label: 'SIT delivery price',
        value: '1.71',
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticPacking', () => {
    const result = makeCalculations('DPK', 99999, testParams.DomesticPacking);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '1.71',
        label: 'Pack price',
        details: ['Origin service schedule: 3', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  it('returns correct data for DomesticUnpacking', () => {
    const result = makeCalculations('DUPK', 99999, testParams.DomesticUnpacking);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '1.71',
        label: 'Unpack price',
        details: ['Destination service schedule: 3', 'Requested pickup: 09 Mar 2020', 'Domestic non-peak'],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: ['Base year: 2'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
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
    const resultFSC = makeCalculations('FSC', 99999, testParams.FuelSurchage);
    expect(resultFSC).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000 lbs'],
      },
      {
        value: '210',
        label: 'Mileage',
        details: ['ZIP 210 to ZIP 910'],
      },
      {
        value: '0.09',
        label: 'Fuel surcharge price (per mi)',
        details: ['EIA diesel: $2.73', 'Weight-based distance multiplier: 0.000417', 'Pickup date: 11 Mar 2020'],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [''],
      },
    ]);
  });

  // it('returns correct data for DomesticMobileHomeFactor', () => {
  //   const result = makeCalculations('?', 99999, testParams.DomesticMobileHomeFactor);
  //   expect(result).toEqual([]);
  // });

  // it('returns correct data for DomesticTowAwayBoatFactor', () => {
  //   const result = makeCalculations('?', 99999, testParams.DomesticTowAwayBoatFactor);
  //   expect(result).toEqual([]);
  // });

  // it('returns correct data for DomesticNTSPackingFactor', () => {
  //   const result = makeCalculations('?', 99999, testParams.DomesticNTSPackingFactor);
  //   expect(result).toEqual([]);
  // });
});

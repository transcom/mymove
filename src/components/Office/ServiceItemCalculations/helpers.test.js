import { makeCalculations } from './helpers';
import testParams from './serviceItemTestParams';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('makeCalculations', () => {
  it('returns correct data for DomesticLongHaul', () => {
    const result = makeCalculations('DLH', 99999, testParams.DomesticLongHaul, testParams.additionalCratingDataDCRT);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Mileage':
          expect(result[i].details).toEqual([{ text: 'ZIP 32210 to ZIP 91910', styles: {} }]);
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticLongHaul for NTS-release', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaul,
      testParams.additionalCratingDataDCRT,
      SHIPMENT_OPTIONS.NTSR,
    );
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Mileage':
          expect(result[i].details).toEqual([{ text: 'ZIP 32210 to ZIP 91910', styles: {} }]);
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticLongHaul with reweigh weight', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaulWithReweigh,
      testParams.additionalCratingDataDCRT,
    );
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Mileage':
          expect(result[i].details).toEqual([{ text: 'ZIP 32210 to ZIP 91910', styles: {} }]);
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticLongHaul weigh reweigh and adjusted weight', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaulWeightWithAdjustedAndReweigh,
      testParams.additionalCratingDataDCRT,
    );
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Mileage':
          expect(result[i].details).toEqual([{ text: 'ZIP 32210 to ZIP 91910', styles: {} }]);
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticLongHaul with no reweigh but billable weight adjusted', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaulWithAdjusted,
      testParams.additionalCratingDataDCRT,
    );
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Mileage':
          expect(result[i].details).toEqual([{ text: 'ZIP 32210 to ZIP 91910', styles: {} }]);
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticShortHaul', () => {
    const result = makeCalculations('DSH', 99999, testParams.DomesticShortHaul);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Mileage':
          expect(result[i].details).toEqual([{ text: 'ZIP 32210 to ZIP 91910', styles: {} }]);
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticOriginPrice', () => {
    const result = makeCalculations('DOP', 99998, testParams.DomesticOriginPrice);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Origin price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.98');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticDestinationPrice', () => {
    const result = makeCalculations('DDP', 99999, testParams.DomesticDestinationPrice);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Destination price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticOrigin1stSIT', () => {
    const result = makeCalculations('DOFSIT', 99999, testParams.DomesticOrigin1stSIT);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Origin price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticDestination1stSIT', () => {
    const result = makeCalculations('DDFSIT', 99999, testParams.DomesticDestination1stSIT);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Destination price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Baseline linehaul price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticOriginAdditionalSIT', () => {
    const result = makeCalculations('DOASIT', 99999, testParams.DomesticOriginAdditionalSIT);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'SIT days invoiced':
          expect(result[i].value).toEqual('2');
          break;
        case 'Additional day SIT price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });
});
describe('returns correct data for DomesticDestinationAdditionalSIT', () => {
  const result = makeCalculations('DDASIT', 99999, testParams.DomesticDestinationAdditionalSIT);
  for (let i = 0; i < result.length; i += 1) {
    switch (result[i].label) {
      case 'Billable weight (cwt)':
        expect(result[i].value).toEqual('85 cwt');
        break;
      case 'SIT days invoiced':
        expect(result[i].value).toEqual('2');
        break;
      case 'Additional day SIT price':
        expect(result[i].value).toEqual('1.71');
        break;
      case 'Price escalation factor':
        expect(result[i].value).toEqual('1.033');
        break;
      case 'Fuel rate adjustment':
        expect(result[i].value).toEqual('$999.99');
        break;
      default:
        break;
    }
  }
});

it('returns correct data for DomesticOriginSITPickup', () => {
  const result = makeCalculations('DOPSIT', 99999, testParams.DomesticOriginSITPickup);
  for (let i = 0; i < result.length; i += 1) {
    switch (result[i].label) {
      case 'Billable weight (cwt)':
        expect(result[i].value).toEqual('85 cwt');
        break;
      case 'Mileage':
        expect(result[i].value).toEqual('29');
        break;
      case 'SIT pickup price':
        expect(result[i].value).toEqual('1.71');
        break;
      case 'Price escalation factor':
        expect(result[i].value).toEqual('1.033');
        break;
      case 'Fuel rate adjustment':
        expect(result[i].value).toEqual('$999.99');
        break;
      default:
        break;
    }
  }
});

describe('DomesticDestinationSITDelivery', () => {
  it('returns the correct data for mileage above 50', () => {
    const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDeliveryLonghaul);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Mileage':
          expect(result[i].value).toEqual('51');
          break;
        case 'SIT pickup price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns the correct data for mileage below 50 with matching ZIP3s', () => {
    const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDeliveryMatchingZip3);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Mileage':
          expect(result[i].value).toEqual('3');
          break;
        case 'SIT delivery price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns the correct data for mileage below 50 with non-matching ZIP3s', () => {
    const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDelivery);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'SIT delivery price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticPacking', () => {
    const result = makeCalculations('DPK', 99999, testParams.DomesticPacking);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Pack price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticNTSPacking', () => {
    const result = makeCalculations('DNPK', 99999, testParams.DomesticNTSPacking);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Pack price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'NTS packing factor':
          expect(result[i].value).toEqual('1.35');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticUnpacking', () => {
    const result = makeCalculations('DUPK', 99999, testParams.DomesticUnpacking);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Unpack price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticCrating', () => {
    const result = makeCalculations('DCRT', 99999, testParams.DomesticCrating, testParams.additionalCratingDataDCRT);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Crating size (cu ft)':
          expect(result[i].value).toEqual('4.00');
          break;
        case 'Crating price (per cu ft)':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticUncrating', () => {
    const result = makeCalculations('DUCRT', 99999, testParams.DomesticUncrating, testParams.additionalCratingDataDCRT);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Crating size (cu ft)':
          expect(result[i].value).toEqual('4.00');
          break;
        case 'Uncrating price (per cu ft)':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticOriginShuttleService', () => {
    const result = makeCalculations('DOSHUT', 99999, testParams.DomesticOriginShuttleService);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Origin price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for DomesticDestinationShuttleService', () => {
    const result = makeCalculations('DDSHUT', 99999, testParams.DomesticDestinationShuttleService);
    for (let i = 0; i < result.length; i += 1) {
      switch (result[i].label) {
        case 'Billable weight (cwt)':
          expect(result[i].value).toEqual('85 cwt');
          break;
        case 'Destination price':
          expect(result[i].value).toEqual('1.71');
          break;
        case 'Price escalation factor':
          expect(result[i].value).toEqual('1.033');
          break;
        case 'Fuel rate adjustment':
          expect(result[i].value).toEqual('$999.99');
          break;
        default:
          break;
      }
    }
  });

  it('returns correct data for NonStandardHHG', () => {
    const result = makeCalculations('?', 99999, testParams.NonStandardHHG);
    expect(result).toEqual([]);
  });

  it('returns correct data for NonStandardUB', () => {
    const result = makeCalculations('?', 99999, testParams.NonStandardUB);
    expect(result).toEqual([]);
  });

  it('FuelSurcharge returns correct data for FSC', () => {
    const resultFSC = makeCalculations('FSC', 99999, testParams.FuelSurchage);
    expect(resultFSC).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [{ text: 'Estimated: 8,000 lbs', styles: {} }],
      },
      {
        value: '210',
        label: 'Mileage',
        details: [{ text: 'ZIP 32210 to ZIP 91910', styles: {} }],
      },
      {
        value: '0.1',
        label: 'Fuel surcharge price (per mi)',
        details: [
          { text: 'EIA diesel: $2.73', styles: {} },
          { text: 'FRA: 2.0', styles: {} },
          { text: 'Weight-based distance multiplier: 0.000417', styles: {} },
          { text: 'Pickup date: 11 Mar 2020', styles: {} },
        ],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('FuelSurcharge returns correct data for DOSFSC', () => {
    const resultFSC = makeCalculations('DOSFSC', 99999, testParams.DomesticOriginSITFuelSurchage);
    expect(resultFSC).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [{ text: 'Estimated: 8,000 lbs', styles: {} }],
      },
      {
        value: '29',
        label: 'Mileage into SIT',
        details: [{ text: 'ZIP 90210 to ZIP 90211', styles: {} }],
      },
      {
        value: '0.0',
        label: 'SIT fuel surcharge price (per mi)',
        details: [
          { text: 'EIA diesel: $2.73', styles: {} },
          { text: 'FRA: 2.0', styles: {} },
          { text: 'Weight-based distance multiplier: 0.000417', styles: {} },
          { text: 'Pickup date: 11 Mar 2020', styles: {} },
        ],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('FuelSurcharge returns correct data for DDSFSC', () => {
    const resultFSC = makeCalculations('DDSFSC', 99999, testParams.DomesticDestinationSITFuelSurchage);
    expect(resultFSC).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [{ text: 'Estimated: 8,000 lbs', styles: {} }],
      },
      {
        value: '29',
        label: 'Mileage out of SIT',
        details: [{ text: 'ZIP 91910 to ZIP 94535', styles: {} }],
      },
      {
        value: '0.0',
        label: 'SIT fuel surcharge price (per mi)',
        details: [
          { text: 'EIA diesel: $2.73', styles: {} },
          { text: 'FRA: 2.0', styles: {} },
          { text: 'Weight-based distance multiplier: 0.000417', styles: {} },
          { text: 'Pickup date: 11 Mar 2020', styles: {} },
        ],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
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
});

import { makeCalculations } from './helpers';
import testParams from './serviceItemTestParams';

import { SHIPMENT_OPTIONS } from 'shared/constants';

function testData(code) {
  let result;
  if (code === 'DCRT' || code === 'DUCRT') {
    result = {
      ...result,
      'Crating size (cu ft)': '4.00',
    };
  }
  if (code === 'DCRT') {
    result = {
      ...result,
      'Crating price (per cu ft)': '1.71',
    };
  } else if (code === 'DUCRT') {
    result = {
      ...result,
      'Uncrating price (per cu ft)': '1.71',
    };
  } else {
    result = {
      ...result,
      'Billable weight (cwt)': '85 cwt',
    };
  }
  if (code === 'DDDSIT') {
    result = {
      ...result,
      Mileage: '51',
      'SIT delivery price': '1.71',
    };
  } else if (code === 'DDDSITb') {
    result = {
      ...result,
      Mileage: '3',
      'SIT delivery price': '1.71',
    };
  } else if (code === 'DDDSITc') {
    result = {
      ...result,
      'SIT delivery price': '1.71',
    };
  } else if (code !== 'DOFSIT' && code !== 'DDFSIT' && code !== 'DOPSIT' && code.includes('SIT')) {
    result = {
      ...result,
      'SIT days invoiced': '2',
      'Additional day SIT price': '1.71',
    };
  }
  if (code === 'DOPSIT') {
    result = {
      ...result,
      Mileage: '29',
      'SIT pickup price': '1.71',
    };
  } else if (code === 'DLH') {
    result = {
      ...result,
      Mileage: '210',
      'Baseline linehaul price': '1.71',
    };
  } else if (code === 'DSH') {
    result = {
      ...result,
      Mileage: '210',
      'Baseline shorthaul price': '1.71',
    };
  }
  if (code === 'DOP' || code === 'DOFSIT') {
    result = {
      ...result,
      'Origin price': '1.71',
    };
  } else if (code === 'DDP') {
    result = {
      ...result,
      'Destination price': '1.71',
    };
  }
  if (!code.includes('FSC')) {
    result = {
      ...result,
      'Price escalation factor': '1.033',
    };
  }
  if (code.includes('FSC')) {
    result = {
      ...result,
      'Total:': '$999.98',
    };
  } else {
    result = {
      ...result,
      'Total:': '$999.99',
    };
  }

  return result;
}

function testAB(a, b) {
  const keys = Object.keys(b);
  for (let j = 0; j < keys.length; j += 1) {
    for (let i = 0; i < a.length; i += 1) {
      if (i < a.length - 1) {
        if (a[i].label === keys[j]) {
          expect(a[i].value).toEqual(b[keys[j]]);
          break;
        }
      } else {
        expect(a[i].value).toEqual(b[keys[j]]);
        break;
      }
    }
  }
}

describe('makeCalculations', () => {
  it('returns correct data for DomesticLongHaul', () => {
    const result = makeCalculations('DLH', 99999, testParams.DomesticLongHaul, testParams.additionalCratingDataDCRT);
    const expected = testData('DLH');

    testAB(result, expected);
  });

  it('returns correct data for DomesticLongHaul for NTS-release', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaul,
      testParams.additionalCratingDataDCRT,
      SHIPMENT_OPTIONS.NTSR,
    );
    const expected = testData('DLH');

    testAB(result, expected);
  });

  it('returns correct data for DomesticLongHaul with reweigh weight', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaulWithReweigh,
      testParams.additionalCratingDataDCRT,
    );
    const expected = testData('DLH');

    testAB(result, expected);
  });

  it('returns correct data for DomesticLongHaul weigh reweigh and adjusted weight', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaulWeightWithAdjustedAndReweigh,
      testParams.additionalCratingDataDCRT,
    );
    const expected = testData('DLH');

    testAB(result, expected);
  });

  it('returns correct data for DomesticLongHaul with no reweigh but billable weight adjusted', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaulWithAdjusted,
      testParams.additionalCratingDataDCRT,
    );
    const expected = testData('DLH');

    testAB(result, expected);
  });

  it('returns correct data for DomesticShortHaul', () => {
    const result = makeCalculations('DSH', 99999, testParams.DomesticShortHaul);
    const expected = testData('DSH');

    testAB(result, expected);
  });

  it('returns correct data for DomesticOriginPrice', () => {
    const result = makeCalculations('DOP', 99999, testParams.DomesticOriginPrice);
    const expected = testData('DOP');

    testAB(result, expected);
  });

  it('returns correct data for DomesticDestinationPrice', () => {
    const result = makeCalculations('DDP', 99999, testParams.DomesticDestinationPrice);
    const expected = testData('DDP');

    testAB(result, expected);
  });

  it('returns correct data for DomesticOrigin1stSIT', () => {
    const result = makeCalculations('DOFSIT', 99999, testParams.DomesticOrigin1stSIT);
    const expected = testData('DOFSIT');

    testAB(result, expected);
  });

  it('returns correct data for DomesticDestination1stSIT', () => {
    const result = makeCalculations('DDFSIT', 99999, testParams.DomesticDestination1stSIT);
    const expected = testData('DDFSIT');

    testAB(result, expected);
  });

  it('returns correct data for DomesticOriginAdditionalSIT', () => {
    const result = makeCalculations('DOASIT', 99999, testParams.DomesticOriginAdditionalSIT);
    const expected = testData('DOASIT');

    testAB(result, expected);
  });
});
describe('returns correct data for DomesticDestinationAdditionalSIT', () => {
  const result = makeCalculations('DDASIT', 99999, testParams.DomesticDestinationAdditionalSIT);
  const expected = testData('DDASIT');

  testAB(result, expected);
});

it('returns correct data for DomesticOriginSITPickup', () => {
  const result = makeCalculations('DOPSIT', 99999, testParams.DomesticOriginSITPickup);
  const expected = testData('DOPSIT');

  testAB(result, expected);
});

describe('DomesticDestinationSITDelivery', () => {
  it('returns the correct data for mileage above 50', () => {
    const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDeliveryLonghaul);
    const expected = testData('DDDSIT');

    testAB(result, expected);
  });

  it('returns the correct data for mileage below 50 with matching ZIP3s', () => {
    const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDeliveryMatchingZip3);
    const expected = testData('DDDSITb');

    testAB(result, expected);
  });

  it('returns the correct data for mileage below 50 with non-matching ZIP3s', () => {
    const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDelivery);
    const expected = testData('DDDSITc');

    testAB(result, expected);
  });

  it('returns correct data for DomesticPacking', () => {
    const result = makeCalculations('DPK', 99999, testParams.DomesticPacking);
    const expected = testData('DPK');

    testAB(result, expected);
  });

  it('returns correct data for DomesticNTSPacking', () => {
    const result = makeCalculations('DNPK', 99999, testParams.DomesticNTSPacking);
    const expected = testData('DNPK');

    testAB(result, expected);
  });

  it('returns correct data for DomesticUnpacking', () => {
    const result = makeCalculations('DUPK', 99999, testParams.DomesticUnpacking);
    const expected = testData('DUPK');

    testAB(result, expected);
  });

  it('returns correct data for DomesticCrating', () => {
    const result = makeCalculations('DCRT', 99999, testParams.DomesticCrating, testParams.additionalCratingDataDCRT);
    const expected = testData('DCRT');

    testAB(result, expected);
  });

  it('returns correct data for DomesticUncrating', () => {
    const result = makeCalculations('DUCRT', 99999, testParams.DomesticUncrating, testParams.additionalCratingDataDCRT);
    const expected = testData('DUCRT');

    testAB(result, expected);
  });

  it('returns correct data for DomesticOriginShuttleService', () => {
    const result = makeCalculations('DOSHUT', 99999, testParams.DomesticOriginShuttleService);
    const expected = testData('DOSHUT');

    testAB(result, expected);
  });

  it('returns correct data for DomesticDestinationShuttleService', () => {
    const result = makeCalculations('DDSHUT', 99999, testParams.DomesticDestinationShuttleService);
    const expected = testData('DDSHUT');

    testAB(result, expected);
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
    const result = makeCalculations('FSC', 99998, testParams.FuelSurchage);
    const expected = testData('FSC');

    testAB(result, expected);
  });

  it('FuelSurcharge returns correct data for DOSFSC', () => {
    const result = makeCalculations('DOSFSC', 99998, testParams.DomesticOriginSITFuelSurchage);
    const expected = testData('DOSFSC');

    testAB(result, expected);
  });

  it('FuelSurcharge returns correct data for DDSFSC', () => {
    const result = makeCalculations('DDSFSC', 99998, testParams.DomesticDestinationSITFuelSurchage);
    const expected = testData('DDSFSC');

    testAB(result, expected);
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

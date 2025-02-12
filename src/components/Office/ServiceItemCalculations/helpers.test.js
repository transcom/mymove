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
  if (code === 'DOP' || code === 'DOFSIT' || code === 'DOSHUT') {
    result = {
      ...result,
      'Origin price': '1.71',
    };
  } else if (code === 'DDP' || code === 'DDFSIT' || code === 'DDSHUT') {
    result = {
      ...result,
      'Destination price': '1.71',
    };
  }
  if (code === 'ISLH') {
    result = {
      ...result,
      'ISLH price': '1.71',
    };
  }
  if (code === 'UBP') {
    result = {
      ...result,
      'International UB price': '1.71',
    };
  }

  // Packing and Unpacking
  if (code === 'IHPK') {
    result = {
      ...result,
      'International Pack price': '1.71',
    };
  } else if (code === 'IHUPK') {
    result = {
      ...result,
      'International Unpack price': '1.71',
    };
  } else if (code === 'IUBPK') {
    result = {
      ...result,
      'International UB Pack price': '1.71',
    };
  } else if (code === 'IUBUPK') {
    result = {
      ...result,
      'International UB Unpack price': '1.71',
    };
  } else if (code.includes('UPK')) {
    result = {
      ...result,
      'Unpack price': '1.71',
    };
  } else if (code.includes('PK')) {
    result = {
      ...result,
      'Pack price': '1.71',
    };
  }

  if (code.includes('DNPK')) {
    result = {
      ...result,
      'NTS packing factor': '1.35',
    };
  }

  // FSC or not
  if (code === 'DOSFSC') {
    result = {
      ...result,
      'Mileage into SIT': '29',
      'SIT mileage factor': '0.012',
      'Total:': '$999.98',
    };
  } else if (code === 'DDSFSC') {
    result = {
      ...result,
      'Mileage out of SIT': '29',
      'SIT mileage factor': '0.012',
      'Total:': '$999.98',
    };
  } else if (code.includes('FSC')) {
    result = {
      ...result,
      Mileage: '210',
      'Mileage factor': '0.088',
      'Total:': '$999.98',
    };
  } else {
    result = {
      ...result,
      'Price escalation factor': '1.033',
      'Total:': '$999.99',
    };
  }

  return result;
}

function testAB(result, expected) {
  const expectedKeys = Object.keys(expected);
  for (let i = 0; i < expectedKeys.length || i < result.length; i += 1) {
    if (i >= expectedKeys.length) {
      expect(result[i].label).toEqual('');
    }
    if (i >= result.length) {
      expect('').toEqual(expectedKeys[i]);
    }
    expect(result[i].label).toEqual(expectedKeys[i]);
    expect(result[i].value).toEqual(expected[expectedKeys[i]]);
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

  it('returns correct data for DomesticDestinationAdditionalSIT', () => {
    const result = makeCalculations('DDASIT', 99999, testParams.DomesticDestinationAdditionalSIT);
    const expected = testData('DDASIT');

    testAB(result, expected);
  });

  it('returns correct data for DomesticOriginSITPickup', () => {
    const result = makeCalculations('DOPSIT', 99999, testParams.DomesticOriginSITPickup);
    const expected = testData('DOPSIT');

    testAB(result, expected);
  });
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
});

describe('Domestic pack, crate, shuttle', () => {
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
});

describe('International', () => {
  it('returns correct data for ISLH', () => {
    const result = makeCalculations('ISLH', 99999, testParams.InternationalShippingAndLinehaul);
    const expected = testData('ISLH');
    testAB(result, expected);
  });

  it('returns correct data for IHPK', () => {
    const result = makeCalculations('IHPK', 99999, testParams.InternationalHHGPack);
    const expected = testData('IHPK');
    testAB(result, expected);
  });

  it('returns correct data for IHUPK', () => {
    const result = makeCalculations('IHUPK', 99999, testParams.InternationalHHGUnpack);
    const expected = testData('IHUPK');
    testAB(result, expected);
  });

  it('returns correct data for POEFSC', () => {
    const result = makeCalculations('POEFSC', 99998, testParams.PortOfEmbarkation);
    const expected = testData('POEFSC');
    testAB(result, expected);
  });

  it('returns correct data for PODFSC', () => {
    const result = makeCalculations('PODFSC', 99998, testParams.PortOfDebarkation);
    const expected = testData('PODFSC');
    testAB(result, expected);
  });
});

describe('Unaccompanied Baggage', () => {
  it('UBP', () => {
    const result = makeCalculations('UBP', 99999, testParams.InternationalUBPrice);
    const expected = testData('UBP');
    testAB(result, expected);
  });

  it('UBP explicit', () => {
    const result = makeCalculations('UBP', 99999, testParams.InternationalUBPrice);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'International UB price': '1.71',
      'Price escalation factor': '1.033',
      'Total:': '$999.99',
    };
    testAB(result, expected);
  });

  it('IUBPK', () => {
    const result = makeCalculations('IUBPK', 99999, testParams.InternationalUBPackPrice);
    const expected = testData('IUBPK');
    testAB(result, expected);
  });

  it('IUBUPK', () => {
    const result = makeCalculations('IUBUPK', 99999, testParams.InternationalUBUnpackPrice);
    const expected = testData('IUBUPK');
    testAB(result, expected);
  });
});

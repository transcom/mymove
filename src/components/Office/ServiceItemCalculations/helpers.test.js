import { makeCalculations } from './helpers';
import testParams from './serviceItemTestParams';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { SERVICE_ITEM_CODES } from 'constants/serviceItems';

const expectedDlh = {
  'Billable weight (cwt)': '85 cwt',
  Mileage: '210',
  'Baseline linehaul price': '1.71',
  'Price escalation factor': '1.033',
  'Total: ': '$999.99',
};

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
    const result = makeCalculations(
      SERVICE_ITEM_CODES.DLH,
      99999,
      testParams.DomesticLongHaul,
      testParams.additionalCratingDataDCRT,
    );

    testAB(result, expectedDlh);
  });

  it('returns correct data for DomesticLongHaul for NTS-release', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.DLH,
      99999,
      testParams.DomesticLongHaul,
      testParams.additionalCratingDataDCRT,
      SHIPMENT_OPTIONS.NTSR,
    );

    testAB(result, expectedDlh);
  });

  it('returns correct data for DomesticLongHaul with reweigh weight', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.DLH,
      99999,
      testParams.DomesticLongHaulWithReweigh,
      testParams.additionalCratingDataDCRT,
    );

    testAB(result, expectedDlh);
  });

  it('returns correct data for DomesticLongHaul weigh reweigh and adjusted weight', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.DLH,
      99999,
      testParams.DomesticLongHaulWeightWithAdjustedAndReweigh,
      testParams.additionalCratingDataDCRT,
    );

    testAB(result, expectedDlh);
  });

  it('returns correct data for DomesticLongHaul with no reweigh but billable weight adjusted', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.DLH,
      99999,
      testParams.DomesticLongHaulWithAdjusted,
      testParams.additionalCratingDataDCRT,
    );

    testAB(result, expectedDlh);
  });

  it('returns correct data for DomesticShortHaul', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DSH, 99999, testParams.DomesticShortHaul);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      Mileage: '210',
      'Baseline shorthaul price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticOriginPrice', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DOP, 99999, testParams.DomesticOriginPrice);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Origin price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticDestinationPrice', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DDP, 99999, testParams.DomesticDestinationPrice);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Destination price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticOrigin1stSIT', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DOFSIT, 99999, testParams.DomesticOrigin1stSIT);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Origin price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticDestination1stSIT', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DDFSIT, 99999, testParams.DomesticDestination1stSIT);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Destination price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticOriginAdditionalSIT', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DOASIT, 99999, testParams.DomesticOriginAdditionalSIT);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'SIT days invoiced': '2',
      'Additional day SIT price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticDestinationAdditionalSIT', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DDASIT, 99999, testParams.DomesticDestinationAdditionalSIT);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'SIT days invoiced': '2',
      'Additional day SIT price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticOriginSITPickup', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DOPSIT, 99999, testParams.DomesticOriginSITPickup);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      Mileage: '29',
      'SIT pickup price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });
});

describe('DomesticDestinationSITDelivery', () => {
  it('returns the correct data for mileage above 50', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.DDDSIT,
      99999,
      testParams.DomesticDestinationSITDeliveryLonghaul,
    );
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      Mileage: '51',
      'SIT delivery price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns the correct data for mileage below 50 with matching ZIP3s', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.DDDSIT,
      99999,
      testParams.DomesticDestinationSITDeliveryMatchingZip3,
    );
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      Mileage: '3',
      'SIT delivery price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns the correct data for mileage below 50 with non-matching ZIP3s', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DDDSIT, 99999, testParams.DomesticDestinationSITDelivery);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'SIT delivery price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });
});

describe('Domestic pack, crate, shuttle', () => {
  it('returns correct data for DomesticPacking', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DPK, 99999, testParams.DomesticPacking);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Pack price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticNTSPacking', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DNPK, 99999, testParams.DomesticNTSPacking);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Pack price': '1.71',
      'NTS packing factor': '1.35',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticUnpacking', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DUPK, 99999, testParams.DomesticUnpacking);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Unpack price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticCrating', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.DCRT,
      99999,
      testParams.DomesticCrating,
      testParams.additionalCratingDataDCRT,
    );
    const expected = {
      'Crating size (cu ft)': '4.00',
      'Crating price (per cu ft)': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticUncrating', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.DUCRT,
      99999,
      testParams.DomesticUncrating,
      testParams.additionalCratingDataDCRT,
    );
    const expected = {
      'Crating size (cu ft)': '4.00',
      'Uncrating price (per cu ft)': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticOriginShuttleService', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DOSHUT, 99999, testParams.DomesticOriginShuttleService);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Origin price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for DomesticDestinationShuttleService', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DDSHUT, 99999, testParams.DomesticDestinationShuttleService);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Destination price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

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
    const result = makeCalculations(SERVICE_ITEM_CODES.FSC, 99998, testParams.FuelSurchage);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      Mileage: '210',
      'Mileage factor': '0.088',
      'Total: ': '$999.98',
    };

    testAB(result, expected);
  });

  it('FuelSurcharge returns correct data for DOSFSC', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DOSFSC, 99998, testParams.DomesticOriginSITFuelSurchage);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Mileage into SIT': '29',
      'SIT mileage factor': '0.012',
      'Total: ': '$999.98',
    };

    testAB(result, expected);
  });

  it('FuelSurcharge returns correct data for DDSFSC', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.DDSFSC, 99998, testParams.DomesticDestinationSITFuelSurchage);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'Mileage out of SIT': '29',
      'SIT mileage factor': '0.012',
      'Total: ': '$999.98',
    };

    testAB(result, expected);
  });
});

describe('International', () => {
  it('returns correct data for ISLH', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.ISLH, 99999, testParams.InternationalShippingAndLinehaul);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'ISLH price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for IHPK', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.IHPK, 99999, testParams.InternationalHHGPack);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'International Pack price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for IHUPK', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.IHUPK, 99999, testParams.InternationalHHGUnpack);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'International Unpack price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for POEFSC', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.POEFSC, 99998, testParams.PortOfEmbarkation);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      Mileage: '210',
      'Mileage factor': '0.088',
      'Total: ': '$999.98',
    };

    testAB(result, expected);
  });

  it('returns correct data for PODFSC', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.PODFSC, 99998, testParams.PortOfDebarkation);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      Mileage: '210',
      'Mileage factor': '0.088',
      'Total: ': '$999.98',
    };

    testAB(result, expected);
  });

  it('returns correct data for ICRT', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.ICRT,
      99999,
      testParams.InternationalCrating,
      testParams.additionalCratingDataDCRT,
    );
    const expected = {
      'Crating size (cu ft)': '4.00',
      'Crating price (per cu ft)': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('returns correct data for IUCRT', () => {
    const result = makeCalculations(
      SERVICE_ITEM_CODES.IUCRT,
      99999,
      testParams.InternationalUncrating,
      testParams.additionalCratingDataDCRT,
    );
    const expected = {
      'Crating size (cu ft)': '4.00',
      'Uncrating price (per cu ft)': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });
});

describe('Unaccompanied Baggage', () => {
  it('UBP', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.UBP, 99999, testParams.InternationalUBPrice);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'International UB price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('IUBPK', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.IUBPK, 99999, testParams.InternationalUBPackPrice);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'International UB Pack price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });

  it('IUBUPK', () => {
    const result = makeCalculations(SERVICE_ITEM_CODES.IUBUPK, 99999, testParams.InternationalUBUnpackPrice);
    const expected = {
      'Billable weight (cwt)': '85 cwt',
      'International UB Unpack price': '1.71',
      'Price escalation factor': '1.033',
      'Total: ': '$999.99',
    };

    testAB(result, expected);
  });
});

import { makeCalculations } from './helpers';
import testParams from './serviceItemTestParams';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('makeCalculations', () => {
  it('returns correct data for DomesticLongHaul', () => {
    const result = makeCalculations('DLH', 99999, testParams.DomesticLongHaul, testParams.additionalCratingDataDCRT);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '32210',
        label: 'Mileage',
        details: [{ text: 'ZIP 322 to ZIP 919', styles: {} }],
      },
      {
        value: '1.71',
        label: 'Baseline linehaul price',
        details: [
          { text: 'Domestic non-peak', styles: {} },
          { text: 'Origin service area: 176', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticLongHaul for NTS-release', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaul,
      testParams.additionalCratingDataDCRT,
      SHIPMENT_OPTIONS.NTSR,
    );
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '32210',
        label: 'Mileage',
        details: [{ text: 'ZIP 322 to ZIP 919', styles: {} }],
      },
      {
        value: '1.71',
        label: 'Baseline linehaul price',
        details: [
          { text: 'Domestic non-peak', styles: {} },
          { text: 'Origin service area: 176', styles: {} },
          { text: 'Actual pickup: 09 Mar 2020', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticLongHaul with reweigh weight', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaulWithReweigh,
      testParams.additionalCratingDataDCRT,
    );
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Reweigh: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Original: 8,500 lbs', styles: {} },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '32210',
        label: 'Mileage',
        details: [{ text: 'ZIP 322 to ZIP 919', styles: {} }],
      },
      {
        value: '1.71',
        label: 'Baseline linehaul price',
        details: [
          { text: 'Domestic non-peak', styles: {} },
          { text: 'Origin service area: 176', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticLongHaul weigh reweigh and adjusted weight', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaulWeightWithAdjustedAndReweigh,
      testParams.additionalCratingDataDCRT,
    );
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Adjusted: 500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Reweigh: 8,500 lbs', styles: {} },
          { text: 'Original: 8,500 lbs', styles: {} },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '32210',
        label: 'Mileage',
        details: [{ text: 'ZIP 322 to ZIP 919', styles: {} }],
      },
      {
        value: '1.71',
        label: 'Baseline linehaul price',
        details: [
          { text: 'Domestic non-peak', styles: {} },
          { text: 'Origin service area: 176', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticLongHaul with no reweigh but billable weight adjusted', () => {
    const result = makeCalculations(
      'DLH',
      99999,
      testParams.DomesticLongHaulWithAdjusted,
      testParams.additionalCratingDataDCRT,
    );
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Adjusted: 500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Original: 8,500 lbs', styles: {} },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '32210',
        label: 'Mileage',
        details: [{ text: 'ZIP 322 to ZIP 919', styles: {} }],
      },
      {
        value: '1.71',
        label: 'Baseline linehaul price',
        details: [
          { text: 'Domestic non-peak', styles: {} },
          { text: 'Origin service area: 176', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticShortHaul', () => {
    const result = makeCalculations('DSH', 99999, testParams.DomesticShortHaul);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '32210',
        label: 'Mileage',
        details: [{ text: 'ZIP 32210 to ZIP 91910', styles: {} }],
      },
      {
        value: '1.71',
        label: 'Baseline shorthaul price',
        details: [
          { text: 'Domestic non-peak', styles: {} },
          { text: 'Origin service area: 176', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticOriginPrice', () => {
    const resultDOP = makeCalculations('DOP', 99998, testParams.DomesticOriginPrice);
    expect(resultDOP).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '1.71',
        label: 'Origin price',
        details: [
          { text: 'Origin service area: 176', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.98',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticDestinationPrice', () => {
    const result = makeCalculations('DDP', 99999, testParams.DomesticDestinationPrice);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '1.71',
        label: 'Destination price',
        details: [
          { text: 'Destination service area: 080', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticOrigin1stSIT', () => {
    const result = makeCalculations('DOFSIT', 99999, testParams.DomesticOrigin1stSIT);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '1.71',
        label: 'Origin price',
        details: [
          { text: 'Origin service area: 176', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticDestination1stSIT', () => {
    const result = makeCalculations('DDFSIT', 99999, testParams.DomesticDestination1stSIT);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '1.71',
        label: 'Destination price',
        details: [
          { text: 'Destination service area: 080', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticOriginAdditionalSIT', () => {
    const result = makeCalculations('DOASIT', 99999, testParams.DomesticOriginAdditionalSIT);
    expect(result).toEqual([
      {
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
        label: 'Billable weight (cwt)',
        value: '85 cwt',
      },
      {
        details: [],
        label: 'SIT days invoiced',
        value: '2',
      },
      {
        details: [
          { text: 'Origin service area: 176', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
        label: 'Additional day SIT price',
        value: '1.71',
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  describe('returns correct data for DomesticDestinationAdditionalSIT', () => {
    const result = makeCalculations('DDASIT', 99999, testParams.DomesticDestinationAdditionalSIT);
    expect(result).toEqual([
      {
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
        label: 'Billable weight (cwt)',
        value: '85 cwt',
      },
      {
        details: [],
        label: 'SIT days invoiced',
        value: '2',
      },
      {
        details: [
          { text: 'Destination service area: 080', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
        label: 'Additional day SIT price',
        value: '1.71',
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticOriginSITPickup', () => {
    const result = makeCalculations('DOPSIT', 99999, testParams.DomesticOriginSITPickup);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '29',
        label: 'Mileage',
        details: [{ text: 'ZIP 90210 to ZIP 90211', styles: {} }],
      },
      {
        value: '1.71',
        label: 'SIT pickup price',
        details: [
          { text: 'Origin SIT schedule: 3', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  describe('DomesticDestinationSITDelivery', () => {
    it('returns the correct data for mileage above 50', () => {
      const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDeliveryLonghaul);
      expect(result).toEqual([
        {
          details: [
            { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
            { text: 'Estimated: 8,000 lbs', styles: {} },
          ],
          label: 'Billable weight (cwt)',
          value: '85 cwt',
        },
        {
          value: '51',
          label: 'Mileage',
          details: [{ text: 'ZIP 91910 to ZIP 94535', styles: {} }],
        },
        {
          details: [
            { text: 'Destination SIT schedule: 3', styles: {} },
            { text: 'Requested pickup: 09 Mar 2020', styles: {} },
            { text: 'Domestic non-peak', styles: {} },
          ],
          label: 'SIT delivery price',
          value: '1.71',
        },
        {
          value: '1.033',
          label: 'Price escalation factor',
          details: [{ text: 'Base year: 2', styles: {} }],
        },
        {
          value: '$999.99',
          label: 'Total amount requested',
          details: [{ text: '', styles: {} }],
        },
      ]);
    });

    it('returns the correct data for mileage below 50 with matching ZIP3s', () => {
      const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDeliveryMatchingZip3);
      expect(result).toEqual([
        {
          details: [
            { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
            { text: 'Estimated: 8,000 lbs', styles: {} },
          ],
          label: 'Billable weight (cwt)',
          value: '85 cwt',
        },
        {
          value: '3',
          label: 'Mileage',
          details: [{ text: 'ZIP 91910 to ZIP 91920', styles: {} }],
        },
        {
          details: [
            { text: 'Destination SIT schedule: 3', styles: {} },
            { text: 'Requested pickup: 09 Mar 2020', styles: {} },
            { text: 'Domestic non-peak', styles: {} },
          ],
          label: 'SIT delivery price',
          value: '1.71',
        },
        {
          value: '1.033',
          label: 'Price escalation factor',
          details: [{ text: 'Base year: 2', styles: {} }],
        },
        {
          value: '$999.99',
          label: 'Total amount requested',
          details: [{ text: '', styles: {} }],
        },
      ]);
    });

    it('returns the correct data for mileage below 50 with non-matching ZIP3s', () => {
      const result = makeCalculations('DDDSIT', 99999, testParams.DomesticDestinationSITDelivery);
      expect(result).toEqual([
        {
          details: [
            { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
            { text: 'Estimated: 8,000 lbs', styles: {} },
          ],
          label: 'Billable weight (cwt)',
          value: '85 cwt',
        },
        {
          details: [
            { text: 'Destination SIT schedule: 3', styles: {} },
            { text: 'Requested pickup: 09 Mar 2020', styles: {} },
            { text: 'Domestic non-peak', styles: {} },
            { text: '<=50 miles', styles: {} },
          ],
          label: 'SIT delivery price',
          value: '1.71',
        },
        {
          value: '1.033',
          label: 'Price escalation factor',
          details: [{ text: 'Base year: 2', styles: {} }],
        },
        {
          value: '$999.99',
          label: 'Total amount requested',
          details: [{ text: '', styles: {} }],
        },
      ]);
    });
  });

  it('returns correct data for DomesticPacking', () => {
    const result = makeCalculations('DPK', 99999, testParams.DomesticPacking);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '1.71',
        label: 'Pack price',
        details: [
          { text: 'Origin service schedule: 3', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticNTSPacking', () => {
    const result = makeCalculations('DNPK', 99999, testParams.DomesticNTSPacking);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '1.71',
        label: 'Pack price',
        details: [
          { text: 'Origin service schedule: 3', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
      },
      {
        value: '1.35',
        label: 'NTS packing factor',
        details: [],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticUnpacking', () => {
    const result = makeCalculations('DUPK', 99999, testParams.DomesticUnpacking);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Original: 8,500 lbs', styles: { fontWeight: 'bold' } },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '1.71',
        label: 'Unpack price',
        details: [
          { text: 'Destination service schedule: 3', styles: {} },
          { text: 'Requested pickup: 09 Mar 2020', styles: {} },
          { text: 'Domestic non-peak', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [{ text: 'Base year: 2', styles: {} }],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticCrating', () => {
    const result = makeCalculations('DCRT', 99999, testParams.DomesticCrating, testParams.additionalCratingDataDCRT);
    expect(result).toEqual([
      {
        value: '4.00',
        label: 'Crating size (cu ft)',
        details: [
          { text: 'Description: Grand piano', styles: {} },
          { text: 'Dimensions: 3x10x6 in', styles: {} },
        ],
      },
      {
        value: '1.71',
        label: 'Crating price (per cu ft)',
        details: [
          { text: 'Service schedule: 3', styles: {} },
          { text: 'Crating date: 09 Mar 2020', styles: {} },
          { text: 'Domestic', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
  });

  it('returns correct data for DomesticUncrating', () => {
    const result = makeCalculations('DUCRT', 99999, testParams.DomesticUncrating, testParams.additionalCratingDataDCRT);
    expect(result).toEqual([
      {
        details: [
          { text: 'Description: Grand piano', styles: {} },
          { text: 'Dimensions: 3x10x6 in', styles: {} },
        ],
        label: 'Crating size (cu ft)',
        value: '4.00',
      },
      {
        details: [
          { text: 'Service schedule: 3', styles: {} },
          { text: 'Uncrating date: 09 Mar 2020', styles: {} },
          { text: 'Domestic', styles: {} },
        ],
        label: 'Uncrating price (per cu ft)',
        value: '1.71',
      },
      {
        details: [],
        label: 'Price escalation factor',
        value: '1.033',
      },
      {
        details: [{ text: '', styles: {} }],
        label: 'Total amount requested',
        value: '$999.99',
      },
    ]);
  });

  it('returns correct data for DomesticOriginShuttleService', () => {
    const result = makeCalculations('DOSHUT', 99999, testParams.DomesticOriginShuttleService);
    expect(result).toEqual([
      {
        details: [
          { text: 'Shuttle weight: 8,500 lbs', styles: {} },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
        label: 'Billable weight (cwt)',
        value: '85 cwt',
      },
      {
        details: [
          { text: 'Service schedule: 3', styles: {} },
          { text: 'Pickup date: 09 Mar 2020', styles: {} },
          { text: 'Domestic', styles: {} },
        ],
        label: 'Origin price',
        value: '1.71',
      },
      {
        details: [],
        label: 'Price escalation factor',
        value: '1.033',
      },
      {
        details: [{ text: '', styles: {} }],
        label: 'Total amount requested',
        value: '$999.99',
      },
    ]);
  });

  it('returns correct data for DomesticDestinationShuttleService', () => {
    const result = makeCalculations('DDSHUT', 99999, testParams.DomesticDestinationShuttleService);
    expect(result).toEqual([
      {
        value: '85 cwt',
        label: 'Billable weight (cwt)',
        details: [
          { text: 'Shuttle weight: 8,500 lbs', styles: {} },
          { text: 'Estimated: 8,000 lbs', styles: {} },
        ],
      },
      {
        value: '1.71',
        label: 'Destination price',
        details: [
          { text: 'Service schedule: 3', styles: {} },
          { text: 'Delivery date: 09 Mar 2020', styles: {} },
          { text: 'Domestic', styles: {} },
        ],
      },
      {
        value: '1.033',
        label: 'Price escalation factor',
        details: [],
      },
      {
        value: '$999.99',
        label: 'Total amount requested',
        details: [{ text: '', styles: {} }],
      },
    ]);
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
        details: [{ text: 'Estimated: 8,000 lbs', styles: {} }],
      },
      {
        value: '32210',
        label: 'Mileage',
        details: [{ text: 'ZIP 322 to ZIP 919', styles: {} }],
      },
      {
        value: '13.43',
        label: 'Fuel surcharge price (per mi)',
        details: [
          { text: 'EIA diesel: $2.73', styles: {} },
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

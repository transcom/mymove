import React from 'react';
import { mount } from 'enzyme';

import ServiceItemCalculations from './ServiceItemCalculations';

const data = [
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yNzc2MTda',
    id: '6c7f1673-1ada-44fe-aa9b-e921d6e15f0e',
    key: 'EscalationCompounded',
    origin: 'PRICER',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'DECIMAL',
    value: '1.033',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yOTc2ODJa',
    id: 'b3ca0c12-fea3-4dd1-b228-30c1cc007452',
    key: 'PriceRateOrFactor',
    origin: 'PRICER',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'DECIMAL',
    value: '1.033',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zMTY5NDha',
    id: '87e77d29-d8c9-4b74-b45f-6842cd3ef970',
    key: 'ServiceAreaOrigin',
    origin: 'SYSTEM',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'STRING',
    value: '176',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zMzU1Njda',
    id: '5a993802-1504-4415-9b18-fdb1fdfd201c',
    key: 'WeightBilledActual',
    origin: 'SYSTEM',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'INTEGER',
    value: '8500',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zNTI1MDZa',
    id: 'b26fcc8f-2c06-4b00-8b51-4715a2eb0f33',
    key: 'ZipDestAddress',
    origin: 'PRIME',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'STRING',
    value: '91910',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yNDYwMDRa',
    id: 'f2a3e73f-6450-43d6-a783-181501cfab22',
    key: 'ContractCode',
    origin: 'SYSTEM',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'STRING',
    value: '1',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yNjY4M1o=',
    id: 'b4ba804d-f661-4df1-a488-11da9668647b',
    key: 'DistanceZip3',
    origin: 'SYSTEM',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'INTEGER',
    value: '210',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yODc3NDla',
    id: '83f24c0d-25ab-465a-b60b-d27bfb77b41a',
    key: 'IsPeak',
    origin: 'PRICER',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'BOOLEAN',
    value: 'FALSE',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zMDY2Nzha',
    id: '0e908b35-e61b-47c5-b4bc-f1649aa1cdc2',
    key: 'RequestedPickupDate',
    origin: 'PRIME',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'DATE',
    value: '2020-03-11',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zMjY2NDVa',
    id: '70abd9bc-afaa-4e4d-ad15-d3e55b57d2fb',
    key: 'WeightActual',
    origin: 'PRIME',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'INTEGER',
    value: '8500',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zNDQxMTda',
    id: '02438e39-de6c-4c64-b817-9932ee319a4c',
    key: 'WeightEstimated',
    origin: 'PRIME',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'INTEGER',
    value: '8000',
  },
  {
    eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zNjA5MTha',
    id: 'dcfa55b2-3106-4e1b-af4a-f19d82b5f446',
    key: 'ZipPickupAddress',
    origin: 'PRIME',
    paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
    type: 'STRING',
    value: '32210',
  },
];

const expectedOutput = [
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
    details: null,
  },
  {
    value: '$10.00',
    label: 'Total amount requested',
    details: [],
  },
];

describe('ServiceItemCalculations', () => {
  const itemCode = 'DLH';
  const totalAmount = 1000;
  const serviceItemCalculations = mount(
    <ServiceItemCalculations itemCode={itemCode} totalAmountRequested={totalAmount} />,
  );
  const siCalcLargeWithData = mount(
    <ServiceItemCalculations serviceItemParams={data} totalAmountRequested={totalAmount} itemCode={itemCode} />,
  );
  const siCalcSmallWithData = mount(
    <ServiceItemCalculations
      serviceItemParams={data}
      tableSize="small"
      totalAmountRequested={totalAmount}
      itemCode={itemCode}
    />,
  );

  it('renders without crashing', () => {
    expect(serviceItemCalculations.length).toBe(1);
  });

  describe('for item code DLH', () => {
    describe('for service item calculations large table', () => {
      it('renders large table styling by default', () => {
        const wrapper = siCalcLargeWithData.find('[data-testid="ServiceItemCalculations"]');
        expect(wrapper.hasClass('ServiceItemCalculationsSmall')).toBe(false);
      });

      it('renders correct icons', () => {
        const wrapper = siCalcLargeWithData;
        const timesIcons = wrapper.find('[icon="times"]');
        const equalsIcons = wrapper.find('[icon="equals"]');

        expect(timesIcons.length).toBe(3);
        expect(equalsIcons.length).toBe(1);
      });

      it('renders correct data', () => {
        const wrapper = siCalcLargeWithData.find('[data-testid="column"]');
        expectedOutput.forEach((obj, index) => {
          expect(wrapper.at(index).find('[data-testid="value"]').text()).toBe(obj.value);
          expect(wrapper.at(index).find('[data-testid="label"]').text()).toBe(obj.label);
          expect(wrapper.at(index).find('[data-testid="details"]').text()).toBe(
            obj.details ? obj.details.join('') : '',
          );
        });
      });
    });

    describe('for service item calculations small table', () => {
      it('renders small table styling', () => {
        expect(
          siCalcSmallWithData.find('[data-testid="ServiceItemCalculations"]').hasClass('ServiceItemCalculationsSmall'),
        ).toBe(true);
      });

      it('renders no icons', () => {
        const wrapper = siCalcSmallWithData;
        const timesIcons = wrapper.find('[icon="times"]');
        const equalsIcons = wrapper.find('[icon="equals"]');

        expect(timesIcons.length).toBe(0);
        expect(equalsIcons.length).toBe(0);
      });

      it('renders correct data', () => {
        const wrapper = siCalcSmallWithData.find('[data-testid="column"]');

        expectedOutput.forEach((obj, index) => {
          expect(wrapper.at(index).find('[data-testid="value"]').text()).toBe(obj.value);
          expect(wrapper.at(index).find('[data-testid="label"]').text()).toBe(obj.label);
          expect(wrapper.at(index).find('[data-testid="details"]').text()).toBe(
            obj.details ? obj.details.join('') : '',
          );
        });
      });
    });
  });
});

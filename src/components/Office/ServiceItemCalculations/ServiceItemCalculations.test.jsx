import React from 'react';
import { mount } from 'enzyme';

import ServiceItemCalculations from './ServiceItemCalculations';
import testParams from './serviceItemTestParams';

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
    <ServiceItemCalculations
      serviceItemParams={testParams.DLHparams}
      totalAmountRequested={totalAmount}
      itemCode={itemCode}
    />,
  );
  const siCalcSmallWithData = mount(
    <ServiceItemCalculations
      serviceItemParams={testParams.DLHparams}
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

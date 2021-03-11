import React from 'react';
import { mount } from 'enzyme';

import ServiceItemCalculations from './ServiceItemCalculations';

const data = [
  {
    value: '85 cwt',
    label: 'Billable weight (cwt)',
    details: ['Shipment weight: 8,500 lbs', 'Estimated: 8,000'],
  },
  {
    value: '2,337',
    label: 'Mileage',
    details: ['Zip 322 to Zip 919'],
  },
  {
    value: '0.03',
    label: 'Baseline linehaul price',
    details: ['Domestic non-peak', 'Origin service area: 176', 'Pickup date: 24 Jan 2020'],
  },
  {
    value: '1.033',
    label: 'Price escalation factor',
    details: null,
  },
  {
    value: '$6.423',
    label: 'Total amount requested',
    details: [],
  },
];

describe('ServiceItemCalculations', () => {
  const serviceItemCalculations = mount(<ServiceItemCalculations calculations={[]} />);
  const siCalcLargeWithData = mount(<ServiceItemCalculations calculations={data} />);
  const siCalcSmallWIthData = mount(<ServiceItemCalculations calculations={data} tableSize="small" />);

  it('renders without crashing', () => {
    expect(serviceItemCalculations.length).toBe(1);
  });

  describe('for service item calculations large table', () => {
    it('renders large table styling by default', () => {
      const wrapper = serviceItemCalculations.find('[data-testid="ServiceItemCalculations"]');
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

      data.forEach((obj, index) => {
        expect(wrapper.at(index).find('[data-testid="value"]').text()).toBe(obj.value);
        expect(wrapper.at(index).find('[data-testid="label"]').text()).toBe(obj.label);
        expect(wrapper.at(index).find('[data-testid="details"]').text()).toBe(obj.details ? obj.details.join(' ') : '');
      });
    });
  });

  describe('for service item calculations small table', () => {
    it('renders small table styling', () => {
      expect(
        siCalcSmallWIthData.find('[data-testid="ServiceItemCalculations"]').hasClass('ServiceItemCalculationsSmall'),
      ).toBe(true);
    });

    it('renders no icons', () => {
      const wrapper = siCalcSmallWIthData;
      const timesIcons = wrapper.find('[icon="times"]');
      const equalsIcons = wrapper.find('[icon="equals"]');

      expect(timesIcons.length).toBe(0);
      expect(equalsIcons.length).toBe(0);
    });

    it('renders correct data', () => {
      const wrapper = siCalcSmallWIthData.find('[data-testid="column"]');

      data.forEach((obj, index) => {
        expect(wrapper.at(index).find('[data-testid="value"]').text()).toBe(obj.value);
        expect(wrapper.at(index).find('[data-testid="label"]').text()).toBe(obj.label);
        expect(wrapper.at(index).find('[data-testid="details"]').text()).toBe(obj.details ? obj.details.join(' ') : '');
      });
    });
  });
});

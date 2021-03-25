import React from 'react';
import { mount } from 'enzyme';

import ServiceItemCalculations from './ServiceItemCalculations';
import testParams from './serviceItemTestParams';

import { SERVICE_ITEM_CODES } from 'constants/serviceItems';

// helper test function that helps test service item calculations based on code
const testServiceItemCalculation = (serviceItemCodeToTest, data, expectedOutput) => {
  const totalAmount = 1000;

  const mountedComponent = mount(
    <ServiceItemCalculations
      serviceItemParams={data}
      totalAmountRequested={totalAmount}
      itemCode={serviceItemCodeToTest}
    />,
  );

  describe(`item code ${serviceItemCodeToTest}`, () => {
    it('renders correct data', () => {
      const wrapper = mountedComponent.find('[data-testid="column"]');
      expectedOutput.forEach((obj, index) => {
        expect(wrapper.at(index).find('[data-testid="value"]').text()).toBe(obj.value);
        expect(wrapper.at(index).find('[data-testid="label"]').text()).toBe(obj.label);
        expect(wrapper.at(index).find('[data-testid="details"]').text()).toBe(obj.details ? obj.details.join('') : '');
      });
    });
  });
};

describe('ServiceItemCalculations', () => {
  const itemCode = 'DLH';
  const totalAmount = 1000;
  const serviceItemCalculationsLarge = mount(
    <ServiceItemCalculations
      itemCode={itemCode}
      totalAmountRequested={totalAmount}
      serviceItemParams={testParams.DomesticLongHaul}
    />,
  );
  const serviceItemCalculationsSmall = mount(
    <ServiceItemCalculations
      itemCode={itemCode}
      totalAmountRequested={totalAmount}
      serviceItemParams={testParams.DomesticLongHaul}
      tableSize="small"
    />,
  );

  it('renders without crashing', () => {
    expect(serviceItemCalculationsLarge.length).toBe(1);
  });

  describe('large table', () => {
    it('renders correct classnames by default', () => {
      const wrapper = serviceItemCalculationsLarge.find('[data-testid="ServiceItemCalculations"]');
      expect(wrapper.hasClass('ServiceItemCalculationsSmall')).toBe(false);
    });

    it('renders icons', () => {
      const wrapper = serviceItemCalculationsLarge;
      const timesIcons = wrapper.find('[icon="times"]');
      const equalsIcons = wrapper.find('[icon="equals"]');

      expect(timesIcons.length).toBe(3);
      expect(equalsIcons.length).toBe(1);
    });
  });

  describe('small table', () => {
    it('renders correct classnames', () => {
      expect(
        serviceItemCalculationsSmall
          .find('[data-testid="ServiceItemCalculations"]')
          .hasClass('ServiceItemCalculationsSmall'),
      ).toBe(true);
    });

    it('renders no icons', () => {
      const wrapper = serviceItemCalculationsSmall;
      const timesIcons = wrapper.find('[icon="times"]');
      const equalsIcons = wrapper.find('[icon="equals"]');

      expect(timesIcons.length).toBe(0);
      expect(equalsIcons.length).toBe(0);
    });
  });

  let expectedOutput = [
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
  testServiceItemCalculation(SERVICE_ITEM_CODES.DLH, testParams.DomesticLongHaul, expectedOutput);

  expectedOutput = [
    {
      value: '85 cwt',
      label: 'Billable weight (cwt)',
      details: ['Shipment weight: 8,500 lbs', 'Estimated: '],
    },
    {
      value: '210',
      label: 'Mileage',
      details: ['Zip 210 to Zip 910'],
    },
    {
      value: '0.09',
      label: 'Fuel surcharge price (per mi)',
      details: ['EIA diesel: $2.73', 'Weight-based distance multiplier: 0.000417', 'Pickup date: 11 Mar 2020'],
    },
    {
      value: '$10.00',
      label: 'Total amount requested',
      details: [''],
    },
  ];
  testServiceItemCalculation(SERVICE_ITEM_CODES.FSC, testParams.FuelSurchage, expectedOutput);
});

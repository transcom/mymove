import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import ServiceItemCalculations from './ServiceItemCalculations';
import testParams from './serviceItemTestParams';

import { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { SHIPMENT_OPTIONS } from 'shared/constants';

// helper test function that helps test service item calculations based on code
const testServiceItemCalculation = (testData) => {
  const [serviceItemCodeToTest, data, additionalData, expectedOutput] = testData;
  const totalAmount = 1000;

  const mountedComponent = mount(
    <ServiceItemCalculations
      serviceItemParams={data}
      additionalServiceItemData={additionalData}
      totalAmountRequested={totalAmount}
      itemCode={serviceItemCodeToTest}
    />,
  );

  const mountedComponentAdditionalData = mount(
    <ServiceItemCalculations
      serviceItemParams={data}
      additionalServiceItemData={additionalData}
      totalAmountRequested={totalAmount}
      itemCode={serviceItemCodeToTest}
    />,
  );

  describe(`item code ${serviceItemCodeToTest}`, () => {
    it('renders correct data', () => {
      const wrapper = additionalData
        ? mountedComponentAdditionalData.find('[data-testid="ServiceItemCalculations"]')
        : mountedComponent.find('[data-testid="ServiceItemCalculations"]');

      expectedOutput.forEach((obj) => {
        expect(wrapper.text()).toContain(obj.label);
        expect(wrapper.text()).toContain(obj.value);
      });
    });
  });
};

describe('ServiceItemCalculations DLH', () => {
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

  const expectedOutput = [
    {
      value: '85 cwt',
      label: 'Billable weight (cwt)',
    },
    {
      value: '210',
      label: 'Mileage',
    },
    {
      value: '1.71',
      label: 'Baseline linehaul price',
    },
    {
      value: '1.033',
      label: 'Price escalation factor',
    },
    {
      value: '$10.00',
      label: 'Total:',
    },
  ];
  testServiceItemCalculation([SERVICE_ITEM_CODES.DLH, testParams.DomesticLongHaul, {}, expectedOutput]);
});

describe('ServiceItemCalculations DCRT', () => {
  const itemCode = SERVICE_ITEM_CODES.DCRT;
  const totalAmount = 1000;
  const serviceItemCalculationsLarge = mount(
    <ServiceItemCalculations
      itemCode={itemCode}
      totalAmountRequested={totalAmount}
      serviceItemParams={testParams.DomesticCrating}
      additionalServiceItemData={testParams.additionalCratingDataDCRT}
    />,
  );
  const serviceItemCalculationsSmall = mount(
    <ServiceItemCalculations
      itemCode={itemCode}
      totalAmountRequested={totalAmount}
      serviceItemParams={testParams.DomesticCrating}
      additionalServiceItemData={testParams.additionalCratingDataDCRT}
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

      expect(timesIcons.length).toBe(2);
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

  const expectedOutput = [
    {
      value: '4.00',
      label: 'Crating size (cu ft)',
    },
    {
      value: '1.71',
      label: 'Crating price (per cu ft)',
    },
    {
      value: '1.033',
      label: 'Price escalation factor',
    },
    {
      value: '$10.00',
      label: 'Total:',
    },
  ];
  testServiceItemCalculation([
    itemCode,
    testParams.DomesticCrating,
    testParams.additionalCratingDataDCRT,
    expectedOutput,
  ]);
});

describe('ServiceItemCalculations DUCRT', () => {
  const itemCode = SERVICE_ITEM_CODES.DUCRT;
  const totalAmount = 1000;
  const serviceItemCalculationsLarge = mount(
    <ServiceItemCalculations
      itemCode={itemCode}
      totalAmountRequested={totalAmount}
      serviceItemParams={testParams.DomesticUncrating}
      additionalServiceItemData={testParams.additionalCratingDataDCRT}
    />,
  );
  const serviceItemCalculationsSmall = mount(
    <ServiceItemCalculations
      itemCode={itemCode}
      totalAmountRequested={totalAmount}
      serviceItemParams={testParams.DomesticUncrating}
      additionalServiceItemData={testParams.additionalCratingDataDCRT}
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

      expect(timesIcons.length).toBe(2);
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

  const expectedOutput = [
    {
      value: '4.00',
      label: 'Crating size (cu ft)',
    },
    {
      value: '1.71',
      label: 'Uncrating price (per cu ft)',
    },
    {
      value: '1.033',
      label: 'Price escalation factor',
    },
    {
      value: '$10.00',
      label: 'Total',
    },
  ];
  testServiceItemCalculation([
    itemCode,
    testParams.DomesticUncrating,
    testParams.additionalCratingDataDCRT,
    expectedOutput,
  ]);
});

describe('shipmentType prop can affect labels', () => {
  it("shows 'Requested pickup' for HHG", () => {
    render(
      <ServiceItemCalculations
        serviceItemParams={testParams.DomesticLongHaul}
        totalAmountRequested={642}
        itemCode={SERVICE_ITEM_CODES.DLH}
        shipmentType={SHIPMENT_OPTIONS.HHG}
      />,
    );

    expect(screen.getByText('Requested pickup: 09 Mar 2020')).toBeInTheDocument();
  });

  it("shows 'Actual pickup' for NTS-release", () => {
    render(
      <ServiceItemCalculations
        serviceItemParams={testParams.DomesticLongHaul}
        totalAmountRequested={642}
        itemCode={SERVICE_ITEM_CODES.DLH}
        shipmentType={SHIPMENT_OPTIONS.NTSR}
      />,
    );

    expect(screen.getByText('Actual pickup: 09 Mar 2020')).toBeInTheDocument();
  });
});

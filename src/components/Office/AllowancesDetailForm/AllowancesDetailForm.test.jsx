import React from 'react';
import { mount } from 'enzyme';
import { Formik } from 'formik';
import { IMaskInput } from 'react-imask';

import AllowancesDetailForm from './AllowancesDetailForm';

const initialValues = {
  authorizedWeight: '11000',
};

const entitlements = {
  authorizedWeight: 11000,
  dependentsAuthorized: true,
  nonTemporaryStorage: false,
  privatelyOwnedVehicle: false,
  proGearWeight: 2000,
  proGearWeightSpouse: 500,
  storageInTransit: 90,
  totalWeight: 11000,
  totalDependents: 2,
};

describe('AllowancesDetailForm', () => {
  const wrapper = mount(
    <Formik initialValues={initialValues} onSubmit={jest.fn()}>
      <form>
        <AllowancesDetailForm entitlements={entitlements} />
      </form>
    </Formik>,
  );

  it('renders the form', () => {
    expect(wrapper.find(AllowancesDetailForm).exists()).toBe(true);
  });

  it('formats weights', () => {
    expect(wrapper.find(IMaskInput).getDOMNode().value).toBe('11,000 lbs');

    // Weight allowance
    expect(wrapper.find('dd').at(0).text()).toBe('11,000 lbs');

    // Pro-gear
    expect(wrapper.find('dd').at(1).text()).toBe('2,000 lbs');

    // Spouse Pro-gear
    expect(wrapper.find('dd').at(2).text()).toBe('500 lbs');
  });

  it('formats days in transit', () => {
    // Storage in-transit
    expect(wrapper.find('dd').at(3).text()).toBe('90 days');
  });

  it('uses defaults for undefined values', () => {
    const wrapperNoProps = mount(
      <Formik initialValues={{ authorizedWeight: '0' }} onSubmit={jest.fn()}>
        <form>
          <AllowancesDetailForm entitlements={{}} />
        </form>
      </Formik>,
    );

    // Authorized weight input
    expect(wrapperNoProps.find(IMaskInput).getDOMNode().value).toBe('0 lbs');

    // Weight allowance
    expect(wrapperNoProps.find('dd').at(0).text()).toBe('0 lbs');

    // Pro-gear
    expect(wrapperNoProps.find('dd').at(1).text()).toBe('0 lbs');

    // Spouse Pro-gear
    expect(wrapperNoProps.find('dd').at(2).text()).toBe('0 lbs');

    // Storage in-transit
    expect(wrapperNoProps.find('dd').at(3).text()).toBe('0 days');
  });
});

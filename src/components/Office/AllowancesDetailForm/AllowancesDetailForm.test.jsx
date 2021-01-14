import React from 'react';
import { mount } from 'enzyme';
import { Formik } from 'formik';
import { IMaskInput } from 'react-imask';

import AllowancesDetailForm from './AllowancesDetailForm';

const initialValues = {
  authorizedWeight: '11000',
};

const rankOptions = [
  { key: 'E_1', value: 'E-1' },
  { key: 'E_2', value: 'E-2' },
  { key: 'E_3', value: 'E-3' },
  { key: 'E_4', value: 'E-4' },
  { key: 'E_5', value: 'E-5' },
  { key: 'E_6', value: 'E-6' },
  { key: 'E_7', value: 'E-7' },
  { key: 'E_8', value: 'E-8' },
  { key: 'E_9', value: 'E-9' },
];

const branchOptions = [
  { key: 'Army', value: 'Army' },
  { key: 'Navy', value: 'Navy' },
  { key: 'Marines', value: 'Marines' },
];

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
        <AllowancesDetailForm entitlements={entitlements} rankOptions={rankOptions} branchOptions={branchOptions} />
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
      <Formik initialValues={{ authorizedWeight: null }} onSubmit={jest.fn()}>
        <form>
          <AllowancesDetailForm entitlements={{}} rankOptions={[]} branchOptions={[]} />
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

    // Dependents authorized
    expect(wrapperNoProps.find(`[name="dependentsAuthorized"]`)).toBeTruthy();
  });
});

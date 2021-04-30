import React from 'react';
import { mount } from 'enzyme';
import { Formik } from 'formik';

import AllowancesDetailForm from './AllowancesDetailForm';

const initialValues = {
  authorizedWeight: '11000',
  proGearWeight: '2000',
  proGearWeightSpouse: '500',
  requiredMedicalEquipmentWeight: '1000',
  organizationalClothingAndIndividualEquipment: true,
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
  requiredMedicalEquipmentWeight: 1000,
  organizationalClothingAndIndividualEquipment: true,
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

  it('formats days in transit', () => {
    // Storage in-transit
    expect(wrapper.find('dd').at(2).text()).toBe('90 days');
  });

  it('renders the pro-gear hints', () => {
    expect(wrapper.find('[data-testid="proGearWeightHint"]').at(0).text()).toBe('Max. 2,000 lbs');
    expect(wrapper.find('[data-testid="proGearWeightSpouseHint"]').at(0).text()).toBe('Max. 500 lbs');
  });

  it('uses defaults for undefined values', () => {
    const wrapperNoProps = mount(
      <Formik initialValues={{ authorizedWeight: null }} onSubmit={jest.fn()}>
        <form>
          <AllowancesDetailForm entitlements={{}} rankOptions={[]} branchOptions={[]} />
        </form>
      </Formik>,
    );

    // Pro-gear
    expect(wrapperNoProps.find(`input[data-testid="proGearWeightInput"]`).getDOMNode().value).toBe('0');

    // Pro-gear spouse
    expect(wrapperNoProps.find(`input[data-testid="proGearWeightSpouseInput"]`).getDOMNode().value).toBe('0');

    // RME
    expect(wrapperNoProps.find(`input[data-testid="rmeInput"]`).getDOMNode().value).toBe('0');

    // Branch
    expect(wrapperNoProps.find(`select[data-testid="branchInput"]`).getDOMNode().value).toBe('');

    // Rank
    expect(wrapperNoProps.find(`select[data-testid="rankInput"]`).getDOMNode().value).toBe('');

    // OCIE
    expect(
      wrapperNoProps.find(`input[name="organizationalClothingAndIndividualEquipment"]`).getDOMNode().checked,
    ).toBeFalsy();

    // Authorized weight text only
    expect(wrapperNoProps.find('dd').at(0).text()).toBe('0 lbs');

    // Weight allowance
    expect(wrapperNoProps.find('dd').at(1).text()).toBe('0 lbs');

    // Storage in-transit
    expect(wrapperNoProps.find('dd').at(2).text()).toBe('0 days');

    // Dependents authorized
    expect(wrapperNoProps.find(`input[name="dependentsAuthorized"]`).getDOMNode().checked).toBeFalsy();
  });

  it('renders editable authorized weight input', () => {
    const wrapperNoProps = mount(
      <Formik initialValues={{ authorizedWeight: null }} onSubmit={jest.fn()}>
        <form>
          <AllowancesDetailForm entitlements={{}} rankOptions={[]} branchOptions={[]} editableAuthorizedWeight />
        </form>
      </Formik>,
    );
    // Authorized weight input
    expect(wrapperNoProps.find('input[data-testid="authorizedWeightInput"]').getDOMNode().value).toBe('0');
  });

  it('renders the title header', () => {
    const wrapperNoProps = mount(
      <Formik initialValues={{ authorizedWeight: null }} onSubmit={jest.fn()}>
        <form>
          <AllowancesDetailForm entitlements={{}} rankOptions={[]} branchOptions={[]} header="Counseling" />
        </form>
      </Formik>,
    );
    // Authorized weight input
    expect(wrapperNoProps.find('[data-testid="header"]').text()).toBe('Counseling');
  });
});

import React from 'react';
import { shallow } from 'enzyme';

import AllowancesTable from './AllowancesTable';

const info = {
  branch: 'Navy',
  rank: 'E-6',
  weightAllowance: 11000,
  authorizedWeight: 12000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
  requiredMedicalEquipmentWeight: 1000,
  organizationalClothingAndIndividualEquipment: true,
};

describe('Allowances Table', () => {
  it('should render the data passed to its props', () => {
    const wrapper = shallow(<AllowancesTable info={info} />);
    expect(wrapper.find({ 'data-testid': 'branchRank' }).text()).toMatch(`${info.branch}, ${info.rank}`);
    expect(wrapper.find({ 'data-testid': 'weightAllowance' }).text()).toMatch('11,000 lbs');
    expect(wrapper.find({ 'data-testid': 'authorizedWeight' }).text()).toMatch('12,000 lbs');
    expect(wrapper.find({ 'data-testid': 'progear' }).text()).toMatch('2,000 lbs');
    expect(wrapper.find({ 'data-testid': 'spouseProgear' }).text()).toMatch('500 lbs');
    expect(wrapper.find({ 'data-testid': 'storageInTransit' }).text()).toMatch('90 days');
    expect(wrapper.find({ 'data-testid': 'dependents' }).text()).toMatch('Authorized');
    expect(wrapper.find({ 'data-testid': 'rme' }).text()).toMatch('1,000 lbs');
    expect(wrapper.find({ 'data-testid': 'ocie' }).text()).toMatch('Authorized');
  });

  it('renders the OCIE as unauthorized passed to its props', () => {
    const testInfo = { ...info, organizationalClothingAndIndividualEquipment: false };
    const wrapper = shallow(<AllowancesTable info={testInfo} />);
    expect(wrapper.find({ 'data-testid': 'ocie' }).text()).toMatch('Unauthorized');
  });

  it('should be able to show edit btn', () => {
    const wrapper = shallow(<AllowancesTable info={info} showEditBtn />);
    expect(wrapper.find('Link').text()).toMatch('Edit Allowances');
    expect(wrapper.find('Link').prop('to')).toBe('allowances');
  });

  it('should be able to hide edit btn', () => {
    const wrapper = shallow(<AllowancesTable info={info} />);
    expect(wrapper.find('Link').exists()).toBe(false);
  });
});

describe('Allowances Table when SIT is 0', () => {
  it('displays an empty string for the SIT allowance', () => {
    info.storageInTransit = 0;
    const wrapper = shallow(<AllowancesTable info={info} />);
    expect(wrapper.find({ 'data-testid': 'storageInTransit' }).text()).toEqual('');
  });
});

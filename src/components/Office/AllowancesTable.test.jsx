import React from 'react';
import { shallow } from 'enzyme';

import AllowancesTable from './AllowancesTable';

const info = {
  branch: 'Navy',
  rank: 'E-6',
  weightAllowance: 11000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
};

describe('Allowances Table', () => {
  it('should render the data passed to its props', () => {
    const wrapper = shallow(<AllowancesTable info={info} />);
    expect(wrapper.find({ 'data-testid': 'branchRank' }).text()).toMatch(`${info.branch}, ${info.rank}`);
    expect(wrapper.find({ 'data-testid': 'weightAllowance' }).text()).toMatch(`${info.weightAllowance} lbs`);
    expect(wrapper.find({ 'data-testid': 'authorizedWeight' }).text()).toMatch(`${info.authorizedWeight} lbs`);
    expect(wrapper.find({ 'data-testid': 'progear' }).text()).toMatch(`${info.progear} lbs`);
    expect(wrapper.find({ 'data-testid': 'spouseProgear' }).text()).toMatch(`${info.spouseProgear} lbs`);
    expect(wrapper.find({ 'data-testid': 'storageInTransit' }).text()).toMatch(`${info.storageInTransit} days`);
    expect(wrapper.find({ 'data-testid': 'dependents' }).text()).toMatch('Authorized');
  });
});

describe('Allowances Table when SIT is 0', () => {
  it('displays an empty string for the SIT allowance', () => {
    info.storageInTransit = 0;
    const wrapper = shallow(<AllowancesTable info={info} />);
    expect(wrapper.find({ 'data-testid': 'storageInTransit' }).text()).toEqual('');
  });
});

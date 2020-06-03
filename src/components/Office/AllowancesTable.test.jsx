import React from 'react';
import { shallow } from 'enzyme';
import AllowancesTable from './AllowancesTable';

const info = {
  branch: 'Navy',
  rank: 'E-6',
  weightAllowance: '11,000 lbs',
  authorizedWeight: '11,000 lbs',
  progear: '2,000 lbs',
  spouseProgear: '500 lbs',
  storageInTransit: '90 days',
  dependents: 'Authorized',
};

describe('Allowances Table', () => {
  it('should render the data passed to its props', () => {
    const wrapper = shallow(<AllowancesTable info={info} />);
    expect(wrapper.find({ 'data-cy': 'branchRank' }).text()).toMatch(`${info.branch}, ${info.rank}`);
    expect(wrapper.find({ 'data-cy': 'weightAllowance' }).text()).toMatch(info.weightAllowance);
    expect(wrapper.find({ 'data-cy': 'authorizedWeight' }).text()).toMatch(info.authorizedWeight);
    expect(wrapper.find({ 'data-cy': 'progear' }).text()).toMatch(info.progear);
    expect(wrapper.find({ 'data-cy': 'spouseProgear' }).text()).toMatch(info.spouseProgear);
    expect(wrapper.find({ 'data-cy': 'storageInTransit' }).text()).toMatch(info.storageInTransit);
    expect(wrapper.find({ 'data-cy': 'dependents' }).text()).toMatch(info.dependents);
  });
});

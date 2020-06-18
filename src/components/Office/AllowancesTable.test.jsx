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
    expect(wrapper.find({ 'data-cy': 'branchRank' }).text()).toMatch(`${info.branch}, ${info.rank}`);
    expect(wrapper.find({ 'data-cy': 'weightAllowance' }).text()).toMatch(`${info.weightAllowance} lbs`);
    expect(wrapper.find({ 'data-cy': 'authorizedWeight' }).text()).toMatch(`${info.authorizedWeight} lbs`);
    expect(wrapper.find({ 'data-cy': 'progear' }).text()).toMatch(`${info.progear} lbs`);
    expect(wrapper.find({ 'data-cy': 'spouseProgear' }).text()).toMatch(`${info.spouseProgear} lbs`);
    expect(wrapper.find({ 'data-cy': 'storageInTransit' }).text()).toMatch(`${info.storageInTransit} days`);
    expect(wrapper.find({ 'data-cy': 'dependents' }).text()).toMatch('Authorized');
  });
});

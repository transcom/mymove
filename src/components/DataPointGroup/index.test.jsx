import React from 'react';
import { shallow } from 'enzyme';

import DataPointGroup from '.';

import DataPoint from 'components/DataPoint';
import { ReactComponent as ArrowRight } from 'shared/icon/arrow-right.svg';

describe('DataPair', () => {
  it('renders a single data point child', () => {
    const headers = ['column 1', 'column 2'];
    const row = ['cell 1', 'cell 2'];
    const wrapper = shallow(
      <DataPointGroup>
        <DataPoint columnHeaders={[headers]} dataRow={[row]} />
      </DataPointGroup>,
    );
    expect(wrapper.find(DataPoint).length).toBe(1);
  });

  it('renders multiple data point children in container', () => {
    const headers = ['column 1', 'column 2'];
    const row = ['cell 1', 'cell 2'];
    const wrapper = shallow(
      <DataPointGroup>
        <DataPoint columnHeaders={[headers]} dataRow={[row]} Icon={ArrowRight} />
        <DataPoint columnHeaders={[headers]} dataRow={[row]} Icon={ArrowRight} />
      </DataPointGroup>,
    );
    expect(wrapper.find(DataPoint).length).toBe(2);
  });
});

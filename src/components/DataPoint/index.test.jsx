import React from 'react';
import { shallow, mount } from 'enzyme';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import DataPoint from '.';

describe('DataPoint', () => {
  it('renders with column header and data row', () => {
    const header = 'This is a datapoint header.';
    const row = 'This is a datapoint row.';
    const wrapper = shallow(<DataPoint columnHeaders={[header]} dataRow={[row]} />);
    expect(wrapper.find('th').text()).toContain(header);
    expect(wrapper.find('td').text()).toContain(row);
  });

  it('renders with an icon', () => {
    const headers = ['column 1', 'column 2'];
    const row = ['cell 1', 'cell 2'];
    const wrapper = mount(
      <DataPoint columnHeaders={headers} dataRow={row} icon={<FontAwesomeIcon icon="arrow-right" />} />,
    );

    expect(wrapper.find('th').at(0).text()).toContain('column 1');
    expect(wrapper.find('th').at(1).text()).toContain('column 2');
    expect(wrapper.find('td').at(0).text()).toContain('cell 1');
    expect(wrapper.find('FontAwesomeIcon').prop('icon')).toBe('arrow-right');
    expect(wrapper.find('td').at(1).text()).toContain('cell 2');
  });
});

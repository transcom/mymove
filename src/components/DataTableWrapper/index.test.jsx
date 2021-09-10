import React from 'react';
import { shallow } from 'enzyme';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import DataTableWrapper from '.';

import DataTable from 'components/DataTable';

describe('DataPair', () => {
  it('renders a single data point child', () => {
    const headers = ['column 1', 'column 2'];
    const row = ['cell 1', 'cell 2'];
    const wrapper = shallow(
      <DataTableWrapper>
        <DataTable columnHeaders={[headers]} dataRow={[row]} />
      </DataTableWrapper>,
    );
    expect(wrapper.find(DataTable).length).toBe(1);
  });

  it('renders multiple data point children in container', () => {
    const headers = ['column 1', 'column 2'];
    const row = ['cell 1', 'cell 2'];
    const wrapper = shallow(
      <DataTableWrapper>
        <DataTable columnHeaders={[headers]} dataRow={[row]} icon={<FontAwesomeIcon icon="arrow-right" />} />
        <DataTable columnHeaders={[headers]} dataRow={[row]} icon={<FontAwesomeIcon icon="arrow-right" />} />
      </DataTableWrapper>,
    );
    expect(wrapper.find(DataTable).length).toBe(2);
  });
});

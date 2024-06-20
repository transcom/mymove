/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import TableQueue from './TableQueue';
import TableCSVExportButton from './TableCSVExportButton';

describe('TableQueue - react table', () => {
  const defaultProps = {
    title: 'TableQueue test',
    useQueries: jest.fn(() => ({ queueResult: {} })),
    handleClick: jest.fn(),
    columns: [],
  };

  const exportButtonProps = {
    title: 'TableQueue With Export',
    useQueries: jest.fn(() => ({ queueResult: {} })),
    handleClick: jest.fn(),
    columns: [],
    showCSVExport: true,
    csvExportQueueFetcher: jest.fn(),
    csvExportQueueFetcherKey: 'moves',
  };

  it('renders without crashing', () => {
    const wrapper = mount(<TableQueue {...defaultProps} />);
    expect(wrapper.find(TableQueue).length).toBe(1);
    expect(wrapper.find(TableCSVExportButton).length).toBe(0);
  });

  it('renders the CSV export button', () => {
    const wrapper = mount(<TableQueue {...exportButtonProps} />);
    expect(wrapper.find(TableCSVExportButton).length).toBe(1);
  });
});

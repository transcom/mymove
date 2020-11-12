/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import TableQueue from './TableQueue';

describe('TableQueue - react table', () => {
  const defaultProps = {
    title: 'TableQueue test',
    useQueries: jest.fn(() => ({ queueResult: {} })),
    handleClick: jest.fn(),
    columns: [],
  };

  it('renders without crashing', () => {
    const wrapper = mount(<TableQueue {...defaultProps} />);
    expect(wrapper.find(TableQueue).length).toBe(1);
  });
});

import React from 'react';
import { mount } from 'enzyme';

import AdminPagination from './AdminPagination';

describe('AdminPagination', () => {
  it('Returns "no results" text when the filter has no matches', () => {
    const props = {
      total: 0,
      isLoading: false,
      count: 0,
      rowsPerPage: 25,
    };
    const pagination = mount(<AdminPagination {...props} />);
    expect(pagination.exists('.no-results')).toEqual(true);
    pagination.setProps({
      total: 1,
    });
    expect(pagination.exists('.no-results')).toEqual(false);
  });
  it('Does not return "no results" text when the component is still loading', () => {
    const props = {
      total: 0,
      isLoading: true,
      count: 0,
      rowsPerPage: 25,
    };
    const pagination = mount(<AdminPagination {...props} />);
    expect(pagination.exists('.no-results')).toEqual(false);
  });
});

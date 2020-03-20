import React from 'react';
import { shallow } from 'enzyme';

import AdminPagination from './AdminPagination';

describe('AdminPagination', () => {
  it('Returns "no results" text when the filter has no matches', () => {
    const props = {
      total: 0,
      isLoading: false,
      count: 0,
      rowsPerPage: 25,
    };
    const pagination = shallow(<AdminPagination {...props} />);
    const noResultsDiv = pagination.find('.no-results');
    expect(noResultsDiv).toHaveLength(1);
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
    const pagination = shallow(<AdminPagination {...props} />);
    // const noResultsDiv = pagination.find('.no-results');
    expect(pagination.find('.no-results')).toHaveLength(0);
  });
});

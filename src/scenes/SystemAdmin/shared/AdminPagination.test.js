import React from 'react';
import { shallow } from 'enzyme';

import AdminPagination from './AdminPagination';

import { useListContext } from 'react-admin';

jest.mock('react-admin', () => ({
  ...jest.requireActual('react-admin'),
  useListContext: jest.fn(),
}));

beforeEach(() => {
  jest.resetAllMocks();
});

describe('AdminPagination', () => {
  it('Returns "no results" text when the filter has no matches', () => {
    useListContext.mockImplementation(() => ({
      total: 0,
      isLoading: false,
    }));

    const pagination = shallow(<AdminPagination />);

    const noResultsDiv = pagination.find('.no-results');
    expect(noResultsDiv).toHaveLength(1);
  });

  it('Does not return "no results" text when the component is still loading', () => {
    useListContext.mockImplementation(() => ({
      total: 0,
      isLoading: true,
    }));
    const pagination = shallow(<AdminPagination />);
    // const noResultsDiv = pagination.find('.no-results');
    expect(pagination.find('.no-results')).toHaveLength(0);
  });
});

/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import Table from './Table';

describe('React table', () => {
  const defaultProps = {
    getTableProps: jest.fn(() => ({})),
    getTableBodyProps: jest.fn(() => ({})),
    prepareRow: jest.fn(),
    headerGroups: [],
    rows: [],
  };

  const createTableComponent = (columns, data, handleClick, showPagination) => {
    const rows = data.map((row, index) => {
      const cells = Object.keys(row).map((v, i) => {
        return {
          key: `cells${i}`,
          column: v.accessor,
          row: v.accessor,
          getCellProps: jest.fn(),
        };
      });
      return { key: `rows${index}`, cells, getRowProps: jest.fn() };
    });

    const headers = columns.map((column, index) => ({
      key: `headers${index}`,
      id: column.accessor,
      render: (val) => val,
      getHeaderProps: jest.fn(),
    }));
    return mount(
      <Table
        {...defaultProps}
        rows={rows}
        headerGroups={[
          {
            headers,
            getHeaderGroupProps: () => ({ key: 'hi' }),
          },
        ]}
        handleClick={handleClick}
        showPagination={showPagination}
      />,
    );
  };

  it('renders without crashing', () => {
    const wrapper = createTableComponent([], []);
    expect(wrapper.find('[data-testid="react-table"]').length).toBe(1);
  });

  it('renders a table with pagination', () => {
    const wrapper = createTableComponent([], [], undefined, true);
    expect(wrapper.find('[data-testid="pagination"]').length).toBe(1);
  });
});

/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import { createHeader, OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID } from './utils';
import TableQueue from './TableQueue';

import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS } from 'constants/queues';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';

const localStorageMock = (() => {
  let store = {};

  return {
    getItem(key) {
      return store[key] || null;
    },
    setItem(key, value) {
      store[key] = value.toString();
    },
    removeItem(key) {
      delete store[key];
    },
    clear() {
      store = {};
    },
  };
})();

Object.defineProperty(window, 'sessionStorage', {
  value: localStorageMock,
});

beforeEach(() => {
  jest.restoreAllMocks();
  window.sessionStorage.clear();
});

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

describe('SessionStorage TableQueue - react table', () => {
  const data = [
    {
      col1: 'Banks, Aaliyah',
      col2: '987654321',
      col3: 'New move',
      col4: 'LCKMAJ',
      col5: 'Navy',
      col6: '3',
      col7: 'NAS Jacksonville',
      col8: 'HAFC',
      col9: 'Garimundi, J (SW)',
    },
    {
      col1: 'Childers, Jamie',
      col2: '987654321',
      col3: 'New move',
      col4: 'XCQ5ZH',
      col5: 'Navy',
      col6: '3',
      col7: 'NAS Jacksonville',
      col8: 'HAFC',
      col9: 'Garimundi, J (SW)',
    },
    {
      col1: 'Clark-Nunez, Sofia',
      col2: '987654321',
      col3: 'New move',
      col4: 'UCAF8Q',
      col5: 'Navy',
      col6: '3',
      col7: 'NAS Jacksonville',
      col8: 'HAFC',
      col9: 'Garimundi, J (SW)',
    },
  ];

  const columns = (isFilterable = false) => [
    createHeader('Customer name', 'col1', { isFilterable }),
    createHeader('DoD ID', 'col2', { isFilterable }),
    createHeader('Status', 'col3', {
      isFilterable,
      Filter: (props) => <MultiSelectCheckBoxFilter options={MOVE_STATUS_OPTIONS} {...props} />,
    }),
    createHeader('Move Code', 'col4', { isFilterable }),
    createHeader('Branch', 'col5', {
      isFilterable,
      Filter: (props) => <SelectFilter options={BRANCH_OPTIONS} {...props} />,
    }),
    createHeader('# of shipments', 'col6', { isFilterable }),
    createHeader('Destination duty location', 'col7', { isFilterable }),
    createHeader('Origin GBLOC', 'col8', { isFilterable }),
    createHeader('Last modified by', 'col9', { isFilterable, Filter: DateSelectFilter }),
  ];

  const defaultProps = {
    title: 'Table queue',
    useQueries: () => ({ queueResult: { data, totalCount: data.length, perPage: 1 } }),
    handleClick: () => {},
    columns: columns(true),
    sessionStorageKey: 'test',
  };

  it('default item is persisted in sessionStorage', () => {
    const getItemSpy = jest.spyOn(window.sessionStorage, 'getItem');
    const setItemSpy = jest.spyOn(window.sessionStorage, 'setItem');

    const wrapper = mount(<TableQueue {...defaultProps} />);
    expect(wrapper.find(TableQueue).length).toBe(1);
    expect(setItemSpy).toBeCalledWith(
      OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID,
      '{"test":{"filters":[],"sortParam":null,"page":null,"pageSize":null}}',
    );
    expect(getItemSpy).toBeCalledWith(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  });
});

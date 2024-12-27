/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { mount } from 'enzyme';

import {
  createHeader,
  OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID,
  TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT,
  setTableQueueFilterSessionStorageValue,
  getTableQueueFilterSessionStorageValue,
} from './utils';
import TableQueue from './TableQueue';
import TableCSVExportButton from './TableCSVExportButton';

import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS } from 'constants/queues';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import DateSelectFilter from 'components/Table/Filters/DateSelectFilter';
import { MockProviders } from 'testUtils';

const localStorageMock = (() => {
  let store = {};

  return {
    getItem(key) {
      return store[key] || null;
    },
    setItem(key, value) {
      // store[key] = value.toString();
      store[key] = value;
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
    const wrapper = mount(
      <MockProviders>
        <TableQueue {...defaultProps} />
      </MockProviders>,
    );
    expect(wrapper.find(TableQueue).length).toBe(1);
    expect(wrapper.find(TableCSVExportButton).length).toBe(0);
  });

  it('renders the CSV export button', () => {
    const wrapper = mount(
      <MockProviders>
        <TableQueue {...exportButtonProps} />
      </MockProviders>,
    );
    expect(wrapper.find(TableCSVExportButton).length).toBe(1);
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
    createHeader('Customer name', 'col1', { id: 'col1', isFilterable }),
    createHeader('DoD ID', 'col2', { isFilterable }),
    createHeader('Status', 'col3', {
      id: 'col3',
      isFilterable,
      Filter: (props) => <MultiSelectCheckBoxFilter options={MOVE_STATUS_OPTIONS} {...props} />,
    }),
    createHeader('Move Code', 'col4', { isFilterable }),
    createHeader('Branch', 'col5', {
      id: 'col5',
      isFilterable,
      Filter: (props) => <SelectFilter options={BRANCH_OPTIONS} {...props} />,
    }),
    createHeader('# of shipments', 'col6', { isFilterable }),
    createHeader('Destination duty location', 'col7', { isFilterable }),
    createHeader('Origin GBLOC', 'col8', { isFilterable }),
    createHeader('Last modified by', 'col9', { isFilterable, Filter: DateSelectFilter }),
  ];

  const testSessionStorageKey = 'default';

  const defaultProps = {
    title: 'Table queue',
    useQueries: () => ({ queueResult: { data, totalCount: data.length, perPage: 1 } }),
    handleClick: () => {},
    columns: columns(true),
    sessionStorageKey: testSessionStorageKey,
  };

  it('default item is persisted in sessionStorage', () => {
    const getItemSpy = jest.spyOn(window.sessionStorage, 'getItem');
    const setItemSpy = jest.spyOn(window.sessionStorage, 'setItem');

    const wrapper = mount(
      <MockProviders>
        <TableQueue {...defaultProps} />
      </MockProviders>,
    );
    expect(wrapper.find(TableQueue).length).toBe(1);
    expect(setItemSpy).toBeCalledWith(
      OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID,
      `{"${testSessionStorageKey}":{"filters":[],"sortParam":null,"page":null,"pageSize":null}}`,
    );
    expect(getItemSpy).toBeCalledWith(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  });

  it('filter shows delete "All" pill button and two pill buttons for status MultiSelectCheckBoxFilter', () => {
    const json = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    json.filters = [{ id: 'col3', value: `${MOVE_STATUS_OPTIONS[0].value},${MOVE_STATUS_OPTIONS[2].value}` }];
    json.sortParam = [{ id: 'col2', desc: false }];
    const cache = {};
    cache[testSessionStorageKey] = json;
    window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(cache));
    const wrapper = mount(
      <MockProviders>
        <TableQueue {...defaultProps} />
      </MockProviders>,
    );
    expect(wrapper.find('button[data-testid="remove-filters-all"]').length).toBe(1);
    expect(wrapper.find(`button[data-testid="remove-filters-col3-${MOVE_STATUS_OPTIONS[0].value}"]`).length).toBe(1);
    expect(wrapper.find(`button[data-testid="remove-filters-col3-${MOVE_STATUS_OPTIONS[0].value}"]`).text()).toContain(
      `Status (${MOVE_STATUS_OPTIONS[0].label}) ×`,
    );
    expect(wrapper.find(`button[data-testid="remove-filters-col3-${MOVE_STATUS_OPTIONS[2].value}"]`).length).toBe(1);
    expect(wrapper.find(`button[data-testid="remove-filters-col3-${MOVE_STATUS_OPTIONS[2].value}"]`).text()).toContain(
      `Status (${MOVE_STATUS_OPTIONS[2].label}) ×`,
    );
  });

  it('one MultiSelectCheckBoxFilter status filter pill button', () => {
    const json = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    json.filters = [{ id: 'col3', value: `${MOVE_STATUS_OPTIONS[0].value}` }];
    json.sortParam = [{ id: 'col2', desc: false }];
    const cache = {};
    cache[testSessionStorageKey] = json;
    window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(cache));
    const wrapper = mount(
      <MockProviders>
        <TableQueue {...defaultProps} />
      </MockProviders>,
    );
    expect(wrapper.find('button[data-testid="remove-filters-all"]').length).toBe(0);
    expect(wrapper.find(`button[data-testid="remove-filters-col3-${MOVE_STATUS_OPTIONS[0].value}"]`).length).toBe(1);
    expect(wrapper.find(`button[data-testid="remove-filters-col3-${MOVE_STATUS_OPTIONS[0].value}"]`).text()).toContain(
      'Status',
    );
  });

  it('one SelectFilter filter pill button', () => {
    const json = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    json.filters = [{ id: 'col5', value: 'SPACE_FORCE' }];
    json.sortParam = [{ id: 'col2', desc: false }];
    const cache = {};
    cache[testSessionStorageKey] = json;
    window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(cache));
    const wrapper = mount(
      <MockProviders>
        <TableQueue {...defaultProps} />
      </MockProviders>,
    );
    expect(wrapper.find('button[data-testid="remove-filters-all"]').length).toBe(0);
    expect(wrapper.find(`button[data-testid="remove-filters-col5-SPACE_FORCE"]`).length).toBe(1);
    expect(wrapper.find(`button[data-testid="remove-filters-col5-SPACE_FORCE"]`).text()).toContain('Branch');
  });

  it('one TextInput filter pill button', () => {
    const json = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    json.filters = [{ id: 'col1', value: 'Foobar' }];
    json.sortParam = [{ id: 'col2', desc: false }];
    const cache = {};
    cache[testSessionStorageKey] = json;
    window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(cache));
    const wrapper = mount(
      <MockProviders>
        <TableQueue {...defaultProps} />
      </MockProviders>,
    );
    expect(wrapper.find('button[data-testid="remove-filters-all"]').length).toBe(0);
    expect(wrapper.find(`button[data-testid="remove-filters-col1"]`).length).toBe(1);
    expect(wrapper.find(`button[data-testid="remove-filters-col1"]`).text()).toContain('Customer name ×');
  });

  it('delete filter via pill button click', async () => {
    window.sessionStorage.clear();
    const filters = [{ id: 'col1', value: 'Foobar' }];
    setTableQueueFilterSessionStorageValue(testSessionStorageKey, filters);
    render(
      <MockProviders>
        <TableQueue {...defaultProps} />
      </MockProviders>,
    );
    const pillButton = await screen.findByTestId('remove-filters-col1');

    await new Promise((r) => {
      setTimeout(r, 1100);
    });

    await userEvent.click(pillButton);
    await waitFor(() => {
      const cacheItem = getTableQueueFilterSessionStorageValue(testSessionStorageKey);
      expect(cacheItem.length).toBe(0);
    });
  });

  it('test delete one of filters. verify ALL button is removed and one filter remains', async () => {
    const filters = [{ id: 'col3', value: `${MOVE_STATUS_OPTIONS[0].value},${MOVE_STATUS_OPTIONS[2].value}` }];
    setTableQueueFilterSessionStorageValue(testSessionStorageKey, filters);
    render(
      <MockProviders>
        <TableQueue {...defaultProps} />
      </MockProviders>,
    );
    const allFilterPillButton = await screen.findByTestId('remove-filters-all');
    expect(allFilterPillButton).toBeInTheDocument();

    const filter1PillButton = await screen.findByTestId(`remove-filters-col3-${MOVE_STATUS_OPTIONS[0].value}`);
    expect(filter1PillButton).toBeInTheDocument();

    const filter2PillButton = await screen.findByTestId(`remove-filters-col3-${MOVE_STATUS_OPTIONS[2].value}`);
    expect(filter2PillButton).toBeInTheDocument();

    await new Promise((r) => {
      setTimeout(r, 1100);
    });

    // delete filter 1
    await userEvent.click(filter1PillButton);

    await waitFor(() => {
      // verify last remaining cache value is expected
      const cacheItem = getTableQueueFilterSessionStorageValue(testSessionStorageKey);
      expect(cacheItem[0].value === MOVE_STATUS_OPTIONS[2].value).toBeTruthy();
      expect(cacheItem.length).toBe(1);

      // verify delete all button is gone
      expect(screen.queryByTestId('remove-filters-all')).not.toBeInTheDocument();

      // verify the one deleted is not rendered again
      expect(screen.queryByTestId(`remove-filters-col3-${MOVE_STATUS_OPTIONS[0].value}`)).not.toBeInTheDocument();

      // verify the one not deleted remains and rendered
      expect(screen.queryByTestId(`remove-filters-col3-${MOVE_STATUS_OPTIONS[2].value}`)).toBeInTheDocument();
    });

    // delete filter 2
    const remainingPillButton = await screen.findByTestId(`remove-filters-col3-${MOVE_STATUS_OPTIONS[2].value}`);
    expect(remainingPillButton).toBeInTheDocument();
    await userEvent.click(remainingPillButton);

    await waitFor(() => {
      // verify all filter are cleared
      const cacheItem = getTableQueueFilterSessionStorageValue(testSessionStorageKey);
      expect(cacheItem.length).toBe(0);

      // verify delete all button is gone
      expect(screen.queryByTestId('remove-filters-all')).not.toBeInTheDocument();

      // verify the one deleted is not rendered again
      expect(screen.queryByTestId(`remove-filters-col3-${MOVE_STATUS_OPTIONS[0].value}`)).not.toBeInTheDocument();

      // verify the one not deleted remains and rendered
      expect(screen.queryByTestId(`remove-filters-col3-${MOVE_STATUS_OPTIONS[2].value}`)).not.toBeInTheDocument();
    });
  });
});

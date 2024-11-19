import React from 'react';
import { Provider } from 'react-redux';
import userEvent from '@testing-library/user-event';
import ReactTable from 'react-table-6';
import { mount } from 'enzyme/build';
import { render, screen } from '@testing-library/react';

import QueueTable from './QueueTable';

import store from 'shared/store';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockRouterProvider } from 'testUtils';

const mockLogOut = jest.fn();
jest.mock('store/auth/actions', () => ({
  ...jest.requireActual('store/auth/actions'),
  logOut: () => mockLogOut,
}));

describe('Shipments column', () => {
  let wrapper;

  it('renders "PPM" when it is a PPM move', (done) => {
    wrapper = mountComponents(
      retrieveMovesStub({
        ppm_status: 'PAYMENT_REQUESTED',
        hhg_status: undefined,
      }),
    );

    setTimeout(() => {
      const move = getMove(wrapper);
      expect(move.shipments).toEqual(SHIPMENT_OPTIONS.PPM);
    });
  });

  it('does not display when the queue type is anything other than "new"', (done) => {
    wrapper = mountComponents(
      retrieveMovesStub({
        ppm_status: undefined,
        hhg_status: undefined,
      }),
      'ppm',
    );

    setTimeout(() => {
      const move = getMove(wrapper);
      expect(move.shipments);
    });
  });
});

describe('Refreshing', () => {
  it('loads the data again', async () => {
    const mockRetrieveMoves = jest.fn();
    render(
      <Provider store={store}>
        <MockRouterProvider>
          <QueueTable queueType="new" retrieveMoves={mockRetrieveMoves} />
        </MockRouterProvider>
      </Provider>,
    );

    // Should have fetched the data once on load
    expect(mockRetrieveMoves).toHaveBeenCalledTimes(1);

    // Click the refresh button
    const refreshButton = await screen.getByTestId('refreshQueue');
    await userEvent.click(refreshButton);

    // Should have fetched the data again
    await expect(mockRetrieveMoves).toHaveBeenCalledTimes(2);
  });
});

describe('on 401 unauthorized error', () => {
  it('force user log out', async () => {
    const mockRetrieveMoves = jest.fn();

    render(
      <Provider store={store}>
        <MockRouterProvider>
          <QueueTable queueType="new" retrieveMoves={mockRetrieveMoves} />
        </MockRouterProvider>
      </Provider>,
    );

    // Should have initially retrieved moves without error and not logged out
    await expect(mockRetrieveMoves).toHaveBeenCalledTimes(1);
    await expect(mockLogOut).toHaveBeenCalledTimes(0);

    // Mock the retrieve moves function to throw a 401 error
    mockRetrieveMoves.mockImplementation(() => {
      const error = new Error('Unauthorized');
      error.status = 401;
      throw error;
    });

    // Click the refresh button
    const refreshButton = screen.getByTestId('refreshQueue');
    await userEvent.click(refreshButton);

    // Should have retreived moves and failed causing a log out
    await expect(mockRetrieveMoves).toHaveBeenCalledTimes(2);
    await expect(mockLogOut).toHaveBeenCalledTimes(1);
  });
});

function retrieveMovesStub(params, throwError) {
  // This is meant as a stub that will act in place of
  // `RetrieveMovesForOffice` from Office/api.js
  return async () => {
    return await new Promise((resolve) => {
      if (throwError) {
        throw throwError;
      }

      resolve([
        {
          id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
          ...params,
        },
      ]);
    });
  };
}

function mountComponents(getMoves, queueType = 'new', mockStore = store) {
  return mount(
    <Provider store={mockStore}>
      <MockRouterProvider>
        <QueueTable queueType={queueType} retrieveMoves={getMoves} />
      </MockRouterProvider>
    </Provider>,
  );
}

function getMove(wrapper) {
  return wrapper.find(ReactTable).state().data[0];
}

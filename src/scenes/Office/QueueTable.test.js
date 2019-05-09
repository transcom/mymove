import React from 'react';
import { Provider } from 'react-redux';
import MockRouter from 'react-mock-router';

import QueueTable from './QueueTable';
import ReactTable from 'react-table';
import store from 'shared/store';
import { mount } from 'enzyme/build';

const push = jest.fn();

describe('Shipments column', () => {
  let wrapper;

  it('renders "PPM" when it is a PPM move', async () => {
    wrapper = mountComponents(
      retrieveMoves({
        ppm_status: 'PAYMENT_REQUESTED',
        hhg_status: undefined,
      }),
    );

    await resolveAllPromises();

    const move = getMove(wrapper);

    expect(move.shipments).toEqual('PPM');
  });

  it('renders "HHG" when it is a HHG move', async () => {
    wrapper = mountComponents(
      retrieveMoves({
        ppm_status: undefined,
        hhg_status: 'APPROVED',
      }),
    );

    await resolveAllPromises();

    const move = getMove(wrapper);

    expect(move.shipments).toEqual('HHG');
  });

  it('renders "HHG, PPM" when it is a combo move', async () => {
    wrapper = mountComponents(
      retrieveMoves({
        ppm_status: 'PAYMENT_REQUESTED',
        hhg_status: 'APPROVED',
      }),
    );

    await resolveAllPromises();

    const move = getMove(wrapper);

    expect(move.shipments).toEqual('HHG, PPM');
  });
});

function retrieveMoves(params) {
  // This is meant as a stub that will act in place of
  // `RetrieveMovesForOffice` from Office/api.js
  return async () => {
    return [
      {
        id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
        ...params,
      },
    ];
  };
}

function mountComponents(getMoves) {
  return mount(
    <Provider store={store}>
      <MockRouter push={push}>
        <QueueTable queueType="new" retrieveMoves={getMoves} />
      </MockRouter>
    </Provider>,
  );
}

function getMove(wrapper) {
  return wrapper.find(ReactTable).state().data[0];
}

function resolveAllPromises() {
  // Forces all promises that are returned inside the
  // component to resolve before the one returned here,
  // effectively giving us the end state of the component
  // as it would be rendered on the page.

  return new Promise(resolve => {
    setTimeout(resolve, 0);
  });
}

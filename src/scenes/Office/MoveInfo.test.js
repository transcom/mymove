import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';

import MoveInfo from './MoveInfo';
import store from 'shared/store';
import { mount } from 'enzyme/build';
import { ReferrerQueueLink } from './MoveInfo';
import { MemoryRouter } from 'react-router';

const dummyFunc = () => {};
const loadDependenciesHasError = null;
const loadDependenciesHasSuccess = false;
const location = {
  pathname: '',
};
const match = {
  params: { moveID: '123456' },
  url: 'www.nino.com',
  path: '/moveIt/moveIt',
};

const push = jest.fn();

describe('Loads MoveInfo', () => {
  it('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(
      <Provider store={store}>
        <MemoryRouter push={push}>
          <MoveInfo
            loadDependenciesHasError={loadDependenciesHasError}
            loadDependenciesHasSuccess={loadDependenciesHasSuccess}
            location={location}
            match={match}
            loadMoveDependencies={dummyFunc}
          />
        </MemoryRouter>
      </Provider>,
      div,
    );
  });
});

let wrapper;
describe('ShipmentInfo tests', () => {
  describe('Shows correct queue to return to', () => {
    it('when a referrer is set in history', () => {
      wrapper = mount(
        <Provider store={store}>
          <MemoryRouter push={jest.fn()}>
            <ReferrerQueueLink history={{ location: { state: { referrerPathname: '/queues/hhg_active' } } }} />
          </MemoryRouter>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('Active HHG Queue');
    });
    it('when no referrer is set', () => {
      wrapper = mount(
        <Provider store={store}>
          <MemoryRouter push={jest.fn()}>
            <ReferrerQueueLink history={{ location: {} }} />
          </MemoryRouter>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('New Moves/Shipments Queue');
    });
  });
});

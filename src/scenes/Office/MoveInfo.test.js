import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import MockRouter from 'react-mock-router';

import MoveInfo from './MoveInfo';
import store from 'shared/store';
import { mount } from 'enzyme/build';
import { ReferrerQueueLink } from './MoveInfo';

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
        <MockRouter push={push}>
          <MoveInfo
            loadDependenciesHasError={loadDependenciesHasError}
            loadDependenciesHasSuccess={loadDependenciesHasSuccess}
            location={location}
            match={match}
            loadMoveDependencies={dummyFunc}
          />
        </MockRouter>
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
          <MockRouter push={jest.fn()}>
            <ReferrerQueueLink history={{ location: { state: { referrerPathname: '/queues/hhg_accepted' } } }} />
          </MockRouter>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('Accepted HHG Queue');
    });
    it('when no referrer is set', () => {
      wrapper = mount(
        <Provider store={store}>
          <MockRouter push={jest.fn()}>
            <ReferrerQueueLink history={{ location: {} }} />
          </MockRouter>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('New Moves/Shipments Queue');
    });
  });
});

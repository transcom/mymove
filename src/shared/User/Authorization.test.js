import React from 'react';
import { mount } from 'enzyme';
import configureStore from 'redux-mock-store';

import Authorization from './Authorization';

const Groot = props => (
  <div>
    <span>I am groot</span>
  </div>
);
describe('Authorization tests', () => {
  describe('When the user is in the role', () => {
    let mockStore = configureStore();
    let initialState = {
      user: {
        features: ['fakeFeature'],
      },
    };
    const RouteWithAuthorization = Authorization(Groot, 'fakeFeature');
    it('the component is rendered', () => {
      let store = mockStore(initialState);
      let wrap = mount(<RouteWithAuthorization store={store} />);
      expect(wrap.containsMatchingElement(<span>I am groot</span>)).toBeTruthy();
    });
  });
  describe('When the user is not in the role', () => {
    let mockStore = configureStore();
    let initialState = {
      user: {
        features: ['somethingElse'],
      },
    };
    const RouteWithAuthorization = Authorization(Groot, 'fakeFeature');
    it('an auth message is rendered', () => {
      let store = mockStore(initialState);
      let wrap = mount(<RouteWithAuthorization store={store} />);
      expect(wrap.containsMatchingElement(<h1>You are not authorized to view this page</h1>)).toBeTruthy();
    });
  });
});

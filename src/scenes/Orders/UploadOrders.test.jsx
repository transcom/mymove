import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { HashRouter as Router } from 'react-router-dom';
import thunk from 'redux-thunk';
import UploadOrders from './UploadOrders';

const defaultProps = {
  pages: [],
  pageKey: '',
  fetchLatestOrders: () => {},
};

const mockStore = configureStore([thunk]);
const initialState = {
  entities: {
    orders: {},
  },
};

const store = mockStore(initialState);

function mountUploadOrders(props = defaultProps) {
  return mount(
    <Provider store={store}>
      <Router>
        <UploadOrders {...props} />
      </Router>
    </Provider>,
  );
}

describe('UploadOrders component', () => {
  it('renders component with no cancel button', () => {
    const wrapper = mountUploadOrders();
    expect(wrapper.find('[data-testid="wizardCancelButton"]').length).toBe(0);
  });
});

/*  react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { HashRouter as Router } from 'react-router-dom';
import thunk from 'redux-thunk';

//  import/no-named-as-default
import UploadOrders from './UploadOrders';

const defaultProps = {
  pages: ['1', '2', '3'],
  pageKey: '2',
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
  it('renders component with next button', () => {
    const wrapper = mountUploadOrders();
    expect(wrapper.find('[data-testid="wizardNextButton"]').length).toBe(1);
  });
});

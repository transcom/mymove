import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import { history, store } from '../../shared/store';

import PaymentRequestReview from './PaymentRequestReview';

describe('PaymentRequestReview', () => {
  const wrapper = mount(
    <Provider store={store}>
      clear
      <ConnectedRouter history={history}>
        <PaymentRequestReview />
      </ConnectedRouter>
    </Provider>,
  );

  it('renders without errors', () => {
    expect(wrapper.find('[data-testid="PaymentRequestReview"]').exists()).toBe(true);
  });
});

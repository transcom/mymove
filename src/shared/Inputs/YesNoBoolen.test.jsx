import YesNoBoolean from './YesNoBoolean';
import { render } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from 'shared/store';

const testProps = {
  input: {
    value: true,
    onChange: jest.fn(),
  },
};

describe('YesNoBoolean', () => {
  describe('with default props', () => {
    it('renders without errors', () => {
      const mockStore = configureStore({});
      render(
        <Provider store={mockStore.store}>
          <YesNoBoolean {...testProps} />
        </Provider>,
      );
    });

    it('renders without input', () => {
      const mockStore = configureStore({});
      render(
        <Provider store={mockStore.store}>
          <YesNoBoolean />
        </Provider>,
      );
    });
  });
});

import React from 'react';
import { mount } from 'enzyme';
import { GridContainer, Alert } from '@trussworks/react-uswds';
import { Provider } from 'react-redux';

import FlashGridContainer from 'containers/FlashGridContainer/FlashGridContainer';
import ScrollToTop from 'components/ScrollToTop';
import ConnectedFlashMessage from 'containers/FlashMessage/FlashMessage';
import { configureStore, history } from 'shared/store';
import { setFlashMessage } from 'store/flash/actions';
import { MockProviders } from 'testUtils';

global.scrollTo = jest.fn();

describe('FlashGridContainer component', () => {
  it('renders the same HTML as USWDS GridContainer if no message is set', () => {
    const flashGridWrapper = mount(
      <MockProviders>
        <FlashGridContainer className="test-class-1 test-class-2" data-testid="test-base-html">
          <h1>Test Container</h1>
          <p>Testing children.</p>
        </FlashGridContainer>
      </MockProviders>,
    );
    const uswdsGridWrapper = mount(
      <GridContainer className="test-class-1 test-class-2" data-testid="test-base-html">
        <h1>Test Container</h1>
        <p>Testing children.</p>
      </GridContainer>,
    );

    const flashGrid = flashGridWrapper.find('FlashGridContainer');
    const uswdsGrid = uswdsGridWrapper.find(GridContainer);

    expect(flashGrid.exists()).toBe(true);
    expect(uswdsGrid.exists()).toBe(true);
    expect(flashGrid.html()).toEqual(uswdsGrid.html());

    expect(flashGrid.hasClass('test-class-1 test-class-2')).toBe(true);
    expect(flashGrid.prop('data-testid')).toEqual('test-base-html');
    // children should be scroll, alert, h1, and p:
    expect(flashGridWrapper.find('div[data-testid="test-base-html"]').children().length).toBe(4);
    expect(flashGridWrapper.find(Alert).exists()).toBe(false);
  });

  it('renders an alert at the top if a message is set', () => {
    const testState = {
      flash: {
        flashMessage: {
          type: 'success',
          message: 'This is a successful message!',
          key: 'TEST_SUCCESS_FLASH',
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={testState}>
        <FlashGridContainer data-testid="test-alert">
          <p>Testing alert.</p>
        </FlashGridContainer>
      </MockProviders>,
    );

    const container = wrapper.find('div[data-testid="test-alert"]');
    expect(container.exists()).toBe(true);
    expect(container.childAt(0).type()).toEqual(ScrollToTop);
    expect(container.childAt(1).type()).toEqual(ConnectedFlashMessage);
    expect(container.childAt(2).type()).toEqual('p');

    const alert = wrapper.find(Alert);
    expect(alert.exists()).toBe(true);
    expect(alert.text()).toEqual('This is a successful message!');
  });

  it('scrolls up to the alert if a new message is set', () => {
    const mockStore = configureStore(history, {});
    const wrapper = mount(
      <Provider store={mockStore.store}>
        <FlashGridContainer data-testid="test-store">
          <p>Testing scroll.</p>
        </FlashGridContainer>
      </Provider>,
    );

    // ScrollToTop should fire at the initial mount
    expect(global.scrollTo).toHaveBeenCalledTimes(1);

    mockStore.store.dispatch(setFlashMessage('TEST_SUCCESS_FLASH', 'success', 'This is a successful message!'));

    // Re-render after changing state with new message and ScrollToTop should fire again
    wrapper.mount();
    expect(global.scrollTo).toHaveBeenCalledTimes(2);

    // Re-render without changing the message state and ScrollToTop should NOT fire
    wrapper.mount();
    expect(global.scrollTo).toHaveBeenCalledTimes(2);
  });

  it('clears the alert if a new element is focused (in and out) and does not scroll', () => {});

  afterEach(() => {
    jest.clearAllMocks();
  });
});

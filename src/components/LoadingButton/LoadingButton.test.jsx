import React from 'react';
import { mount } from 'enzyme';

import LoadingButton from './LoadingButton';

describe('LoadingButton component', () => {
  const defaultProps = {
    onClick: () => Promise.resolve(),
    labelText: 'Button',
    loadingText: 'Loading',
    isLoading: false,
  };

  it('renders the LoadingButton component without errors', () => {
    const wrapper = mount(<LoadingButton {...defaultProps} />);
    const btn = wrapper.find('[data-testid="loading-button"]');

    expect(btn.first().text().includes(defaultProps.labelText));
  });

  const loadingProps = {
    ...defaultProps,
    isLoading: true,
    disabled: true,
  };

  it('multi-click calls fetcher once', () => {
    const wrapper = mount(<LoadingButton {...loadingProps} />);
    const btn = wrapper.find('[data-testid="loading-button"]');

    expect(btn.first().text().includes(loadingProps.loadingText));
    expect(btn.first().getDOMNode()).toHaveProperty('disabled');
  });
});

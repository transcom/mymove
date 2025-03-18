// StyledLine.test.js
import React from 'react';
import { render } from '@testing-library/react';

import { StyledLine } from './StyledLine';

describe('StyledLine', () => {
  it('renders with default styles when no props are provided', () => {
    const { container } = render(<StyledLine />);
    const div = container.firstChild;

    expect(div).toHaveStyle('width: 75%');
    expect(div).toHaveStyle('background-color: #565c65');

    expect(div.className).toBeTruthy();
  });

  it('renders with provided inline styles', () => {
    const customWidth = '50%';
    const customColor = '#ff0000';
    const { container } = render(<StyledLine width={customWidth} color={customColor} />);
    const div = container.firstChild;

    expect(div).toHaveStyle(`width: ${customWidth}`);
    expect(div).toHaveStyle(`background-color: ${customColor}`);
  });

  it('uses the custom className if provided, overriding the default', () => {
    const customClass = 'my-custom-line';
    const { container } = render(<StyledLine className={customClass} />);
    const div = container.firstChild;

    // the component renders the className as: `${className || styles.styledLine}`.
    // so when customClass is provided, it should be the only class applied.
    expect(div.className).toBe(customClass);
  });
});

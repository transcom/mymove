import React from 'react';
import { render } from '@testing-library/react';

import AccessibilityStatement from './AccessibilityStatement';

describe('Accessibility Statement page', () => {
  it('has the correct title', () => {
    render(<AccessibilityStatement />);
    expect(document.title).toContain('508 Compliance');
  });
});

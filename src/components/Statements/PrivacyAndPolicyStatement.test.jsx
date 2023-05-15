import React from 'react';
import { render } from '@testing-library/react';

import PrivacyPolicy from './PrivacyAndPolicyStatement';

describe('Privacy Policy page', () => {
  it('has the correct title', () => {
    render(<PrivacyPolicy />);
    expect(document.title).toContain('Privacy & Security Policy');
  });
});

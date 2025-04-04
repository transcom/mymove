import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import { render } from '@testing-library/react';

import Footer from './index';

import { milmoveHelpDesk } from 'shared/constants';

describe('Footer', () => {
  it('has helpdesk email', () => {
    render(
      <BrowserRouter>
        <Footer />
      </BrowserRouter>,
    );

    const obj = document.getElementById('helpMeLink');

    Object.keys(obj).forEach((key) => {
      if (key === 'href') {
        expect(obj[key]).toContain(milmoveHelpDesk);
      }
    });
  });
});

import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import { render } from '@testing-library/react';

import Footer from './index';

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
        expect(obj[key]).toContain('mailto:usarmy.scott.sddc.mbx.G6-SRC-MilMove-HD@army.mil');
      }
    });
  });
});

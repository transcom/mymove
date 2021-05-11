import React from 'react';
import { render, screen } from '@testing-library/react';

import DetailsPanel from './DetailsPanel';

describe('DetailsPanel', () => {
  it('renders a title', () => {
    render(
      <DetailsPanel title="My title">
        <p>Some child content</p>
      </DetailsPanel>,
    );
    expect(screen.getByRole('heading', { level: 2 }).textContent).toEqual('My title');
  });

  it('renders child content', () => {
    render(
      <DetailsPanel title="My title">
        <p>Some child content</p>
      </DetailsPanel>,
    );
    expect(screen.getByText('Some child content')).toBeInTheDocument();
  });

  it('renders an edit button', () => {
    render(
      <DetailsPanel title="My title" editButton={<a href="/some-link">Edit link</a>}>
        <p>Some child content</p>
      </DetailsPanel>,
    );
    expect(screen.getByRole('link').textContent).toEqual('Edit link');
  });
});

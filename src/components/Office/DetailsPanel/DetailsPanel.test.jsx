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

  it('renders a title with a tag', () => {
    render(
      <DetailsPanel title="My title" tag="NEW">
        <p>Some child content</p>
      </DetailsPanel>,
    );

    expect(screen.getByText('My title')).toBeInTheDocument();
    expect(screen.getByTestId('detailsPanelTag')).toBeInTheDocument();
  });

  it('renders child content', () => {
    render(
      <DetailsPanel title="My title">
        <p>Some child content</p>
      </DetailsPanel>,
    );
    expect(screen.getByText('Some child content')).toBeInTheDocument();
  });
});

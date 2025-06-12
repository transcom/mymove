import React from 'react';
import { render, screen } from '@testing-library/react';

import DetailsPanel from './DetailsPanel';

import ButtonDropdown from 'components/ButtonDropdown/ButtonDropdown';

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

  it('renders buttons', () => {
    render(
      <DetailsPanel
        title="My title"
        tag="NEW"
        editButton={
          <button type="button" href="#" className="usa-button usa-button--secondary">
            Edit
          </button>
        }
      >
        <p>Some child content</p>
      </DetailsPanel>,
    );

    expect(screen.getByText('My title')).toBeInTheDocument();
    expect(screen.getByTestId('detailsPanelTag')).toBeInTheDocument();
    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  it('renders combobox', () => {
    render(
      <DetailsPanel
        title="My title"
        tag="NEW"
        editButton={
          <ButtonDropdown>
            <option value="">Dropdown Button</option>
            <option>Option 1</option>
            <option>Option 2</option>
            <option>Option 3</option>
            <option>Option 4</option>
          </ButtonDropdown>
        }
      >
        <p>Some child content</p>
      </DetailsPanel>,
    );

    expect(screen.getByText('My title')).toBeInTheDocument();
    expect(screen.getByTestId('detailsPanelTag')).toBeInTheDocument();
    expect(screen.queryByRole('combobox')).toBeInTheDocument();
  });

  it('renders button and combobox', () => {
    render(
      <DetailsPanel
        title="My title"
        tag="NEW"
        editButton={
          <ButtonDropdown>
            <option value="">Dropdown Button</option>
            <option>Option 1</option>
            <option>Option 2</option>
            <option>Option 3</option>
            <option>Option 4</option>
          </ButtonDropdown>
        }
        reviewButton={<button type="button" label="Review Button" />}
      >
        <p>Some child content</p>
      </DetailsPanel>,
    );

    expect(screen.getByText('My title')).toBeInTheDocument();
    expect(screen.getByTestId('detailsPanelTag')).toBeInTheDocument();
    expect(screen.queryByRole('combobox')).toBeInTheDocument();
    expect(screen.getByRole('button')).toBeInTheDocument();
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

import React from 'react';
import { render, screen } from '@testing-library/react';
import { Tag } from '@trussworks/react-uswds';

import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

describe('WeightDisplay', () => {
  it('renders without crashing', () => {
    render(<WeightDisplay heading="heading test" />);

    expect(screen.getByText('heading test')).toBeInTheDocument();
  });

  it('renders with weight value', () => {
    render(<WeightDisplay heading="heading test" weightValue={1234} />);

    expect(screen.getByText('1,234 lbs')).toBeInTheDocument();
  });

  it('renders with edit button when user has permissions', () => {
    render(
      <MockProviders permissions={[permissionTypes.updateBillableWeight]}>
        <WeightDisplay heading="heading test" value={1234} onEdit={jest.fn()} />
      </MockProviders>,
    );

    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  it('renders with no edit button when user does not have permissions', () => {
    render(<WeightDisplay heading="heading test" value={1234} />);

    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('renders with react element as a child', () => {
    render(
      <WeightDisplay heading="heading test" value={1234}>
        <Tag>tag passed in</Tag>
      </WeightDisplay>,
    );

    expect(screen.getByText('tag passed in')).toBeInTheDocument();
  });

  it('renders with text as a child', () => {
    render(
      <WeightDisplay heading="heading test" value={1234}>
        text passed in
      </WeightDisplay>,
    );

    expect(screen.getByText('text passed in')).toBeInTheDocument();
  });

  it('edit button is clicked', () => {
    const mockEditBtn = jest.fn();
    render(<WeightDisplay heading="heading test" value={1234} showEditBtn onEdit={mockEditBtn} />);
    screen.getByRole('button').click();

    expect(mockEditBtn).toHaveBeenCalledTimes(1);
  });
});

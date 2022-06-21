import React from 'react';
import { render, screen } from '@testing-library/react';

import Restricted from 'components/Restricted/Restricted';
import PermissionProvider from 'components/Restricted/PermissionProvider';

describe('Restricted', () => {
  const testPermission = 'test.permission';
  const badPermission = 'bad.permission';

  it('renders children when there are matching permissions', () => {
    render(
      <PermissionProvider permissions={[testPermission]}>
        <Restricted to={testPermission}>
          <div>Child Component</div>
        </Restricted>
      </PermissionProvider>,
    );

    expect(screen.getByText('Child Component')).toBeInTheDocument();
  });

  it('renders the fallback when there are mismatched permissions', () => {
    render(
      <PermissionProvider permissions={[testPermission]}>
        <Restricted to={badPermission} fallback={<div>Fallback Component</div>}>
          <div>Child Component</div>
        </Restricted>
      </PermissionProvider>,
    );

    expect(screen.queryByText('Child Component')).not.toBeInTheDocument();
    expect(screen.getByText('Fallback Component')).toBeInTheDocument();
  });

  it('does not render children when there are no permissions provided', () => {
    render(
      <Restricted to={testPermission}>
        <div>Child Component</div>
      </Restricted>,
    );

    expect(screen.queryByText('Child Component')).not.toBeInTheDocument();
  });

  it('does not render children when there are mismatched permissions', () => {
    render(
      <PermissionProvider permissions={[testPermission]}>
        <Restricted to={badPermission}>
          <div>Child Component</div>
        </Restricted>
      </PermissionProvider>,
    );

    expect(screen.queryByText('Child Component')).not.toBeInTheDocument();
  });

  it('does not render the fallback when there are matched permissions', () => {
    render(
      <PermissionProvider permissions={[testPermission]}>
        <Restricted to={testPermission} fallback={<div>Fallback Component</div>}>
          <div>Child Component</div>
        </Restricted>
      </PermissionProvider>,
    );

    expect(screen.getByText('Child Component')).toBeInTheDocument();
    expect(screen.queryByText('Fallback Component')).not.toBeInTheDocument();
  });

  it('renders children if current user matches the restricted access user', () => {
    const testUserId = 'testUserId';

    render(
      <PermissionProvider currentUserId={testUserId}>
        <Restricted user={testUserId}>
          <div>Child Component</div>
        </Restricted>
      </PermissionProvider>,
    );

    expect(screen.getByText('Child Component')).toBeInTheDocument();
  });

  it('does not render children if current user does not match the restricted access user', () => {
    const userId1 = 'user1';
    const userId2 = 'user2';

    render(
      <PermissionProvider currentUserId={userId1}>
        <Restricted user={userId2}>
          <div>Child Component</div>
        </Restricted>
      </PermissionProvider>,
    );

    expect(screen.queryByText('Child Component')).not.toBeInTheDocument();
  });
});

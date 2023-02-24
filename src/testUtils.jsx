import React, { Suspense } from 'react';
import { node, shape, arrayOf, string, bool } from 'prop-types';
import { Provider } from 'react-redux';
import { RouterProvider, createMemoryRouter, generatePath } from 'react-router-dom';
import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { configureStore } from 'shared/store';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PermissionProvider from 'components/Restricted/PermissionProvider';

// Helper function to create a react-router memory router with the provided options
const createMockRouter = ({ path, params, routes, initialEntries, children }) => {
  const mockRoutes = [
    {
      path: path || '/',
      element: <Suspense fallback={<LoadingPlaceholder />}>{children}</Suspense>,
    },
  ];

  if (routes && routes.length > 0) mockRoutes.push(...routes);

  const router = createMemoryRouter(mockRoutes, {
    initialEntries: [initialEntries || generatePath(path, params)] || ['/'],
  });

  return router;
};

export const ReactQueryWrapper = ({ children, client }) => {
  const queryClient = client || new QueryClient();
  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
};

/**
 * Render the `ui` with mock routing in place, setup using the `options`. Most common options are `path` and `params`.
 * Enables rendered components to use routing hooks like `useParams` and `useLocation`.
 *
 * @param {*} ui - The component to be rendered using react-testing-library
 * @param {*} options - Routing options used to create a mock router. Prefered options are `path` and `params`. Other options are for extended use cases and supporting existiong patterns.
 * @returns {*} - The result of the render call from react-testing-library
 */
export const renderWithRouter = (ui, options) => {
  return render(<MockRouting {...options}>{ui}</MockRouting>);
};

export const renderWithProviders = (ui, options) => {
  return render(
    <ReactQueryWrapper>
      <MockProviders {...options}>{ui}</MockProviders>
    </ReactQueryWrapper>,
  );
};

/**
 * Renders class components with a `routing` prop. For use with components that use the `withRouter` HOC.
 * Adds a mock router prop to the component using the provided path and params in options to mock behavior of `withRouter`.
 *
 * @param {*} ui - The component to be rendered using react-testing-library
 * @param {*} options - Routing options used to create a mock router. Prefered options are `path` and `params`. Other options are for extended use cases and supporting existiong patterns.
 * @returns {*} - The result of the render call from react-testing-library
 */
export const renderWithRouterProp = (ui, options) => {
  const path = options?.path || '/';
  const params = options?.params || {};
  const navigate = options?.navigate || jest.fn();
  const search = options?.search || '';

  const pathname = generatePath(path, params);
  const router = { location: { pathname, search }, params, navigate };

  if (options?.includeProviders) return renderWithProviders(React.cloneElement(ui, { router: { ...router } }));
  return renderWithRouter(React.cloneElement(ui, { router: { ...router } }));
};

/** Wrap the provided children with a mock router using the provided options. */
export const MockRouting = ({ children, path, params, initialEntries, routes }) => {
  const mockRouter = createMockRouter({ path, params, routes, initialEntries, children });

  return <RouterProvider router={mockRouter} />;
};

MockRouting.propTypes = {
  children: node.isRequired,
  params: shape({}),
  path: string,
  routes: arrayOf(
    shape({
      path: string,
      element: node,
      children: arrayOf(shape({})),
      caseSensitive: bool,
    }),
  ),
};

MockRouting.defaultProps = {
  params: {},
  path: '/',
  routes: null,
};

/** Wrap the three most common mock providers (permission, redux, and router) around the provided children */
export const MockProviders = ({
  children,
  initialState,
  initialEntries,
  path,
  params,
  routes,
  permissions,
  currentUserId,
  client,
}) => {
  const mockRouter = createMockRouter({ path, params, routes, initialEntries, children });
  const mockStore = configureStore(initialState);

  return (
    <ReactQueryWrapper client={client}>
      <PermissionProvider permissions={permissions} currentUserId={currentUserId}>
        <Provider store={mockStore.store}>
          <RouterProvider router={mockRouter} />
        </Provider>
      </PermissionProvider>
    </ReactQueryWrapper>
  );
};

MockProviders.propTypes = {
  children: node.isRequired,
  initialState: shape({}),
  params: shape({}),
  path: string,
  permissions: arrayOf(string),
  currentUserId: string,
};

MockProviders.defaultProps = {
  initialState: {},
  params: {},
  path: '/',
  permissions: [],
  currentUserId: null,
};

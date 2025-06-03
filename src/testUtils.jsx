import React from 'react';
import { node, shape, arrayOf, string, bool } from 'prop-types';
import { Provider } from 'react-redux';
import { RouterProvider, createMemoryRouter, generatePath } from 'react-router-dom';
import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { configureStore } from 'shared/store';
import PermissionProvider from 'components/Restricted/PermissionProvider';
import SelectedGblocProvider from 'components/Office/GblocSwitcher/SelectedGblocProvider';

/** Helper function to create a react-router `MemoryRouter` with the provided options */
const createMockRouter = ({ path, params, routes, children }) => {
  // Add the current path (or /) to the routes array
  const mockRoutes = [
    {
      path: path || '/',
      element: children,
    },
  ];

  // include any additional routes provided
  if (routes && routes.length > 0) mockRoutes.push(...routes);

  const router = createMemoryRouter(mockRoutes, {
    initialEntries: [generatePath(path, params)] || ['/'],
  });

  return router;
};

/** Wrap the provided children with a react query client provider
 * @param {object} client - A react query client to use. If not provided, a new client will be created.
 * @returns {object} - The react query client provider
 * */
export const ReactQueryWrapper = ({ children, client }) => {
  const queryClient = client || new QueryClient();
  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
};

/** Wrap the provided children with a mock router using the provided options. */
export const MockRouterProvider = ({ children, path, params, routes }) => {
  const mockRouter = createMockRouter({ path, params, routes, children });

  return <RouterProvider router={mockRouter} />;
};

MockRouterProvider.propTypes = {
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

MockRouterProvider.defaultProps = {
  params: {},
  path: '/',
  routes: null,
};

/** Wrap the four most common mock providers (permission, redux, and router) around the provided children */
export const MockProviders = ({
  children,
  initialState, // redux
  path, // routing
  params, // routing
  routes, // routing
  permissions, // permissions
  currentUserId, // permissions
  client, // react query
}) => {
  const mockRouter = createMockRouter({ path, params, routes, children });
  const mockStore = configureStore(initialState);

  return (
    <ReactQueryWrapper client={client}>
      <PermissionProvider permissions={permissions} currentUserId={currentUserId}>
        <SelectedGblocProvider>
          <Provider store={mockStore.store}>
            <RouterProvider router={mockRouter} />
          </Provider>
        </SelectedGblocProvider>
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
  routes: arrayOf(
    shape({
      path: string,
      element: node,
      children: arrayOf(shape({})),
      caseSensitive: bool,
    }),
  ),
};

MockProviders.defaultProps = {
  initialState: {},
  params: {},
  path: '/',
  permissions: [],
  currentUserId: null,
  routes: null,
};

/**
 * Render the `ui` with mock routing in place, setup using the `options`. Most common options are `path` and `params`.
 * Enables rendered components to use routing hooks like `useParams` and `useLocation`.
 *
 * @param {*} ui - The component to be rendered using react-testing-library
 * @param {*} options - Routing options used to create a mock router. Prefered options are `path` and `params`. Other options are for extended use cases and supporting existing patterns.
 * @returns {*} - The result of the render call from react-testing-library
 */
export const renderWithRouter = (ui, options) => {
  return render(<MockRouterProvider {...options}>{ui}</MockRouterProvider>);
};

export const renderWithProviders = (ui, options) => {
  return render(<MockProviders {...options}>{ui}</MockProviders>);
};

/**
 * Renders class components with a `router` prop. For use with components that use the `withRouter` HOC.
 * Sets up mock routing and adds a mock `router` prop to the component using the provided path and params in options to mock behavior of `withRouter`.
 *
 * @param {*} ui - The component to be rendered using react-testing-library
 * @param {*} options - Routing options used to create a mock router. Prefered options are `path` and `params`. Other options are for extended use cases and supporting existing patterns.
 * @returns {*} - The result of the render call from react-testing-library
 */
export const renderWithRouterProp = (ui, options) => {
  const path = options?.path || '/';
  const params = options?.params || {};
  const navigate = options?.navigate || jest.fn();
  const search = options?.search || '';

  const pathname = generatePath(path, params);
  const router = { location: { pathname, search }, params, navigate };

  const routingOptions = { ...options, path, params, navigate, search };
  if (options?.includeProviders)
    return renderWithProviders(React.cloneElement(ui, { router: { ...router } }), routingOptions);

  return renderWithRouter(React.cloneElement(ui, { router: { ...router } }), routingOptions);
};

export const mockPage = (path, name) => {
  return jest.mock(path, () => {
    // Create component name from path, if not provided (e.g. 'MoveQueue' -> 'Move Queue')
    const componentName =
      name ||
      path
        .substring(path.lastIndexOf('/') + 1)
        .replace(/([A-Z])/g, ' $1')
        .trim();

    return () => <div>{`Mock ${componentName} Component`}</div>;
  });
};

export const flushPromises = () => Promise.resolve().then();

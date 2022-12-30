import React, { Suspense } from 'react';
import { node, shape, arrayOf, string, bool } from 'prop-types';
import { Provider } from 'react-redux';
import { RouterProvider, createMemoryRouter, generatePath } from 'react-router-dom';
import { render } from '@testing-library/react';

import { configureStore } from 'shared/store';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PermissionProvider from 'components/Restricted/PermissionProvider';

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

/**
 * Render the `ui` with mock routing in place, setup using the `options`.
 * Enables rendered components to use routing hooks like `useParams` and `useLocation`.
 *
 * @param {*} ui - The component to be rendered using react-testing-library
 * @param {*} options - Routing options used to create a mock router. See `MockRouting` for more detail on how these are used.
 * @returns {*} - The result of the render call from react-testing-library
 */
export const renderWithRouter = (ui, options) => {
  return render(<MockRouting {...options}>{ui}</MockRouting>);
};

/**
 * For use rendering class components that make use of the `withRouter` HOC.
 * Adds a mock router prop to the component using the provided path and params in options to mock behavior of `withRouter`.
 *
 * @param {*} ui - The component to be rendered using react-testing-library
 * @param {*} options - Routing options used to create a mock router. Both `path` and `params` are required. See `MockRouting` for more detail on how these are used.
 * @returns {*} - The result of the render call from react-testing-library
 */
export const renderWithRouterProp = (ui, options) => {
  // if (!options) throw new Error('renderWithRouterProp requires options to be passed in');
  // if (!options.path) throw new Error('renderWithRouterProp requires a path to be included in the options');
  // if (!options.params) throw new Error('renderWithRouterProp requires params to be included in the options');

  // TODO: Do this cleaner and more distinct. There could be options that are not path/params (like navigate)
  const { path, params } = options && options.path && options.params ? options : { path: '/', params: {} };

  const pathname = generatePath(path, params);
  const router = { location: { pathname }, params, navigate: options.navigate || jest.fn() };

  return renderWithRouter(React.cloneElement(ui, { router: { ...router }, ...options }), options);
};

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
}) => {
  const mockRouter = createMockRouter({ path, params, routes, initialEntries, children });
  const mockStore = configureStore(initialState);

  return (
    <PermissionProvider permissions={permissions} currentUserId={currentUserId}>
      <Provider store={mockStore.store}>
        <RouterProvider router={mockRouter} />
      </Provider>
    </PermissionProvider>
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

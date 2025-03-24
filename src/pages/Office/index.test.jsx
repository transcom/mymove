import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { MemoryRouter } from 'react-router';

import OfficeApp from './index';

// import { roleTypes } from 'constants/userRoles';
import { configureStore } from 'shared/store';
// import { isBooleanFlagEnabled } from 'utils/featureFlags';

// Mock Redux actions to prevent actual API calls
jest.mock('store/auth/actions', () => ({
  loadUser: jest.fn(() => async () => {}),
}));

jest.mock('store/onboarding/actions', () => ({
  initOnboarding: jest.fn(() => async () => {}),
}));

jest.mock('shared/Swagger/ducks', () => ({
  loadInternalSchema: jest.fn(() => async () => {}),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const defaultState = {
  auth: {
    hasErrored: false,
    hasSucceeded: true,
    isLoading: false,
    isLoggedIn: false,
    underMaintenance: false,
  },
  swaggerInternal: {
    hasErrored: false,
    hasSucceeded: true,
    isLoading: false,
  },
  generalState: {
    moveId: '',
    showLoadingSpinner: false,
    loadingSpinnerMessage: null,
  },
};

const loggedOutState = {
  ...defaultState,
  auth: {
    ...defaultState.auth,
    activeRole: null,
    isLoggedIn: false,
  },
};

// const loggedInState = {
//   ...defaultState,
//   auth: {
//     ...defaultState.auth,
//     isLoggedIn: true,
//     activeRole: roleTypes.TOO,
//   },
// };

const renderWithState = (state, path) => {
  const mockStore = configureStore({ ...state });

  const minProps = {
    loadPublicSchema: jest.fn(),
    loadInternalSchema: jest.fn(),
    loadUser: jest.fn(),
  };

  return render(
    <MemoryRouter initialEntries={[path]}>
      <Provider store={mockStore.store}>
        <OfficeApp {...minProps} />
      </Provider>
    </MemoryRouter>,
  );
};

// const createMockStore = (role) => {
//   // Otherwise, use logged in state with the provided role
//   const state = {
//     auth: {
//       activeRole: role,
//       isLoading: false,
//       isLoggedIn: true,
//     },
//     swaggerInternal: {
//       hasErrored: false,
//       hasSucceeded: true,
//       isLoading: false,
//     },
//     generalState: {
//       moveId: '',
//       showLoadingSpinner: true,
//       loadingSpinnerMessage: 'test message',
//     },
//     entities: {
//       user: {
//         userId123: {
//           id: 'userId123',
//           roles: [{ roleType: role }],
//         },
//       },
//     },
//   };

//   return configureStore(state);
// };

// Render the OfficeApp component with routing and Redux setup for the provided route and role
// const renderOfficeAppAtRoute = (route, role) => {
//   isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
//   const mockStore = createMockStore(role);
//   const userRoles = role ? [{ roleType: role }] : [];
//   render(
//     <MemoryRouter initialEntries={[route]}>
//       <Provider store={mockStore.store}>
//         <OfficeApp
//           router={{ location: { pathname: route } }}
//           loadInternalSchema={jest.fn()}
//           loadPublicSchema={jest.fn()}
//           loadUser={jest.fn()}
//           hasRecentError={false}
//           activeRole={role || null}
//           userRoles={userRoles}
//           traceId=""
//           loginIsLoading={!!role}
//           userIsLoggedIn={!!role}
//           hqRoleFlag
//           gsrRoleFlag
//         />
//       </Provider>
//     </MemoryRouter>,
//   );
// };

describe('OfficeApp', () => {
  it('renders Sign In page when logged out', async () => {
    renderWithState(loggedOutState, '/sign-in');
    await waitFor(() => expect(screen.getByText(/sign in/i)).toBeInTheDocument());
  });
  it('displays maintenance page when under maintenance', async () => {
    const updatedState = {
      ...loggedOutState,
      auth: {
        ...loggedOutState.auth,
        underMaintenance: true,
      },
    };

    renderWithState(updatedState, '/');
    await waitFor(() =>
      expect(screen.getByText(/This system is currently undergoing maintenance/i)).toBeInTheDocument(),
    );
  });
  // it('shows the loading spinner when props are set to show it', async () => {
  //   const updatedState = {
  //     ...loggedInState,
  //     generalState: {
  //       ...loggedInState.generalState,
  //       showLoadingSpinner: true,
  //       loadingSpinnerMessage: 'test message',
  //     },
  //   };

  //   renderOfficeAppAtRoute('/', roleTypes.TOO);
  //   await waitFor(
  //     () => expect(screen.getByTestId('loading-spinner')).toBeInTheDocument(),
  //     expect(screen.getByText(/test message/i)).toBeInTheDocument(),
  //   );

  // });

  //   it('shows the loading spinner when props are set to show it', async () => {
  //     const updatedState = {
  //       ...loggedInState,
  //       generalState: {
  //         ...loggedInState.generalState,
  //         showLoadingSpinner: true,
  //         loadingSpinnerMessage: 'test message',
  //       },
  //     };

  //     renderWithState(updatedState, '/');
  //     await waitFor(
  //       () => expect(screen.getByTestId('loading-spinner')).toBeInTheDocument(),
  //       expect(screen.getByText(/test message/i)).toBeInTheDocument(),
  //     );
  //   });

  //   it('handles the Invalid Permissions URL for logged in user', async () => {
  //     renderWithState(loggedInState, '/invalid-permissions');

  //     expect(screen.getByText('Skip to content')).toBeInTheDocument();
  //     expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument();

  //     await waitFor(() =>
  //       expect(screen.getByText(/You do not have permission to access this site/i)).toBeInTheDocument(),
  //     );
  //   });

  //   it('shows the server error for logged in user', async () => {
  //     renderWithState(loggedInState, '/server_error');

  //     await waitFor(() => expect(screen.getByText('We are experiencing an internal server error')));
  //   });

  //   it('handles the forbidden URL for logged in user', async () => {
  //     renderWithState(loggedInState, '/forbidden');

  //     await waitFor(() => expect(screen.getByText('You are forbidden to use this endpoint')));
  //   });

  //   it('renders Footer component', async () => {
  //     renderWithState(loggedInState, '/');
  //     await waitFor(() => expect(screen.getByText(/Military OneSource/i)).toBeInTheDocument());
  //   });
});

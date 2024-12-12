import React, { useContext } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider } from 'react-redux';
import { createStore } from 'redux';

import GblocSwitcher from './GblocSwitcher';
import SelectedGblocProvider from './SelectedGblocProvider';
import SelectedGblocContext from './SelectedGblocContext';

import { appReducer } from 'appReducer';

const queryClient = new QueryClient();

const mockedState = {
  auth: {
    activeRole: 'services_counselor',
    isLoggedIn: true,
    hasSucceeded: true,
    hasErrored: false,
    isLoading: false,
    underMaintenance: false,
  },
  entities: {
    user: {
      'bf65095f-a70b-4e7e-b02c-136015fb417b': {
        office_user: {
          transportation_office: {
            gbloc: 'USMC',
            name: 'PPSO DMO Camp Lejeune - USMC ',
          },
          transportation_office_assignments: [
            {
              primaryOffice: true,
              transportationOffice: {
                gbloc: 'USMC',
                name: 'PPSO DMO Camp Lejeune - USMC ',
              },
            },
            {
              primaryOffice: false,
              transportationOffice: {
                gbloc: 'KKFA',
                name: 'JPPSO - North Central (KKFA) - USAF',
              },
            },
          ],
        },
      },
    },
  },
};

const store = createStore(appReducer(), mockedState);

const withDecorators = (Story) => {
  return (
    <SelectedGblocProvider>
      <Provider store={store}>
        <QueryClientProvider client={queryClient}>
          <Story />
        </QueryClientProvider>
      </Provider>
    </SelectedGblocProvider>
  );
};

export default {
  title: 'Office Components/GblocSwitcher',
  component: GblocSwitcher,
  decorators: [withDecorators],
};

const SelectedGblocDisplayer = () => {
  const { selectedGbloc } = useContext(SelectedGblocContext);
  return (
    <div>
      <b>Selected GBLOC provided by the SelectedGblocProvider: {selectedGbloc}</b>
    </div>
  );
};

export const defaultGblocSwitcher = () => {
  return (
    <>
      <div style={{ width: '110px' }}>
        <GblocSwitcher />
      </div>
      <SelectedGblocDisplayer />
    </>
  );
};

import React, { useContext } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider } from 'react-redux';

import GblocSwitcher from './GblocSwitcher';
import SelectedGblocProvider from './SelectedGblocProvider';
import SelectedGblocContext from './SelectedGblocContext';

import { configureStore } from 'shared/store';

const queryClient = new QueryClient();

const withQueryClient = (Story) => {
  const store = configureStore({
    auth: { activeRole: 'services_counselor' },
    entities: {
      user: { 'bf65095f-a70b-4e7e-b02c-136015fb417b': { officeUser: { transportation_office: { gbloc: 'KKFA' } } } },
    },
  });
  store.getState = () => {
    return {
      auth: { activeRole: 'services_counselor' },
      entities: {
        user: { 'bf65095f-a70b-4e7e-b02c-136015fb417b': { officeUser: { transportation_office: { gbloc: 'KKFA' } } } },
      },
    };
  };
  store.subscribe = () => {};
  return (
    <Provider store={store}>
      <QueryClientProvider client={queryClient}>
        <Story />
      </QueryClientProvider>
    </Provider>
  );
};

export default {
  title: 'Office Components/GblocSwitcher',
  component: GblocSwitcher,
  decorators: [withQueryClient],
};

const SelectedGblocDisplayer = () => {
  const { selectedGbloc } = useContext(SelectedGblocContext);
  return (
    <div>
      <b>Selected GBLOC provided by the SelectedGblocProvider: {selectedGbloc}</b>
    </div>
  );
};

// gblocsOverride is needed due to react-redux not meshing well with Storybook
export const defaultGblocSwitcher = () => {
  return (
    <SelectedGblocProvider>
      <div style={{ width: '110px' }}>
        <GblocSwitcher gblocsOverride={['KKFA', 'AGFM']} />
      </div>
      <SelectedGblocDisplayer />
    </SelectedGblocProvider>
  );
};

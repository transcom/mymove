import React, { useContext } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import GblocSwitcher from './GblocSwitcher';
import SelectedGblocProvider from './SelectedGblocProvider';
import SelectedGblocContext from './SelectedGblocContext';

const queryClient = new QueryClient();

const withQueryClient = (Story) => (
  <QueryClientProvider client={queryClient}>
    <Story />
  </QueryClientProvider>
);

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

export const defaultGblocSwitcher = () => {
  return (
    <SelectedGblocProvider>
      <div style={{ width: '110px' }}>
        <GblocSwitcher officeUsersDefaultGbloc="AGFM" />
      </div>

      <SelectedGblocDisplayer />
    </SelectedGblocProvider>
  );
};

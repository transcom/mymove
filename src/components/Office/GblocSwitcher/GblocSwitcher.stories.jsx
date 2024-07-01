import React, { useContext } from 'react';

import GblocSwitcher from './GblocSwitcher';
import SelectedGblocProvider from './SelectedGblocProvider';
import SelectedGblocContext from './SelectedGblocContext';

export default {
  title: 'Office Components/GblocSwitcher',
  component: GblocSwitcher,
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
        <GblocSwitcher officeUsersDefaultGbloc="KKFA" />
      </div>

      <SelectedGblocDisplayer />
    </SelectedGblocProvider>
  );
};

import React, { useMemo, useState } from 'react';
import PropTypes from 'prop-types';

import SelectedGblocContext, {
  SELECTED_GBLOC_SESSION_STORAGE_KEY,
} from 'components/Office/GblocSwitcher/SelectedGblocContext';

const SelectedGblocProvider = ({ children }) => {
  const [selectedGbloc, setSelectedGbloc] = useState(undefined);
  const handleGblocChange = (value) => {
    setSelectedGbloc(value);
    window.sessionStorage.setItem(SELECTED_GBLOC_SESSION_STORAGE_KEY, value);
  };

  const getterAndSetter = useMemo(
    () => ({
      selectedGbloc,
      handleGblocChange,
    }),
    [selectedGbloc],
  );

  return <SelectedGblocContext.Provider value={getterAndSetter}>{children}</SelectedGblocContext.Provider>;
};

SelectedGblocProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

export default SelectedGblocProvider;

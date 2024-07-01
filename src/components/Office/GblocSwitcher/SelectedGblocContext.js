import { createContext } from 'react';

export const SELECTED_GBLOC_SESSION_STORAGE_KEY = 'selected_gbloc';

const SelectedGblocContext = createContext(undefined);
SelectedGblocContext.displayName = 'SelectedGblocContext';

export default SelectedGblocContext;

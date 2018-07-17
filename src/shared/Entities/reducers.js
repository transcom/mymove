import { combineReducers } from 'redux';

import moveDocuments, {
  STATE_KEY as MOVEDOCUMENTS_STATE_KEY,
} from './modules/moveDocuments';
import documentModel, {
  STATE_KEY as DOCUMENTS_STATE_KEY,
} from './modules/documents';
import uploads, { STATE_KEY as UPLOADS_STATE_KEY } from './modules/uploads';

const reducer = combineReducers({
  [MOVEDOCUMENTS_STATE_KEY]: moveDocuments,
  [DOCUMENTS_STATE_KEY]: documentModel,
  [UPLOADS_STATE_KEY]: uploads,
});

export default reducer;

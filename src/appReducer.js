import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import { routerReducer } from 'react-router-redux';

import userReducer from 'shared/User/ducks';
import swaggerReducer from 'shared/Swagger/ducks';

import { feedbackReducer } from 'scenes/Feedback/ducks';
import { moveReducer } from 'scenes/Moves/ducks';
import { ppmReducer } from 'scenes/Moves/Ppm/ducks';
import issuesReducer from 'scenes/SubmittedFeedback/ducks';
import { shipmentsReducer } from 'scenes/Shipments/ducks';
import dd1299Reducer from 'scenes/DD1299/ducks';
import { signedCertificationReducer } from 'scenes/Legalese/ducks';
import { documentReducer } from 'shared/Uploader/ducks';

export const appReducer = combineReducers({
  user: userReducer,
  // swagger: swaggerReducer,
  submittedIssues: issuesReducer,
  submittedMoves: moveReducer,
  ppm: ppmReducer,
  shipments: shipmentsReducer,
  router: routerReducer,
  form: formReducer,
  feedback: feedbackReducer,
  signedCertification: signedCertificationReducer,
  DD1299: dd1299Reducer,
  upload: documentReducer,
});

export default appReducer;

import { connect } from 'react-redux';
import { get } from 'lodash';
import { SignedCertification } from './index';
import { loadCertificationText, signAndSubmitForApproval } from './ducks';
import { getFormValues } from 'redux-form';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';

const formName = 'signature-form';

function mapStateToProps(state) {
  return {
    schema: get(state, 'swaggerInternal.spec.definitions.CreateSignedCertificationPayload', {}),
    hasLoggedInUser: state.loggedInUser.hasSucceeded,
    values: getFormValues(formName)(state),
    ...state.signedCertification,
    has_sit: get(state.ppm, 'currentPpm.has_sit', false),
    has_advance: get(state.ppm, 'currentPpm.has_requested_advance', false),
    selectedMoveType: get(state.moves.currentMove, 'selected_move_type', null),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      loadCertificationText,
      signAndSubmitForApproval,
      push,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(SignedCertification);

import { connect } from 'react-redux';
import { get } from 'lodash';
import { SignedCertification } from './index';
import { getFormValues } from 'redux-form';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';
import { loadCertificationText, signAndSubmitPpm } from './ducks';
import { isHHGPPMComboMove } from 'scenes/Moves/Ppm/ducks';
import { selectGetCurrentUserIsSuccess } from 'shared/Data/users';

const formName = 'signature-form';

function mapStateToProps(state) {
  return {
    schema: get(state, 'swaggerInternal.spec.definitions.CreateSignedCertificationPayload', {}),
    hasLoggedInUser: selectGetCurrentUserIsSuccess(state),
    values: getFormValues(formName)(state),
    ...state.signedCertification,
    has_sit: get(state.ppm, 'currentPpm.has_sit', false),
    has_advance: get(state.ppm, 'currentPpm.has_requested_advance', false),
    selectedMoveType: get(state.moves.currentMove, 'selected_move_type', null),
    ppmId: get(state, 'ppm.currentPpm.id'),
    isHHGPPMComboMove: isHHGPPMComboMove(state),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      loadCertificationText,
      signAndSubmitForApproval: signAndSubmitPpm,
      push,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(SignedCertification);

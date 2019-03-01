import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { push } from 'react-router-redux';
import PropTypes from 'prop-types';
import { getFormValues } from 'redux-form';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { selectGetCurrentUserIsSuccess } from 'shared/Data/users';
import CertificationText from './CertificationText';
import Alert from 'shared/Alert';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { formatSwaggerDate } from 'shared/formatters';
import WizardHeader from 'scenes/Moves/WizardHeader';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import reviewGray from 'shared/icon/review-gray.svg';
import './index.css';

import { loadCertificationText, signAndSubmitForApproval } from './ducks';

const formName = 'signature-form';
const SignatureWizardForm = reduxifyWizardForm(formName);

export class SignedCertification extends Component {
  componentDidMount() {
    const { hasLoggedInUser, certificationText, has_advance, has_sit, selectedMoveType } = this.props;
    if (hasLoggedInUser && !certificationText) {
      this.props.loadCertificationText(has_sit, has_advance, selectedMoveType);
      return;
    }
  }

  handleSubmit = () => {
    const pendingValues = this.props.values;
    const { latestSignedCertification } = this.props;

    if (latestSignedCertification) {
      return this.props.push('/');
    }

    if (pendingValues) {
      const moveId = this.props.match.params.moveId;
      const { certificationText, ppmId } = this.props;

      return this.props
        .signAndSubmitForApproval(moveId, certificationText, pendingValues.signature, pendingValues.date, ppmId)
        .then(() => this.props.push('/'));
    }
  };
  print() {
    window.print();
  }
  render() {
    const { hasSubmitError, pages, pageKey, latestSignedCertification, isHHGPPMComboMove } = this.props;
    const today = formatSwaggerDate(new Date());
    const initialValues = {
      date: get(latestSignedCertification, 'date', today),
      signature: get(latestSignedCertification, 'signature', null),
    };
    return (
      <div>
        {isHHGPPMComboMove && (
          <WizardHeader
            icon={reviewGray}
            title="Review"
            right={
              <ProgressTimeline>
                <ProgressTimelineStep name="Move Setup" completed />
                <ProgressTimelineStep name="Review" current />
              </ProgressTimeline>
            }
          />
        )}
        <div className="legalese">
          {this.props.certificationText && (
            <SignatureWizardForm
              handleSubmit={this.handleSubmit}
              className={formName}
              pageList={pages}
              pageKey={pageKey}
              initialValues={initialValues}
              discardOnBack
            >
              <div className="usa-width-one-whole">
                <div>
                  <h2>Now for the official part...</h2>
                  <span className="box_top">
                    <p className="instructions">
                      Before officially booking your move, please carefully read and then sign the following.
                    </p>
                    <a className="pdf" onClick={this.print}>
                      Print
                    </a>
                  </span>

                  <CertificationText certificationText={this.props.certificationText} />

                  <div className="signature-box">
                    <h3>SIGNATURE</h3>
                    <p>
                      In consideration of said household goods or mobile homes being shipped at Government expense,{' '}
                      <strong>I hereby agree to the certifications stated above.</strong>
                    </p>
                    <div className="signature-fields">
                      <SwaggerField
                        className="signature"
                        fieldName="signature"
                        swagger={this.props.schema}
                        required
                        disabled={!!initialValues.signature}
                      />
                      <SwaggerField
                        className="signature-date"
                        fieldName="date"
                        swagger={this.props.schema}
                        required
                        disabled
                      />
                    </div>
                  </div>

                  {hasSubmitError && (
                    <Alert type="error" heading="Server Error">
                      There was a problem saving your signature. Please reload the page.
                    </Alert>
                  )}
                </div>
              </div>
            </SignatureWizardForm>
          )}
        </div>
      </div>
    );
  }
}

SignedCertification.propTypes = {
  signAndSubmitForApproval: PropTypes.func.isRequired,
  match: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
  ppmId: PropTypes.string,
};

function mapStateToProps(state) {
  return {
    schema: get(state, 'swaggerInternal.spec.definitions.CreateSignedCertificationPayload', {}),
    hasLoggedInUser: selectGetCurrentUserIsSuccess(state),
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

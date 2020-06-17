import { get } from 'lodash';
import moment from 'moment';
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
import './index.css';
import { createSignedCertification } from 'shared/Entities/modules/signed_certifications';
import { selectActivePPMForMove, loadPPMs } from 'shared/Entities/modules/ppms';
import { submitMoveForApproval } from 'shared/Entities/modules/moves';
import { ppmStandardLiability, storageLiability, ppmAdvance, additionalInformation } from './legaleseText';
import { showSubmitSuccessBanner, removeSubmitSuccessBanner } from './ducks';

const formName = 'signature-form';
const SignatureWizardForm = reduxifyWizardForm(formName);

export class SignedCertification extends Component {
  state = {
    hasMoveSubmitError: false,
  };

  componentDidMount() {
    this.props.loadPPMs(this.props.moveId);
  }

  getCertificationText(hasSit, hasRequestedAdvance) {
    const txt = [ppmStandardLiability];
    if (hasSit) txt.push(storageLiability);
    if (hasRequestedAdvance) txt.push(ppmAdvance);
    txt.push(additionalInformation);
    return txt.join('');
  }

  submitCertificate = () => {
    const signatureTime = moment().format();
    const { currentPpm, moveId } = this.props;
    const certificate = {
      certification_text: this.getCertificationText(currentPpm.has_sit, currentPpm.has_requested_advance),
      date: signatureTime,
      signature: 'CHECKBOX',
      personally_procured_move_id: currentPpm.id,
      certification_type: 'PPM_PAYMENT',
    };
    return this.props.createSignedCertification(moveId, certificate);
  };

  handleSubmit = () => {
    const pendingValues = this.props.values;
    const { latestSignedCertification } = this.props;
    const submitDate = moment().format();
    if (latestSignedCertification) {
      return this.props.push('/');
    }

    if (pendingValues) {
      const moveId = this.props.match.params.moveId;
      Promise.all([this.submitCertificate(), this.props.submitMoveForApproval(moveId, submitDate)])
        .then(() => {
          this.props.showSubmitSuccessBanner();
          setTimeout(() => this.props.removeSubmitSuccessBanner(), 10000);
          this.props.push('/');
        })
        .catch(() => this.setState({ hasMoveSubmitError: true }));
    }
  };

  print() {
    window.print();
  }

  render() {
    const { hasSubmitError, pages, pageKey, latestSignedCertification, currentPpm } = this.props;
    const today = formatSwaggerDate(new Date());
    const initialValues = {
      date: get(latestSignedCertification, 'date', today),
      signature: get(latestSignedCertification, 'signature', null),
    };
    const certificationText = this.getCertificationText(currentPpm.has_sit, currentPpm.has_requested_advance);
    return (
      <div>
        <div className="legalese">
          {certificationText && (
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
                    <a className="usa-link pdf" onClick={this.print}>
                      Print
                    </a>
                  </span>

                  <CertificationText certificationText={certificationText} />

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

                  {(hasSubmitError || this.state.hasMoveSubmitError) && (
                    <Alert type="error" heading="Server Error">
                      There was a problem saving your signature.
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
  match: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
  ppmId: PropTypes.string,
};

function mapStateToProps(state, ownProps) {
  const { moveId } = ownProps.match.params;
  return {
    moveId: moveId,
    schema: get(state, 'swaggerInternal.spec.definitions.CreateSignedCertificationPayload', {}),
    hasLoggedInUser: selectGetCurrentUserIsSuccess(state),
    values: getFormValues(formName)(state),
    ...state.signedCertification,
    currentPpm: selectActivePPMForMove(state, moveId),
    tempPpmId: get(state.ppm, 'currentPpm.id', null),
    has_sit: get(state.ppm, 'currentPpm.has_sit', false),
    has_advance: get(state.ppm, 'currentPpm.has_requested_advance', false),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      createSignedCertification,
      loadPPMs,
      submitMoveForApproval,
      showSubmitSuccessBanner,
      removeSubmitSuccessBanner,
      push,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(SignedCertification);

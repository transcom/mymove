import { get } from 'lodash';
import moment from 'moment';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import PropTypes from 'prop-types';
import { getFormValues } from 'redux-form';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { selectGetCurrentUserIsSuccess } from 'shared/Data/users';
import CertificationText from './CertificationText';
import Alert from 'shared/Alert';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { formatSwaggerDate } from 'shared/formatters';
import './index.scss';
import { createSignedCertification } from 'shared/Entities/modules/signed_certifications';
import { SIGNED_CERT_OPTIONS } from 'shared/constants';
import { getPPMsForMove, submitMoveForApproval } from 'services/internalApi';
import { updatePPMs, updateMove } from 'store/entities/actions';
import { completeCertificationText } from './legaleseText';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { setFlashMessage } from 'store/flash/actions';
import { selectCurrentPPM } from 'store/entities/selectors';

const formName = 'signature-form';
const SignatureWizardForm = reduxifyWizardForm(formName);

export class SignedCertification extends Component {
  state = {
    hasMoveSubmitError: false,
  };

  componentDidMount() {
    getPPMsForMove(this.props.moveId).then((response) => this.props.updatePPMs(response));
  }

  handleSubmit = () => {
    const { currentPpm, moveId, values } = this.props;
    const landingPath = '/';
    const submitDate = moment().format();
    const certificate = {
      certification_text: completeCertificationText,
      date: submitDate,
      signature: values.signature,
      personally_procured_move_id: currentPpm.id,
      certification_type: SIGNED_CERT_OPTIONS.SHIPMENT,
    };

    if (values) {
      submitMoveForApproval(moveId, certificate)
        .then((response) => {
          // Update Redux with new data
          this.props.updateMove(response);
          this.props.setFlashMessage('MOVE_SUBMIT_SUCCESS', 'success', 'You’ve submitted your move request.');
          this.props.push(landingPath);
        })
        .catch(() => this.setState({ hasMoveSubmitError: true }));
    }
  };

  print() {
    window.print();
  }

  render() {
    const { hasSubmitError, pages, pageKey } = this.props;
    const today = formatSwaggerDate(new Date());
    const initialValues = {
      date: today,
      signature: null,
    };
    const certificationText = completeCertificationText;
    const instructionsText = (
      <>
        <p>
          Please read this agreement, type your name in the <strong>Signature</strong> field to sign it, then tap the{' '}
          <strong>Complete</strong> button.
        </p>
        <p>This agreement covers the shipment of your personal property.</p>
      </>
    );
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
                  <h1>Now for the official part...</h1>
                  <div className="instructions">{instructionsText}</div>
                  <SectionWrapper>
                    <span className="box_top">
                      <a className="usa-link pdf" onClick={this.print}>
                        Print
                      </a>
                    </span>

                    <CertificationText certificationText={completeCertificationText} />

                    <div className="signature-box">
                      <h3>SIGNATURE</h3>
                      <p>
                        In consideration of said household goods or mobile homes being shipped at Government expense, I
                        hereby agree to the certifications stated above.
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

                      {(hasSubmitError || this.state.hasMoveSubmitError) && (
                        <Alert type="error" heading="Server Error">
                          There was a problem saving your signature.
                        </Alert>
                      )}
                    </div>
                  </SectionWrapper>
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
  setFlashMessage: PropTypes.func,
};

function mapStateToProps(state, ownProps) {
  const { moveId } = ownProps.match.params;
  return {
    moveId: moveId,
    schema: get(state, 'swaggerInternal.spec.definitions.CreateSignedCertificationPayload', {}),
    hasLoggedInUser: selectGetCurrentUserIsSuccess(state),
    values: getFormValues(formName)(state),
    currentPpm: selectCurrentPPM(state) || {},
  };
}

const mapDispatchToProps = {
  createSignedCertification,
  updatePPMs,
  updateMove,
  push,
  setFlashMessage,
};

export default connect(mapStateToProps, mapDispatchToProps)(SignedCertification);

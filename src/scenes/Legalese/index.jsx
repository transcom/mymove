import { get } from 'lodash';
import moment from 'moment';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
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
import { selectActivePPMForMove, loadPPMs } from 'shared/Entities/modules/ppms';
import { submitMoveForApproval } from 'shared/Entities/modules/moves';
import { completeCertificationText } from './legaleseText';
import { showSubmitSuccessBanner, removeSubmitSuccessBanner } from './ducks';
import SectionWrapper from 'components/Customer/SectionWrapper';

const formName = 'signature-form';
const SignatureWizardForm = reduxifyWizardForm(formName);

export class SignedCertification extends Component {
  state = {
    hasMoveSubmitError: false,
  };

  componentDidMount() {
    this.props.loadPPMs(this.props.moveId);
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
      this.props
        .submitMoveForApproval(moveId, certificate)
        .then(() => {
          this.props.showSubmitSuccessBanner();
          setTimeout(() => this.props.removeSubmitSuccessBanner(), 10000);
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
                  <p className="instructions">{instructionsText}</p>
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
};

function mapStateToProps(state, ownProps) {
  const { moveId } = ownProps.match.params;
  return {
    moveId: moveId,
    schema: get(state, 'swaggerInternal.spec.definitions.CreateSignedCertificationPayload', {}),
    hasLoggedInUser: selectGetCurrentUserIsSuccess(state),
    values: getFormValues(formName)(state),
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

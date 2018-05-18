import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { push } from 'react-router-redux';
import PropTypes from 'prop-types';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import CertificationText from './CertificationText';
import Alert from 'shared/Alert';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import './index.css';

import {
  loadCertificationText,
  loadLatestCertification,
  signAndSubmitForApproval,
} from './ducks';

const validateSignatureForm = (values, form) => {
  // all the validation is taken care of by the schema fields
  return {};
};

const formName = 'signature-form';
const SignatureWizardForm = reduxifyWizardForm(formName, validateSignatureForm);

export class SignedCertification extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Submit SignedCertification';
    this.props.loadLatestCertification(this.props.match.params.moveId);
  }

  componentDidUpdate() {
    const { getCertificationSuccess, certificationText } = this.props;
    if (getCertificationSuccess && !certificationText) {
      this.props.loadCertificationText();
      return;
    }
  }

  handleSubmit = () => {
    const pendingValues = this.props.formData.values;
    const { latestSignedCertification } = this.props;

    if (latestSignedCertification) {
      return this.props.push('/');
    }

    if (pendingValues) {
      const moveId = this.props.match.params.moveId;

      this.props
        .signAndSubmitForApproval(
          moveId,
          this.props.certificationText,
          pendingValues.signature,
          pendingValues.date,
        )
        .then(() => this.props.push('/'));
    }
  };

  render() {
    const {
      hasSubmitError,
      pages,
      pageKey,
      latestSignedCertification,
    } = this.props;
    const today = new Date(Date.now()).toISOString().split('T')[0];
    const initialValues = {
      date: get(latestSignedCertification, 'date', today),
      signature: get(latestSignedCertification, 'signature', null),
    };
    return (
      <div>
        {this.props.certificationText && (
          <SignatureWizardForm
            handleSubmit={this.handleSubmit}
            className={formName}
            pageList={pages}
            pageKey={pageKey}
            hasSucceeded={false}
            initialValues={initialValues}
          >
            <div className="usa-grid">
              <h2>Now for the official part...</h2>
              <span className="box_top">
                <p className="instructions">
                  Before officially booking your move, please carefully read and
                  then sign the following.
                </p>
                <a className="pdf Todo">Printer Friendly PDF</a>
              </span>

              <CertificationText
                certificationText={this.props.certificationText}
              />

              <div className="signature-box">
                <h3>SIGNATURE</h3>
                <p>
                  In consideration of said household goods or mobile homes being
                  shipped at Government expense,{' '}
                  <strong>
                    I hereby agree to the certifications stated above.
                  </strong>
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
                  There was a problem saving your signature. Please reload the
                  page.
                </Alert>
              )}
            </div>
          </SignatureWizardForm>
        )}
      </div>
    );
  }
}

SignedCertification.propTypes = {
  signAndSubmitForApproval: PropTypes.func.isRequired,
  loadLatestCertification: PropTypes.func.isRequired,
  match: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return {
    schema: get(
      state,
      'swagger.spec.definitions.CreateSignedCertificationPayload',
      {},
    ),
    formData: state.form[formName],
    ...state.signedCertification,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      loadCertificationText,
      loadLatestCertification,
      signAndSubmitForApproval,
      push,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(
  SignedCertification,
);

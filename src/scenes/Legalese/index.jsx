import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { push } from 'react-router-redux';
import PropTypes from 'prop-types';
import { getFormValues } from 'redux-form';
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

const formName = 'signature-form';
const SignatureWizardForm = reduxifyWizardForm(formName);

export class SignedCertification extends Component {
  componentDidMount() {
    this.props.loadLatestCertification(this.props.match.params.moveId);
  }

  componentDidUpdate() {
    const {
      getCertificationSuccess,
      hasLoggedInUser,
      certificationText,
      has_advance,
      has_sit,
    } = this.props;
    if (hasLoggedInUser && getCertificationSuccess && !certificationText) {
      this.props.loadCertificationText(has_sit, has_advance);
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
  print() {
    window.print();
  }
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
      <div className="legalese">
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
                <a className="pdf" onClick={this.print}>
                  Print
                </a>
              </span>

              <CertificationText
                certificationText={this.props.certificationText}
              />

              <div className="signature-box">
                <h3>SIGNATURE</h3>
                <p>
                  I agree that I have read and understand the above
                  notifications.
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
    hasLoggedInUser: state.loggedInUser.hasSucceeded,
    values: getFormValues(formName)(state),
    ...state.signedCertification,
    has_sit: get(state.ppm, 'currentPpm.has_sit', false),
    has_advance: get(state.ppm, 'currentPpm.has_requested_advance', false),
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

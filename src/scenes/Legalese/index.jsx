import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import CertificationText from './CertificationText';
import SignatureForm from './SignatureForm';
import Alert from 'shared/Alert';
import './index.css';

import { loadCertificationText, createSignedCertification } from './ducks';

export class SignedCertification extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Submit SignedCertification';
    this.props.loadCertificationText();
  }

  handleSubmit = values => {
    const certRequest = {
      moveId: this.props.match.params.moveId,
      createSignedCertificationPayload: {
        certification_text: this.props.certificationText,
        signature: values.signature,
        date: values.date,
      },
    };
    this.props.createSignedCertification(certRequest);
  };

  render() {
    const { hasSubmitError } = this.props;
    const today = new Date(Date.now()).toISOString().split('T')[0];
    const initialValues = {
      date: today,
    };
    return (
      <div className="usa-grid">
        <h2>Now for the official part...</h2>
        <span className="box_top">
          <p className="instructions">
            Before officially booking your move, please carefully read and then
            sign the following.
          </p>
          <a className="pdf">Printer Friendly PDF</a>
        </span>

        <CertificationText certificationText={this.props.certificationText} />
        <SignatureForm
          onSubmit={this.handleSubmit}
          initialValues={initialValues}
        />

        {hasSubmitError && (
          <Alert type="error" heading="Server Error">
            There was a problem saving your signature. Please reload the page.
          </Alert>
        )}
      </div>
    );
  }
}

SignedCertification.propTypes = {
  createSignedCertification: PropTypes.func.isRequired,
  loadCertificationText: PropTypes.func.isRequired,
  match: PropTypes.object.isRequired,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return state.signedCertification;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { loadCertificationText, createSignedCertification },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(
  SignedCertification,
);

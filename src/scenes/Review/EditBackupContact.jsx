import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { push } from 'react-router-redux';
import { reduxForm } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import {
  updateBackupContact,
  indexBackupContacts,
} from 'scenes/ServiceMembers/ducks';

import './Review.css';
import profileImage from './images/profile.png';

const editBackupContactFormName = 'edit_backup_contact';

let EditBackupContactForm = props => {
  const { onCancel, schema, handleSubmit, submitting, valid } = props;
  return (
    <form onSubmit={handleSubmit}>
      <img src={profileImage} alt="" /> Backup Contact
      <hr />
      <h3 className="sm-heading">Edit Backup Contact:</h3>
      <SwaggerField fieldName="name" swagger={schema} required />
      <SwaggerField fieldName="email" swagger={schema} required />
      <SwaggerField fieldName="telephone" swagger={schema} />
      <button type="submit" disabled={submitting || !valid}>
        Save
      </button>
      <button
        type="button"
        className="usa-button-secondary"
        disabled={submitting}
        onClick={onCancel}
      >
        Cancel
      </button>
    </form>
  );
};
EditBackupContactForm = reduxForm({
  form: editBackupContactFormName,
})(EditBackupContactForm);

class EditBackupContact extends Component {
  returnToReview = () => {
    const reviewAddress = `/moves/${this.props.match.params.moveId}/review`;
    this.props.push(reviewAddress);
  };

  componentDidUpdate = prevProps => {
    // Once service member loads, load the backup contact.
    if (this.props.serviceMember && !prevProps.serviceMember) {
      this.props.indexBackupContacts(this.props.serviceMember.id);
    }
  };

  updateContact = fieldValues => {
    if (fieldValues.telephone === '') {
      fieldValues.telephone = null;
    }
    return this.props
      .updateBackupContact(fieldValues.id, fieldValues)
      .then(() => {
        // This promise resolves regardless of error.
        if (!this.props.hasSubmitError) {
          this.returnToReview();
        } else {
          window.scrollTo(0, 0);
        }
      });
  };

  render() {
    const { error, schema, backupContacts } = this.props;

    let backupContact = null;
    if (backupContacts) {
      backupContact = backupContacts[0];
    }

    return (
      <div className="usa-grid">
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">
          <EditBackupContactForm
            initialValues={backupContact}
            onSubmit={this.updateContact}
            onCancel={this.returnToReview}
            schema={schema}
          />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    backupContacts: get(state, 'serviceMember.currentBackupContacts'),
    serviceMember: get(state, 'loggedInUser.loggedInUser.service_member'),
    move: get(state, 'moves.currentMove'),
    error: get(state, 'serviceMember.error'),
    hasSubmitError: get(state, 'serviceMember.updateBackupContactError'),
    schema: get(
      state,
      'swagger.spec.definitions.CreateServiceMemberBackupContactPayload',
      {},
    ),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { push, updateBackupContact, indexBackupContacts },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(EditBackupContact);
